/*
 * Transform all content streams in all pages in a list of pdf files.
 *
 * Run as: go run pdf_grayscale_convert_bench -g output [-d][-k][-a] testdata/*.pdf > blah
 *
 * This will transform all .pdf file in testdata and write the results to output.
 *
 * See the other command line options in the top of main()
 *      -o processDir - Temporary processing directory (default compare.pdfs)
 *      -d: Debug level logging
 *      -a: Keep converting PDF files after failures
 *      -min <val>: Minimum PDF file size to test
 *      -max <val>: Maximum PDF file size to test
 *      -r <name>: Name of results file
 *
 * The grayscale transform
 *	- converts PDF files into our internal representation
 *	- transforms the internal representation to grayscale
 *	- converts the internal representation back to a PDF file
 *	- checks that the output PDF file is grayscale
 *
 * Meanings:
 * fail - Failed for any reason
 * bad - Grayscale conversion process failed
 * pass+fail should equal total files.
 */

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
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
	"strings"
	"time"

	common "github.com/unidoc/unipdf/v3/common"
	pdfcontent "github.com/unidoc/unipdf/v3/contentstream"
	pdfcore "github.com/unidoc/unipdf/v3/core"
	pdf "github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/ps"
)

const usage = `Usage:
pdf_grayscale_convert_bench -g <output directory> [-r <results>][-d][-k][-a][-min <val>][-max <val>] <file1> <file2> ...
-g <outputDir> - Converted grayscale PDFs are written to outputDir
-r <resultsPath> - Results are written to resultsPath
-o <processDir> - Temporary processing directory (default compare.pdfs)
-d: Debug level logging
-a: Keep converting PDF files after failures
-min <val>: Minimum PDF file size to test
-max <val>: Maximum PDF file size to test
-k: Keep temp PNG files used for PDF grayscale test
`

// Ignore CCITTFaxDecode, JBIG2 - that are always grayscale.
// Change with -ignoregrayfilters=false
var ignoreGrayFilters = true

func initUniDoc(debug bool) {
	pdf.SetPdfCreator("pdf_grayscale_convert_bench test suite")

	logLevel := common.LogLevelInfo
	if debug {
		logLevel = common.LogLevelDebug
		//logLevel = common.LogLevelTrace
	}
	common.SetLogger(common.ConsoleLogger{LogLevel: logLevel})
}

// imageThreshold represents the threshold for image difference in image comparisons
type imageThreshold struct {
	fracPixels float64 // Fraction of pixels in page raster that may differ
	mean       float64 // Max mean difference on scale 0..255 for pixels that differ
}

