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

	unicommon "github.com/unidoc/unipdf/v3/common"
	pdfcore "github.com/unidoc/unipdf/v3/core"
	pdf "github.com/unidoc/unipdf/v3/model"
)

func main() {
	// Debug log mode.
	//unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

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
	unicommon.Log.Debug("Input PDF: %v", inputPath)

	pdfWriter := pdf.NewPdfWriter()

	// Read the input pdf file.
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		// Read the page.
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return err
		}

		// New text annotation.
		textAnnotation := pdf.NewPdfAnnotationText()
		textAnnotation.Contents = pdfcore.MakeString(annotationText)
		// The rect specifies the location of the markup.
		textAnnotation.Rect = pdfcore.MakeArray(pdfcore.MakeInteger(20), pdfcore.MakeInteger(100), pdfcore.MakeInteger(10+50), pdfcore.MakeInteger(100+50))

		// Add to the page annotations.
		page.AddAnnotation(textAnnotation.PdfAnnotation)

		err = pdfWriter.AddPage(page)
		if err != nil {
			unicommon.Log.Error("Failed to add page: %s", err)
			return err
		}
	}

	fWrite, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer fWrite.Close()

	err = pdfWriter.Write(fWrite)
	if err != nil {
		return err
	}

	return nil
}
