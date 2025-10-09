/*
 * Show contents of all objects in a PDF file. Handy for debugging UniPDF programs
 *
 * Run as: go run pdf_all_objects.go input.pdf
 */

package main

import (
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
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Syntax: go run pdf_all_objects.go [options] input.pdf")
		os.Exit(1)
	}

	inputPath := args[0]

	fmt.Printf("Input file: %s\n", inputPath)
	readerOpts := model.NewReaderOpts()
	readerOpts.Password = opt.pdfPassword

	pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, readerOpts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	err = pdfReader.PrintPdfObjects(os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
