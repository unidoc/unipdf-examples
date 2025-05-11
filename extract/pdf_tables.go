/*
 * Extract all tables from the specified pages of one or more PDF files.
 *
 * Run as: go run pdf_tables.go input.pdf
 */

package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime/pprof"
	"sort"
	"strings"
	"time"

	"golang.org/x/text/unicode/norm"

	"github.com/bmatcuk/doublestar"
	"github.com/unidoc/unipdf/v4/common"
	"github.com/unidoc/unipdf/v4/common/license"
	"github.com/unidoc/unipdf/v4/extractor"
	"github.com/unidoc/unipdf/v4/model"
	"github.com/unidoc/unipdf/v4/pdfutil"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}
}

const (
	usage = "Usage: go run pdf_tables.go [options] <file1.pdf> <file2.pdf> ...\n"
)

func main() {
	var (
		firstPage, lastPage     int
		width, height           int
		csvDir                  string
		debug, trace, doProfile bool
		verbose                 int
	)
	flag.StringVar(&csvDir, "o", "./outcsv", `Output CSVs (default outtext). Set to "" to not save.`)
	flag.IntVar(&firstPage, "f", -1, "First page.")
	flag.IntVar(&lastPage, "l", 100000, "Last page.")
	flag.IntVar(&width, "w", 0, "Minimum table width.")
	flag.IntVar(&height, "h", 0, "Minimum table height.")
	flag.IntVar(&verbose, "v", 1, `Verbosity level
    <empty>                (level 0)
    N pages M tables       (level 1) default
    page I: M tables       (level 2)
      table J: W x H       (level 3)
         table content     (level 4)
	`)
	flag.BoolVar(&debug, "d", false, "Print debugging information.")
	flag.BoolVar(&trace, "e", false, "Print detailed debugging information.")
	flag.BoolVar(&doProfile, "p", false, "Save profiling information.")
	makeUsage(usage)
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	if trace {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelTrace))
	} else if debug {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
	} else {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelInfo))
	}

	makeDir("CSV directory", csvDir)

	pathList, err := patternsToPaths(os.Args[1:])
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d PDF files\n", len(pathList))

	if doProfile {
		f, err := os.Create("cpu.profile")
		if err != nil {
			panic(fmt.Errorf("could not create CPU profile: err=%w", err))
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			panic(fmt.Errorf("could not start CPU profile: err=%w", err))
		}
		defer pprof.StopCPUProfile()
	}

	for i, inPath := range pathList {
		t0 := time.Now()
		result, err := extractTables(inPath, firstPage, lastPage)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		duration := time.Since(t0).Seconds()
		numPages := len(result.pageTables)
		result = result.filter(width, height)
		fmt.Printf("%3d of %d: %4.1f MB %3d pages %4.1f sec %q %s",
			i+1, len(pathList), fileSizeMB(inPath), numPages, duration, inPath, result.describe(verbose))
		csvRoot := changeDirExt(csvDir, filepath.Base(inPath), "", "")
		if err := result.saveCSVFiles(csvRoot); err != nil {
			fmt.Printf("Failed to write %q: %v\n", csvRoot, err)
			continue
		}
	}
}

// extractTables extracts tables from pages `firstPage` to `lastPage` in PDF file `inPath`.
func extractTables(inPath string, firstPage, lastPage int) (docTables, error) {
	f, err := os.Open(inPath)
	if err != nil {
		return docTables{}, fmt.Errorf("Could not open %q err=%w", inPath, err)
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReaderLazy(f)
	if err != nil {
		return docTables{}, fmt.Errorf("NewPdfReaderLazy failed. %q err=%w", inPath, err)
	}
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return docTables{}, fmt.Errorf("GetNumPages failed. %q err=%w", inPath, err)
	}

	if firstPage < 1 {
		firstPage = 1
	}
	if lastPage > numPages {
		lastPage = numPages
	}

	result := docTables{pageTables: make(map[int][]stringTable)}
	for pageNum := firstPage; pageNum <= lastPage; pageNum++ {
		tables, err := extractPageTables(pdfReader, pageNum)
		if err != nil {
			return docTables{}, fmt.Errorf("extractPageTables failed. inPath=%q pageNum=%d err=%w",
				inPath, pageNum, err)
		}
		result.pageTables[pageNum] = tables
	}
	return result, nil
}

