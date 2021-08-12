/*
 * Extract images from a PDF file. Passes through each page, goes through the content stream and finds instances of both
 * XObject Images and inline images. Also handles images referred within XObject Form content streams.
 * The output files are saved as a zip archive.
 *
 * Run as: go run pdf_extract_images.go input.pdf output.zip
 */

package main

import (
	"archive/zip"
	"fmt"
	"image/jpeg"
	"os"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/extractor"
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
	if len(os.Args) < 3 {
		fmt.Printf("Syntax: go run pdf_extract_images.go input.pdf output.zip\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	fmt.Printf("Input file: %s\n", inputPath)
	err := extractImagesToArchive(inputPath, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// Extracts images and properties of a PDF specified by inputPath.
// The output images are stored into a zip archive whose path is given by outputPath.
func extractImagesToArchive(inputPath, outputPath string) error {
	pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, nil)
	if err != nil {
		return err
	}
	defer f.Close()

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}
	fmt.Printf("PDF Num Pages: %d\n", numPages)

	// Prepare output archive.
	zipf, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer zipf.Close()
	zipw := zip.NewWriter(zipf)

	totalImages := 0
	for i := 0; i < numPages; i++ {
		fmt.Printf("-----\nPage %d:\n", i+1)

		page, err := pdfReader.GetPage(i + 1)
		if err != nil {
			return err
		}

		pextract, err := extractor.New(page)
		if err != nil {
			return err
		}

		pimages, err := pextract.ExtractPageImages(nil)
		if err != nil {
			return err
		}

		fmt.Printf("%d Images\n", len(pimages.Images))
		for idx, img := range pimages.Images {
			fmt.Printf("Image %d - X: %.2f Y: %.2f, Width: %.2f, Height: %.2f\n",
				totalImages+idx+1, img.X, img.Y, img.Width, img.Height)
			fname := fmt.Sprintf("p%d_%d.jpg", i+1, idx)

			gimg, err := img.Image.ToGoImage()
			if err != nil {
				return err
			}

			imgf, err := zipw.Create(fname)
			if err != nil {
				return err
			}
			opt := jpeg.Options{Quality: 100}
			err = jpeg.Encode(imgf, gimg, &opt)
			if err != nil {
				return err
			}
		}
		totalImages += len(pimages.Images)
	}
	fmt.Printf("Total: %d images\n", totalImages)

	// Make sure to check the error on Close.
	err = zipw.Close()
	if err != nil {
		return err
	}

	return nil
}
