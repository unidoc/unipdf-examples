/*
 * Annotate/mark up pages of a PDF file.
 * Add a rectangle annotation to a specified location on a page.
 *
 * Run as: go run pdf_annotate_add_rectangle.go input.pdf output.pdf <x> <y> <width> <height>
 * The x, y, width and height coordinates are in the PDF coordinate's system, where 0,0 is in the lower left corner.
 */

package main

import (
	"fmt"
	"os"

	"strconv"

	unicommon "github.com/unidoc/unidoc/common"
	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func init() {
	// Debug log mode.
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))
}

func main() {
	if len(os.Args) < 7 {
		fmt.Printf("go run pdf_annotate_add_rectangle.go input.pdf output.pdf <x> <y> <width> <height>\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	x, err := strconv.ParseInt(os.Args[3], 10, 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	y, err := strconv.ParseInt(os.Args[4], 10, 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	width, err := strconv.ParseInt(os.Args[5], 10, 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	height, err := strconv.ParseInt(os.Args[6], 10, 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err = annotatePdfAddRectAnnotation(inputPath, outputPath, x, y, width, height)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

// Annotate pdf file.
func annotatePdfAddRectAnnotation(inputPath string, outputPath string, x, y, width, height int64) error {
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
		pageNum := i + 1

		// Read the page.
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return err
		}

		// Rectangle annotation.
		rectAnnotation := pdf.NewPdfAnnotationSquare()
		rectAnnotation.C = pdfcore.MakeArrayFromFloats([]float64{1.0, 0.0, 0.0}) // Red border.
		rectAnnotation.IC = pdfcore.MakeArrayFromIntegers([]int{})               // No fill.
		bs := pdf.NewBorderStyle()
		bs.SetBorderWidth(3) // Width: 3 points.
		rectAnnotation.BS = bs.ToPdfObject()

		// The rect specifies the location of the markup.
		rectAnnotation.Rect = pdfcore.MakeArrayFromIntegers([]int{int(x), int(y), int(x + width), int(y + height)})

		// Add to the page annotations.
		page.Annotations = append(page.Annotations, rectAnnotation.PdfAnnotation)

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
