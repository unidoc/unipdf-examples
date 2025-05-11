/*
 * Fill PDF form via JSON input data and flatten the output PDF.
*
* Run as: go run pdf_form_fill_json.go input.pdf fill.json [output.pdf].
*/

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v4/annotator"
	"github.com/unidoc/unipdf/v4/common/license"
	"github.com/unidoc/unipdf/v4/fjson"
	"github.com/unidoc/unipdf/v4/model"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}
}

// Example of filling PDF formdata with a form.
func main() {
	if len(os.Args) < 2 {
		fmt.Printf("List and fill values in PDF form, flatten\n")
		fmt.Printf("Usage: go run pdf_form_fill_json.go input.pdf fill.json [output.pdf]\n\n")
		fmt.Printf("To get a list of fields and values from a PDF file as JSON:\n")
		fmt.Printf("  go run pdf_form_fill_json.go input.pdf > formdata.json\n\n")
		fmt.Printf("To fill a PDF with form data from a JSON file:\n")
		fmt.Printf("  go run pdf_form_fill_json.go input.pdf formdata.json output.pdf\n")
		os.Exit(1)
	}

	var (
		inputPath    string
		filljsonPath string
		outputPath   string
	)
	inputPath = os.Args[1]
	if len(os.Args) > 3 {
		filljsonPath = os.Args[2]
		outputPath = os.Args[3]
	}

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

	err := fillFields(inputPath, filljsonPath, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Success, output written to %s\n", outputPath)
}

// fillFields loads field data from `jsonPath` and used to fill in form data in `inputPath` and outputs
// as PDF in `outputPath`. The output PDF form is flattened.
func fillFields(inputPath, jsonPath, outputPath string) error {
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

	// Populate the form data.
	err = pdfReader.AcroForm.Fill(fdata)
	if err != nil {
		return err
	}

	// Flatten form.
	fieldAppearance := annotator.FieldAppearance{OnlyIfMissing: true, RegenerateTextFields: true}

	// NOTE: To customize certain styles try:
	// style := fieldAppearance.Style()
	// style.CheckmarkGlyph = "a22"
	// style.AutoFontSizeFraction = 0.70
	// fieldAppearance.SetStyle(style)
	//
	// or for specifying a full set of appearance styles:
	// fieldAppearance.SetStyle(annotator.AppearanceStyle{
	//     CheckmarkGlyph:       "a22",
	//     AutoFontSizeFraction: 0.70,
	//     FillColor:            model.NewPdfColorDeviceGray(0.8),
	//     BorderColor:          model.NewPdfColorDeviceRGB(1, 0, 0),
	//     BorderSize:           2.0,
	//     AllowMK:              false,
	// })

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
