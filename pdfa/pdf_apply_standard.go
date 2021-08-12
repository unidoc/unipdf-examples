/*
 * PDF optimization (compression) example.
 *
 * Run as: go run pdf_apply_standard.go <input.pdf> <output.pdf>
 */

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/model/pdfa"
)

const licenseKey = `
-----BEGIN UNIDOC LICENSE KEY-----
Free trial license keys are available at: https://unidoc.io/
-----END UNIDOC LICENSE KEY-----
`

const usage = "Usage: %s INPUT_PDF_PATH OUTPUT_PDF_PATH\n"

func init() {
	err := license.SetLicenseKey(licenseKey, `Company Name`)
	if err != nil {
		panic(err)
	}
}

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Printf(usage, os.Args[0])
		return
	}
	inputPath := args[1]
	outputPath := args[2]

	// Initialize starting time.
	start := time.Now()

	// Create reader.
	reader, file, err := model.NewPdfReaderFromFile(inputPath, nil)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
	defer file.Close()

	// Generate a PDFWriter from PDFReader.
	pdfWriter, err := reader.ToWriter(nil)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Apply standard PDF/A-1B.
	pdfWriter.ApplyStandard(pdfa.NewProfile1B(nil))

	// Create output file.
	err = pdfWriter.WriteToFile(outputPath)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	duration := float64(time.Since(start)) / float64(time.Millisecond)
	fmt.Printf("Processing time: %.2f ms\n", duration)
}

