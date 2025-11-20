/**
 * This is a sample Go program that demonstrates how to use the UniPDF library
 * to perform OCR on an image using an HTTP OCR service. The program sends an image
 * to the configured OCR endpoint and displays the extracted text.
 *
 * Run as: go run ocr_sample.go input.jpg
 */
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v4/ocr"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run ocr_sample.go input.jpg\n")
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error opening file: %s", err)
		return
	}
	defer func() { _ = f.Close() }()

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
		fmt.Printf("Error extracting text: %s", err)
		return
	}

	fmt.Printf("Extracted text: %s", result)
}
