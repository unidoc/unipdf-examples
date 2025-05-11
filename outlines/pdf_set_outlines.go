/*
 * Applies outlines to a PDF file. The files are read from a JSON formatted file,
 * which can be created via pdf_get_outlines which outputs outlines for an input PDF file
 * in the JSON format.
 *
 * Run as: go run pdf_set_outlines.go input.pdf outlines.json output.pdf
 */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Usage: go run pdf_set_outlines.go input.pdf outlines.json output.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outlinesPath := os.Args[2]
	outPath := os.Args[3]

	fmt.Printf("Input file: %s\n", inputPath)
	fmt.Printf("Outlines file (JSON): %s\n", outlinesPath)
	fmt.Printf("Output file: %s\n", outPath)

	err := applyOutlines(inputPath, outlinesPath, outPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func applyOutlines(inputPath, outlinesPath, outPath string) error {
	data, err := ioutil.ReadFile(outlinesPath)
	if err != nil {
		return err
	}

	var newOutlines model.Outline
	err = json.Unmarshal(data, &newOutlines)
	if err != nil {
		return err
	}

	pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, nil)
	if err != nil {
		return err
	}
	defer f.Close()

	// Don't copy document outlines.
	opt := &model.ReaderToWriterOpts{
		SkipOutlines: true,
	}

	// Generate a PdfWriter instance from existing PdfReader.
	pdfWriter, err := pdfReader.ToWriter(opt)
	if err != nil {
		return err
	}

	// Add the new document outline.
	pdfWriter.AddOutlineTree(newOutlines.ToOutlineTree())

	return pdfWriter.WriteToFile(outPath)
}
