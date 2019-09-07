/*
 * Passthrough benchmark for UniPDF: loads input PDF files and writes them back out.
 *
 * Run as: go run pdf_passthrough_bench.go ...
 *
 * This will perform the passthrough benchmark on all pdf files and write results to stdout.
 *
 * See the other command line options in the top of main()
 *     -o <processPath> - Temporary output file path (default /tmp/test.pdf)
 *     -odir <outputdir> - Output directory path (Optional, overrides -o)
 *     -d: Debug level logging
 *     -a: Run all tests. Don't stop at first failure (This flag isn't used here)
 *     -gsv: Validate with ghostscript
 *     -hang: Hang when completed (no exit) - for memory profiling
 *     -rmlist: Print out a list of files to rm to make fully compliant
 *     -optimize: Use Use Pdf compression and optimization
 *     -pprof: Run with profiling enabled.
 *     -lazy: Use lazy loading.
 *
 * The passthrough benchmark
 * - Loads the input PDF with unipdf
 * - Writes the output PDF
 * - Runs ghostscript (gs) on both input and output and checks for errors
 * - Invalid if unipdf returns an error or if gs on output has more errors than gs on input PDF.
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
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/model/optimize"
)

// Results for single pdf.
type benchmarkResult struct {
	path         string
	passed       bool
	processTime  float64
	sizeMB       float64
	outputSizeMB float64
	errorMessage string
	rmList       bool
}

// Total results.
type benchmarkResults []benchmarkResult

const usage = `Usage:
pdf_passthrough_bench [options] <file1> <file2> ... > results
Options:
-o <processPath> - Temporary output file path (default /tmp/test.pdf)
-odir <outputdir> - Output directory path (Optional, overrides -o)
-d: Debug level logging
-gsv: Validate with ghostscript
-hang: Hang when completed (no exit) - for memory profiling
-rmlist: Print out a list of files to rm to make fully compliant
-optimize: Use Pdf compression and optimization
-lazy: Use lazy loading.

Example: pdf_passthrough_bench -gsv ~/pdfdb/* >results_YYYY_MM_DD
`

type benchParams struct {
	runAllTests  bool
	processPath  string
	outputDir    string
	gsValidation bool
	hangOnExit   bool
	printRmList  bool
	optimize     bool
	lazyLoading  bool
	profilePath  string
	loglevel     string
}

func main() {
	params := benchParams{}

	fmt.Printf("UniPDF version %s\n", common.Version)

	//params.debug = false       // Write debug level info to stdout?
	params.runAllTests = false // Don't stop when a PDF file fails to process?
	params.processPath = ""    // Transformed PDFs are written here
	params.outputDir = ""      // Alternatively, can store output files in an output directory.
	params.gsValidation = false
	params.hangOnExit = false
	params.printRmList = false
	params.optimize = false
	params.loglevel = "info"

	flag.StringVar(&params.loglevel, "loglevel", "info", "Set loglevel: info (default), debug, trace, none")
	flag.BoolVar(&params.gsValidation, "gsv", false, "Enable ghostscript validation")
	flag.BoolVar(&params.runAllTests, "a", false, "Run all tests. Don't stop at first failure")
	flag.BoolVar(&params.hangOnExit, "hang", false, "Hang when completed without exiting (memory profiling)")
	flag.BoolVar(&params.printRmList, "rmlist", false, "Print rm list at end")
	flag.StringVar(&params.processPath, "o", "/tmp/test.pdf", "Temporary output file path")
	flag.StringVar(&params.outputDir, "odir", "", "Output directory (optional)")
	flag.BoolVar(&params.optimize, "optimize", false, "Use Pdf compression and optimization")
	flag.BoolVar(&params.lazyLoading, "lazy", false, "Use lazy loading")
	flag.StringVar(&params.profilePath, "pprof", "", "Pprof output for profiling (optional)")

	flag.Parse()
	args := flag.Args()
	if len(args) < 1 || (len(params.processPath) == 0 && len(params.outputDir) == 0) {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}

	if len(params.profilePath) > 0 {
		fmt.Printf("Profiling to %s\n", params.profilePath)
		f, err := os.Create(params.profilePath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		err = pprof.StartCPUProfile(f)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		defer pprof.StopCPUProfile()
	}

	fmt.Printf("With GS validation: %t\n", params.gsValidation)
	fmt.Printf("With compression and optimization: %t\n", params.optimize)

	switch params.loglevel {
	case "none":
		common.SetLogger(common.DummyLogger{})
	case "info":
		common.SetLogger(common.NewConsoleLogger(common.LogLevelInfo))
	case "debug":
		common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
	case "trace":
		common.SetLogger(common.NewConsoleLogger(common.LogLevelTrace))
	default:
		fmt.Printf("Unknown loglevel: %v\n", params.loglevel)
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

	fmt.Printf("Complete\n")

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
			fi, err := os.Stat(path)
			if err != nil {
				panic(err)
			}

			// One level of directories supported.
			if fi.Mode().IsDir() {
				innerFiles, err := ioutil.ReadDir(path)
				if err != nil {
					panic(err)
				}
				for _, f := range innerFiles {
					if f.Mode().IsRegular() {
						pathList = append(pathList, filepath.Join(path, f.Name()))
					}
				}
				continue
			}
			if !fi.Mode().IsRegular() {
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
		cmd = exec.Command(ghostscriptName(), "-dBATCH", "-dNODISPLAY", "-dNOPAUSE", "-dPDFSTOPONERROR", option, path)
	} else {
		cmd = exec.Command(ghostscriptName(), "-dBATCH", "-dNODISPLAY", "-dNOPAUSE", "-dPDFSTOPONERROR", path)
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
func testPassthroughSinglePdf(inputPath string, params benchParams) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var reader *model.PdfReader
	if params.lazyLoading {
		fmt.Printf("Lazy loading used\n")
		reader, err = model.NewPdfReaderLazy(file)
	} else {
		reader, err = model.NewPdfReader(file)
	}
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

	writer := model.NewPdfWriter()
	if params.optimize {
		optimizer := optimize.New(optimize.Options{
			CombineDuplicateStreams:         true,
			CombineDuplicateDirectObjects:   true,
			ImageUpperPPI:                   100.0,
			ImageQuality:                    90,
			CombineIdenticalIndirectObjects: true,
			UseObjectStreams:                true,
			CompressStreams:                 true,
		})
		writer.SetOptimizer(optimizer)
	}
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

	// By default: uses processPath, unless if output dir is specified, uses the basename of `path` and outputs
	// to the output dir.
	outputPath := params.processPath
	if len(params.outputDir) > 0 {
		outputPath = filepath.Join(params.outputDir, filepath.Base(inputPath))
	}

	common.Log.Debug("Write the file")
	of, err := os.Create(outputPath)
	if err != nil {
		common.Log.Debug("Failed to create file (%s)", err)
		return err
	}
	defer of.Close()

	err = writer.Write(of)
	if err != nil {
		common.Log.Debug("WriteFile error")

		return err
	}

	// GS validation of input, output pdfs.
	if params.gsValidation {
		common.Log.Debug("Validating input file")
		_, inputWarnings := validatePdf(inputPath, "")
		common.Log.Debug("Validating output file")

		err, warnings := validatePdf(outputPath, "")
		if err != nil && warnings > inputWarnings {
			common.Log.Error("Input warnings %d vs output %d", inputWarnings, warnings)
			return fmt.Errorf("GS: Invalid PDF input %d/ output %d warnings", inputWarnings, warnings)
		}
		common.Log.Debug("Valid PDF!")
	}

	return nil
}

// TestSinglePdf tests a single pdf file.
func TestSinglePdf(target string, params benchParams) error {
	err := testPassthroughSinglePdf(target, params)

	if err != nil {
		// If unidoc fails the file, check the input file.  Do not count as error
		// if ghostscript has issues with the file.  Ensure not processed by ghostscript already (GS: error prefix).
		if params.gsValidation && !strings.HasPrefix(err.Error(), "GS: ") {
			fmt.Printf("Error, lets doa ghostscript check of input\n")
			err, _ := validatePdf(target, "")
			fmt.Println(err)
			if err != nil {
				common.Log.Debug("GS fails processing input file %s - not counting as problematic", target)
				return nil
			}
		}
	}

	return err
}

// printResults prints a summary of the benchmark results.
func (this benchmarkResults) printResults(params benchParams) {
	succeeded := 0
	total := 0
	totalInputSize := 0.0
	totalOutputSize := 0.0
	var totalTime float64 = 0.0

	for _, result := range this {
		if result.passed {
			succeeded++
			totalTime += result.processTime
		}
		total++
		totalInputSize += result.sizeMB
		totalOutputSize += result.outputSizeMB

		baseName := filepath.Base(result.path)
		if len(baseName) > 30 {
			baseName = baseName[0:30]
		}
		if !result.passed {
			// Only print ones that failed.
			fmt.Printf("%30s\t%.1f\t%v\t%.1f\tError: %s\n", baseName, result.sizeMB,
				result.passed, result.processTime, result.errorMessage)
		} else if params.optimize && result.outputSizeMB > result.sizeMB {
			fmt.Printf("%30s\t%.1fM\t%v\t%.1fs\t%.4f -> %.4f\t%.4f\n", baseName, result.sizeMB,
				result.passed, result.processTime, result.sizeMB, result.outputSizeMB, result.sizeMB/result.outputSizeMB)
		}
	}

	fmt.Printf("----------------------\n")
	fmt.Printf("Total: %d\n", total)
	fmt.Printf("Successes: %d\n", succeeded)
	fmt.Printf("Failed: %d\n", total-succeeded)
	fmt.Printf("Total time: %.1f secs (%.2f per file)\n", totalTime, totalTime/float64(succeeded))
	if params.optimize {
		fmt.Printf("Total input files size: %.3f MB\nTotal output files size: %.3f MB\nTotal compression ratio: %.3f\n",
			totalInputSize, totalOutputSize, totalInputSize/totalOutputSize)
	}

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

// getFileSize returns the file size in MB.
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

// isDirectory returns true if the path refers to a directory.
func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir(), err
}

// benchmarkPDFs runs a benchmark on a list of PDF files specified by path with a specified set of parameters.
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
		if benchmark.passed && params.optimize {
			outputPath := params.processPath
			fmt.Printf("Output path: %s\n", outputPath)
			if len(params.outputDir) > 0 {
				outputPath = filepath.Join(params.outputDir, filepath.Base(path))
			}
			fmt.Printf("Output path 2: %s\n", outputPath)
			outputFileSizeMB, err := getFileSize(outputPath)
			if err != nil {
				fmt.Printf("Error in getFileSize: %v\n", err)
				return err
			}
			benchmark.outputSizeMB = outputFileSizeMB
			fmt.Printf("Input file size: %.3f MB\nOutput file size: %.3f MB\nCompression ratio: %.3f\n",
				benchmark.sizeMB, benchmark.outputSizeMB, benchmark.sizeMB/benchmark.outputSizeMB)
		}
		benchmarkResults = append(benchmarkResults, benchmark)
	}

	benchmarkResults.printResults(params)

	return nil
}
