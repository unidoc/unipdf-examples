/*
 * Create a `Creator` from `model.PdfReader`.
 *
 * Run as: go run pdf_creator_from_reader.go input.pdf output.pdf
 */

package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/core"
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
	if len(os.Args) < 3 {
		fmt.Printf("go run pdf_creator_from_reader.go input.pdf output.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	err := modifyPage(inputPath, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

// Annotate pdf file.
func modifyPage(inputPath string, outputPath string) error {
	watermarkPath := "logo.png"

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

	imgFile, err := os.Open(watermarkPath)
	if err != nil {
		return err
	}
	defer imgFile.Close()

	// Load the image with default handler.
	img, err := model.ImageHandling.Read(imgFile)
	if err != nil {
		return err
	}

	opts := &creator.CreatorFromReaderOpts{
		PageProcessCallback: func(c *creator.Creator, pageNum int, page *model.PdfPage) error {
			// Add text annotation.
			textAnnotation := model.NewPdfAnnotationText()
			textAnnotation.Contents = core.MakeString("Test using Annotation")
			// The rect specifies the location of the markup.
			textAnnotation.Rect = core.MakeArray(core.MakeInteger(20), core.MakeInteger(100), core.MakeInteger(10+50), core.MakeInteger(100+50))

			page.AddAnnotation(textAnnotation.PdfAnnotation)

			watermarkImg, err := c.NewImage(img)
			if err != nil {
				return errors.New(fmt.Sprintf("Error: %v\n", err))
			}

			watermarkImg.ScaleToWidth(c.Context().PageWidth)
			watermarkImg.SetPos(0, (c.Context().PageHeight-watermarkImg.Height())/2)
			watermarkImg.SetOpacity(0.5)

			err = c.Draw(watermarkImg)
			if err != nil {
				return errors.New(fmt.Sprintf("Error: %v\n", err))
			}

			return nil
		},
	}

	pdfCreator, err := creator.NewFromReader(pdfReader, opts)
	if err != nil {
		return err
	}

	err = pdfCreator.WriteToFile(outputPath)
	if err != nil {
		return err
	}

	return nil
}
