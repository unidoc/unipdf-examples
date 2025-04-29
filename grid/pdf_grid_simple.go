/*
 * This example showcases PDF grid features using unipdf's creator package.
 * The output is saved as unipdf-grid-simple.pdf which illustrates how to
 * create a basic grid.
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

	// create grid with two columns
	grid := c.NewGrid(2)

	// Creating first row:
	row := grid.NewRow()
	AddCell("Company", c, row, true, true)
	AddCell("UniDoc", c, row, true, false)

	// Creating second row:
	row = grid.NewRow()
	AddCell("Programming language", c, row, true, true)
	AddCell("Golang", c, row, true, false)

	// Creating third row:
	row = grid.NewRow()
	AddCell("Library", c, row, true, true)
	AddCell("UniPDF", c, row, true, false)

	// Creating fourth row:
	row = grid.NewRow()
	AddCell("Version", c, row, true, true)
	AddCell("4.0", c, row, true, false)

	err := c.Draw(grid)
	if err != nil {
		log.Fatal(err)
	}

	// Write to output file.
	if err := c.WriteToFile("unipdf-grid-simple.pdf"); err != nil {
		log.Fatal(err)
	}
}

func AddCell(text string, c *creator.Creator, row *creator.GridRow, isBorder, isBackground bool) error {
	cell, err := row.NewCell()
	if err {
		return err
	}
	// it is possible to add any content to cell
	p := c.NewStyledParagraph()
	p.SetText(text)
	p.SetMargins(5, 5, 5, 5)
	p.SetFontSize(14)
	// for cell it is possible to modify border, background color and other properties
	if isBorder {
		cell.SetBorder(CellBorderSideAll, CellBorderStyleSingle, 1)
	}
	if isBackground {
		cell.SetBackgroundColor(ColorBlue)
	}
	cell.SetContent(p)
	return nil
}
