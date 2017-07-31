/*
 * Inspect PDF object types. This example shows the capability of assessing the object types in PDF files.
 *
 * Run as: go run pdf_inspect.go input.pdf
 */

package main

import (
	"errors"
	"fmt"
	"os"
	"sort"

	pdf "github.com/unidoc/unidoc/pdf/model"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Syntax: go run pdf_inspector.go input.pdf\n")
		os.Exit(1)
	}

	// Enable debug-level logging.
	//unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	inputPath := os.Args[1]

	fmt.Printf("Input file: %s\n", inputPath)
	err := inspectPdf(inputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func inspectPdf(inputPath string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
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
			return errors.New("Unable to decrypt password protected file - need to specify pass to Decrypt")
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	fmt.Printf("PDF Num Pages: %d\n", numPages)

	objTypes, err := pdfReader.Inspect()
	if err != nil {
		return err
	}

	// Sort object types alphabetically.
	keys := []string{}
	for key, _ := range objTypes {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Output.
	fmt.Printf("Object types:\n")
	for _, key := range keys {
		fmt.Printf("- %s: %d instances\n", key, objTypes[key])
	}

	// Identify potentially risky content.
	isMalicious := false
	if count, has := objTypes["JavaScript"]; has {
		fmt.Printf("! Potentially malicious file - has %d Javascript objects\n", count)
		isMalicious = true
	}
	if count, has := objTypes["Flash"]; has {
		fmt.Printf("! Potentially malicious file - has %d Flash rich media objects\n", count)
		isMalicious = true
	}
	if count, has := objTypes["Video"]; has {
		fmt.Printf("! Potentially malicious file - has %d video objects\n", count)
		isMalicious = true
	}
	if !isMalicious {
		fmt.Printf("Most likely harmless - No javascript or rich media objects.\n")
	}

	return nil
}
