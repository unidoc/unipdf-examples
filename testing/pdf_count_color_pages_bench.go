/*
 * Detect the number of pages and the color pages (1-offset) all pages in a list of PDF files.
 * Compares these results to running Ghostscript on the PDF files and reports an error if the results don't match.
 *
 * Run as: ./pdf_count_color_pages_bench [-o processDir] [-d] [-a] testdata/*.pdf > blah
 *
 * The main results are written to stderr so you will see them in your console.
 * Detailed information is written to stdout and you will see them in blah.
 *
 *  See the other command line options in the top of main()
 *      -o processDir - Temporary processing directory (default compare.pdfs)
 *      -d: Debug level logging
 *      -a: Keep converting PDF files after failures
 *      -min <val>: Minimum PDF file size to test
 *      -max <val>: Maximum PDF file size to test
 *      -r <name>: Name of results file
 */

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"image"
	"io"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"time"

	common "github.com/unidoc/unipdf/v3/common"
	pdfcontent "github.com/unidoc/unipdf/v3/contentstream"
	pdfcore "github.com/unidoc/unipdf/v3/core"
	pdf "github.com/unidoc/unipdf/v3/model"
)

func initUniDoc(debug bool) error {
	logLevel := common.LogLevelInfo
	if debug {
		logLevel = common.LogLevelDebug
	}
	common.SetLogger(common.ConsoleLogger{LogLevel: logLevel})

	return nil
}

const usage = `Usage:
pdf_count_color_pages_bench [-o <processDir>] [-d][-a][-min <val>][-max <val>] <file1> <file2> ...
-d: Debug level logging
-a: Keep converting PDF files after failures
-min <val>: Minimum PDF file size to test
-max <val>: Maximum PDF file size to test
-r <name>: Name of results file
`

func main() {
	debug := false            // Write debug level info to stdout?
	strict := true            // panic immediately a page color detection error occurs"
	compareGrayscale := false // Do PDF raster comparison on grayscale rasters?
	runAllTests := false      // Don't stop when a PDF file fails to process?
	var minSize int64 = -1    // Minimum size for an input PDF to be processed.
	var maxSize int64 = -1    // Maximum size for an input PDF to be processed.
	var results string        // Results file
	var outputDir string

	flag.StringVar(&outputDir, "o", "compare.pdfs", "Set output dir for ghostscript")
	flag.BoolVar(&debug, "d", false, "Enable debug logging")
	flag.BoolVar(&strict, "s", false, "Enable strict checking")
	flag.BoolVar(&compareGrayscale, "g", false, "Do PDF raster comparison on grayscale rasters")
	flag.BoolVar(&runAllTests, "a", false, "Run all tests. Don't stop at first failure")
	flag.Int64Var(&minSize, "min", -1, "Minimum size of files to process (bytes)")
	flag.Int64Var(&maxSize, "max", -1, "Maximum size of files to process (bytes)")
	flag.StringVar(&results, "r", "", "Results file")

	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}

	err := initUniDoc(debug)
	if err != nil {
		os.Exit(1)
	}
	compDir := makeUniqueDir(outputDir)
	fmt.Fprintf(os.Stderr, "compDir=%#q\n", compDir)
	defer removeDir(compDir)

	writers := []io.Writer{os.Stderr}
	if len(results) > 0 {
		f, err := os.OpenFile(results, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		writers = append(writers, f)
	}

	pdfList, err := patternsToPaths(args)
	if err != nil {
		common.Log.Error("patternsToPaths failed. args=%#q err=%v", args, err)
		os.Exit(1)
	}
	pdfList = sortFiles(pdfList, minSize, maxSize)
	passFiles := []string{}
	badFiles := []string{}
	failFiles := []string{}

	for idx, inputPath := range pdfList {

		colorPagesIn, err := pdfColorPages(inputPath, compDir)
		if err != nil {
			common.Log.Error("PDF is damaged. err=%v\n\tinputPath=%#q", err, inputPath)
			continue
		}
		var strictColorPages []int = nil
		if strict {
			strictColorPages = colorPagesIn
		}

		_, name := filepath.Split(inputPath)
		inputSize := fileSize(inputPath)
		report(writers, "%3d of %d %#-30q  (%6d)", idx, len(pdfList), name, inputSize)

		result := "pass"
		t0 := time.Now()

		numPages, colorPages, err := describePdf(inputPath, strictColorPages)
		dt := time.Since(t0)
		if err != nil {
			common.Log.Error("describePdf failed. err=%v", err)
			result = "bad"
		}
		report(writers, " %d pages %d color %.3f sec", numPages, len(colorPages), dt.Seconds())

		if result == "pass" {
			if !equalSlices(colorPagesIn, colorPages) {
				common.Log.Error("pdfColorPages: \ncolorPagesIn=%d %v\ncolorPages  =%d %v",
					len(colorPagesIn), colorPagesIn, len(colorPages), colorPages)
				fp := sliceDiff(colorPages, colorPagesIn)
				fn := sliceDiff(colorPagesIn, colorPages)
				if len(fp) > 0 {
					common.Log.Error("False positives=%d %+v", len(fp), fp)
				}
				if len(fn) > 0 {
					common.Log.Error("False negatives=%d %+v", len(fn), fn)
				}
				result = "fail"
			}
		}
		report(writers, ", %s\n", result)

		switch result {
		case "pass":
			passFiles = append(passFiles, inputPath)
		case "fail":
			failFiles = append(failFiles, inputPath)
		case "bad":
			badFiles = append(badFiles, inputPath)
		}

		if result != "pass" {
			if runAllTests {
				continue
			}
			break
		}

	}

	report(writers, "%d files %d bad %d pass %d fail\n", len(pdfList), len(badFiles), len(passFiles), len(failFiles))
	report(writers, "%d bad\n", len(badFiles))
	for i, path := range badFiles {
		report(writers, "%3d %#q\n", i, path)
	}
	report(writers, "%d pass\n", len(passFiles))
	for i, path := range passFiles {
		report(writers, "%3d %#q\n", i, path)
	}
	report(writers, "%d fail\n", len(failFiles))
	for i, path := range failFiles {
		report(writers, "%3d %#q\n", i, path)
	}
}

