/*
 * This example showcases PDF tables features using unipdf's creator package.
 * The output is saved as unipdf-subtables-tables.pdf which illustrates how to
 * configure subtables in tables
 */

package main

import (
	"fmt"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
	"log"
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
	if err := subtables(c, font, fontBold); err != nil {
		log.Fatal(err)
	}

	// Write to output file.
	if err := c.WriteToFile("unipdf-subtables-tables.pdf"); err != nil {
		log.Fatal(err)
	}
}

func subtables(c *creator.Creator, font, fontBold *model.PdfFont) error {
	// Create subchapter.
	ch := c.NewChapter("Subtables")
	ch.SetMargins(0, 0, 30, 0)
	ch.GetHeading().SetFont(font)
	ch.GetHeading().SetFontSize(13)
	ch.GetHeading().SetColor(creator.ColorRGBFrom8bit(72, 86, 95))

	// Create subchapter description.
	desc := c.NewStyledParagraph()
	desc.SetMargins(0, 0, 10, 0)
	desc.Append("Large tables can be tedious to construct. In order to make the process more manageable, the table component allows building tables from subtables. If subtables do not fit in the current configuration of the table, the table is automatically expanded.")

	ch.Add(desc)

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

	ch.Add(table)

	// Draw chapter.
	if err := c.Draw(ch); err != nil {
		return err
	}

	return nil
}
