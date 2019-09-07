/*
 * Fill PDF form via JSON input data and flatten the output PDF.
*
* Run as: go run pdf_form_fill_json.go input.pdf fill.json [output.pdf].
*/

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/fjson"
	"github.com/unidoc/unipdf/v3/model"
)

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

	// Enable debug-level logging.
	common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))

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

	// Write out.
	pdfWriter := model.NewPdfWriter()
	pdfWriter.SetForms(nil)

	for _, p := range pdfReader.PageList {
		err := pdfWriter.AddPage(p)
		if err != nil {
			return err
		}
	}

	fout, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer fout.Close()

	err = pdfWriter.Write(fout)
	return err
}
