/*
 * Add images to a PDF file, one image per page.
 *
 * Run as: go run pdf_images_to_pdf.go output.pdf img1.jpg img2.jpg img3.png ...
 */

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/creator"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: go run pdf_images_to_pdf.go output.pdf img1.jpg img2.jpg ...\n")
		os.Exit(1)
	}

	outputPath := os.Args[1]
	inputPaths := os.Args[2:len(os.Args)]

	err := imagesToPdf(inputPaths, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

// Images to PDF.
func imagesToPdf(inputPaths []string, outputPath string) error {
	c := creator.New()

	for _, imgPath := range inputPaths {
		common.Log.Debug("Image: %s", imgPath)

		img, err := c.NewImageFromFile(imgPath)
		if err != nil {
			common.Log.Debug("Error loading image: %v", err)
			return err
		}
		img.ScaleToWidth(612.0)

		// Use page width of 612 points, and calculate the height proportionally based on the image.
		// Standard PPI is 72 points per inch, thus a width of 8.5"
		height := 612.0 * img.Height() / img.Width()
		c.SetPageSize(creator.PageSize{612, height})
		c.NewPage()
		img.SetPos(0, 0)
		_ = c.Draw(img)
	}

	err := c.WriteToFile(outputPath)
	return err
}
