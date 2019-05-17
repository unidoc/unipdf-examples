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

	unicommon "github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Outputs Document Information Dictionary information from PDF files\n")
		fmt.Printf("Usage: go run pdf_metadata_get_docinfo.go input1.pdf [input2.pdf] ...\n")
		os.Exit(1)
	}

	// Enable debug-level logging.
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

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

	trailerDict, err := pdfReader.GetTrailer()
	if err != nil {
		return err
	}
	if trailerDict == nil {
		fmt.Printf("No trailer dictionary -> No DID dictionary\n")
		// Note: not returning an error - it is not guaranteed for every PDF to have the DID dictionary.
		return nil
	}

	// XXX/FIXME: Much of the clunky type casting and tracing is being improved in v3.

	var infoDict *core.PdfObjectDictionary

	infoObj := trailerDict.Get("Info")
	switch t := infoObj.(type) {
	case *core.PdfObjectReference:
		infoRef := t
		infoObj, err = pdfReader.GetIndirectObjectByNumber(int(infoRef.ObjectNumber))
		infoObj = core.TraceToDirectObject(infoObj)
		if err != nil {
			return err
		}
		infoDict, _ = infoObj.(*core.PdfObjectDictionary)
	case *core.PdfObjectDictionary:
		infoDict = t
	}

	if infoDict == nil {
		fmt.Printf("DID dictionary not present\n")
		return nil
	}

	di := pdfDocInfo{
		Filename: path.Base(inputPath),
		NumPages: numPages,
	}

	if str, has := infoDict.Get("Title").(*core.PdfObjectString); has {
		di.Title = str.String()
	}

	if str, has := infoDict.Get("Author").(*core.PdfObjectString); has {
		di.Author = str.String()
	}

	if str, has := infoDict.Get("Keywords").(*core.PdfObjectString); has {
		di.Keywords = str.String()
	}

	if str, has := infoDict.Get("Creator").(*core.PdfObjectString); has {
		di.Creator = str.String()
	}

	if str, has := infoDict.Get("Producer").(*core.PdfObjectString); has {
		di.Producer = str.String()
	}

	if str, has := infoDict.Get("CreationDate").(*core.PdfObjectString); has {
		di.CreationDate = str.String()
	}

	if str, has := infoDict.Get("ModDate").(*core.PdfObjectString); has {
		di.ModDate = str.String()
	}

	if name, has := infoDict.Get("Trapped").(*core.PdfObjectName); has {
		di.Trapped = name.String()
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
}
