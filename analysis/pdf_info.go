/*
 * Prints basic PDF info: number of pages and encryption status.
 *
 * Run as: go run pdf_info.go input1.pdf [input2.pdf] ...
 */

package main

import (
	"fmt"
	"os"

	unicommon "github.com/unidoc/unipdf/v3/common"
	unipdf "github.com/unidoc/unipdf/v3/model"
)

type PdfProperties struct {
	IsEncrypted bool
	CanView     bool // Is the document viewable without password?
	NumPages    int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Print out basic properties of PDF files\n")
		fmt.Printf("Usage: go run pdf_info.go input.pdf [input2.pdf] ...\n")
		os.Exit(1)
	}

	// Enable debug-level logging.
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	for _, inputPath := range os.Args[1:len(os.Args)] {
		fmt.Printf("Input file: %s\n", inputPath)

		ret, err := getPdfProperties(inputPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf(" Num Pages: %d\n", ret.NumPages)
		fmt.Printf(" Is Encrypted: %t\n", ret.IsEncrypted)
		fmt.Printf(" Is Viewable (without pass): %t\n", ret.CanView)
	}
}

func getPdfProperties(inputPath string) (*PdfProperties, error) {
	ret := PdfProperties{}

	f, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	pdfReader, err := unipdf.NewPdfReader(f)
	if err != nil {
		return nil, err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return nil, err
	}

	ret.IsEncrypted = isEncrypted
	ret.CanView = true

	// Try decrypting with an empty one.
	if isEncrypted {
		auth, err := pdfReader.Decrypt([]byte(""))
		if err != nil {
			return nil, err
		}
		ret.CanView = auth
		return &ret, nil
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return nil, err
	}
	ret.NumPages = numPages

	return &ret, nil
}
