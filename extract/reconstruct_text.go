// Example that illustrates the accuracy of the text extraction, by first extracting
// all TextMarks and then reconstructing the text by writing out the text page-by-page
// to a new PDF with the creator package.
// Only retains the text.
//
// Useful to check accuracy of text extraction properties.

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

const licenseKey = `
-----BEGIN UNIDOC LICENSE KEY-----
Free trial license keys are available at: https://unidoc.io/
-----END UNIDOC LICENSE KEY-----
`

func init() {
	// Enable debug-level logging.
	// common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))

	err := license.SetLicenseKey(licenseKey, `Company Name`)
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Syntax: reconstruct_text <file.pdf>\n")
		os.Exit(1)
	}

	pdfPath := os.Args[1]
	err := reconstruct(pdfPath)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully written to reconst.pdf\n")
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

		textmarks := pageText.Marks()
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
	}

	return c.WriteToFile("reconstr.pdf")
}
