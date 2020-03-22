/*
 * This example showcases PDF tables features using unipdf's creator package.
 * The output is saved as unipdf-header-imgs-tables.pdf which illustrates how to
 * add headers on top of tables. It also highlights how to add images
 */

package main

import (
	"fmt"
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
	if err := tableHeaders(c, font, fontBold); err != nil {
		log.Fatal(err)
	}

	// Write to output file.
	if err := c.WriteToFile("unipdf-header-imgs-tables.pdf"); err != nil {
		log.Fatal(err)
	}
}

func tableHeaders(c *creator.Creator, font, fontBold *model.PdfFont) error {
	// Create subchapter.
	ch := c.NewChapter("Headers")
	ch.SetMargins(0, 0, 30, 0)
	ch.GetHeading().SetFont(font)
	ch.GetHeading().SetFontSize(13)
	ch.GetHeading().SetColor(creator.ColorRGBFrom8bit(72, 86, 95))

	// Create subchapter description.
	desc := c.NewStyledParagraph()
	desc.SetMargins(0, 0, 10, 0)
	desc.Append("Table rows can be configured to become headers which are automatically repeated on every new page the table spans. This example also showcases the usage of images inside table cells.")

	ch.Add(desc)

	// Load table image.
	img, err := c.NewImageFromFile("./unidoc-logo.png")
	if err != nil {
		return err
	}
	img.SetMargins(2, 2, 2, 2)
	img.ScaleToWidth(30)

	// Create table.
	table := c.NewTable(4)
	table.SetColumnWidths(0.1, 0.3, 0.4, 0.2)
	table.SetMargins(0, 0, 10, 0)

	drawCell := func(text string, font *model.PdfFont, align creator.CellHorizontalAlignment, color, bgColor creator.Color, colspan int) {
		p := c.NewStyledParagraph()
		chunk := p.Append(text)
		chunk.Style.Font = font
		chunk.Style.Color = color

		cell := table.MultiColCell(colspan)
		cell.SetBackgroundColor(bgColor)
		cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
		cell.SetHorizontalAlignment(align)
		cell.SetContent(p)
	}

	drawCell("Header", fontBold, creator.CellHorizontalAlignmentCenter, creator.ColorWhite, creator.ColorBlue, 4)
	drawCell("This is the subheader", fontBold, creator.CellHorizontalAlignmentCenter, creator.ColorBlack, creator.ColorWhite, 4)
	table.SetHeaderRows(1, 2)

	// Draw table content.
	for i := 0; i < 62; i++ {
		num := i + 1

		color := creator.ColorBlack
		bgColor := creator.ColorWhite
		if num%2 == 0 {
			color = creator.ColorRGBFromHex("#fefefe")
			bgColor = creator.ColorRGBFromHex("#999")
		}

		// Draw image cell.
		cell := table.NewCell()
		cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
		cell.SetHorizontalAlignment(creator.CellHorizontalAlignmentCenter)
		cell.SetContent(img)

		drawCell(fmt.Sprintf("Product #%d", num), font, creator.CellHorizontalAlignmentLeft, color, bgColor, 1)
		drawCell(fmt.Sprintf("Description #%d", num), font, creator.CellHorizontalAlignmentCenter, color, bgColor, 1)
		drawCell(fmt.Sprintf("$%d", num*10), font, creator.CellHorizontalAlignmentRight, color, bgColor, 1)
	}

	ch.Add(table)

	// Draw chapter.
	if err := c.Draw(ch); err != nil {
		return err
	}

	return nil
}
