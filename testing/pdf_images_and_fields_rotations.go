/*
 * Add images to a PDF file, one image per page.
 *
 * Run as: go run pdf_images_and_fields_rotations.go output.pdf
 */

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/fjson"
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

type imagePathAndAngle struct {
	FilePath string
	Rotation float64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run pdf_images_and_fields_rotations.go output.pdf ...\n")
		os.Exit(1)
	}

	outputPath := os.Args[1]

	imagePaths := []imagePathAndAngle{
		imagePathAndAngle{
			FilePath: "./images/1.jpg",
			Rotation: 90,
		},
		imagePathAndAngle{
			FilePath: "./images/2.jpg",
			Rotation: 180,
		},
		imagePathAndAngle{
			FilePath: "./images/3.jpg",
			Rotation: 270,
		},
		imagePathAndAngle{
			FilePath: "./images/4.jpg",
			Rotation: 0,
		},
	}

	c := creator.New()

	err := imagesToCreator(c, imagePaths)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	pdfFormFieldsPath := "./form_fields_test.pdf"

	if err := c.WriteToFile(outputPath); err != nil {
		panic(err)
	}

	c = creator.New()

	pdfWriter, err := addFormFieldsToPdfWriter(pdfFormFieldsPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	filePdf, err := os.Open(outputPath)
	if err != nil {
		panic(err)
	}
	defer filePdf.Close()

	pdfReader, err := model.NewPdfReader(filePdf)
	if err != nil {
		panic(err)
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		panic(err)
	}

	for i := 1; i <= numPages; i++ {
		page, err := pdfReader.GetPage(i)
		if err != nil {
			panic(err)
		}
		pdfWriter.AddPage(page)
	}

	if err := pdfWriter.WriteToFile(outputPath); err != nil {
		panic(err)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

// Images to creator.
func imagesToCreator(c *creator.Creator, imagePaths []imagePathAndAngle) error {
	for _, imgPath := range imagePaths {
		img, err := c.NewImageFromFile(imgPath.FilePath)
		if err != nil {
			common.Log.Debug("Error loading image: %v", err)
			return err
		}
		img.ScaleToWidth(100)
		// Set image angle rotation.
		img.SetAngle(imgPath.Rotation)
		img.SetMargins(10, 10, 10, 10)

		// Optionally, set an encoder for the image. If none is specified, the
		// encoder defaults to core.FlateEncoder, which applies lossless compression
		// to the image stream. However, core.FlateEncoder tends to produce large
		// image streams which results in large output file sizes.
		// However, the encoder can be changed to core.DCTEncoder, which applies
		// lossy compression (this type of compression is used by JPEG images) in
		// order to reduce the output file size.
		encoder := core.NewDCTEncoder()
		// The default quality is 75. There is not much difference in the image
		// quality between 75 and 100 but the size difference when compressing the
		// image stream is signficant.
		// encoder.Quality = 100
		img.SetEncoder(encoder)

		if err := c.Draw(img); err != nil {
			panic(err)
		}
	}
	return nil
}

// Fields to creator.

func addFormFieldsToPdfWriter(pdfFormPath string) (*model.PdfWriter, error) {
	f, err := os.Open(pdfFormPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return nil, err
	}
	acroForm := pdfReader.AcroForm
	fields := acroForm.AllFields()

	// Get and modify fields appearance.
	for _, field := range fields {
		fmt.Printf("field: %v\n", field.T)

		angle := int64(0)
		ctx := field.GetContext()
		switch t := ctx.(type) {
		case *model.PdfFieldButton:
			angle = 90
			if t.IsPush() {
				angle = 270
			}
		case *model.PdfFieldText:
			angle = 180
		case *model.PdfFieldChoice:
			angle = 90
		default:
			fmt.Printf(" Unknown Field Type\n")
			continue
		}

		fmt.Printf(" Annotations: %d\n", len(field.Annotations))
		for j, wa := range field.Annotations {
			// Get MK dictionary.
			if mkDict, has := core.GetDict(wa.MK); has {
				// R object for rotation value of field appearance.
				rotateName := core.MakeName("R")
				rotateVal := core.MakeInteger(angle)
				mkDict.Set(*rotateName, rotateVal)
				wa.MK = mkDict
			}
			field.Annotations[j] = wa
		}

		fieldAppearance := annotator.FieldAppearance{OnlyIfMissing: false, RegenerateTextFields: true}

		// Style MK override true:
		style := fieldAppearance.Style()
		style.AllowMK = true
		fieldAppearance.SetStyle(style)

		// We Extract Fields Data from the fileJson Path.
		fieldsData, err := fjson.LoadFromJSONFile("./form_fields_test.json")
		if err != nil {
			return nil, err
		}

		// Populate the form data.
		err = pdfReader.AcroForm.FillWithAppearance(fieldsData, fieldAppearance)
		if err != nil {
			return nil, err
		}
	}

	opt := &model.ReaderToWriterOpts{
		SkipAcroForm: false,
	}

	return pdfReader.ToWriter(opt)
}