func main() {
	debug := false       // Write debug level info to stdout?
	runAllTests := false // Don't stop when a PDF file fails to process?
	compRoot := ""
	var minSize int64 = -1   // Minimum size for an input PDF to be processed
	var maxSize int64 = -1   // Maximum size for an input PDF to be processed
	results := ""            // Results are written here
	outputDir := ""          // Transformed PDFs are written here
	keep := false            // Keep the rasters used for PDF comparison
	ignoreGrayFilters = true // Ignore CCITTFaxDecode, JBIG2 - that are always grayscale.

	flag.BoolVar(&debug, "d", false, "Enable debug logging")
	flag.BoolVar(&runAllTests, "a", false, "Run all tests. Don't stop at first failure")
	flag.StringVar(&compRoot, "o", "compare.pdfs", "Set output dir for ghostscript")
	flag.Int64Var(&minSize, "min", -1, "Minimum size of files to process (bytes)")
	flag.Int64Var(&maxSize, "max", -1, "Maximum size of files to process (bytes)")
	flag.StringVar(&outputDir, "g", "", "Output directory")
	flag.StringVar(&results, "r", "", "Results file")
	flag.BoolVar(&keep, "k", false, "Keep the rasters used for PDF comparison")
	flag.BoolVar(&ignoreGrayFilters, "ignoregrayfilters", true, "Ignore gray filters (CCITTFaxDecode, JPXDecode)")

	flag.Parse()
	args := flag.Args()
	if len(args) < 1 || len(outputDir) == 0 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}

	initUniDoc(debug)
	compDir := makeUniqueDir(compRoot)
	fmt.Printf("compDir=%#q\n", compDir)
	if !keep {
		defer removeDir(compDir)
	}
	defer removeDir(compDir)

	writers := []io.Writer{os.Stdout}
	if len(results) > 0 {
		f, err := os.OpenFile(results, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		writers = append(writers, f)
	}

	err := os.MkdirAll(outputDir, 0777)
	if err != nil {
		common.Log.Error("MkdirAll failed. outputDir=%#q err=%v", outputDir, err)
		os.Exit(1)
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
	failErrors := []string{}
	passTotalTime := float64(0)

	startT := time.Now()

	for idx, inputPath := range pdfList {

		_, name := filepath.Split(inputPath)
		inputSize := fileSize(inputPath)

		report(writers, "%3d of %d %#-30q  (%6d->", idx, len(pdfList), name, inputSize)
		outputPath := modifyPath(inputPath, outputDir)

		t0 := time.Now()
		result := "pass"

		// 1. Transforms the pdf to grayscale pdf.
		numPages, err := transformPdfFile(inputPath, outputPath)
		dt := time.Since(t0)
		if err != nil {
			common.Log.Error("transformPdfFile failed. err=%v", err)
			failFiles = append(failFiles, inputPath)
			failErrors = append(failErrors, fmt.Sprintf("%v", err))
			result = "bad"
		}

		// Reason for failure.
		errStr := ""

		// 2. Runs pdftops on the transformed file to validate if OK.
		if result == "pass" {
			outputSize := fileSize(outputPath)
			report(writers, "%6d %3d%%) %d pages %.3f sec => %#q",
				outputSize, int(float64(outputSize)/float64(inputSize)*100.0+0.5),
				numPages, dt.Seconds(), outputPath)

			err = runPdfToPs(outputPath, compDir)
			if err != nil {
				common.Log.Error("Transform has damaged PDF. err=%v\n\tinputPath=%#q\n\toutputPath=%#q",
					err, inputPath, outputPath)
				result = "fail"
				errStr = "Transform -> damaged PDF"
			}
		}

		// 3. Checks if the transformed PDF has color pixels.
		if result == "pass" {
			isColorOut, colorPagesOut, err := isPdfColor(outputPath, compDir, true, keep)

			if err != nil || isColorOut {
				if err != nil {
					common.Log.Error("Transform has damaged PDF. err=%v\n\tinputPath=%#q\n\toutputPath=%#q",
						err, inputPath, outputPath)
					errStr = fmt.Sprintf("isPdfColor check failed: %v", err)
				} else {
					common.Log.Error("isPdfColor: %d Color pages", len(colorPagesOut))
					errStr = fmt.Sprintf("color fail: %d color pages / %d total", len(colorPagesOut), numPages)
				}
				result = "fail"
			}
		}
		report(writers, ", %s\n", result)

		switch result {
		case "pass":
			passFiles = append(passFiles, inputPath)
			passTotalTime += dt.Seconds()
		case "fail":
			// Fail here means that it does not have color pixels.
			failFiles = append(failFiles, inputPath)
			failErrors = append(failErrors, errStr)
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

	totalDur := time.Since(startT)

	report(writers, "%d files %d bad %d pass %d fail\n", len(pdfList), len(badFiles), len(passFiles), len(failFiles))
	report(writers, "total duration (everything): %.0f seconds\n", totalDur.Seconds())
	if len(passFiles) > 0 {
		avgTime := passTotalTime / float64(len(passFiles))
		report(writers, "total processing time (pass only): %.0f seconds (%.2f sec per file)\n", passTotalTime, avgTime)
	}
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
		report(writers, "%3d %#q - %s\n", i, path, failErrors[i])
	}

	// Make a frequency count of error messages.
	errFreqMap := map[string]int{}
	for _, errstr := range failErrors {
		if _, has := errFreqMap[errstr]; has {
			errFreqMap[errstr]++
		} else {
			errFreqMap[errstr] = 1
		}
	}
	type errFreqEntry struct {
		errmsg string
		count  int
	}
	entries := []errFreqEntry{}
	for errstr, count := range errFreqMap {
		entries = append(entries, errFreqEntry{errstr, count})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].count < entries[j].count
	})
	report(writers, "Error frequencies:\n")
	for _, entry := range entries {
		report(writers, "%d - %s\n", entry.count, entry.errmsg)
	}
}

type ObjCounts struct {
	xobjNameSubtype map[string]string
}

