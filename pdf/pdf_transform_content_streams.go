/*
 * "Transform" all content streams in all pages in a list of pdf files.

 * "Transform" is in quotes because this program doesn't perform any transforms. It converts PDF
 *   files into our internal representation then converts the internal representation back to a PDF
 *   file and checks that the output PDF file is the same as the input PDF file.
 *
 * Run as: ./pdf_transform_content_streams -o output testdata/*.pdf > blah
 *
 *   This will transform all .pdfs file in testdata and write the results to output.
 *   The main results are written to stderr so you will see them in your console.
 *	 Detailed information is written to stdout and you will see them in blah.
 *
 *  See the other command line options in the top of main()
 *		-a tests all the input files. The default behaviour is stop at the first failure. Use this
 *			to find out how many of your corpus files this program works for.
 *		-x will transform without parsing content streams. Use this to see which failures are due to
 *			problems in the content parsing code.
 *			Running -a then -a -x will tell you how well this code is performing on your corpus
 *			and which failures are due to content parsing.
 *
 *
 *	Currently failing files in PETER's corpus of 332 PDF files
 *		Radon_Transform.pdf
 *		ESCP-R reference_151008.pdf
 *		lda.pdf
 *		BLUEBOOK.pdf
 *		pdf_hacks.pdf
 */

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
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

