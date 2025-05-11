/*
 * Example that illustrates the accuracy of the text extraction, by first extracting
 * all TextMarks and then reconstructing the text by writing out the text page-by-page
 * to a new PDF with the creator package.
 * Only retains the text.
 *
 * Useful to check accuracy of text extraction properties.
 *
 * Run as: go run extract_text_bound.go input.pdf
 */

package main

import (
	"errors"
	"fmt"
	"os"
	"unicode"

	"github.com/unidoc/unipdf/v4/common/license"

	"github.com/unidoc/unipdf/v4/extractor"
	"github.com/unidoc/unipdf/v4/model"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}
}

type PdfWordData struct {
	Page     int
	Text     string
	Bounds   model.PdfRectangle
	Font     string
	FontSize float64
}

type PageData struct {
	Number int
	Words  []*PdfWordData
}

var pageDataList []*PageData

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Syntax: extract_text_bound <file.pdf>\n")
		os.Exit(1)
	}

	pdfPath := os.Args[1]
	err := reconstruct(pdfPath)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
}

func reconstruct(pdfPath string) error {
	f, err := os.Open(pdfPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfr, err := model.NewPdfReaderLazy(f)
	if err != nil {
		return err
	}

	for pageNum := 1; pageNum <= len(pdfr.PageList); pageNum++ {
		err = extractWordsDataOnPage(pdfr, pageNum)
		if err != nil {
			return err
		}
	}

	extractResult := ""

	for _, data := range pageDataList {
		extractResult += fmt.Sprintf("Page %d\n", data.Number)
		extractResult += "==========\n"

		for _, w := range data.Words {
			extractResult += fmt.Sprintf("- %s (%f %f %f %f)\n", w.Text, w.Bounds.Llx, w.Bounds.Lly, w.Bounds.Urx, w.Bounds.Ury)
		}

		extractResult += "\n"
	}

	outFile, err := os.Create("extract_boundary.txt")

	_, err = outFile.WriteString(extractResult)
	if err != nil {
		return err
	}

	return nil
}

func extractWordsDataOnPage(pdfReader *model.PdfReader, pageNumber int) error {
	if pdfReader == nil {
		return fmt.Errorf("It is impossible to extract words before a pdf is loaded\n")
	}

	page, err := pdfReader.GetPage(pageNumber)
	if err != nil {
		return fmt.Errorf("UniDoc pdfReader.GetNumPages failed. pageNum=%d err=%v\n", pageNumber, err)
	}

	ex, err := extractor.New(page)
	if err != nil {
		return fmt.Errorf("UniDoc pdfReader.extractor failed to create: pageNum=%d err=%v\n", pageNumber, err)
	}
	pageText, _, _, err := ex.ExtractPageText()
	if err != nil {
		return fmt.Errorf("UniDoc pdfReader failed to ExtractPageText: pageNum=%d err=%v\n", pageNumber, err)
	}

	trackingWord := false
	textMarks := pageText.Marks()
	var textMarkArrays []*extractor.TextMarkArray
	curMarkArray := new(extractor.TextMarkArray)

	// getTolerancePoint for better subscript/superscript text poisition.
	tolerancePoint := getTolerancePoint(textMarks)

	for _, textMark := range textMarks.Elements() {
		runes := []rune(textMark.Text)
		if len(runes) != 1 {
			fmt.Printf("Rune of length %v found -- %v -- Meta: %v\n", len(runes), runes, textMark.Meta)
			continue
		}
		if unicode.IsSpace(runes[0]) || textMark.Meta == true {
			if trackingWord {
				trackingWord = false
				textMarkArrays = append(textMarkArrays, curMarkArray)
				curMarkArray = new(extractor.TextMarkArray)
			}
		} else if curMarkArray.Elements() != nil &&
			curMarkArray.Elements()[curMarkArray.Len()-1].BBox.Lly > textMark.BBox.Lly+tolerancePoint {
			// If current char is at a new line then the word is splitted into multiple line
			// store each part data separately.
			if trackingWord {
				trackingWord = false
				textMarkArrays = append(textMarkArrays, curMarkArray)
				curMarkArray = new(extractor.TextMarkArray)
				curMarkArray.Append(textMark)
			}
		} else {
			if !trackingWord {
				trackingWord = true
			}
			curMarkArray.Append(textMark)
		}
	}

	pageData := &PageData{
		Number: pageNumber,
		Words:  make([]*PdfWordData, 0, len(textMarkArrays)),
	}

	for idx, textMarkArray := range textMarkArrays {
		wordData, err := extractSingleWordData(textMarkArray, pageData.Number, 1000, 1000)
		if err != nil {
			fmt.Printf("extractWordsOnPage[%v] has a nil word at index %v\n", pageData.Number, idx)
			continue
		}

		pageData.Words = append(pageData.Words, wordData)
	}

	pageDataList = append(pageDataList, pageData)

	return nil
}

func extractSingleWordData(textMarkArray *extractor.TextMarkArray, pageNumber int, pageWidth, pageHeight float64) (*PdfWordData, error) {
	wordData := new(PdfWordData)
	wordData.Page = pageNumber
	wordString := make([]rune, textMarkArray.Len(), textMarkArray.Len())

	markArrayBBox, ok := textMarkArray.BBox()
	if !ok {
		return nil, errors.New("extractSingleWord failed: There was a problem generating the bbox for textMark")
	}

	wordData.Bounds = markArrayBBox

	for idx, mark := range textMarkArray.Elements() {
		if idx == 0 {
			// Obtain the Pdf ObjectID of the first charracter in the word to assign to PdfWordData
			//wordData.ObjectID = mark.
			// TODO: [DP-6] Determine the other Font data that needs to be extracted with the words
			wordData.Font = mark.Font.String()
			wordData.FontSize = mark.FontSize
		}
		runes := []rune(mark.Text)
		if len(runes) != 1 {
			//log.Printf("ERROR: Rune of length %v found -- %v -- Meta: %v\n", len(runes), runes, mark.Meta)
			continue
		}
		if unicode.IsSpace(runes[0]) {
			continue
		}
		wordString[idx] = runes[0]
	}

	wordData.Text = string(wordString)

	return wordData, nil
}

// getTolerancePoint calculate tolerance point based on smallest font size.
func getTolerancePoint(textMarkArray *extractor.TextMarkArray) float64 {
	minFontSize := textMarkArray.Elements()[0].FontSize
	for _, textMark := range textMarkArray.Elements() {
		if textMark.FontSize < minFontSize && textMark.FontSize > 0 {
			minFontSize = textMark.FontSize
		}
	}
	return minFontSize
}
