/*
* Fill PDF form via json data where field names and values along side text colors for each field are provided.
*
* Run as: go run form-fill-custom-text-color.go.
 */
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/common/license"
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

// FieldData represents the field data with name and text color attributes.
type FieldData struct {
	Name      string `json:"name"`
	TextColor string `json:"text_color"`
}

func main() {
	inputPath := "sample_form.pdf"
	jsonDataPath := "sample_form_rich.json"
	outputPath := "filled_form_rich.pdf"

	fdata, err := fjson.LoadFromJSONFile(jsonDataPath)
	if err != nil {
		fmt.Printf("failed to load json file. Error : %v\n", err)
		os.Exit(1)
	}

	f, err := os.Open(inputPath)
	if err != nil {
		fmt.Printf("failed to open input file. Error : %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		fmt.Printf("failed to read pdf file. Error : %v\n", err)
		os.Exit(1)
	}

	fieldColor, err := loadFieldColors(jsonDataPath)
	if err != nil {
		fmt.Printf("failed to load field colors. Error : %v\n", err)
		os.Exit(1)
	}

	fieldAppearance := annotator.FieldAppearance{OnlyIfMissing: true, RegenerateTextFields: true}

	fieldAppearance.SetStyle(annotator.AppearanceStyle{
		AutoFontSizeFraction: 0.70,
		FillColor:            model.NewPdfColorDeviceRGB(1, 1, 1),
		BorderColor:          model.NewPdfColorDeviceRGB(0, 0, 0),
		BorderSize:           2.0,
		AllowMK:              false,
		TextColor:            model.NewPdfColorDeviceRGB(0.5, 0.8, 0.8), // Default text color.
		FieldColors:          fieldColor,                                // This specifies the text color for each field.
	})

	pdfReader.AcroForm.FillWithAppearance(fdata, fieldAppearance)
	pdfWriter, err := pdfReader.ToWriter(nil)
	if err != nil {
		fmt.Printf("failed to convert reader to writer pdf file. Error : %v\n", err)
		os.Exit(1)
	}

	// Write to file.
	err = pdfWriter.WriteToFile(outputPath)
	if err != nil {
		fmt.Printf("failed to write to output file. Error : %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Success, output written to %s\n", outputPath)
}

// loadFieldColors loads field colors from json file and returns a map of field name to PdfColor.
func loadFieldColors(jsonPath string) (map[string]model.PdfColor, error) {
	file, err := os.Open(jsonPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var fields []FieldData
	err = json.Unmarshal(byteValue, &fields)
	if err != nil {
		return nil, err
	}

	fieldColors := make(map[string]model.PdfColor)
	for _, field := range fields {
		if field.TextColor == "" {
			continue
		}
		tc := creator.ColorRGBFromHex(field.TextColor)
		r, g, b := tc.ToRGB()
		tcRGB := model.NewPdfColorDeviceRGB(r, g, b)
		fieldColors[field.Name] = tcRGB
	}

	return fieldColors, nil
}
