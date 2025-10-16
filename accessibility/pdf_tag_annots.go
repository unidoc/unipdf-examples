// This example demonstrates how to create a PDF with text annotations
// using the UniPDF library. The PDF will contain multiple text annotations
// with different properties, and the annotations will be properly tagged
// in the document structure tree for accessibility compliance.
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
	outputPath := "pdf_tag_annots.pdf"

	err := createPdfWithTextAnnotations(outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

// Create a new PDF with text annotations.
func createPdfWithTextAnnotations(outputPath string) error {
	// Create a new Creator.
	c := creator.New()

	// Construct the StructTreeRoot.
	str := model.NewStructTreeRoot()

	// Construct base K dictionary.
	docK := model.NewKDictionary()
	docK.S = core.MakeName(string(model.StructureTypeDocument))

	str.AddKDict(docK)

	// Create a new page with standard letter size.
	page := model.NewPdfPage()
	mediaBox := core.MakeArrayFromFloats([]float64{0, 0, 612, 792})
	mediaRect, err := model.NewPdfRectangle(*mediaBox)
	if err != nil {
		return fmt.Errorf("failed to create MediaBox rectangle: %w", err)
	}
	page.MediaBox = mediaRect

	// Add first text annotation.
	textAnnotation1 := model.NewPdfAnnotationText()
	textAnnotation1.Contents = core.MakeString("This is a sample text annotation!")
	textAnnotation1.Rect = core.MakeArray(
		core.MakeInteger(100), // x1
		core.MakeInteger(600), // y1 (from bottom of page)
		core.MakeInteger(150), // x2 (x1 + width)
		core.MakeInteger(650), // y2 (y1 + height)
	)
	// Set annotation properties.
	textAnnotation1.Open = core.MakeBool(false) // Closed by default
	textAnnotation1.Name = core.MakeName("Comment")

	docK.AddKChild(textAnnotation1.GenerateKDict())

	// Add second text annotation.
	textAnnotation2 := model.NewPdfAnnotationText()
	textAnnotation2.Contents = core.MakeString("Another annotation with more detailed information about this section.")
	textAnnotation2.Rect = core.MakeArray(
		core.MakeInteger(300), // x1
		core.MakeInteger(550), // y1
		core.MakeInteger(350), // x2
		core.MakeInteger(600), // y2
	)
	textAnnotation2.Open = core.MakeBool(true) // Open by default
	textAnnotation2.Name = core.MakeName("Note")

	docK.AddKChild(textAnnotation2.GenerateKDict())

	// Add third text annotation.
	textAnnotation3 := model.NewPdfAnnotationText()
	textAnnotation3.Contents = core.MakeString("You can click on these yellow icons to view the annotation content.")
	textAnnotation3.Rect = core.MakeArray(
		core.MakeInteger(450), // x1
		core.MakeInteger(400), // y1
		core.MakeInteger(500), // x2
		core.MakeInteger(450), // y2
	)
	textAnnotation3.Open = core.MakeBool(false)
	textAnnotation3.Name = core.MakeName("Help")

	docK.AddKChild(textAnnotation3.GenerateKDict())

	// Add annotations to the page.
	page.AddAnnotation(textAnnotation1.PdfAnnotation)
	page.AddAnnotation(textAnnotation2.PdfAnnotation)
	page.AddAnnotation(textAnnotation3.PdfAnnotation)

	// Add StructTreeRoot to the page.
	c.SetStructTreeRoot(str)

	// Add the page to the writer.
	err := c.AddPage(page)
	if err != nil {
		return err
	}

	// Write the PDF to file.
	err = c.WriteToFile(outputPath)
	if err != nil {
		return err
	}

	return nil
}