// describePdf reads PDF `inputPath` and returns number of pages, slice of color page numbers (1-offset)
func describePdf(inputPath string, strictColorPages []int) (int, []int, error) {

	f, err := os.Open(inputPath)
	if err != nil {
		return 0, []int{}, err
	}
	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return 0, []int{}, err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return 0, []int{}, err
	}
	if isEncrypted {
		_, err = pdfReader.Decrypt([]byte(""))
		if err != nil {
			return 0, []int{}, err
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return numPages, []int{}, err
	}

	colorPages := []int{}

	for i := 0; i < numPages; i++ {
		pageNum := i + 1
		page := pdfReader.PageList[i]
		common.Log.Debug("^^^^page %d", pageNum)

		desc := fmt.Sprintf("%s:page%d", filepath.Base(inputPath), pageNum)
		colored, err := isPageColored(page, desc, false)
		if err != nil {
			return numPages, colorPages, err
		}
		if strictColorPages != nil {
			strictColored := contains(strictColorPages, pageNum)
			if colored != strictColored {
				common.Log.Error("$$$$$$# Mismatch: Redoing pageNum=%d strictColored=%t colored=%t",
					pageNum, strictColored, colored)
				colored, _ = isPageColored(page, desc, true)
				common.Log.Error("$$$$$$* Mismatch: Redid inputPath=%q pageNum=%d strictColored=%t colored=%t",
					inputPath, pageNum, strictColored, colored)
				panic("Done")
			}
		}
		if colored {
			colorPages = append(colorPages, pageNum)
		}

	}

	return numPages, colorPages, nil
}

// isPageColored returns true if `page` contains color. It also references
// XObject Images and Forms to _possibly_ record if they contain color
func isPageColored(page *pdf.PdfPage, desc string, debug bool) (bool, error) {
	// For each page, we go through the resources and look for the images.

	contents, err := page.GetAllContentStreams()
	if err != nil {
		common.Log.Error("GetAllContentStreams failed. err=%v", err)
		return false, err
	}

	if debug {
		fmt.Println("\n===============***================")
		fmt.Printf("%s\n", desc)
		fmt.Println("===============+++================")
		fmt.Printf("%s\n", contents)
		fmt.Println("==================================")
	}

	colored, err := isContentStreamColored(contents, page.Resources, debug)
	if debug {
		common.Log.Info("colored=%t err=%v", colored, err)
	}
	if err != nil {
		common.Log.Error("isContentStreamColored failed. err=%v", err)
		return false, err
	}
	return colored, nil
}

