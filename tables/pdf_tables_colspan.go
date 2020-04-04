/*
 * This example showcases PDF tables features using unipdf's creator package.
 * The output is saved as unipdf-colSpan-tables.pdf which illustrates how to
 * configure columns that span multiple columns.
 */

package main

import (
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

	// Create Column Span example
	if err := columnSpan(c, font, fontBold); err != nil {
		log.Fatal(err)
	}

	// Write to output file.
	if err := c.WriteToFile("unipdf-colSpan-tables.pdf"); err != nil {
		log.Fatal(err)
	}
}

func columnSpan(c *creator.Creator, font, fontBold *model.PdfFont) error {
	c.NewPage()
	// Create subchapter.
	ch := c.NewChapter("Column span")
	ch.SetMargins(0, 0, 30, 0)
	ch.GetHeading().SetFont(font)
	ch.GetHeading().SetFontSize(13)
	ch.GetHeading().SetColor(creator.ColorRGBFrom8bit(72, 86, 95))

	// Create subchapter description.
	desc := c.NewStyledParagraph()
	desc.SetMargins(0, 0, 10, 0)
	desc.Append("Table content can be configured to span a specified number of cells.")

	ch.Add(desc)

	// Create table.
	table := c.NewTable(5)
	table.SetMargins(0, 0, 10, 0)

	drawCell := func(text string, font *model.PdfFont, colspan int, color, bgColor creator.Color) {
		p := c.NewStyledParagraph()
		p.SetMargins(2, 2, 0, 0)
		chunk := p.Append(text)
		chunk.Style.Font = font
		chunk.Style.Color = color

		cell := table.MultiColCell(colspan)
		cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
		cell.SetBackgroundColor(bgColor)
		cell.SetContent(p)
		cell.SetHorizontalAlignment(creator.CellHorizontalAlignmentCenter)
		cell.SetIndent(0)
	}

	// Draw table content.

	// Colspan 1 + 1 + 1 + 1 + 1.
	drawCell("1", fontBold, 1, creator.ColorWhite, creator.ColorBlue)
	drawCell("1", fontBold, 1, creator.ColorWhite, creator.ColorBlue)
	drawCell("1", fontBold, 1, creator.ColorWhite, creator.ColorBlue)
	drawCell("1", fontBold, 1, creator.ColorWhite, creator.ColorBlue)
	drawCell("1", fontBold, 1, creator.ColorWhite, creator.ColorBlue)

	// Colspan 2 + 3.
	drawCell("2", fontBold, 2, creator.ColorWhite, creator.ColorRed)
	drawCell("3", fontBold, 3, creator.ColorWhite, creator.ColorRed)

	// Colspan 4 + 1.
	drawCell("4", fontBold, 4, creator.ColorBlack, creator.ColorGreen)
	drawCell("1", fontBold, 1, creator.ColorBlack, creator.ColorGreen)

	// Colspan 2 + 2 + 1.
	drawCell("2", fontBold, 2, creator.ColorBlack, creator.ColorYellow)
	drawCell("2", fontBold, 2, creator.ColorBlack, creator.ColorYellow)
	drawCell("1", fontBold, 1, creator.ColorBlack, creator.ColorYellow)

	// Colspan 5.
	drawCell("5", fontBold, 5, creator.ColorWhite, creator.ColorBlack)

	// Colspan 1 + 2 + 1 + 1.
	drawCell("1", fontBold, 1, creator.ColorBlack, creator.ColorYellow)
	drawCell("2", fontBold, 2, creator.ColorBlack, creator.ColorYellow)
	drawCell("1", fontBold, 1, creator.ColorBlack, creator.ColorYellow)
	drawCell("1", fontBold, 1, creator.ColorBlack, creator.ColorYellow)

	// Colspan 1 + 4.
	drawCell("1", fontBold, 1, creator.ColorBlack, creator.ColorGreen)
	drawCell("4", fontBold, 4, creator.ColorBlack, creator.ColorGreen)

	// Colspan 3 + 2.
	drawCell("3", fontBold, 3, creator.ColorWhite, creator.ColorRed)
	drawCell("2", fontBold, 2, creator.ColorWhite, creator.ColorRed)

	// Colspan 1 + 2 + 2.
	drawCell("1", fontBold, 2, creator.ColorBlack, creator.ColorYellow)
	drawCell("2", fontBold, 2, creator.ColorBlack, creator.ColorYellow)
	drawCell("2", fontBold, 1, creator.ColorBlack, creator.ColorYellow)

	// Colspan 1 + 1 + 1 + 2.
	drawCell("2", fontBold, 1, creator.ColorWhite, creator.ColorBlue)
	drawCell("1", fontBold, 1, creator.ColorWhite, creator.ColorBlue)
	drawCell("1", fontBold, 1, creator.ColorWhite, creator.ColorBlue)
	drawCell("2", fontBold, 2, creator.ColorWhite, creator.ColorBlue)

	ch.Add(table)
	// Draw chapter.
	if err := c.Draw(ch); err != nil {
		return err
	}

	return nil
}
