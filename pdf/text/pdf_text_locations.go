/*
 * PDF to text: Extract all text for each page of a pdf file.
 *
 * Run as: go run pdf_extract_text.go input.pdf
 */

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/unidoc/unidoc/common"
	"github.com/unidoc/unidoc/pdf/extractor"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

const (
	usage                 = "Usage: go run pdf_render_text.go testdata/*.pdf\n"
	badFilesPath          = "bad.files"
	defaultNormalizeWidth = 60
)

var ErrBadText = errors.New("could not decode text")

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run pdf_text_locations.go input.pdf\n")
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
	var debug, trace bool
	maxPages := 5
	maxText := 500
	maxLocations := 10

	flag.BoolVar(&debug, "d", false, "Print debugging information.")
	flag.BoolVar(&trace, "e", false, "Print detailed debugging information.")
	flag.IntVar(&maxPages, "p", maxPages, "Maximum number of pages to extract.")
	flag.IntVar(&maxText, "t", maxText, "Maximum number of characters of text to show per page.")
	flag.IntVar(&maxLocations, "l", maxLocations, "Maximum number of locations to show per page.")

	makeUsage(usage)

	makeUsage(usage)
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}
	if trace {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelTrace))
	} else if debug {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
	} else {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelInfo))
	}

	inputPath := args[0]

	err := outputPdfTextLocations(inputPath, maxPages, maxText, maxLocations)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// outputPdfTextLocations prints out the text of the PDF file and the text locationsto stdout.
func outputPdfTextLocations(inputPath string, maxPages, maxText, maxLocations int) error {
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

	fmt.Println(separator)
	fmt.Printf("PDF text location extraction: %d pages\n", numPages)
	fmt.Println(separator)

	for pageNum := 1; pageNum <= numPages; pageNum++ {
		if maxPages >= 0 && pageNum > maxPages {
			break
		}
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return err
		}

		ex, err := extractor.New(page)
		if err != nil {
			return err
		}

		pageText, _, _, err := ex.ExtractPageText()
		if err != nil {
			return err
		}
		text, locations := pageText.ToTextLocation()

		fmt.Println(separText)
		fmt.Printf("Page %d: %d chars %d locations\n", pageNum, len(text), len(locations))
		fmt.Printf("\"%s\"\n", truncate(text, maxText))
		fmt.Println(separLocs)
		for i, loc := range locations {
			if maxLocations >= 0 && i >= maxLocations {
				break
			}
			fmt.Printf("%6d: %s\n", i, loc)
		}
		fmt.Println(separator)
	}

	fmt.Printf("PDF text location extraction: %d pages\n", numPages)

	return nil
}

const (
	separator = "---------------------------------------------------------------"
	separText = "--------------------------- TEXT ------------------------------"
	separLocs = "------------------------- LOCATIONS ---------------------------"
)

// makeUsage updates flag.Usage to include usage message `msg`.
func makeUsage(msg string) {
	usage := flag.Usage
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, msg)
		usage()
	}
}

func truncate(text string, n int) string {
	if len(text) < n {
		return text
	}
	return text[:n]
}
