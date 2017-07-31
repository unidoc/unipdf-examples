/*
 * Add watermark image to each page of a PDF file.
 *
 * Run as: go run pdf_watermark_image.go input.pdf watermark.jpg output.pdf
 */

package main

import (
	"fmt"
	"os"

	unicommon "github.com/unidoc/unidoc/common"
	"github.com/unidoc/unidoc/pdf/creator"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func main() {
	// Enable console-level debug-mode logging when debugging:
	//unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	if len(os.Args) < 4 {
		fmt.Printf("go run pdf_watermark_image.go input.pdf watermark.jpg output.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	watermarkPath := os.Args[2]
	outputPath := os.Args[3]

	err := addWatermarkImage(inputPath, outputPath, watermarkPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

// Watermark pdf file based on an image.
func addWatermarkImage(inputPath string, outputPath string, watermarkPath string) error {
	unicommon.Log.Debug("Input PDF: %v", inputPath)
	unicommon.Log.Debug("Watermark image: %s", watermarkPath)

	c := creator.New()

	watermarkImg, err := creator.NewImageFromFile(watermarkPath)
	if err != nil {
		return err
	}

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

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return err
	}

	// If encrypted, try to decrypt with an empty password.
	if isEncrypted {
		// Fails, try fallback with empty password.
		_, err = pdfReader.Decrypt([]byte(""))
		if err != nil {
			return err
		}
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

		// Add to creator.
		c.AddPage(page)

		watermarkImg.ScaleToWidth(c.Context().PageWidth)
		watermarkImg.SetPos(0, (c.Context().PageHeight-watermarkImg.Height())/2)
		watermarkImg.SetOpacity(0.5)
		_ = c.Draw(watermarkImg)
	}

	err = c.WriteToFile(outputPath)
	return err
}
