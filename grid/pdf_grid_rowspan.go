/*
 * This example showcases PDF grid features using unipdf's creator package.
 * The output is saved as unipdf-grid-rowspan.pdf which illustrates how to
 * create a grid with a cell occupying several rows.
 */

package main

import (
	"log"
	"os"

	"github.com/unidoc/unipdf/v4/common/license"
	"github.com/unidoc/unipdf/v4/creator"
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

	// create grid with three columns
	grid := c.NewGrid(3)

	// add first row
	row := grid.NewRow()

	// add cell which will occupy three rows
	cell, err := row.NewMultiCell(1, 3)
	p := c.NewStyledParagraph()
	p.SetText("Rowspan = 3")
	p.SetMargins(5, 5, 5, 5)
	p.SetFontSize(12)
	cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
	cell.SetContent(p)

	AddTextCell("Row: 0 Cell: 1", c, row, true, false)
	AddTextCell("Row: 0 Cell: 2", c, row, true, false)

	// add second row
	row = grid.NewRow()
	AddTextCell("Row: 1 Cell: 1", c, row, true, false)
	AddTextCell("Row: 1 Cell: 2", c, row, true, false)

	// add third row
	row = grid.NewRow()
	AddTextCell("Row: 2 Cell: 1", c, row, true, false)
	AddTextCell("Row: 2 Cell: 2", c, row, true, false)

	err = c.Draw(grid)
	if err != nil {
		log.Fatal(err)
	}

	// Write to output file.
	if err := c.WriteToFile("unipdf-grid-rowspan.pdf"); err != nil {
		log.Fatal(err)
	}
}

func AddTextCell(text string, c *creator.Creator, row *creator.GridRow, isBorder, isBackground bool) error {
	cell, err := row.NewCell()
	if err != nil {
		return err
	}
	// it is possible to add any content to cell
	p := c.NewStyledParagraph()
	p.SetText(text)
	p.SetMargins(5, 5, 5, 5)
	p.SetFontSize(14)
	// for cell it is possible to modify border, background color and other properties
	if isBorder {
		cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
	}
	if isBackground {
		cell.SetBackgroundColor(creator.ColorBlue)
	}
	cell.SetContent(p)
	return nil
}
