/*
 * This example showcases PDF grid features using unipdf's creator package.
 * The output is saved as unipdf-grid-rowspan.pdf which illustrates how to
 * create a grid with a cell occupying several rows.
 */

package main

import (
	"log"
	"os"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
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
	cell.SetBorder(CellBorderSideAll, CellBorderStyleSingle, 1)
	cell.SetContent(p)

	AddCell("Row: 0 Cell: 1", c, row, true, false)
	AddCell("Row: 0 Cell: 2", c, row, true, false)

	// add second row
	row = grid.NewRow()
	AddCell("Row: 1 Cell: 1", c, row, true, false)
	AddCell("Row: 1 Cell: 2", c, row, true, false)

	// add third row
	row = grid.NewRow()
	AddCell("Row: 2 Cell: 1", c, row, true, false)
	AddCell("Row: 2 Cell: 2", c, row, true, false)

	err = c.Draw(grid)
	if err != nil {
		log.Fatal(err)
	}

	// Write to output file.
	if err := c.WriteToFile("unipdf-grid-rowspan.pdf"); err != nil {
		log.Fatal(err)
	}
}