// isPatternCS returns true if `colorspace` represents a Pattern colorspace.
func isPatternCS(cs pdf.PdfColorspace) bool {
	_, isPattern := cs.(*pdf.PdfColorspaceSpecialPattern)
	return isPattern
}

// isContentStreamColored returns true if `contents` contains any color object
func isContentStreamColored(contents string, resources *pdf.PdfPageResources, debug bool) (bool, error) {
	cstreamParser := pdfcontent.NewContentStreamParser(contents)
	operations, err := cstreamParser.Parse()
	if err != nil {
		return false, err
	}

	colored := false                                    // Has a colored mark been detected in the stream?
	coloredPatterns := map[pdfcore.PdfObjectName]bool{} // List of already detected patterns. Re-use for subsequent detections.
	coloredShadings := map[string]bool{}                // List of already detected shadings. Re-use for subsequent detections.

	// The content stream processor keeps track of the graphics state and we can make our own handlers to process
	// certain commands using the AddHandler method. In this case, we hook up to color related operands, and for image
	// and form handling.
	processor := pdfcontent.NewContentStreamProcessor(*operations)
	// Add handlers for colorspace related functionality.
	processor.AddHandler(pdfcontent.HandlerConditionEnumAllOperands, "",
		func(op *pdfcontent.ContentStreamOperation, gs pdfcontent.GraphicsState,
			resources *pdf.PdfPageResources) error {
			if colored {
				return nil
			}
			operand := op.Operand
			switch operand {
			case "SC", "SCN": // Set stroking color.  Includes pattern colors.
				if isPatternCS(gs.ColorspaceStroking) {
					op := pdfcontent.ContentStreamOperation{}
					op.Operand = operand
					op.Params = []pdfcore.PdfObject{}

					patternColor, ok := gs.ColorStroking.(*pdf.PdfColorPattern)
					if !ok {
						return errors.New("Invalid stroking color type")
					}

					if patternColor.Color != nil {
						if isColorColored(patternColor.Color) {
							if debug {
								common.Log.Info("op=%s hasCol=%t", op, true)
							}
							colored = true
							return nil
						}
					}

					if hasCol, ok := coloredPatterns[patternColor.PatternName]; ok {
						// Already processed, need not change anything, except underlying color if used.
						if hasCol {
							if debug {
								common.Log.Info("op=%s hasCol=%t", op, hasCol)
							}
							colored = true
						}
						return nil
					}

					// Look up the pattern name and convert it.
					pattern, found := resources.GetPatternByName(patternColor.PatternName)
					if !found {
						return errors.New("Undefined pattern name")
					}
					hasCol, err := isPatternColored(pattern, debug)
					if err != nil {
						common.Log.Error("isPatternColored failed. err=%v", err)
						return err
					}
					coloredPatterns[patternColor.PatternName] = hasCol
					colored = colored || hasCol
					if debug {
						common.Log.Info("op=%s hasCol=%t", op, hasCol)
					}

				} else {
					hasCol := isColorColored(gs.ColorStroking)
					colored = colored || hasCol
					if debug {
						common.Log.Info("op=%s ColorspaceStroking=%T ColorStroking=%#v hasCol=%t",
							op, gs.ColorspaceStroking, gs.ColorStroking, hasCol)
					}
				}
				return nil
			case "sc", "scn": // Set non-stroking color.
				if isPatternCS(gs.ColorspaceNonStroking) {
					op := pdfcontent.ContentStreamOperation{}
					op.Operand = operand
					op.Params = []pdfcore.PdfObject{}
					patternColor, ok := gs.ColorNonStroking.(*pdf.PdfColorPattern)
					if !ok {
						return errors.New("Invalid stroking color type")
					}
					if patternColor.Color != nil {
						hasCol := isColorColored(patternColor.Color)
						colored = colored || hasCol
						if debug {
							common.Log.Info("op=%#v hasCol=%t", op, hasCol)
						}
					}
					if hasCol, ok := coloredPatterns[patternColor.PatternName]; ok {
						// Already processed, need not change anything, except underlying color if used.
						colored = colored || hasCol
						if debug {
							common.Log.Info("op=%#v hasCol=%t", op, hasCol)
						}
						return nil
					}

					// Look up the pattern name and convert it.
					pattern, found := resources.GetPatternByName(patternColor.PatternName)
					if !found {
						return errors.New("Undefined pattern name")
					}
					hasCol, err := isPatternColored(pattern, debug)
					if err != nil {
						common.Log.Debug("Unable to convert pattern to grayscale: %v", err)
						return err
					}
					coloredPatterns[patternColor.PatternName] = hasCol
				} else {
					hasCol := isColorColored(gs.ColorNonStroking)
					colored = colored || hasCol
					if debug {
						common.Log.Info("op=%s ColorspaceNonStroking=%T ColorNonStroking=%#v hasCol=%t",
							op, gs.ColorspaceNonStroking, gs.ColorNonStroking, hasCol)
					}

				}
				return nil
			case "RG", "K": // Set RGB or CMYK stroking color.
				hasCol := isColorColored(gs.ColorStroking)
				if debug {
					common.Log.Info("op=%s ColorspaceStroking=%T ColorStroking=%#v hasCol=%t",
						op, gs.ColorspaceStroking, gs.ColorStroking, hasCol)
				}
				colored = colored || hasCol
				return nil
			case "rg", "k": // Set RGB or CMYK as non-stroking color.
				hasCol := isColorColored(gs.ColorNonStroking)
				colored = colored || hasCol
				if debug {
					common.Log.Info("op=%s ColorspaceStroking=%T ColorStroking=%#v hasCol=%t",
						op, gs.ColorspaceStroking, gs.ColorStroking, hasCol)
				}
				return nil
			case "sh": // Paints the shape and color defined by shading dict.
				if len(op.Params) != 1 {
					return errors.New("Params to sh operator should be 1")
				}
				shname, ok := op.Params[0].(*pdfcore.PdfObjectName)
				if !ok {
					return errors.New("sh parameter should be a name")
				}
				if hasCol, has := coloredShadings[string(*shname)]; has {
					// Already processed, no need to do anything.
					colored = colored || hasCol
					if debug {
						common.Log.Info("hasCol=%t", hasCol)
					}
					return nil
				}

				shading, found := resources.GetShadingByName(*shname)
				if !found {
					common.Log.Error("Shading not defined in resources. shname=%#q", string(*shname))
					return errors.New("Shading not defined in resources")
				}
				hasCol, err := isShadingColored(shading)
				if err != nil {
					return err
				}
				coloredShadings[string(*shname)] = hasCol
			}
			return nil
		})

	// Add handler for image related handling.  Note that inline images are completely stored with a ContentStreamInlineImage
	// object as the parameter for BI.
	processor.AddHandler(pdfcontent.HandlerConditionEnumOperand, "BI",
		func(op *pdfcontent.ContentStreamOperation, gs pdfcontent.GraphicsState, resources *pdf.PdfPageResources) error {
			if colored {
				return nil
			}
			if len(op.Params) != 1 {
				err := errors.New("invalid number of parameters")
				common.Log.Error("BI error. err=%v")
				return err
			}
			// Inline image.
			iimg, ok := op.Params[0].(*pdfcontent.ContentStreamInlineImage)
			if !ok {
				common.Log.Error("Invalid handling for inline image")
				return errors.New("Invalid inline image parameter")
			}
			if debug {
				common.Log.Info("iimg=%s", iimg)
			}
			img, err := iimg.ToImage(resources)
			if err != nil {
				common.Log.Error("Error converting inline image to image: %v", err)
				return err
			}

			if debug {
				common.Log.Info("img=%v %d", img.ColorComponents, img.BitsPerComponent)
			}

			if img.ColorComponents <= 1 {
				return nil
			}

			cs, err := iimg.GetColorSpace(resources)
			if err != nil {
				common.Log.Error("Error getting color space for inline image: %v", err)
				return err
			}
			rgbImg, err := cs.ImageToRGB(*img)
			if err != nil {
				common.Log.Error("Error converting image to rgb: %v", err)
				return err
			}
			hasCol := isRgbImageColored(rgbImg, debug)
			colored = colored || hasCol
			if debug {
				common.Log.Info("hasCol=%t", hasCol)
			}

			return nil
		})

	// Handler for XObject Image and Forms.
	processedXObjects := map[string]bool{} // Keep track of processed XObjects to avoid repetition.

	processor.AddHandler(pdfcontent.HandlerConditionEnumOperand, "Do",
		func(op *pdfcontent.ContentStreamOperation, gs pdfcontent.GraphicsState, resources *pdf.PdfPageResources) error {
			if colored {
				return nil
			}

			if len(op.Params) < 1 {
				common.Log.Error("Invalid number of params for Do object")
				return errors.New("Range check")
			}

			// XObject.
			name := op.Params[0].(*pdfcore.PdfObjectName)
			common.Log.Debug("Name=%#v=%#q", name, string(*name))

			// Only process each one once.
			hasCol, has := processedXObjects[string(*name)]
			common.Log.Debug("name=%q has=%t hasCol=%t processedXObjects=%+v", *name, has, hasCol, processedXObjects)
			if has {
				colored = colored || hasCol
				return nil
			}
			processedXObjects[string(*name)] = false

			_, xtype := resources.GetXObjectByName(*name)
			common.Log.Debug("xtype=%+v pdf.XObjectTypeImage=%v", xtype, pdf.XObjectTypeImage)

			if xtype == pdf.XObjectTypeImage {
				ximg, err := resources.GetXObjectImageByName(*name)
				if err != nil {
					common.Log.Error("Error w/GetXObjectImageByName : %v", err)
					return err
				}
				if debug {
					common.Log.Info("!!Filter=%s ColorSpace=%s ImageMask=%v wxd=%dx%d",
						ximg.Filter.GetFilterName(), ximg.ColorSpace,
						ximg.ImageMask, *ximg.Width, *ximg.Height)
				}
				if _, isIndexed := ximg.ColorSpace.(*pdf.PdfColorspaceSpecialIndexed); !isIndexed {
					if ximg.ColorSpace.GetNumComponents() == 1 {
						return nil
					}
				}
				switch ximg.Filter.GetFilterName() {
				// TODO: Add JPEG2000 encoding/decoding. Until then we assume JPEG200 images are color
				case "JPXDecode":
					processedXObjects[string(*name)] = true
					colored = true
					return nil
				// These filters are only used with grayscale images
				case "CCITTDecode", "JBIG2Decode", "RunLengthDecode":
					return nil
				}

				img, err := ximg.ToImage()
				if err != nil {
					common.Log.Error("Error w/ToImage: %v", err)
					return err
				}

				rgbImg, err := ximg.ColorSpace.ImageToRGB(*img)
				if err != nil {
					common.Log.Error("Error ImageToRGB: %v", err)
					return err
				}

				if debug {
					common.Log.Info("img: ColorComponents=%d wxh=%dx%d", img.ColorComponents, img.Width, img.Height)
					common.Log.Info("ximg: ColorSpace=%T=%s mask=%v", ximg.ColorSpace, ximg.ColorSpace, ximg.Mask)
					common.Log.Info("rgbImg: ColorComponents=%d wxh=%dx%d", rgbImg.ColorComponents, rgbImg.Width, rgbImg.Height)
				}

				hasCol := isRgbImageColored(rgbImg, debug)
				processedXObjects[string(*name)] = hasCol
				colored = colored || hasCol
				if debug {
					common.Log.Info("hasCol=%t", hasCol)
				}

			} else if xtype == pdf.XObjectTypeForm {
				common.Log.Debug(" XObject Form: %s", *name)

				// Go through the XObject Form content stream.
				xform, err := resources.GetXObjectFormByName(*name)
				if err != nil {
					fmt.Printf("Error : %v\n", err)
					return err
				}

				formContent, err := xform.GetContentStream()
				if err != nil {
					common.Log.Error("err=%v")
					return err
				}

				// Process the content stream in the Form object too:
				// XXX/TODO/Consider: Use either form resources (priority) and fall back to page resources alternatively if not found.
				// Have not come into cases where needed yet.
				formResources := xform.Resources
				if formResources == nil {
					formResources = resources
				}

				// Process the content stream in the Form object too:
				hasCol, err := isContentStreamColored(string(formContent), formResources, debug)
				if err != nil {
					common.Log.Error("err=%v", err)
					return err
				}
				processedXObjects[string(*name)] = hasCol
				colored = colored || hasCol
				if debug {
					common.Log.Info("hasCol=%t", hasCol)
				}

			}

			return nil
		})

	err = processor.Process(resources)
	if err != nil {
		common.Log.Error("processor.Process returned: err=%v", err)
		return false, err
	}

	return colored, nil
}

