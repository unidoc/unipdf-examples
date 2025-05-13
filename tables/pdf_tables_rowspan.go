/*
 * This example showcases PDF tables features using unipdf's creator package.
 * The output is saved as unipdf-rowSpan-tables.pdf which illustrates how to
 * configure rows that span multiple columns.
 */

package main

import (
	"fmt"
	"log"
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

	// Create Row Span example
	if err := rowSpan(c, font, fontBold); err != nil {
		log.Fatal(err)
	}

	// Write to output file.
	if err := c.WriteToFile("unipdf-rowSpan-tables.pdf"); err != nil {
		log.Fatal(err)
	}
}

func rowSpan(c *creator.Creator, font, fontBold *model.PdfFont) error {
	c.NewPage()
	// Create subchapter.
	ch := c.NewChapter("Row span")
	ch.SetMargins(0, 0, 30, 0)

	heading := ch.GetHeading()
	heading.SetFont(font)
	heading.SetFontSize(13)
	heading.SetFontColor(creator.ColorRGBFrom8bit(72, 86, 95))

	// Create subchapter description.
	desc := c.NewStyledParagraph()
	desc.SetMargins(0, 0, 10, 0)
	desc.Append("Table content can be configured to span a specified number of cells.")

	ch.Add(desc)

	// Create table.
	table := c.NewTable(5)
	table.SetMargins(0, 0, 10, 0)

	drawCell := func(font *model.PdfFont, rowspan int, color, bgColor creator.Color) {
		var cell *creator.TableCell

		if rowspan == 1 {
			cell = table.NewCell()
		} else {
			cell = table.MultiRowCell(rowspan)
		}

		p := c.NewStyledParagraph()
		p.SetText(fmt.Sprintf("%d - %d", table.CurRow(), table.CurCol()))

		cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
		cell.SetBackgroundColor(bgColor)
		cell.SetContent(p)
		cell.SetHorizontalAlignment(creator.CellHorizontalAlignmentCenter)
		cell.SetIndent(0)
	}

	// Draw table content.

	// Rowspan 1.
	drawCell(fontBold, 1, creator.ColorWhite, creator.ColorBlue)
	drawCell(fontBold, 1, creator.ColorWhite, creator.ColorBlue)
	drawCell(fontBold, 1, creator.ColorWhite, creator.ColorBlue)
	drawCell(fontBold, 1, creator.ColorWhite, creator.ColorBlue)
	drawCell(fontBold, 1, creator.ColorWhite, creator.ColorBlue)

	// Rowspan 1, 2, 3, 1, 1.
	drawCell(fontBold, 1, creator.ColorWhite, creator.ColorRed)
	drawCell(fontBold, 2, creator.ColorWhite, creator.ColorRed)
	drawCell(fontBold, 3, creator.ColorWhite, creator.ColorRed)
	drawCell(fontBold, 1, creator.ColorWhite, creator.ColorRed)
	drawCell(fontBold, 1, creator.ColorWhite, creator.ColorRed)

	// Rowspan 1
	drawCell(fontBold, 1, creator.ColorWhite, creator.ColorGreen)
	drawCell(fontBold, 1, creator.ColorWhite, creator.ColorGreen)
	drawCell(fontBold, 1, creator.ColorWhite, creator.ColorGreen)

	// Rowspan 1
	drawCell(fontBold, 1, creator.ColorWhite, creator.ColorYellow)
	drawCell(fontBold, 1, creator.ColorWhite, creator.ColorYellow)
	drawCell(fontBold, 1, creator.ColorWhite, creator.ColorYellow)
	drawCell(fontBold, 1, creator.ColorWhite, creator.ColorYellow)

	ch.Add(table)
	// Draw chapter.
	if err := c.Draw(ch); err != nil {
		return err
	}

	return nil
}
