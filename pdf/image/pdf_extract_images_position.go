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
 * Run as: go run pdf_extract_native_images.go input.pdf output.folder
 */

package main

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

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
		fmt.Printf("Syntax: go run pdf_extract_images.go input.pdf output.folder\n")
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
func extractImagesToArchive(inputPath, outputDir string) error {
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

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.Mkdir(outputDir, 0777)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not create outputDir=%q err=%v\n", outputDir, err)
			return err
		}
	}

	for pageNum := 1; pageNum <= numPages; pageNum++ {

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return err
		}

		// List images on the page.
		images, report, err := extractImagesOnPage(page)
		if err != nil {
			return err
		}

		if report != nil {
			fmt.Printf("-----\nPage %d:\n", pageNum)
			fmt.Printf("%s\n", strings.Join(report, ""))
		}

		for idx, imgData := range images {
			img := imgData.img
			fname := fmt.Sprintf("page%d_img%d", pageNum, idx+1)
			if !imgData.inline {
				fname += "_xobj"
			}

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
			outpuPath := filepath.Join(outputDir, fname)

			fmt.Printf("  Converting to go image: page %d img %d ⇾ %q %s\n",
				pageNum, idx+1, outpuPath, describe(img))

			imgf, err := os.OpenFile(outpuPath, os.O_RDWR|os.O_CREATE, 0755)
			if err != nil {
				return err
			}
			defer imgf.Close()

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
	inline bool
}

func count(data []imageData, inline bool) int {
	n := 0
	for _, d := range data {
		if d.inline == inline {
			n++
		}
	}
	return n
}

// extractImagesOnPage returns a slice of all images on page `page`.
func extractImagesOnPage(page *pdf.PdfPage) ([]imageData, []string, error) {
	contents, err := page.GetAllContentStreams()
	if err != nil {
		return nil, nil, err
	}
	return extractImagesInContentStream(contents, page.Resources)
}

// extractImagesInContentStream returns a slice of all images in content stream `contents`.
func extractImagesInContentStream(contents string, resources *pdf.PdfPageResources) ([]imageData, []string, error) {
	var images []imageData
	cstreamParser := pdfcontent.NewContentStreamParser(contents)
	operations, err := cstreamParser.Parse()
	if err != nil {
		return nil, nil, err
	}

	processedXObjects := map[string]bool{}
	var report []string

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
				return nil, report, err
			}

			cs, err := iimg.GetColorSpace(resources)
			if err != nil {
				return nil, report, err
			}
			if cs == nil {
				// Default if not specified?
				cs = pdf.NewPdfColorspaceDeviceGray()
			}
			images = append(images, imageData{img, nil, true})
			report = append(report, fmt.Sprintf(" Inline Cs: %T (%d inline + %d xobjimages)\n",
				cs, count(images, true), count(images, false)))
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
				report = append(report, fmt.Sprintf(" XObject Image: %s\n", *name))

				ximg, err := resources.GetXObjectImageByName(*name)
				if err != nil {
					return nil, report, err
				}

				report = append(report, fmt.Sprintf(" ximg=%dx%d ", *ximg.Width, *ximg.Height))
				report = append(report, fmt.Sprintf("ximg.Filter=%T=%s\n",
					ximg.Filter, ximg.Filter.GetFilterName()))

				img, err := ximg.ToImage()
				if err != nil {
					return nil, report, err
				}
				images = append(images, imageData{img, ximg.Filter, false})
				xObjectImages++
			} else if xtype == pdf.XObjectTypeForm {
				// Go through the XObject Form content stream.
				xform, err := resources.GetXObjectFormByName(*name)
				if err != nil {
					return nil, report, err
				}

				formContent, err := xform.GetContentStream()
				if err != nil {
					return nil, report, err
				}

				// Process the content stream in the Form object too:
				formResources := xform.Resources
				if formResources == nil {
					formResources = resources
				}

				// Process the content stream in the Form object too:
				formImages, formReport, err := extractImagesInContentStream(string(formContent), formResources)
				if err != nil {
					return nil, report, err
				}
				images = append(images, formImages...)
				report = append(report, formReport...)
			}
		}
	}

	return images, report, nil
}

func describe(img *pdf.Image) string {
	desc := fmt.Sprintf("%dx%d cpts=%d bpp=%d",
		img.Width, img.Height, img.ColorComponents, img.BitsPerComponent)
	if len(desc) > 100 {
		panic("ddd")
	}
	return desc
}
