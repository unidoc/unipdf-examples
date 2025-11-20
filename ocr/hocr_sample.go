/**
 * This is a sample Go program that demonstrates how to use the UniPDF library
 * to perform OCR on an image using an HTTP OCR service that returns HOCR formatted
 * output. The program parses the HOCR response and extracts word-level information
 * including bounding boxes and confidence scores.
 *
 * Run as: go run hocr_sample.go input.jpg
 */
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/stefanhengl/gohocr"
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
		FormFields: map[string]string{
			"format": "hocr",
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

	// Parse JSON response to extract the "result" field.
	var jsonObj map[string]interface{}
	if err := json.Unmarshal(result, &jsonObj); err != nil {
		fmt.Printf("Error parsing JSON response: %s", err)
		return
	}

	content := jsonObj["result"].(string)
	fmt.Printf("Extracted text: %s", content)

	content, err = strconv.Unquote(content)
	if err != nil {
		fmt.Printf("Error unquoting content: %s", err)
		return
	}

	contentBytes := []byte(content)

	data, err := gohocr.Parse(contentBytes)
	if err != nil {
		fmt.Printf("Error parsing HOCR data: %s", err)
		return
	}

	for _, v := range data.Words {
		fmt.Printf("Word: %s, Title: %f\n", v.Content, v.Title)
	}
}
