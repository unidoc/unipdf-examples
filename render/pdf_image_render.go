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

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/render"
)

func main() {
	common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))

	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s OUTPUT_DIR INPUT.pdf...\n", os.Args[0])
		os.Exit(1)
	}
	outDir := os.Args[1]

	for _, filename := range os.Args[2:] {
		// Create reader.
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("Could not open input file: %v\n", err)
		}
		defer file.Close()

		reader, err := model.NewPdfReader(file)
		if err != nil {
			log.Fatalf("Could not create reader: %v\n", err)
		}

		// Check if file is encrypted.
		isEncrypted, err := reader.IsEncrypted()
		if err != nil {
			log.Fatalf("Could not read file info: %v\n", err)
		}

		// Attempt to decrypt using an empty password.
		if isEncrypted {
			auth, err := reader.Decrypt([]byte(""))
			if err != nil {
				log.Fatalf("Could not decrypt input file: %v\n", err)
			}
			if !auth {
				log.Fatalf("Could not decrypt input file")
			}
		}

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
