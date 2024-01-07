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
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/contentstream"
	"github.com/unidoc/unipdf/v3/core"
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

	var form *model.PdfAcroForm

	// Generate a new AcroForm instead of copying from the source PDF.
	opt := &model.ReaderToWriterOpts{
		SkipAcroForm: true,
		PageProcessCallback: func(pageNum int, page *model.PdfPage) error {
			if pageNum == 1 {
				form = createForm(page)
			}

			return nil
		},
	}

	// Generate a PdfWriter instance from existing PdfReader.
	pdfWriter, err := pdfReader.ToWriter(opt)
	if err != nil {
		return err
	}

	// Set new AcroForm.
	err = pdfWriter.SetForms(form)
	if err != nil {
		return err
	}

	return pdfWriter.WriteToFile(outputPath)
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

	addSubmitButton(page, form)

	return form
}

func addSubmitButton(page *model.PdfPage, form *model.PdfAcroForm) error {
	submitUrl := "https://url-to-backend"
	btnX := 400.0
	btnY := 400.0
	btnW := 50.0
	btnH := 20.0

	// Construct a new pdf field for the submit button.
	field := model.NewPdfField()
	btnField := &model.PdfFieldButton{}
	field.SetContext(btnField)
	btnField.PdfField = field

	btnField.T = core.MakeString("btnSubmit")
	btnField.SetType(model.ButtonTypePush)

	btnField.V = core.MakeName("Off")
	btnField.TU = core.MakeString("submit")

	// Define an Action that will be executed when the button is clicked.
	submitAction := model.NewPdfActionSubmitForm()
	submitAction.Flags = core.MakeInteger(5) // This will submit all fields value using HTTP Form format.
	submitAction.F = model.NewPdfFilespec()
	submitAction.F.F = core.MakeString(submitUrl)
	submitAction.F.FS = core.MakeName("URL")

	widgetText := core.MakeDict()
	widgetText.Set(*core.MakeName("CA"), core.MakeString("Submit"))

	fontData, err := model.NewStandard14Font("Helvetica")
	if err != nil {
		return err
	}
	fontName := core.MakeName("Helv")

	// Construct button.
	cc := contentstream.NewContentCreator()
	cc.Add_q()
	cc.Add_g(0.7)
	cc.Add_re(0, 0, btnW, btnH)
	cc.Add_f()
	cc.Add_Q()
	cc.Add_q()
	cc.Add_BT()
	cc.Add_Tf(*fontName, 10)
	cc.Add_g(0)
	cc.Add_Td(5, 5)
	cc.Add_Tj(*core.MakeString("Submit"))
	cc.Add_ET()
	cc.Add_Q()

	xform := model.NewXObjectForm()
	xform.SetContentStream(cc.Bytes(), core.NewRawEncoder())
	xform.BBox = core.MakeArrayFromFloats([]float64{0, 0, btnW, btnH})
	xform.Resources = model.NewPdfPageResources()
	xform.Resources.SetFontByName(*fontName, fontData.ToPdfObject())

	appearance := core.MakeDict()
	appearance.Set("N", xform.ToPdfObject())

	// Construct a widget annotation.
	btnWidget := model.NewPdfAnnotationWidget()
	btnWidget.Rect = core.MakeArrayFromFloats([]float64{btnX, btnY, btnX + btnW, btnY + btnH})
	btnWidget.P = page.ToPdfObject()
	btnWidget.F = core.MakeInteger(4)
	btnWidget.Parent = btnField.ToPdfObject()
	btnWidget.A = submitAction.ToPdfObject()
	btnWidget.MK = widgetText
	btnWidget.AP = appearance

	btnField.Annotations = append(btnField.Annotations, btnWidget)

	// Add widget to existing form.
	*form.Fields = append(*form.Fields, btnField.PdfField)
	page.AddAnnotation(btnField.Annotations[0].PdfAnnotation)

	return nil
}
