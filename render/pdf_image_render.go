/*
 * Render PDF files to images.
 *
 * Renders all pages of all input files to PNG images, and saves them in the
 * specified output directory.
 *
 * Run as: go run pdf_image_render.go OUTPUT_DIR INPUT.pdf...
 */

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/render"
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
		fmt.Printf("Usage: %s OUTPUT_DIR INPUT.pdf...\n", os.Args[0])
		os.Exit(1)
	}
	outDir := os.Args[1]

	for _, filename := range os.Args[2:] {
		// Create reader.
		reader, f, err := model.NewPdfReaderFromFile(filename, nil)
		if err != nil {
			log.Fatalf("Could not create reader: %v\n", err)
		}
		defer f.Close()

		// Get total number of pages.
		numPages, err := reader.GetNumPages()
		if err != nil {
			log.Fatalf("Could not retrieve number of pages: %v\n", err)
		}

		// Render pages.
		basename := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))

		device := render.NewImageDevice()
		for i := 1; i <= numPages; i++ {
			// Get page.
			page, err := reader.GetPage(i)
			if err != nil {
				log.Fatalf("Could not retrieve page: %v\n", err)
			}

			// Render page to PNG file.
			// RenderToPath chooses the image format by looking at the extension
			// of the specified filename. Only PNG and JPEG files are supported
			// currently.
			outFilename := filepath.Join(outDir, fmt.Sprintf("%s_%d.png", basename, i))
			if err = device.RenderToPath(page, outFilename); err != nil {
				log.Fatalf("Image rendering error: %v\n", err)
			}

			// Alternatively, an image.Image instance can be obtained by using
			// the Render method of the image device, which can then be encoded
			// and saved in any format.
			// image, err := device.Render(page)
		}
	}
}
