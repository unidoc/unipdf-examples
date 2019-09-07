/*
 * PDF to CSV: Example illustrating capability to extract TextMarks from PDF, and grouping together into
 * words, rows and columns for CSV data extraction.
 *
 * Includes debugging capabilities such as outputing a marked up PDF showing bounding boxes of marks,
 * words, lines and columns.
 *
 * Run as: go run pdf_to_csv.go -m all -mf markup.pdf table.pdf table.csv
 * - Outputs debug markup including: marks, words, lines, columns to markup.pdf
 * - The table data is outputed to table.csv with UTF-8 encoding.
 */

package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/contentstream"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

var saveParams saveMarkedupParams

func main() {
	// Make sure to enter a valid license key.
	// Otherwise text is truncated and a watermark added to the text.
	// License keys are available via: https://unidoc.io
	/*
			license.SetLicenseKey(`
		-----BEGIN UNIDOC LICENSE KEY-----
		...key contents...
		-----END UNIDOC LICENSE KEY-----
		`, "Customer Name")
	*/
	// Alternatively license can be loaded via UNIPDF_LICENSE_PATH and UNIPDF_CUSTOMER_NAME environment variables,
	// where UNIPDF_LICENSE_PATH points to the file containing the license key and the UNIPDF_CUSTOMER_NAME
	// the explicitly specified customer name to which the key is licensed.
	var (
		loglevel   string
		saveMarkup string
		markupPath string
	)
	flag.StringVar(&loglevel, "l", "info", "Set log level (default: info)")
	flag.StringVar(&saveMarkup, "m", "none", "Save markup (none/marks/words/lines/columns/all)")
	flag.StringVar(&markupPath, "mf", "/tmp/markup.pdf", "Output markup path (default /tmp/markup.pdf)")
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		fmt.Printf("Syntax: go run pdf_to_csv.go [options] <file.pdf> <output.csv>\n")
		os.Exit(1)
	}

	switch loglevel {
	case "trace":
		common.SetLogger(common.NewConsoleLogger(common.LogLevelTrace))
	case "debug":
		common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
	default:
		common.SetLogger(common.NewConsoleLogger(common.LogLevelInfo))
	}

	saveParams = saveMarkedupParams{}
	switch saveMarkup {
	case "marks":
		saveParams.markupType = "marks"
	case "words":
		saveParams.markupType = "words"
	case "lines":
		saveParams.markupType = "lines"
	case "columns":
		saveParams.markupType = "columns"
	case "all":
		saveParams.markupType = "all"
	default:
		saveParams.markupType = "none"
	}
	saveParams.markupOutputPath = markupPath

	inPath := args[0]
	outPath := args[1]
	err := extractTableData(inPath, outPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// extractTableData extracts tabular information from PDF file `inPath` and outputs
// the data as CSV file to `outPath`.
func extractTableData(inPath string, outPath string) error {
	f, err := os.Open(inPath)
	if err != nil {
		return fmt.Errorf("Could not open %q err=%v", inPath, err)
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReaderLazy(f)
	if err != nil {
		return fmt.Errorf("NewPdfReaderLazy failed. %q err=%v", inPath, err)
	}
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return fmt.Errorf("GetNumPages failed. %q err=%v", inPath, err)
	}

	saveParams.pdfReader = pdfReader
	saveParams.markups = map[int][][]model.PdfRectangle{}

	var csvData bytes.Buffer
	for pageNum := 1; pageNum <= numPages; pageNum++ {
		if pageNum > 3 {
			break
		}
		saveParams.curPage = pageNum

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return fmt.Errorf("GetNumPages failed. %q pageNum=%d err=%v", inPath, pageNum, err)
		}

		mbox, err := page.GetMediaBox()
		if err != nil {
			return err
		}
		if page.Rotate != nil && *page.Rotate == 90 {
			// TODO: This is a "hack" to change the perspective of the extractor to account for the rotation.
			contents, err := page.GetContentStreams()
			if err != nil {
				return err
			}

			cc := contentstream.NewContentCreator()
			cc.Translate(mbox.Width()/2, mbox.Height()/2)
			cc.RotateDeg(-90)
			cc.Translate(-mbox.Width()/2, -mbox.Height()/2)
			rotateOps := cc.Operations().String()
			contents = append([]string{rotateOps}, contents...)

			page.Duplicate()
			err = page.SetContentStreams(contents, core.NewRawEncoder())
			if err != nil {
				return err

			}
			page.Rotate = nil
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

		group := []model.PdfRectangle{}
		for _, mark := range textMarks.Elements() {
			group = append(group, mark.BBox)
		}
		saveParams.markups[pageNum] = append(saveParams.markups[pageNum], group)

		pageCSV, err := pageMarksToCSV(textMarks)
		if err != nil {
			common.Log.Debug("Error grouping text: %v", err)
			return err
		}
		csvData.WriteString(pageCSV)
	}

	if saveParams.markupType != "none" {
		err = saveMarkedupPDF(saveParams)
		if err != nil {
			return fmt.Errorf("Failed to save marked up pdf: %v", err)
		}
	}

	return ioutil.WriteFile(outPath, csvData.Bytes(), 0666)
}

func rectUnion(b1, b2 model.PdfRectangle) model.PdfRectangle {
	return model.PdfRectangle{
		Llx: math.Min(b1.Llx, b2.Llx),
		Lly: math.Min(b1.Lly, b2.Lly),
		Urx: math.Max(b1.Urx, b2.Urx),
		Ury: math.Max(b1.Ury, b2.Ury),
	}
}

func bboxArea(bbox model.PdfRectangle) float64 {
	return math.Abs(bbox.Urx-bbox.Llx) * math.Abs(bbox.Ury-bbox.Lly)
}

// Measure of the difference between areas of `bbox1` and `bbox2` individually
// and that of the union of the two.
func overlaps(bbox1, bbox2 model.PdfRectangle) float64 {
	union := rectUnion(bbox1, bbox2)
	a := bboxArea(union)
	b := bboxArea(bbox1) + bboxArea(bbox2)
	diff := (a - b) / (a + b)
	return diff
}

// Measure of the vertical overlap of `bbox1` and `bbox2`, when the difference is 0
// then they are exactly on top of each other, and there is overlap when < 0.
func lineOverlap(bbox1, bbox2 model.PdfRectangle) float64 {
	union := rectUnion(bbox1, bbox2)
	a := math.Abs(union.Ury - union.Lly)
	b := math.Abs(bbox1.Ury-bbox1.Lly) + math.Abs(bbox2.Ury-bbox2.Lly)
	diff := (a - b) / (a + b)
	return diff
}

// Measure of the horizontal overlap of `bbox1` and `bbox2`, when the difference is 0
// then they are exactly on next to each other, and there is overlap when < 0.
func columnOverlap(bbox1, bbox2 model.PdfRectangle) float64 {
	union := rectUnion(bbox1, bbox2)
	a := math.Abs(union.Urx - union.Llx)
	b := math.Abs(bbox1.Urx-bbox1.Llx) + math.Abs(bbox2.Urx-bbox2.Llx)
	diff := (a - b) / (a + b)
	return diff
}

// Identify lines. - segment words into lines
func identifyLines(words []segmentationWord) [][]segmentationWord {
	lines := [][]segmentationWord{}

	for _, word := range words {
		wbbox, ok := word.BBox()
		if !ok {
			continue
		}

		match := false
		for i, line := range lines {
			firstWord := line[0]
			firstBBox, ok := firstWord.BBox()
			if !ok {
				continue
			}

			overlap := lineOverlap(wbbox, firstBBox)
			common.Log.Debug("'%s'/'%s' overlap: %v [%+v/%+v]", word.String(), firstWord.String(), overlap, wbbox, firstBBox)
			if overlap < 0 {
				lines[i] = append(lines[i], word)
				match = true
				break
			}
		}
		if !match {
			lines = append(lines, []segmentationWord{word})
		}
	}
	sort.SliceStable(lines, func(i, j int) bool {
		bboxi, _ := lines[i][0].BBox()
		bboxj, _ := lines[j][0].BBox()
		return bboxi.Lly >= bboxj.Lly
	})
	for li := range lines {
		sort.SliceStable(lines[li], func(i, j int) bool {
			bboxi, _ := lines[li][i].BBox()
			bboxj, _ := lines[li][j].BBox()
			return bboxi.Llx < bboxj.Llx
		})
	}

	// Save the line bounding boxes for markup output.
	lineGroups := []model.PdfRectangle{}
	for li, line := range lines {
		var lineRect model.PdfRectangle
		common.Log.Trace("Line %d: ", li+1)
		for i, word := range line {
			wbbox, ok := word.BBox()
			if !ok {
				continue
			}
			common.Log.Trace("'%s' / ", word.String())

			if i == 0 {
				lineRect = wbbox
			} else {
				lineRect = rectUnion(lineRect, wbbox)
			}
		}
		lineGroups = append(lineGroups, lineRect)
	}
	saveParams.markups[saveParams.curPage] = append(saveParams.markups[saveParams.curPage], lineGroups)
	return lines
}

// Identify columns.
func identifyColumns(words []segmentationWord) []model.PdfRectangle {
	columns := [][]segmentationWord{}
	for _, word := range words {
		wbbox, ok := word.BBox()
		if !ok {
			continue
		}

		match := false
		bestOverlap := 1.0
		bestColumn := 0
		for i, column := range columns {
			firstWord := column[0]
			firstBBox, ok := firstWord.BBox()
			if !ok {
				continue
			}

			overlap := columnOverlap(wbbox, firstBBox)
			common.Log.Debug("column: '%s'/'%s' overlap: %v [%+v/%+v]", word.String(), firstWord.String(), overlap, wbbox, firstBBox)
			if overlap < 0.0 {
				if overlap < bestOverlap {
					bestOverlap = overlap
					bestColumn = i
				}
				match = true
			}
		}
		if match {
			columns[bestColumn] = append(columns[bestColumn], word)
		} else {
			columns = append(columns, []segmentationWord{word})
		}
	}
	sort.SliceStable(columns, func(i, j int) bool {
		bboxi, _ := columns[i][0].BBox()
		bboxj, _ := columns[j][0].BBox()
		return bboxi.Llx < bboxj.Llx
	})
	for li := range columns {
		sort.SliceStable(columns[li], func(i, j int) bool {
			bboxi, _ := columns[li][i].BBox()
			bboxj, _ := columns[li][j].BBox()
			return bboxi.Lly >= bboxj.Lly
		})
	}

	colGroups := []model.PdfRectangle{}
	for li, column := range columns {
		var colRect model.PdfRectangle
		common.Log.Debug("Column %d: ", li+1)
		for i, word := range column {
			wbbox, ok := word.BBox()
			if !ok {
				continue
			}
			common.Log.Debug("    '%s' / ", word.String())

			if i == 0 {
				colRect = wbbox
			} else {
				colRect = rectUnion(colRect, wbbox)
			}
		}
		common.Log.Debug("Column %d: Bbox: %+v", li+1, colRect)
		colGroups = append(colGroups, colRect)
	}

	// Filter by combining overlapping columns.
	filtered := []model.PdfRectangle{}
	for i := 0; i < len(colGroups); {
		colgroup := colGroups[i]
		j := i + 1
		for ; j < len(colGroups); j++ {
			overlap := columnOverlap(colgroup, colGroups[j])
			common.Log.Debug("COLUMN overlap %d/%d: %v (%+v/%+v)", i+1, j+1, overlap, colgroup, colGroups[j])
			if overlap > 0.0 {
				break
			}
			colgroup = rectUnion(colgroup, colGroups[j])
		}
		i = j
		filtered = append(filtered, colgroup)
	}

	saveParams.markups[saveParams.curPage] = append(saveParams.markups[saveParams.curPage], filtered)
	return filtered
}

// getLineTableTextData converts the lines of words into table strings cells by accounting for
// distribution of lines into columns as specified by `columnBBoxes`.
func getLineTableTextData(lines [][]segmentationWord, columnBBoxes []model.PdfRectangle) [][]string {
	tabledata := [][]string{}
	for _, line := range lines {
		linedata := make([]string, len(columnBBoxes))
		for _, word := range line {
			wordBBox, ok := word.BBox()
			if !ok {
				continue
			}

			bestColumn := 0
			bestOverlap := 1.0
			for icol, colBBox := range columnBBoxes {
				overlap := columnOverlap(wordBBox, colBBox)
				if overlap < bestOverlap {
					bestOverlap = overlap
					bestColumn = icol
				}
			}
			linedata[bestColumn] += word.String()
		}
		tabledata = append(tabledata, linedata)
	}
	return tabledata
}

// segmentationWord represents a word that has been segmented in PDF text.
type segmentationWord struct {
	ma *extractor.TextMarkArray
}

func (w segmentationWord) Elements() []extractor.TextMark {
	return w.ma.Elements()
}

func (w segmentationWord) BBox() (model.PdfRectangle, bool) {
	return w.ma.BBox()
}

func (w segmentationWord) String() string {
	if w.ma == nil {
		return ""
	}

	var buf bytes.Buffer
	for _, m := range w.Elements() {
		buf.WriteString(m.Text)
	}
	return buf.String()
}

// pageMarksToCSV converts textMarks from a single page into CSV by grouping the marks into
// words, lines and columns and then writing the table cells data as CSV output.
func pageMarksToCSV(textMarks *extractor.TextMarkArray) (string, error) {
	// STEP - Form words.
	// Group the closest text marks that are overlapping.
	words := []segmentationWord{}
	word := segmentationWord{ma: &extractor.TextMarkArray{}}
	var lastMark extractor.TextMark
	isFirst := true
	for i, mark := range textMarks.Elements() {
		if mark.Text == "" {
			continue
		}

		common.Log.Debug("Mark %d - '%s' (% X)", i, mark.Text, mark.Text)
		if isFirst {
			word = segmentationWord{ma: &extractor.TextMarkArray{}}
			word.ma.Append(mark)
			lastMark = mark
			isFirst = false
			continue
		}
		common.Log.Debug(" - overlaps: %f", overlaps(mark.BBox, lastMark.BBox))
		overlap := overlaps(mark.BBox, lastMark.BBox)
		if overlap > 0.1 {
			if len(strings.TrimSpace(word.String())) > 0 {
				common.Log.Debug("Appending word: '%s' (%d chars) (%d elements)", word.String(), len(word.String()), len(word.Elements()))
				words = append(words, word)
			}
			word = segmentationWord{ma: &extractor.TextMarkArray{}}
		}
		word.ma.Append(mark)
		lastMark = mark
	}
	if len(strings.TrimSpace(word.String())) > 0 {
		common.Log.Debug("Appending word: '%s' (%d chars) (%d elements)", word.String(), len(word.String()), len(word.Elements()))
		words = append(words, word)
	}

	// Include the words in the markup.
	{
		wbboxes := []model.PdfRectangle{}
		for _, word := range words {
			wbbox, ok := word.BBox()
			if !ok {
				continue
			}
			wbboxes = append(wbboxes, wbbox)
		}
		saveParams.markups[saveParams.curPage] = append(saveParams.markups[saveParams.curPage], wbboxes)
	}

	lines := identifyLines(words)

	// Filter out words in lines with only 1 column.
	tableLines := [][]segmentationWord{}
	for _, line := range lines {
		if len(line) <= 1 {
			continue
		}
		tableLines = append(tableLines, line)
	}

	tableWords := []segmentationWord{}
	for _, line := range tableLines {
		for _, word := range line {
			tableWords = append(tableWords, word)
		}
	}

	columnBBoxes := identifyColumns(tableWords)

	tabledata := getLineTableTextData(lines, columnBBoxes)

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	err := w.WriteAll(tabledata)
	if err != nil {
		return "", err
	}
	w.Flush()
	return buf.String(), nil
}

type saveMarkedupParams struct {
	pdfReader        *model.PdfReader
	markups          map[int][][]model.PdfRectangle
	curPage          int
	markupType       string
	markupOutputPath string
}

// Saves a marked up PDF with the original with certain groups highlighted: marks, words, lines, columns.
func saveMarkedupPDF(params saveMarkedupParams) error {
	var pageNums []int
	for pageNum, _ := range params.markups {
		pageNums = append(pageNums, pageNum)
	}
	sort.Ints(pageNums)
	if len(pageNums) == 0 {
		return nil
	}

	// Make a new PDF creator.
	c := creator.New()
	for _, pageNum := range pageNums {
		common.Log.Debug("Page %d - %d marks", pageNum, len(params.markups[pageNum]))
		page, err := params.pdfReader.GetPage(pageNum)
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
			return fmt.Errorf("AddPage failed err=%v ", err)
		}

		colors := map[int]string{
			0: "hide", // marks
			1: "hide", // words
			2: "hide", // lines
			3: "hide", // columns
		}

		switch saveParams.markupType {
		case "marks":
			colors[0] = "#0000ff"
		case "words":
			colors[1] = "#00ff00"
		case "lines":
			colors[2] = "#ff0000"
		case "columns":
			colors[3] = "#f0f000"
		case "all":
			colors[0] = "#0000ff"
			colors[1] = "#00ff00"
			colors[2] = "#ff0000"
			colors[3] = "#f0f000"
		}

		for gi, group := range params.markups[pageNum] {
			if colors[gi] == "hide" {
				continue
			}
			borderColor := creator.ColorRGBFromHex(colors[gi])
			for i, r := range group {
				common.Log.Trace("Mark %d", i+1)
				rect := c.NewRectangle(r.Llx, h-r.Lly, r.Urx-r.Llx, -(r.Ury - r.Lly))
				rect.SetBorderColor(borderColor)
				rect.SetBorderWidth(1.0)
				err = c.Draw(rect)
				if err != nil {
					return fmt.Errorf("Draw failed. pageNum=%d err=%v", pageNum, err)
				}
			}
		}
	}

	c.SetOutlineTree(params.pdfReader.GetOutlineTree())
	if err := c.WriteToFile(saveParams.markupOutputPath); err != nil {
		return fmt.Errorf("WriteToFile failed. err=%v", err)
	}

	common.Log.Info("Saved marked-up PDF file: %v", saveParams.markupOutputPath)
	return nil
}
