/*
 * Unlocks PDF files, tries to decrypt encrypted documents with the given password,
 * if that fails it tries an empty password as best effort.
 *
 * Run as: go run pdf_unlock.go input.pdf <password> output.pdf
 */

package main

import (
	"fmt"
	"os"

	pdf "github.com/unidoc/unipdf/v3/model"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Usage: go run pdf_unlock.go input.pdf <password> output.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	password := os.Args[2]
	outputPath := os.Args[3]

	err := unlockPdf(inputPath, outputPath, password)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

func unlockPdf(inputPath string, outputPath string, password string) error {
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

	// Try decrypting both with given password and an empty one if that fails.
	if isEncrypted {
		auth, err := pdfReader.Decrypt([]byte(password))
		if err != nil {
			return err
		}
		if !auth {
			return fmt.Errorf("Wrong password")
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	for i := 0; i < numPages; i++ {
		pageNum := i + 1

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
