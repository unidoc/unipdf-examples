/*
 *	Attach files into PDF document either by passing the file name or by passing the file content.
 */

package main

import (
	"fmt"
	"os"

	"github.com/gabriel-vasile/mimetype"
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
	inputPath := "minimal.pdf"
	outputPath := "output.pdf"
	xmlFile := "dummy.xml"

	err := addAttachment(inputPath, xmlFile, outputPath)
	if err != nil {
		fmt.Printf("Failed to add attachment: %v", err)

		os.Exit(1)
	}

	fmt.Printf("Success, output in %s\n", outputPath)
}

func addAttachment(inputPath string, attachment string, outputPath string) error {
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

	pdfWriter, err := pdfReader.ToWriter(nil)
	if err != nil {
		return err
	}

	for i := 0; i < 5; i++ {
		// Instantiate embedded file object from existing file path.
		emFile, err := model.NewEmbeddedFile(attachment)
		if err != nil {
			return err
		}

		// Overwrite attachment name.
		emFile.Name = fmt.Sprintf("dummy%d.xml", i+1)
		emFile.Description = fmt.Sprintf("Sample Attachment %d", i+1)
		emFile.Relationship = model.RelationshipData

		err = pdfWriter.AttachFile(emFile)
		if err != nil {
			return err
		}
	}

	return pdfWriter.WriteToFile(outputPath)
}

func addAttachmentFromContent(inputPath string, attachment string, outputPath string) error {
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

	pdfWriter, err := pdfReader.ToWriter(nil)
	if err != nil {
		return err
	}

	// Get attachment file content.
	content, err := os.ReadFile(attachment)
	if err != nil {
		return err
	}

	for i := 0; i < 2; i++ {
		// Instantiate embedded file object from existing file content.
		eFile, err := model.NewEmbeddedFileFromContent(content)
		if err != nil {
			return err
		}

		eFile.Name = fmt.Sprintf("dummy%d.xml", i+1)
		eFile.Description = fmt.Sprintf("Sample Attachment %d", i+1)
		eFile.FileType = mimetype.Detect(content).String()
		eFile.Relationship = model.RelationshipData

		err = pdfWriter.AttachFile(eFile)
		if err != nil {
			return err
		}
	}

	return pdfWriter.WriteToFile(outputPath)
}
