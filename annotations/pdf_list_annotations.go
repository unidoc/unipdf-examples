/*
 * Lists annotations in a PDF file.
 *
 * Run as: go run pdf_list_annotations.go input.pdf [input2.pdf] ...
 */

package main

import (
	"fmt"
	"os"

	pdf "github.com/unidoc/unipdf/v3/model"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run pdf_list_annotations.go input.pdf [input2.pdf] ...\n")
		os.Exit(1)
	}

	// Enable debug-level logging.
	//unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))
	//unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelTrace))

	for _, inputPath := range os.Args[1:len(os.Args)] {
		err := listAnnotations(inputPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}
}

func listAnnotations(inputPath string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	fmt.Printf("Input file: %s\n", inputPath)

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
			fmt.Printf(" Encrypted! Need to modify code to decrypt with your password.\n")
			return nil
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}
	for i := 0; i < numPages; i++ {
		page, err := pdfReader.GetPage(i + 1)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		fmt.Printf("-- Page %d\n", i+1)
		annotations, err := page.GetAnnotations()
		if err != nil {
			return err
		}
		printAnnotations(annotations)
	}

	return nil
}

func printAnnotations(annotations []*pdf.PdfAnnotation) {
	for idx, annotation := range annotations {
		fmt.Printf(" %d. %s\n", idx+1, annotation.String())
	}
}
