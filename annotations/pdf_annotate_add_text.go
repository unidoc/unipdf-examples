/*
 * Annotate/mark up pages of a PDF file.
 * Adds a text annotation with a user specified string to a fixed location on every page.
 *
 * Run as: go run pdf_annotate_add_text.go input.pdf output.pdf text
 */

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/core"
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
		fmt.Printf("go run pdf_annotate_add_text.go input.pdf output.pdf text\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]
	annotationText := os.Args[3]

	err := annotatePdfAddText(inputPath, outputPath, annotationText)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

// Annotate pdf file.
func annotatePdfAddText(inputPath string, outputPath string, annotationText string) error {
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

	// Process each page using the following callback
	// when generating PdfWriter.
	opt := &model.ReaderToWriterOpts{
		PageProcessCallback: func(pageNum int, page *model.PdfPage) error {
			// New text annotation.
			textAnnotation := model.NewPdfAnnotationText()
			textAnnotation.Contents = core.MakeString(annotationText)
			// The rect specifies the location of the markup.
			textAnnotation.Rect = core.MakeArray(core.MakeInteger(20), core.MakeInteger(100), core.MakeInteger(10+50), pdfcore.MakeInteger(100+50))

			// Add to the page annotations.
			page.AddAnnotation(textAnnotation.PdfAnnotation)

			return nil
		},
	}

	// Generate a PdfWriter instance from existing PdfReader.
	pdfWriter, err := pdfReader.ToWriter(opt)
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
