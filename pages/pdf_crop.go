/*
 * Crop pages in a PDF file. Crops the view to a certain percentage  of the original.
 * The percentage specifies the trim-off percentage, both width- and heightwise.
 *
 * Run as: go run pdf_crop.go input.pdf <percentage> output.pdf
 */

package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

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
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return err
	}

	// Try decrypting both with given password and an empty one if that fails.
	if isEncrypted {
		auth, err := pdfReader.Decrypt([]byte(""))
		if err != nil {
			return err
		}
		if !auth {
			return errors.New("Unable to decrypt pdf with empty pass")
		}
	}

	opts := &pdf.ReaderToWriterOpts{
		PageCallback: func(pageNum int, page *pdf.PdfPage) {
			bbox, err := page.GetMediaBox()
			if err != nil {
				fmt.Println(err)
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
		},
	}

	pdfWriter, err := pdfReader.ToWriter(opts)
	if err != nil {
		return err
	}

	err = pdfWriter.WriteToFile(outputPath)
	if err != nil {
		return err
	}

	return nil
}
