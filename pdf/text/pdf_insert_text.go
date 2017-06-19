/*
 * Insert text to a specific page, location in a PDF file.
 * If unsure about position, try getting the dimensions of a PDF with pdf_page_info.go first.
 *
 * Run as: go run pdf_insert_text.go input.pdf page xpos ypos "text" output.pdf
 */

package main

import (
	"fmt"
	"os"
	"strconv"

	unicommon "github.com/unidoc/unidoc/common"
	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func main() {
	if len(os.Args) < 7 {
		fmt.Printf("Usage: go run pdf_insert_text.go input.pdf page xpos ypos \"text\" output.pdf\n")
		os.Exit(1)
	}

	// Use debug logging.
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

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
	if pageNum <= 0 || pageNum > numPages {
		return fmt.Errorf("Page number out of range (%d/%d)", pageNum, numPages)
	}

	// Load the pages.
	pages := []*pdf.PdfPage{}
	for i := 0; i < numPages; i++ {
		page, err := pdfReader.GetPage(i + 1)
		if err != nil {
			return err
		}

		pages = append(pages, page)
	}

	// Add the image to the selected page.
	selPage := pages[pageNum-1]

	fmt.Printf("Page: %+v\n", selPage)
	fmt.Printf("Page mediabox: %+v\n", selPage.MediaBox)

	// Find a free name for the font.
	num := 1
	fontName := pdfcore.PdfObjectName(fmt.Sprintf("Font%d", num))
	for selPage.HasFontByName(fontName) {
		num++
		fontName = pdfcore.PdfObjectName(fmt.Sprintf("Img%d", num))
	}

	// Create the font dictionary using one of the standard 14 fonts.
	fontDict := &pdfcore.PdfObjectDictionary{
		"Type":     pdfcore.MakeName("Font"),
		"Subtype":  pdfcore.MakeName("Type1"),
		"BaseFont": pdfcore.MakeName("Helvetica"),
	}

	// Add to the page resources.
	selPage.AddFont(fontName, fontDict)

	fontSize := float64(16)

	// NextTextWriter..  textwriter package.
	creator := pdfcontent.NewContentCreator()
	creator.
		Add_BT().
		Add_Tf(fontName, fontSize).
		Add_Tm(1, 0, 0, 1, xPos, yPos).
		Add_Tj(pdfcore.PdfObjectString(text)).
		Add_ET()

	fmt.Printf("Content Str: %s\n", creator.Bytes())
	selPage.AddContentStreamByString(string(creator.Bytes()))

	// Write output.
	pdfWriter := pdf.NewPdfWriter()
	for _, page := range pages {
		err = pdfWriter.AddPage(page)
		if err != nil {
			unicommon.Log.Error("Failed to add page: %s", err)
			return err
		}
	}

	fWrite, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer fWrite.Close()

	err = pdfWriter.Write(fWrite)
	if err != nil {
		return err
	}

	return nil
}
