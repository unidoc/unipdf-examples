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
 * Run as: ./pdf_transform_content_streams -o output [-d] testdata/*.pdf > blah
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
 *
 *	Currently failing files in PETER's corpus of 332 PDF files
 *		Radon_Transform.pdf
 *		ESCP-R reference_151008.pdf
 *		lda.pdf
 *		BLUEBOOK.pdf
 *		pdf_hacks.pdf

 27 fail with latest v2 code
	  0 `lda.pdf`
	  1 `power_law_bins.pdf`
	  2 `rules_of_ml.pdf`
	  3 `compression.kdd06(1).pdf`
	  4 `climate_science_words.pdf`
	  5 `joi160132.pdf`
	  6 `1560401.pdf`
	  7 `twitter.pdf`
	  8 `pearson_science_8_sb_chapter_5_unit_5.2.pdf`
	  9 `1512.03547v2.pdf`
	 10 `BLUEBOOK.pdf`
	 11 `cvxopt_1306.0057v1.pdf`
	 12 `scan_alan_2016-03-30-10-38-15.pdf`
	 13 `a0w20000000dikuAAA.pdf`
	 14 `2013-12-12_rg_final_report.pdf`
	 15 `Parsing-Probabilistic.pdf`
	 16 `PhysRevLett.118.060401.pdf`
	 17 `art%3A10.1186%2Fs13673-015-0039-9.pdf`
	 18 `dark-internet-mail-environment-march-2015.pdf`
	 19 `2015-09-16-T23-39-51_ec2-user_ip-172-31-6-72_jim.pdf`
	 20 `Hierarchical Detection of Hard Exudates.pdf`
	 21 `WhatIsEnergy.pdf`
	 22 `Physics_Sample_Chapter_3.pdf`
	 23 `day3_TemporalImageProcessing.pdf`
	 24 `Lesson_054_handout.pdf`
	 25 `talk_Simons_part1_pdf.pdf`
	 26 `nips-tutorial-policy-optimization-Schulman-Abbeel.pdf`
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
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	unicommon "github.com/unidoc/unidoc/common"
	// unilicense "github.com/unidoc/unidoc/license"
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
	unicommon.DebugOutput = debug
	unicommon.SetLogger(unicommon.ConsoleLogger{})

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
	enabled:        false,
	testResultPath: "xform.test.results.csv",
	imageInfoPath:  "xform.image.info.csv",
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
	unipdf.ValidatingOperations = true

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
	defer testStats.save()

	for idx, inputPath := range pdfList {

		testStats.save()
		_, name := filepath.Split(inputPath)
		inputSize := fileSize(inputPath)

		fmt.Fprintf(os.Stderr, "%3d of %d %#-30q  (%6d->", idx,
			len(pdfList), name, inputSize)
		outputPath := modifyPath(inputPath, outputDir)

		t0 := time.Now()
		numPages, err := transformPdfFile(inputPath, outputPath, noContentTransforms, doGrayscaleTransform)
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
			e := testResult{
				name:     path.Base(inputPath),
				colorIn:  isColorIn,
				colorOut: isColorOut,
				numPages: numPages,
				duration: dt.Seconds(),
			}
			testStats.addTestResult(e, true)

			if err != nil || isColorOut {
				if err != nil {
					unicommon.Log.Error("Transform has damaged PDF. err=%v\n\tinputPath=%#q\n\toutputPath=%#q",
						err, inputPath, outputPath)
				} else {
					unicommon.Log.Error("Output PDF is color.\n\tinputPath=%#q\n\toutputPath=%#q",
						inputPath, outputPath)
					unicommon.Log.Error("isPdfColor: %d Color pages", len(colorPagesOut))
					for _, p := range colorPagesOut {
						unicommon.Log.Error("isPdfColor: page %d", p+1)
					}
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

	fmt.Fprintf(os.Stderr, "%d files %d bad %d failed\n",
		len(pdfList), len(badFiles), len(failFiles))
	fmt.Fprintf(os.Stderr, "%d bad\n", len(badFiles))
	for i, path := range badFiles {
		fmt.Fprintf(os.Stderr, "%3d %#q\n", i, path)
	}
	fmt.Fprintf(os.Stderr, "%d fail\n", len(failFiles))
	for i, path := range failFiles {
		fmt.Fprintf(os.Stderr, "%3d %#q\n", i, path)
	}

	printOpCounts("all operation in PDF", allOpCounts)
}

// transformPdfFile transforms PDF `inputPath` and writes the resulting PDF to `outputPath`
// If `noContentTransforms` is true then stream contents are not parsed
func transformPdfFile(inputPath, outputPath string, noContentTransforms, doGrayscaleTransform bool) (int, error) {
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
		unicommon.Log.Debug("^^^^ page %d=%s", pageNum, page.String())

		if !noContentTransforms {
			err = transformPageContents(page, inputPath, pageNum, doGrayscaleTransform)
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

	// unicommon.Log.Error("allAlternateColorSpaces: %d", len(allAlternateColorSpaces))
	// for k, v := range allAlternateColorSpaces {
	// 	fmt.Fprintf(os.Stderr, "%#15q: %d\n", k, v)
	// }

	return numPages, nil
}

// transformPageContents
//	- parses the contents of streams in `page` into a slice of operations
//	- converts the slice of operations into a stream
//	- replaces the streams in `page` with the new stream

func transformPageContents(page *unipdf.PdfPage, inputPath string, pageNum int,
	doGrayscaleTransform bool) error {
	cstream, err := page.GetAllContentStreams()
	if err != nil {
		return err
	}

	desc := fmt.Sprintf("%s:page%d", filepath.Base(inputPath), pageNum)
	cstreamOut, err := transformString(page, cstream, desc, doGrayscaleTransform)
	if err != nil {
		return err
	}
	return page.SetContentStreams([]string{cstreamOut}, nil)
}

func transformString(page *unipdf.PdfPage, cstream, desc string, doGrayscaleTransform bool) (string,
	error) {
	unicommon.Log.Info("transformString: desc=%s", desc)

	operations, err := parseStreamContents(cstream, desc)
	if err != nil {
		return "", err
	}

	if doGrayscaleTransform {
		printOperations(fmt.Sprintf("%s: before", desc), operations)
		printOpCounts(fmt.Sprintf("%s: before", desc), getOpCounts(operations))
		if err := transformColorToGrayscale(page, desc, &operations); err != nil {
			return "", err
		}
		printOpCounts(fmt.Sprintf("%s: after ", desc), getOpCounts(operations))
	}

	opStrings := []string{}
	for _, op := range operations {
		opStrings = append(opStrings, op.DefaultWriteString())
	}
	cstreamOut := strings.Join(opStrings, " ")

	return cstreamOut, nil
}

// parsePageContents parses the contents of streams in `page` and returns them as a slice of operations
func parseStreamContents(cstream string, desc string) ([]*unipdf.ContentStreamOperation, error) {
	cstreamParser := unipdf.NewContentStreamParser(cstream)
	unicommon.Log.Debug("transformContentStream: %s cstream=\n'%s'\nXXXXXX", desc, cstream)
	return cstreamParser.Parse()
}

// // writePageContents converts `operations` into a stream and replaces the streams in `page` with it
// func writePageContents(page *unipdf.PdfPage, pageNum int,
// 	operations []*unipdf.ContentStreamOperation) error {

// 	opStrings := []string{}
// 	for _, op := range operations {
// 		opStrings = append(opStrings, op.DefaultWriteString())
// 	}
// 	cstreamOut := strings.Join(opStrings, " ")

// 	return page.SetContentStreams([]string{cstreamOut}, nil)
// }

var allAlternateColorSpaces = map[string]int{}

func transformStreamGrayscale(page *unipdf.PdfPage, cstream, desc string) (string, error) {
	unicommon.Log.Info("transformStreamGrayscale: desc=%s", desc)
	return transformString(page, cstream, desc, true)
}

var stackDepth = 0

// transformColorToGrayscale transforms color pages to grayscale
func transformColorToGrayscale(page *unipdf.PdfPage, desc_ string,
	pOperations *[]*unipdf.ContentStreamOperation) (err error) {

	stackDepth++
	defer func() { stackDepth-- }()
	rubric := fmt.Sprintf("transformColorToGrayscale: <%d> %s.", stackDepth, desc_)
	unicommon.Log.Info("****%s", rubric)

	xobjImgs := []string{}
	xobjCSs := []string{}

	currentColorSpace := "DeviceGray"
	currentSeparate := false

	for _, op := range *pOperations {
		h := fmt.Sprintf("%s op=%s", rubric, op)
		var vals []float64
		switch op.Operand {
		case "cs", "CS":
			name, err := op.GetNameParam()
			if err != nil {
				return err
			}
			xobjCSs = append(xobjCSs, name)

			_, currentColorSpace, currentSeparate, _ = page.GetColorSpace(name)
			unicommon.Log.Info("%s currentColorSpace=%#q currentSeparate=%t",
				h, currentColorSpace, currentSeparate)
			switch currentColorSpace {
			case "DeviceRGB", "DeviceCMYK":
				unicommon.Log.Info("### %s: name=%#q", h, name)
				if err = op.SetNameParam("DeviceGray"); err != nil {
					return err
				}
				unicommon.Log.Info("##@ %s: ", op)
			}
		case "Do":
			name, err := op.GetNameParam()
			if err != nil {
				return err
			}
			xobjImgs = append(xobjImgs, name)
		case "sc", "SC", "scn", "SCN":
			unicommon.Log.Info("#@: %s currentColorSpace=%#q currentSeparate=%t",
				h, currentColorSpace, currentSeparate)
			if currentSeparate {
				if vals, err = op.GetFloatParams(1); err != nil {
					return err
				}
				unicommon.Log.Info("#!# %s: vals=%v", h, vals)
				if err = op.SetOpFloatParams(op.Operand, []float64{1.0 - vals[0]}); err != nil {
					return err
				}
			} else {
				switch currentColorSpace {
				case "DeviceRGB":
					if vals, err = op.GetFloatParams(3); err != nil {
						return err
					}
					unicommon.Log.Info("#!# %s: vals=%v", h, vals)
					if err = op.SetOpFloatParams(op.Operand, []float64{rgbToGray(vals)}); err != nil {
						return err
					}
				case "DeviceCMKY":
					if vals, err = op.GetFloatParams(4); err != nil {
						return err
					}
					unicommon.Log.Info("#!^ %s: vals=%v", h, vals)
					if err = op.SetOpFloatParams(op.Operand, []float64{cmykToGray(vals)}); err != nil {
						return err
					}
				}
			}
		case "rg", "RG":
			if vals, err = op.GetFloatParams(3); err != nil {
				return err
			}
			if err = op.SetOpFloatParams(grayOp[op.Operand], []float64{rgbToGray(vals)}); err != nil {
				return err
			}
			unicommon.Log.Info("#!# %s: vals=%v", h, vals)
		case "k", "K":
			if vals, err = op.GetFloatParams(4); err != nil {
				return err
			}
			if err = op.SetOpFloatParams(grayOp[op.Operand], []float64{cmykToGray(vals)}); err != nil {
				return err
			}
			unicommon.Log.Info("#!# %s: vals=%v", h, vals)
		}
	}

	for _, name := range xobjCSs {
		colorSpace, alternate, _, err := page.GetColorSpace(name)
		if err != nil {
			return err
		}
		unicommon.Log.Debug("%s colorspace: %#q=%s %s", rubric, name, alternate, colorSpace)
		_, ok := allAlternateColorSpaces[alternate]
		if !ok {
			allAlternateColorSpaces[alternate] = 1
		} else {
			allAlternateColorSpaces[alternate]++
		}
	}

	for _, name := range xobjImgs {
		/*imageSummary*/ _, err := page.ConvertXObjectToGray(name, desc_, transformStreamGrayscale)
		if err != nil {
			return err
		}
		// if imageSummary != nil {
		// 	e := imageInfo{
		// 		fileName:   path.Base(inputPath),
		// 		pageNum:    pageNum,
		// 		xobjName:   name,
		// 		length:     imageSummary.Length,
		// 		filter1:    imageSummary.Filter1,
		// 		filter2:    imageSummary.Filter2,
		// 		colorSpace: imageSummary.ColorSpace,
		// 	}
		// 	testStats.addImageInfo(e, true)
		// }
	}
	return nil
}

