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

	// Don't copy document info into the new PDF.
	opt := &model.ReaderToWriterOpts{
		SkipInfo: true,
	}

	// Generate a PdfWriter instance from existing PdfReader.
	defaultPdfWriter, err := pdfReader.ToWriter(opt)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	customPdfWriter := defaultPdfWriter

	// Write new PDF with default author name.
	err = defaultPdfWriter.WriteToFile("gen_pdf_default_author.pdf")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Write new PDF with custom information dictionary.
	pdfInfo := &model.PdfInfo{}
	pdfInfo.Author = core.MakeString("UniPDF Tester 2")
	pdfInfo.Subject = core.MakeString("PDF Example with custom information dictionary")
	pdfInfo.AddCustomInfo("custom_info", "This is an optional custom info")

	customPdfWriter.SetDocInfo(pdfInfo)

	err = customPdfWriter.WriteToFile("gen_pdf_custom_info.pdf")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
