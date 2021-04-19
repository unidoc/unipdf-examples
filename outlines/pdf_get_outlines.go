/*
 * Retrieves outlines (bookmarks) from a PDF file and prints out in JSON format.
 * Note: The JSON output can be used with the related pdf_set_outlines.go example to
 * apply outlines to a PDF file.
 *
 * Run as: go run pdf_get_outlines.go input.pdf > outlines.json
 */

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/model"
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
		fmt.Printf("Usage:  go run pdf_get_outlines.go input.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]

	fmt.Printf("Input file: %s\n", inputPath)

	err := getOutlines(inputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func getOutlines(inputPath string) error {
	readerOpts := model.NewReaderOpts()
	readerOpts.LazyLoad = false

	pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, readerOpts)
	if err != nil {
		return err
	}
	defer f.Close()

	outlines, err := pdfReader.GetOutlines()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(outlines, "", "    ")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", data)

	return nil
}
