/*
 * Creates outlines for a PDF file.
 *
 * Run as: go run pdf_create_outlines.go input.pdf output.pdf
 */

package main

import (
	"fmt"
	"os"

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
	if len(os.Args) < 3 {
		fmt.Printf("Usage: go run pdf_create_outlines.go input.pdf output.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outPath := os.Args[2]

	fmt.Printf("Input file: %s\n", inputPath)
	fmt.Printf("Output file: %s\n", outPath)

	pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	// Check number of PDF pages.
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if numPages < 3 {
		fmt.Printf("Too short PDF file. At least 3 pages are needed for this outline example.")
		os.Exit(1)
	}

	// Don't copy document outlines.
	opt := &model.ReaderToWriterOpts{
		SkipOutlines: true,
	}

	// Generate a PdfWriter instance from existing PdfReader.
	pdfWriter, err := pdfReader.ToWriter(opt)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Add the new document outlines.
	outline := model.NewOutline()
	for i := 0; i < 3; i++ {
		page, _ := pdfReader.GetPage(i + 1)
		_, y, _ := page.Size()
		// Create outline for the top left corner of each page. Note that PDF y coordinate goes from bottom to top.
		outline.Add(model.NewOutlineItem(fmt.Sprintf("page%d", i+1), model.NewOutlineDest(int64(i), 0, y)))
	}

	pdfWriter.AddOutlineTree(outline.ToOutlineTree())

	// Write output file.
	err = pdfWriter.WriteToFile(outPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