// isPatternColored returns true if `pattern` contains color (tiling or shading pattern).
func isPatternColored(pattern *pdf.PdfPattern, debug bool) (bool, error) {
	// Case 1: Colored tiling patterns.  Need to process the content stream and replace.
	if pattern.IsTiling() {
		tilingPattern := pattern.GetAsTilingPattern()
		if tilingPattern.IsColored() {
			// A colored tiling pattern can use color operators in its stream, need to process the stream.
			content, err := tilingPattern.GetContentStream()
			if err != nil {
				return false, err
			}
			colored, err := isContentStreamColored(string(content), tilingPattern.Resources, debug)
			return colored, err
		}
	} else if pattern.IsShading() {
		// Case 2: Shading patterns.  Need to create a new colorspace that can map from N=3,4 colorspaces to grayscale.
		shadingPattern := pattern.GetAsShadingPattern()
		colored, err := isShadingColored(shadingPattern.Shading)
		return colored, err
	}
	common.Log.Error("isPatternColored. pattern is neither tiling nor shading")
	return false, nil
}

// isShadingColored returns true if `shading` is a colored colorspace
func isShadingColored(shading *pdf.PdfShading) (bool, error) {
	cs := shading.ColorSpace
	if cs.GetNumComponents() == 1 {
		// Grayscale colorspace
		return false, nil
	} else if cs.GetNumComponents() == 3 {
		// RGB colorspace
		return true, nil
	} else if cs.GetNumComponents() == 4 {
		// CMYK colorspace
		return true, nil
	} else {
		err := errors.New("Unsupported pattern colorspace for color detection")
		common.Log.Error("isShadingColored: colorpace N=%d err=%v", cs.GetNumComponents(), err)
		return false, err
	}
}

