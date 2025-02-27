/**
 * This example demonstrates how to set the language identifier for the document and its content.
 * The language identifier is used to specify the language of the document and its content.
 * It is useful for screen readers and other assistive technologies to provide the correct pronunciation and interpretation of the content.
 *
 * Usage:
 * go run pdf_set_language_identifier.go
 */

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
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
	// Set the language identifier for the document to US English.
	c.SetLanguage("en-US")

	// Set the marked content flag for the document.
	c.SetPdfWriterAccessFunc(func(w *model.PdfWriter) error {
		w.SetCatalogMarkInfo(core.MakeDictMap(map[string]core.PdfObject{
			"Marked": core.MakeBool(true),
		}))

		return nil
	})

	// Construct the StructTreeRoot.
	str := model.NewStructTreeRoot()

	// Construct base K dictionary.
	docK := model.NewKDictionary()
	docK.S = core.MakeName(model.StructureTypeDocument)

	str.AddKDict(docK)

	// Create a new paragraph.
	p := c.NewParagraph("Hello World")
	p.SetPos(100, 100)

	// Set marked content identifier for the paragraph.
	pMarkedContent := p.SetMarkedContentID(0)

	// Set the language identifier for the paragraph to British English.
	pMarkedContent.Lang = core.MakeString("en-GB")

	// Set the page number for the paragraph.
	pMarkedContent.SetPageNumber(int64(c.Context().Page))

	c.Draw(p)

	docK.AddKChild(pMarkedContent)

	// Create a new styled paragraph.
	sp := c.NewStyledParagraph()

	// Set the text for the styled paragraph.
	// It's "Hello World" in Indonesian.
	sp.SetText("Halo Dunia")
	sp.SetPos(100, 200)
	pMarkedContent = sp.SetMarkedContentID(1)

	// Set the language identifier for the styled paragraph to Indonesian.
	pMarkedContent.Lang = core.MakeString("id-ID")
	pMarkedContent.SetPageNumber(int64(c.Context().Page))

	c.Draw(sp)

	docK.AddKChild(pMarkedContent)

	c.SetStructTreeRoot(str)

	err := c.WriteToFile("pdf_set_language_identifier.pdf")
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}
}
