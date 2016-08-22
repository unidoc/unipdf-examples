/*
 * List all content streams for all pages in a pdf file.
 *
 * Run as: go run pdf_print_content_streams.go input.pdf
 */

package main

import (
	"fmt"
	"os"

	unicommon "github.com/unidoc/unidoc/common"
	unilicense "github.com/unidoc/unidoc/license"
	unipdf "github.com/unidoc/unidoc/pdf"
)

func initUniDoc(licenseKey string) error {
	if len(licenseKey) > 0 {
		err := unilicense.SetLicenseKey(licenseKey)
		if err != nil {
			return err
		}
	}

	// To make the library log we just have to initialise the logger which satisfies
	// the unicommon.Logger interface, unicommon.DummyLogger is the default and
	// does not do anything. Very easy to implement your own.
	unicommon.SetLogger(unicommon.DummyLogger{})
	//unicommon.SetLogger(unicommon.ConsoleLogger{})

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run pdf_list_content_streams.go input.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]

	err := initUniDoc("")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err = listContentStreams(inputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func listContentStreams(inputPath string) error {
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

	fmt.Printf("--------------------\n")
	fmt.Printf("Content streams:\n")
	fmt.Printf("--------------------\n")
	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPageAsPdfPage(pageNum)
		if err != nil {
			return err
		}

		contentStreams, err := page.GetContentStreams()
		if err != nil {
			return err
		}
		fmt.Printf("Page %d has %d content streams:\n", pageNum, len(contentStreams))
		for idx, cstream := range contentStreams {
			fmt.Printf("Page %d - content stream %d:\n", pageNum, idx+1)
			fmt.Printf("%s\n", cstream)

			cstreamParser := unipdf.NewContentStreamParser(cstream)
			operations, err := cstreamParser.Parse()
			if err != nil {
				return err
			}
			fmt.Printf("=== Full list\n")
			for idx, op := range operations {
				fmt.Printf("Operation %d: %s\n", idx+1, op.Operand)
			}
			fmt.Printf("=== Text related\n")
			inText := false
			for idx, op := range operations {
				if op.Operand == "BT" {
					inText = true
				} else if op.Operand == "ET" {
					inText = false
				}
				if op.Operand == "TJ" {
					if len(op.Params) < 1 {
						continue
					}
					paramList := op.Params[0].(*unipdf.PdfObjectArray)
					for _, obj := range *paramList {
						if strObj, ok := obj.(*unipdf.PdfObjectString); ok {
							fmt.Printf("%s", *strObj)
						}
					}
					fmt.Printf("\n")
				}

				if inText == true {
					fmt.Printf("%d. %v\n", idx+1, op)
				}
			}
		}
	}

	return nil
}
