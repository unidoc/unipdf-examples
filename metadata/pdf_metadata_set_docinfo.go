/*
 * Generate multiple copy of template pdf file which contains different
 * Document Information Dictionary value.
 *
 * Run as: go run pdf_metadata_set_docinfo.go template.pdf
 */

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
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
	if len(os.Args) < 2 {
		fmt.Printf("Set Document Information Dictionary information in PDF file\n")
		fmt.Printf("Usage: go run pdf_metadata_set_docinfo.go template.pdf\n")
		os.Exit(1)
	}

	author := "UniPDF Tester"
	model.SetPdfAuthor(author)

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	defaultPdfWriter := model.NewPdfWriter()

	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		err = defaultPdfWriter.AddPage(page)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}

	customPdfWriter := defaultPdfWriter

	// Write new PDF with default author name
	defaultFWrite, err := os.Create("gen_pdf_default_author.pdf")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer defaultFWrite.Close()

	err = defaultPdfWriter.Write(defaultFWrite)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Write new PDF with custom information dictionary
	pdfInfo := &model.PdfInfo{}
	pdfInfo.Author = core.MakeString("UniPDF Tester 2")
	pdfInfo.Subject = core.MakeString("PDF Example with custom information dictionary")

	customFWrite, err := os.Create("gen_pdf_custom_info.pdf")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	defer customFWrite.Close()

	customPdfWriter.SetDocInfo(pdfInfo)

	err = customPdfWriter.Write(customFWrite)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
