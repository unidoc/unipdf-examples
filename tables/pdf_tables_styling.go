/*
 * This example showcases PDF tables features using unipdf's creator package.
 * The output is saved as unipdf-style-tables.pdf which illustrates how to
 * configure style in table.
 */

package main

import (
	"github.com/unidoc/unipdf/v3/contentstream/draw"
	"log"

	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

func main() {
	c := creator.New()
	c.SetPageMargins(50, 50, 50, 50)

	// Create report fonts.
	// UniPDF supports a number of font-families, which can be accessed using model.
	font, err := model.NewStandard14Font(model.HelveticaName)
	if err != nil {
		log.Fatal(err)
	}

	fontBold, err := model.NewStandard14Font(model.HelveticaBoldName)
	if err != nil {
		log.Fatal(err)
	}

	// Generate styling content chapter.
	if err := stylingContent(c, font, fontBold); err != nil {
		log.Fatal(err)
	}

	// Write to output file.
	if err := c.WriteToFile("unipdf-style-tables.pdf"); err != nil {
		log.Fatal(err)
	}
}

func stylingContent(c *creator.Creator, font, fontBold *model.PdfFont) error {
	c.NewPage()

	// Create chapter.
	ch := c.NewChapter("Styling content")
	ch.SetMargins(0, 0, 50, 0)
	ch.GetHeading().SetFont(font)
	ch.GetHeading().SetFontSize(18)
	ch.GetHeading().SetColor(creator.ColorRGBFrom8bit(72, 86, 95))

	// Draw subchapters.
	contentBorders(c, ch, font, fontBold)
	contentBackground(c, ch, font, fontBold)

	// Draw chapter.
	if err := c.Draw(ch); err != nil {
		return err
	}

	return nil
}

func contentBorders(c *creator.Creator, ch *creator.Chapter, font, fontBold *model.PdfFont) {
	// Create subchapter.
	sc := ch.NewSubchapter("Cell borders")
	sc.SetMargins(0, 0, 30, 0)
	sc.GetHeading().SetFont(font)
	sc.GetHeading().SetFontSize(13)
	sc.GetHeading().SetColor(creator.ColorRGBFrom8bit(72, 86, 95))

	// Create subchapter description.
	desc := c.NewStyledParagraph()
	desc.SetMargins(0, 0, 10, 0)
	desc.Append("Customizable cell border properties:\n\n")
	desc.Append("\u2022 Border side: left, right, top, bottom, all\n")
	desc.Append("\u2022 Border style: single or double\n")
	desc.Append("\u2022 Border line style: solid or dashed\n")
	desc.Append("\u2022 Border color\n")
	desc.Append("\u2022 Border width\n")

	sc.Add(desc)

	// Create table.
	table := c.NewTable(2)
	table.SetMargins(0, 0, 10, 0)

	drawCell := func(text string, font *model.PdfFont, borderStyle creator.CellBorderStyle, borderSide creator.CellBorderSide, borderWidth float64, borderColor creator.Color, lineStyle draw.LineStyle) {
		p := c.NewStyledParagraph()
		chunk := p.Append(text)
		chunk.Style.Font = font

		cell := table.NewCell()
		cell.SetBorder(borderSide, borderStyle, borderWidth)
		cell.SetBorderColor(borderColor)
		cell.SetBorderLineStyle(lineStyle)
		cell.SetHorizontalAlignment(creator.CellHorizontalAlignmentCenter)
		cell.SetContent(p)
	}

	// Draw table header.
	drawCell("Border right single width 2", fontBold, creator.CellBorderStyleSingle, creator.CellBorderSideRight, 2, creator.ColorBlack, draw.LineStyleSolid)
	drawCell("Border right double", fontBold, creator.CellBorderStyleDouble, creator.CellBorderSideRight, 1, creator.ColorBlack, draw.LineStyleSolid)

	// Draw table content.
	drawCell("Border top single width 2", fontBold, creator.CellBorderStyleSingle, creator.CellBorderSideTop, 2, creator.ColorBlack, draw.LineStyleSolid)
	drawCell("Border bottom single width 2", fontBold, creator.CellBorderStyleSingle, creator.CellBorderSideBottom, 2, creator.ColorBlack, draw.LineStyleSolid)

	drawCell("Border all double", fontBold, creator.CellBorderStyleDouble, creator.CellBorderSideAll, 1, creator.ColorBlack, draw.LineStyleSolid)
	drawCell("No border", fontBold, creator.CellBorderStyleSingle, creator.CellBorderSideAll, 0, creator.ColorBlack, draw.LineStyleSolid)

	drawCell("Border bottom single green", fontBold, creator.CellBorderStyleSingle, creator.CellBorderSideBottom, 1, creator.ColorGreen, draw.LineStyleSolid)
	drawCell("Border top double red", fontBold, creator.CellBorderStyleDouble, creator.CellBorderSideTop, 1, creator.ColorRed, draw.LineStyleSolid)

	drawCell("Border all single yellow", fontBold, creator.CellBorderStyleSingle, creator.CellBorderSideAll, 1, creator.ColorYellow, draw.LineStyleSolid)
	drawCell("Border right double dashed", fontBold, creator.CellBorderStyleDouble, creator.CellBorderSideRight, 1, creator.ColorBlack, draw.LineStyleDashed)

	drawCell("Border bottom double solid", fontBold, creator.CellBorderStyleDouble, creator.CellBorderSideBottom, 1, creator.ColorBlack, draw.LineStyleSolid)
	drawCell("Border bottom double dashed green", fontBold, creator.CellBorderStyleDouble, creator.CellBorderSideBottom, 1, creator.ColorGreen, draw.LineStyleDashed)

	drawCell("Border left double blue", fontBold, creator.CellBorderStyleDouble, creator.CellBorderSideLeft, 1, creator.ColorBlue, draw.LineStyleSolid)
	drawCell("Border right double red", fontBold, creator.CellBorderStyleDouble, creator.CellBorderSideRight, 1, creator.ColorRed, draw.LineStyleSolid)

	drawCell("Border all double yellow", fontBold, creator.CellBorderStyleDouble, creator.CellBorderSideAll, 1, creator.ColorYellow, draw.LineStyleSolid)
	drawCell("Border all double dashed blue", fontBold, creator.CellBorderStyleDouble, creator.CellBorderSideAll, 1, creator.ColorBlue, draw.LineStyleDashed)

	sc.Add(table)
}

func contentBackground(c *creator.Creator, ch *creator.Chapter, font, fontBold *model.PdfFont) {
	// Create subchapter.
	sc := ch.NewSubchapter("Cell background")
	sc.SetMargins(0, 0, 30, 0)
	sc.GetHeading().SetFont(font)
	sc.GetHeading().SetFontSize(13)
	sc.GetHeading().SetColor(creator.ColorRGBFrom8bit(72, 86, 95))

	// Create subchapter description.
	desc := c.NewStyledParagraph()
	desc.SetMargins(0, 0, 10, 0)
	desc.Append("The background color of the cells is also customizable.")

	sc.Add(desc)

	// Create table.
	table := c.NewTable(4)
	table.SetMargins(0, 0, 10, 0)

	drawCell := func(text string, font *model.PdfFont, bgColor creator.Color) {
		p := c.NewStyledParagraph()
		p.SetMargins(2, 2, 0, 0)
		chunk := p.Append(text)
		chunk.Style.Font = font
		chunk.Style.Color = creator.ColorWhite

		cell := table.NewCell()
		cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
		cell.SetBackgroundColor(bgColor)
		cell.SetContent(p)
		cell.SetHorizontalAlignment(creator.CellHorizontalAlignmentCenter)
		cell.SetIndent(0)
	}

	// Draw table content.
	for i := 0; i < 15; i++ {
		drawCell("Content", fontBold, creator.ColorRGBFrom8bit(byte(i*20), byte(i*7), byte(i*4)))
		drawCell("Content", fontBold, creator.ColorRGBFrom8bit(byte(i*10), byte(i*20), byte(i*4)))
		drawCell("Content", fontBold, creator.ColorRGBFrom8bit(byte(i*15), byte(i*6), byte(i*9)))
		drawCell("Content", fontBold, creator.ColorRGBFrom8bit(byte(i*6), byte(i*7), byte(i*25)))
	}

	sc.Add(table)
}
