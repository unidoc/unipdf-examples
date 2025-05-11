/*
 * Outputs information from the Document Information Dictionary for the PDF.
 *
 * Run as: go run pdf_metadata_get_docinfo.go input1.pdf [input2.pdf] ...
 */

package main

import (
	"fmt"
	"os"
	"path"

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
		fmt.Printf("Outputs Document Information Dictionary information from PDF files\n")
		fmt.Printf("Usage: go run pdf_metadata_get_docinfo.go input1.pdf [input2.pdf] ...\n")
		os.Exit(1)
	}

	fmt.Printf("Document Information Dictionary analysis\n")
	for _, inputPath := range os.Args[1:len(os.Args)] {
		fmt.Printf("Input file: %s\n", inputPath)

		err := printPdfDocInfo(inputPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}
}

func printPdfDocInfo(inputPath string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	pdfInfo, err := pdfReader.GetPdfInfo()
	if err != nil {
		return err
	}

	di := pdfDocInfo{
		Filename: path.Base(inputPath),
		NumPages: numPages,
	}

	if pdfInfo.Title != nil {
		di.Title = pdfInfo.Title.Decoded()
	}

	if pdfInfo.Author != nil {
		di.Author = pdfInfo.Author.Decoded()
	}

	if pdfInfo.Subject != nil {
		di.Subject = pdfInfo.Subject.Decoded()
	}

	if pdfInfo.Keywords != nil {
		di.Keywords = pdfInfo.Keywords.Decoded()
	}

	if pdfInfo.Creator != nil {
		di.Creator = pdfInfo.Creator.Decoded()
	}

	if pdfInfo.Producer != nil {
		di.Producer = pdfInfo.Producer.Decoded()
	}

	if pdfInfo.CreationDate != nil {
		di.CreationDate = pdfInfo.CreationDate.ToGoTime().String()
	}

	if pdfInfo.ModifiedDate != nil {
		di.ModDate = pdfInfo.ModifiedDate.ToGoTime().String()
	}

	if pdfInfo.Trapped != nil {
		di.Trapped = pdfInfo.Trapped.String()
	}

	customInfoKeys := pdfInfo.CustomKeys()
	di.CustomInfo = make(map[string]string, len(customInfoKeys))

	for _, key := range customInfoKeys {
		di.CustomInfo[key] = pdfInfo.GetCustomInfo(key).Decoded()
	}

	di.print()

	return nil
}

// pdfDocInfo is a summary of PDF document information, including Document Information Dictionary infromation.
type pdfDocInfo struct {
	Filename string
	NumPages int

	Title        string
	Author       string
	Subject      string
	Keywords     string
	Creator      string
	Producer     string
	CreationDate string
	ModDate      string
	Trapped      string
	CustomInfo   map[string]string
}

// print prints a summary of the PDF document information.
func (di pdfDocInfo) print() {
	fmt.Printf("Filename: %s\n", di.Filename)
	fmt.Printf("  Pages: %d\n", di.NumPages)
	fmt.Printf("  Title: %s\n", di.Title)
	fmt.Printf("  Author: %s\n", di.Author)
	fmt.Printf("  Subject: %s\n", di.Subject)
	fmt.Printf("  Keywords: %s\n", di.Keywords)
	fmt.Printf("  Creator: %s\n", di.Creator)
	fmt.Printf("  Producer: %s\n", di.Producer)
	fmt.Printf("  CreationDate: %s\n", di.CreationDate)
	fmt.Printf("  ModDate: %s\n", di.ModDate)
	fmt.Printf("  Trapped: %s\n", di.Trapped)

	if di.CustomInfo != nil {
		for k, v := range di.CustomInfo {
			fmt.Printf("  %s: %s\n", k, v)
		}
	}
}