var grayOp = map[string]string{
	"rg": "g",
	"RG": "G",
	"k":  "g",
	"K":  "G",
}

// rgbToGray returns the grayscale equivalent of the r,g,b values in `vals`
func rgbToGray(vals []float64) float64 {
	r, g, b := vals[0], vals[1], vals[2]
	return (r + g + b) / 3.0
}

// cmykToGray returns the grayscale equivalent of the c,m,y,k values in `vals`
func cmykToGray(vals []float64) float64 {
	c, m, y, k := vals[0], vals[1], vals[2], vals[3]
	a := 1.0 - k - (c+m+y)/3.0
	if a < 0.0 {
		a = 0.0
	}
	return a
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

// runGhostscript runs Ghostscript on file `pdf` to create file one png file per page in directory
// `outputDir`
func runGhostscript(pdf, outputDir string, grayscale bool) error {
	unicommon.Log.Debug("runGhostscript: pdf=%#q outputDir=%#q", pdf, outputDir)
	outputPath := filepath.Join(outputDir, "doc-%03d.png")
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

	if !showPages {
		colorPages, err := colorDirectoryPages("*.png", dir)
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
		isColor, err := isColorImage(path)
		if isColor || err != nil {
			return isColor, err
		}
	}
	return false, nil
}

// isColorDirectory returns lists pages  the image files that match `mask` in directories `dir`
// that have any color pixels. It
func colorDirectoryPages(mask, dir string) ([]int, error) {
	pattern := filepath.Join(dir, mask)
	files, err := filepath.Glob(pattern)
	if err != nil {
		unicommon.Log.Error("isColorDirectory: Glob failed. pattern=%#q err=%v", pattern, err)
		return nil, err
	}

	colorPages := []int{}
	for i, path := range files {
		color, err := isColorImage(path)
		if err != nil {
			return colorPages, err
		}
		if color {
			colorPages = append(colorPages, i)
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
	isColor := imgIsColor(img)
	if isColor {
		markedPath := fmt.Sprintf("%s.marked.png", path)
		markedImg := imgMarkColor(img)
		if err := writeImage(markedPath, markedImg); err != nil {
			return isColor, err
		}
	}
	return isColor, nil
}

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
			if rg+gb > 4*255 {
				// if r != g || g != b {
				return true
			}
		}
	}
	return false
}

