/*
 * Add watermark to each page of a PDF file.
 *
 * Run as: go run pdf_watermark_image.go input.pdf output.pdf watermark.jpg
 */

package main

import (
	"fmt"
	"os"

	unicommon "github.com/unidoc/unidoc/common"
	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func init() {
	// Debug log level.
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))
}

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("go run pdf_watermark_image.go input.pdf output.pdf watermark.jpg\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]
	watermarkPath := os.Args[3]

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

	pdfWriter := pdf.NewPdfWriter()

	// Open the watermark image file.
	reader, err := os.Open(watermarkPath)
	if err != nil {
		unicommon.Log.Error("Error opening file: %s", err)
		return err
	}
	defer reader.Close()

	watermarkImg, err := pdf.ImageHandling.Read(reader)
	if err != nil {
		unicommon.Log.Error("Error loading image: %s", err)
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

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	imgName := pdfcore.PdfObjectName("Imw0")
	ximg, err := pdf.NewXObjectImageFromImage(imgName, watermarkImg, nil)
	if err != nil {
		unicommon.Log.Error("Failed to create xobject image: %s", err)
		return err
	}
	err = ximg.Compress()
	if err != nil {
		unicommon.Log.Error("Failed to compress: %v", err)
		return err
	}

	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		// Read the page.
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return err
		}

		wmOpt := pdf.WatermarkImageOptions{}
		wmOpt.Alpha = 0.5
		wmOpt.FitToWidth = true
		wmOpt.PreserveAspectRatio = true

		err = page.AddWatermarkImage(ximg, wmOpt)
		if err != nil {
			return err
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
