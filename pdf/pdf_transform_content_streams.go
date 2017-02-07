/*
 * List all content streams for all pages in a pdf file.
 *
 * Run as: go run pdf_print_content_streams.go input.pdf
 */

package main

import (
	"flag"
	"fmt"
	"os"
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
	flag.BoolVar(&debug, "d", false, "Enable debug logging")
	flag.StringVar(&outputDir, "o", "", "Output directory")
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

	os.MkdirAll(outputDir, 0777)

	fmt.Printf("pdfList=%d %#q\n", len(pdfList), pdfList)

	for idx, inputPath := range pdfList {
		outputPath := modifyPath(inputPath, outputDir)
		fmt.Fprintf(os.Stderr, "inputPath %3d of %d %#q => %#q\n", idx, len(pdfList),
			inputPath, outputPath)
		err = transformContentStreams(inputPath, outputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "transformContentStreams failed. err=%v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Input : Size %10d: %#q\n", fileSize(inputPath), inputPath)
		fmt.Printf("Output: Size %10d: %#q\n", fileSize(outputPath), outputPath)
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

		contentStreams, err := page.GetContentStreams()
		if err != nil {
			return err
		}
		xobjDict := unipdf.XObjectImageMap{}

		contentStreamsOut := []string{}
		for idx, cstream := range contentStreams {
			unicommon.Log.Debug("transformContentStream=%s\n", cstream)

			cstreamParser := unipdf.NewContentStreamParser(cstream, xobjDict)
			operations, err := cstreamParser.Parse()
			if err != nil {
				return err
			}

			opStrings := []string{}
			for i, op := range operations {
				// fmt.Printf("Operation %d: %s - Params: %+v\n", idx+1, op.Operand, op.Params)
				opCounts[op.Operand]++
				allOpCounts[op.Operand]++
				s := op.DefaultWriteString()
				s = fmt.Sprintf("\t%-40s %% command %d of %d [%d] %s\n", s, i, len(operations),
					len(op.Params), op.Operand)
				opStrings = append(opStrings, s)
			}

			cstreamOut := strings.Join(opStrings, " ")
			contentStreamsOut = append(contentStreamsOut, cstreamOut)

			fmt.Printf("Page %d - content stream %d: %d => %d\n", pageNum, idx+1,
				len(cstream), len(cstreamOut))

		}

		for name, ximg := range xobjDict {
			unicommon.Log.Debug("transformContentStreams: name=%#q ximg=%T", name, ximg)
			err = page.AddImageResource(name, ximg)
			if err != nil {
				return err
			}
		}

		fmt.Printf("Page %d has %d content streams\n", pageNum, len(contentStreamsOut))
		err = page.SetContentStreams(contentStreamsOut)
		if err != nil {
			return err
		}

		err = pdfWriter.AddPage(page)
		if err != nil {
			return err
		}
	}

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
