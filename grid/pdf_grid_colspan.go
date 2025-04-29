/*
 * This example showcases PDF grid features using unipdf's creator package.
 * The output is saved as unipdf-grid-colspan.pdf which illustrates how to
 * create a grid with a cell occupying several columns.
 */

package main

import (
	"log"
	"os"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
)

const loremShort = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce pulvinar consectetur augue, et molestie erat porttitor id. Integer id elementum justo. Vestibulum ut luctus arcu. Nam varius nibh non vulputate condimentum. Etiam molestie velit at ex blandit condimentum. Maecenas vulputate velit quis maximus mattis. Donec dolor velit, vehicula non est suscipit, rutrum congue tortor. Morbi facilisis sed metus non volutpat."

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

	// add cell occupying three columns
	cell, err := row.NewMultiCell(3, 1)
	if err != nil {
		log.Fatal(err)
	}

	// we set the cell content to image
	imgData, err := os.ReadFile("logo.png")
	if err != nil {
		log.Fatal(err)
	}
	img, err := c.NewImageFromData(imgData)
	if err != nil {
		log.Fatal(err)
	}

	cell.SetHorizontalAlignment(CellHorizontalAlignmentCenter)
	cell.SetVerticalAlignment(CellVerticalAlignmentMiddle)

	cell.SetContent(img)

	// add second row
	row = grid.NewRow()
	AddLoremIpsumCell(c, row)
	AddLoremIpsumCell(c, row)
	AddLoremIpsumCell(c, row)

	err = c.Draw(grid)
	if err != nil {
		log.Fatal(err)
	}

	// Write to output file.
	if err := c.WriteToFile("unipdf-grid-colspan.pdf"); err != nil {
		log.Fatal(err)
	}
}

func AddLoremIpsumCell(c *creator.Creator, row *creator.GridRow) error {
	cell, err := row.NewCell()
	if err {
		return err
	}
	// it is possible to add any content to cell
	p := c.NewStyledParagraph()
	p.SetText(loremShort)
	p.SetMargins(5, 5, 5, 5)
	p.SetFontSize(12)
	// for cell it is possible to modify border, background color and other properties
	cell.SetContent(p)
	return nil
}
