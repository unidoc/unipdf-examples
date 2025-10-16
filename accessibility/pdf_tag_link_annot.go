// This example demonstrates how to create accessible links in PDF documents
// following best practices for PDF/UA compliance:
//
// 1. Use tooltips (annotation Contents field) for all links
// 2. Only set Alt text when the visible text isn't descriptive enough
// 3. Properly associate annotations with structure elements
// 4. Use the same value for MCID and StructParent to ensure proper association

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
	docK.S = core.MakeName(string(model.StructureTypeDocument))

	str.AddKDict(docK)

	c.NewPage()

	// Add a title to the document
	heading := c.NewStyledParagraph()
	heading.SetMargins(0, 0, 0, 20)
	heading.SetFontSize(20)
	chunk := heading.Append("Accessible Links in PDF Documents")
	chunk.Style.FontSize = 20

	if err := c.Draw(heading); err != nil {
		log.Fatalf("failed to draw heading: %v", err)
	}

	// Example 1: Link with descriptive text (no Alt needed)
	p1 := c.NewStyledParagraph()
	p1.SetMargins(0, 0, 0, 10)
	p1.Append("Example 1: Link with descriptive text (no Alt needed): ")

	// The link text is descriptive, so no separate Alt text is needed
	// Only the tooltip is set (same as the URL for simplicity)
	url := "https://www.example.com/pdf-specification"
	linkText := "Download PDF Specification"
	tooltip := url // Use the URL as tooltip, visible on hover
	altText := ""  // No Alt needed, the link text is descriptive

	// Using mcid 1 for this element - this value is used for both the structure element's MCID
	// and the annotation's StructParent value to ensure proper structure tree association
	_, annot := p1.AddExternalLinkWithTag(linkText, url, creator.LinkTagOptions{
		Tooltip: tooltip,
		AltText: altText,
		MCID:    1,
	})

	if err := c.Draw(p1); err != nil {
		log.Fatalf("failed to draw paragraph: %v", err)
	}

	docK.AddKChild(annot)

	// Example 2: Generic link text that needs Alt text
	p2 := c.NewStyledParagraph()
	p2.SetMargins(0, 0, 0, 10)
	p2.Append("Example 2: Generic link text that needs Alt text: ")

	// "Click here" is not descriptive, so Alt text is needed
	url = "https://www.example.com/pdf-specification"
	linkText = "Click here"
	tooltip = url                          // Use the URL as tooltip, visible on hover
	altText = "Download PDF Specification" // Alt text for screen readers - important for accessibility

	// Using mcid 2 for this element
	_, annot = p2.AddExternalLinkWithTag(linkText, url, creator.LinkTagOptions{
		Tooltip: tooltip,
		AltText: altText,
		MCID:    2,
	})

	if err := c.Draw(p2); err != nil {
		log.Fatalf("failed to draw paragraph: %v", err)
	}

	docK.AddKChild(annot)

	// Example 3: Link that opens email
	p3 := c.NewStyledParagraph()
	p3.SetMargins(0, 0, 0, 10)
	p3.Append("Example 3: Email link with descriptive tooltip: ")

	url = "mailto:info@example.com"
	linkText = "Contact us"
	tooltip = "Send email to info@example.com" // More descriptive tooltip
	altText = ""                               // Link text is clear enough

	// Using mcid 3 for this element
	_, annot = p3.AddExternalLinkWithTag(linkText, url, creator.LinkTagOptions{
		Tooltip: tooltip,
		AltText: altText,
		MCID:    3,
	})

	if err := c.Draw(p3); err != nil {
		log.Fatalf("failed to draw paragraph: %v", err)
	}

	docK.AddKChild(annot)

	// Example 4: Internal link to another page (would work if document had multiple pages)
	p4 := c.NewStyledParagraph()
	p4.SetMargins(0, 0, 0, 10)
	p4.Append("Example 4: Internal link with descriptive text and tooltip: ")

	linkText = "See page 1"
	tooltip = "Jump to first page of document" // Clear tooltip explaining destination
	altText = ""                               // Link text is clear enough

	// Using mcid 4 for this element
	_, annot = p4.AddInternalLinkWithTag(linkText, 1, 50, 50, 0, creator.LinkTagOptions{
		Tooltip: tooltip,
		AltText: altText,
		MCID:    4,
	})

	if err := c.Draw(p4); err != nil {
		log.Fatalf("failed to draw paragraph: %v", err)
	}

	docK.AddKChild(annot)

	c.SetStructTreeRoot(str)

	err := c.WriteToFile("pdf_tag_link_annot.pdf")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
