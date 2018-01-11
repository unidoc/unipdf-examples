/*
 * Transform all content streams in all pages in a list of pdf files.
 *
 * There are currently 2 transforms implemented.
 *	- Identity transform. No -t command line option
 *	- Grayscale transform. -t command line option
 *
 * The identity transform
 *	- converts PDF files into our internal representation
 *	- converts the internal representation back to a PDF file
 *	- checks that the output PDF file is the same as the input PDF file
 *
 * The grayscale transform
 *	- converts PDF files into our internal representation
 *	- transforms the internal representation to grayscale
 *	- converts the internal representation back to a PDF file
 *	- checks that the output PDF file is grayscale
 *
 * Run as: ./pdf_transform_content_streams -o output [-d] [-t] testdata/*.pdf > blah
 *
 * This will transform all .pdf file in testdata and write the results to output.
 * The main results are written to stderr so you will see them in your console.
 * Detailed information is written to stdout and you will see them in blah.
 *
 *  See the other command line options in the top of main()
 *		-a tests all the input files. The default behavior is stop at the first failure. Use this
 *			to find out how many of your corpus files this program works for.
 *		-x will transform without parsing content streams. Use this to see which failures are due to
 *			problems in the content parsing code.
 *			Running -a then -a -x will tell you how well this code is performing on your corpus
 *			and which failures are due to content parsing.
 *			(-x will disable -g)
 */

package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	// unilicense "github.com/unidoc/unidoc/license"
	unicommon "github.com/unidoc/unidoc/common"
	unicontent "github.com/unidoc/unidoc/pdf/contentstream"
	unipdf "github.com/unidoc/unidoc/pdf/model"
)

func initUniDoc(licenseKey string, debug bool) error {
	// PETER: I can't find github.com/unidoc/unidoc/license so I have comment out the license code
	//        in this example program.
	// if len(licenseKey) > 0 {
	// 	err := unilicense.SetLicenseKey(licenseKey)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// To make the library log we just have to initialise the logger which satisfies
	// the unicommon.Logger interface, unicommon.DummyLogger is the default and
	// does not do anything. Very easy to implement your own.
	// unicommon.SetLogger(unicommon.DummyLogger{})
	unicommon.SetLogger(unicommon.ConsoleLogger{LogLevel: unicommon.LogLevelDebug})

	return nil
}

// imageThreshold represents the threshold for image difference in image comparisons
type imageThreshold struct {
	fracPixels float64 // Fraction of pixels in page raster that may differ
	mean       float64 // Max mean difference on scale 0..255 for pixels that differ
}

// identityThreshold is the imageThreshold for identity transforms in this program.
var identityThreshold = imageThreshold{
	fracPixels: 1.0e-4, // Fraction of pixels in page raster that may differ
	mean:       10.0,   // Max mean difference on scale 0..255 for pixels that differ
}

var testStats = statistics{
	enabled:        true,
	testResultPath: "xform.test.results.csv",
	// imageInfoPath:  "xform.image.info.csv",
}

var allOpCounts = map[string]int{}

