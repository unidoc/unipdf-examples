/*
 * Add image with alternate text to a PDF file.
 *
 * This example showcases how to add an image to a PDF file and set alternate text for the image.
 * Alternate text is useful for screen readers and other assistive technologies to provide a description of the image.
 *
 * The example demonstrates how to construct a `StructTreeRoot` object and add alternate text for images.
 *
 * Usage:
 * go run pdf_add_image_alt_text.go
 */

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
	var (
		mcid           int64 = 0
		structTreeRoot       = model.NewStructTreeRoot()
		imageFile            = "logo.png"
		kIds                 = []string{
			"86dfd8f4-b09b-41bc-981a-8b77de9aa251",
			"7c3ffc57-120a-4b44-827a-3f515ffa87b7",
			"b78b019d-cfb1-465f-b025-7dfda40d0b58",
			"06296b7a-818b-451a-9b59-5d4a20f59756",
		}
		outputPath = "image_alt_text.pdf"
	)

	c := creator.New()

	// Construct base K dictionary.
	docK := model.NewKDictionary()
	docK.S = core.MakeName(string(model.StructureTypeDocument))
	// Manually set optional ID for the K object
	docK.ID = core.MakeString(kIds[0])
	// Or generate the ID automatically using the following:
	// docK.GenerateRandomID()

	// Add K dictionary to the struct tree root.
	structTreeRoot.AddKDict(docK)

	// Create a child K dictionary.
	pageMarkedContentSection := model.NewKDictionary()

	// Set the structure type to Section.
	pageMarkedContentSection.S = core.MakeName(string(model.StructureTypeSection))
	pageMarkedContentSection.ID = core.MakeString(kIds[1])

	// Add as a child
	docK.AddKChild(pageMarkedContentSection)

	// Add first image
	err := addImage(c, imageFile, 0, 10, mcid, "An image alt text", kIds[2], pageMarkedContentSection)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Update the MCID to be used for the next image.
	mcid++

	// Add second image.
	err = addImage(c, imageFile, 0, 400, mcid, "A second image alt text", kIds[3], pageMarkedContentSection)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Set the struct tree root.
	c.SetStructTreeRoot(structTreeRoot)

	err = c.WriteToFile(outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func addImage(c *creator.Creator, imageFile string, x, y float64, mcid int64, altText, kID string, pageMarkedContentSection *model.KDict) error {
	img, err := c.NewImageFromFile(imageFile)
	if err != nil {
		return err
	}

	img.SetPos(x, y)
	err = c.Draw(img)
	if err != nil {
		return err
	}

	// Add the image to the marked content section.
	img.SetMarkedContentID(mcid)

	altKdictEntry, err := img.GenerateKDict()
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	// Set the alternate text.
	altKdictEntry.Alt = core.MakeString(altText)

	// Set the page number.
	altKdictEntry.SetPageNumber(int64(c.Context().Page))

	// Set the bounding box.
	altKdictEntry.SetBoundingBox(x, y, img.Width(), img.Height())

	// Set the ID.
	altKdictEntry.ID = core.MakeString(kID)

	// Add the K dictionary entry to the marked content section.
	pageMarkedContentSection.AddKChild(altKdictEntry)

	return nil
}
