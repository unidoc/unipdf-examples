/*
 * Annotate/mark up pages of a PDF file.
 * Add an ink annotation with shape of a checkmark on a page.
 *
 * Run as: go run pdf_annotate_add_ink.go input.pdf output.pdf
 */

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/contentstream/draw"
	"github.com/unidoc/unipdf/v3/core"
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
	if len(os.Args) < 3 {
		fmt.Printf("go run pdf_annotate_add_ink.go input.pdf output.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	paths := []draw.Path{
		{
			Points: []draw.Point{
				draw.NewPoint(361.019, 638.415),
				draw.NewPoint(361.488, 637.946),
				draw.NewPoint(363.361, 635.605),
				draw.NewPoint(365.233, 632.796),
				draw.NewPoint(368.043, 629.518),
				draw.NewPoint(371.321, 625.772),
				draw.NewPoint(373.662, 622.494),
				draw.NewPoint(375.535, 620.621),
				draw.NewPoint(376.471, 619.216),
				draw.NewPoint(376.94, 618.28),
				draw.NewPoint(377.408, 618.28),
				draw.NewPoint(377.876, 618.748),
				draw.NewPoint(379.281, 621.089),
				draw.NewPoint(382.09, 627.177),
				draw.NewPoint(386.305, 636.073),
				draw.NewPoint(390.987, 646.843),
				draw.NewPoint(395.201, 655.272),
				draw.NewPoint(397.543, 660.891),
				draw.NewPoint(398.947, 664.168),
				draw.NewPoint(399.884, 666.041),
				draw.NewPoint(400.352, 666.51),
			},
		},
	}

	err := annotatePdfAddInk(inputPath, outputPath, paths)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

// Annotate pdf file.
func annotatePdfAddInk(inputPath string, outputPath string, paths []draw.Path) error {
	// Read the input pdf file.
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return err
	}

	// Process each page using the following callback
	// when generating PdfWriter.
	opt := &model.ReaderToWriterOpts{
		PageProcessCallback: func(pageNum int, page *model.PdfPage) error {
			inkAnnot, err := annotator.CreateInkAnnotation(annotator.InkAnnotationDef{
				Paths: paths,
				Color: model.NewPdfColorDeviceRGB(0.89, 0.13, 0.21),
			})
			if err != nil {
				return err
			}

			switch t := inkAnnot.GetContext().(type) {
			case *model.PdfAnnotationInk:
				t.T = core.MakeName("InkAnnot")

				creationDate, err := model.NewPdfDateFromTime(time.Now())
				if err != nil {
					return err
				}

				t.CreationDate = creationDate.ToPdfObject()
				t.M = creationDate.ToPdfObject()
				t.P = page.ToPdfObject()
			}

			// Add to the page annotations.
			page.AddAnnotation(inkAnnot)

			return nil
		},
	}

	// Generate a PdfWriter instance from existing PdfReader.
	pdfWriter, err := pdfReader.ToWriter(opt)
	if err != nil {
		return err
	}

	// Write to file.
	err = pdfWriter.WriteToFile(outputPath)
	if err != nil {
		return err
	}

	return nil
}
