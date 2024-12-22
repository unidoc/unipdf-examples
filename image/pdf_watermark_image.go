/*
 * Add watermark image to each page of a PDF file.
 *
 * Run as: go run pdf_watermark_image.go input.pdf watermark.jpg output.pdf
 */

package main

import (
	"fmt"
	"image"
	"os"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/model"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}
}

func main() {
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
	common.Log.Debug("Input PDF: %v", inputPath)
	common.Log.Debug("Watermark image: %s", watermarkPath)

	// Create the watermark image from file.
	wImgFile, err := os.Open(watermarkPath)
	if err != nil {
		return err
	}
	defer wImgFile.Close()

	watermarkImg, _, err := image.Decode(wImgFile)
	if err != nil {
		return err
	}

	image, err := model.DefaultImageHandler{}.NewImageFromGoImage(watermarkImg)
	if err != nil {
		return err
	}

	xImage, err := model.NewXObjectImageFromImage(image, nil, nil)

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

	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		// Read the page.
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return err
		}

		// Add watermark with options.
		page.AddWatermarkImage(xImage, model.WatermarkImageOptions{Alpha: 0.5, FitToWidth: true})
	}

	// Generate a PDFWriter from PDFReader.
	pdfWriter, err := pdfReader.ToWriter(nil)
	if err != nil {
		return err
	}

	err = pdfWriter.WriteToFile(outputPath)
	if err != nil {
		return err
	}
	return nil
}
