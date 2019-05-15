/*
 * PDF to text: Extract all text for each page of a pdf file.
 *
 * Run as: go run pdf_extract_text.go input.pdf
 */

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unidoc/common"
	"github.com/unidoc/unidoc/pdf/extractor"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run pdf_extract_text.go input.pdf\n")
		os.Exit(1)
	}

	// Make sure to enter a valid license key.
	// Otherwise text is truncated and a watermark added to the text.
	// License keys are available via: https://unidoc.io
	/*
					license.SetLicenseKey(`
		-----BEGIN UNIDOC LICENSE KEY-----
		eyJsaWNlbnNlX2lkIjoiMjA0YWIxMjgtZGY5Yy00ZWE3LTdlM2UtNzJiYTk4OWFhNGZmIiwiY3VzdG9tZXJfaWQiOiJlMDRiODNjZC0zOTYzLTQxNDktNjljOC03MDU0MTM0OWUyMWMiLCJjdXN0b21lcl9uYW1lIjoiUGFwZXJjdXQiLCJ0eXBlIjoiY29tbWVyY2lhbCIsInRpZXIiOiJidXNpbmVzc191bmxpbWl0ZWQiLCJmZWF0dXJlcyI6WyJ1bmlkb2MiLCJ1bmlkb2MtY2xpIl0sImNyZWF0ZWRfYXQiOjE0ODU0NzUxOTksImV4cGlyZXNfYXQiOjE1MTcwMTExOTksImNyZWF0b3JfbmFtZSI6IlVuaURvYyBTdXBwb3J0IiwiY3JlYXRvcl9lbWFpbCI6InN1cHBvcnRAdW5pZG9jLmlvIn0=
		+
		JYUUjfjjpek96Rh2LoPy4LbWEHT5X46PxLyNkMyF74L/eNeLR55vcvvi2MIUtZBamCbay+YjmqZu5n6IJQWVDrImdC3b7OthoSdGMvfNSjOSuQcoV/mFpkMYin34Uwe7KM6EebzCuX2LF/LTPpdL6iYHtiWxTnF3yZwFqSgJLa8NSSSElfVLidbfQHYJSu52FTcqqWaqIjT51YiZB0Pq54YDP/jS10sRDYDe3sOpI1bfFplYkcdxPX1tK0AQKbvYCDcNbbnoKhk0EZAVSmI+kh5TdKzUn3BpQc7MP+koGrAePc3ddZF6pNzaiW1CJiO7/TmRzQioEq3Rp/h1XYkKXw==
		-----END UNIDOC LICENSE KEY-----
		`)
	*/

	// For debugging.
	common.SetLogger(common.NewConsoleLogger(common.LogLevelInfo))

	inputPath := os.Args[1]

	err := outputPdfText(inputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// outputPdfText prints out contents of PDF file to stdout.
func outputPdfText(inputPath string) error {
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

	fmt.Printf("--------------------\n")
	fmt.Printf("PDF to text extraction:\n")
	fmt.Printf("--------------------\n")
	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return err
		}

		ex, err := extractor.New(page)
		if err != nil {
			return err
		}

		text, err := ex.ExtractText()
		if err != nil {
			return err
		}

		fmt.Println("------------------------------")
		fmt.Printf("Page %d:\n", pageNum)
		fmt.Printf("\"%s\"\n", text)
		fmt.Println("------------------------------")
	}

	return nil
}
