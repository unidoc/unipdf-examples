/*
 * Outputs multiple pages (4) per page to an output PDF from an input PDF.
 * Showcases page templating by loading pages as Blocks and manipulating with the creator package.
 *
 * Run as: go run pdf_4up.go <input.pdf> <output.pdf>
 */

package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

func init() {
	// Use debug-mode log level.
	common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: go run pdf_4up.go <input.pdf> <output.pdf>\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	err := multiplePagesPerPage(inputPath, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

// Load an input PDF and output as n-pages per page in the output.
func multiplePagesPerPage(inputPath, outputPath string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
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

		pos := i % 4
		if pos == 0 {
			c.NewPage()
		}

		pageWidth := c.Context().PageWidth
		pageHeight := c.Context().PageHeight
		block.ScaleToWidth(0.3 * pageWidth)

		var xPos, yPos float64
		switch pos {
		case 0:
			xPos, yPos = 0.1*pageWidth, 0.2*pageHeight
		case 1:
			xPos, yPos = 0.6*pageWidth, 0.2*pageHeight
		case 2:
			xPos, yPos = 0.1*pageWidth, 0.6*pageHeight
		case 3:
			xPos, yPos = 0.6*pageWidth, 0.6*pageHeight
		}
		block.SetPos(xPos, yPos)

		blockWidth, blockHeight := block.RotatedSize()
		dx := blockWidth - block.Width()
		dy := blockHeight - block.Height()

		rect := c.NewRectangle(xPos-dx/2, yPos-dy/2, blockWidth, blockHeight)
		rect.SetBorderWidth(1.0)
		rect.SetBorderColor(creator.ColorBlack)

		err = c.Draw(block)
		if err != nil {
			return err
		}
		err = c.Draw(rect)
		if err != nil {
			return err
		}
	}

	err = c.WriteToFile(outputPath)
	return err
}
