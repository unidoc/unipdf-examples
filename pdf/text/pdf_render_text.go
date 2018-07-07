/*
 * PDF to text: Extract all text for each page of a pdf file.
 *
 * Run as: go run pdf_extract_text.go testdata/*.pdf
 */

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/unidoc/unidoc/common"
	"github.com/unidoc/unidoc/pdf/extractor"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

const (
	usage        = "Usage: go run pdf_render_text.go testdata/*.pdf\n"
	badFilesPath = "bad.files"
)

func main() {
	// Make sure to enter a valid license key.
	// Otherwise text is truncated and a watermark added to the text.
	// License keys are available via: https://unidoc.io
	/*
			license.SetLicenseKey(`
		-----BEGIN UNIDOC LICENSE KEY-----
		...key contents...
		-----END UNIDOC LICENSE KEY-----
		`)
	*/
	var showHelp, debug, trace bool
	var filesPath string
	flag.BoolVar(&showHelp, "h", false, "Show this help message.")
	flag.BoolVar(&debug, "d", false, "Print debugging information.")
	flag.BoolVar(&trace, "e", false, "Print detailed debugging information.")
	flag.StringVar(&filesPath, "@", "",
		"File containing list of files to process. Usually a 'bad.files' from a previous test run.")
	makeUsage(usage)

	flag.Parse()
	args := flag.Args()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}
	if len(args) < 1 && filesPath == "" {
		flag.Usage()
		os.Exit(1)
	}
	if trace {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelTrace))
	} else if debug {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
	}

	files := args[:]
	sort.Strings(files)
	if filesPath != "" {
		if filesPath == badFilesPath {
			fmt.Fprintf(os.Stderr, "Setting files to %s will overwrite %s. Try a different name",
				badFilesPath, badFilesPath)
			os.Exit(1)
		}
		var err error
		files, err = filesFromPreviousRun(filesPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse %q. err=%v", filesPath, err)
			os.Exit(1)
		}
	}

	fBad, err := os.OpenFile(badFilesPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create %s. err=%v", badFilesPath, err)
		os.Exit(1)
	}
	defer fBad.Close()
	fmt.Fprintln(fBad, "========================= ^^^ =========================")

	errorCounts := map[error]int{}

	for i, inputPath := range files {
		fmt.Println("========================= ^^^ =========================")
		pdfReader, numPages, err := getReader(inputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\to====> Pdf File %3d of %d %q err=%v\n",
				i+1, len(files), inputPath, err)
			continue
		}

		version := pdfReader.PdfVersion()
		// We are currently not interested in old PDF files. If you are, comment out these lines.
		if version == "1.0" || version == "1.1" || version == "1.2" {
			continue
		}

		fmt.Fprintf(os.Stderr, "Pdf File %3d of %d (%3s) %3d pages %q \n",
			i+1, len(files), pdfReader.PdfVersion(), numPages, inputPath)

		err = outputPdfText(inputPath, pdfReader, numPages)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\tx====> Pdf File %3d of %d %q err=%v\n",
				i+1, len(files), inputPath, err)
			fmt.Fprintf(fBad, "%q version=%s pages=%d err=%v\n", inputPath, version, numPages, err)
		}
		fmt.Println("========================= ~~~ =========================")
		if err != nil {
			errorCounts[err]++
		}
	}
	fmt.Fprintf(os.Stderr, "Done %d files \n", len(files))
	if len(errorCounts) > 0 {
		fmt.Fprintln(os.Stderr, "=== Error counts ===")
		for err, n := range errorCounts {
			fmt.Fprintf(os.Stderr, "%-30s %d (%.0f%%)\n", err, n, 100.0*float64(n)/float64(len(files)))
		}
	}
}

// "/Users/pcadmin/testdata/seqdec120813.pdf\" version=1.3 pages=45 err=Could not decode text"
var reFilename = regexp.MustCompile(`^\s*"(.+?)"\s*version=([\d\.]+)\s*(?:pages=(\d+)\s*)?err=(.+?)\s*$`)

