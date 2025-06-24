/*
 * Extract vector lines and other paths for each page of a PDF file.
 *
 * Run as: go run pdf_extract_lines.go input.pdf
 */

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v4/common/license"
	"github.com/unidoc/unipdf/v4/extractor"
	"github.com/unidoc/unipdf/v4/model"
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
		fmt.Printf("Usage: go run pdf_extract_lines.go input.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]

	err := outputPdfLines(inputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// outputPdfLines prints out lines of PDF file to stdout.
func outputPdfLines(inputPath string) error {
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

	// Iterate through pages.

	fmt.Printf("--------------------\n")
	fmt.Printf("PDF lines extraction:\n")
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

		fmt.Println("------------------------------")
		fmt.Printf("Page %d:\n", pageNum)

		// Extract stroke paths from the current page.
		paths, err := ex.ExtractStrokePaths()
		if err != nil {
			return err
		}

		// Print debugging info.
		for i, path := range paths {
			fmt.Printf("Path %d:\n", i)
			for j, point := range path.Points {
				fmt.Printf("Point %d: %f %f \n", j, point.X, point.Y)
			}
		}
	}

	return nil
}
