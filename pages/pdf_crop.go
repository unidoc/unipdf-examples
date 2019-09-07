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

	pdf "github.com/unidoc/unipdf/v3/model"
)

func main() {
	// When debugging: log to console.
	//unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

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
	pdfWriter := pdf.NewPdfWriter()

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

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return err
		}

		bbox, err := page.GetMediaBox()
		if err != nil {
			return err
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

		err = pdfWriter.AddPage(page)
		if err != nil {
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
