/*
 * This example showcases PDF grid features using unipdf's creator package.
 * The output is saved as unipdf-grid-wrapping.pdf which illustrates how to
 * create a grid with a row automatically wrapped between pages.
 */

package main

import (
	"fmt"
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

	// create grid with four columns
	grid := c.NewGrid(4)
	// add 25 rows
	for i := 0; i < 25; i++ {
		row := grid.NewRow()
		for j := 0; j < 4; j++ {
			cell, err := row.NewCell()
			if err != nil {
				log.Fatal(err)
			}

			p := c.NewStyledParagraph()

			p.SetText(fmt.Sprintf("Row: %d Cell: %d", cell.row, cell.col))
			p.SetMargins(5, 5, 5, 5)
			p.SetFontSize(12)
			cell.SetBorder(CellBorderSideAll, CellBorderStyleSingle, 1)
			cell.SetContent(p)
		}
	}

	row := grid.NewRow()
	// add a row which will be wrapped due to its long content
	cell, err := row.NewMultiCell(2, 1)
	if err != nil {
		log.Fatal(err)
	}
	p := c.NewStyledParagraph()
	p.SetText(loremShort)
	p.SetMargins(5, 5, 5, 5)
	p.SetFontSize(12)
	cell.SetBorder(CellBorderSideAll, CellBorderStyleSingle, 1)
	cell.SetContent(p)

	cell, err = row.NewCell()
	if err != nil {
		log.Fatal(err)
	}
	cell.SetBorder(CellBorderSideAll, CellBorderStyleSingle, 1)
	cell, err = row.NewCell()
	if err != nil {
		log.Fatal(err)
	}
	cell.SetBorder(CellBorderSideAll, CellBorderStyleSingle, 1)

	err = c.Draw(grid)
	if err != nil {
		log.Fatal(err)
	}

	// Write to output file.
	if err := c.WriteToFile("unipdf-grid-wrapping.pdf"); err != nil {
		log.Fatal(err)
	}
}
