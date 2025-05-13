/*
 * This example showcases PDF report generation in landscape mode with unidoc's creator package.
 * The output is saved as unidoc-report-landscape.pdf which illustrates some of the features
 * of the creator.
 */

/*
 * NOTE: This example depends on github.com/boombuler/barcode, MIT licensed,
 *       and github.com/unidoc/unichart, MIT licensed,
 *       and the Roboto font (Roboto-Bold.ttf, Roboto-Regular.ttf), Apache-2 licensed.
 */

package main

import (
	"fmt"
	"image"
	"math"
	"os"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"

	"github.com/unidoc/unichart"
	"github.com/unidoc/unichart/dataset"
	"github.com/unidoc/unipdf/v4/common"
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
	err := RunPdfReport("unidoc-report-landscape.pdf")
	if err != nil {
		panic(err)
	}
}

func RunPdfReport(outputPath string) error {
	robotoFontRegular, err := model.NewPdfFontFromTTFFile("./Roboto-Regular.ttf")
	if err != nil {
		return err
	}

	robotoFontPro, err := model.NewPdfFontFromTTFFile("./Roboto-Bold.ttf")
	if err != nil {
		return err
	}

	c := creator.New()

	// Switch the page size constant from the available predefined size for a landscape page
	c.SetPageSize(creator.PageSize{creator.PageSizeA4[1], creator.PageSizeA4[0]})

	c.SetPageMargins(50, 50, 100, 70)

	// Generate the table of contents.
	c.AddTOC = true
	toc := c.TOC()
	hstyle := c.NewTextStyle()
	hstyle.Color = creator.ColorRGBFromArithmetic(0.2, 0.2, 0.2)
	hstyle.FontSize = 28
	toc.SetHeading("Table of Contents", hstyle)
	lstyle := c.NewTextStyle()
	lstyle.FontSize = 14
	toc.SetLineStyle(lstyle)

	logoImg, err := c.NewImageFromFile("./unidoc-logo.png")
	if err != nil {
		return err
	}

	logoImg.ScaleToHeight(25)
	logoImg.SetPos(58, 20)

	DoDocumentControl(c, robotoFontRegular, robotoFontPro)

	DoFeatureOverview(c, robotoFontRegular, robotoFontPro)

	// Setup a front page (always placed first).
	c.CreateFrontPage(func(args creator.FrontpageFunctionArgs) {
		DoFirstPage(c, robotoFontRegular, robotoFontPro)
	})

	// Draw a header on each page.
	c.DrawHeader(func(block *creator.Block, args creator.HeaderFunctionArgs) {
		// Draw the header on a block. The block size is the size of the page's top margins.
		block.Draw(logoImg)
	})

	// Draw footer on each page.
	c.DrawFooter(func(block *creator.Block, args creator.FooterFunctionArgs) {
		// Draw the on a block for each page.
		p := c.NewStyledParagraph()
		p.SetText("unidoc.io")
		p.SetFont(robotoFontRegular)
		p.SetFontSize(8)
		p.SetPos(50, 20)
		p.SetFontColor(creator.ColorRGBFrom8bit(63, 68, 76))
		block.Draw(p)

		strPage := fmt.Sprintf("Page %d of %d", args.PageNum, args.TotalPages)
		p = c.NewStyledParagraph()
		p.SetText(strPage)
		p.SetFont(robotoFontRegular)
		p.SetFontSize(8)
		p.SetPos(750, 20)
		p.SetFontColor(creator.ColorRGBFrom8bit(63, 68, 76))
		block.Draw(p)
	})

	err = c.WriteToFile(outputPath)
	if err != nil {
		return err
	}

	return nil
}

