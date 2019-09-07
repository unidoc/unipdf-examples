/*
 * Annotate/mark up pages of a PDF file.
 * Add a line with arrowhead between two specified points on a page.
 *
 * Run as: go run pdf_annotate_add_line.go input.pdf <page> <x1> <y1> <x2> <y2> output.pdf
 * The x, y coordinates are in the PDF coordinate's system, where 0,0 is in the lower left corner.
 * The line properties are further defined within the example code and can be adjusted.
 */

package main

import (
	"fmt"
	"os"
	"strconv"

	unicommon "github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/contentstream/draw"
	pdf "github.com/unidoc/unipdf/v3/model"
)

func main() {
	// Debug log mode.
	//unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	if len(os.Args) < 8 {
		fmt.Printf("go run pdf_annotate_add_line.go input.pdf <page> <x1> <y1> <x2> <y2> output.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]

	pageNum, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	x1, err := strconv.ParseFloat(os.Args[3], 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	y1, err := strconv.ParseFloat(os.Args[4], 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	x2, err := strconv.ParseFloat(os.Args[5], 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	y2, err := strconv.ParseFloat(os.Args[6], 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	outputPath := os.Args[7]

	// Draw a 3pt line with arrowhead.
	lineDef := annotator.LineAnnotationDef{}
	lineDef.X1, lineDef.Y1 = x1, y1
	lineDef.X2, lineDef.Y2 = x2, y2
	lineDef.LineColor = pdf.NewPdfColorDeviceRGB(1.0, 0.0, 0.0) // Red.
	lineDef.Opacity = 0.50
	lineDef.LineEndingStyle1 = draw.LineEndingStyleNone
	lineDef.LineEndingStyle2 = draw.LineEndingStyleArrow
	lineDef.LineWidth = 3.0

	err = annotatePdfAddLineAnnotation(inputPath, pageNum, outputPath, lineDef)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

// Annotate pdf file.
func annotatePdfAddLineAnnotation(inputPath string, pageNum int64, outputPath string, lineDef annotator.LineAnnotationDef) error {
	unicommon.Log.Debug("Input PDF: %v", inputPath)

	pdfWriter := pdf.NewPdfWriter()

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

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	for i := 0; i < numPages; i++ {
		// Read the page.
		page, err := pdfReader.GetPage(int(i + 1))
		if err != nil {
			return err
		}

		if int(pageNum) == (i + 1) {
			lineAnnotation, err := annotator.CreateLineAnnotation(lineDef)
			if err != nil {
				return err
			}

			// Add to the page annotations.
			page.AddAnnotation(lineAnnotation)
		}

		err = pdfWriter.AddPage(page)
		if err != nil {
			unicommon.Log.Error("Failed to add page: %s", err)
			return err
		}
	}

	fWrite, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer fWrite.Close()

	err = pdfWriter.Write(fWrite)
	if err != nil {
		return err
	}

	return nil
}