// isColorColored returns true if `color` is not gray
func isColorColored(color pdf.PdfColor) bool {
	switch color.(type) {
	case *pdf.PdfColorDeviceGray:
		return false
	case *pdf.PdfColorDeviceRGB:
		col := color.(*pdf.PdfColorDeviceRGB)
		r, g, b := col.R(), col.G(), col.B()
		return visible(r-g, r-b, g-b)
	case *pdf.PdfColorDeviceCMYK:
		col := color.(*pdf.PdfColorDeviceCMYK)
		c, m, y := col.C(), col.M(), col.Y()
		return visible(c-m, c-y, m-y)
	case *pdf.PdfColorCalGray:
		return false
	case *pdf.PdfColorCalRGB:
		col := color.(*pdf.PdfColorCalRGB)
		a, b, c := col.A(), col.B(), col.C()
		return visible(a-b, a-c, b-c)
	case *pdf.PdfColorLab:
		col := color.(*pdf.PdfColorLab)
		a, b := col.A(), col.B()
		return visible(a, b)
	}
	common.Log.Error("isColorColored: Unknown color %T %s", color, color)
	panic("Unknown color type")
}

// isRgbImageColored returns true if `img` contains any color pixels
func isRgbImageColored(img pdf.Image, debug bool) bool {

	samples := img.GetSamples()
	maxVal := math.Pow(2, float64(img.BitsPerComponent)) - 1

	for i := 0; i < len(samples); i += 3 {
		// Normalized data, range 0-1.
		r := float64(samples[i]) / maxVal
		g := float64(samples[i+1]) / maxVal
		b := float64(samples[i+2]) / maxVal
		if visible(r-g, r-b, g-b) {
			if debug {
				common.Log.Info("@@ colored pixel: i=%d rgb=%.3f %.3f %.3f", i, r, g, b)
				common.Log.Info("                 delta rgb=%.3f %.3f %.3f", r-g, r-b, g-b)
				common.Log.Info("            colorTolerance=%.3f", colorTolerance)
			}
			return true
		}
	}
	return false
}

