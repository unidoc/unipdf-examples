/*
 * Show contents of all objects in a PDF file. Handy for debugging UniDoc programs
 *
 * Run as: go run pdf_all_objects.go input.pdf
 */

package main

import (
	"errors"
	"fmt"
	"os"

	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Syntax: go run pdf_all_objects.go input.pdf")
		os.Exit(1)
	}

	// Enable debug-level logging.
	//unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	inputPath := os.Args[1]

	fmt.Printf("Input file: %s\n", inputPath)
	err := inspectPdf(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
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

	objNums := pdfReader.GetObjectNums()

	// Output.
	fmt.Printf("%d PDF objects:\n", len(objNums))
	for i, objNum := range objNums {
		obj, err := pdfReader.GetIndirectObjectByNumber(objNum)
		if err != nil {
			return err
		}
		fmt.Println("=========================================================")
		fmt.Printf("%3d: %d 0 %T\n", i, objNum, obj)
		if stream, is := obj.(*pdfcore.PdfObjectStream); is {
			decoded, err := pdfcore.DecodeStream(stream)
			if err != nil {
				return err
			}
			fmt.Printf("Decoded:\n%s\n", decoded)
		} else if indObj, is := obj.(*pdfcore.PdfIndirectObject); is {
			fmt.Printf("%T\n", indObj.PdfObject)
			fmt.Printf("%s\n", indObj.PdfObject.String())
		}

	}

	return nil
}
