/*
 * Draw a line in a new PDF file.
 *
 * Run as: go run pdf_draw_line.go <x1> <y1> <x2> <y2> output.pdf
 * The x, y coordinates start from the upper left corner at (0,0) and increase going right, down respectively.
 */

package main

import (
	"fmt"
	"os"
	"strconv"

	unicommon "github.com/unidoc/unidoc/common"
	"github.com/unidoc/unidoc/pdf/creator"
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
	// New creator with default properties (pagesize letter default).
	c := creator.New()

	c.NewPage()
	line := creator.NewLine(x1, y1, x2, y2)
	line.SetLineWidth(1.5)
	r, g, b, _ := creator.ColorRGBFromHex("#ff0000")
	line.SetColorRGB(r, g, b)
	c.Draw(line)

	err := c.WriteToFile(outputPath)
	return err
}
