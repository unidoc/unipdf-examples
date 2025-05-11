/*
 * Rotate pages in a PDF file using global flag instead of
 * rotating each page one by one.
 * Degrees needs to be a multiple of 90.
 *
 * Run as: go run pdf_rotate.go input.pdf <angle> output.pdf
 * The angle is specified in degrees.
 */

package main

import (
	"fmt"
	"os"
	"strconv"

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
		fmt.Printf("Usage: go run pdf_rotate.go input.pdf <angle> output.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[3]

	degrees, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		fmt.Printf("Invalid degrees: %v\n", err)
		os.Exit(1)
	}
	if degrees%90 != 0 {
		fmt.Printf("Degrees needs to be a multiple of 90\n")
		os.Exit(1)
	}

	err = rotatePdf(inputPath, degrees, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

// Rotate all pages by degrees.
func rotatePdf(inputPath string, degrees int64, outputPath string) error {
	pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, nil)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfWriter, err := pdfReader.ToWriter(&model.ReaderToWriterOpts{})
	if err != nil {
		return nil
	}

	// Rotate all page degrees.
	err = pdfWriter.SetRotation(degrees)
	if err != nil {
		return nil
	}

	pdfWriter.WriteToFile(outputPath)

	return err
}