func main() {
	debug := false               // Write debug level info to stdout?
	noContentTransforms := false // Don't parse stream contents?
	runAllTests := false         // Don't stop when a PDF file fails to process?
	outputDir := ""              // Transformed PDFs are written here
	var minSize int64 = -1       // Minimum size for an input PDF to be processed.
	var maxSize int64 = -1       // Maximum size for an input PDF to be processed.
	flag.BoolVar(&debug, "d", false, "Enable debug logging")
	flag.BoolVar(&noContentTransforms, "x", false, "Don't transform streams")
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

	pdfList := []string{}
	for _, a := range args {
		files, err := filepath.Glob(a)
		if err != nil {
			panic(err)
		}
		pdfList = append(pdfList, files...)
	}

	err := initUniDoc("", debug)
	if err != nil {
		os.Exit(1)
	}

	err = os.MkdirAll(outputDir, 0777)
	if err != nil {
		unicommon.Log.Error("MkdirAll failed. outputDir=%#q err=%v", outputDir, err)
		os.Exit(1)
	}

	compDir := makeUniqueDir("compare.pdfs")
	fmt.Fprintf(os.Stderr, "compDir=%#q\n", compDir)
	defer os.RemoveAll(compDir)

	pdfList = sortFiles(pdfList, minSize, maxSize)
	fmt.Printf("pdfList=%d %#q\n", len(pdfList), pdfList)
	badFiles := []string{}
	failFiles := []string{}

	for idx, inputPath := range pdfList {
		_, name := filepath.Split(inputPath)
		inputSize := fileSize(inputPath)
		fmt.Fprintf(os.Stderr, "inputPath %3d of %d %#-30q  (%6d->", idx,
			len(pdfList), name, inputSize)
		outputPath := modifyPath(inputPath, outputDir)

		t0 := time.Now()
		numPages, err := transformContentStreams(inputPath, outputPath, noContentTransforms)
		dt := time.Since(t0)
		if err != nil {
			unicommon.Log.Error("transformContentStreams failed. err=%v", err)
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

		equal, badInput, err := pdfsEqual(inputPath, outputPath, compDir)
		if badInput {
			unicommon.Log.Error("Bad input PDF. inputPath=%#q err=%v", inputPath, err)
			badFiles = append(badFiles, inputPath)
			continue
		}
		if !equal || err != nil {
			unicommon.Log.Error("Transform has changed PDF. inputPath=%q outputPath=%q",
				inputPath, outputPath)
			failFiles = append(failFiles, inputPath)
			if runAllTests {
				continue
			}
			os.Exit(1)
		}
	}

	fmt.Fprintf(os.Stderr, "%d files %d bad %d failed\n",
		len(pdfList), len(badFiles), len(failFiles))
	fmt.Fprintf(os.Stderr, "%d bad\n", len(badFiles))
	for i, path := range badFiles {
		fmt.Fprintf(os.Stderr, "%d %#q\n", i, path)
	}
	fmt.Fprintf(os.Stderr, "%d fail\n", len(failFiles))
	for i, path := range failFiles {
		fmt.Fprintf(os.Stderr, "f%d %#q\n", i, path)
	}

	fmt.Printf("%d ops ---------^^^----------------\n", len(allOpCounts))
	for i, k := range sortCounts(allOpCounts) {
		fmt.Printf("%3d: %6#q %5d\n", i, k, allOpCounts[k])
	}

}

var allOpCounts = map[string]int{}

func transformContentStreams(inputPath, outputPath string, noContentTransforms bool) (int, error) {
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

		if !noContentTransforms {
			err = transformPageContents(page, pageNum)
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

func transformPageContents(page *unipdf.PdfPage, pageNum int) error {
	// fmt.Fprintln(os.Stderr, "$$$")

	opCounts := map[string]int{}

	cstream, err := page.GetAllContentStreams()
	if err != nil {
		return err
	}

	cstreamParser := unipdf.NewContentStreamParser(cstream, page)
	unicommon.Log.Debug("transformContentStream: pageNum=%d cstream=\n'%s'\nXXXXXX",
		pageNum, cstream)
	operations, err := cstreamParser.Parse()
	if err != nil {
		return err
	}

	opStrings := []string{}
	for _, op := range operations {
		opCounts[op.Operand]++
		allOpCounts[op.Operand]++
		s := op.DefaultWriteString()
		opStrings = append(opStrings, s)
	}

	cstreamOut := strings.Join(opStrings, " ")

	fmt.Printf("Page %d - content stream %d: %d => %d\n", pageNum, len(cstream), len(cstreamOut))

	// for name, ximg := range xobjDict {
	// 	unicommon.Log.Debug("transformContentStreams: name=%#q ximg=%T", name, ximg)
	// 	err = page.AddImageResource(name, ximg)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	err = page.SetContentStreams([]string{cstreamOut})
	if err != nil {
		return err
	}

	fmt.Printf("%d ops -------------------------\n", len(opCounts))
	for i, k := range sortCounts(opCounts) {
		fmt.Printf("%3d: %6#q %5d\n", i, k, opCounts[k])
	}
	return nil
}

func fileSize(path string) int64 {
	fi, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	return fi.Size()
}

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
	fmt.Printf("fdList=%d\n", len(fdList))

	fdList = fdList[i0:i1]

	fmt.Printf("fdList=%d\n", len(fdList))

	outList := make([]string, len(fdList))
	for i, fd := range fdList {
		outList[i] = fd.path
	}
	fmt.Printf("outList=%d\n", len(outList))

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

// pdfsEqual compares PDF files `path1` and `path2` and returns equal, bad1, err where
// 	equal: PDF files are equal
//  bad1: error is in path1
//  err: the error

func pdfsEqual(path1, path2, temp string) (bool, bool, error) {
	dir1 := filepath.Join(temp, "1")
	dir2 := filepath.Join(temp, "2")
	err := os.MkdirAll(dir1, 0777)
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir1)
	err = os.MkdirAll(dir2, 0777)
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir2)

	err = runGhostscript(path1, dir1)
	if err != nil {
		return false, true, err
	}
	err = runGhostscript(path2, dir2)
	if err != nil {
		return false, false, err
	}

	equal, err := directoriesEqual(dir1, dir2, temp)
	if err != nil {
		panic(err)
	}

	return equal, false, nil
}

// runPs2Pdf runs Ghostscript  on file `pdf` to create file one png file per page in directory
// `outputDir`
func runGhostscript(pdf, outputDir string) error {
	unicommon.Log.Debug("runGhostscript: pdf=%#q outputDir=%#q", pdf, outputDir)
	outputPath := filepath.Join(outputDir, "doc-%03d.png")
	output := fmt.Sprintf("-sOutputFile=%s", outputPath)
	cmd := exec.Command("gs",
		"-dSAFER",
		"-dBATCH",
		"-dNOPAUSE",
		"-r150",
		"-sDEVICE=pnggray",
		"-dTextAlphaBits=4",
		output,
		pdf)
	unicommon.Log.Debug("runGhostscript: cmd=%#q", cmd.Args)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		unicommon.Log.Error("runGhostscript: pdf=%q err=%v\nstdout=\n%s\nstderr=\n%s\n",
			pdf, err, stdout, stderr)
		panic(err)
	}

	return err
}

// directoriesEqual compares files that match `mask` in directories `dir1` and `dir2` and returns
// true if they are the same.
func directoriesEqual(mask, dir1, dir2 string) (bool, error) {
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
		equal, err := filesEqual(files1[i], files2[i])
		if !equal || err != nil {
			return false, nil
		}
	}
	return true, nil
}

// filesEqual compares files `path1` and `path2` and returns true if they are the same.
func filesEqual(path1, path2 string) (bool, error) {
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
