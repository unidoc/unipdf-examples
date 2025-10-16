// This example demonstrates how to create a tagged PDF with a table
// using the UniPDF library. The table will be properly tagged in the document
// structure tree for accessibility compliance.
//
// The program creates a PDF document containing:
// - A table with three columns and multiple rows
// - Proper tagging structure for accessibility compliance
// - Custom cell alignment and borders
// - Structure tree root with K dictionaries for the table and its cells
//
// This is useful for creating accessible PDF documents that can be properly
// interpreted by screen readers and other assistive technologies.

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

	drawTable(c, docK)

	// Set the struct tree root.
	c.SetStructTreeRoot(structTreeRoot)

	err := c.WriteToFile("pdf_tag_table.pdf")
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}
}

func drawTable(c *creator.Creator, rootKObj *model.KDict) error {
	// Create table.
	table := c.NewTable(3)
	table.SetMargins(0, 0, 10, 0)

	drawCell := func(text string, align creator.CellHorizontalAlignment) {
		p := c.NewStyledParagraph()
		p.SetText(text)

		cell := table.NewCell()
		cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
		cell.SetHorizontalAlignment(align)
		cell.SetContent(p)
	}

	// Draw table header.
	drawCell("Align left", creator.CellHorizontalAlignmentLeft)
	drawCell("Align center", creator.CellHorizontalAlignmentCenter)
	drawCell("Align right", creator.CellHorizontalAlignmentRight)

	// Draw table content.
	for i := 0; i < 5; i++ {
		num := i + 1

		drawCell(fmt.Sprintf("Product #%d", num), creator.CellHorizontalAlignmentLeft)
		drawCell(fmt.Sprintf("Description #%d", num), creator.CellHorizontalAlignmentCenter)
		drawCell(fmt.Sprintf("$%d", num*10), creator.CellHorizontalAlignmentRight)
	}

	// Tag the table with the root K dictionary.
	// This associates the table with the document structure tree.
	// The table will be a child of the document structure tree.
	table.AddTag(rootKObj)

	return c.Draw(table)
}
