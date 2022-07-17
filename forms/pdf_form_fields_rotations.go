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
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/core"
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
		fmt.Printf("Usage: go run pdf_form_fields_rotations.go output.pdf\n")
		os.Exit(1)
	}

	outputPath := os.Args[1]

	pdfFormFieldsPath := "./form_fields_rotations.pdf"
	if err := fillFormFields(pdfFormFieldsPath, outputPath); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

// Fill form fields and write to file.
func fillFormFields(pdfFormPath, outputPath string) error {
	f, err := os.Open(pdfFormPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return err
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
	}

	fieldAppearance := annotator.FieldAppearance{OnlyIfMissing: true, RegenerateTextFields: true}

	// Style MK override true:
	style := fieldAppearance.Style()
	style.AllowMK = true
	fieldAppearance.SetStyle(style)

	// We Extract Fields Data from the fileJson Path.
	fieldsData, err := fjson.LoadFromJSONFile("./form_fields_rotations.json")
	if err != nil {
		return err
	}

	// Fill image.
	if err = fieldsData.SetImageFromFile("image1_af_image", "./images/1.jpg", nil); err != nil {
		return err
	}

	// Populate the form data.
	err = pdfReader.AcroForm.FillWithAppearance(fieldsData, fieldAppearance)
	if err != nil {
		return err
	}

	opt := &model.ReaderToWriterOpts{
		SkipAcroForm: false,
	}

	pdfWriter, err := pdfReader.ToWriter(opt)
	if err != nil {
		return err
	}

	if err := pdfWriter.WriteToFile(outputPath); err != nil {
		return err
	}

	return nil
}