// ColorTolerance is the smallest color component that is visible on a typical mid-range color laser printer
// cpts have values in range 0.0-1.0
const colorTolerance = 3.1 / 255.0

// visible returns true if any of color component `cpts` is visible on a typical mid-range color laser printer
// cpts have values in range 0.0-1.0
func visible(cpts ...float64) bool {
	for _, x := range cpts {
		if math.Abs(x) > colorTolerance {
			return true
		}
	}
	return false
}

// report writes Sprintf formatted `format` ... to all writers in `writers`
func report(writers []io.Writer, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	for _, w := range writers {
		if _, err := io.WriteString(w, msg); err != nil {
			common.Log.Error("report: write to %#v failed msg=%s err=%v", w, msg, err)
		}
	}
}

// equalSlices returns true if `a` and `b` are identical
func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, x := range a {
		if x != b[i] {
			return false
		}
	}
	return true
}

// sortFiles returns the paths of the files in `pathList` sorted by ascending size.
// If minSize > 0 then only files of this size or larger are returned.
// If maxSize > 0 then only files of this size or smaller are returned.
func sortFiles(pathList []string, minSize, maxSize int64) []string {
	n := len(pathList)
	fdList := make([]FileData, n)
	for i, path := range pathList {
		fi, err := os.Stat(path)
		if err != nil {
			panic(err)
		}
		fdList[i].path = path
		fdList[i].FileInfo = fi
	}

	sort.Stable(byFile(fdList))

	i0 := 0
	i1 := n
	if minSize >= 0 {
		i0 = sort.Search(len(fdList), func(i int) bool { return fdList[i].Size() >= minSize })
	}
	if maxSize >= 0 {
		i1 = sort.Search(len(fdList), func(i int) bool { return fdList[i].Size() >= maxSize })
	}
	fdList = fdList[i0:i1]

	outList := make([]string, len(fdList))
	for i, fd := range fdList {
		outList[i] = fd.path
	}

	return outList
}

