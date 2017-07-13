/*
 * Extract images from a PDF file. Passes through each page, goes through the content stream and finds instances of both
 * XObject Images and inline images. Also handles images referred within XObject Form content streams.
 * The output files are saved as a zip archive.
 *
 * Run as: go run pdf_extract_images.go input.pdf output.zip
 */

package main

import (
	"archive/zip"
	"fmt"
	"image/jpeg"
	"os"

	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

var xObjectImages = 0
var inlineImages = 0

func main() {
	// Enable debug-level console logging, when debuggingn:
	//unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	if len(os.Args) < 3 {
		fmt.Printf("Syntax: go run pdf_extract_images.go input.pdf output.zip\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	fmt.Printf("Input file: %s\n", inputPath)
	err := extractImagesToArchive(inputPath, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("-- Summary\n")
	fmt.Printf("%d XObject images extracted\n", xObjectImages)
	fmt.Printf("%d inline images extracted\n", inlineImages)
	fmt.Printf("Total %d images\n", xObjectImages+inlineImages)
}

// Extracts images and properties of a PDF specified by inputPath.
// The output images are stored into a zip archive whose path is given by outputPath.
func extractImagesToArchive(inputPath, outputPath string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return err
	}

	// Try decrypting with an empty one.
	if isEncrypted {
		auth, err := pdfReader.Decrypt([]byte(""))
		if err != nil {
			// Encrypted and we cannot do anything about it.
			return err
		}
		if !auth {
			fmt.Println("Need to decrypt with password")
			return nil
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}
	fmt.Printf("PDF Num Pages: %d\n", numPages)

	// Prepare output archive.
	zipf, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer zipf.Close()
	zipw := zip.NewWriter(zipf)

	for i := 0; i < numPages; i++ {
		fmt.Printf("-----\nPage %d:\n", i+1)

		page, err := pdfReader.GetPage(i + 1)
		if err != nil {
			return err
		}

		// List images on the page.
		rgbImages, err := extractImagesOnPage(page)
		if err != nil {
			return err
		}
		_ = rgbImages

		for idx, img := range rgbImages {
			fname := fmt.Sprintf("p%d_%d.jpg", i+1, idx)

			gimg, err := img.ToGoImage()
			if err != nil {
				return err
			}

			imgf, err := zipw.Create(fname)
			if err != nil {
				return err
			}
			opt := jpeg.Options{Quality: 100}
			err = jpeg.Encode(imgf, gimg, &opt)
			if err != nil {
				return err
			}
		}
	}

	// Make sure to check the error on Close.
	err = zipw.Close()
	if err != nil {
		return err
	}

	return nil
}

func extractImagesOnPage(page *pdf.PdfPage) ([]*pdf.Image, error) {
	contents, err := page.GetAllContentStreams()
	if err != nil {
		return nil, err
	}

	return extractImagesInContentStream(contents, page.Resources)
}

func extractImagesInContentStream(contents string, resources *pdf.PdfPageResources) ([]*pdf.Image, error) {
	rgbImages := []*pdf.Image{}
	cstreamParser := pdfcontent.NewContentStreamParser(contents)
	operations, err := cstreamParser.Parse()
	if err != nil {
		return nil, err
	}

	processedXObjects := map[string]bool{}

	// Range through all the content stream operations.
	for _, op := range *operations {
		if op.Operand == "BI" && len(op.Params) == 1 {
			// BI: Inline image.

			iimg, ok := op.Params[0].(*pdfcontent.ContentStreamInlineImage)
			if !ok {
				continue
			}

			img, err := iimg.ToImage(resources)
			if err != nil {
				return nil, err
			}

			cs, err := iimg.GetColorSpace(resources)
			if err != nil {
				return nil, err
			}
			if cs == nil {
				// Default if not specified?
				cs = pdf.NewPdfColorspaceDeviceGray()
			}
			fmt.Printf("Cs: %T\n", cs)

			rgbImg, err := cs.ImageToRGB(*img)
			if err != nil {
				return nil, err
			}

			rgbImages = append(rgbImages, &rgbImg)
			inlineImages++
		} else if op.Operand == "Do" && len(op.Params) == 1 {
			// Do: XObject.
			name := op.Params[0].(*pdfcore.PdfObjectName)

			// Only process each one once.
			_, has := processedXObjects[string(*name)]
			if has {
				continue
			}
			processedXObjects[string(*name)] = true

			_, xtype := resources.GetXObjectByName(*name)
			if xtype == pdf.XObjectTypeImage {
				fmt.Printf(" XObject Image: %s\n", *name)

				ximg, err := resources.GetXObjectImageByName(*name)
				if err != nil {
					return nil, err
				}

				img, err := ximg.ToImage()
				if err != nil {
					return nil, err
				}

				rgbImg, err := ximg.ColorSpace.ImageToRGB(*img)
				if err != nil {
					return nil, err
				}
				rgbImages = append(rgbImages, &rgbImg)
				xObjectImages++
			} else if xtype == pdf.XObjectTypeForm {
				// Go through the XObject Form content stream.
				xform, err := resources.GetXObjectFormByName(*name)
				if err != nil {
					return nil, err
				}

				formContent, err := xform.GetContentStream()
				if err != nil {
					return nil, err
				}

				// Process the content stream in the Form object too:
				formResources := xform.Resources
				if formResources == nil {
					formResources = resources
				}

				// Process the content stream in the Form object too:
				formRgbImages, err := extractImagesInContentStream(string(formContent), formResources)
				if err != nil {
					return nil, err
				}
				rgbImages = append(rgbImages, formRgbImages...)
			}
		}
	}

	return rgbImages, nil
}
