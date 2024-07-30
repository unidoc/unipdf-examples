/*
 * Create a new form with customized text field colors.
 * The example shows how to create a form that has customized text field colors.
 * The form has submit and reset buttons.
 * Run as: go pdf_form_with_text_color.go
 */
package main

import (
	"fmt"
	"os"
	"path/filepath"

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
		Label       string
		Name        string
		SampleInput string
		Rect        []float64
	}{
		{Label: "Full Name", Name: "full_name", SampleInput: "Enter Full Name", Rect: []float64{123.97, 619.02, 343.99, 633.6}},
		{Label: "Address 1", Name: "address_line_1", SampleInput: "Enter Address 1", Rect: []float64{123.97, 596.82, 343.99, 611.4}},
		{Label: "Address 2", Name: "address_line_2", SampleInput: "Enter Address 2", Rect: []float64{123.97, 574.28, 343.99, 588.86}},
	}

	c := creator.New()
	page := c.NewPage()

	_, pageHeight, err := page.Size()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	form := model.NewPdfAcroForm()
	fields := core.MakeArray()

	for _, fdef := range textFieldsDef {
		opt := annotator.TextFieldOptions{
			TextColor: "#0000FF", // Set text color to Blue.
			Value:     fdef.SampleInput,
		}
		textf, err := annotator.NewTextField(page, fdef.Name, fdef.Rect, opt)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		*form.Fields = append(*form.Fields, textf.PdfField)
		page.AddAnnotation(textf.Annotations[0].PdfAnnotation)

		y := pageHeight - fdef.Rect[1]

		p := c.NewParagraph(fdef.Label)
		p.SetPos(fdef.Rect[0]-80, y-10)
		err = c.Draw(p)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		line := c.NewLine(fdef.Rect[0], y, fdef.Rect[2], y)
		err = c.Draw(line)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fields.Append(textf.ToPdfObject())
	}

	// Add Submit button
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
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	*form.Fields = append(*form.Fields, btnSubmitField.PdfField)
	page.AddAnnotation(btnSubmitField.Annotations[0].PdfAnnotation)

	// Add Reset button
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
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Add widget to existing form.
	*form.Fields = append(*form.Fields, btnResetField.PdfField)
	page.AddAnnotation(btnResetField.Annotations[0].PdfAnnotation)

	c.SetForms(form)

	model.SetPdfProducer(`UniDoc v3.61.0 (Unlicensed) - http://unidoc.io`)
	defer model.SetPdfProducer("")

	outPath := filepath.Join(".", "form_field_with_colored_text_fields.pdf")
	err = c.WriteToFile(outPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