type FileData struct {
	path string
	os.FileInfo
}

// byFile sorts slices of FileData by some file attribute, currently size.
type byFile []FileData

func (x byFile) Len() int { return len(x) }

func (x byFile) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

func (x byFile) Less(i, j int) bool {
	si, sj := x[i].Size(), x[j].Size()
	if si != sj {
		return si < sj
	}
	return x[i].path < x[j].path
}

const (
	gsImageFormat  = "doc-%03d.png"
	gsImagePattern = `doc-(\d+).png$`
)

var gsImageRegex = regexp.MustCompile(gsImagePattern)

// runGhostscript runs Ghostscript on file `pdf` to create file one png file per page in directory
// `outputDir`
func runGhostscript(pdf, outputDir string) error {
	common.Log.Debug("runGhostscript: pdf=%#q outputDir=%#q", pdf, outputDir)
	outputPath := filepath.Join(outputDir, gsImageFormat)
	output := fmt.Sprintf("-sOutputFile=%s", outputPath)

	cmd := exec.Command(
		ghostscriptName(),
		"-dSAFER",
		"-dBATCH",
		"-dNOPAUSE",
		"-r150",
		fmt.Sprintf("-sDEVICE=png16m"),
		"-dTextAlphaBits=1",
		"-dGraphicsAlphaBits=1",
		output,
		pdf)
	common.Log.Debug("runGhostscript: cmd=%#q", cmd.Args)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		common.Log.Error("runGhostscript: Could not process pdf=%q err=%v\nstdout=\n%s\nstderr=\n%s\n",
			pdf, err, stdout, stderr)
	}
	return err
}

// ghostscriptName returns the name of the Ghostscript binary on this OS
func ghostscriptName() string {
	if runtime.GOOS == "windows" {
		return "gswin64c.exe"
	}
	return "gs"
}

// pdfColorPages returns a list of the (1-offset) page numbers of the colored pages in PDF at `path`
func pdfColorPages(path, temp string) ([]int, error) {
	dir := filepath.Join(temp, "color")
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		panic(err)
	}
	defer removeDir(dir)

	err = runGhostscript(path, dir)
	if err != nil {
		return nil, err
	}

	return colorDirectoryPages("*.png", dir)
}

