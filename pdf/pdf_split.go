/*
 * Splits PDF files, tries to decrypt encrypted documents with an empty password
 * as best effort.
 *
 * Run as: go run pdf_split.go input.pdf <page_from> <page_to> output.pdf
 * To get only page 1 and 2 from input.pdf and save as output.pdf run: go run pdf_split.go input.pdf 1 2 output.pdf
 */

package main

import (
	"fmt"
	"os"
	"strconv"

	unicommon "github.com/unidoc/unidoc/common"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func init() {
	// To make the library log we just have to initialise the logger which satisfies
	// the unicommon.Logger interface, unicommon.DummyLogger is the default and
	// does not do anything. Very easy to implement your own.
	// unicommon.SetLogger(unicommon.DummyLogger{})

	// Use debug-level console logger.
	//unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelTrace))
}

func main() {
	if len(os.Args) < 5 {
		fmt.Printf("Usage: go run pdf_split.go input.pdf <page_from> <page_to> output.pdf\n")
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
		return err
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
