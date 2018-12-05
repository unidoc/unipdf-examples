/*
 * Extract images from a PDF file. Passes through each page, goes through the content stream and
 * finds instances of both XObject Images and inline images. Also handles images referred within
 * XObject Form content streams.
 * The output files are saved as a zip archive.
 *
 * NOTE(peterwilliams97): Unlike pdf_extract_images.go,
 *       1) Losslessly compressed PDF images are saved in PNG format. (Lossily compressed PDF images
 *          still are saved in JPEG format)
 *       2) Images are saved in the same color space as they occur in PDF files.
 *
 * XXX(peterwilliams97): This file fixes an apparent problem in the UniDoc resampling code with
 *       handling 1 bit per component images. There is an additional problem that got6.DecodeBytes()
 *       returns CCITTFax images as 8 bits per pixel while PDF expects these images to be 1 bit per
 *       pixel. I tried modifying got.6 to fix this and found that ResampleBytes() didn't work with
 *       1 bit per pixel, so I just set imgData.img.BitsPerComponent = 8 for CCITTFaxEncoder images.
 *
 * TODO: Handle JBIG images.
 *       Handle CCITTFaxEncoder inline images?
 *       Handle CCITTFaxEncoder Group 3 images?
 *       Change got.6 to return 1 bit images.
 *       Save images in orientation they appear in the PDF file.
 *
 * Run as: go run pdf_extract_native_images.go input.pdf output.zip
 */

package main

import (
	"archive/zip"
	"fmt"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/unidoc/unidoc/common"
	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

var xObjectImages = 0
var inlineImages = 0

func main() {
	// Enable debug-level console logging, when debuggingn:
	//unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))
	common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))

	if len(os.Args) < 3 {
		fmt.Printf("Syntax: go run pdf_extract_images.go input.pdf output.zip\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	fmt.Printf("Input file: %s\n", inputPath)
	err := extractImagesToArchive(inputPath, outputPath)
	if err != nil {
		fmt.Printf("ERROR: Could not process inputPath=%q outputPath=%q err=%v\n",
			inputPath, outputPath, err)
		os.Exit(1)
	}

	fmt.Printf("-- Summary\n")
	fmt.Printf("%d XObject images extracted\n", xObjectImages)
	fmt.Printf("%d inline images extracted\n", inlineImages)
	fmt.Printf("Total %d images\n", xObjectImages+inlineImages)
}

// extractImagesToArchive extracts images and properties of a PDF specified by `inputPath`.
// The output images are stored into a zip archive whose path is given by `outputPath`.
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
	defer func() {
		// Make sure to check the error on Close.
		err2 := zipw.Close()
		if err == nil {
			err = err2
		}
	}()

	for i := 0; i < numPages; i++ {
		fmt.Printf("-----\nPage %d:\n", i+1)

		page, err := pdfReader.GetPage(i + 1)
		if err != nil {
			return err
		}

		// List images on the page.
		images, err := extractImagesOnPage(page)
		if err != nil {
			return err
		}

		for idx, imgData := range images {
			img := imgData.img
			fname := fmt.Sprintf("p%d_%d", i+1, idx)

			var lossless bool // Is image compressed losslessly?

			switch imgData.filter.(type) {
			case *pdfcore.FlateEncoder:
				lossless = true
			case *pdfcore.CCITTFaxEncoder:
				lossless = true
				// XXX(peterwilliams97) Hack to work around got6.DecodeBytes() returning an 8 bits
				// per component raster and sampling.ResampleBytes() not working for 1 bits per
				// pixel
				imgData.img.BitsPerComponent = 8
			}

			gimg, err := img.ToGoImage()
			if err != nil {
				return err
			}

			if lossless {
				fname += ".png"
			} else {
				fname += ".jpg"
			}

			fmt.Printf("Converting to go image: page %d img %d -> %q %s\n", i+1, idx+1, fname, img)

			imgf, err := zipw.Create(fname)
			if err != nil {
				return err
			}

			if lossless {
				err = png.Encode(imgf, gimg)
			} else {
				opt := jpeg.Options{Quality: 100}
				err = jpeg.Encode(imgf, gimg, &opt)
			}
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type imageData struct {
	img    *pdf.Image
	filter pdfcore.StreamEncoder
}

// extractImagesOnPage returns a slice of all images on page `page`.
func extractImagesOnPage(page *pdf.PdfPage) ([]imageData, error) {
	contents, err := page.GetAllContentStreams()
	if err != nil {
		return nil, err
	}
	return extractImagesInContentStream(contents, page.Resources)
}

// extractImagesInContentStream returns a slice of all images in content stream `contents`.
func extractImagesInContentStream(contents string, resources *pdf.PdfPageResources) ([]imageData, error) {
	images := []imageData{}
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

			images = append(images, imageData{img: img})
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

				fmt.Printf("ximg=%s\n", ximg)
				fmt.Printf("ximg.Filter=%T=%s\n", ximg.Filter, ximg.Filter.GetFilterName())

				img, err := ximg.ToImage()
				if err != nil {
					return nil, err
				}
				images = append(images, imageData{img, ximg.Filter})
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
				formImages, err := extractImagesInContentStream(string(formContent), formResources)
				if err != nil {
					return nil, err
				}
				images = append(images, formImages...)
			}
		}
	}

	return images, nil
}
