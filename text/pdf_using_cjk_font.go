/*
 * Insert text using a CJK font.
 *
 * Run as: go run pdf_using_cjk_font.go
 */

package main

import (
	"fmt"

	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

func main() {
	outputFile := "output.pdf"

	err := genPdfFile(outputFile)

	if err != nil {
		fmt.Println(err.Error())
	}
}

func genPdfFile(outputFile string) error {
	compositeFontRegular, err := model.NewCompositePdfFontFromTTFFile("./rounded-mplus-1p-regular.ttf")

	if err != nil {
		return err
	}

	c := creator.New()
	c.SetPageMargins(50, 50, 100, 70)
	c.EnableFontSubsetting(compositeFontRegular)

	c.CreateFrontPage(func(args creator.FrontpageFunctionArgs) {
		writeContent(c, compositeFontRegular)
	})

	return c.WriteToFile(outputFile)
}

func writeContent(c *creator.Creator, compositeFont *model.PdfFont) {
	p := c.NewParagraph("こんにちは世界")
	p.SetFont(compositeFont)
	p.SetFontSize(48)
	p.SetMargins(85, 0, 150, 0)
	p.SetColor(creator.ColorRGBFrom8bit(56, 68, 77))
	c.Draw(p)

	p = c.NewParagraph("UniPDFへようこそ")
	p.SetFont(compositeFont)
	p.SetFontSize(48)
	p.SetMargins(85, 0, 0, 0)
	p.SetColor(creator.ColorRGBFrom8bit(56, 68, 77))
	c.Draw(p)

	p = c.NewParagraph("Welcome To UniPDF")
	p.SetFont(compositeFont)
	p.SetFontSize(30)
	p.SetMargins(85, 0, 0, 0)
	p.SetColor(creator.ColorRGBFrom8bit(45, 148, 215))
	c.Draw(p)
}