// Generates the front page.
func DoFirstPage(c *creator.Creator, fontRegular *model.PdfFont, fontBold *model.PdfFont) {
	helvetica, _ := model.NewStandard14Font("Helvetica")
	helveticaBold, _ := model.NewStandard14Font("Helvetica-Bold")

	p := c.NewStyledParagraph()
	p.SetText("UniDoc")
	p.SetFont(helvetica)
	p.SetFontSize(48)
	p.SetMargins(85, 0, 150, 0)
	p.SetFontColor(creator.ColorRGBFrom8bit(56, 68, 77))
	c.Draw(p)

	p = c.NewStyledParagraph()
	p.SetText("Example Report")
	p.SetFont(helveticaBold)
	p.SetFontSize(30)
	p.SetMargins(85, 0, 0, 0)
	p.SetFontColor(creator.ColorRGBFrom8bit(45, 148, 215))
	c.Draw(p)

	t := time.Now().UTC()
	dateStr := t.Format("1 Jan, 2006 15:04")

	p = c.NewStyledParagraph()
	p.SetText(dateStr)
	p.SetFont(helveticaBold)
	p.SetFontSize(12)
	p.SetMargins(90, 0, 5, 0)
	p.SetFontColor(creator.ColorRGBFrom8bit(56, 68, 77))
	c.Draw(p)
}

// Document control page.
func DoDocumentControl(c *creator.Creator, fontRegular *model.PdfFont, fontBold *model.PdfFont) {
	ch := c.NewChapter("Document control")
	ch.SetMargins(0, 0, 40, 0)

	heading := ch.GetHeading()
	heading.SetFont(fontRegular)
	heading.SetFontSize(18)
	heading.SetFontColor(creator.ColorRGBFrom8bit(72, 86, 95))

	sc := ch.NewSubchapter("Issuer details")

	heading = ch.GetHeading()
	heading.SetFont(fontRegular)
	heading.SetFontSize(18)
	heading.SetFontColor(creator.ColorRGBFrom8bit(72, 86, 95))

	issuerTable := c.NewTable(2)
	issuerTable.SetMargins(0, 0, 30, 0)

	pColor := creator.ColorRGBFrom8bit(72, 86, 95)
	bgColor := creator.ColorRGBFrom8bit(56, 68, 67)

	p := c.NewStyledParagraph()
	p.SetText("Issuer")
	p.SetFont(fontBold)
	p.SetFontSize(10)
	p.SetFontColor(creator.ColorWhite)

	cell := issuerTable.NewCell()
	cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
	cell.SetBackgroundColor(bgColor)
	cell.SetContent(p)

	p = c.NewStyledParagraph()
	p.SetText("UniDoc")
	p.SetFont(fontRegular)
	p.SetFontSize(10)
	p.SetFontColor(pColor)

	cell = issuerTable.NewCell()
	cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
	cell.SetContent(p)

	p = c.NewStyledParagraph()
	p.SetText("Address")
	p.SetFont(fontBold)
	p.SetFontColor(creator.ColorWhite)

	cell = issuerTable.NewCell()
	cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
	cell.SetBackgroundColor(bgColor)
	cell.SetContent(p)

	p = c.NewStyledParagraph()
	p.SetText("Klapparstig 16, 101 Reykjavik, Iceland")
	p.SetFont(fontRegular)
	p.SetFontSize(10)
	p.SetFontColor(pColor)
	cell = issuerTable.NewCell()
	cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
	cell.SetContent(p)

	p = c.NewStyledParagraph()
	p.SetText("Email")
	p.SetFont(fontBold)
	p.SetFontSize(10)
	p.SetFontColor(creator.ColorWhite)

	cell = issuerTable.NewCell()
	cell.SetBackgroundColor(bgColor)
	cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
	cell.SetContent(p)

	p = c.NewStyledParagraph()
	p.SetText("sales@unidoc.io")
	p.SetFont(fontRegular)
	p.SetFontSize(10)
	p.SetFontColor(pColor)

	cell = issuerTable.NewCell()
	cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
	cell.SetContent(p)

	p = c.NewStyledParagraph()
	p.SetText("Web")
	p.SetFont(fontBold)
	p.SetFontSize(10)
	p.SetFontColor(creator.ColorWhite)

	cell = issuerTable.NewCell()
	cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
	cell.SetBackgroundColor(bgColor)
	cell.SetContent(p)

	p = c.NewStyledParagraph()
	p.SetText("unidoc.io")
	p.SetFont(fontRegular)
	p.SetFontSize(10)
	p.SetFontColor(pColor)

	cell = issuerTable.NewCell()
	cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
	cell.SetContent(p)

	p = c.NewStyledParagraph()
	p.SetText("Author")
	p.SetFont(fontBold)
	p.SetFontSize(10)
	p.SetFontColor(creator.ColorWhite)

	cell = issuerTable.NewCell()
	cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
	cell.SetBackgroundColor(bgColor)
	cell.SetContent(p)

	p = c.NewStyledParagraph()
	p.SetText("UniDoc report generator")
	p.SetFont(fontRegular)
	p.SetFontSize(10)
	p.SetFontColor(pColor)

	cell = issuerTable.NewCell()
	cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
	cell.SetContent(p)

	sc.Add(issuerTable)

	// 1.2 - Document history
	sc = ch.NewSubchapter("Document History")
	sc.SetMargins(0, 0, 5, 0)

	heading = sc.GetHeading()
	heading.SetFont(fontRegular)
	heading.SetFontSize(18)
	heading.SetFontColor(pColor)

	histTable := c.NewTable(3)
	histTable.SetMargins(0, 0, 30, 50)

	histCols := []string{"Date Issued", "UniDoc Version", "Type/Change"}
	for _, histCol := range histCols {
		p = c.NewStyledParagraph()
		p.SetText(histCol)
		p.SetFont(fontBold)
		p.SetFontSize(10)
		p.SetFontColor(creator.ColorWhite)

		cell = histTable.NewCell()
		cell.SetBackgroundColor(bgColor)
		cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
		cell.SetHorizontalAlignment(creator.CellHorizontalAlignmentCenter)
		cell.SetVerticalAlignment(creator.CellVerticalAlignmentMiddle)
		cell.SetContent(p)
	}

	dateStr := common.ReleasedAt.Format("1 Jan, 2006 15:04")

	histVals := []string{dateStr, common.Version, "First issue"}
	for _, histVal := range histVals {
		p = c.NewStyledParagraph()
		p.SetText(histVal)
		p.SetFont(fontRegular)
		p.SetFontSize(10)
		p.SetFontColor(pColor)

		cell = histTable.NewCell()
		cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
		cell.SetHorizontalAlignment(creator.CellHorizontalAlignmentCenter)
		cell.SetVerticalAlignment(creator.CellVerticalAlignmentMiddle)
		cell.SetContent(p)
	}

	sc.Add(histTable)

	err := c.Draw(ch)
	if err != nil {
		panic(err)
	}
}

