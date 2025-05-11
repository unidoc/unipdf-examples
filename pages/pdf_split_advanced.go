/*
 * Advanced PDF split example: Takes into account optional content - OCProperties (rarely used).
 *
 * Run as: go run pdf_split_advanced.go input.pdf <page_from> <page_to> output.pdf
 * To get only page 1 and 2 from input.pdf and save as output.pdf run: go run pdf_split.go input.pdf 1 2 output.pdf
 */

package main

import (
	"fmt"
	"os"
	"strconv"

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
	if len(os.Args) < 5 {
		fmt.Printf("Usage: go run pdf_split_advanced.go input.pdf <page_from> <page_to> output.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]

	strSplitFrom := os.Args[2]
	splitFrom, err := strconv.Atoi(strSplitFrom)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	strSplitTo := os.Args[3]
	splitTo, err := strconv.Atoi(strSplitTo)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	outputPath := os.Args[4]

	err = splitPdf(inputPath, outputPath, splitFrom, splitTo)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

func splitPdf(inputPath string, outputPath string, pageFrom int, pageTo int) error {
	pdfWriter := model.NewPdfWriter()
	pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, nil)
	if err != nil {
		return err
	}
	defer f.Close()

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	if numPages < pageTo {
		return fmt.Errorf("numPages (%d) < pageTo (%d)", numPages, pageTo)
	}

	// Keep the OC properties intact (optional content).
	// Rarely used but can be relevant in certain cases.
	ocProps, err := pdfReader.GetOCProperties()
	if err != nil {
		return err
	}
	pdfWriter.SetOCProperties(ocProps)

	for i := pageFrom; i <= pageTo; i++ {
		pageNum := i

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return err
		}

		err = pdfWriter.AddPage(page)
		if err != nil {
			return err
		}
	}

	fWrite, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer fWrite.Close()

	err = pdfWriter.Write(fWrite)
	if err != nil {
		return err
	}

	return nil
}
