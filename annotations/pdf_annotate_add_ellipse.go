/*
 * Annotate/mark up pages of a PDF file.
 * Add a circle/ellipse annotation to a specified location on a page.
 *
 * Run as: go run pdf_annotate_add_ellipse.go input.pdf <page> <x> <y> <width> <height> output.pdf
 * The x, y, width and height coordinates are in the PDF coordinate's system, where 0,0 is in the lower left corner.
 */

package main

import (
	"fmt"
	"os"

	"strconv"

	"github.com/unidoc/unipdf/v3/annotator"
	unicommon "github.com/unidoc/unipdf/v3/common"
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
	// unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	err := license.SetLicenseKey(licenseKey, `Company Name`)
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) < 8 {
		fmt.Printf("go run pdf_annotate_add_ellipse.go input.pdf <page> <x> <y> <xRad> <yRad> output.pdf\n")
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

	err = annotatePdfAddEllipseAnnotation(inputPath, pageNum, outputPath, x, y, width, height)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

// Annotate pdf file.
func annotatePdfAddEllipseAnnotation(inputPath string, targetPageNum int64, outputPath string, x, y, width, height float64) error {
	unicommon.Log.Debug("Input PDF: %v", inputPath)

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

	// Process each page using the following callback
	// when generating PdfWriter.
	opt := &pdf.ReaderToWriterOpts{
		PageCallback: func(pageNum int, page *pdf.PdfPage) {
			// Add only to the specific page.
			if int(targetPageNum) == pageNum {
				// Define a semi-transparent yellow rectangle with black borders at the specified location.

				circDef := annotator.CircleAnnotationDef{}
				circDef.X = x
				circDef.Y = y
				circDef.Width = width
				circDef.Height = height
				circDef.Opacity = 0.5 // Semi transparent.
				circDef.FillEnabled = true
				circDef.FillColor = pdf.NewPdfColorDeviceRGB(1, 1, 0) // Yellow fill.
				circDef.BorderEnabled = true
				circDef.BorderWidth = 15
				circDef.BorderColor = pdf.NewPdfColorDeviceRGB(0, 0, 0) // Black border.

				circAnnotation, err := annotator.CreateCircleAnnotation(circDef)
				if err != nil {
					fmt.Println(err)

					return
				}

				// Add to the page annotations.
				page.AddAnnotation(circAnnotation)
			}
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