func main() {
	debug := false                // Write debug level info to stdout?
	keep := false                 // Keep the rasters used for PDF comparison"
	noContentTransforms := false  // Don't parse stream contents?
	doGrayscaleTransform := false // Apply the grayscale transform?
	compareGrayscale := false     // Do PDF raster comparison on grayscale rasters?
	runAllTests := false          // Don't stop when a PDF file fails to process?
	outputDir := ""               // Transformed PDFs are written here
	var minSize int64 = -1        // Minimum size for an input PDF to be processed.
	var maxSize int64 = -1        // Maximum size for an input PDF to be processed.
	flag.BoolVar(&debug, "d", false, "Enable debug logging")
	flag.BoolVar(&keep, "k", false, "Keep the rasters used for PDF comparison")
	flag.BoolVar(&noContentTransforms, "x", false, "Don't transform streams")
	flag.BoolVar(&doGrayscaleTransform, "t", false, "Do grayscale transform")
	flag.BoolVar(&compareGrayscale, "g", false, "Do PDF raster comparison on grayscale rasters")
	flag.BoolVar(&runAllTests, "a", false, "Run all tests. Don't stop at first failure")
	flag.StringVar(&outputDir, "o", "", "Output directory")
	flag.Int64Var(&minSize, "min", -1, "Minimum size of files to process (bytes)")
	flag.Int64Var(&maxSize, "max", -1, "Maximum size of files to process (bytes)")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 || len(outputDir) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s -o <output directory> [-d] <file1> <file2> ...\n",
			os.Args[0])
		os.Exit(1)
	}

	err := initUniDoc("", debug)
	if err != nil {
		os.Exit(1)
	}
	unicontent.ValidatingOperations = true

	err = os.MkdirAll(outputDir, 0777)
	if err != nil {
		unicommon.Log.Error("MkdirAll failed. outputDir=%#q err=%v", outputDir, err)
		os.Exit(1)
	}

	compDir := makeUniqueDir("compare.pdfs")
	fmt.Fprintf(os.Stderr, "compDir=%#q\n", compDir)
	if !keep {
		defer removeDir(compDir)
	}

	pdfList, err := patternsToPaths(args)
	if err != nil {
		unicommon.Log.Error("patternsToPaths failed. args=%#q err=%v", args, err)
		os.Exit(1)
	}
	pdfList = sortFiles(pdfList, minSize, maxSize)
	badFiles := []string{}
	failFiles := []string{}

	if err = testStats.load(); err != nil {
		unicommon.Log.Error("stats.load failed. testStats=%+v err=%v", testStats, err)
		os.Exit(1)
	}
	defer testStats._save()

	for idx, inputPath := range pdfList {

		_, name := filepath.Split(inputPath)
		inputSize := fileSize(inputPath)

		fmt.Fprintf(os.Stderr, "%3d of %d %#-30q  (%6d->", idx,
			len(pdfList), name, inputSize)
		outputPath := modifyPath(inputPath, outputDir)

		objCounts := ObjCounts{xobjNameSubtype: map[string]string{}}
		t0 := time.Now()
		numPages, err := transformPdfFile(inputPath, outputPath, noContentTransforms,
			doGrayscaleTransform, &objCounts)
		dt := time.Since(t0)
		if err != nil {
			unicommon.Log.Error("transformPdfFile failed. err=%v", err)
			failFiles = append(failFiles, inputPath)
			if runAllTests {
				continue
			}
			os.Exit(1)
		}

		outputSize := fileSize(outputPath)
		fmt.Fprintf(os.Stderr, "%6d %3d%%) %d pages %.3f sec => %#q\n",
			outputSize, int(float64(outputSize)/float64(inputSize)*100.0+0.5),
			numPages, dt.Seconds(), outputPath)

		if doGrayscaleTransform {
			isColorIn, _, err := isPdfColor(inputPath, compDir, false, keep)
			isColorOut, colorPagesOut, err := isPdfColor(outputPath, compDir, true, keep)
			xobjForm := 0
			xobjImg := 0
			for _, subtype := range objCounts.xobjNameSubtype {
				if subtype == "Form" {
					xobjForm++
				} else if subtype == "Image" {
					xobjImg++
				}
			}

			e := testResult{
				name:     path.Base(inputPath),
				colorIn:  isColorIn,
				colorOut: isColorOut,
				numPages: numPages,
				duration: dt.Seconds(),
				xobjForm: xobjForm,
				xobjImg:  xobjImg,
			}
			testStats.addTestResult(e, true)

			if err != nil || isColorOut {
				if err != nil {
					unicommon.Log.Error("Transform has damaged PDF. err=%v\n\tinputPath=%#q\n\toutputPath=%#q",
						err, inputPath, outputPath)
				} else {
					unicommon.Log.Error("isPdfColor: %d Color pages", len(colorPagesOut))
				}
				failFiles = append(failFiles, inputPath)
				if runAllTests {
					continue
				}
				os.Exit(1)
			}
		} else {
			equal, badInput, err := pdfsEqual(inputPath, outputPath, identityThreshold, compDir, true, keep)
			if badInput {
				unicommon.Log.Error("Bad input PDF. inputPath=%#q err=%v", inputPath, err)
				badFiles = append(badFiles, inputPath)
				continue
			}
			if err != nil || !equal {
				if err != nil {
					unicommon.Log.Error("Transform has damaged PDF. err=%v\n\tinputPath=%#q\n\toutputPath=%#q",
						err, inputPath, outputPath)
				} else {
					unicommon.Log.Error("Transform has changed PDF.\n\tinputPath=%#q\n\toutputPath=%#q",
						inputPath, outputPath)
				}
				failFiles = append(failFiles, inputPath)
				if runAllTests {
					continue
				}
				os.Exit(1)
			}
		}
	}

	fmt.Fprintf(os.Stderr, "%d files %d bad %d failed\n", len(pdfList), len(badFiles), len(failFiles))
	fmt.Fprintf(os.Stderr, "%d bad\n", len(badFiles))
	for i, path := range badFiles {
		fmt.Fprintf(os.Stderr, "%3d %#q\n", i, path)
	}
	fmt.Fprintf(os.Stderr, "%d fail\n", len(failFiles))
	for i, path := range failFiles {
		fmt.Fprintf(os.Stderr, "%3d %#q\n", i, path)
	}

	printOpCounts("operations in all PDFs", allOpCounts)
	printCsCounts("color spaces in all PDFs", allCsCounts)
}

type ObjCounts struct {
	xobjNameSubtype map[string]string
}