// extractPageTables extracts the tables from (1-offset) page number `pageNum` in opened
// PdfReader `pdfReader.
func extractPageTables(pdfReader *model.PdfReader, pageNum int) ([]stringTable, error) {
	page, err := pdfReader.GetPage(pageNum)
	if err != nil {
		return nil, err
	}
	if err := pdfutil.NormalizePage(page); err != nil {
		return nil, err
	}

	ex, err := extractor.New(page)
	if err != nil {
		return nil, err
	}
	pageText, _, _, err := ex.ExtractPageText()
	if err != nil {
		return nil, err
	}
	tables := pageText.Tables()
	stringTables := make([]stringTable, len(tables))
	for i, table := range tables {
		stringTables[i] = asStringTable(table)
	}
	return stringTables, nil
}

// docTables describes the tables in a document.
type docTables struct {
	pageTables map[int][]stringTable
}

// stringTable is the strings in TextTable.
type stringTable [][]string

func (r docTables) saveCSVFiles(csvRoot string) error {
	for _, pageNum := range r.pageNumbers() {
		for i, table := range r.pageTables[pageNum] {
			csvPath := fmt.Sprintf("%s.page%d.table%d.csv", csvRoot, pageNum, i+1)
			contents := table.csv()
			if err := ioutil.WriteFile(csvPath, []byte(contents), 0666); err != nil {
				return fmt.Errorf("failed to write csvPath=%q err=%w", csvPath, err)
			}
		}
	}
	return nil
}

// wh returns the width and height of table `t`.
func (t stringTable) wh() (int, int) {
	if len(t) == 0 {
		return 0, 0
	}
	return len(t[0]), len(t)
}

// csv returns `t` in CSV format.
func (t stringTable) csv() string {
	w, h := t.wh()
	b := new(bytes.Buffer)
	csvwriter := csv.NewWriter(b)
	for y, row := range t {
		if len(row) != w {
			err := fmt.Errorf("table = %d x %d row[%d]=%d %q", w, h, y, len(row), row)
			panic(err)
		}
		csvwriter.Write(row)
	}
	csvwriter.Flush()
	return b.String()
}

func (r *docTables) String() string {
	return r.describe(1)
}

// describe returns a string describing the tables in `r`.
//
//	                            (level 0)
//	%d pages %d tables          (level 1)
//	  page %d: %d tables        (level 2)
//	    table %d: %d x %d       (level 3)
//	        contents            (level 4)
func (r *docTables) describe(level int) string {
	if level == 0 || r.numTables() == 0 {
		return "\n"
	}
	var sb strings.Builder
	pageNumbers := r.pageNumbers()
	fmt.Fprintf(&sb, "%d pages %d tables\n", len(pageNumbers), r.numTables())
	if level <= 1 {
		return sb.String()
	}
	for _, pageNum := range r.pageNumbers() {
		tables := r.pageTables[pageNum]
		if len(tables) == 0 {
			continue
		}
		fmt.Fprintf(&sb, "   page %d: %d tables\n", pageNum, len(tables))
		if level <= 2 {
			continue
		}
		for i, table := range tables {
			w, h := table.wh()
			fmt.Fprintf(&sb, "      table %d: %d x %d\n", i+1, w, h)
			if level <= 3 || len(table) == 0 {
				continue
			}
			for _, row := range table {
				cells := make([]string, len(row))
				for i, cell := range row {
					if len(cell) > 0 {
						cells[i] = fmt.Sprintf("%q", cell)
					}
				}
				fmt.Fprintf(&sb, "        [%s]\n", strings.Join(cells, ", "))
			}
		}
	}
	return sb.String()
}

func (r *docTables) pageNumbers() []int {
	pageNums := make([]int, len(r.pageTables))
	i := 0
	for pageNum := range r.pageTables {
		pageNums[i] = pageNum
		i++
	}
	sort.Ints(pageNums)
	return pageNums
}

func (r *docTables) numTables() int {
	n := 0
	for _, tables := range r.pageTables {
		n += len(tables)
	}
	return n
}

// filter returns the tables in `r` that are at least `width` cells wide and `height` cells high.
func (r docTables) filter(width, height int) docTables {
	filtered := docTables{pageTables: make(map[int][]stringTable)}
	for pageNum, tables := range r.pageTables {
		var filteredTables []stringTable
		for _, table := range tables {
			if len(table[0]) >= width && len(table) >= height {
				filteredTables = append(filteredTables, table)
			}
		}
		if len(filteredTables) > 0 {
			filtered.pageTables[pageNum] = filteredTables
		}
	}
	return filtered
}