func imgMarkColor(imgIn image.Image) image.Image {
	img := image.NewNRGBA(imgIn.Bounds())
	black := color.RGBA{0, 0, 0, 255}
	// white := color.RGBA{255, 255, 255, 255}
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b, _ := imgIn.At(x, y).RGBA()
			if r != g || g != b {
				img.Set(x, y, black)
				// unicommon.Log.Error("^^^ (%d, %d) %d %d %d", x, y, r, g, b)
			}
		}
	}
	return img
}

// isColorImage returns true if image file `path` contains color
func showColorImage(path string) (bool, error) {
	img, err := readImage(path)
	if err != nil {
		return false, err
	}

	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if r != g || g != b {
				return true, nil
			}
		}
	}
	return false, nil
}

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
	imageInfoList  []imageInfo
	testResultMap  map[string]int
	imageInfoMap   map[string]int
}

func (s *statistics) load() error {
	if !s.enabled {
		return nil
	}
	s.testResultList = []testResult{}
	s.testResultMap = map[string]int{}
	s.imageInfoList = []imageInfo{}
	s.imageInfoMap = map[string]int{}

	testResultList, err := testResultRead(s.testResultPath)
	if err != nil {
		return err
	}
	for _, e := range testResultList {
		s.addTestResult(e, true)
	}

	imageInfoList, err := imageInfoRead(s.imageInfoPath)
	if err != nil {
		return err
	}
	for _, e := range imageInfoList {
		s.addImageInfo(e, true)
	}

	return nil
}

