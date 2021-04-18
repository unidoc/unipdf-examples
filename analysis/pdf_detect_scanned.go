/*
 * Detect scanned PDF files by looking through the object types and determining whether it is likely to be a scanned file.
 *
 * Run as: go run pdf_detect_scanned.go input.pdf
 */

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/model"
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
		fmt.Printf("Syntax: go run pdf_detect_scanned.go input1.pdf input2.pdf ...\n")
		os.Exit(1)
	}

	for _, inputPath := range os.Args[1:] {
		err := detectScanned(inputPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

func detectScanned(inputPath string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return err
	}

	if isEncrypted {
		// Decrypt if needed.  Put your password in the empty string below.
		auth, err := pdfReader.Decrypt([]byte(""))
		if err != nil {
			return err
		}
		if !auth {
			fmt.Printf("%s - Unable to access (encrypted)\n", inputPath)
			return nil
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}
	fmt.Printf("%s (%d pages) - ", inputPath, numPages)
	objTypes, err := pdfReader.Inspect()
	if err != nil {
		return err
	}

	fontObjs, ok := objTypes["Font"]
	if !ok || fontObjs < 2 {
		fmt.Printf("SCANNED!\n")
	} else {
		fmt.Printf("not scanned (has text objects)\n")
	}

	return nil
}