// transformPdfFile transforms PDF `inputPath` and writes the resulting PDF to `outputPath`
// If `noContentTransforms` is true then stream contents are not parsed
func transformPdfFile(inputPath, outputPath string, noContentTransforms, doGrayscaleTransform bool,
	objCounts *ObjCounts) (int, error) {

	docCsCounts = map[string]int{}

	f, err := os.Open(inputPath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	pdfReader, err := unipdf.NewPdfReader(f)
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

	pdfWriter := unipdf.NewPdfWriter()

	for i := 0; i < numPages; i++ {
		pageNum := i + 1
		page := pdfReader.PageList[i]
		unicommon.Log.Debug("^^^^page %d", pageNum)

		if !noContentTransforms {
			desc := fmt.Sprintf("%s:page%d", filepath.Base(inputPath), pageNum)
			err = transformPdfPage(page, desc, doGrayscaleTransform, objCounts)
			if err != nil {
				return numPages, err
			}
		}

		err = pdfWriter.AddPage(page)
		if err != nil {
			return numPages, err
		}
	}

	fWrite, err := os.Create(outputPath)
	if err != nil {
		return numPages, err
	}
	defer fWrite.Close()
	err = pdfWriter.Write(fWrite)

	return numPages, nil
}

// transformPdfPage
//	- parses the contents of streams in `page` into a slice of operations
//	- converts the slice of operations into a stream
//	- replaces the streams in `page` with the new stream
func transformPdfPage(page *unipdf.PdfPage, desc string, doGrayscaleTransform bool,
	objCounts *ObjCounts) error {

	nameSubtype, err := unipdf.GetXObjectSubtypes(page)
	if err != nil {
		return nil
	}
	unicommon.Log.Info("nameSubtype=%+v", nameSubtype)
	for name, subtype := range nameSubtype {
		objCounts.xobjNameSubtype[name] = subtype
	}

	err = transformPdfPageContent(page, desc, doGrayscaleTransform, objCounts)
	if err != nil {
		return err
	}

	err = transformXObjects(page, desc, doGrayscaleTransform)
	if err != nil {
		return err
	}

	err = transformColorspaces(page, desc, doGrayscaleTransform)
	if err != nil {
		return err
	}

	return nil
}

func transformColorspaces(page unipdf.PdfFormPage, desc string, doGrayscaleTransform bool) error {
	unicommon.Log.Info("desc=%s doGrayscaleTransform=%t", desc, doGrayscaleTransform)

	resources, err := page.GetResources()
	if err != nil {
		return err
	}
	colorspaces := resources.ColorSpace
	if colorspaces == nil {
		return nil
	}
	names := colorspaces.Names
	sort.Strings(names)
	for _, name := range names {
		obj := colorspaces.Colorspaces[name]
		// For indexed color spaces, we just modify the lookup table to grayscale
		if cs, isIndexed := obj.(*unipdf.PdfColorspaceSpecialIndexed); isIndexed {
			err := cs.ColorspaceToGray()
			if err != nil {
				return err
			}
		}
	}
	// err = transformImageXObjects(page, desc, doGrayscaleTransform)
	// if err != nil {
	// 	return err
	// }
	// err = transformFormXObjects(page, desc, doGrayscaleTransform)
	// if err != nil {
	// 	return err
	// }

	// xobjs, err = unipdf.GetXObjects(page)
	// if err != nil {
	// 	return nil
	// }
	// unicommon.Log.Info("+XObjects=%s", xobjs)
	return nil
}

func transformColorspaces(page unipdf.PdfFormPage, desc string, doGrayscaleTransform bool) error {
	unicommon.Log.Info("desc=%s doGrayscaleTransform=%t", desc, doGrayscaleTransform)

	xobjs, err := unipdf.GetXObjects(page)
	if err != nil {
		return err
	}
	colorspaces := resources.ColorSpace
	if colorspaces == nil {
		return nil
	}
	names := colorspaces.Names
	sort.Strings(names)
	for _, name := range names {
		obj := colorspaces.Colorspaces[name]
		// For indexed color spaces, we just modify the lookup table to grayscale
		if cs, isIndexed := obj.(*unipdf.PdfColorspaceSpecialIndexed); isIndexed {
			err := cs.ColorspaceToGray()
			if err != nil {
				return err
			}
		}
	}
	// err = transformImageXObjects(page, desc, doGrayscaleTransform)
	// if err != nil {
	// 	return err
	// }
	// err = transformFormXObjects(page, desc, doGrayscaleTransform)
	// if err != nil {
	// 	return err
	// }

	// xobjs, err = unipdf.GetXObjects(page)
	// if err != nil {
	// 	return nil
	// }
	// unicommon.Log.Info("+XObjects=%s", xobjs)

	return nil
}

func transformXObjects(page unipdf.PdfFormPage, desc string, doGrayscaleTransform bool) error {
	unicommon.Log.Info("desc=%s doGrayscaleTransform=%t", desc, doGrayscaleTransform)

	xobjs, err := unipdf.GetXObjects(page)
	if err != nil {
		return err
	}
	unicommon.Log.Info("-XObjects=%s", xobjs)

	// err = transformImageXObjects(page, desc, doGrayscaleTransform)
	// if err != nil {
	// 	return err
	// }
	err = transformFormXObjects(page, desc, doGrayscaleTransform)
	if err != nil {
		return err
	}

	xobjs, err = unipdf.GetXObjects(page)
	if err != nil {
		return nil
	}
	unicommon.Log.Info("+XObjects=%s", xobjs)

	return nil
}

func transformImageXObjects(page unipdf.PdfFormPage, desc string, doGrayscaleTransform bool) error {
	unicommon.Log.Info("desc=%s doGrayscaleTransform=%t", desc, doGrayscaleTransform)

	nameXimgMap, err := unipdf.GetXImageResourceMap(page)
	if err != nil {
		unicommon.Log.Error("No resource map. err=%v", err)
		return err
	}

	names := []string{}
	for name := range nameXimgMap {
		names = append(names, name)
	}
	sort.Strings(names)

	unicommon.Log.Info("nameXimgMap=%d %+v", len(nameXimgMap), names)
	for _, name := range names {
		ximg := nameXimgMap[name]
		if ximg == nil {
			panic("nil Image XObject")
		}
		unicommon.Log.Info("Converting image XObject %#q to gray ximg=%s", name, ximg)
		if err := ximg.ToGray(); err != nil {
			return err
		}
	}

	err = unipdf.SetXImageResourceMap(page, nameXimgMap)
	if err != nil {
		unicommon.Log.Error("SetXImageResourceMap failed. err=%v", err)
		return err
	}

	return nil
}

func transformFormXObjects(page unipdf.PdfFormPage, desc string, doGrayscaleTransform bool) error {
	unicommon.Log.Info("desc=%s doGrayscaleTransform=%t", desc, doGrayscaleTransform)

	nameXformMap, err := unipdf.GetXFormResourceMap(page)
	if err != nil {
		unicommon.Log.Error("No resource map. err=%v", err)
		return err
	}

	names := []string{}
	for name := range nameXformMap {
		names = append(names, name)
	}
	sort.Strings(names)

	unicommon.Log.Info("nameXformMap=%d %+v", len(nameXformMap), names)
	for _, name := range names {
		xform := nameXformMap[name]
		formDesc := fmt.Sprintf("%s:form['%#q']", desc, name)
		unicommon.Log.Info("Converting form XObject %#q to gray", name)
		if err := transformPdfPageContent(xform, formDesc, doGrayscaleTransform, nil); err != nil {
			return err
		}
		if err := transformXObjects(xform, formDesc, doGrayscaleTransform); err != nil {
			return err
		}
	}
	err = unipdf.SetXFormResourceMap(page, nameXformMap)
	if err != nil {
		unicommon.Log.Error("SetXFormResourceMap failed. err=%v", err)
		return err
	}

	return nil
}

func transformPdfPageContent(page unipdf.PdfFormPage, // *unipdf.PdfPage,
	desc string, doGrayscaleTransform bool,
	objCounts *ObjCounts) error {
	cstream, err := page.GetAllContentStreams()
	if err != nil {
		return err
	}
	cstreamOut, err := transformString(page, cstream, desc, doGrayscaleTransform)
	if err != nil {
		return err
	}
	err = page.SetContentStream(cstreamOut, nil)
	if err != nil {
		return err
	}
	return nil
}

// transformString applies the transform (currently identity or grayscale conversion) to string
// and updates `page`
func transformString(page unipdf.PdfFormPage, //*unipdf.PdfPage,
	cstream, desc string, doGrayscaleTransform bool) (string, error) {
	unicommon.Log.Debug("desc=%s doGrayscaleTransform=%t", desc, doGrayscaleTransform)

	operations, err := parseStreamContents(cstream, desc)
	if err != nil {
		printOperations(fmt.Sprintf("%s: ERROR incomplete", desc), operations)
		return "", err
	}

	printOperations(fmt.Sprintf("%s: original", desc), operations)
	printOpCounts(fmt.Sprintf("%s: original", desc), getOpCounts(operations))
	if doGrayscaleTransform {
		if err := transformColorToGrayscale(page, desc, &operations); err != nil {
			return "", err
		}
		printOpCounts(fmt.Sprintf("%s: transformed ", desc), getOpCounts(operations))
	}
	printCsCounts(desc, docCsCounts)
	printCsCounts("color spaces in all PDFs", allCsCounts)

	opStrings := []string{}
	for _, op := range operations {
		opStrings = append(opStrings, op.DefaultWriteString())
	}
	cstreamOut := strings.Join(opStrings, " ")

	return cstreamOut, nil
}

// parsePageContents parses the contents of streams in `cstream` and returns them as a slice of
// operations
func parseStreamContents(cstream string, desc string) ([]*unicontent.ContentStreamOperation, error) {
	cstreamParser := unicontent.NewContentStreamParser(cstream)
	unicommon.Log.Debug("%s cstream=\n'%s'\nXXXXXX", desc, cstream)
	return cstreamParser.Parse()
}

// Per-document color space counts
var docCsCounts = map[string]int{}

// Color space counts for all documents
var allCsCounts = map[string]int{}

// transformColorToGrayscale transforms color pages to grayscale
func transformColorToGrayscale(page unipdf.PdfFormPage, //*unipdf.PdfPage,
	desc string,
	pOperations *[]*unicontent.ContentStreamOperation) (err error) {

	unicommon.Log.Debug("%s", desc)

	resources, err := page.GetResources()
	if err != nil {
		panic(err) // !@#$ REmove in production code
	}

	var colorspaceMap map[string]unipdf.PdfColorspace
	if resources.ColorSpace != nil {
		colorspaceMap = resources.ColorSpace.Colorspaces
	}
	gs := graphicState{}
	gs.colorspaceStroke, _ = unipdf.NewPdfColorspaceFromPdfObject(nil)
	gs.colorspaceFill, _ = unipdf.NewPdfColorspaceFromPdfObject(nil)

	gsStack := graphicStateStack{}

	noContentColor := false

	// op0 := unicontent.ContentStreamOperation{Operand: "sc"}
	// op0 := unicontent.ContentStreamOperation{Operand: "SC"}
	badOps := map[string]bool{
	// "l":  true,
	// "BX": true,
	// "EX": true,
	// "c":  true,
	// "m":  true,
	// "Do": true,
	// // "gs": true,
	// "cm":  true,
	// "scn": true,
	// "TD":  true,
	// "BDC": true,
	// "EMC": true,
	// "sh": true, // !@#$ Fix shading
	}
	//  0:   `re`   304
	//  1:    `f`   203
	//  2:   `Tw`   122
	//  3:   `Tc`   118
	//  4:   `Tm`    88
	//  5:   `BT`    85
	//  6:   `ET`    85
	//  7:   `Tj`    83
	//  8:   `TJ`    82
	//  9:    `Q`    68
	// 10:    `q`    68
	// 11:    `n`    67
	// 12:    `W`    67
	// 13:   `Td`    64
	// 14:  `scn`    58
	// 15:  `BDC`    55
	// 16:  `EMC`    55
	// 17:   `f*`    34
	// 18:   `Tf`    34
	// 19:   `TD`     3
	// 20:   `cm`     1
	// 21:   `cs`     1
	// 22:   `Do`     1
	// 23:   `gs`     1
	removeOps(pOperations, badOps, false)

	for i, op := range *pOperations {
		unicommon.Log.Debug("i=%d op=%s", i, op)
		var vals []float64
		switch op.Operand {
		case "q":
			gsStack.push(gs)
		case "Q":
			gs = gsStack.pop()
			unicommon.Log.Debug("gs=%+v", gs)
		case "cs", "CS":
			name, err := op.GetNameParam()
			if err != nil {
				return err
			}
			colorspace, fromMap := colorspaceMap[name]
			if !fromMap {
				colorspace, err = unipdf.NewPdfColorspaceFromPdfObject(op.Params[0])
				if err != nil {
					unicommon.Log.Error("No color space name=%s err=%v", op.Params[0], err)
					return err
				}
			}
			if unipdf.PdfColorspaceHasColor(colorspace) {
				if err = op.SetNameParam("DeviceGray"); err != nil {
					return err
				}
			}
			if isUpper(op.Operand) {
				gs.colorspaceStroke = colorspace
			} else {
				gs.colorspaceFill = colorspace
			}
			unicommon.Log.Debug("gs=%s fromMap=%t", gs, fromMap)
			csName := fmt.Sprintf("%#T", colorspace)
			docCsCounts[csName]++
			allCsCounts[csName]++

		case "sc", "SC", "scn", "SCN":
			if noContentColor {
				*op = unicontent.ContentStreamOperation{}
			} else {
				var colorspace unipdf.PdfColorspace
				if isUpper(op.Operand) {
					colorspace = gs.colorspaceStroke
				} else {
					colorspace = gs.colorspaceFill
				}
				unicommon.Log.Debug("^^^ op=%s hasColor=%t gs=%s",
					op, unipdf.PdfColorspaceHasColor(colorspace), gs)

				// !@#$ SemanticStates201.pdf  HH factsheet.pdf
				if _, ok := colorspace.(*unipdf.PdfColorspaceSpecialPattern); ok {
					*op = unicontent.ContentStreamOperation{}
					// op.SetOpFloatParams(op.Operand, []float64{0.0})
				} else if unipdf.PdfColorspaceHasColor(colorspace) {
					if vals, err = op.GetFloatParams(colorspace.GetNumComponents()); err != nil {
						unicommon.Log.Error("Wrong # params. colorspace=%#T err=%v", colorspace, err)
						return err
					}
					gray, err := unipdf.FloatsToGray(colorspace, vals)
					if err != nil {
						return err
					}
					if err = op.SetOpFloatParams(op.Operand, []float64{gray}); err != nil {
						return err
					}
				}
			}

		case "rg", "RG", "k", "K":
			if noContentColor {
				*op = unicontent.ContentStreamOperation{}
			} else {
				cs := opColorspace[op.Operand]
				if vals, err = op.GetFloatParams(cs.GetNumComponents()); err != nil {
					unicommon.Log.Error("Wrong # params. cs=%#T err=%v", cs, err)
					return err
				}
				gray, err := unipdf.FloatsToGray(cs, vals)
				if err != nil {
					return err
				}
				if err = op.SetOpFloatParams(opGrayOp[op.Operand], []float64{gray}); err != nil {
					return err
				}

				csName := fmt.Sprintf("%#T", cs)
				docCsCounts[csName]++
				allCsCounts[csName]++
			}

		case "g", "G":
			// For instrumentation only
			cs := opColorspace[op.Operand]

			csName := fmt.Sprintf("%#T", cs)
			docCsCounts[csName]++
			allCsCounts[csName]++
		}
	}

	return nil
}

/*
 * Simple graphics state stack
 */
type graphicState struct {
	colorspaceStroke unipdf.PdfColorspace
	colorspaceFill   unipdf.PdfColorspace
}

func (gs graphicState) String() string {
	return fmt.Sprintf("{colorspaceStroke=%T colorspaceFill=%T}",
		gs.colorspaceStroke, gs.colorspaceFill)
}

type graphicStateStack []graphicState

func (gsStack *graphicStateStack) push(gs graphicState) {
	*gsStack = append(*gsStack, gs)
	unicommon.Log.Debug("gsStack=%d", len(*gsStack))
}

func (gsStack *graphicStateStack) pop() graphicState {
	unicommon.Log.Debug("gsStack=%d", len(*gsStack))
	gs := (*gsStack)[len(*gsStack)-1]
	*gsStack = (*gsStack)[:len(*gsStack)-1]
	return gs
}

func isUpper(s string) bool {
	return strings.ToUpper(s) == s
}

var (
	opGrayOp = map[string]string{
		"rg": "g",
		"RG": "G",
		"k":  "g",
		"K":  "G",
	}
	opColorspace = map[string]unipdf.PdfColorspace{
		"g":  unipdf.NewPdfColorspaceDeviceGray(),
		"G":  unipdf.NewPdfColorspaceDeviceGray(),
		"rg": unipdf.NewPdfColorspaceDeviceRGB(),
		"RG": unipdf.NewPdfColorspaceDeviceRGB(),
		"k":  unipdf.NewPdfColorspaceDeviceCMYK(),
		"K":  unipdf.NewPdfColorspaceDeviceCMYK(),
	}
)

func removeOps(pOperations *[]*unicontent.ContentStreamOperation, badOps map[string]bool, removeAll bool) {
	filtered := []*unicontent.ContentStreamOperation{}
	if !removeAll {
		for _, op := range *pOperations {
			if !badOps[op.Operand] {
				filtered = append(filtered, op)
			}
		}
	}
	*pOperations = filtered
}

// totalCounts returns the keys of map `counts` sorted by count
func totalCounts(counts map[string]int) (total int) {
	for _, n := range counts {
		total += n
	}
	return
}

// sortCounts returns the keys of map `counts` sorted by count
func sortCounts(counts map[string]int) []string {
	wordCounts = counts
	keys := []string{}
	for k := range wordCounts {
		keys = append(keys, k)
	}
	sort.Sort(byCount(keys))
	return keys
}

var wordCounts map[string]int

// byCount sorts slices of string by their wordCount
type byCount []string

func (x byCount) Len() int { return len(x) }

func (x byCount) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

func (x byCount) Less(i, j int) bool {
	d := wordCounts[x[i]] - wordCounts[x[j]]
	if d != 0 {
		return d > 0
	}
	li, lj := strings.ToLower(x[i]), strings.ToLower(x[j])
	if li != lj {
		return li < lj
	}
	return x[i] < x[j]
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
		unicommon.Log.Error("modifyPath: Cannot modify path to itself. inputPath=%#q outputDir=%#q",
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

// pdfsEqual compares PDF files `path1` and `path2`
// `threshold` holds the threshold values for image comparison
// If `grayscale` is true then the comparison is on grayscale and color hue may differ
// If `keep` is true then the comparison rasters are retained
// Returns equal, bad1, err where
// 	equal: PDF files are equal
//  bad1: error is in path1
//  err: the error
func pdfsEqual(path1, path2 string, threshold imageThreshold,
	temp string, grayscale, keep bool) (bool, bool, error) {
	dir1 := filepath.Join(temp, "1")
	dir2 := filepath.Join(temp, "2")
	err := os.MkdirAll(dir1, 0777)
	if err != nil {
		panic(err)
	}
	if !keep {
		defer removeDir(dir1)
	}
	err = os.MkdirAll(dir2, 0777)
	if err != nil {
		panic(err)
	}
	if !keep {
		defer removeDir(dir2)
	}

	err = runGhostscript(path1, dir1, grayscale)
	if err != nil {
		return false, true, err
	}
	err = runGhostscript(path2, dir2, grayscale)
	if err != nil {
		return false, false, err
	}

	equal, err := directoriesEqual("*.png", dir1, dir2, threshold)
	return equal, false, nil
}

var (
	gsImageFormat  = "doc-%03d.png"
	gsImagePattern = `doc-(\d+).png$`
	gsImageRegex   = regexp.MustCompile(gsImagePattern)
)

// runGhostscript runs Ghostscript on file `pdf` to create file one png file per page in directory
// `outputDir`
func runGhostscript(pdf, outputDir string, grayscale bool) error {
	unicommon.Log.Debug("runGhostscript: pdf=%#q outputDir=%#q", pdf, outputDir)
	outputPath := filepath.Join(outputDir, gsImageFormat)
	output := fmt.Sprintf("-sOutputFile=%s", outputPath)
	pngDevices := map[bool]string{
		false: "png16m",
		true:  "pnggray",
	}
	cmd := exec.Command(
		ghostscriptName(),
		"-dSAFER",
		"-dBATCH",
		"-dNOPAUSE",
		"-r150",
		fmt.Sprintf("-sDEVICE=%s", pngDevices[grayscale]),
		"-dTextAlphaBits=1",
		"-dGraphicsAlphaBits=1",
		output,
		pdf)
	unicommon.Log.Debug("runGhostscript: cmd=%#q", cmd.Args)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		unicommon.Log.Error("runGhostscript: Could not process pdf=%q err=%v\nstdout=\n%s\nstderr=\n%s\n",
			pdf, err, stdout, stderr)
	}
	return err
}

// ghostscriptName returns the name of the Ghostscipt binary on this OS
func ghostscriptName() string {
	if runtime.GOOS == "windows" {
		return "gswin64c.exe"
	}
	return "gs"
}

// directoriesEqual compares image files that match `mask` in directories `dir1` and `dir2` and
// returns true if they are the same within `threshold`.
func directoriesEqual(mask, dir1, dir2 string, threshold imageThreshold) (bool, error) {
	pattern1 := filepath.Join(dir1, mask)
	pattern2 := filepath.Join(dir2, mask)
	files1, err := filepath.Glob(pattern1)
	if err != nil {
		panic(err)
	}
	files2, err := filepath.Glob(pattern2)
	if err != nil {
		panic(err)
	}
	if len(files1) != len(files2) {
		return false, nil
	}
	n := len(files1)
	for i := 0; i < n; i++ {
		equal, err := filesEqual(files1[i], files2[i], threshold)
		if !equal || err != nil {
			return equal, err
		}
	}
	return true, nil
}

// filesEqual compares files `path1` and `path2` and returns true if they are the same within
// `threshold`
func filesEqual(path1, path2 string, threshold imageThreshold) (bool, error) {
	equal, err := filesBinaryEqual(path1, path2)
	if equal || err != nil {
		return equal, err
	}
	return imagesEqual(path1, path2, threshold)
}

// filesBinaryEqual compares files `path1` and `path2` and returns true if they are identical.
func filesBinaryEqual(path1, path2 string) (bool, error) {
	f1, err := ioutil.ReadFile(path1)
	if err != nil {
		panic(err)
	}
	f2, err := ioutil.ReadFile(path2)
	if err != nil {
		panic(err)
	}
	return bytes.Equal(f1, f2), nil
}

// imagesEqual compares files `path1` and `path2` and returns true if they are the same within
// `threshold`
func imagesEqual(path1, path2 string, threshold imageThreshold) (bool, error) {
	img1, err := readImage(path1)
	if err != nil {
		return false, err
	}
	img2, err := readImage(path2)
	if err != nil {
		return false, err
	}

	w1, h1 := img1.Bounds().Max.X, img1.Bounds().Max.Y
	w2, h2 := img2.Bounds().Max.X, img2.Bounds().Max.Y
	if w1 != w2 || h1 != h2 {
		unicommon.Log.Error("compareImages: Different dimensions. img1=%dx%d img2=%dx%d",
			w1, h1, w2, h2)
		return false, nil
	}

	// `different` contains the grayscale distance (scale 0...255) between pixels in img1 and
	// img2 for pixels that differ between the two images
	different := []float64{}
	for x := 0; x < w1; x++ {
		for y := 0; y < h1; y++ {
			r1, g1, b1, _ := img1.At(x, y).RGBA()
			r2, g2, b2, _ := img2.At(x, y).RGBA()
			if r1 != r2 || g1 != g2 || b1 != b2 {
				d1, d2, d3 := float64(r1)-float64(r2), float64(g1)-float64(g2), float64(b1)-float64(b2)
				// Euclidean distance between pixels in rgb space with scale 0..0xffff
				distance := math.Sqrt(d1*d1 + d2*d2 + d3*d3)
				// Convert scale to 0..0xff and take average of r,g,b values to get grayscale value
				distance = distance / float64(0xffff) * float64(0xff) / 3.0
				different = append(different, distance)
			}
		}
	}
	if len(different) == 0 {
		return true, nil
	}

	fracPixels := float64(len(different)) / float64(w1*h1)
	mean := meanFloatSlice(different)
	equal := fracPixels <= threshold.fracPixels && mean <= threshold.mean

	n := len(different)
	if n > 10 {
		n = 10
	}
	unicommon.Log.Error("compareImages: Different pixels. different=%d/(%dx%d)=%e mean=%.1f %.0f",
		len(different), w1, h1, fracPixels, mean, different[:n])

	return equal, nil
}

// meanFloatSlice returns the mean of the elements of `vals`
func meanFloatSlice(vals []float64) float64 {
	if len(vals) == 0 {
		return 0.0
	}
	var total float64 = 0.0
	for _, v := range vals {
		total += v
	}
	return total / float64(len(vals))
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

	err = runGhostscript(path, dir, false)
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
		unicommon.Log.Error("isColorDirectory: Glob failed. pattern=%#q err=%v", pattern, err)
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
		unicommon.Log.Error("isColorDirectory: Glob failed. pattern=%#q err=%v", pattern, err)
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
		// unicommon.Log.Error("isColorDirectory:  path=%#q", path)
		isColor, err := isColorImage(path, keep)
		// unicommon.Log.Error("isColorDirectory: isColor=%t path=%#q", isColor, path)
		if err != nil {
			panic(err)
			return colorPages, err
		}
		if isColor {
			colorPages = append(colorPages, pageNum)
			// unicommon.Log.Error("isColorDirectory: colorPages=%d %d", len(colorPages), colorPages)
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
		unicommon.Log.Error("markedPath=%#q %s", markedPath, summary)
		err = writeImage(markedPath, markedImg)
	}
	return isColor, err
}

const colorThreshold = 5.0

// imgIsColor returns true if image `img` contains color
func imgIsColor(img image.Image) bool {
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rg := int(r) - int(g)
			gb := int(g) - int(b)
			if rg < 0 {
				rg = -rg
			}
			if gb < 0 {
				gb = -gb
			}
			rgb := float64(rg+gb) / float64(0xFFFF) * float64(0xFF)
			if rgb > colorThreshold {
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
			r, g, b, _ := imgIn.At(x, y).RGBA()
			rg := int(r) - int(g)
			gb := int(g) - int(b)
			if rg < 0 {
				rg = -rg
			}
			if gb < 0 {
				gb = -gb
			}
			rgb := float64(rg+gb) / float64(0xFFFF) * float64(0xFF)
			if rgb > colorThreshold {
				img.Set(x, y, black)
				data = append(data, rgb)
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

// // isColorImage returns true if image file `path` contains color
// func showColorImage(path string) (bool, error) {
// 	img, err := readImage(path)
// 	if err != nil {
// 		return false, err
// 	}

// 	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
// 	for x := 0; x < w; x++ {
// 		for y := 0; y < h; y++ {
// 			r, g, b, _ := img.At(x, y).RGBA()
// 			if r != g || g != b {
// 				return true, nil
// 			}
// 		}
// 	}
// 	return false, nil
// }

// readImage reads image file `path` and returns its contents as an Image.
func readImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		unicommon.Log.Error("readImage: Could not open file. path=%#q err=%v", path, err)
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	return img, err
}

// readImage reads image file `path` and returns its contents as an Image.
func writeImage(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		unicommon.Log.Error("writeImage: Could not create file. path=%#q err=%v", path, err)
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
			unicommon.Log.Error("patternsToPaths: Glob failed. pattern=%#q err=%v", pattern, err)
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

type statistics struct {
	enabled        bool
	testResultPath string
	imageInfoPath  string
	testResultList []testResult
	testResultMap  map[string]int
}

func (s *statistics) load() error {
	if !s.enabled {
		return nil
	}
	s.testResultList = []testResult{}
	s.testResultMap = map[string]int{}

	testResultList, err := testResultRead(s.testResultPath)
	if err != nil {
		return err
	}
	for _, e := range testResultList {
		s.addTestResult(e, true)
	}

	return nil
}

func (s *statistics) _save() error {
	if !s.enabled {
		return nil
	}
	return testResultWrite(s.testResultPath, s.testResultList)
}

func (s *statistics) addTestResult(e testResult, force bool) {
	if !s.enabled {
		return
	}
	i, ok := s.testResultMap[e.name]
	if !ok {
		s.testResultList = append(s.testResultList, e)
		s.testResultMap[e.name] = len(s.testResultList) - 1
	} else {
		s.testResultList[i] = e
	}
	if force {
		s._save()
	}
}

type testResult struct {
	name     string
	colorIn  bool
	colorOut bool
	numPages int
	duration float64
	xobjImg  int
	xobjForm int
}

func testResultRead(path string) ([]testResult, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []testResult{}, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)

	results := []testResult{}
	for i := 0; ; i++ {
		row, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			unicommon.Log.Error("testResultRead: i=%d err=%v", i, err)
			return results, err
		}
		if i == 0 {
			continue
		}
		e := testResult{
			name:     row[0],
			colorIn:  toBool(row[1]),
			colorOut: toBool(row[2]),
			numPages: toInt(row[3]),
			duration: toFloat(row[4]),
			xobjImg:  toInt(row[5]),
			xobjForm: toInt(row[6]),
		}
		results = append(results, e)
	}
	return results, nil
}

func testResultWrite(path string, results []testResult) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)

	if err := w.Write([]string{"name", "colorIn", "colorOut", "numPages", "duration",
		"imageXobj", "formXobj"}); err != nil {
		return err
	}
	for i, e := range results {
		row := []string{
			e.name,
			fmt.Sprintf("%t", e.colorIn),
			fmt.Sprintf("%t", e.colorOut),
			fmt.Sprintf("%d", e.numPages),
			fmt.Sprintf("%.3f", e.duration),
			fmt.Sprintf("%d", e.xobjImg),
			fmt.Sprintf("%d", e.xobjForm),
		}
		if err := w.Write(row); err != nil {
			unicommon.Log.Error("testResultWrite: Error writing record. i=%d path=%#q err=%v",
				i, path, err)
		}
	}

	w.Flush()
	return w.Error()
}

func toBool(s string) bool {
	return strings.ToLower(strings.TrimSpace(s)) == "true"
}

func toInt(s string) int {
	s = strings.TrimSpace(s)
	x, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return x
}

func toFloat(s string) float64 {
	s = strings.TrimSpace(s)
	x, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return x
}

// getOpCounts returns a map of operand: number of occurrences of operand in `operations`
func getOpCounts(operations []*unicontent.ContentStreamOperation) map[string]int {
	opCounts := map[string]int{}
	for _, op := range operations {
		opCounts[op.Operand]++
		allOpCounts[op.Operand]++
	}
	return opCounts
}

func printOperations(description string, operations []*unicontent.ContentStreamOperation) {
	fmt.Printf("%d operations --------^^^--------%#q\n", len(operations), description)
	for i, op := range operations {
		fmt.Printf("%8d: %s\n", i, op)
	}
}

// printOpCounts prints `opCounts` in descending order of occurrences of operand
func printOpCounts(description string, opCounts map[string]int) {
	fmt.Printf("%d ops -------------------------%#q\n", len(opCounts), description)
	for i, k := range sortCounts(opCounts) {
		fmt.Printf("\t%3d: %#6q %5d\n", i, k, opCounts[k])
	}
}

// printCsCounts prints `csCounts` in descending order of occurrences of color spaces
func printCsCounts(description string, csCounts map[string]int) {
	fmt.Printf("%d Color Spaces -------------------------%#q\n", len(csCounts), description)
	for i, k := range sortCounts(csCounts) {
		fmt.Printf("\t%3d: %#6q %5d\n", i, k, csCounts[k])
	}
}
