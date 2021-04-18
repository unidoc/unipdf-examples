/*
 * Extract text at certain location.
 *
 * Run as: go run pdf_extract_location.go input.pdf <page> <x1> <y1> <x2> <y2>
 */

package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/extractor"
	pdf "github.com/unidoc/unipdf/v3/model"
)

const licenseKey = `
-----BEGIN UNIDOC LICENSE KEY-----
Free trial license keys are available at: https://unidoc.io/
-----END UNIDOC LICENSE KEY-----
`

func init() {
	// Enable debug-level logging.
	// common.SetLogger(common.NewConsoleLogger(common.LogLevelTrace))

	err := license.SetLicenseKey(licenseKey, `Company Name`)
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) < 7 {
		fmt.Printf("Usage: go run pdf_extract_location.go input.pdf page x1 y1 x2 y2\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	page, _ := strconv.Atoi(os.Args[2])
	x1, _ := strconv.ParseFloat(os.Args[3], 64)
	y1, _ := strconv.ParseFloat(os.Args[4], 64)
	x2, _ := strconv.ParseFloat(os.Args[5], 64)
	y2, _ := strconv.ParseFloat(os.Args[6], 64)

	err := extractTextAtLocation(inputPath, page, x1, y1, x2, y2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// extractTextAtLocation prints out contents of PDF file at certain location to stdout.
func extractTextAtLocation(inputPath string, pageNum int, x1, y1, x2, y2 float64) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	totalPage, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	if pageNum > totalPage {
		return errors.New(fmt.Sprintf("Page number %d is not available", pageNum))
	}

	fmt.Printf("--------------------\n")
	fmt.Printf("PDF to text extraction:\n")
	fmt.Printf("--------------------\n")

	page, err := pdfReader.GetPage(pageNum)
	if err != nil {
		return err
	}

	ex, err := extractor.New(page)
	if err != nil {
		return err
	}

	pageText, _, _, err := ex.ExtractPageText()
	if err != nil {
		return err
	}

	pageText.ApplyArea(pdf.PdfRectangle{
		Llx: x1, Lly: y1,
		Urx: x2, Ury: y2,
	})

	text := pageText.Text()

	fmt.Println("------------------------------")
	fmt.Printf("Page %d:\n", pageNum)
	fmt.Printf("Location %.f, %.f, %.f, %.f\n", x1, y1, x2, y2)
	fmt.Printf("\"%s\"\n", text)
	fmt.Println("------------------------------")

	return nil
}