// asStringTable returns TextTable `table` as a stringTable.
func asStringTable(table extractor.TextTable) stringTable {
	cells := make(stringTable, table.H)
	for y, row := range table.Cells {
		cells[y] = make([]string, table.W)
		for x, cell := range row {
			cells[y][x] = cell.Text
		}
	}
	return normalizeTable(cells)
}

// normalizeTable returns `cells` with each cell normalized.
func normalizeTable(cells stringTable) stringTable {
	for y, row := range cells {
		for x, cell := range row {
			cells[y][x] = normalize(cell)
		}
	}
	return cells
}

// normalize returns a version of `text` that is NFKC normalized and has reduceSpaces() applied.
func normalize(text string) string {
	return reduceSpaces(norm.NFKC.String(text))
}

// reduceSpaces returns `text` with runs of spaces of any kind (spaces, tabs, line breaks, etc)
// reduced to a single space.
func reduceSpaces(text string) string {
	text = reSpace.ReplaceAllString(text, " ")
	return strings.Trim(text, " \t\n\r\v")
}

var reSpace = regexp.MustCompile(`(?m)\s+`)

// patternsToPaths returns the file paths matched by the patterns in `patternList`.
func patternsToPaths(patternList []string) ([]string, error) {
	var pathList []string
	common.Log.Debug("patternList=%d", len(patternList))
	for i, pattern := range patternList {
		pattern = expandUser(pattern)
		files, err := doublestar.Glob(pattern)
		if err != nil {
			common.Log.Error("PatternsToPaths: Glob failed. pattern=%#q err=%v", pattern, err)
			return pathList, err
		}
		common.Log.Debug("patternList[%d]=%q %d matches", i, pattern, len(files))
		for _, filename := range files {
			ok, err := regularFile(filename)
			if err != nil {
				common.Log.Error("PatternsToPaths: regularFile failed. pattern=%#q err=%v", pattern, err)
				return pathList, err
			}
			if !ok {
				continue
			}
			pathList = append(pathList, filename)
		}
	}
	// pathList = StringUniques(pathList)
	sort.Strings(pathList)
	return pathList, nil
}

// homeDir is the current user's home directory.
var homeDir = getHomeDir()

// getHomeDir returns the current user's home directory.
func getHomeDir() string {
	usr, _ := user.Current()
	return usr.HomeDir
}

// expandUser returns `filename` with "~"" replaced with user's home directory.
func expandUser(filename string) string {
	return strings.Replace(filename, "~", homeDir, -1)
}

// regularFile returns true if file `filename` is a regular file.
func regularFile(filename string) (bool, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return false, err
	}
	return fi.Mode().IsRegular(), nil
}

// fileSizeMB returns the size of file `filename` in megabytes.
func fileSizeMB(filename string) float64 {
	fi, err := os.Stat(filename)
	if err != nil {
		panic(err)
	}
	return float64(fi.Size()) / 1024.0 / 1024.0
}

// makeUsage updates flag.Usage to include usage message `msg`.
func makeUsage(msg string) {
	usage := flag.Usage
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, msg)
		usage()
	}
}

// makeDir creates `outDir`. Name is the name of `outDir` in the calling code.
func makeDir(name, outDir string) {
	if outDir == "." || outDir == ".." {
		panic(fmt.Errorf("%s=%q not allowed", name, outDir))
	}
	if outDir == "" {
		return
	}

	outDir, err := filepath.Abs(outDir)
	if err != nil {
		panic(fmt.Errorf("Abs failed. %s=%q err=%w", name, outDir, err))
	}
	if err := os.MkdirAll(outDir, 0777); err != nil {
		panic(fmt.Errorf("Couldn't create %s=%q err=%w", name, outDir, err))
	}
}

// changeDirExt inserts `qualifier` into `filename` before its extension then changes its
// directory to `dirName` and extension to `extName`,
func changeDirExt(dirName, filename, qualifier, extName string) string {
	if dirName == "" {
		return ""
	}
	base := filepath.Base(filename)
	ext := filepath.Ext(base)
	base = base[:len(base)-len(ext)]
	if len(qualifier) > 0 {
		base = fmt.Sprintf("%s.%s", base, qualifier)
	}
	filename = fmt.Sprintf("%s%s", base, extName)
	path := filepath.Join(dirName, filename)
	common.Log.Debug("changeDirExt(%q,%q,%q)->%q", dirName, base, extName, path)
	return path
}
