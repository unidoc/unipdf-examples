/*
 * List all content streams for all pages in a pdf file.
 *
 * Run as: go run pdf_print_content_streams.go input.pdf
 */

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

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
	debug := false
	outputDir := ""
	var minSize int64 = -1
	var maxSize int64 = -1
	flag.BoolVar(&debug, "d", false, "Enable debug logging")
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
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err = os.MkdirAll(outputDir, 0777)
	if err != nil {
		unicommon.Log.Error("MkdirAll failed. outputDir=%#q err=%v", outputDir, err)
	}

	fmt.Printf("pdfList=%d %#q\n", len(pdfList), pdfList)

	for idx, inputPath := range sortFiles(pdfList, minSize, maxSize) {
		_, name := filepath.Split(inputPath)
		inputSize := fileSize(inputPath)
		fmt.Fprintf(os.Stderr, "inputPath %3d of %d %#-30q  (%6d->", idx,
			len(pdfList), name, inputSize)
		outputPath := modifyPath(inputPath, outputDir)
		err = transformContentStreams(inputPath, outputPath)

		outputSize := fileSize(outputPath)
		fmt.Fprintf(os.Stderr, "%6d %3d%%) => %#q\n",
			outputSize, int(float64(outputSize)/float64(inputSize)*100.0+0.5), outputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "transformContentStreams failed. err=%v\n", err)
			os.Exit(1)
		}

		equal, err := pdfsEqual(inputPath, outputPath, "compare.pdfs")
		if !equal || err != nil {
			unicommon.Log.Error("Transform has changed PDF. inputPath=%q outputPath=%q",
				inputPath, outputPath)
			return
		}
	}

	fmt.Printf("%d files\n", len(pdfList))
	fmt.Printf("%d ops ---------^^^----------------\n", len(allOpCounts))
	for i, k := range sortCounts(allOpCounts) {
		fmt.Printf("%3d: %6q %5d\n", i, k, allOpCounts[k])
	}
}

var allOpCounts = map[string]int{}

func transformContentStreams(inputPath, outputPath string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	defer f.Close()

	pdfReader, err := unipdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return err
	}

	if isEncrypted {
		_, err = pdfReader.Decrypt([]byte(""))
		if err != nil {
			return err
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	fmt.Println("--------------------")
	fmt.Println("Content streams:")
	fmt.Println("--------------------")

	opCounts := map[string]int{}

	pdfWriter := unipdf.NewPdfWriter()

	for i := 0; i < numPages; i++ {
		pageNum := i + 1
		page := pdfReader.PageList[i]

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
			// fmt.Printf("Operation %d: %s - Params: %+v\n", idx+1, op.Operand, op.Params)
			opCounts[op.Operand]++
			allOpCounts[op.Operand]++
			s := op.DefaultWriteString()
			// s = fmt.Sprintf("\t%-40s %% command %d of %d [%d] %s\n", s, i, len(operations),
			// 	len(op.Params), op.Operand)

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

		err = pdfWriter.AddPage(page)
		if err != nil {
			return err
		}
	}

	pageDict, err := pdfWriter.Pages()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d pages -----------@@@-----------\n", len(*pageDict))
	fmt.Printf("%s\n", *pageDict)

	fWrite, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer fWrite.Close()

	err = pdfWriter.Write(fWrite)
	if err != nil {
		return err
	}

	fmt.Printf("%d ops -------------------------\n", len(opCounts))
	for i, k := range sortCounts(opCounts) {
		fmt.Printf("%3d: %6q %5d\n", i, k, opCounts[k])
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
	name = fmt.Sprintf("%08d_%s", fileSize(inputPath), name)

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
	fdList = fdList[i0:i1]

	outList := make([]string, n)
	for i, fd := range fdList {
		outList[i] = fd.path
	}

	return outList
}

type FileData struct {
	path string
	os.FileInfo
}

// byName sorts slices of PdfObjectName. It is needed because sort.Strings(keys) gives a typecheck
// error, which I find strange because a PdfObjectName is a string.
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

func pdfsEqual(path1, path2, temp string) (bool, error) {
	dir1 := filepath.Join(temp, "1")
	dir2 := filepath.Join(temp, "2")
	err := os.MkdirAll(dir1, 0777)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(dir2, 0777)
	if err != nil {
		panic(err)
	}
	err = runGhostscript(path1, dir1)
	if err != nil {
		panic(err)
	}
	err = runGhostscript(path2, dir2)
	if err != nil {
		panic(err)
	}
	equal, err := directoriesEqual(dir1, dir2, temp)
	if err != nil {
		panic(err)
	}
	os.RemoveAll(dir1)
	os.RemoveAll(dir2)
	return equal, nil
}

// runPs2Pdf runs Ghostscript  on file `pdf` to create file one png file per page in directory
// `outputDir`
func runGhostscript(pdf, outputDir string) error {
	unicommon.Log.Info("runGhostscript: pdf=%#q outputDir=%#q", pdf, outputDir)
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
	unicommon.Log.Info("runGhostscript: cmd=%#q", cmd.Args)

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