// colorDirectoryPages returns a list of the (1-offset) page numbers of the image files that match `mask` in
// directories `dir` that have any color pixels.
func colorDirectoryPages(mask, dir string) ([]int, error) {
	pattern := filepath.Join(dir, mask)
	files, err := filepath.Glob(pattern)
	if err != nil {
		common.Log.Error("colorDirectoryPages: Glob failed. pattern=%#q err=%v", pattern, err)
		return nil, err
	}

	colorPages := []int{}
	for _, path := range files {
		matches := gsImageRegex.FindStringSubmatch(path)
		if len(matches) == 0 {
			continue
		}
		pageNum, err := strconv.Atoi(matches[1])
		if err != nil {
			panic(err)
			return colorPages, err
		}
		isColor, err := isColorImage(path)
		if err != nil {
			panic(err)
			return colorPages, err
		}
		if isColor {
			colorPages = append(colorPages, pageNum)
		}
	}
	return colorPages, nil
}

// isColorImage returns true if image file `path` contains color
func isColorImage(path string) (bool, error) {
	img, err := readImage(path)
	if err != nil {
		return false, err
	}
	return imgIsColor(img), nil
}

const F = 255.0 / float64(0xFFFF)

// colorThreshold is the min difference of r,g,b values for which a pixel is considered to be color
// Color components are in range 0-0xFFFF
// We make this 5 on a 0-255 scale as a guess
const colorThreshold = 5.0 // / F

// imgIsColor returns true if image `img` contains color
func imgIsColor(img image.Image) bool {
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			rr, gg, bb, _ := img.At(x, y).RGBA()
			r, g, b := float64(rr)*F, float64(gg)*F, float64(bb)*F
			if math.Abs(r-g) > colorThreshold || math.Abs(r-b) > colorThreshold || math.Abs(g-b) > colorThreshold {
				// fmt.Printf("@@@ xy=%d %d rgb=%.3f %.3f %.3f =%.3f %.3f %.3f\n", x, y, r, g, b, r-g, r-b, g-b)
				return true
			}
		}
	}
	return false
}

// readImage reads image file `path` and returns its contents as an Image.
func readImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		common.Log.Error("readImage: Could not open file. path=%#q err=%v", path, err)
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	return img, err
}

// makeUniqueDir creates a new directory inside `baseDir`
func makeUniqueDir(baseDir string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 1000; i++ {
		dir := filepath.Join(baseDir, fmt.Sprintf("dir.%03d", i))
		if _, err := os.Stat(dir); err != nil {
			if os.IsNotExist(err) {
				if err := os.MkdirAll(dir, 0777); err != nil {
					panic(err)
				}
				return dir
			}
		}
		time.Sleep(time.Duration(r.Float64() * float64(time.Second)))
	}
	panic("Cannot create new directory")
}

// removeDir removes directory `dir` and its contents
func removeDir(dir string) error {
	err1 := os.RemoveAll(dir)
	err2 := os.Remove(dir)
	if err1 != nil {
		return err1
	}
	return err2
}

// patternsToPaths returns a list of files matching the patterns in `patternList`
func patternsToPaths(patternList []string) ([]string, error) {
	pathList := []string{}
	for _, pattern := range patternList {
		files, err := filepath.Glob(pattern)
		if err != nil {
			common.Log.Error("patternsToPaths: Glob failed. pattern=%#q err=%v", pattern, err)
			return pathList, err
		}
		for _, path := range files {
			if !regularFile(path) {
				fmt.Fprintf(os.Stderr, "Not a regular file. %#q\n", path)
				continue
			}
			pathList = append(pathList, path)
		}
	}
	return pathList, nil
}

// regularFile returns true if file `path` is a regular file
func regularFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	return fi.Mode().IsRegular()
}

// fileSize returns the size of file `path` in bytes
func fileSize(path string) int64 {
	fi, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	return fi.Size()
}

// sliceDiff returns the elements in `a` that aren't in `b`
func sliceDiff(a, b []int) []int {
	mb := map[int]bool{}
	for _, x := range b {
		mb[x] = true
	}
	ab := []int{}
	for _, x := range a {
		if _, ok := mb[x]; !ok {
			ab = append(ab, x)
		}
	}
	return ab
}

// contains returns true if `s` contains `e`
func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
