/*
 * Add JPG images from given folder to a PDF file, one image per page using lazy mode.
 *
 * Run as: go run pdf_images_to_pdf_lazy.go output.pdf image_folder
 */

package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
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
		fmt.Printf("Usage: go run pdf_images_to_pdf_lazy.go output.pdf image_folder\n")
		os.Exit(1)
	}

	outputPath := os.Args[1]
	inputPath := os.Args[2]

	err := imageFolderToPdf(inputPath, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

// Images to PDF.
func imageFolderToPdf(inputPath string, outputPath string) error {
	c := creator.New()

	files, err := os.ReadDir(inputPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, imgPath := range files {
		if !strings.HasSuffix(imgPath.Name(), ".jpg") {
			continue
		}
		common.Log.Debug("Image: %s", imgPath)

		path := path.Join(inputPath, imgPath.Name())

		img, err := c.NewImageFromFile(path)
		if err != nil {
			common.Log.Debug("Error loading image: %v", err)
			return err
		}
		// Set lazy mode for image allows to reduce memory consumption
		// Lazy mode is helpful for creating PDF with a lot of images
		img.SetLazy(true)

		// Use page width of 612 points, and calculate the height proportionally based on the image.
		// Standard PPI is 72 points per inch, thus a width of 8.5"
		pageWidth := 612.0
		img.ScaleToWidth(pageWidth)

		pageHeight := pageWidth * img.Height() / img.Width()
		c.SetPageSize(creator.PageSize{pageWidth, pageHeight})
		c.NewPage()
		img.SetPos(0, 0)
		_ = c.Draw(img)
	}

	err = c.WriteToFile(outputPath)
	return err
}
