/*
 * Convert a PDF to grayscale in a vectorized fashion, including images and all content.
 *
 *
 * Run as: go run pdf_grayscale.go color.pdf output.pdf
 */

package main

import (
	"fmt"

	"os"

	"github.com/unidoc/unipdf/v4/common/license"
	"github.com/unidoc/unipdf/v4/model"
	"github.com/unidoc/unipdf/v4/pdfutil"
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
	if len(os.Args) < 3 {
		fmt.Printf("Syntax: go run pdf_grayscale_transform.go input.pdf output.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	err := convertPdfToGrayscale(inputPath, outputPath)
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Completed, see output %s\n", outputPath)
}

func convertPdfToGrayscale(inputPath, outputPath string) error {
	pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, nil)
	if err != nil {
		return err
	}
	defer f.Close()

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}
	fmt.Printf("PDF Num Pages: %d\n", numPages)

	// Process each page using the following callback
	// when generating PdfWriter.
	opts := &model.ReaderToWriterOpts{
		PageProcessCallback: func(pageNum int, page *model.PdfPage) error {
			fmt.Printf("Processing page %d/%d\n", pageNum, numPages)

			err = pdfutil.ConvertPageToGrayscale(page)
			if err != nil {
				return err
			}

			return nil
		},
	}

	// Generate a PdfWriter instance from existing PdfReader.
	pdfWriter, err := pdfReader.ToWriter(opts)
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
