/*
 * Combthrough benchmark combs through all indirect objects in input PDF files. Error if fails to process.
 *
 * Run as: go run pdf_combthrough_bench.go ...
 *
 * This will perform the combthrough benchmark on all pdf files and write results to stdout.
 *
 * See the other command line options in the top of main()
 *      -d: Debug level logging
 *      -a: Keep processing PDF files after failures
 *      -min <val>: Minimum PDF file size to test
 *      -max <val>: Maximum PDF file size to test
 *      -r <name>: Name of results file
 *
 * The combthrough benchmark
 * - Loads the input PDF with unidoc
 * - Combs through all indirect objects, decoding streams etc. (same as pdf_all_objects.go).
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/unidoc/unidoc/pdf/core"

	common "github.com/unidoc/unidoc/common"
	unipdf "github.com/unidoc/unidoc/pdf/model"
)

// Results for single pdf.
type benchmarkResult struct {
	path         string
	passed       bool
	processTime  float64
	sizeMB       float64
	errorMessage string
	rmList       bool
}

// Total results.
type benchmarkResults []benchmarkResult

func initUniDoc(debug bool) error {
	if debug {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelTrace))
		//common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
	} else {
		common.SetLogger(common.DummyLogger{})
	}

	return nil
}

const usage = `Usage:
pdf_combthrough_bench [options] <file1> <file2> ... > results
Options:
-d: Debug level logging
-hang: Hang when completed (no exit) - for memory profiling
-rmlist: Print out a list of files to rm to make fully compliant
-password: Specify a password for opening PDF

Example: pdf_combthrough_bench ~/pdfdb/* >results_YYYY_MM_DD
`

type benchParams struct {
	debug       bool
	runAllTests bool
	processPath string
	hangOnExit  bool
	printRmList bool
	password    string
}

func main() {
	params := benchParams{}

	params.debug = false       // Write debug level info to stdout?
	params.runAllTests = false // Don't stop when a PDF file fails to process?
	params.processPath = ""    // Transformed PDFs are written here
	params.hangOnExit = false
	params.printRmList = false

	flag.BoolVar(&params.debug, "d", false, "Enable debug logging")
	flag.BoolVar(&params.runAllTests, "a", false, "Run all tests. Don't stop at first failure")
	flag.BoolVar(&params.hangOnExit, "hang", false, "Hang when completed without exiting (memory profiling)")
	flag.BoolVar(&params.printRmList, "rmlist", false, "Print rm list at end")
	flag.StringVar(&params.processPath, "o", "/tmp/test.pdf", "Temporary output file path")
	flag.StringVar(&params.password, "password", "", "PDF Password (empty default)")

	flag.Parse()
	args := flag.Args()
	if len(args) < 1 || len(params.processPath) == 0 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}

	err := initUniDoc(params.debug)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	pdfList, err := patternsToPaths(args)
	if err != nil {
		common.Log.Error("patternsToPaths failed. args=%#q err=%v", args, err)
		os.Exit(1)
	}

	err = benchmarkPDFs(pdfList, params)
	if err != nil {
		common.Log.Error("benchmarkPDFs failed err=%v", err)
		os.Exit(1)
	}

	if params.hangOnExit {
		// Endless loop.
		for {
			time.Sleep(100 * time.Millisecond)
		}
	}
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

// testCombSinglePdf tests loading a pdf file, and combs through all the indirect objects/streams
// and decodes the streams. Returns an error on fail.
func testCombSinglePdf(path string, params benchParams) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader, err := unipdf.NewPdfReader(file)
	if err != nil {
		common.Log.Debug("Reader create error %s\n", err)
		return err
	}

	isEncrypted, err := reader.IsEncrypted()
	if err != nil {
		return err
	}

	if isEncrypted {
		valid, err := reader.Decrypt([]byte(params.password))
		if err != nil {
			common.Log.Debug("Fail to decrypt: %v", err)
			return err
		}

		if !valid {
			return fmt.Errorf("Unable to access, encrypted")
		}
	}

	numPages, err := reader.GetNumPages()
	if err != nil {
		common.Log.Debug("Failed to get number of pages")
		return err
	}

	if numPages < 1 {
		common.Log.Debug("Empty pdf - nothing to be done!")
		return nil
	}

	objNums := reader.GetObjectNums()

	// Output.
	for _, objNum := range objNums {
		obj, err := reader.GetIndirectObjectByNumber(objNum)
		if err != nil {
			return err
		}
		if stream, is := obj.(*core.PdfObjectStream); is {
			_, err := core.DecodeStream(stream)
			if err != nil {
				return err
			}
		} else if indObj, is := obj.(*core.PdfIndirectObject); is {
			if indObj.PdfObject == nil {
				return fmt.Errorf("indirect object with no object")
			}
		}
	}

	return nil
}

// Test a single pdf file.
func TestSinglePdf(target string, params benchParams) error {
	err := testCombSinglePdf(target, params)
	return err
}

// Print the summary of the benchmark results.
func (br benchmarkResults) printResults(params benchParams) {
	succeeded := 0
	total := 0
	var totalTime float64 = 0.0

	for _, result := range br {
		if result.passed {
			succeeded++
			totalTime += result.processTime
		}
		total++

		if !result.passed {
			// Only print ones that failed.
			fmt.Printf("%s\t%.1f\t%v\t%.1f\t%s\n", result.path, result.sizeMB,
				result.passed, result.processTime, result.errorMessage)
		}
	}

	// Frequency of errors.
	errmap := map[string]int{}
	for _, result := range br {
		if !result.passed {
			errmap[result.errorMessage]++
		}
	}
	keys := []string{}
	for errmsg := range errmap {
		keys = append(keys, errmsg)
	}
	sort.Slice(keys, func(i, j int) bool {
		return errmap[keys[i]] >= errmap[keys[j]]
	})
	for _, key := range keys {
		fmt.Printf("'%s' - %d files\n", key, errmap[key])
	}

	fmt.Printf("----------------------\n")
	fmt.Printf("Total: %d\n", total)
	fmt.Printf("Successes: %d\n", succeeded)
	fmt.Printf("Failed: %d\n", total-succeeded)
	fmt.Printf("Total time: %.1f secs (%.2f per file)\n", totalTime, totalTime/float64(succeeded))

	// Print list to remove
	if params.printRmList {
		for _, result := range br {
			if !result.passed {
				// Only print ones that failed.
				fmt.Printf("rm \"%s\"\n", result.path)
			}

		}
	}
}

// Get file size in MB.
func getFileSize(path string) (float64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	finfo, err := file.Stat()
	if err != nil {
		return 0, err
	}

	sizeMB := float64(finfo.Size()) / 1024 / 1024
	return sizeMB, nil
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir(), err
}

func benchmarkPDFs(paths []string, params benchParams) error {
	benchmarkResults := benchmarkResults{}

	for _, path := range paths {
		benchmark := benchmarkResult{}
		benchmark.path = path

		fileSizeMB, err := getFileSize(path)
		if err != nil {
			return err
		}
		benchmark.sizeMB = fileSizeMB

		fmt.Printf("Testing %s\n", path)
		start := time.Now()
		err = TestSinglePdf(path, params)
		elapsed := time.Since(start)
		benchmark.processTime = elapsed.Seconds()
		if err == nil {
			benchmark.passed = true
			fmt.Printf("%s - pass\n", path)
		} else {
			benchmark.passed = false
			benchmark.errorMessage = fmt.Sprintf("%s", err)
			fmt.Printf("%s - fail %s\n", path, err)
		}

		benchmarkResults = append(benchmarkResults, benchmark)
	}

	benchmarkResults.printResults(params)

	return nil
}
