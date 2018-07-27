/*
 * PDF to text: Extract all text for each page of a pdf file.
 *
 * Run as: go run pdf_extract_text.go testdata/*.pdf
 */

package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/unidoc/unidoc/common"
	"github.com/unidoc/unidoc/pdf/extractor"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

const (
	usage                 = "Usage: go run pdf_render_text.go testdata/*.pdf\n"
	badFilesPath          = "bad.files"
	defaultNormalizeWidth = 60
)

var ErrBadText = errors.New("could not decode text")

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
	var showHelp, debug, trace, verbose bool
	var filesPath string
	var threshold float64
	var width int
	flag.BoolVar(&showHelp, "h", false, "Show this help message.")
	flag.BoolVar(&debug, "d", false, "Print debugging information.")
	flag.BoolVar(&verbose, "v", false, "Print extra page information.")
	flag.BoolVar(&trace, "e", false, "Print detailed debugging information.")
	flag.Float64Var(&threshold, "t", 1.0,
		"Missclassification theshold. percentage of missclassified characters above this "+
			"threshold are treated as errors.")
	flag.IntVar(&width, "w", 1.0,
		"Normalize text (remove runs of space) with this target output width.")
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

	errorCounts := map[string]int{}

	for i, inputPath := range files {
		sizeMB := fileSizeMB(inputPath)
		if verbose {
			fmt.Println("========================= ^^^ =========================")
		}
		t0 := time.Now()
		pdfReader, numPages, err := getReader(inputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\to====> Pdf File %3d of %d %q err=%v\n",
				i+1, len(files), inputPath, err)
			continue
		}

		version := pdfReader.PdfVersion()
		fmt.Fprintf(os.Stderr, "Pdf File %3d of %d (%3s) %4.1f MB %3d pages %q ",
			i+1, len(files), pdfReader.PdfVersion(), sizeMB, numPages, inputPath)

		numChars, numMisses, err := outputPdfText(inputPath, pdfReader, numPages, width, verbose)
		dt := time.Since(t0)
		if err == nil {
			err = missclassificationError(threshold, numChars, numMisses)
		}

		fmt.Fprintf(os.Stderr, "%3.1f sec %d chars %d misses (%.1f%%)\n", dt.Seconds(),
			numChars, numMisses, percentage(numChars, numMisses))
		if err != nil {
			fmt.Fprintf(os.Stderr, "\tx====> Pdf File %3d of %d %q numChars=%d numMisses=%d err=%v\n",
				i+1, len(files), inputPath, numChars, numMisses, err)
			fmt.Fprintf(fBad, "%q version=%s MB=%.1f pages=%d secs=%.1f numChars=%d numMisses=%d err=%v\n",
				inputPath, version, sizeMB, numPages, dt.Seconds(), numChars, numMisses, err)
		}
		if verbose {
			fmt.Println("========================= ~~~ =========================")
		}
		if err != nil {
			errorCounts[err.Error()]++
		}
	}
	fmt.Fprintf(os.Stderr, "Done %d files \n", len(files))
	if len(errorCounts) > 0 {
		fmt.Fprintln(os.Stderr, "=== Error counts ===")
		for err, n := range errorCounts {
			fmt.Fprintf(os.Stderr, "%-30s %d (%.0f%%)\n", err, n, 100.0*float64(n)/float64(len(files)))
		}
		fmt.Fprintf(os.Stderr, "badFilesPath=%q\n", badFilesPath)
		fmt.Fprintf(os.Stderr, "threshold=%.1f%%\n", threshold)
	}

	fmt.Fprintf(os.Stderr, "badFilesPath=%q\n", badFilesPath)
	fmt.Fprintf(os.Stderr, "threshold=%.1f%%\n", threshold)
}

// missclassificationError returns an error if the percentage of misclassified characters exceeds
// `threshold`.
func missclassificationError(threshold float64, numChars, numMisses int) error {
	if numChars == 0 || numMisses == 0 {
		return nil
	}
	if threshold*float64(numChars) >= 1.0 && percentage(numChars, numMisses) < threshold {
		return nil
	}
	return ErrBadText
}

// percentage returns the percentage of missclassified characters.
func percentage(numChars, numMisses int) float64 {
	if numChars == 0 {
		return 0.0
	}
	return 100.0 * float64(numMisses) / float64(numChars)
}

// "~/testdata/2005JE002531.pdf" version=1.3 MB=1.3 pages=18 secs=0.0 numChars=484 numMisses=1 err=Could not decode text (ToUnicode)
var reFilename = regexp.MustCompile(`^\s*"(.+?)"\s*` +
	`version=([\d\.]+)\s+` +
	`MB=([\d\.]+)\s+` +
	`pages=(\d+)\s+` +
	`secs=([\d\.]+)\s+` +
	`numChars=(\d+)\s+` +
	`numMisses=(\d+)\s+` +
	`err=(.+?)\s*$`)

type testResult struct {
	filename string
	version  string
	mbytes   float64
	pages    int
	seconds  float64
	chars    int
	misses   int
	err      string
}

var testHeader = []string{"filename", "version", "mbytes", "pages", "seconds", "chars", "misses",
	"error"}

