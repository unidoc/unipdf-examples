/*
 * Rotates the contents of a PDF file in accordance with each page's Rotate entry and then sets Rotate to 0.
 * I.e. flattens the rotation.  Will look the same in viewer, but when working with the PDF, the upper left
 * corner will be the origin (in unidoc coordinate system).
 *
 * Run as: go run pdf_rotate_flatten.go <input.pdf> <output.pdf>
 */

package main

import (
	"errors"
	"fmt"
	"os"

	unicommon "github.com/unidoc/unidoc/common"
	"github.com/unidoc/unidoc/pdf/creator"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func init() {
	// Use debug-mode log level.
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))
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

	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return err
	}

	// Try decrypting both with given password and an empty one if that fails.
	if isEncrypted {
		auth, err := pdfReader.Decrypt([]byte(""))
		if err != nil {
			return err
		}
		if !auth {
			return errors.New("Unable to decrypt pdf with empty pass")
		}
	}

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

		rotateDeg := int64(0)
		if page.Rotate != nil && *page.Rotate != 0 {
			rotateDeg = 360 - *page.Rotate
		}

		// Rotate the page block if needed.
		if rotateDeg != 0 {
			block.SetAngle(float64(rotateDeg))
		}

		// Set page size in creator.
		// Account for translation that is needed when rotating about the upper left corner.
		if rotateDeg == 90 || rotateDeg == 270 {
			// Swap width and height.
			c.SetPageSize(creator.PageSize{block.Height(), block.Width()})
			block.SetPos(0, block.Width())
		} else {
			c.SetPageSize(creator.PageSize{block.Width(), block.Height()})
			block.SetPos(0, 0)
		}

		c.NewPage()
		err = c.Draw(block)
		if err != nil {
			return err
		}
	}

	err = c.WriteToFile(outputPath)
	return err
}
