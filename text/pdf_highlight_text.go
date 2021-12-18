/*
 * Highlight text: Hightlight target texts inside the PDF file.
 *
 * Run as: go run pdf_highlight_text.go inputFile.pdf outputFile.pdf term
 */

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

func init() {
	err := license.SetMeteredKey(`put your license key here`)
	if err != nil {
		panic(err)
	}
}
func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Usage: go run pdf_highlight_text.go <input.pdf> <output.pdf> <term> \n")
		os.Exit(0)
	}
	inputPath := os.Args[1]
	outputPath := os.Args[2]
	term := os.Args[3]

	err := hightlightWords(inputPath, outputPath, term)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully highlighted the word %s and created %s\n", term, outputPath)
}
func getBoundingBoxes(page *model.PdfPage, term string) ([]*model.PdfRectangle, error) {
	boundingBox := []*model.PdfRectangle{}
	ex, _ := extractor.New(page)
	pageText, _, _, err := ex.ExtractPageText()
	if err != nil {
		return nil, err
	}
	textMarks := pageText.Marks()
	text := pageText.Text()
	indexes := indexAllSubstrings(text, term)
	if len(indexes) == 0 {
		return nil, nil
	}
	for _, start := range indexes {
		end := start + len(term)
		spanMarks, err := textMarks.RangeOffset(start, end)
		if err != nil {
			return nil, err
		}
		bbox, ok := spanMarks.BBox()
		if !ok {
			return nil, fmt.Errorf("spanMarks.BBox has no bounding box. spanMarks=%s", spanMarks)
		}
		boundingBox = append(boundingBox, &bbox)
	}
	return boundingBox, nil

}
func indexAllSubstrings(text, term string) []int {
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
func hightlightWords(inputPath, outputPath, term string) error {
	reader, f, err := model.NewPdfReaderFromFile(inputPath, nil)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	numPages, err := reader.GetNumPages()
	if err != nil {
		return err
	}
	cr := creator.New()
	for n := 1; n <= numPages; n++ {
		page, err := reader.GetPage(n)
		if err != nil {
			return err
		}
		bBoxes, err := getBoundingBoxes(page, term)
		if err != nil {
			return err
		}

		mediaBox, err := page.GetMediaBox()
		if err != nil {
			return err
		}
		if page.MediaBox == nil {
			page.MediaBox = mediaBox
		}

		if err := cr.AddPage(page); err != nil {
			return err
		}
		h := mediaBox.Ury
		for _, bbox := range bBoxes {
			rect := cr.NewRectangle(bbox.Llx, h-bbox.Lly, bbox.Urx-bbox.Llx, -(bbox.Ury - bbox.Lly))
			rect.SetFillColor(creator.ColorRGBFromHex("#FFFF00"))
			rect.SetBorderWidth(0)
			rect.SetFillOpacity(0.5)
			if err := cr.Draw(rect); err != nil {
				return nil
			}
		}

	}
	cr.SetOutlineTree(reader.GetOutlineTree())

	if err := cr.WriteToFile(outputPath); err != nil {
		return fmt.Errorf("failed to write the output file %s", outputPath)
	}
	return nil
}
