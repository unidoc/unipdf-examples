package main

import (
	"fmt"
	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/fjson"
	"github.com/unidoc/unipdf/v3/model"
	"os"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`)) //
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Syntax: go run pdf_update_existing_fields.go sample_form.pdf sample_form2.pdf\n")
	}
	inputPath := os.Args[1]
	outputPath := os.Args[2]

	err := updatePdfFields(inputPath, outputPath) //
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func updatePdfFields(inputPath, outputPath string) error { //
	f, err := os.Open(inputPath)
	if err != nil {
		fmt.Println("Error 1", err.Error())
		return err
	}
	fmt.Println("Input File: ", inputPath)

	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		fmt.Println("Error 2", err.Error())
		return err
	}

	acroForm := pdfReader.AcroForm

	if acroForm == nil {
		fmt.Println("No Form Data Found")
		return nil
	}

	fieldFallBacks := make(map[string]*annotator.AppearanceFont)
	fieldAppearance := annotator.FieldAppearance{OnlyIfMissing: false, RegenerateTextFields: true}

	fmt.Printf(" AcroForm (%p)\n", acroForm)
	fmt.Printf(" NeedAppearances: %v\n", acroForm.NeedAppearances)
	fmt.Printf(" SigFlags: %v\n", acroForm.SigFlags)
	fmt.Printf(" CO: %v\n", acroForm.CO)
	fmt.Printf(" DR: %v\n", acroForm.DR)
	fmt.Printf(" DA: %v\n", acroForm.DA)
	fmt.Printf(" Q: %v\n", acroForm.Q)
	fmt.Printf(" XFA: %v\n", acroForm.XFA)
	if acroForm.Fields != nil {
		fmt.Printf(" #Fields: %d\n", len(acroForm.AllFields()))
	} else {
		fmt.Printf("No fields set\n")
	}
	fmt.Printf(" =====\n")

	updatedFields := []struct {
		Name string
		Flag model.FieldFlag
		Font model.StdFontName
	}{
		{Name: "Name1", Flag: model.FieldFlagClear, Font: model.CourierBoldObliqueName},
		{Name: "Name2", Flag: model.FieldFlagClear, Font: model.CourierBoldObliqueName},
		{Name: "Name3", Flag: model.FieldFlagClear, Font: model.CourierBoldObliqueName},
		{Name: "Name4", Flag: model.FieldFlagClear, Font: model.CourierBoldObliqueName},
		{Name: "Name5", Flag: model.FieldFlagClear, Font: model.CourierBoldObliqueName},
		{Name: "Name6", Flag: model.FieldFlagMultiline, Font: model.CourierBoldObliqueName},
		{Name: "Name7", Flag: model.FieldFlagDoNotScroll, Font: model.CourierBoldObliqueName},
		{Name: "Name8", Flag: model.FieldFlagCombo, Font: model.CourierBoldObliqueName},
		{Name: "Name9", Flag: model.FieldFlagMultiline, Font: model.CourierBoldObliqueName},
		{Name: "Name10", Flag: model.FieldFlagDoNotScroll, Font: model.CourierBoldObliqueName},
		{Name: "Name11", Flag: model.FieldFlagDoNotSpellCheck, Font: model.CourierBoldObliqueName},
		{Name: "Name12", Flag: model.FieldFlagDoNotSpellCheck, Font: model.CourierBoldObliqueName},
		{Name: "Name13", Flag: model.FieldFlagDoNotSpellCheck, Font: model.CourierBoldObliqueName},
		{Name: "Name14", Flag: model.FieldFlagMultiline, Font: model.CourierBoldObliqueName},
	}

	fields := acroForm.AllFields()

	if len(fields) != len(updatedFields) {
		return fmt.Errorf("error -> the names provided are not complete")
	}
	var newFields []*model.PdfField

	for index, field := range fields {
		name := updatedFields[index]
		objectString := core.MakeString(name.Name)
		field.T = objectString
		field.SetFlag(name.Flag)
		font, err := model.NewStandard14Font(name.Font)
		if err != nil {
			return err
		}
		fieldFallBacks[name.Name] = &annotator.AppearanceFont{
			Name: font.FontDescriptor().FontName.String(),
			Font: font,
			Size: 0,
		}
		newFields = append(newFields, field)
	}
	defaultFontReplacement, err := model.NewPdfFontFromTTFFile("./DoHyeon-Regular.ttf") //model.NewStandard14Font(model.HelveticaObliqueName)
	style := fieldAppearance.Style()
	style.Fonts = &annotator.AppearanceFontStyle{
		Fallback: &annotator.AppearanceFont{
			Font: defaultFontReplacement,
			Name: defaultFontReplacement.FontDescriptor().FontName.String(),
			Size: 0,
		},
		FieldFallbacks: fieldFallBacks,
		ForceReplace:   true,
	}

	acroForm.Fields = &newFields

	fieldAppearance.SetStyle(style)
	fdata, err := fjson.LoadFromPDFFile(inputPath)
	if err != nil {
		fmt.Println("Error Here 1 -> ", err.Error())
		return err
	}

	err = pdfReader.AcroForm.FillWithAppearance(fdata, fieldAppearance)
	if err != nil {
		fmt.Println("Error 6: ", err.Error())
		return err
	}

	//err = pdfReader.FlattenFields(true, fieldAppearance)
	//if err != nil {
	//	fmt.Println("Error 7: ", err.Error())
	//	return err
	//}

	// The document AcroForm field is no longer needed.
	opt := &model.ReaderToWriterOpts{
		SkipAcroForm: false,
	}

	pdfWriter, err := pdfReader.ToWriter(opt)
	if err != nil {
		fmt.Println("Error PDF Writer", err.Error())
		return err
	}

	err = pdfWriter.WriteToFile(outputPath)
	if err != nil {
		return err
	}

	return nil
}
