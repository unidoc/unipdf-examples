/*
 * Detect the number of pages and the color pages (1-offset) all pages in a list of PDF files.
 * Compares these results to running Ghostscript on the PDF files and reports an error if the results don't match.
 *
 * Run as: ./pdf_detect -o output [-d] [-a] testdata/*.pdf > blah
 *
 * The main results are written to stderr so you will see them in your console.
 * Detailed information is written to stdout and you will see them in blah.
 *
 *  See the other command line options in the top of main()
 *      -d Write debug level logs to stdout
 *		-a Tests all the input files. The default behavior is stop at the first failure. Use this
 *			to find out how many of your corpus files this program works for.
 */

package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
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

	common "github.com/unidoc/unidoc/common"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func initUniDoc(debug bool) error {

	pdf.SetPdfCreator("Peter Williams")

	// To make the library log we just have to initialise the logger which satisfies
	// the common.Logger interface, common.DummyLogger is the default and
	// does not do anything. Very easy to implement your own.
	// common.SetLogger(common.DummyLogger{})
	logLevel := common.LogLevelInfo
	if debug {
		logLevel = common.LogLevelDebug
	}
	common.SetLogger(common.ConsoleLogger{LogLevel: logLevel})

	return nil
}

const usage = `Usage:
%s -o <output directory> [-d][-a][-min <val>][-max <val>] <file1> <file2> ...
-d: Debug level logging
-a: Keep converting PDF files after failures
-min: Minimum PDF file size to test
-max: Maximum PDF file size to test
`

func main() {
	debug := false            // Write debug level info to stdout?
	keep := false             // Keep the rasters used for PDF comparison"
	compareGrayscale := false // Do PDF raster comparison on grayscale rasters?
	runAllTests := false      // Don't stop when a PDF file fails to process?
	var minSize int64 = -1    // Minimum size for an input PDF to be processed.
	var maxSize int64 = -1    // Maximum size for an input PDF to be processed.
	flag.BoolVar(&debug, "d", false, "Enable debug logging")
	flag.BoolVar(&compareGrayscale, "g", false, "Do PDF raster comparison on grayscale rasters")
	flag.BoolVar(&runAllTests, "a", false, "Run all tests. Don't stop at first failure")
	flag.Int64Var(&minSize, "min", -1, "Minimum size of files to process (bytes)")
	flag.Int64Var(&maxSize, "max", -1, "Maximum size of files to process (bytes)")

	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		panic("args")
		os.Exit(1)
	}

	err := initUniDoc(debug)
	if err != nil {
		os.Exit(1)
	}
	compDir := makeUniqueDir("compare.pdfs")
	fmt.Fprintf(os.Stderr, "compDir=%#q\n", compDir)
	if !keep {
		defer removeDir(compDir)
	}

	pdfList, err := patternsToPaths(args)
	if err != nil {
		common.Log.Error("patternsToPaths failed. args=%#q err=%v", args, err)
		os.Exit(1)
	}
	pdfList = sortFiles(pdfList, minSize, maxSize)
	badFiles := []string{}
	failFiles := []string{}

	for idx, inputPath := range pdfList {

		_, name := filepath.Split(inputPath)
		inputSize := fileSize(inputPath)

		fmt.Fprintf(os.Stderr, "%3d of %d %#-30q  (%6d)", idx, len(pdfList), name, inputSize)

		t0 := time.Now()
		numPages, colorPages, err := describePdf(inputPath)
		dt := time.Since(t0)
		if err != nil {
			common.Log.Error("describePdf failed. err=%v", err)
			failFiles = append(failFiles, inputPath)
			if runAllTests {
				continue
			}
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, " %d pages %d color %.3f sec\n", numPages, len(colorPages), dt.Seconds())

		colorPagesIn, err := pdfColorPages(inputPath, compDir)

		if err != nil || !equalSlices(colorPagesIn, colorPages) {
			if err != nil {
				common.Log.Error("PDF is damaged. err=%v\n\tinputPath=%#q", err, inputPath)
			} else {
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
			}
			failFiles = append(failFiles, inputPath)
			if runAllTests {
				continue
			}
			os.Exit(1)
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

	fmt.Printf("%d colorspaces\n", len(pdf.AllColorspaces))
	for cs, n := range pdf.AllColorspaces {
		fmt.Printf("%15q: %4d\n", cs, n)
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

// describePdf reads PDF `inputPath` and returns number of pages, slice of color page numbers (1-offset)
func describePdf(inputPath string) (int, []int, error) {

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
		if colored {
			colorPages = append(colorPages, pageNum)
		}
	}

	return numPages, colorPages, nil
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

// colorThreshold is the total difference of r,g,b values for which a pixel is considered to be color
// Color components are in range 0-0xFFFF
// We make this 10x the PDF color threshold as guess
const colorThreshold = pdf.ColorTolerance * float64(0xFFFF) * 10.0

// imgIsColor returns true if image `img` contains color
func imgIsColor(img image.Image) bool {
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			rr, gg, bb, _ := img.At(x, y).RGBA()
			r, g, b := float64(rr), float64(gg), float64(bb)
			if math.Abs(r-g) > colorThreshold || math.Abs(r-b) > colorThreshold || math.Abs(g-b) > colorThreshold {
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

// sliceDiff returns the elements in a that aren't in b
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
