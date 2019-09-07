/*
 * Create and apply a new form to an existing PDF.
 * The example shows how to load template1.pdf and add an interactive form to it and save it
 * as template1_with_form.pdf.
 *
 * Run as: go run pdf_form_add.go
 */

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/model"
)

func main() {
	inputPath := `template1.pdf`
	outputPath := `template1_with_form.pdf`

	err := addFormToPdf(inputPath, outputPath)
	if err != nil {
		fmt.Printf("Failed to add form: %#v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Success, output in %s\n", outputPath)
}

// addFormToPdf adds the form to the PDF specified by `inputPath` and outputs to `outputPath`.
func addFormToPdf(inputPath string, outputPath string) error {
	// Read the input pdf file.
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	pdfWriter := model.NewPdfWriter()

	// Load the pages.
	for i := 0; i < numPages; i++ {
		page, err := pdfReader.GetPage(i + 1)
		if err != nil {
			return err
		}

		if i == 0 {
			err = pdfWriter.SetForms(createForm(page))
			if err != nil {
				return err
			}
		}

		err = pdfWriter.AddPage(page)
		if err != nil {
			return err
		}
	}

	of, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer of.Close()

	return pdfWriter.Write(of)
}

// textFieldsDef is a list of text fields to add to the form. The Rect field specifies the coordinates of the
// field.
var textFieldsDef = []struct {
	Name string
	Rect []float64
}{
	{Name: "full_name", Rect: []float64{123.97, 619.02, 343.99, 633.6}},
	{Name: "address_line_1", Rect: []float64{142.86, 596.82, 347.3, 611.4}},
	{Name: "address_line_2", Rect: []float64{143.52, 574.28, 347.96, 588.86}},
	{Name: "age", Rect: []float64{95.15, 551.75, 125.3, 566.33}},
	{Name: "city", Rect: []float64{96.47, 506.35, 168.37, 520.93}},
	{Name: "country", Rect: []float64{114.69, 483.82, 186.59, 498.4}},
}

// checkboxFieldDefs is a list of checkboxes to add to the form.
var checkboxFieldDefs = []struct {
	Name    string
	Rect    []float64
	Checked bool
}{
	{Name: "male", Rect: []float64{113.7, 525.57, 125.96, 540.15}, Checked: true},
	{Name: "female", Rect: []float64{157.44, 525.24, 169.7, 539.82}, Checked: false},
}

// choiceFieldDefs is a list of comboboxes to add to the form with specified options.
var choiceFieldDefs = []struct {
	Name    string
	Rect    []float64
	Options []string
}{
	{
		Name:    "fav_color",
		Rect:    []float64{144.52, 461.61, 243.92, 476.19},
		Options: []string{"Black", "Blue", "Green", "Orange", "Red", "White", "Yellow"},
	},
}

// createForm creates the form and fields to be placed on the `page`.
func createForm(page *model.PdfPage) *model.PdfAcroForm {
	form := model.NewPdfAcroForm()

	// Add ZapfDingbats font.
	zapfdb := model.NewStandard14FontMustCompile(model.ZapfDingbatsName)
	form.DR = model.NewPdfPageResources()
	form.DR.SetFontByName(`ZaDb`, zapfdb.ToPdfObject())

	for _, fdef := range textFieldsDef {
		opt := annotator.TextFieldOptions{}
		textf, err := annotator.NewTextField(page, fdef.Name, fdef.Rect, opt)
		if err != nil {
			panic(err)
		}

		*form.Fields = append(*form.Fields, textf.PdfField)
		page.AddAnnotation(textf.Annotations[0].PdfAnnotation)
	}

	for _, cbdef := range checkboxFieldDefs {
		opt := annotator.CheckboxFieldOptions{}
		checkboxf, err := annotator.NewCheckboxField(page, cbdef.Name, cbdef.Rect, opt)
		if err != nil {
			panic(err)
		}

		*form.Fields = append(*form.Fields, checkboxf.PdfField)
		page.AddAnnotation(checkboxf.Annotations[0].PdfAnnotation)
	}

	for _, chdef := range choiceFieldDefs {
		opt := annotator.ComboboxFieldOptions{Choices: chdef.Options}
		comboboxf, err := annotator.NewComboboxField(page, chdef.Name, chdef.Rect, opt)
		if err != nil {
			panic(err)
		}

		*form.Fields = append(*form.Fields, comboboxf.PdfField)
		page.AddAnnotation(comboboxf.Annotations[0].PdfAnnotation)
	}

	return form
}
