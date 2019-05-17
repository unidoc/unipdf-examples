/*
 * Insert text to a specific page, location in a PDF file.
 * If unsure about position, try getting the dimensions of a PDF with pdf/pages/pdf_page_info.go first or start with
 * 0,0 (upper left corner) and increase to move right, down.
 *
 * Run as: go run pdf_insert_text.go input.pdf <page> <xpos> <ypos> "text" output.pdf
 */

package main

import (
	"fmt"
	"os"
	"strconv"

	//unicommon "github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/creator"
	pdf "github.com/unidoc/unipdf/v3/model"
)

func main() {
	if len(os.Args) < 7 {
		fmt.Printf("Usage: go run pdf_insert_text.go input.pdf <page> <xpos> <ypos> \"text\" output.pdf\n")
		os.Exit(1)
	}

	// When debugging, log to console:
	//unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	inputPath := os.Args[1]
	pageNumStr := os.Args[2]
	textStr := os.Args[5]
	outputPath := os.Args[6]

	xPos, err := strconv.ParseFloat(os.Args[3], 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	yPos, err := strconv.ParseFloat(os.Args[4], 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("xPos: %d, yPos: %d\n", xPos, yPos)
	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err = addTextToPdf(inputPath, outputPath, textStr, pageNum, xPos, yPos)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

func addTextToPdf(inputPath string, outputPath string, text string, pageNum int, xPos float64, yPos float64) error {
	// Read the input pdf file.
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	c := creator.New()

	// Load the pages.
	for i := 0; i < numPages; i++ {
		page, err := pdfReader.GetPage(i + 1)
		if err != nil {
			return err
		}

		err = c.AddPage(page)
		if err != nil {
			return err
		}

		if i == pageNum || pageNum == -1 {
			p := c.NewParagraph(text)
			// Change to times bold font (default is helvetica).
			timesBold, err := pdf.NewStandard14Font("Times-Bold")
			if err != nil {
				panic(err)
			}
			p.SetFont(timesBold)
			p.SetPos(xPos, yPos)

			_ = c.Draw(p)
		}

	}

	err = c.WriteToFile(outputPath)
	return err
}
