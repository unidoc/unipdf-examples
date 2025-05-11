/*
 * Insert unicode text and currency symbols.
 *
 * Run as: go run pdf_using_unicode_font.go
 */

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v4/common/license"
	"github.com/unidoc/unipdf/v4/creator"
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

	// Subset the font.
	// Composite fonts usually quite big and in turn would enlarge the document size if we embed all the runes/glyphs
	// This setting will embed the runes/glyphs that are used in the document
	// and in turn would reduce the document size by a lot.
	//
	// For example, in this case, the PDF file output size is reduced from 1.4 MB to 74 KB
	c.EnableFontSubsetting(compositeFontRegular)

	c.CreateFrontPage(func(args creator.FrontpageFunctionArgs) {
		writeContent(c, compositeFontRegular)
	})

	addCurrencyPage(c, compositeFontRegular)

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

func addCurrencyPage(c *creator.Creator, compositeFont *model.PdfFont) {
	c.NewPage()

	currencyText := "\u00a3 (GBP - Pound Sterling)\n" +
		"\u20ac (EUR - Euro)\n" +
		"\u20b9 (INR - Indian Rupee)\n" +
		"\u20aa (ILS - New Israeli Shekel)\n" +
		"\u20a9 (KRW - Won)\n" +
		"\u002e\u0645\u002e\u062f\u002e (MAD - Moroccan Dirham)\n" +
		"\u20b1 (PHP - Philippine Peso)\n" +
		"\uFDFC (SAR - Saudi Riyal)\n" +
		"\u0e3f (THB - Baht)\n" +
		"\u20ba (TRY - Turkish Lira)\n" +
		"\u5143 (TWD - New Taiwan Dollar)\n" +
		"\u20ab (VND - Dong)\n"

	p := c.NewParagraph(currencyText)
	p.SetFont(compositeFont)
	p.SetFontSize(20)
	p.SetMargins(85, 0, 150, 0)
	p.SetColor(creator.ColorRGBFrom8bit(56, 68, 77))
	c.Draw(p)
}
