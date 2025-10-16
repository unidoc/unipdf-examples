// This example demonstrates how to create a PDF with form fields
// using the UniPDF library. The PDF will contain multiple text fields
// along with submit and reset buttons, and the form fields will be
// properly tagged in the document structure tree for accessibility compliance.
//
// The example covers best practices for PDF/UA compliance:
// 1. Each form field has an associated label with a tooltip.
// 2. The document structure tree is constructed with K dictionaries
//    to represent the hierarchical structure of the form elements.
// 3. Each label is associated with its corresponding form field using
//    marked content IDs (MCID).
// 4. The submit button is configured to submit the form data to a specified URL.
// 5. The reset button is configured to reset the specified fields to their default values.

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/unidoc/unipdf/v4/annotator"
	"github.com/unidoc/unipdf/v4/common/license"
	"github.com/unidoc/unipdf/v4/contentstream/draw"
	"github.com/unidoc/unipdf/v4/core"
	"github.com/unidoc/unipdf/v4/creator"
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

func main() {
	textFieldsDef := []struct {
		Label   string
		Name    string
		Rect    []float64
		Tooltip string
	}{
		{Label: "Full Name", Name: "full_name", Rect: []float64{123.97, 619.02, 343.99, 633.6}, Tooltip: "Enter your full name"},
		{Label: "Address 1", Name: "address_line_1", Rect: []float64{123.97, 596.82, 343.99, 611.4}, Tooltip: "Enter your primary address"},
		{Label: "Address 2", Name: "address_line_2", Rect: []float64{123.97, 574.28, 343.99, 588.86}, Tooltip: "Enter your secondary address (optional)"},
	}

	c := creator.New()
	page := c.NewPage()
	_, pageHeight, err := page.Size()
	if err != nil {
		log.Fatal(err)
	}

	// Construct the StructTreeRoot.
	str := model.NewStructTreeRoot()

	// Construct base K dictionary.
	docK := model.NewKDictionary()
	docK.S = core.MakeName(string(model.StructureTypeDocument))

	str.AddKDict(docK)

	form := model.NewPdfAcroForm()
	fields := core.MakeArray()

	// Create text fields and it's label
	for idx, fdef := range textFieldsDef {
		opt := annotator.TextFieldOptions{}
		textf, err := annotator.NewTextField(page, fdef.Name, fdef.Rect, opt)
		if err != nil {
			log.Fatal(err)
		}

		textf.DV = core.MakeString("")           // Set default value for the field.
		textf.V = core.MakeString("")            // Set current value for the field.
		textf.TU = core.MakeString(fdef.Tooltip) // Set tooltip for the field.

		*form.Fields = append(*form.Fields, textf.PdfField)
		page.AddAnnotation(textf.Annotations[0].PdfAnnotation)

		y := pageHeight - fdef.Rect[1]

		p := c.NewStyledParagraph()
		p.SetText(fdef.Label)
		p.SetPos(fdef.Rect[0]-80, y-10)

		// Tag the label paragraph and generate its K dictionary.
		// This will be used to associate the label with the form field.
		p.SetStructureType(model.StructureTypeForm)

		// Set unique ID for each field label.
		p.SetMarkedContentID(int64(idx))

		k, err := p.GenerateKDict()
		if err != nil {
			fmt.Errorf("Error: %v", err)
		}

		k.Alt = core.MakeString(fdef.Tooltip) // Set alternative text for the label.

		docK.AddKChild(k)

		err = c.Draw(p)
		if err != nil {
			log.Fatal(err)
		}

		line := c.NewLine(fdef.Rect[0], y, fdef.Rect[2], y)
		err = c.Draw(line)
		if err != nil {
			log.Fatal(err)
		}

		fields.Append(textf.ToPdfObject())
	}

	err = addSubmitButton(page, form)
	if err != nil {
		log.Fatal(err)
	}

	err = addResetButton(page, form, fields)
	if err != nil {
		log.Fatal(err)
	}

	c.SetForms(form)
	c.SetStructTreeRoot(str)

	err = c.WriteToFile("pdf_tag_form.pdf")
	if err != nil {
		log.Fatal(err)
	}
}

// Add Submit button that will submit all fields value.
func addSubmitButton(page *model.PdfPage, form *model.PdfAcroForm) error {
	optSubmit := annotator.FormSubmitActionOptions{
		Url: "https://unidoc.io",
		Rectangle: draw.Rectangle{
			X:         400.0,
			Y:         400.0,
			Width:     50.0,
			Height:    20.0,
			FillColor: model.NewPdfColorDeviceRGB(0.0, 1.0, 0.0),
		},
		Label:      "Submit",
		LabelColor: model.NewPdfColorDeviceRGB(1.0, 0.0, 0.0),
	}

	btnSubmitField, err := annotator.NewFormSubmitButtonField(page, optSubmit)
	if err != nil {
		return err
	}

	*form.Fields = append(*form.Fields, btnSubmitField.PdfField)
	page.AddAnnotation(btnSubmitField.Annotations[0].PdfAnnotation)

	return nil
}

// Add Reset button that would reset the specified fields to it's default value.
func addResetButton(page *model.PdfPage, form *model.PdfAcroForm, fields *core.PdfObjectArray) error {
	optReset := annotator.FormResetActionOptions{
		Rectangle: draw.Rectangle{
			X:         100.0,
			Y:         400.0,
			Width:     50.0,
			Height:    20.0,
			FillColor: model.NewPdfColorDeviceGray(0.5),
		},
		Label:      "Reset",
		LabelColor: model.NewPdfColorDeviceGray(1.0),
		Fields:     fields,
	}

	btnResetField, err := annotator.NewFormResetButtonField(page, optReset)
	if err != nil {
		return err
	}

	// Add widget to existing form.
	*form.Fields = append(*form.Fields, btnResetField.PdfField)
	page.AddAnnotation(btnResetField.Annotations[0].PdfAnnotation)

	return nil
}
