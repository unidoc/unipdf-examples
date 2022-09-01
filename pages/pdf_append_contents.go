/*
 * Basic append contents of PDF files.
 * Simply loads all pages for each file, add it to block then draw block to the output file.
 *
 * Run as: go run pdf_append_contents.go output.pdf input1.pdf input2.pdf input3.pdf ...
 */

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}
	common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
}

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Requires at least 3 arguments: output_path and 2 input paths\n")
		fmt.Printf("Usage: go run pdf_append_contents.go output.pdf input1.pdf input2.pdf input3.pdf ...\n")
		os.Exit(0)
	}

	outputPath := ""
	inputPaths := []string{}

	// Sanity check the input arguments.
	for i, arg := range os.Args {
		if i == 0 {
			continue
		} else if i == 1 {
			outputPath = arg
			continue
		}

		inputPaths = append(inputPaths, arg)
	}

	err := appendPdfContents(inputPaths, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

func appendPdfContents(inputPaths []string, outputPath string) error {
	c := creator.New()

	// You need to costumize the page size if you want the contents fit into single pages.
	c.SetPageSize(creator.PageSize{500 * creator.PPMM, 600 * creator.PPMM})

	for _, inputPath := range inputPaths {
		pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, nil)
		if err != nil {
			return err
		}
		defer f.Close()

		numPages, err := pdfReader.GetNumPages()
		if err != nil {
			return err
		}

		for i := 0; i < numPages; i++ {
			pageNum := i + 1

			page, err := pdfReader.GetPage(pageNum)
			if err != nil {
				return err
			}

			block, err := creator.NewBlockFromPage(page)
			if err != nil {
				return err
			}

			if err := c.Draw(block); err != nil {
				return err
			}
		}
	}

	err := c.WriteToFile(outputPath)
	if err != nil {
		return err
	}
	return nil
}
