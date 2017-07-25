/*
 * Draw a line in a new PDF file.
 *
 * Run as: go run pdf_draw_line.go <x1> <y1> <x2> <y2> output.pdf
 * The dimensions of the PDF file are 595.276 x 841.89 px.
 * The x, y coordinates are in the PDF coordinate's system, where 0,0 is in the lower left corner.
 */

package main

import (
	"fmt"
	"os"
	"strconv"

	unicommon "github.com/unidoc/unidoc/common"
	"github.com/unidoc/unidoc/pdf/contentstream/draw"
	pdfcore "github.com/unidoc/unidoc/pdf/core"
	unipdf "github.com/unidoc/unidoc/pdf/model"
)

func init() {
	// Debug log mode.
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))
}

func main() {
	if len(os.Args) < 6 {
		fmt.Printf("go run pdf_draw_line.go <x1> <y1> <x2> <y2> output.pdf\n")
		os.Exit(1)
	}

	x1, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	y1, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	x2, err := strconv.ParseFloat(os.Args[3], 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	y2, err := strconv.ParseFloat(os.Args[4], 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	outputPath := os.Args[5]

	err = drawPdfLineToFile(x1, y1, x2, y2, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

func drawPdfLineToFile(x1, y1, x2, y2 float64, outputPath string) error {
	pdfWriter := unipdf.NewPdfWriter()

	// TODO: Inteface enhancments: what about unipdf.NewPdfPage(x, y, width, height) ?
	page := unipdf.NewPdfPage()
	bbox := unipdf.PdfRectangle{0, 0, 595.276, 841.89}
	page.MediaBox = &bbox

	// Define line path and style.
	line := draw.Line{
		X1:               x1,
		Y1:               y1,
		X2:               x2,
		Y2:               y2,
		Opacity:          1.0,
		LineEndingStyle1: draw.LineEndingStyleNone,
		LineEndingStyle2: draw.LineEndingStyleArrow,
		LineColor:        unipdf.NewPdfColorDeviceRGB(1, 0, 0),
		LineWidth:        2.0,
	}

	content, linebbox, err := line.Draw("")
	if err != nil {
		return err
	}

	fmt.Printf("Line bbox: %v\n", linebbox)

	page.SetContentStreams([]string{string(content)}, pdfcore.NewFlateEncoder())

	err = pdfWriter.AddPage(page)
	if err != nil {
		return err
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