func (r *testResult) String() string {
	return fmt.Sprintf("version=%s pages=%3d mbytes=%3.1f seconds=%2.1f chars=%d misses=%d "+
		"err=%-20s %q",
		r.version, r.pages, r.mbytes, r.seconds, r.chars, r.misses, r.err, r.filename)
}

func (r *testResult) asStrings() []string {
	return []string{
		r.filename,
		r.version,
		fmt.Sprintf("%.1f", r.mbytes),
		fmt.Sprintf("%d", r.pages),
		fmt.Sprintf("%.1f", r.seconds),
		fmt.Sprintf("%d", r.chars),
		fmt.Sprintf("%d", r.misses),
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

		mbytes, err := strconv.ParseFloat(groups[3], 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Bad MB line %d in %q: line=%q\n", i, filename, line)
			continue
		}
		pages, err := strconv.Atoi(groups[4])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Bad pages line %d in %q: line=%q\n", i, filename, line)
			continue
		}
		seconds, err := strconv.ParseFloat(groups[5], 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Bad secs line %d in %q: line=%q\n", i, filename, line)
			continue
		}
		chars, err := strconv.Atoi(groups[6])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Bad chars line %d in %q: line=%q\n", i, filename, line)
			continue
		}
		misses, err := strconv.Atoi(groups[7])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Bad misses line %d in %q: line=%q\n", i, filename, line)
			continue
		}

		r := testResult{
			filename: groups[1],
			version:  groups[2],
			mbytes:   mbytes,
			pages:    pages,
			seconds:  seconds,
			chars:    chars,
			misses:   misses,
			err:      groups[8],
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
		if ri.seconds != rj.seconds {
			return ri.seconds < rj.seconds
		}

		if ri.err != rj.err {
			return ri.err < rj.err
		}
		if ri.misses != rj.misses {
			return ri.misses > rj.misses
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

// saveAsCsv saves `results` as a CSV file
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
		return nil, 0, err
	}
	defer f.Close()

	pdfReader, err = pdf.NewPdfReader(f)
	if err != nil {
		return nil, 0, err
	}
	numPages, err = pdfReader.GetNumPages()
	return pdfReader, numPages, err
}

// outputPdfText prints out text of PDF file `inputPath` to stdout.
// `pdfReader` is a previously opened PdfReader of `inputPath`
func outputPdfText(inputPath string, pdfReader *pdf.PdfReader, numPages, width int, verbose bool) (int, int, error) {
	numChars, numMisses := 0, 0
	for pageNum := 1; pageNum <= numPages; pageNum++ {

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return 0, 0, err
		}
		ex, err := extractor.New(page)
		if err != nil {
			return 0, 0, err
		}
		text, nChars, nMisses, err := ex.ExtractText2()
		numChars += nChars
		numMisses += nMisses
		if err != nil {
			return numChars, numMisses, err
		}
		if width != 0 {
			text = normalizeText(text, width)
		}

		if verbose {
			fmt.Printf("Page %d of %d: %q\n", pageNum, numPages, inputPath)
		}
		fmt.Printf("%s\n", text)
		if verbose {
			fmt.Println("------------------------- ... -------------------------")
		}
	}

	return numChars, numMisses, nil
}

func makeUsage(msg string) {
	usage := flag.Usage
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, msg)
		usage()
	}
}

// fileSizeMB returns the size of file `path` in megabytes
func fileSizeMB(path string) float64 {
	fi, err := os.Stat(path)
	if err != nil {
		return -1.0
	}
	return float64(fi.Size()) / 1024.0 / 1024.0
}

// normalizeText returns `text` with runs of spaces of any kind (spaces, tabs, line breaks, etc)
// reduced to a single space. `width` is the target line width.
func normalizeText(text string, width int) string {
	if width < 0 {
		width = defaultNormalizeWidth
	}
	return splitLines(reduceSpaces(text), width)
}

// reduceSpaces returns `text` with runs of spaces of any kind (spaces, tabs, line breaks, etc)
// reduced to a single space.
func reduceSpaces(text string) string {
	text = reSpace.ReplaceAllString(text, " ")
	return strings.Trim(text, " \t\n\r\v")
}

var reSpace = regexp.MustCompile(`(?m)\s+`)

// splitLines inserts line breaks in string `text`. `width` is the target line width.
func splitLines(text string, width int) string {
	runes := []rune(text)
	if len(runes) < 2 {
		return text
	}
	lines := []string{}
	chars := []rune{}
	for i := 0; i < len(runes)-1; i++ {
		r, r1 := runes[i], runes[i+1]
		chars = append(chars, r)
		if (len(chars) >= width && unicode.IsSpace(r)) || (r == '.' && unicode.IsSpace(r1)) {
			lines = append(lines, string(chars))
			chars = []rune{}
		}
	}
	chars = append(chars, runes[len(runes)-1])
	if len(chars) > 0 {
		lines = append(lines, string(chars))
	}
	for i, ln := range lines {
		lines[i] = strings.Trim(ln, " \t\n\r\v")
	}
	return strings.Join(lines, "\n")
}
