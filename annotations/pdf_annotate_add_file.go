/*
 * Annotate/mark up pages of a PDF file.
 * Adds a file attachment annotation placed in certain position in a page.
 *
 * Run as: go run pdf_annotate_add_file.go input.pdf output.pdf file_to_be_attached
 */

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/common/license"
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
	if len(os.Args) < 4 {
		fmt.Printf("go run pdf_annotate_add_file.go input.pdf output.pdf file_to_be_attached\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]
	attachmentPath := os.Args[3]

	err := annotatePdfAddFile(inputPath, outputPath, attachmentPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

func annotatePdfAddFile(inputPath string, outputPath string, attachmentPath string) error {
	// Read the input pdf file.
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return err
	}

	page, err := pdfReader.GetPage(1)
	if err != nil {
		return err
	}

	// Instantiate embedded file object from existing file path.
	emFile, err := model.NewEmbeddedFile(attachmentPath)
	if err != nil {
		return err
	}

	// Overwrite attachment name.
	emFile.Name = "dummy.xml"
	emFile.Description = "Sample Attachment"
	emFile.Relationship = model.RelationshipData

	creationTime := time.Now()

	// Get the page media box.
	box, err := page.GetMediaBox()
	if err != nil {
		return err
	}

	// Create a file attachment annotation definition.
	fileAnnotDef := annotator.FileAnnotationDef{
		// Position of the pin image.
		X: box.Urx - 20,
		Y: box.Ury - 20,

		// Size of the pin image.
		Width:  10,
		Height: 15,

		// Color of the pin image.
		Color: model.NewPdfColorDeviceRGB(1.0, 0, 0),

		EmbeddedFile: emFile,
		Description:  emFile.Description,
		IconName:     "Paperclip",
		CreationDate: &creationTime,
	}

	// Create a file attachment annotation.
	fileAnnot, err := annotator.CreateFileAttachmentAnnotation(fileAnnotDef)
	if err != nil {
		return err
	}

	// Add annotation to a page.
	page.AddAnnotation(fileAnnot)

	// Generate a PdfWriter instance from existing PdfReader.
	pdfWriter, err := pdfReader.ToWriter(nil)
	if err != nil {
		return err
	}

	// Write to file.
	err = pdfWriter.WriteToFile(outputPath)
	if err != nil {
		return err
	}

	return nil
}
