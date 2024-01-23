/*
 * Create a new form with a submit and reset button.
 * The example shows how to add a submit button and a reset button to a form.
 * The submit button could be used to send the form data to a backend server,
 * and the reset button would reset the form fields value.
 *
 * Run as: go run pdf_form_action.go
 */

package main

import (
	"log"
	"os"

	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/contentstream/draw"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/creator"
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
	textFieldsDef := []struct {
		Label string
		Name  string
		Rect  []float64
	}{
		{Label: "Full Name", Name: "full_name", Rect: []float64{123.97, 619.02, 343.99, 633.6}},
		{Label: "Address 1", Name: "address_line_1", Rect: []float64{123.97, 596.82, 343.99, 611.4}},
		{Label: "Address 2", Name: "address_line_2", Rect: []float64{123.97, 574.28, 343.99, 588.86}},
	}

	c := creator.New()
	page := c.NewPage()
	_, pageHeight, err := page.Size()
	if err != nil {
		log.Fatal(err)
	}

	form := model.NewPdfAcroForm()
	fields := core.MakeArray()

	// Create text fields and it's label
	for _, fdef := range textFieldsDef {
		opt := annotator.TextFieldOptions{}
		textf, err := annotator.NewTextField(page, fdef.Name, fdef.Rect, opt)
		if err != nil {
			log.Fatal(err)
		}

		*form.Fields = append(*form.Fields, textf.PdfField)
		page.AddAnnotation(textf.Annotations[0].PdfAnnotation)

		y := pageHeight - fdef.Rect[1]

		p := c.NewParagraph(fdef.Label)
		p.SetPos(fdef.Rect[0]-80, y-10)
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

	err = c.WriteToFile("form_with_action_button.pdf")
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