func (s *statistics) save() error {
	if !s.enabled {
		return nil
	}
	err1 := testResultWrite(s.testResultPath, s.testResultList)
	err2 := imageInfoWrite(s.imageInfoPath, s.imageInfoList)
	if err1 != nil {
		return err1
	}
	return err2
}

func (s *statistics) addTestResult(e testResult, force bool) {
	if !s.enabled {
		return
	}
	i, ok := s.testResultMap[e.name]
	if !ok {
		s.testResultList = append(s.testResultList, e)
		s.testResultMap[e.name] = len(s.testResultList) - 1
	} else if force {
		s.testResultList[i] = e
	}
}

func (s *statistics) addImageInfo(e imageInfo, force bool) {
	if !s.enabled {
		return
	}
	k := imageInfoKey(e)
	i, ok := s.imageInfoMap[k]
	if !ok {
		s.imageInfoList = append(s.imageInfoList, e)
		s.imageInfoMap[k] = i
	} else if force {
		s.imageInfoList[i] = e
	}
}

func imageInfoKey(e imageInfo) string {
	return fmt.Sprintf("%s:page%d:%s", e.fileName, e.pageNum, e.xobjName)
}

type testResult struct {
	name     string
	colorIn  bool
	colorOut bool
	numPages int
	duration float64
}

