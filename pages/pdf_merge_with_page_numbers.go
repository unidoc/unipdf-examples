/*
 * Merging of PDF files and add footer with page numbers to each pages.
 * Simply loads all pages for each file, add page with page numbers using creator and writes to the output file.
 *
 * Run as: go run pdf_merge_with_page_numbers.go output.pdf input1.pdf input2.pdf input3.pdf ...
 */

package main

import (
	"fmt"
	"os"

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
}

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Requires at least 3 arguments: output_path and 2 input paths\n")
		fmt.Printf("Usage: go run pdf_merge_with_page_numbers.go output.pdf input1.pdf input2.pdf input3.pdf ...\n")
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

	err := mergePdf(inputPaths, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

func mergePdf(inputPaths []string, outputPath string) error {
	c := creator.New()

	for fi, inputPath := range inputPaths {
		fileNo := fi + 1
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

			err = c.AddPage(page)
			if err != nil {
				return err
			}

			// Draw paragraph to page.
			p := c.NewStyledParagraph()
			p.SetTextAlignment(creator.TextAlignmentCenter)

			chunk := p.Append(fmt.Sprintf("I am extra content of file no. %d and page %d", fileNo, pageNum))
			chunk.Style.FontSize = 8
			if err := c.Draw(p); err != nil {
				return err
			}
		}
	}

	// Draw footer with page numbers and total pages to pdf pages.
	c.DrawFooter(func(block *creator.Block, args creator.FooterFunctionArgs) {
		p := c.NewStyledParagraph()
		p.SetTextAlignment(creator.TextAlignmentCenter)

		chunk := p.Append(fmt.Sprintf("Page %d of %d", args.PageNum, args.TotalPages))
		chunk.Style.FontSize = 8
		chunk.Style.Color = creator.ColorRGBFrom8bit(63, 68, 76)

		block.Draw(p)
	})

	err := c.WriteToFile(outputPath)
	if err != nil {
		return err
	}

	return nil
}
