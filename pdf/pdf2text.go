/*
 * Example for UniDoc v2.
 * PDF to text: Extract all text for each page of a pdf file.
 *
 * N.B. Only outputs character codes as seen in the content stream.  Need to account for encoding to get readable
 * text in many cases.
 *
 * Run as: go run pdf_extract_text.go input.pdf
 */

package main

import (
	"bytes"
	"fmt"
	"os"

	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run pdf_extract_text.go input.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]

	text, err := getContentStreams(inputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("----------------------------------------")
	fmt.Printf("Text %d chars\n", len(text))
	fmt.Println("----------------------------------------")
	fmt.Printf("%s\n", text)
	fmt.Println("----------------------------------------")
}

func getContentStreams(inputPath string) (string, error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return "", err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return "", err
	}

	if isEncrypted {
		_, err = pdfReader.Decrypt([]byte(""))
		if err != nil {
			return "", err
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return "", err
	}

	var b bytes.Buffer

	fmt.Printf("--------------------\n")
	fmt.Printf("PDF to text extraction:\n")
	fmt.Printf("--------------------\n")

	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return "", err
		}

		contentStreams, err := page.GetContentStreams()
		if err != nil {
			return b.String(), err
		}

		// If the value is an array, the effect shall be as if all of the streams in the array were concatenated,
		// in order, to form a single stream.
		pageContentStr := ""
		for _, cstream := range contentStreams {
			pageContentStr += cstream
		}

		fmt.Printf("Page %d - content streams %d: %d chars\n", pageNum, len(contentStreams), b.Len())
		cstreamParser := pdfcontent.NewContentStreamParser(pageContentStr)
		txt, err := cstreamParser.ExtractText()
		if err != nil {
			return b.String(), err
		}
		b.WriteString(txt)
		b.WriteString(" ")
		fmt.Printf("\"%s\"\n", txt)
	}

	return b.String(), nil
}
