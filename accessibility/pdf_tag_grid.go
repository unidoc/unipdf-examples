// This example demonstrates how to create a tagged PDF with a grid (table)
// using the UniPDF library. The grid will be properly tagged in the document
// structure tree for accessibility compliance.
//

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v4/common/license"
	"github.com/unidoc/unipdf/v4/core"
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

	structTreeRoot := model.NewStructTreeRoot()

	// Construct base K dictionary.
	docK := model.NewKDictionary()
	docK.S = core.MakeName(string(model.StructureTypeDocument))
	docK.ID = core.MakeString(docK.GenerateRandomID())

	// Add K dictionary to the struct tree root.
	structTreeRoot.AddKDict(docK)

	drawGrid(c, docK)

	// Set the struct tree root.
	c.SetStructTreeRoot(structTreeRoot)

	err := c.WriteToFile("pdf_tag_grid.pdf")
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}
}

func drawGrid(c *creator.Creator, rootKObj *model.KDict) error {
	// Create table.
	grid := c.NewGrid(3)
	grid.SetMargins(0, 0, 10, 0)

	var currentRow *creator.GridRow

	drawCell := func(text string, align creator.CellHorizontalAlignment) {
		p := c.NewStyledParagraph()
		p.SetText(text)

		cell, err := currentRow.NewCell()
		if err != nil {
			fmt.Printf("Error creating new cell: %v\n", err)
			return
		}

		cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
		cell.SetHorizontalAlignment(align)
		cell.SetContent(p)
	}

	currentRow = grid.NewRow()
	currentRow.SetSection(creator.GridRowSectionHeader) // Set the first row as a header row.

	// Draw table header.
	drawCell("Align left", creator.CellHorizontalAlignmentLeft)
	drawCell("Align center", creator.CellHorizontalAlignmentCenter)
	drawCell("Align right", creator.CellHorizontalAlignmentRight)

	// Draw table content.
	for i := 0; i < 5; i++ {
		num := i + 1

		currentRow = grid.NewRow()

		drawCell(fmt.Sprintf("Product #%d", num), creator.CellHorizontalAlignmentLeft)
		drawCell(fmt.Sprintf("Description #%d", num), creator.CellHorizontalAlignmentCenter)
		drawCell(fmt.Sprintf("$%d", num*10), creator.CellHorizontalAlignmentRight)
	}

	// Draw tabke footer.
	currentRow = grid.NewRow()
	currentRow.SetSection(creator.GridRowSectionFooter) // Set the last row as a footer row.

	drawCell("Total", creator.CellHorizontalAlignmentLeft)
	drawCell("", creator.CellHorizontalAlignmentCenter)
	drawCell("$50", creator.CellHorizontalAlignmentRight)

	// Tag the grid with the root K dictionary.
	// This associates the grid with the document structure tree.
	// The grid will be a child of the document structure tree.
	grid.AddTag(rootKObj)

	return c.Draw(grid)
}
