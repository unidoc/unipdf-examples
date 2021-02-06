/*
 * Duplicate a PDF file.
 *
 * Run as: go run pdf_duplicate_doc.go input.pdf output.pdf
 */

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/common/license"
	pdf "github.com/unidoc/unipdf/v3/model"
)

const licenseKey = `
-----BEGIN UNIDOC LICENSE KEY-----
Free trial license keys are available at: https://unidoc.io/
-----END UNIDOC LICENSE KEY-----
`

func init() {
	// Enable debug-level logging.
	// common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))

	err := license.SetLicenseKey(licenseKey, `Company Name`)
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("go run pdf_duplicate_doc.go input.pdf output.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	err := duplicateDoc(inputPath, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

func duplicateDoc(inputPath string, outputPath string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return nil
	}

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return nil
	}

	pdfWriter, err := pdfReader.ToWriter(nil)
	if err != nil {
		return nil
	}

	err = pdfWriter.WriteToFile(outputPath)
	if err != nil {
		return err
	}

	return nil
}
