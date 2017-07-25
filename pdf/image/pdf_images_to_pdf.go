/*
 * Add images to a PDF file, one image per page.
 *
 * Run as: go run pdf_images_to_pdf.go output.pdf img1.jpg img2.jpg img3.png ...
 */

package main

import (
	"fmt"
	"os"

	unicommon "github.com/unidoc/unidoc/common"
	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: go run pdf_add_images.go output.pdf img1.jpg img2.jpg ...\n")
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
	pdfWriter := pdf.NewPdfWriter()

	unicommon.Log.Debug("Inputs: %v", inputPaths)

	// Make the document structure.
	for idx, imgPath := range inputPaths {
		unicommon.Log.Debug("Image: %s", imgPath)
		// Open the image file.
		reader, err := os.Open(imgPath)
		if err != nil {
			unicommon.Log.Error("Error opening file: %s", err)
			return err
		}
		defer reader.Close()

		img, err := pdf.ImageHandling.Read(reader)
		if err != nil {
			unicommon.Log.Error("Error loading image: %s", err)
			return err
		}

		// Use page width of 612 points, and calculate the height proportionally based on the image.
		// Standard PPI is 72 points per inch.
		height := 612 * float64(img.Height) / float64(img.Width)

		// Make a page.
		page := pdf.NewPdfPage()
		bbox := pdf.PdfRectangle{0, 0, 612, height}
		page.MediaBox = &bbox

		// Name for the image object.
		imgName := pdfcore.PdfObjectName(fmt.Sprintf("Im%d", idx+1))

		// Use flate decoding for the images.
		encoder := pdfcore.NewFlateEncoder()

		// Create an XObject Image for the PDF.
		ximg, err := pdf.NewXObjectImageFromImage(img, nil, encoder)
		if err != nil {
			unicommon.Log.Error("Failed to create xobject image: %s", err)
			return err
		}

		// Add to the page resources.
		err = page.AddImageResource(imgName, ximg)
		if err != nil {
			unicommon.Log.Error("Failed to create xobject image: %s", err)
			return err
		}

		// Create a normal graphics state.
		gsName := pdfcore.PdfObjectName(fmt.Sprintf("GS0"))
		gs0 := pdfcore.PdfObjectDictionary{}
		gs0[pdfcore.PdfObjectName("BM")] = pdfcore.MakeName("Normal")
		page.AddExtGState(gsName, &gs0)

		// Content stream to load the image.
		creator := pdfcontent.NewContentCreator()
		creator.
			Add_q().
			Add_gs(gsName).
			Add_cm(612, 0, 0, height, 0, 0).
			Add_Do(imgName).
			Add_Q()
		page.AddContentStreamByString(creator.String())

		err = pdfWriter.AddPage(page)
		if err != nil {
			unicommon.Log.Error("Failed to add page: %s", err)
			return err
		}
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
