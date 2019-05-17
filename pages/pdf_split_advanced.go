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

	//unicommon "github.com/unidoc/unipdf/v3/common"
	pdf "github.com/unidoc/unipdf/v3/model"
)

func init() {
	// When debugging: use debug-level console logger.
	//unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))
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
	pdfWriter := pdf.NewPdfWriter()

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

	if isEncrypted {
		_, err = pdfReader.Decrypt([]byte(""))
		if err != nil {
			return err
		}
	}

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
