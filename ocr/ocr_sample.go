/**
 * This is a sample Go program that demonstrates how to use the UniPDF library
 * to perform OCR on an image using an HTTP OCR service. The program sends an image
 * to the configured OCR endpoint and displays the extracted text.
 *
 * This example uses https://github.com/unidoc/ocrserver as the OCR service.
 * However, UniPDF's OCR API is designed to support other OCR services that accept
 * image uploads via HTTP and return text or HOCR formatted results.
 *
 * Run as: go run ocr_sample.go input.jpg
 */
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v4/common/license"
	"github.com/unidoc/unipdf/v4/ocr"
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
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run ocr_sample.go input.jpg\n")
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	// Configure OCR service options.
	opts := ocr.OCROptions{
		Url:           "http://localhost:8080/file",
		Method:        "POST",
		FileFieldName: "file",
		Headers: map[string]string{
			"Accept": "application/json",
		},
		TimeoutSeconds: 30,
	}

	// Create OCR client.
	client := ocr.NewHTTPOCRService(opts)

	result, err := client.ExtractText(context.Background(), f, "image.jpg")
	if err != nil {
		fmt.Printf("Error extracting text: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Extracted text: %s\n", string(result))
}
