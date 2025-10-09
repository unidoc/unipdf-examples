/*
 * Initiate empty PDF Writer, clean the Annots from pages and show analysis for newly created PDF.
 *
 * Run as: go run pdf_modify_contents.go input.pdf output.pdf
 */

package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

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

type cmdOptions struct {
	pdfPassword string
}

func main() {
	var opt cmdOptions
	flag.StringVar(&opt.pdfPassword, "password", "", "PDF Password (empty default)")
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Syntax: go run main.go [options] input.pdf output.pdf")
		os.Exit(1)
	}

	inputPath := args[0]
	outputPath := args[1]
	fmt.Printf("Input file: %s\n", inputPath)

	readerOpts := model.NewReaderOpts()
	readerOpts.Password = opt.pdfPassword

	pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, readerOpts)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	pdfWriter := model.NewPdfWriter()

	// Removes annotations from each pages.
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		panic(err)
	}
	for i := 0; i < numPages; i++ {
		page, err := pdfReader.GetPage(i + 1)
		if err != nil {
			panic(err)
		}

		if page.Annots != nil {
			fmt.Printf("-- Removing Annots on Page %d\n", i+1)
			page.Annots = nil
		}

		err = pdfWriter.AddPage(page)
		if err != nil {
			panic(err)
		}
	}

	outDoc := bytes.NewBuffer(nil)
	err = pdfWriter.Write(outDoc)
	if err != nil {
		panic(err)
	}

	pdfReader, err = model.NewPdfReader(bytes.NewReader(outDoc.Bytes()))
	if err != nil {
		panic(err)
	}

	err = pdfReader.PrintPdfObjects(os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(outputPath, outDoc.Bytes(), 0777)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Output file: %s\n", outputPath)
}
