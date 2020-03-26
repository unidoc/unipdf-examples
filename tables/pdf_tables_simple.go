/*
 * This example showcases PDF tables features using unipdf's creator package.
 * The output is saved as unipdf-simple-tables.pdf which illustrates how to
 * create a basic table.
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

	// Generate basic usage chapter.
	if err := basicUsage(c, font, fontBold); err != nil {
		log.Fatal(err)
	}

	// Write to output file.
	if err := c.WriteToFile("unipdf-simple-tables.pdf"); err != nil {
		log.Fatal(err)
	}
}

func basicUsage(c *creator.Creator, font, fontBold *model.PdfFont) error {
	// Create chapter.
	ch := c.NewChapter("Basic usage")
	ch.SetMargins(0, 0, 50, 0)
	ch.GetHeading().SetFont(font)
	ch.GetHeading().SetFontSize(18)
	ch.GetHeading().SetColor(creator.ColorRGBFrom8bit(72, 86, 95))

	// Draw subchapters.
	contentAlignH(c, ch, font, fontBold)
	contentAlignV(c, ch, font, fontBold)
	contentWrapping(c, ch, font, fontBold)

	// Draw chapter.
	if err := c.Draw(ch); err != nil {
		return err
	}

	return nil
}

func contentAlignH(c *creator.Creator, ch *creator.Chapter, font, fontBold *model.PdfFont) {
	// Create subchapter.
	sc := ch.NewSubchapter("Content horizontal alignment")
	sc.SetMargins(0, 0, 30, 0)
	sc.GetHeading().SetFont(font)
	sc.GetHeading().SetFontSize(13)
	sc.GetHeading().SetColor(creator.ColorRGBFrom8bit(72, 86, 95))

	// Create subchapter description.
	desc := c.NewStyledParagraph()
	desc.SetMargins(0, 0, 10, 0)
	desc.Append("Cell content can be aligned horizontally left, right or it can be centered.")

	sc.Add(desc)

	// Create table.
	table := c.NewTable(3)
	table.SetMargins(0, 0, 10, 0)

	drawCell := func(text string, font *model.PdfFont, align creator.CellHorizontalAlignment) {
		p := c.NewStyledParagraph()
		p.Append(text).Style.Font = font

		cell := table.NewCell()
		cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
		cell.SetHorizontalAlignment(align)
		cell.SetContent(p)
	}

	// Draw table header.
	drawCell("Align left", fontBold, creator.CellHorizontalAlignmentLeft)
	drawCell("Align center", fontBold, creator.CellHorizontalAlignmentCenter)
	drawCell("Align right", fontBold, creator.CellHorizontalAlignmentRight)

	// Draw table content.
	for i := 0; i < 5; i++ {
		num := i + 1

		drawCell(fmt.Sprintf("Product #%d", num), font, creator.CellHorizontalAlignmentLeft)
		drawCell(fmt.Sprintf("Description #%d", num), font, creator.CellHorizontalAlignmentCenter)
		drawCell(fmt.Sprintf("$%d", num*10), font, creator.CellHorizontalAlignmentRight)
	}

	sc.Add(table)
}

func contentAlignV(c *creator.Creator, ch *creator.Chapter, font, fontBold *model.PdfFont) {
	// Create subchapter.
	sc := ch.NewSubchapter("Content vertical alignment")
	sc.SetMargins(0, 0, 30, 0)
	sc.GetHeading().SetFont(font)
	sc.GetHeading().SetFontSize(13)
	sc.GetHeading().SetColor(creator.ColorRGBFrom8bit(72, 86, 95))

	// Create subchapter description.
	desc := c.NewStyledParagraph()
	desc.SetMargins(0, 0, 10, 0)
	desc.Append("Cell content can be positioned vertically at the top, bottom or in the middle of the cell.")

	sc.Add(desc)

	// Create table.
	table := c.NewTable(3)
	table.SetMargins(0, 0, 10, 0)

	drawCell := func(text string, font *model.PdfFont, fontSize float64, align creator.CellVerticalAlignment) {
		p := c.NewStyledParagraph()
		chunk := p.Append(text)
		chunk.Style.Font = font
		chunk.Style.FontSize = fontSize

		cell := table.NewCell()
		cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
		cell.SetVerticalAlignment(align)
		cell.SetHorizontalAlignment(creator.CellHorizontalAlignmentCenter)
		cell.SetContent(p)
	}

	// Draw table header.
	drawCell("Align top", fontBold, 10, creator.CellVerticalAlignmentMiddle)
	drawCell("Align bottom", fontBold, 10, creator.CellVerticalAlignmentMiddle)
	drawCell("Align middle", fontBold, 10, creator.CellVerticalAlignmentMiddle)

	// Draw table content.
	for i := 0; i < 5; i++ {
		num := i + 1
		fontSize := float64(num) * 2

		drawCell(fmt.Sprintf("Product #%d", num), font, fontSize, creator.CellVerticalAlignmentTop)
		drawCell(fmt.Sprintf("$%d", num*10), font, fontSize, creator.CellVerticalAlignmentBottom)
		drawCell(fmt.Sprintf("Description #%d", num), font, fontSize, creator.CellVerticalAlignmentMiddle)
	}

	sc.Add(table)
}

func contentWrapping(c *creator.Creator, ch *creator.Chapter, font, fontBold *model.PdfFont) {
	// Create subchapter.
	sc := ch.NewSubchapter("Content wrapping")
	sc.SetMargins(0, 0, 30, 0)
	sc.GetHeading().SetFont(font)
	sc.GetHeading().SetFontSize(13)
	sc.GetHeading().SetColor(creator.ColorRGBFrom8bit(72, 86, 95))

	// Create subchapter description.
	desc := c.NewStyledParagraph()
	desc.SetMargins(0, 0, 10, 0)
	desc.Append("Cell text content is automatically broken into lines, depeding on the cell size.")

	sc.Add(desc)

	// Create table.
	table := c.NewTable(4)
	table.SetColumnWidths(0.25, 0.2, 0.25, 0.3)
	table.SetMargins(0, 0, 10, 0)

	drawCell := func(text string, font *model.PdfFont, align creator.TextAlignment) {
		p := c.NewStyledParagraph()
		p.SetTextAlignment(align)
		p.SetMargins(2, 2, 0, 0)
		p.Append(text).Style.Font = font

		cell := table.NewCell()
		cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
		cell.SetContent(p)
		cell.SetIndent(0)
	}

	// Draw table header.
	drawCell("Align left", fontBold, creator.TextAlignmentLeft)
	drawCell("Align center", fontBold, creator.TextAlignmentCenter)
	drawCell("Align right", fontBold, creator.TextAlignmentRight)
	drawCell("Align justify", fontBold, creator.TextAlignmentCenter)

	// Draw table content.
	content := "Maecenas tempor nibh gravida nunc laoreet, ut rhoncus justo ultricies. Mauris nec purus sit amet purus tincidunt efficitur tincidunt non dolor. Aenean nisl eros, volutpat vitae dictum id, facilisis ac felis. Integer lacinia, turpis at fringilla posuere, erat tortor ultrices orci, non tempor neque mauris ac neque. Morbi blandit ante et lacus ornare, ut vulputate massa dictum."

	drawCell(content, font, creator.TextAlignmentLeft)
	drawCell(content, font, creator.TextAlignmentCenter)
	drawCell(content, font, creator.TextAlignmentRight)
	drawCell(content, font, creator.TextAlignmentJustify)

	sc.Add(table)
}