// transformPdfFile transforms PDF `inputPath` and writes the resulting PDF to `outputPath`
// Returns: number of pages in `inputPath`
func transformPdfFile(inputPath, outputPath string) (int, error) {

	f, err := os.Open(inputPath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return 0, err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return 0, err
	}
	if isEncrypted {
		_, err = pdfReader.Decrypt([]byte(""))
		if err != nil {
			return 0, err
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return numPages, err
	}

	pdfWriter := pdf.NewPdfWriter()

	for i := 0; i < numPages; i++ {
		pageNum := i + 1
		page := pdfReader.PageList[i]
		common.Log.Debug("^^^^page %d", pageNum)

		desc := fmt.Sprintf("%s:page%d", filepath.Base(inputPath), pageNum)
		err = convertPageToGrayscale(page, desc)
		if err != nil {
			return numPages, err
		}

		err = pdfWriter.AddPage(page)
		if err != nil {
			return numPages, err
		}
		//  break !@#$ Single page mode
	}

	fWrite, err := os.Create(outputPath)
	if err != nil {
		return numPages, err
	}
	defer fWrite.Close()
	err = pdfWriter.Write(fWrite)

	return numPages, nil
}

// =================================================================================================
// Page transform code goes here
// =================================================================================================

// convertPageToGrayscale replaces color objects on the page with grayscale ones. It also references
// XObject Images and Forms to convert those to grayscale.
func convertPageToGrayscale(page *pdf.PdfPage, desc string) error {
	// For each page, we go through the resources and look for the images.
	contents, err := page.GetAllContentStreams()
	if err != nil {
		common.Log.Error("GetAllContentStreams failed. err=%v", err)
		return err
	}

	grayContent, err := transformContentStreamToGrayscale(contents, page.Resources)
	if err != nil {
		common.Log.Error("transformContentStreamToGrayscale failed. err=%v", err)
		return err
	}
	page.SetContentStreams([]string{string(grayContent)}, pdfcore.NewFlateEncoder())

	return nil
}

// isPatternCS returns true if `colorspace` represents a Pattern colorspace.
func isPatternCS(cs pdf.PdfColorspace) bool {
	_, isPattern := cs.(*pdf.PdfColorspaceSpecialPattern)
	return isPattern
}

// transformContentStreamToGrayscale `contents` converted to grayscale.
func transformContentStreamToGrayscale(contents string, resources *pdf.PdfPageResources) ([]byte, error) {
	cstreamParser := pdfcontent.NewContentStreamParser(contents)
	operations, err := cstreamParser.Parse()
	if err != nil {
		return nil, err
	}
	processedOperations := &pdfcontent.ContentStreamOperations{}

	transformedPatterns := map[pdfcore.PdfObjectName]bool{} // List of already transformed patterns. Avoid multiple conversions.
	transformedShadings := map[pdfcore.PdfObjectName]bool{} // List of already transformed shadings. Avoid multiple conversions.

	// Make a map of pattern colorspaces. Once processed, changes UnderlyingCS to DeviceGray.
	patternColorspaces := map[*pdf.PdfColorspaceSpecialPattern]bool{}
	defer func() {
		// At the very end, set the UnderlyingCS to DeviceGray.
		// TODO: Needs to be global?  I.e. do this at the very end of the execution.
		for patternCS, _ := range patternColorspaces {
			patternCS.UnderlyingCS = pdf.NewPdfColorspaceDeviceGray()
		}
	}()

	// The content stream processor keeps track of the graphics state and we can make our own handlers to process
	// certain commands using the AddHandler method. In this case, we hook up to color related operands, and for image
	// and form handling.
	processor := pdfcontent.NewContentStreamProcessor(*operations)
	// Add handlers for colorspace related functionality.
	processor.AddHandler(pdfcontent.HandlerConditionEnumAllOperands, "",
		func(op *pdfcontent.ContentStreamOperation, gs pdfcontent.GraphicsState, resources *pdf.PdfPageResources) error {
			operand := op.Operand
			switch operand {
			case "CS": // Set colorspace operands (stroking).
				if isPatternCS(gs.ColorspaceStroking) {
					// If referring to a pattern colorspace with an external definition, need to update the definition.
					// If has an underlying colorspace, then go and change it to DeviceGray.
					// Needs to be specified externally in the colorspace resources.

					csname := op.Params[0].(*pdfcore.PdfObjectName)
					if *csname != "Pattern" {
						// Update if referring to an external colorspace in resources.
						cs, ok := resources.GetColorspaceByName(*csname)
						if !ok {
							common.Log.Debug("Undefined colorspace for pattern (%s)", csname)
							return errors.New("Colorspace not defined")
						}

						patternCS, ok := cs.(*pdf.PdfColorspaceSpecialPattern)
						if !ok {
							return errors.New("Type error")
						}

						// Mark for processing.
						patternColorspaces[patternCS] = true
						resources.SetColorspaceByName(*csname, patternCS)
					}
					*processedOperations = append(*processedOperations, op)
					return nil
				}

				op := pdfcontent.ContentStreamOperation{}
				op.Operand = operand
				op.Params = []pdfcore.PdfObject{pdfcore.MakeName("DeviceGray")}
				*processedOperations = append(*processedOperations, &op)
				return nil
			case "cs": // Set colorspace operands (non-stroking).
				if isPatternCS(gs.ColorspaceNonStroking) {
					// If referring to a pattern colorspace with an external definition, need to update the definition.
					// If has an underlying colorspace, then go and change it to DeviceGray.
					// Needs to be specified externally in the colorspace resources.

					csname := op.Params[0].(*pdfcore.PdfObjectName)
					if *csname != "Pattern" {
						// Update if referring to an external colorspace in resources.
						cs, ok := resources.GetColorspaceByName(*csname)
						if !ok {
							common.Log.Debug("Undefined colorspace for pattern (%s)", csname)
							return errors.New("Colorspace not defined")
						}

						patternCS, ok := cs.(*pdf.PdfColorspaceSpecialPattern)
						if !ok {
							return errors.New("Type error")
						}

						// Mark for processing.
						patternColorspaces[patternCS] = true
						resources.SetColorspaceByName(*csname, patternCS)
					}
					*processedOperations = append(*processedOperations, op)
					return nil
				}

				op := pdfcontent.ContentStreamOperation{}
				op.Operand = operand
				op.Params = []pdfcore.PdfObject{pdfcore.MakeName("DeviceGray")}
				*processedOperations = append(*processedOperations, &op)
				return nil

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
						color, err := gs.ColorspaceStroking.ColorToRGB(patternColor.Color)
						if err != nil {
							common.Log.Error("err=%v", err)
							return err
						}
						rgbColor := color.(*pdf.PdfColorDeviceRGB)
						grayColor := rgbColor.ToGray()

						op.Params = append(op.Params, pdfcore.MakeFloat(grayColor.Val()))
					}

					if _, has := transformedPatterns[patternColor.PatternName]; has {
						// Already processed, need not change anything, except underlying color if used.
						op.Params = append(op.Params, &patternColor.PatternName)
						*processedOperations = append(*processedOperations, &op)
						return nil
					}
					transformedPatterns[patternColor.PatternName] = true

					// Look up the pattern name and convert it.
					pattern, found := resources.GetPatternByName(patternColor.PatternName)
					if !found {
						return errors.New("Undefined pattern name")
					}

					grayPattern, err := convertPatternToGray(pattern)
					if err != nil {
						common.Log.Debug("Unable to convert pattern to grayscale: %v", err)
						return err
					}
					resources.SetPatternByName(patternColor.PatternName, grayPattern.ToPdfObject())

					op.Params = append(op.Params, &patternColor.PatternName)
					*processedOperations = append(*processedOperations, &op)
				} else {
					color, err := gs.ColorspaceStroking.ColorToRGB(gs.ColorStroking)
					if err != nil {
						common.Log.Error("Error with ColorToRGB: %v", err)
						return err
					}
					rgbColor := color.(*pdf.PdfColorDeviceRGB)
					grayColor := rgbColor.ToGray()

					op := pdfcontent.ContentStreamOperation{}
					op.Operand = operand
					op.Params = []pdfcore.PdfObject{pdfcore.MakeFloat(grayColor.Val())}
					*processedOperations = append(*processedOperations, &op)
				}

				return nil
			case "sc", "scn": // Set nonstroking color.
				if isPatternCS(gs.ColorspaceNonStroking) {
					op := pdfcontent.ContentStreamOperation{}
					op.Operand = operand
					op.Params = []pdfcore.PdfObject{}
					patternColor, ok := gs.ColorNonStroking.(*pdf.PdfColorPattern)
					if !ok {
						return errors.New("Invalid stroking color type")
					}

					if patternColor.Color != nil {
						color, err := gs.ColorspaceNonStroking.ColorToRGB(patternColor) //.Color)
						if err != nil {
							common.Log.Error("err=%v", err)
							return err
						}
						rgbColor := color.(*pdf.PdfColorDeviceRGB)
						grayColor := rgbColor.ToGray()

						op.Params = append(op.Params, pdfcore.MakeFloat(grayColor.Val()))
					}

					if _, has := transformedPatterns[patternColor.PatternName]; has {
						// Already processed, need not change anything, except underlying color if used.
						op.Params = append(op.Params, &patternColor.PatternName)
						*processedOperations = append(*processedOperations, &op)
						return nil
					}
					transformedPatterns[patternColor.PatternName] = true

					// Look up the pattern name and convert it.
					pattern, found := resources.GetPatternByName(patternColor.PatternName)
					if !found {
						return errors.New("Undefined pattern name")
					}

					grayPattern, err := convertPatternToGray(pattern)
					if err != nil {
						common.Log.Debug("Unable to convert pattern to grayscale: %v", err)
						return err
					}
					resources.SetPatternByName(patternColor.PatternName, grayPattern.ToPdfObject())
					op.Params = append(op.Params, &patternColor.PatternName)
					*processedOperations = append(*processedOperations, &op)
				} else {
					color, err := gs.ColorspaceNonStroking.ColorToRGB(gs.ColorNonStroking)
					if err != nil {
						common.Log.Error("err=%v", err)
						return err
					}
					rgbColor := color.(*pdf.PdfColorDeviceRGB)
					grayColor := rgbColor.ToGray()

					op := pdfcontent.ContentStreamOperation{}
					op.Operand = operand
					op.Params = []pdfcore.PdfObject{pdfcore.MakeFloat(grayColor.Val())}

					*processedOperations = append(*processedOperations, &op)
				}
				return nil
			case "RG", "K": // Set RGB or CMYK stroking color.
				color, err := gs.ColorspaceStroking.ColorToRGB(gs.ColorStroking)
				if err != nil {
					common.Log.Error("err=%v", err)
					return err
				}
				rgbColor := color.(*pdf.PdfColorDeviceRGB)
				grayColor := rgbColor.ToGray()

				op := pdfcontent.ContentStreamOperation{}
				op.Operand = "G"
				op.Params = []pdfcore.PdfObject{pdfcore.MakeFloat(grayColor.Val())}

				*processedOperations = append(*processedOperations, &op)
				return nil
			case "rg", "k": // Set RGB or CMYK as nonstroking color.
				color, err := gs.ColorspaceNonStroking.ColorToRGB(gs.ColorNonStroking)
				if err != nil {
					common.Log.Error("err=%v", err)
					return err
				}
				rgbColor := color.(*pdf.PdfColorDeviceRGB)
				grayColor := rgbColor.ToGray()

				op := pdfcontent.ContentStreamOperation{}
				op.Operand = "g"
				op.Params = []pdfcore.PdfObject{pdfcore.MakeFloat(grayColor.Val())}

				*processedOperations = append(*processedOperations, &op)
				return nil
			case "sh": // Paints the shape and color defined by shading dict.
				if len(op.Params) != 1 {
					return errors.New("Params to sh operator should be 1")
				}
				shname, ok := op.Params[0].(*pdfcore.PdfObjectName)
				if !ok {
					return errors.New("sh parameter should be a name")
				}
				if _, has := transformedShadings[*shname]; has {
					// Already processed, no need to do anything.
					*processedOperations = append(*processedOperations, op)
					return nil
				}
				transformedShadings[*shname] = true

				shading, found := resources.GetShadingByName(*shname)
				if !found {
					common.Log.Error("Shading not defined in resources. shname=%#q", string(*shname))
					return errors.New("Shading not defined in resources")
				}

				grayShading, err := convertShadingToGray(shading)
				if err != nil {
					return err
				}

				resources.SetShadingByName(*shname, grayShading.GetContext().ToPdfObject())
			}
			*processedOperations = append(*processedOperations, op)

			return nil
		})
	// Add handler for image related handling.  Note that inline images are completely stored with a ContentStreamInlineImage
	// object as the parameter for BI.
	processor.AddHandler(pdfcontent.HandlerConditionEnumOperand, "BI",
		func(op *pdfcontent.ContentStreamOperation, gs pdfcontent.GraphicsState, resources *pdf.PdfPageResources) error {
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

			cs, err := iimg.GetColorSpace(resources)
			if err != nil {
				common.Log.Error("Error getting color space for inline image: %v", err)
				return err
			}

			if cs.GetNumComponents() == 1 {
				return nil
			}

			encoder, err := iimg.GetEncoder()
			if err != nil {
				common.Log.Error("Error getting encoder for inline image: %v", err)
				return err
			}

			if ignoreGrayFilters {
				switch encoder.GetFilterName() {
				// TODO: Add JPEG2000 encoding/decoding. Until then we assume JPEG200 images are color
				//case "JPXDecode":
				//	return nil
				// These filters are only used with grayscale images
				case "CCITTDecode", "CCITTFaxDecode", "JBIG2Decode":
					return nil
				}
			}

			img, err := iimg.ToImage(resources)
			if err != nil {
				common.Log.Error("Error converting inline image to image: %v", err)
				return err
			}
			rgbImg, err := cs.ImageToRGB(*img)
			if err != nil {
				common.Log.Error("Error converting image to rgb: %v", err)
				return err
			}
			rgbColorSpace := pdf.NewPdfColorspaceDeviceRGB()
			grayImage, err := rgbColorSpace.ImageToGray(rgbImg)
			if err != nil {
				common.Log.Error("Error converting img to gray: %v", err)
				return err
			}

			// Update the XObject image.
			// Use same encoder as input data.  Make sure for DCT filter it is updated to 1 color component.

			if dctEncoder, is := encoder.(*pdfcore.DCTEncoder); is {
				dctEncoder.ColorComponents = 1
			}

			grayInlineImg, err := pdfcontent.NewInlineImageFromImage(grayImage, encoder)
			if err != nil {
				if err == pdfcore.ErrUnsupportedEncodingParameters {
					// Unsupported encoding parameters, revert to a basic flate encoder without predictor.
					encoder = pdfcore.NewFlateEncoder()
				}
				// Try again, fail on error.
				grayInlineImg, err = pdfcontent.NewInlineImageFromImage(grayImage, encoder)
				if err != nil {
					common.Log.Error("Error making a new inline image object: %v", err)
					return err
				}
			}

			// Replace inline image data with the gray image.
			pOp := pdfcontent.ContentStreamOperation{}
			pOp.Operand = "BI"
			pOp.Params = []pdfcore.PdfObject{grayInlineImg}
			*processedOperations = append(*processedOperations, &pOp)

			return nil
		})

	// Handler for XObject Image and Forms.
	processedXObjects := map[string]bool{} // Keep track of processed XObjects to avoid repetition.

	processor.AddHandler(pdfcontent.HandlerConditionEnumOperand, "Do",
		func(op *pdfcontent.ContentStreamOperation, gs pdfcontent.GraphicsState, resources *pdf.PdfPageResources) error {
			if len(op.Params) < 1 {
				common.Log.Error("Invalid number of params for Do object")
				return errors.New("Range check")
			}

			// XObject.
			name := op.Params[0].(*pdfcore.PdfObjectName)
			common.Log.Debug("Name=%#v=%#q", name, string(*name))

			// Only process each one once.
			_, has := processedXObjects[string(*name)]
			if has {
				return nil
			}
			processedXObjects[string(*name)] = true

			_, xtype := resources.GetXObjectByName(*name)
			common.Log.Debug("xtype=%+v pdf.XObjectTypeImage=%v", xtype, pdf.XObjectTypeImage)

			if xtype == pdf.XObjectTypeImage {

				ximg, err := resources.GetXObjectImageByName(*name)
				if err != nil {
					common.Log.Error("Error w/GetXObjectImageByName : %v", err)
					return err
				}

				if ignoreGrayFilters {
					switch ximg.Filter.GetFilterName() {
					// TODO: Add JPEG2000 encoding/decoding. Until then we assume JPEG200 images are color
					//case "JPXDecode":
					//	return nil
					// These filters are only used with grayscale images
					case "CCITTDecode", "CCITTFaxDecode", "JBIG2Decode":
						return nil
					}
				}

				// Hacky workaround for Szegedy_Going_Deeper_With_2015_CVPR_paper.pdf that has a colored image
				// that is completely masked
				if ximg.Filter.GetFilterName() == "RunLengthDecode" && ximg.SMask != nil {
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

				rgbColorSpace := pdf.NewPdfColorspaceDeviceRGB()
				grayImage, err := rgbColorSpace.ImageToGray(rgbImg)
				if err != nil {
					common.Log.Error("Error ImageToGray: %v", err)
					return err
				}

				// Update the XObject image.
				// Use same encoder as input data.  Make sure for DCT filter it is updated to 1 color component.
				encoder := ximg.Filter
				if dctEncoder, is := encoder.(*pdfcore.DCTEncoder); is {
					dctEncoder.ColorComponents = 1
				}

				ximgGray, err := pdf.NewXObjectImageFromImage(&grayImage, nil, encoder)
				if err != nil {
					if err == pdfcore.ErrUnsupportedEncodingParameters {
						// Unsupported encoding parameters, revert to a basic flate encoder without predictor.
						encoder = pdfcore.NewFlateEncoder()
					}

					// Try again, fail if error.
					ximgGray, err = pdf.NewXObjectImageFromImage(&grayImage, nil, encoder)
					if err != nil {
						common.Log.Error("Error creating image: %v", err)
						return err
					}
				}

				// Update the entry.
				err = resources.SetXObjectImageByName(*name, ximgGray)
				if err != nil {
					common.Log.Error("Failed setting x object: %v (%s)", err, string(*name))
					return err
				}
			} else if xtype == pdf.XObjectTypeForm {
				common.Log.Debug(" XObject Form: %s", *name)

				// Go through the XObject Form content stream.
				xform, err := resources.GetXObjectFormByName(*name)
				if err != nil {
					common.Log.Error("err=%v", err)
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
				grayContent, err := transformContentStreamToGrayscale(string(formContent), formResources)
				if err != nil {
					common.Log.Error("err=%v", err)
					return err
				}

				xform.SetContentStream(grayContent, nil)

				// Update the resource entry.
				resources.SetXObjectFormByName(*name, xform)
			}

			return nil
		})

	err = processor.Process(resources)
	if err != nil {
		common.Log.Error("processor.Process returned: err=%v", err)
		return nil, err
	}

	return processedOperations.Bytes(), nil
}

// convertPatternToGray converts `pattern` to grayscale (tiling or shading pattern).
func convertPatternToGray(pattern *pdf.PdfPattern) (*pdf.PdfPattern, error) {
	// Case 1: Colored tiling patterns.  Need to process the content stream and replace.
	if pattern.IsTiling() {
		tilingPattern := pattern.GetAsTilingPattern()

		if tilingPattern.IsColored() {
			// A colored tiling pattern can use color operators in its stream, need to process the stream.

			content, err := tilingPattern.GetContentStream()
			if err != nil {
				return nil, err
			}

			grayContents, err := transformContentStreamToGrayscale(string(content), tilingPattern.Resources)
			if err != nil {
				return nil, err
			}

			tilingPattern.SetContentStream(grayContents, nil)

			// Update in-memory pdf objects.
			_ = tilingPattern.ToPdfObject()
		}
	} else if pattern.IsShading() {
		// Case 2: Shading patterns.  Need to create a new colorspace that can map from N=3,4 colorspaces to grayscale.
		shadingPattern := pattern.GetAsShadingPattern()

		grayShading, err := convertShadingToGray(shadingPattern.Shading)
		if err != nil {
			return nil, err
		}
		shadingPattern.Shading = grayShading

		// Update in-memory pdf objects.
		_ = shadingPattern.ToPdfObject()
	}

	return pattern, nil
}

// convertShadingToGray converts `shading` to grayscale.
// This one is slightly involved as a shading defines a color as function of position, i.e. color(x,y) = F(x,y).
// Since the function can be challenging to change, we define new DeviceN colorspace with a color conversion
// function.
func convertShadingToGray(shading *pdf.PdfShading) (*pdf.PdfShading, error) {
	cs := shading.ColorSpace

	if cs.GetNumComponents() == 1 {
		// Already grayscale, should be fine. No action taken.

		// Make sure is device gray.
		shading.ColorSpace = pdf.NewPdfColorspaceDeviceGray()

		return shading, nil
	} else if cs.GetNumComponents() == 3 {
		// Create a new DeviceN colorspace that converts R,G,B -> Grayscale
		// Use: gray := 0.3*R + 0.59G + 0.11B
		// PS program: { 0.11 mul exch 0.59 mul add exch 0.3 mul add }.
		transformFunc := &pdf.PdfFunctionType4{}
		transformFunc.Domain = []float64{0, 1, 0, 1, 0, 1}
		transformFunc.Range = []float64{0, 1}
		rgbToGrayPsProgram := ps.NewPSProgram()
		rgbToGrayPsProgram.Append(ps.MakeReal(0.11))
		rgbToGrayPsProgram.Append(ps.MakeOperand("mul"))
		rgbToGrayPsProgram.Append(ps.MakeOperand("exch"))
		rgbToGrayPsProgram.Append(ps.MakeReal(0.59))
		rgbToGrayPsProgram.Append(ps.MakeOperand("mul"))
		rgbToGrayPsProgram.Append(ps.MakeOperand("add"))
		rgbToGrayPsProgram.Append(ps.MakeOperand("exch"))
		rgbToGrayPsProgram.Append(ps.MakeReal(0.3))
		rgbToGrayPsProgram.Append(ps.MakeOperand("mul"))
		rgbToGrayPsProgram.Append(ps.MakeOperand("add"))
		transformFunc.Program = rgbToGrayPsProgram

		// Define the DeviceN colorspace that performs the R,G,B -> Gray conversion for us.
		transformcs := pdf.NewPdfColorspaceDeviceN()
		transformcs.AlternateSpace = pdf.NewPdfColorspaceDeviceGray()
		transformcs.ColorantNames = pdfcore.MakeArray(pdfcore.MakeName("R"), pdfcore.MakeName("G"), pdfcore.MakeName("B"))
		transformcs.TintTransform = transformFunc

		// Replace the old colorspace with the new.
		shading.ColorSpace = transformcs

		return shading, nil
	} else if cs.GetNumComponents() == 4 {
		// Create a new DeviceN colorspace that converts C,M,Y,K -> Grayscale.
		// Use: gray = 1.0 - min(1.0, 0.3*C + 0.59*M + 0.11*Y + K)  ; where BG(k) = k simply.
		// PS program: {exch 0.11 mul add exch 0.59 mul add exch 0.3 mul add dup 1.0 ge { pop 1.0 } if}
		transformFunc := &pdf.PdfFunctionType4{}
		transformFunc.Domain = []float64{0, 1, 0, 1, 0, 1, 0, 1}
		transformFunc.Range = []float64{0, 1}

		cmykToGrayPsProgram := ps.NewPSProgram()
		cmykToGrayPsProgram.Append(ps.MakeOperand("exch"))
		cmykToGrayPsProgram.Append(ps.MakeReal(0.11))
		cmykToGrayPsProgram.Append(ps.MakeOperand("mul"))
		cmykToGrayPsProgram.Append(ps.MakeOperand("add"))
		cmykToGrayPsProgram.Append(ps.MakeOperand("exch"))
		cmykToGrayPsProgram.Append(ps.MakeReal(0.59))
		cmykToGrayPsProgram.Append(ps.MakeOperand("mul"))
		cmykToGrayPsProgram.Append(ps.MakeOperand("add"))
		cmykToGrayPsProgram.Append(ps.MakeOperand("exch"))
		cmykToGrayPsProgram.Append(ps.MakeReal(0.30))
		cmykToGrayPsProgram.Append(ps.MakeOperand("mul"))
		cmykToGrayPsProgram.Append(ps.MakeOperand("add"))
		cmykToGrayPsProgram.Append(ps.MakeOperand("dup"))
		cmykToGrayPsProgram.Append(ps.MakeReal(1.0))
		cmykToGrayPsProgram.Append(ps.MakeOperand("ge"))
		// Add sub procedure.
		subProc := ps.NewPSProgram()
		subProc.Append(ps.MakeOperand("pop"))
		subProc.Append(ps.MakeReal(1.0))
		cmykToGrayPsProgram.Append(subProc)
		cmykToGrayPsProgram.Append(ps.MakeOperand("if"))
		transformFunc.Program = cmykToGrayPsProgram

		// Define the DeviceN colorspace that performs the R,G,B -> Gray conversion for us.
		transformcs := pdf.NewPdfColorspaceDeviceN()
		transformcs.AlternateSpace = pdf.NewPdfColorspaceDeviceGray()
		transformcs.ColorantNames = pdfcore.MakeArray(pdfcore.MakeName("C"), pdfcore.MakeName("M"), pdfcore.MakeName("Y"), pdfcore.MakeName("K"))
		transformcs.TintTransform = transformFunc

		// Replace the old colorspace with the new.
		shading.ColorSpace = transformcs

		return shading, nil
	} else {
		common.Log.Debug("Cannot convert to shading pattern grayscale, color space N = %d", cs.GetNumComponents())
		return nil, errors.New("Unsupported pattern colorspace for grayscale conversion")
	}
}

// modifyPath returns `inputPath` with its directory replaced by `outputDir`
func modifyPath(inputPath, outputDir string) string {
	_, name := filepath.Split(inputPath)
	// name = fmt.Sprintf("%08d_%s", fileSize(inputPath), name)

	outputPath := filepath.Join(outputDir, name)
	in, err := filepath.Abs(inputPath)
	if err != nil {
		panic(err)
	}
	out, err := filepath.Abs(outputPath)
	if err != nil {
		panic(err)
	}
	if strings.ToLower(in) == strings.ToLower(out) {
		common.Log.Error("modifyPath: Cannot modify path to itself. inputPath=%#q outputDir=%#q",
			inputPath, outputDir)
		panic("Don't write over test files")
	}
	return outputPath
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
	common.Log.Trace("runGhostscript: pdf=%#q outputDir=%#q", pdf, outputDir)
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
	common.Log.Trace("runGhostscript: cmd=%#q", cmd.Args)

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

// runPdfToPs runs pdftops on file `pdf` to create a PostScript file in directory `outputDir`
func runPdfToPs(pdf, outputDir string) error {
	common.Log.Trace("pdf=%#q outputDir=%#q", pdf, outputDir)
	ps := changeDir(pdf, outputDir)
	cmd := exec.Command("pdftops", pdf, ps)
	common.Log.Trace("cmd=%#q", cmd.Args)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		common.Log.Error("Could not process pdf=%q err=%v\nstdout=\n%s\nstderr=\n%s\n",
			pdf, err, stdout, stderr)
	}
	return err
}

// isPdfColor returns true if PDF files `path` has color marks on any page
// If `keep` is true then the page rasters are retained
func isPdfColor(path, temp string, showPages, keep bool) (bool, []int, error) {
	dir := filepath.Join(temp, "color")
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		panic(err)
	}
	if !keep {
		defer removeDir(dir)
	}

	err = runGhostscript(path, dir)
	if err != nil {
		return false, nil, err
	}

	if showPages {
		colorPages, err := colorDirectoryPages("*.png", dir, keep)
		return len(colorPages) > 0, colorPages, err
	}

	isColor, err := isColorDirectory("*.png", dir)
	return isColor, nil, err
}

