package main

import (
	"fmt"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/fjson"
	"os"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`)) //license.SetMeteredKey("08fddff07958eaa43970d606bcda148243769221cb3e160bd1e24b2460e3ced4")
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) == 2 {
		fmt.Printf("Syntax: go run pdf_fdata_check.go sample_form2.pdf.pdf\n") // output.pdf
	}
	inputPath := os.Args[1]
	outputPath := "" //os.Args[2]

	err := CheckFData(inputPath, outputPath) //
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func CheckFData(input, output string) error {
	fdata, err := fjson.LoadFromPDFFile(input)
	if err != nil {
		fmt.Println("Error Here 1 -> ", err.Error())
		return err
	}
	values, err := fdata.FieldValues()
	if err != nil {
		fmt.Println("Error Here 2 -> ", err.Error())
		return err
	}
	for k, v := range values {
		fmt.Println(k, "->", v)
	}

	//f, err := os.Open(input)
	//if err != nil {
	//	fmt.Println("Error 3: ", err.Error())
	//	return err
	//}
	//defer f.Close()
	//
	//pdfReader, err := model.NewPdfReader(f)
	//if err != nil {
	//	fmt.Println("Error 4: ", err.Error())
	//	return err
	//}
	//
	//fieldAppearance := annotator.FieldAppearance{OnlyIfMissing: true, RegenerateTextFields: false}
	//
	//// set font using standard font
	//defaultFontReplacement, err := model.NewPdfFontFromTTFFile("./DoHyeon-Regular.ttf") //model.NewStandard14Font(model.HelveticaObliqueName)
	//
	//// set font using ttf font file
	//fontReplacement, err := model.NewPdfFontFromTTFFile("./DoHyeon-Regular.ttf")
	//
	//// use composite ttf font file
	//// refer to `text/pdf_using_cjk_font.go` example file for more information
	//cjkFont, err := model.NewCompositePdfFontFromTTFFile("./DoHyeon-Regular.ttf")
	//
	//if err != nil {
	//	fmt.Println("Error 5: ", err.Error())
	//	return err
	//}
	//
	//// replace email field's font using `fontReplacement`
	//// and set the other field's font using `defaultFontReplacement`
	//style := fieldAppearance.Style()
	//style.Fonts = &annotator.AppearanceFontStyle{
	//	Fallback: &annotator.AppearanceFont{
	//		Font: defaultFontReplacement,
	//		Name: defaultFontReplacement.FontDescriptor().FontName.String(),
	//		Size: 0,
	//	},
	//	FieldFallbacks: map[string]*annotator.AppearanceFont{
	//		"email4": {
	//			Font: fontReplacement,
	//			Name: fontReplacement.FontDescriptor().FontName.String(),
	//			Size: 0,
	//		},
	//		"address5[addr_line1]": {
	//			Font: cjkFont,
	//			Name: cjkFont.FontDescriptor().FontName.String(),
	//			Size: 0,
	//		},
	//		"address5[addr_line2]": {
	//			Font: cjkFont,
	//			Name: cjkFont.FontDescriptor().FontName.String(),
	//			Size: 0,
	//		},
	//		"address5[city]": {
	//			Font: cjkFont,
	//			Name: cjkFont.FontDescriptor().FontName.String(),
	//			Size: 0,
	//		},
	//		"address5[state]": {
	//			Font: cjkFont,
	//			Name: cjkFont.FontDescriptor().FontName.String(),
	//			Size: 0,
	//		},
	//		"address5[postal]": {
	//			Font: cjkFont,
	//			Name: cjkFont.FontDescriptor().FontName.String(),
	//			Size: 0,
	//		},
	//	},
	//	ForceReplace: true,
	//}
	//
	//fieldAppearance.SetStyle(style)
	//
	//// Populate the form data.
	//err = pdfReader.AcroForm.FillWithAppearance(fdata, fieldAppearance)
	//if err != nil {
	//	fmt.Println("Error 6: ", err.Error())
	//	return err
	//}
	//
	//// Flatten form.
	////err = pdfReader.FlattenFields(false, fieldAppearance)
	////if err != nil {
	////	fmt.Println("Error 7: ", err.Error())
	////	return err
	////}
	//
	//// The document AcroForm field is no longer needed.
	//opt := &model.ReaderToWriterOpts{
	//	SkipAcroForm: false,
	//}
	//
	//// Generate a PdfWriter instance from existing PdfReader.
	//pdfWriter, err := pdfReader.ToWriter(opt)
	//if err != nil {
	//	fmt.Println("Error 8: ", err.Error())
	//	return err
	//}
	//
	//// Subset the composite font file to reduce pdf file size.
	//// Refer to `text/pdf_using_cjk_font.go` example file for more information
	//err = cjkFont.SubsetRegistered()
	//if err != nil {
	//	fmt.Println("Error 9: ", err.Error())
	//	return err
	//}
	//
	//// Write to file.
	//err = pdfWriter.WriteToFile(output)
	return nil
}
