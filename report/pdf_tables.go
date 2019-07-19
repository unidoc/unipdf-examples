/*
 * This example showcases PDF tables features using unipdf's creator package.
 * The output is saved as unipdf-tables.pdf which illustrates some of the
 * features of the creator.
 */

package main

import (
	"fmt"
	"log"

	"github.com/unidoc/unipdf/v3/contentstream/draw"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

func main() {
	// For development:
	//common.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	// Create report fonts.
	font, err := model.NewStandard14Font("Helvetica")
	if err != nil {
		log.Fatal(err)
	}

	fontBold, err := model.NewStandard14Font("Helvetica-Bold")
	if err != nil {
		log.Fatal(err)
	}

	c := creator.New()
	c.SetPageMargins(50, 50, 50, 50)

	// Generate front page.
	drawFrontPage(c, font, fontBold)

	// Generate footer for pages.
	drawFooter(c, font, fontBold)

	// Customize table of contents style.
	customizeTOC(c, font, fontBold)

	// Generate basic usage chapter.
	if err := basicUsage(c, font, fontBold); err != nil {
		log.Fatal(err)
	}

	// Generate styling content chapter.
	if err := stylingContent(c, font, fontBold); err != nil {
		log.Fatal(err)
	}

	// Generate advanced usage chapter.
	if err := advancedUsage(c, font, fontBold); err != nil {
		log.Fatal(err)
	}

	// Write to output file.
	if err := c.WriteToFile("unipdf-tables.pdf"); err != nil {
		log.Fatal(err)
	}
}

func drawFrontPage(c *creator.Creator, font, fontBold *model.PdfFont) {
	c.CreateFrontPage(func(args creator.FrontpageFunctionArgs) {
		p := c.NewStyledParagraph()
		p.SetMargins(0, 0, 300, 0)
		p.SetTextAlignment(creator.TextAlignmentCenter)

		chunk := p.Append("UniPDF")
		chunk.Style.Font = font
		chunk.Style.FontSize = 56
		chunk.Style.Color = creator.ColorRGBFrom8bit(56, 68, 77)

		chunk = p.Append("\n")

		chunk = p.Append("Table features")
		chunk.Style.Font = fontBold
		chunk.Style.FontSize = 40
		chunk.Style.Color = creator.ColorRGBFrom8bit(45, 148, 215)

		c.Draw(p)
	})
}

func drawFooter(c *creator.Creator, font, fontBold *model.PdfFont) {
	c.DrawFooter(func(block *creator.Block, args creator.FooterFunctionArgs) {
		p := c.NewStyledParagraph()
		p.SetTextAlignment(creator.TextAlignmentCenter)

		chunk := p.Append(fmt.Sprintf("Page %d of %d", args.PageNum, args.TotalPages))
		chunk.Style.Font = font
		chunk.Style.FontSize = 8
		chunk.Style.Color = creator.ColorRGBFrom8bit(63, 68, 76)

		block.Draw(p)
	})
}

func customizeTOC(c *creator.Creator, font, fontBold *model.PdfFont) {
	// Enable automatic table of contents generation.
	c.AddTOC = true

	// Customize table of contents heading and its style.
	toc := c.TOC()
	toc.Heading().SetMargins(0, 0, 50, 0)

	hstyle := c.NewTextStyle()
	hstyle.Color = creator.ColorRGBFromArithmetic(0.2, 0.2, 0.2)
	hstyle.Font = fontBold
	hstyle.FontSize = 28
	toc.SetHeading("Table of Contents", hstyle)

	// Customize the style of the lines of the table of contents.
	lstyle := c.NewTextStyle()
	lstyle.Font = font
	lstyle.FontSize = 14
	toc.SetLineStyle(lstyle)
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

func advancedUsage(c *creator.Creator, font, fontBold *model.PdfFont) error {
	c.NewPage()

	// Create chapter.
	ch := c.NewChapter("Advanced usage")
	ch.SetMargins(0, 0, 50, 0)
	ch.GetHeading().SetFont(font)
	ch.GetHeading().SetFontSize(18)
	ch.GetHeading().SetColor(creator.ColorRGBFrom8bit(72, 86, 95))

	// Draw subchapters.
	columnSpan(c, ch, font, fontBold)
	tableHeaders(c, ch, font, fontBold)
	subtables(c, ch, font, fontBold)

	// Draw chapter.
	if err := c.Draw(ch); err != nil {
		return err
	}

	return nil
}

func columnSpan(c *creator.Creator, ch *creator.Chapter, font, fontBold *model.PdfFont) {
	// Create subchapter.
	sc := ch.NewSubchapter("Column span")
	sc.SetMargins(0, 0, 30, 0)
	sc.GetHeading().SetFont(font)
	sc.GetHeading().SetFontSize(13)
	sc.GetHeading().SetColor(creator.ColorRGBFrom8bit(72, 86, 95))

	// Create subchapter description.
	desc := c.NewStyledParagraph()
	desc.SetMargins(0, 0, 10, 0)
	desc.Append("Table content can be configured to span a specified number of cells.")

	sc.Add(desc)

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

	sc.Add(table)
}

func tableHeaders(c *creator.Creator, ch *creator.Chapter, font, fontBold *model.PdfFont) error {
	// Create subchapter.
	sc := ch.NewSubchapter("Headers")
	sc.SetMargins(0, 0, 30, 0)
	sc.GetHeading().SetFont(font)
	sc.GetHeading().SetFontSize(13)
	sc.GetHeading().SetColor(creator.ColorRGBFrom8bit(72, 86, 95))

	// Create subchapter description.
	desc := c.NewStyledParagraph()
	desc.SetMargins(0, 0, 10, 0)
	desc.Append("Table rows can be configured to become headers which are automatically repeated on every new page the table spans. This example also showcases the usage of images inside table cells.")

	sc.Add(desc)

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

	sc.Add(table)
	return nil
}

func subtables(c *creator.Creator, ch *creator.Chapter, font, fontBold *model.PdfFont) {
	// Create subchapter.
	sc := ch.NewSubchapter("Subtables")
	sc.SetMargins(0, 0, 30, 0)
	sc.GetHeading().SetFont(font)
	sc.GetHeading().SetFontSize(13)
	sc.GetHeading().SetColor(creator.ColorRGBFrom8bit(72, 86, 95))

	// Create subchapter description.
	desc := c.NewStyledParagraph()
	desc.SetMargins(0, 0, 10, 0)
	desc.Append("Large tables can be tedious to construct. In order to make the process more manageable, the table component allows building tables from subtables. If subtables do not fit in the current configuration of the table, the table is automatically expanded.")

	sc.Add(desc)

	// Create table.
	table := c.NewTable(6)
	table.SetMargins(0, 0, 10, 0)

	headerColor := creator.ColorRGBFrom8bit(255, 255, 0)
	footerColor := creator.ColorRGBFrom8bit(0, 255, 0)

	generateSubtable := func(rows, cols, index int, rightBorder bool) *creator.Table {
		subtable := c.NewTable(cols)

		// Add header row.
		sp := c.NewStyledParagraph()
		sp.Append(fmt.Sprintf("Header of subtable %d", index)).Style.Font = font

		cell := subtable.MultiColCell(cols)
		cell.SetContent(sp)
		cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
		cell.SetHorizontalAlignment(creator.CellHorizontalAlignmentCenter)
		cell.SetVerticalAlignment(creator.CellVerticalAlignmentMiddle)
		cell.SetBackgroundColor(headerColor)

		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				sp = c.NewStyledParagraph()
				sp.Append(fmt.Sprintf("%d-%d", i+1, j+1))
				cell = subtable.NewCell()
				cell.SetContent(sp)

				if j == 0 {
					cell.SetBorder(creator.CellBorderSideLeft, creator.CellBorderStyleSingle, 1)
				}
				if rightBorder && j == cols-1 {
					cell.SetBorder(creator.CellBorderSideRight, creator.CellBorderStyleSingle, 1)
				}
			}
		}

		// Add footer row.
		sp = c.NewStyledParagraph()
		sp.Append(fmt.Sprintf("Footer of subtable %d", index))

		cell = subtable.MultiColCell(cols)
		cell.SetContent(sp)
		cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
		cell.SetHorizontalAlignment(creator.CellHorizontalAlignmentCenter)
		cell.SetVerticalAlignment(creator.CellVerticalAlignmentMiddle)
		cell.SetBackgroundColor(footerColor)

		subtable.SetRowHeight(1, 30)
		subtable.SetRowHeight(subtable.Rows(), 40)
		return subtable
	}

	// Add subtable 1 on row 1, col 1 (4x4)
	table.AddSubtable(1, 1, generateSubtable(4, 4, 1, false))

	// Add subtable 2 on row 1, col 5 (4x4)
	// Table will be expanded to 8 columns because the subtable does not fit.
	table.AddSubtable(1, 5, generateSubtable(4, 4, 2, true))

	// Add subtable 3 on row 7, col 1 (4x4)
	table.AddSubtable(7, 1, generateSubtable(4, 4, 3, false))

	// Add subtable 4 on row 7, col 5 (4x4)
	table.AddSubtable(7, 5, generateSubtable(4, 4, 4, true))

	// Add subtable 5 on row 13, col 3 (4x4)
	table.AddSubtable(13, 3, generateSubtable(4, 4, 5, true))

	// Add subtable 6 on row 13, col 1 (3x2)
	table.AddSubtable(13, 1, generateSubtable(3, 2, 6, false))

	// Add subtable 7 on row 13, col 7 (3x2)
	table.AddSubtable(13, 7, generateSubtable(3, 2, 7, true))

	// Add subtable 8 on row 18, col 1 (3x2)
	table.AddSubtable(18, 1, generateSubtable(3, 2, 8, false))

	// Add subtable 9 on row 19, col 3 (2x4)
	table.AddSubtable(19, 3, generateSubtable(2, 4, 9, true))

	// Add subtable 10 on row 18, col 7 (3x2)
	table.AddSubtable(18, 7, generateSubtable(3, 2, 10, true))

	sc.Add(table)
}