// isColorDirectory returns true if any of the image files that match `mask` in directories `dir`
// has any color pixels
func isColorDirectory(mask, dir string) (bool, error) {
	pattern := filepath.Join(dir, mask)
	files, err := filepath.Glob(pattern)
	if err != nil {
		common.Log.Error("isColorDirectory: Glob failed. pattern=%#q err=%v", pattern, err)
		return false, err
	}

	for _, path := range files {
		isColor, err := isColorImage(path, false)
		if isColor || err != nil {
			return isColor, err
		}
	}
	return false, nil
}

// colorDirectoryPages returns a lists of the page numbers of the image files that match `mask` in
// directories `dir` that have any color pixels.
func colorDirectoryPages(mask, dir string, keep bool) ([]int, error) {
	pattern := filepath.Join(dir, mask)
	files, err := filepath.Glob(pattern)
	if err != nil {
		common.Log.Error("isColorDirectory: Glob failed. pattern=%#q err=%v", pattern, err)
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
		isColor, err := isColorImage(path, keep)
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
func isColorImage(path string, keep bool) (bool, error) {
	img, err := readImage(path)
	if err != nil {
		return false, err
	}
	isColor := imgIsColor(img)
	if isColor && keep {
		markedPath := fmt.Sprintf("%s.marked.png", path)
		markedImg, summary := imgMarkColor(img)
		common.Log.Error("markedPath=%#q %s", markedPath, summary)
		err = writeImage(markedPath, markedImg)
	}
	return isColor, err
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
				return true
			}
		}
	}
	return false
}

