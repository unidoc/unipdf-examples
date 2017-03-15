/*
 * Detect scanned PDF files by looking through the object types and determining whether it is likely to be a scanned file.
 *
 * Run as: go run pdf_detect_scanned.go input.pdf
 */

package main

import (
	"fmt"
	"os"

	unicommon "github.com/unidoc/unidoc/common"
	unipdf "github.com/unidoc/unidoc/pdf/model"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Syntax: go run pdf_detect_scanned.go input1.pdf input2.pdf ...\n")
		os.Exit(1)
	}

	// Enable debug-level logging.
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

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

	pdfReader, err := unipdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return err
	}

	// Try decrypting with an empty one.
	if isEncrypted {
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
