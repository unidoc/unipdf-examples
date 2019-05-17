/*
* Merge form data from FDF file to output PDF - flattened.
*
* Run as: go run pdf_form_fill_fdf_merge.go template.pdf input.fdf output.pdf
 */

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/fdf"
	"github.com/unidoc/unipdf/v3/model"
)

// Example of merging fdf data into a form.
func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Merge in form data from FDF to output PDF - flattened\n")
		fmt.Printf("Usage: go run pdf_form_fill_fdf_merge.go template.pdf input.fdf output.pdf\n")
		os.Exit(1)
	}

	// Enable debug-level logging.
	common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))

	templatePath := os.Args[1]
	fdfPath := os.Args[2]
	outputPath := os.Args[3]

	err := fdfMerge(templatePath, fdfPath, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Success, output written to %s\n", outputPath)
}

// fdfMerge loads template PDF in `templatePath` and FDF form data from `fdfPath` and fills into the fields,
// flattens and outputs as a PDF to `outputPath`.
func fdfMerge(templatePath, fdfPath, outputPath string) error {
	fdfData, err := fdf.LoadFromPath(fdfPath)
	if err != nil {
		return err
	}

	f, err := os.Open(templatePath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return err
	}

	// Populate the form data.
	err = pdfReader.AcroForm.Fill(fdfData)
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
