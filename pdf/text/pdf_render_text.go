/*
 * PDF to text: Extract all text for each page of a pdf file.
 *
 * Run as: go run pdf_extract_text.go testdata/*.pdf
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
		fmt.Printf("Usage: go run pdf_render_text.go testdata/*.pdf\n")
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
		fmt.Println("========================= ^^^ =========================")
		pdfReader, numPages, err := getReader(inputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\to====> Pdf File %3d of %d err=%v\n",
				i+1, len(files), err)
			continue
		}

		version := pdfReader.PdfVersion()
		// We are currently not interested in old PDF files. If you are, comment out these lines.
		if version == "1.0" || version == "1.1" || version == "1.2" {
			continue
		}

		fmt.Fprintf(os.Stderr, "Pdf File %3d of %d (%3s) %3d pages %q \n",
			i+1, len(files), pdfReader.PdfVersion(), numPages, inputPath)

		err = outputPdfText(inputPath, pdfReader, numPages)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\tx====> Pdf File %3d of %d %q err=%v\n",
				i+1, len(files), inputPath, err)
		}
		fmt.Println("========================= ~~~ =========================")
	}
	fmt.Fprintf(os.Stderr, "Done %d files\n", len(files))
}

// getReader returns a PdfReader and the number of pages for PDF file `inputPath`.
func getReader(inputPath string) (pdfReader *pdf.PdfReader, numPages int, err error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return
	}
	defer f.Close()

	pdfReader, err = pdf.NewPdfReader(f)
	if err != nil {
		return
	}
	numPages, err = pdfReader.GetNumPages()
	return
}

// outputPdfText prints out text of PDF file `inputPath` to stdout.
// `pdfReader` is a previously opened PdfReader of `inputPath`
func outputPdfText(inputPath string, pdfReader *pdf.PdfReader, numPages int) error {
	for pageNum := 1; pageNum <= numPages; pageNum++ {

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

		fmt.Printf("Page %d of %d: %q\n", pageNum, numPages, inputPath)
		fmt.Printf("%s\n", text)
		fmt.Println("------------------------- ... -------------------------")
	}

	return nil
}
