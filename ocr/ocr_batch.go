/**
 * This is a sample Go program that demonstrates how to use the UniPDF library
 * to perform batch OCR processing on multiple images using an HTTP OCR service.
 * The program processes multiple image files concurrently and displays the extracted
 * text results along with a summary of successful and failed operations.
 *
 * This example uses https://github.com/unidoc/ocrserver as the OCR service.
 * However, UniPDF's OCR API is designed to support other OCR services that accept
 * image uploads via HTTP and return text or HOCR formatted results.
 *
 * Run as: go run ocr_batch.go image1.jpg image2.png [image3.jpg ...]
 */
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

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
		fmt.Printf("Usage: go run ocr_batch.go image1.jpg image2.png [image3.jpg ...]\n")
		os.Exit(1)
	}

	// Get list of image files from command line arguments
	filePaths := os.Args[1:]

	// Validate that all files exist
	for _, filePath := range filePaths {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Printf("Error: File does not exist: %s\n", filePath)
			os.Exit(1)
		}
	}

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
	client := ocr.NewOCRHTTPClient(opts)

	fmt.Printf("Processing %d files...\n", len(filePaths))

	// Batch process files.
	results, errors := client.BatchProcessFiles(context.Background(), filePaths)

	// Display results
	for i, filePath := range filePaths {
		filename := filepath.Base(filePath)
		fmt.Printf("\n--- Results for %s ---\n", filename)

		if errors[i] != nil {
			fmt.Printf("Error processing %s: %s\n", filename, errors[i])
			continue
		}

		fmt.Printf("Extracted text from %s:\n%s\n", filename, string(results[i]))
	}

	// Summary
	successCount := 0
	errorCount := 0
	for _, err := range errors {
		if err != nil {
			errorCount++
		} else {
			successCount++
		}
	}

	fmt.Printf("\n--- Summary ---\n")
	fmt.Printf("Successfully processed: %d files\n", successCount)
	fmt.Printf("Failed to process: %d files\n", errorCount)
	fmt.Printf("Total files: %d\n", len(filePaths))
}