// Chapter giving an overview of features.
// TODO: Add code snippets and show more styles and options.
func DoFeatureOverview(c *creator.Creator, fontRegular *model.PdfFont, fontBold *model.PdfFont) {
	// Ensure that the chapter starts on a new page.
	c.NewPage()

	ch := c.NewChapter("Feature overview")

	chapterFont := fontRegular
	chapterFontColor := creator.ColorRGBFrom8bit(72, 86, 95)
	chapterFontSize := 18.0

	normalFont := fontRegular
	normalFontColor := creator.ColorRGBFrom8bit(72, 86, 95)
	normalFontSize := 10.0

	bgColor := creator.ColorRGBFrom8bit(56, 68, 67)

	heading := ch.GetHeading()
	heading.SetFont(chapterFont)
	heading.SetFontSize(chapterFontSize)
	heading.SetFontColor(chapterFontColor)

	p := c.NewStyledParagraph()
	p.SetText("This chapter demonstrates a few of the features of UniDoc that can be used for report generation.")
	p.SetFont(normalFont)
	p.SetFontSize(normalFontSize)
	p.SetFontColor(normalFontColor)
	p.SetMargins(0, 0, 5, 0)
	ch.Add(p)

	// Paragraphs.
	sc := ch.NewSubchapter("Paragraphs")

	heading = sc.GetHeading()
	heading.SetMargins(0, 0, 20, 0)
	heading.SetFont(chapterFont)
	heading.SetFontSize(chapterFontSize)
	heading.SetFontColor(chapterFontColor)

	p = c.NewStyledParagraph()
	p.SetText("Paragraphs are used to represent text, as little as a single character, a word or " +
		"multiple words forming multiple sentences. UniDoc handles automatically wrapping those across lines and pages, making " +
		"it relatively easy to work with. They can also be left, center, right aligned or justified as illustrated below:")
	p.SetFont(normalFont)
	p.SetFontSize(normalFontSize)
	p.SetFontColor(normalFontColor)
	p.SetMargins(0, 0, 5, 0)
	sc.Add(p)

	// Example paragraphs:
	loremTxt := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt" +
		"ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut " +
		"aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore" +
		"eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt " +
		"mollit anim id est laborum."
	alignments := []creator.TextAlignment{creator.TextAlignmentLeft, creator.TextAlignmentCenter,
		creator.TextAlignmentRight, creator.TextAlignmentJustify}
	for j := 0; j < 4; j++ {
		p = c.NewStyledParagraph()
		p.SetText(loremTxt)
		p.SetFont(normalFont)
		p.SetFontSize(normalFontSize)
		p.SetFontColor(normalFontColor)
		p.SetMargins(20, 0, 10, 10)
		p.SetTextAlignment(alignments[j%4])

		sc.Add(p)
	}

	sc = ch.NewSubchapter("Tables")
	// Mock table: Priority table.
	priTable := c.NewTable(2)
	priTable.SetMargins(40, 40, 10, 0)
	// Column headers:
	tableCols := []string{"Priority", "Items fulfilled / available"}
	for _, tableCol := range tableCols {
		p = c.NewStyledParagraph()
		p.SetText(tableCol)
		p.SetFont(fontBold)
		p.SetFontSize(10)
		p.SetFontColor(creator.ColorWhite)

		cell := priTable.NewCell()
		cell.SetBackgroundColor(bgColor)
		cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
		cell.SetContent(p)
	}
	items := [][]string{
		[]string{"High", "52/80"},
		[]string{"Medium", "32/100"},
		[]string{"Low", "10/90"},
	}
	for _, lineItems := range items {
		for _, item := range lineItems {
			p = c.NewStyledParagraph()
			p.SetText(item)
			p.SetFont(fontBold)
			p.SetFontSize(10)
			p.SetFontColor(creator.ColorWhite)

			cell := priTable.NewCell()
			cell.SetBackgroundColor(bgColor)
			cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
			cell.SetContent(p)
		}
	}
	sc.Add(priTable)

	sc = ch.NewSubchapter("Images")

	heading = sc.GetHeading()
	heading.SetMargins(0, 0, 20, 0)
	heading.SetFont(chapterFont)
	heading.SetFontSize(chapterFontSize)
	heading.SetFontColor(chapterFontColor)

	p = c.NewStyledParagraph()
	p.SetText("Images can be loaded from multiple file formats, example from a PNG image:")
	p.SetFont(normalFont)
	p.SetFontSize(normalFontSize)
	p.SetFontColor(normalFontColor)
	p.SetMargins(0, 0, 5, 5)
	sc.Add(p)

	// Show logo.
	img, err := c.NewImageFromFile("./unidoc-logo.png")
	if err != nil {
		panic(err)
	}
	img.ScaleToHeight(50)
	sc.Add(img)

	sc = ch.NewSubchapter("QR Codes / Barcodes")

	heading = sc.GetHeading()
	heading.SetMargins(0, 0, 20, 0)
	heading.SetFont(chapterFont)
	heading.SetFontSize(chapterFontSize)
	heading.SetFontColor(chapterFontColor)

	p = c.NewStyledParagraph()
	p.SetText("Example of a QR code generated with package github.com/boombuler/barcode:")
	p.SetFont(normalFont)
	p.SetFontSize(normalFontSize)
	p.SetFontColor(normalFontColor)
	p.SetMargins(0, 0, 5, 5)
	sc.Add(p)

	qrCode, _ := makeQrCodeImage("HELLO", 40, 5)
	img, err = c.NewImageFromGoImage(qrCode)
	if err != nil {
		panic(err)
	}
	img.SetWidth(40)
	img.SetHeight(40)
	sc.Add(img)

	sc = ch.NewSubchapter("Graphing / Charts")

	heading = sc.GetHeading()
	heading.SetMargins(0, 0, 20, 0)
	heading.SetFont(chapterFont)
	heading.SetFontSize(chapterFontSize)
	heading.SetFontColor(chapterFontColor)

	p = c.NewStyledParagraph()
	p.SetText("Graphs can be generated via packages such as github.com/unidoc/unichart as illustrated " +
		"in the following plot:")
	p.SetFont(normalFont)
	p.SetFontSize(normalFontSize)
	p.SetFontColor(normalFontColor)
	p.SetMargins(0, 0, 5, 0)
	sc.Add(p)

	chart := &unichart.PieChart{
		Values: []dataset.Value{
			{Value: 70, Label: "Compliant"},
			{Value: 30, Label: "Non-Compliant"},
		},
	}
	chart.SetWidth(175)
	chart.SetHeight(175)

	// Create unipdf chart component.
	chartComponent := creator.NewChart(chart)
	sc.Add(chartComponent)

	sc = ch.NewSubchapter("Headers and footers")

	heading = sc.GetHeading()
	heading.SetMargins(0, 0, 20, 0)
	heading.SetFont(chapterFont)
	heading.SetFontSize(chapterFontSize)
	heading.SetFontColor(chapterFontColor)

	p = c.NewStyledParagraph()
	p.SetText("Convenience functions are provided to generate headers and footers, see: " +
		"https://godoc.org/github.com/unidoc/unipdf/creator#Creator.DrawHeader and " +
		"https://godoc.org/github.com/unidoc/unipdf/creator#Creator.DrawFooter " +
		"They both set a function that accepts a block which the header/footer is drawn on for each page. " +
		"More information is provided in the arguments, allowing to skip header/footer on specific pages and " +
		"showing page number and count.")
	p.SetFont(normalFont)
	p.SetFontSize(normalFontSize)
	p.SetFontColor(normalFontColor)
	p.SetMargins(0, 0, 5, 0)
	sc.Add(p)

	sc = ch.NewSubchapter("Table of contents generation")

	heading = sc.GetHeading()
	heading.SetMargins(0, 0, 20, 0)
	heading.SetFont(chapterFont)
	heading.SetFontSize(chapterFontSize)
	heading.SetFontColor(chapterFontColor)

	p = c.NewStyledParagraph()
	p.SetText("A convenience function is provided to generate table of contents " +
		"as can be seen on https://godoc.org/github.com/unidoc/unipdf/creator#Creator.CreateTableOfContents and " +
		"in our example code on unidoc.io.")
	p.SetFont(normalFont)
	p.SetFontSize(normalFontSize)
	p.SetFontColor(normalFontColor)
	p.SetMargins(0, 0, 5, 0)
	sc.Add(p)

	c.Draw(ch)
}

// Helper function to make the QR code image with a specified oversampling factor.
// The oversampling specifies how many pixels/point. Standard PDF resolution is 72 points/inch.
func makeQrCodeImage(text string, width float64, oversampling int) (image.Image, error) {
	qrCode, err := qr.Encode(text, qr.M, qr.Auto)
	if err != nil {
		return nil, err
	}

	pixelWidth := oversampling * int(math.Ceil(width))
	qrCode, err = barcode.Scale(qrCode, pixelWidth, pixelWidth)
	if err != nil {
		return nil, err
	}

	return qrCode, nil
}
