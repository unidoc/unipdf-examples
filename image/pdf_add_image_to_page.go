/*
 * Insert an image to a PDF file.
 *
 * Adds image to a specific page of a PDF.  xPos and yPos define the upper left corner of the image location, and width
 * is the width of the image in PDF coordinates (height/width ratio is maintained).
 *
 * Example go run pdf_add_image_to_page.go /tmp/input.pdf 1 /tmp/image.jpg 0 0 100 /tmp/output.pdf
 * adds the image to the upper left corner of the page (0,0).  The width is 100 (typical page width 612 with defaults).
 *
 * Syntax: go run pdf_add_image_to_page.go input.pdf <page> image.jpg <xpos> <ypos> <width> output.pdf
 */

package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

func main() {
	if len(os.Args) < 8 {
		fmt.Printf("Usage: go run pdf_add_image_to_page.go input.pdf <page> image.jpg <xpos> <ypos> <width> output.pdf\n")
		os.Exit(1)
	}

	// Use debug logging.
	common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))

	inputPath := os.Args[1]
	pageNumStr := os.Args[2]
	imagePath := os.Args[3]

	xPos, err := strconv.ParseFloat(os.Args[4], 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	yPos, err := strconv.ParseFloat(os.Args[5], 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	iwidth, err := strconv.ParseFloat(os.Args[6], 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	outputPath := os.Args[7]

	fmt.Printf("xPos: %d, yPos: %d\n", xPos, yPos)
	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err = addImageToPdf(inputPath, outputPath, imagePath, pageNum, xPos, yPos, iwidth)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

// Add image to a specific page of a PDF.  xPos and yPos define the upper left corner of the image location, and iwidth
// is the width of the image in PDF document dimensions (height/width ratio is maintained).
func addImageToPdf(inputPath string, outputPath string, imagePath string, pageNum int, xPos float64, yPos float64, iwidth float64) error {
	c := creator.New()

	// Prepare the image.
	img, err := c.NewImageFromFile(imagePath)
	if err != nil {
		return err
	}
	img.ScaleToWidth(iwidth)
	img.SetPos(xPos, yPos)

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

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	// Load the pages.
	for i := 0; i < numPages; i++ {
		page, err := pdfReader.GetPage(i + 1)
		if err != nil {
			return err
		}

		// Add the page.
		err = c.AddPage(page)
		if err != nil {
			return err
		}

		// If the specified page, or -1, apply the image to the page.
		if i+1 == pageNum || pageNum == -1 {
			_ = c.Draw(img)
		}
	}

	err = c.WriteToFile(outputPath)
	return err
}
