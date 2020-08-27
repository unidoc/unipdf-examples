/*
 * Fill PDF form via JSON file data and use custom font replacement.
*
* Run as: go run pdf_form_fill_custom_font.go.
*/

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/fjson"
	"github.com/unidoc/unipdf/v3/model"
)

// Example of filling PDF formdata with a form using custom font.
func main() {
	// Enable debug-level logging.
	// common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))

	inputPdf := "sample_form.pdf"
	fillJSONFile := "formdata.json"
	outputFile := "output.pdf"

	err := fillFields(inputPdf, fillJSONFile, outputFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Success, output written to %s\n", outputFile)
}

// fillFields loads field data from `jsonPath` to fill in
// form data in `inputPath` and outputs as PDF in `outputPath`.
// The output PDF form is flattened.
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

	// set custom font
	fieldAppearance := annotator.FieldAppearance{OnlyIfMissing: true, RegenerateTextFields: true}

	// set font using standard font
	defaultFontReplacement, err := model.NewStandard14Font(model.HelveticaObliqueName)

	// set font using ttf font file
	fontReplacement, err := model.NewPdfFontFromTTFFile("./DoHyeon-Regular.ttf")

	if err != nil {
		log.Fatalf("Error %s", err)
	}

	// replace email field's font using `fontReplacement`
	// and set the other field's font using `defaultFontReplacement`
	style := fieldAppearance.Style()
	style.Fonts = &annotator.AppearanceFontStyle{
		Fallback: &annotator.AppearanceFont{
			Font: defaultFontReplacement,
			Name: defaultFontReplacement.FontDescriptor().FontName.String(),
			Size: 0,
		},
		FieldFallbacks: map[string]*annotator.AppearanceFont{
			"email4": {
				Font: fontReplacement,
				Name: fontReplacement.FontDescriptor().FontName.String(),
				Size: 0,
			},
		},
		ForceReplace: true,
	}

	fieldAppearance.SetStyle(style)

	// Populate the form data.
	err = pdfReader.AcroForm.FillWithAppearance(fdata, fieldAppearance)
	if err != nil {
		return err
	}

	// Flatten form.
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
