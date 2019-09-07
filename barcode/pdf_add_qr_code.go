/*
 * Create a barcode or a QR code and insert on a specific location in a PDf file.
 * If unsure about position, try getting the dimensions of a PDF with unipdf-examples/pages/pdf_page_info.go first,
 * or just start with 0,0 and increase to move right, down.
 *
 * Run as: go run pdf_add_qr_code.go input.pdf <page> <qrtext> <xpos> <ypos> <width> output.pdf
 * - The x and y positions are relative to the upper left corner of the page.
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
	"github.com/boombuler/barcode/qr"

	unicommon "github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

func main() {
	if len(os.Args) < 8 {
		fmt.Printf("Usage: go run pdf_add_qr_code.go input.pdf <page> <qrtext> <xpos> <ypos> <width> output.pdf\n")
		os.Exit(1)
	}

	// Use debug logging.
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	inputPath := os.Args[1]
	pageNumStr := os.Args[2]
	textStr := os.Args[3]
	outputPath := os.Args[7]

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

	err = addQrCodeToPdf(inputPath, outputPath, textStr, pageNum, xPos, yPos, width)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)

}

// Prepare the QR code. The oversampling ratio specifies how many pixels/point to use.  The default resolution of
// PDFs is 72PPI (points per inch). A higher PPI allows higher resolution QR code generation which is particularly
// important if the document is scaled (zoom in).
func makeQrCode(contentStr string, width float64, oversampling int) (goimage.Image, error) {
	qrCode, err := qr.Encode(contentStr, qr.M, qr.Auto)
	if err != nil {
		return nil, err
	}

	// Prepare the qr code image.
	pixelWidth := oversampling * int(math.Ceil(width))
	qrCode, err = barcode.Scale(qrCode, pixelWidth, pixelWidth)
	if err != nil {
		return nil, err
	}

	return qrCode, err
}

// Add image to a specific page of a PDF.  xPos and yPos define the lower left corner of the image location, and iwidth
// is the width of the image in PDF coordinates (height/width ratio is maintained).
func addQrCodeToPdf(inputPath string, outputPath string, qrContentStr string, pageNum int, xPos float64, yPos float64, width float64) error {
	qrCode, err := makeQrCode(qrContentStr, width, 5)
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
			// Apply the QR code to the specified page or all pages if -1.
			img, err := c.NewImageFromGoImage(qrCode)
			if err != nil {
				return err
			}
			img.SetWidth(width)
			img.SetHeight(width)
			img.SetPos(xPos, yPos)
			err = c.Draw(img)
			if err != nil {
				return err
			}
		}
	}

	err = c.WriteToFile(outputPath)
	if err != nil {
		return err
	}

	return nil
}
