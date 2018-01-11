/*
 * Passthrough benchmark for UniDoc, loads input PDF files and writes them back out. Has ghostscript validation.
 *
 * Run as: go run pdf_passhtrough_bench.go ...
 *
 * This will perform the passthrough benchmark on all pdf files and write results to stdout.
 *
 * See the other command line options in the top of main()
 *      -o processDir - Temporary processing directory (default compare.pdfs)
 *      -d: Debug level logging
 *      -a: Keep converting PDF files after failures
 *      -min <val>: Minimum PDF file size to test
 *      -max <val>: Maximum PDF file size to test
 *      -r <name>: Name of results file
 *
 * The passthrough benchmark
 * - Loads the input PDF with unidoc
 * - Writes the output PDF
 * - Runs ghostscript (gs) on both input and output and checks for errors
 * - Invalid if unidoc returns an error or if gs on output has more errors than gs on input PDF.
 */

package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

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
pdf_passthrough_bench [options] <file1> <file2> ... > results
Options:
-o <processPath> - Temporary output file path (default /tmp/test.pdf)
-d: Debug level logging
-gsv: Validate with ghostscript
-hang: Hang when completed (no exit) - for memory profiling
-rmlist: Print out a list of files to rm to make fully compliant

Example: pdf_passthrough_bench -gsv ~/pdfdb/* >results_YYYY_MM_DD
`

type benchParams struct {
	debug        bool
	runAllTests  bool
	processPath  string
	gsValidation bool
	hangOnExit   bool
	printRmList  bool
}

func main() {
	params := benchParams{}

	params.debug = false       // Write debug level info to stdout?
	params.runAllTests = false // Don't stop when a PDF file fails to process?
	params.processPath = ""    // Transformed PDFs are written here
	params.gsValidation = false
	params.hangOnExit = false
	params.printRmList = false

	flag.BoolVar(&params.debug, "d", false, "Enable debug logging")
	flag.BoolVar(&params.gsValidation, "gsv", false, "Enable ghostscript validation")
	flag.BoolVar(&params.runAllTests, "a", false, "Run all tests. Don't stop at first failure")
	flag.BoolVar(&params.hangOnExit, "hang", false, "Hang when completed without exiting (memory profiling)")
	flag.BoolVar(&params.printRmList, "rmlist", false, "Print rm list at end")
	flag.StringVar(&params.processPath, "o", "/tmp/test.pdf", "Temporary output file path")

	flag.Parse()
	args := flag.Args()
	if len(args) < 1 || len(params.processPath) == 0 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}

	fmt.Printf("With GS validation: %t\n", params.gsValidation)

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

// validatePdf a pdf file using Ghostscript, returns an error if unable to execute.
// Also returns the number of output warnings, which can be used as some sort of measure
// of validity, especially when comparing with a transformed version of same file.
func validatePdf(path string, password string) (error, int) {
	common.Log.Debug("Validating: %s", path)

	var cmd *exec.Cmd
	if len(password) > 0 {
		option := fmt.Sprintf("-sPDFPassword=%s", password)
		cmd = exec.Command(ghostscriptName(), "-dBATCH", "-dNODISPLAY", "-dNOPAUSE", option, path)
	} else {
		cmd = exec.Command(ghostscriptName(), "-dBATCH", "-dNODISPLAY", "-dNOPAUSE", path)
	}

	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()
	if err != nil {
		common.Log.Debug("%s", out.String())
		common.Log.Debug("%s", errOut.String())
		common.Log.Error("GS failed with error %s", err)
		return fmt.Errorf("GS failed with error (%s)", err), 0
	}

	outputErr := errOut.String()
	warnings := strings.Count(outputErr, "****")
	common.Log.Error(": - %d warnings %s", warnings, outputErr)

	if warnings > 1 {
		if len(outputErr) > 80 {
			outputErr = outputErr[:80] // Trim the output.
		}
		common.Log.Error("Invalid - %d warnings %s", warnings, outputErr)
		return fmt.Errorf("Invalid - %d warnings (%s)", warnings, outputErr), warnings
	}

	// Valid if no error.
	return nil, 0
}

// ghostscriptName returns the name of the Ghostscript binary on this OS
func ghostscriptName() string {
	if runtime.GOOS == "windows" {
		return "gswin64c.exe"
	}
	return "gs"
}

// testPassthroughSinglePdf tests loading a pdf file, and writing it back out (passthrough).
func testPassthroughSinglePdf(path string, params benchParams) error {
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
		valid, err := reader.Decrypt([]byte(""))
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

	writer := unipdf.NewPdfWriter()

	ocProps, err := reader.GetOCProperties()
	if err != nil {
		return err
	}
	writer.SetOCProperties(ocProps)

	for j := 0; j < numPages; j++ {
		page, err := reader.GetPage(j + 1)
		if err != nil {
			common.Log.Debug("Get page error %s", err)
			return err
		}

		// Load and set outlines (table of contents).
		outlineTree := reader.GetOutlineTree()

		err = writer.AddPage(page)
		if err != nil {
			common.Log.Debug("Add page error %s", err)
			return err
		}

		writer.AddOutlineTree(outlineTree)
	}

	// Copy the forms over to the new document also.
	if reader.AcroForm != nil {
		err = writer.SetForms(reader.AcroForm)
		if err != nil {
			common.Log.Debug("Add forms error %s", err)
			return err
		}
	}

	common.Log.Debug("Write the file")
	file, err = os.Create(params.processPath)
	if err != nil {
		common.Log.Debug("Failed to create file (%s)", err)
		return err
	}
	defer file.Close()

	err = writer.Write(file)
	if err != nil {
		common.Log.Debug("WriteFile error")
		return err
	}

	// GS validation of input, output pdfs.
	if params.gsValidation {
		common.Log.Debug("Validating input file")
		_, inputWarnings := validatePdf(path, "")
		common.Log.Debug("Validating output file")

		err, warnings := validatePdf(params.processPath, "")
		if err != nil && warnings > inputWarnings {
			common.Log.Error("Input warnings %d vs output %d", inputWarnings, warnings)
			return fmt.Errorf("Invalid PDF input %d/ output %d warnings", inputWarnings, warnings)
		}
		common.Log.Debug("Valid PDF!")
	}

	return nil
}

// Test a single pdf file.
func TestSinglePdf(target string, params benchParams) error {
	err := testPassthroughSinglePdf(target, params)
	return err
}

// Print the summary of the benchmark results.
func (this benchmarkResults) printResults(params benchParams) {
	succeeded := 0
	total := 0
	var totalTime float64 = 0.0

	for _, result := range this {
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

	fmt.Printf("----------------------\n")
	fmt.Printf("Total: %d\n", total)
	fmt.Printf("Successes: %d\n", succeeded)
	fmt.Printf("Failed: %d\n", total-succeeded)
	fmt.Printf("Total time: %.1f secs (%.2f per file)\n", totalTime, totalTime/float64(succeeded))

	// Print list to remove
	if params.printRmList {

		for _, result := range this {
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