func imgMarkColor(imgIn image.Image) (image.Image, string) {
	img := image.NewNRGBA(imgIn.Bounds())
	black := color.RGBA{0, 0, 0, 255}
	// white := color.RGBA{255, 255, 255, 255}
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	data := []float64{}
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			rr, gg, bb, _ := img.At(x, y).RGBA()
			r, g, b := float64(rr)*F, float64(gg)*F, float64(bb)*F
			if math.Abs(r-g) > colorThreshold || math.Abs(r-b) > colorThreshold || math.Abs(g-b) > colorThreshold {
				img.Set(x, y, black)
				data = append(data, math.Abs(r-g)+math.Abs(r-b)+math.Abs(g-b))
			}
		}
	}
	return img, summarizeSeries(data)
}

func summarizeSeries(data []float64) string {
	n := len(data)
	total := 0.0
	min := +1e20
	max := -1e20
	for _, x := range data {
		total += x
		if x < min {
			min = x
		}
		if x > max {
			max = x
		}
	}
	mean := total / float64(n)
	return fmt.Sprintf("n=%d min=%.3f mean=%.3f max=%.3f", n, min, mean, max)
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

// writeImage writes image `img` to file `path`
func writeImage(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		common.Log.Error("writeImage: Could not create file. path=%#q err=%v", path, err)
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
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
				fmt.Printf("Not a regular file. %#q\n", path)
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

func changeDir(path, dir string) string {
	_, name := filepath.Split(path)
	return filepath.Join(dir, name)
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
