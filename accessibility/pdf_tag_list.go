// This example demonstrates how to create a tagged PDF document with nested lists
// using the UniPDF library. This example shows how to properly structure a PDF
// with accessibility tags by creating a document structure tree with K dictionaries
// for hierarchical list elements.
//
// The program creates a PDF document containing:
// - A main list with two levels of nested sublists
// - Proper tagging structure for accessibility compliance
// - Custom list markers and indentation
// - Structure tree root with K dictionaries for each list level
//
// This is useful for creating accessible PDF documents that can be properly
// interpreted by screen readers and other assistive technologies.
//

package main

import (
	"fmt"
	"log"
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

	structTreeRoot := model.NewStructTreeRoot()

	// Construct base K dictionary.
	docK := model.NewKDictionary()
	docK.S = core.MakeName(string(model.StructureTypeDocument))
	docK.ID = core.MakeString(docK.GenerateRandomID())

	// Add K dictionary to the struct tree root.
	structTreeRoot.AddKDict(docK)

	list := c.NewList()

	// Tag the list and generate its K dictionary.
	// This will be used as the parent K dictionary for nested lists.
	list.AddTag(docK)
	listKDict, err := list.GenerateKDict()
	if err != nil {
		log.Fatalf("failed to generate K dictionary for list: %v", err)
	}

	currentList := list
	currentKDict := listKDict
	for i := 0; i < 2; i++ {
		l := c.NewList()
		l.Marker().Text = "- "
		l.SetIndent(10)
		l.AddTag(currentKDict)

		l.AddTextItem(fmt.Sprintf("Nesting level %d", i+1))
		l.AddTextItem(fmt.Sprintf("Nesting level %d", i+1))

		currentList.Add(l)
		currentList = l

		// Generate the K dictionary for the sublist, which will be used as the parent
		// for the next level of nesting.
		currentKDict, err = l.GenerateKDict()
		if err != nil {
			log.Fatalf("failed to generate K dictionary for sublist: %v", err)
		}
	}

	if err := c.Draw(list); err != nil {
		log.Fatalf("failed to draw list: %v", err)
	}

	// Set the struct tree root.
	c.SetStructTreeRoot(structTreeRoot)

	if err := c.WriteToFile("pdf_tag_list.pdf"); err != nil {
		log.Fatalf("failed to write pdf: %v", err)
	}
}
