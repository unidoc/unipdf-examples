/*
 * Rotate certain page in a PDF file.  Degrees needs to be a multiple of 90.
 *
 * Run as: go run pdf_page_rotate.go input.pdf page <angle> output.pdf
 * The angle is specified in degrees.
 */

package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

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
	if len(os.Args) < 5 {
		fmt.Printf("Usage: go run pdf_page_rotate.go input.pdf <page> <angle> output.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[4]

	page, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Invalid page: %v\n", err)
		os.Exit(1)
	}
	if page < 1 {
		fmt.Println("Invalid page number specified")
		os.Exit(1)
	}

	degrees, err := strconv.ParseInt(os.Args[3], 10, 64)
	if err != nil {
		fmt.Printf("Invalid degrees: %v\n", err)
		os.Exit(1)
	}
	if degrees%90 != 0 {
		fmt.Printf("Degrees needs to be a multiple of 90\n")
		os.Exit(1)
	}

	err = rotatePage(inputPath, page, degrees, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

// Rotate all pages by 90 degrees.
func rotatePage(inputPath string, pageNum int, degrees int64, outputPath string) error {
	readerOpts := model.NewReaderOpts()
	readerOpts.LazyLoad = false

	pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, readerOpts)
	if err != nil {
		return err
	}
	defer f.Close()

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	if pageNum > numPages {
		return errors.New("Invalid page number specified")
	}

	pdfWriter, err := pdfReader.ToWriter(&model.ReaderToWriterOpts{
		PageProcessCallback: func(index int, page *model.PdfPage) error {
			if index == pageNum {
				page.Rotate = &degrees
			}

			return nil
		},
	})

	if err != nil {
		return err
	}

	pdfWriter.WriteToFile(outputPath)

	return err
}
