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
	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func init() {
	// To make the library log we just have to initialise the logger which satisfies
	// the unicommon.Logger interface, unicommon.DummyLogger is the default and
	// does not do anything. Very easy to implement your own.
	unicommon.SetLogger(unicommon.DummyLogger{})

	//Or can use a debug-level console logger:
	//unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))
	//unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelTrace))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run pdf_list_content_streams.go input.pdf\n")
		os.Exit(1)
	}

	for i := 1; i < len(os.Args); i++ {
		inputPath := os.Args[i]
		fmt.Println(inputPath)

		err := listContentStreams(inputPath)
		if err != nil {
			fmt.Println(inputPath)
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}
}

func listContentStreams(inputPath string) error {
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
		fmt.Println("Is encrypted!")
		ok, err := pdfReader.Decrypt([]byte(""))
		if err != nil {
			return err
		}
		if !ok {
			fmt.Println("Unable to decrypt with empty string - skipping")
			return nil
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
		fmt.Printf("Page %d\n", pageNum)

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return err
		}

		contentStreams, err := page.GetContentStreams()
		if err != nil {
			return err
		}
		fmt.Printf("Page %d has %d content streams:\n", pageNum, len(contentStreams))

		pageContentStr := ""

		// If the value is an array, the effect shall be as if all of the streams in the array were concatenated,
		// in order, to form a single stream.
		for _, cstream := range contentStreams {
			pageContentStr += cstream
		}
		fmt.Printf("%s\n", pageContentStr)

		cstreamParser := pdfcontent.NewContentStreamParser(pageContentStr)
		operations, err := cstreamParser.Parse()
		if err != nil {
			return err
		}

		fmt.Printf("=== Full list\n")
		for idx, op := range operations {
			fmt.Printf("Operation %d: %s - Params: %v\n", idx+1, op.Operand, op.Params)
		}
	}

	return nil
}