type testResult struct {
	filename string
	version  string
	pages    int
	err      string
}

var testHeader = []string{"filename", "version", "pages", "error"}

func (r *testResult) String() string {
	return fmt.Sprintf("version=%s pages=%3d err=%-20s %q",
		r.version, r.pages, r.err, r.filename)
}

func (r *testResult) asStrings() []string {
	return []string{
		r.filename,
		r.version,
		fmt.Sprintf("%d", r.pages),
		r.err,
	}
}

var ignoredErrors = map[string]bool{
	"Unsupported font": true,
}

// filesFromPreviousRun returns the files that failed in a previous run.
func filesFromPreviousRun(filename string) ([]string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ReadFile failed. filename=%q\n", filename)
		return nil, err
	}
	data := string(b)
	lines := strings.Split(data, "\n")
	fileResult := map[string]testResult{}
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		groups := reFilename.FindStringSubmatch(line)
		if groups == nil {
			fmt.Fprintf(os.Stderr, "Bad line %d in %q: line=%q\n", i, filename, line)
			continue
		}
		pages := -1
		if groups[3] != "" {
			pages, err = strconv.Atoi(groups[3])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Bad line %d in %q: line=%q\n", i, filename, line)
				continue
			}
		}
		r := testResult{
			filename: groups[1],
			version:  groups[2],
			pages:    pages,
			err:      groups[4],
		}
		if _, err := os.Stat(r.filename); err != nil {
			fmt.Fprintf(os.Stderr, "Non-existant i=%d.\n\tgroups=%+v\n\tline=%q\n", i, groups, line)
			return nil, err
		}
		fileResult[r.filename] = r
	}

	results := []testResult{}
	for _, r := range fileResult {
		results = append(results, r)
	}

	sort.Slice(results, func(i, j int) bool {
		ri, rj := results[i], results[j]
		if ri.err != rj.err {
			return ri.err < rj.err
		}
		if ri.pages != rj.pages {
			return ri.pages < rj.pages
		}
		if ri.version != rj.version {
			return ri.version < rj.version
		}
		return ri.filename < rj.filename
	})
	fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^   Results   ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
	for i, r := range results {
		fmt.Printf("%3d of %d: %s\n", i+1, len(results), r.String())
	}
	fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ End Results ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")

	err = saveAsCsv(filename+".csv", results)
	if err != nil {
		return nil, err
	}

	files := []string{}
	for _, r := range results {
		if _, ok := ignoredErrors[r.err]; ok {
			continue
		}
		files = append(files, r.filename)
	}
	return files, nil
}

func saveAsCsv(filename string, results []testResult) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	err = writer.Write(testHeader)
	if err != nil {
		return err
	}
	for _, r := range results {
		err := writer.Write(r.asStrings())
		if err != nil {
			return err
		}
	}
	return nil
}

// getReader returns a PdfReader and the number of pages for PDF file `inputPath`.
func getReader(inputPath string) (pdfReader *pdf.PdfReader, numPages int, err error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return
	}
	defer f.Close()

	pdfReader, err = pdf.NewPdfReader(f)
	if err != nil {
		return
	}
	numPages, err = pdfReader.GetNumPages()
	return
}

// outputPdfText prints out text of PDF file `inputPath` to stdout.
// `pdfReader` is a previously opened PdfReader of `inputPath`
func outputPdfText(inputPath string, pdfReader *pdf.PdfReader, numPages int) error {
	for pageNum := 1; pageNum <= numPages; pageNum++ {

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return err
		}
		ex, err := extractor.New(page)
		if err != nil {
			return err
		}
		text, err := ex.ExtractText()
		if err != nil {
			return err
		}

		fmt.Printf("Page %d of %d: %q\n", pageNum, numPages, inputPath)
		fmt.Printf("%s\n", text)
		fmt.Println("------------------------- ... -------------------------")
	}

	return nil
}

func makeUsage(msg string) {
	usage := flag.Usage
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, msg)
		usage()
	}
}
