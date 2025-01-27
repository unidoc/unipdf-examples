/*
 * An example of adding text watermark to pdf pages.
 *
 * Run as: go run pdf_add_text_watermark.go <input.pdf> <watermark text> <output.pdf>
 */
package main

import (
	"fmt"
	"image/color"
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
		fmt.Printf("Usage: go run pdf_add_text_watermark.go input.pdf watermark output.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	watermark := os.Args[2]
	outputPath := os.Args[3]

	file, err := os.Open(inputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	pdfReader, err := model.NewPdfReader(file)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	c := creator.New()

	totalPageNumb, err := pdfReader.GetNumPages()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	for i := 1; i <= totalPageNumb; i++ {
		page, err := pdfReader.GetPage(i)
		if err != nil {
			fmt.Printf("Error: failed to get page %d. %v\n", i, err)
			os.Exit(1)
		}
		options := model.WatermarkTextOptions{
			Alpha:     0.3,
			FontSize:  40,
			FontPath:  "Roboto-Regular.ttf",
			FontColor: color.RGBA{R: 255, G: 0, B: 0, A: 1},
			Angle:     30,
		}

		err = page.AddWatermarkText(watermark, options)
		if err != nil {
			fmt.Printf("Error: failed to add watermark on page %d. %v\n", i, err)
		}

		c.AddPage(page)
	}

	c.WriteToFile(outputPath)
	fmt.Print("Watermark added successfully.\n")
}
