/*
 * Markup PDF text: Mark up locations of substrings of extracted text in a PDF file.
 *
 * Run as: go run pdf_text_locations.go file.pdf term
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/extractor"
	pdf "github.com/unidoc/unipdf/v3/model"
)

const (
	// markupDir is the directory where the marked-up PDFs are saved.
	// The PDFs in markupDir can be viewed in a PDF viewer to check that they correct.
	markupDir = "marked.up"

	usage = `
	Usage: go run pdf_text_locations.go file.pdf term

	Finds all instances of term in file.pdf
	Saves marked-up PDF to marked.up/file.pdf
	Saves bounding box coordinates to marked.up/file.json
`
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
	var debug bool
	flag.BoolVar(&debug, "d", false, "Enable debug logging")
	makeUsage(usage)
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}
	if debug {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
	} else {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelInfo))
	}

	inPath := args[0]
	term := args[1]

	err := markTextLocations(inPath, term)
	if err != nil {
		fmt.Fprintf(os.Stderr, "TextLocations failed. inPath=%q term=%q err=%v\n",
			inPath, term, err)
	}
}

// markTextLocations finds all instances of `term` in the text extracted from PDF file `inPath` and
// saves a PDF file marked-up with boxes around the instances of `term` and a JSON file with the
//  box coordinates.
func markTextLocations(inPath, term string) error {
	f, err := os.Open(inPath)
	if err != nil {
		return fmt.Errorf("Could not open %q err=%v", inPath, err)
	}
	defer f.Close()
	common.Log.Info("Searching %q for %q", inPath, term)
	pdfReader, err := pdf.NewPdfReaderLazy(f)
	if err != nil {
		return fmt.Errorf("NewPdfReaderLazy failed. %q err=%v", inPath, err)
	}
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return fmt.Errorf("GetNumPages failed. %q err=%v", inPath, err)
	}
	l := createMarkupList(inPath, pdfReader)

	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return fmt.Errorf("GetNumPages failed. %q pageNum=%d err=%v", inPath, pageNum, err)
		}
		ex, err := extractor.New(page)
		if err != nil {
			return fmt.Errorf("NewPdfReaderLazy failed. %q pageNum=%d err=%v", inPath, pageNum, err)
		}
		pageText, _, _, err := ex.ExtractPageText()
		if err != nil {
			return fmt.Errorf("ExtractPageText failed. %q pageNum=%d err=%v", inPath, pageNum, err)

		}
		text := pageText.Text()
		textMarks := pageText.Marks()
		common.Log.Debug("pageNum=%d text=%d textMarks=%d", pageNum, len(text), textMarks.Len())
		matches, err := getMatches(text, textMarks, term)
		if err != nil {
			return fmt.Errorf("getMatches failed. %q pageNum=%d err=%v", inPath, pageNum, err)
		}
		if matches != nil {
			l.pageMatches[pageNum] = matches
		}
	}
	err = l.saveOutputPdf()
	if err != nil {
		return fmt.Errorf("saveOutputPdf failed. %q  err=%v", inPath, err)
	}
	return nil
}

// getMatches returns the matches (bounding box + offset) on the PDF page described by `textMarks`
// that correspond to/ all the instances of `term` in `text`, where `text` and `textMarks` are the
// extracted text returned by text := pageText.Text and textMarks := pageText.Marks().
func getMatches(text string, textMarks *extractor.TextMarkArray, term string) ([]match, error) {
	indexes := indexAll(text, term)
	if len(indexes) == 0 {
		return nil, nil
	}
	matches := make([]match, len(indexes))
	for i, start := range indexes {
		end := start + len(term)
		spanMarks, err := textMarks.RangeOffset(start, end)
		if err != nil {
			return nil, err
		}
		bbox, ok := spanMarks.BBox()
		if !ok {
			return nil, fmt.Errorf("spanMarks.BBox has no bounding box. spanMarks=%s", spanMarks)
		}
		matches[i] = match{
			Term:        term,
			OffsetRange: [2]int{start, end},
			BBox:        bbox,
		}
	}
	return matches, nil
}

// indexAll returns the indices of all instances of `term` in `text`
func indexAll(text, term string) []int {
	if len(term) == 0 {
		return nil
	}
	var indexes []int
	for start := 0; start < len(text); {
		i := strings.Index(text[start:], term)
		if i < 0 {
			return indexes
		}
		indexes = append(indexes, start+i)
		start += i + len(term)
	}
	return indexes
}

// markupList saves the results of text searches so they can be used to mark-up a PDF with search
// matches and show the search term that was matched.
// Marked up results are saved in markupDir if markupPDFs is true.
// The PDFs in markupDir can be viewed in a PDF viewer to check that they correct.
type markupList struct {
	inPath      string          // Name of input PDF to be searced searched.
	pageMatches map[int][]match // {pageNum: matches on page}
	pdfReader   *pdf.PdfReader  // Reader for input PDF
	pageNum     int             // (1-offset) Page number being worked on.
}

// match is a match of search term `Term` on a page. `BBox` is the bounding box around the matched
// term on the PDF page
type match struct {
	Term        string
	BBox        pdf.PdfRectangle
	OffsetRange [2]int
}

// String returns a description of `l`.
func (l markupList) String() string {
	return fmt.Sprintf("Term found on %d pages with input page numbers %v",
		len(l.pageMatches), l.pageNums())
}

// createMarkupList returns an initialized markupList for saving match results to so the bounding
// boxes can be checked for accuracy in a PDF viewer.
func createMarkupList(inPath string, pdfReader *pdf.PdfReader) *markupList {
	return &markupList{
		inPath:      inPath,
		pdfReader:   pdfReader,
		pageMatches: map[int][]match{},
	}
}

// pageNums returns the (1-offset) page numbers in `l` of pages that have searc matches
func (l *markupList) pageNums() []int {
	var nums []int
	for pageNum, matches := range l.pageMatches {
		if len(matches) == 0 {
			continue
		}
		nums = append(nums, pageNum)
	}
	sort.Ints(nums)
	return nums
}

// saveOutputPdf is called to mark-up a PDF file with the locations of text.
// `l` contains the input PDF, the pages, search terms and bounding boxes to mark.
func (l *markupList) saveOutputPdf() error {
	if len(l.pageNums()) == 0 {
		common.Log.Info("No marked-up PDFs to save")
		return nil
	}
	common.Log.Info("%s", l)

	os.Mkdir(markupDir, 0777)
	outPath := filepath.Join(markupDir, filepath.Base(l.inPath))
	ext := path.Ext(outPath)
	metaPath := outPath[:len(outPath)-len(ext)] + ".json"

	// Make a new PDF creator.
	c := creator.New()

	for _, pageNum := range l.pageNums() {
		common.Log.Debug("saveOutputPdf: %q pageNum=%d", l.inPath, pageNum)
		page, err := l.pdfReader.GetPage(pageNum)
		if err != nil {
			return fmt.Errorf("saveOutputPdf: Could not get page  pageNum=%d. err=%v", pageNum, err)
		}
		mediaBox, err := page.GetMediaBox()
		if err != nil {
			return fmt.Errorf("saveOutputPdf: Could not get MediaBox  pageNum=%d. err=%v", pageNum, err)
		}
		if page.MediaBox == nil {
			// Deal with MediaBox inherited from Parent.
			common.Log.Info("MediaBox: %v -> %v", page.MediaBox, mediaBox)
			page.MediaBox = mediaBox
		}
		h := mediaBox.Ury

		if err := c.AddPage(page); err != nil {
			return fmt.Errorf("AddPage failed %s:%d err=%v ", l.String(), pageNum, err)
		}

		for _, m := range l.pageMatches[pageNum] {
			r := m.BBox
			rect := c.NewRectangle(r.Llx, h-r.Lly, r.Urx-r.Llx, -(r.Ury - r.Lly))
			rect.SetBorderColor(creator.ColorRGBFromHex("#0000ff")) // Blue border.
			rect.SetBorderWidth(1.0)
			if err := c.Draw(rect); err != nil {
				return fmt.Errorf("Draw failed. pageNum=%d match=%v err=%v", pageNum, m, err)
			}
		}
	}

	c.SetOutlineTree(l.pdfReader.GetOutlineTree())
	if err := c.WriteToFile(outPath); err != nil {
		return fmt.Errorf("WriteToFile failed. err=%v", err)
	}
	common.Log.Info("Saved marked-up PDF file: %q", outPath)
	b, err := json.MarshalIndent(l.pageMatches, "", "\t")
	if err != nil {
		return fmt.Errorf("MarshalIndent failed. err=%v", err)
	}
	err = ioutil.WriteFile(metaPath, b, 0666)
	if err != nil {
		return fmt.Errorf("WriteFile failed. metaPath=%q err=%v", metaPath, err)
	}
	common.Log.Info("Saved bounding box locations file: %q", metaPath)
	return nil
}

// makeUsage updates flag.Usage to include usage message `msg`.
func makeUsage(msg string) {
	usage := flag.Usage
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, msg)
		usage()
	}
}