type imageInfo struct {
	fileName   string
	pageNum    int
	xobjName   string
	length     int
	filter1    string
	filter2    string
	colorSpace string
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

	if err := w.Write([]string{"name", "colorIn", "colorOut", "numPages", "duration"}); err != nil {
		return err
	}
	for i, e := range results {
		row := []string{
			e.name,
			fmt.Sprintf("%t", e.colorIn),
			fmt.Sprintf("%t", e.colorOut),
			fmt.Sprintf("%d", e.numPages),
			fmt.Sprintf("%.3f", e.duration),
		}
		if err := w.Write(row); err != nil {
			unicommon.Log.Error("testResultWrite: Error writing record. i=%d path=%#q err=%v",
				i, path, err)
		}
	}

	w.Flush()
	return w.Error()
}

func imageInfoRead(path string) ([]imageInfo, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []imageInfo{}, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)

	results := []imageInfo{}
	for i := 0; ; i++ {
		row, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			unicommon.Log.Error("imageInfoRead: i=%d err=%v", i, err)
			return results, err
		}
		if i == 0 {
			continue
		}
		e := imageInfo{
			fileName:   row[0],
			pageNum:    toInt(row[1]),
			xobjName:   row[2],
			length:     toInt(row[3]),
			filter1:    row[4],
			filter2:    row[5],
			colorSpace: row[6],
		}
		results = append(results, e)
	}
	return results, nil
}

func imageInfoWrite(path string, results []imageInfo) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)

	if err := w.Write([]string{"fileName", "pageNum", "xobjName", "length",
		"filter1", "filter2", "colorSpace"}); err != nil {
		return err
	}
	for i, e := range results {
		if len(e.xobjName) == 0 {
			panic("YYYY")
		}
		row := []string{
			e.fileName,
			fmt.Sprintf("%d", e.pageNum),
			e.xobjName,
			fmt.Sprintf("%d", e.length),
			e.filter1,
			e.filter2,
			e.colorSpace,
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
func getOpCounts(operations []*unipdf.ContentStreamOperation) map[string]int {
	opCounts := map[string]int{}
	for _, op := range operations {
		opCounts[op.Operand]++
		allOpCounts[op.Operand]++
	}
	return opCounts
}

func printOperations(description string, operations []*unipdf.ContentStreamOperation) {
	fmt.Printf("%d operations --------^^^--------%#q\n", len(operations), description)
	for _, op := range operations {
		fmt.Printf("%s\n", op)
	}
}

// printOpCounts prints `opCounts` in descending order of occurrences of operand
func printOpCounts(description string, opCounts map[string]int) {
	fmt.Printf("%d ops -------------------------%#q\n", len(opCounts), description)
	for i, k := range sortCounts(opCounts) {
		fmt.Printf("%3d: %#6q %5d\n", i, k, opCounts[k])
	}
}
