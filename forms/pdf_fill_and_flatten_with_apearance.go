/*
 * Fill PDF form JSON input data and flatten it with a customized appearance to output PDF.
*
* Run as: go run pdf_fill_and_flatten_with_apearance.go.
*/

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/common/license"
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

func main() {
	inputPath := "./sample_form.pdf"
	jsonDataPath := "./sample_form.json"
	outputPath := "./sample_form_output.pdf"

	// Output path not specified: Export list of fields and data as JSON format.
	if len(outputPath) == 0 {
		fdata, err := fjson.LoadFromPDFFile(inputPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		if fdata == nil {
			fmt.Printf("No data\n")
			return
		}
		fjson, err := fdata.JSON()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", fjson)
		return
	}

	err := fillFieldsWithAppearance(inputPath, jsonDataPath, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Success, output written to %s\n", outputPath)
}

// fillFieldsWithAppearance loads field data from `jsonPath` and to fills the data to form fields in a PDF data provided via `inputPath` and outputs
// as a flattened PDF in `outputPath`. Customized field appearance can be given either during filling or during flattening via the
// `annotator.AppearanceStyle` object.
func fillFieldsWithAppearance(inputPath, jsonPath, outputPath string) error {
	fdata, err := fjson.LoadFromJSONFile(jsonPath)
	if err != nil {
		return err
	}

	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return err
	}

	fieldAppearance := annotator.FieldAppearance{OnlyIfMissing: true, RegenerateTextFields: true}

	// specify a full set of appearance styles
	fieldAppearance.SetStyle(annotator.AppearanceStyle{
		AutoFontSizeFraction: 0.70,
		FillColor:            model.NewPdfColorDeviceRGB(1, 1, 1),
		BorderColor:          model.NewPdfColorDeviceRGB(0, 0, 0),
		BorderSize:           2.0,
		AllowMK:              false,
		TextColor:            model.NewPdfColorDeviceRGB(0.5, 0.8, 0.8),
	})

	// Populate the form data.
	// pdfReader.AcroForm.FillWithAppearance(fdata, fieldAppearance) // uncomment this line to fill with appearance
	pdfReader.AcroForm.Fill(fdata)

	// Flatten form. with field appearance
	err = pdfReader.FlattenFields(true, fieldAppearance)
	if err != nil {
		return err
	}

	// Generate a PdfWriter instance from existing PdfReader.
	pdfWriter, err := pdfReader.ToWriter(nil)
	if err != nil {
		return err
	}

	// Write to file.
	err = pdfWriter.WriteToFile(outputPath)
	return err
}

func getFont(path string) (*model.PdfFont, error) {
	font, err := model.NewCompositePdfFontFromTTFFile(path)
	if err != nil {
		return nil, err
	}
	return font, nil
}
