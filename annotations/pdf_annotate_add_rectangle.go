/*
 * Annotate/mark up pages of a PDF file.
 * Add a rectangle annotation to a specified location on a page.
 *
 * Run as: go run pdf_annotate_add_rectangle.go input.pdf <page> <x> <y> <width> <height> output.pdf
 * The x, y, width and height coordinates are in the PDF coordinate's system, where 0,0 is in the lower left corner.
 */

package main

import (
	"fmt"
	"os"
	"strconv"

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
	if len(os.Args) < 8 {
		fmt.Printf("go run pdf_annotate_add_rectangle.go input.pdf <page> <x> <y> <width> <height> output.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]

	pageNum, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	x, err := strconv.ParseFloat(os.Args[3], 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	y, err := strconv.ParseFloat(os.Args[4], 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	width, err := strconv.ParseFloat(os.Args[5], 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	height, err := strconv.ParseFloat(os.Args[6], 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	outputPath := os.Args[7]

	err = annotatePdfAddRectAnnotation(inputPath, pageNum, outputPath, x, y, width, height)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

// Annotate pdf file.
func annotatePdfAddRectAnnotation(inputPath string, targetPageNum int64, outputPath string, x, y, width, height float64) error {
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
			// Add only to the specific page.
			if int(targetPageNum) == pageNum {
				// Define a semi-transparent yellow rectangle with black borders at the specified location.
				rectDef := annotator.RectangleAnnotationDef{}
				rectDef.X = x
				rectDef.Y = y
				rectDef.Width = width
				rectDef.Height = height
				rectDef.Opacity = 0.5 // Semi transparent.
				rectDef.FillEnabled = false
				rectDef.FillColor = model.NewPdfColorDeviceRGB(1, 1, 0) // Yellow fill.
				rectDef.BorderEnabled = true
				rectDef.BorderWidth = 30
				rectDef.BorderColor = model.NewPdfColorDeviceRGB(0, 0, 0) // Black border.

				rectAnnotation, err := annotator.CreateRectangleAnnotation(rectDef)
				if err != nil {
					return err
				}

				// Add to the page annotations.
				page.AddAnnotation(rectAnnotation)
			}

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
