/*
 * Create a barcode and insert on a specific location in a PDf file.
 * If unsure about position, try getting the dimensions of a PDF with unipdf-examples/pages/pdf_page_info.go first,
 * or just start with 0,0 and increase to move right, down.
 *
 * This example demonstrates an EAN barcode.  It is worth noting that the barcode package supports multiple other
 * types, see: https://github.com/boombuler/barcode.
 *
 * Run as: go run pdf_add_barcode.go input.pdf <page> <code> <xpos> <ypos> <width> output.pdf
 * - The x and y positions are relative to the upper left corner of the page.
 * - If page number is set to -1, the barcode is applied to all pages.
 *
 * As an example of adding EAN-13 barcode with a width of 100:
 * Example: go run pdf_add_barcode.go myfile1.pdf -1 "123456789012" 100 100 100 myfile1_out.pdf
 */
/*
 * NOTE: This example depends on github.com/boombuler/barcode, MIT licensed.
 */

package main

import (
	"fmt"
	goimage "image"
	"math"
	"os"
	"strconv"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/ean"

	unicommon "github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

func main() {
	if len(os.Args) < 8 {
		fmt.Printf("Usage: go run pdf_add_barcode.go input.pdf <page> <code> <xpos> <ypos> <width> output.pdf\n")
		os.Exit(1)
	}

	// Use debug logging.
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	inputPath := os.Args[1]
	pageNumStr := os.Args[2]
	codeStr := os.Args[3]
	outputPath := os.Args[7]

	allowedLengths := map[int]bool{7: true, 8: true, 12: true, 13: true}
	if _, ok := allowedLengths[len(codeStr)]; !ok {
		fmt.Printf("Code must be either 7 characters (EAN-8) or 12 characters long (EAN-13) if provided without checksum\n")
		fmt.Printf("Code must be either 8 characters (EAN-8) or 13 characters long (EAN-13) if provided with checksum\n")
		os.Exit(1)
	}

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
	width, err := strconv.ParseFloat(os.Args[6], 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err = addBarcodeToPdf(inputPath, outputPath, codeStr, pageNum, xPos, yPos, width)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)

}

// Prepare the barcode. The oversampling ratio specifies how many pixels/point to use.  The default resolution of
// PDFs is 72PPI (points per inch). A higher PPI allows higher resolution QR code generation which is particularly
// important if the document is scaled (zoom in).
func makeBarcode(codeStr string, width float64, oversampling int) (goimage.Image, error) {
	bcode, err := ean.Encode(codeStr)
	if err != nil {
		return nil, err
	}

	// Prepare the code image.
	pixelWidth := oversampling * int(math.Ceil(width))
	bcodeImg, err := barcode.Scale(bcode, pixelWidth, pixelWidth)
	if err != nil {
		return nil, err
	}

	return bcodeImg, err
}

// Add barcode to a specific page of a PDF.  xPos and yPos define the lower left corner of the image location, and iwidth
// is the width of the image in PDF coordinates (height/width ratio is maintained).
func addBarcodeToPdf(inputPath string, outputPath string, codeStr string, pageNum int, xPos float64, yPos float64, width float64) error {
	bcodeImg, err := makeBarcode(codeStr, width, 5)
	if err != nil {
		return err
	}

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

	// Make a new PDF creator.
	c := creator.New()

	// Load the pages and add to creator.  Apply the QR code to the specified page.
	for i := 0; i < numPages; i++ {
		page, err := pdfReader.GetPage(i + 1)
		if err != nil {
			return err
		}

		err = c.AddPage(page)
		if err != nil {
			return err
		}

		if i+1 == pageNum || pageNum == -1 {
			// Apply the barcode to the specified page or all pages if -1.
			img, err := c.NewImageFromGoImage(bcodeImg)
			if err != nil {
				return err
			}
			img.ScaleToWidth(width)
			img.SetPos(xPos, yPos)
			_ = c.Draw(img)

			// Add the code below.
			p := c.NewParagraph(codeStr)
			p.SetWidth(width)
			p.SetTextAlignment(creator.TextAlignmentCenter)
			p.SetPos(xPos, yPos+img.Height())
			_ = c.Draw(p)
		}
	}

	err = c.WriteToFile(outputPath)
	if err != nil {
		return err
	}

	return nil
}
