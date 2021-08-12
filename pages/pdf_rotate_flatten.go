/*
 * Rotates the contents of a PDF file in accordance with each page's Rotate entry and then sets Rotate to 0.
 * I.e. flattens the rotation.  Will look the same in viewer, but when working with the PDF, the upper left
 * corner will be the origin (in unidoc coordinate system).
 *
 * Run as: go run pdf_rotate_flatten.go <input.pdf> <output.pdf>
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
	if len(os.Args) < 3 {
		fmt.Printf("Usage: go run pdf_rotate_flatten.go <input.pdf> <output.pdf>\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	err := rotateFlattenPdf(inputPath, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

// Flatten the PDF's rotation flags.  For each page rotate page contents with page.Rotate, then set page.Rotate to 0.
func rotateFlattenPdf(inputPath, outputPath string) error {
	pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, nil)
	if err != nil {
		return err
	}
	defer f.Close()

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	c := creator.New()
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

		var rotateDeg int64
		if page.Rotate != nil && *page.Rotate != 0 {
			rotateDeg = -*page.Rotate
		}

		// Rotate the page block. If the page is not rotated, this is a no-op.
		block.SetAngle(float64(rotateDeg))
		w, h := block.RotatedSize()
		block.SetPos((w-block.Width())/2, (h-block.Height())/2)

		c.SetPageSize(creator.PageSize{w, h})
		c.NewPage()
		if err = c.Draw(block); err != nil {
			return err
		}
	}

	err = c.WriteToFile(outputPath)
	return err
}
