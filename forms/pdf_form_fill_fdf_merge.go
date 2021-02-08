/*
* Merge form data from FDF file to output PDF - flattened.
*
* Run as: go run pdf_form_fill_fdf_merge.go template.pdf input.fdf output.pdf
 */

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/fdf"
	"github.com/unidoc/unipdf/v3/model"
)

const licenseKey = `
-----BEGIN UNIDOC LICENSE KEY-----
Free trial license keys are available at: https://unidoc.io/
-----END UNIDOC LICENSE KEY-----
`

func init() {
	// Enable debug-level logging.
	// unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	err := license.SetLicenseKey(licenseKey, `Company Name`)
	if err != nil {
		panic(err)
	}
}

// Example of merging fdf data into a form.
func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Merge in form data from FDF to output PDF - flattened\n")
		fmt.Printf("Usage: go run pdf_form_fill_fdf_merge.go template.pdf input.fdf output.pdf\n")
		os.Exit(1)
	}

	templatePath := os.Args[1]
	fdfPath := os.Args[2]
	outputPath := os.Args[3]

	flatten := false
	err := fdfMerge(templatePath, fdfPath, outputPath, flatten)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Success, output written to %s\n", outputPath)
}

// fdfMerge loads template PDF in `templatePath` and FDF form data from `fdfPath` and fills into the fields,
// flattens (if flatten specified) and outputs as a PDF to `outputPath`.
func fdfMerge(templatePath, fdfPath, outputPath string, flatten bool) error {
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
	if flatten {
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
	} else {
		fieldAppearance := annotator.FieldAppearance{OnlyIfMissing: true, RegenerateTextFields: true}
		for _, field := range pdfReader.AcroForm.AllFields() {
			for _, wa := range field.Annotations {
				// appgen generates the appearance based on the form/field/annotation and other settings
				// based on the implementation (for example may only generate appearance if none set).
				apDict, err := fieldAppearance.GenerateAppearanceDict(pdfReader.AcroForm, field, wa)
				if err != nil {
					return err
				}
				wa.AP = apDict
				// Force update of the widget appearance.
				_ = wa.ToPdfObject()
			}
		}
	}

	// Don't copy AcroForm when flattening.
	opt := &model.ReaderToWriterOpts{
		SkipAcroForm: flatten,
	}

	// Generate a PdfWriter instance from existing PdfReader.
	pdfWriter, err := pdfReader.ToWriter(opt)
	if err != nil {
		return err
	}

	// Write to file.
	err = pdfWriter.WriteToFile(outputPath)
	return err
}
