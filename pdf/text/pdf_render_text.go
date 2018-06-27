/*
 * PDF to text: Extract all text for each page of a pdf file.
 *
 * Run as: go run pdf_extract_text.go input.pdf
 */

package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/unidoc/unidoc/common"
	"github.com/unidoc/unidoc/pdf/extractor"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run pdf_render_text.go input.pdf\n")
		os.Exit(1)
	}

	// Make sure to enter a valid license key.
	// Otherwise text is truncated and a watermark added to the text.
	// License keys are available via: https://unidoc.io
	/*
			license.SetLicenseKey(`
		-----BEGIN UNIDOC LICENSE KEY-----
		...key contents...
		-----END UNIDOC LICENSE KEY-----
		`)
	*/

	// For debugging.
	common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
	files := os.Args[1:]
	sort.Strings(files)

	for i, inputPath := range files {
		fmt.Println("======================== ^^^ ========================")
		fmt.Printf("Pdf File %3d of %d %q\n", i+1, len(files), inputPath)
		err := outputPdfText(inputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Pdf File %3d of %d %q err=%v\n",
				i+1, len(files), inputPath, err)
		}
		fmt.Println("======================== ||| ========================")
	}
	fmt.Fprintf(os.Stderr, "Done %d files\n", len(files))
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

	fmt.Println("---------------------------------------")
	fmt.Printf("PDF text rendering: %q\n", inputPath)
	fmt.Println("---------------------------------------")
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

		fmt.Printf("Page %d: %q\n", pageNum, inputPath)
		fmt.Printf("%s\n", text)
		fmt.Println("---------------------------------------")
	}

	return nil
}
