/*
 * Crop pages in a PDF file. Crops the view to a certain percentage  of the original.
 * The percentage specifies the trim-off percentage, both width- and heightwise.
 *
 * Run as: go run pdf_crop.go input.pdf <percentage> output.pdf
 */

package main

import (
	"fmt"
	"os"
	"strconv"

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
		fmt.Printf("Usage: go run pdf_crop.go input.pdf <percentage> output.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	percentageStr := os.Args[2]
	outputPath := os.Args[3]

	percentage, err := strconv.ParseInt(percentageStr, 10, 32)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	if percentage < 0 || percentage > 100 {
		fmt.Printf("Percentage should be in the range 0 - 100 (%)\n")
		os.Exit(1)
	}

	err = cropPdf(inputPath, outputPath, percentage)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

// Crop all pages by a given percentage.
func cropPdf(inputPath string, outputPath string, percentage int64) error {
	readerOpts := model.NewReaderOpts()
	readerOpts.LazyLoad = false

	pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, readerOpts)
	if err != nil {
		return err
	}
	defer f.Close()

	// Process each page using the following callback
	// when generating PdfWriter from PdfReader.
	opts := &model.ReaderToWriterOpts{
		PageProcessCallback: func(pageNum int, page *model.PdfPage) error {
			bbox, err := page.GetMediaBox()
			if err != nil {
				return err
			}

			// Zoom in on the page middle, with a scaled width and height.
			width := (*bbox).Urx - (*bbox).Llx
			height := (*bbox).Ury - (*bbox).Lly
			newWidth := width * float64(percentage) / 100.0
			newHeight := height * float64(percentage) / 100.0
			(*bbox).Llx += newWidth / 2
			(*bbox).Lly += newHeight / 2
			(*bbox).Urx -= newWidth / 2
			(*bbox).Ury -= newHeight / 2

			page.MediaBox = bbox

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
