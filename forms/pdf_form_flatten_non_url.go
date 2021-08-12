/*
 * Partially flatten form data in PDF files by using a callback to filter a non URL annotations.
 *
 * Run as: go run pdf_form_flatten_non_url.go pdf_file
 */

package main

import (
	"fmt"
	"os"

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
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run pdf_form_flatten_non_url.go <input.pdf>\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := fmt.Sprintf("flatten_non_url_%s", inputPath)

	err := nonUrlFlattenPdf(inputPath, outputPath)
	if err != nil {
		fmt.Printf("%s - Error: %v\n", inputPath, err)
	}
}

// nonUrlFlattenPdf flattens non url annotations.
func nonUrlFlattenPdf(inputPath, outputPath string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return err
	}

	// Define annotation filter to exclude URL from flattening process.
	nonUrlFlattenOpts := model.FieldFlattenOpts{
		AnnotFilterFunc: func(pa *model.PdfAnnotation) bool {
			switch pa.GetContext().(type) {
			case *model.PdfAnnotationLink:
				return false
			}

			return true
		},
	}

	fieldAppearance := annotator.FieldAppearance{}
	err = pdfReader.FlattenFieldsWithOpts(fieldAppearance, &nonUrlFlattenOpts)
	if err != nil {
		return err
	}

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
