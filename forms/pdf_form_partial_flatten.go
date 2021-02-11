/*
 * Partially flatten form data in PDF file by using a callback function to filter Fields.
 *
 * Run as: go run pdf_form_partial_flatten.go
 */

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/model"
)

const licenseKey = `
-----BEGIN UNIDOC LICENSE KEY-----
Free trial license keys are available at: https://unidoc.io/
-----END UNIDOC LICENSE KEY-----
`

func init() {
	// Enable debug-level logging.
	// unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	err := license.SetLicenseKey(licenseKey, `Company Name`)
	if err != nil {
		panic(err)
	}
}

func main() {
	inputFile := "sample_form.pdf"
	outputFile := "partial_flattened_sample_form.pdf"
	fieldToFlatten := []string{
		"email4",
		"address5[city]",
	}

	err := partialFlattenPdf(inputFile, outputFile, fieldToFlatten)
	if err != nil {
		fmt.Printf("%s - Error: %v\n", inputFile, err)
		os.Exit(1)
	}
}

// partialFlattenPdf partially flattens annotations and forms moving
// the appearance stream to the page contents so cannot be modified.
func partialFlattenPdf(inputPath, outputPath string, fieldToFlatten []string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return err
	}

	fieldFlattenOpts := model.FieldFlattenOpts{
		FilterFunc: func(pf *model.PdfField) bool {
			for _, fName := range fieldToFlatten {
				if pf.T.String() == fName {
					return true
				}
			}

			return false
		},
	}

	fieldAppearance := annotator.FieldAppearance{}
	err = pdfReader.FlattenFieldsWithOpts(fieldAppearance, &fieldFlattenOpts)
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
	return err
}
