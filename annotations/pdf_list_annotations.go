/*
 * Lists annotations in a PDF file.
 *
 * Run as: go run pdf_list_annotations.go input.pdf [input2.pdf] ...
 */

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v4/common/license"
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
		fmt.Printf("Usage: go run pdf_list_annotations.go input.pdf [input2.pdf] ...\n")
		os.Exit(1)
	}

	for _, inputPath := range os.Args[1:len(os.Args)] {
		err := listAnnotations(inputPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}
}

func listAnnotations(inputPath string) error {
	pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, nil)
	if err != nil {
		return err
	}
	defer f.Close()

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

func printAnnotations(annotations []*model.PdfAnnotation) {
	for idx, annotation := range annotations {
		fmt.Printf(" %d. %s\n", idx+1, annotation.String())
	}
}
