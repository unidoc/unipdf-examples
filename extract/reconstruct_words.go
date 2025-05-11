/*
 * Example that illustrates the accuracy of the text extraction, by first extracting
 * all TextMarks and then reconstructing the text by writing out the text page-by-page
 * to a new PDF with the creator package.
 * Only retains the text.
 *
 * Useful to check accuracy of text extraction properties.
 * Expands upon reconstruct_text.go to show word placements.
 *
 * Run as: go run reconstruct_words.go input.pdf
 */

package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"

	"github.com/unidoc/unipdf/v4/common/license"
	"github.com/unidoc/unipdf/v4/creator"
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

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Syntax: reconstruct_words <file.pdf>\n")
		os.Exit(1)
	}

	pdfPath := os.Args[1]
	err := reconstruct(pdfPath)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully written to reconstr_words.pdf\n")
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

	c := creator.New()

	for pageNum := 1; pageNum <= len(pdfr.PageList); pageNum++ {
		page, err := pdfr.GetPage(pageNum)
		if err != nil {
			return err
		}

		extr, err := extractor.New(page)
		if err != nil {
			return err
		}
		pageText, _, _, err := extr.ExtractPageText()
		if err != nil {
			return err
		}

		// Start on a new page.
		c.NewPage()
		fmt.Printf("Page %d\n", pageNum)

		text := pageText.Text()
		textmarks := pageText.Marks()
		fmt.Printf("%s\n", text)

		reWord := regexp.MustCompile(`(?m)\S+`)
		offsets := reWord.FindAllStringIndex(text, -1)
		wordTextMarkArrays := make([]*extractor.TextMarkArray, len(offsets))
		for i, o := range offsets {
			wordTextMarkArrays[i], err = textmarks.RangeOffset(o[0], o[1])
			if err != nil {
				return err
			}
			wordTextMarkArrays[i].BBox()
		}

		// Reconstruct the text, each single TextMark drawn at a time with creator.Paragraph.
		for _, tm := range textmarks.Elements() {
			if tm.Font == nil {
				continue
			}
			fmt.Printf("%s\n", tm.Text)
			// Reconstruct by drawing each glyph from textmarks with the creator package.
			para := c.NewParagraph(tm.Original)
			para.SetFont(tm.Font)
			para.SetFontSize(tm.FontSize)
			r, g, b, _ := tm.StrokeColor.RGBA()
			rf, gf, bf := float64(r)/0xffff, float64(g)/0xffff, float64(b)/0xffff
			para.SetColor(creator.ColorRGBFromArithmetic(rf, gf, bf))
			// Convert to PDF coordinate system.
			yPos := c.Context().PageHeight - (tm.BBox.Lly + tm.BBox.Height())
			para.SetPos(tm.BBox.Llx, yPos) // Upper left corner.
			c.Draw(para)
		}

		// Draw bounding boxes around the words identified.
		for _, wtma := range wordTextMarkArrays {
			bbox, ok := wtma.BBox()
			if !ok {
				continue
			}
			var wbuf bytes.Buffer
			for _, el := range wtma.Elements() {
				wbuf.WriteString(el.Text)
			}
			word := wbuf.String() // Might be nice to have wtma.Text() output this?

			fmt.Printf("Word: '%s' - %s - bbox: %+v\n", word, wtma.String(), bbox)
			x := bbox.Llx
			y := c.Context().PageHeight - (bbox.Lly + bbox.Height())
			rect := c.NewRectangle(x, y, bbox.Width(), bbox.Height())
			rect.SetBorderColor(creator.ColorRGBFromHex("#ff0000"))
			c.Draw(rect)
		}
	}

	return c.WriteToFile("reconstr_words.pdf")
}
