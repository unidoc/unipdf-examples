/*
 * Insert unicode text and currency symbols.
 *
 * Run as: go run pdf_using_unicode_font.go
 */

package main

import (
	"fmt"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

const licenseKey = `
-----BEGIN UNIDOC LICENSE KEY-----
Free trial license keys are available at: https://unidoc.io/
-----END UNIDOC LICENSE KEY-----
`

func init() {
	// Enable debug-level logging.
	// unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	err := license.SetLicenseKey(licenseKey, `Company Name`)
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
