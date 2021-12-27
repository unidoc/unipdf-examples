package main

//func init() {
//	// Make sure to load your metered License API key prior to using the library.
//	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
//	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
//	if err != nil {
//		panic(err)
//	}
//}
//
//func main() {
//	if len(os.Args) < 3 {
//		fmt.Printf("Syntax: go run pdf_update_existing_fields.go sample_form.pdf sample_form2.pdf\n")
//	}
//	inputPath := os.Args[1]
//	outputPath := os.Args[2]
//
//	err := updatePdfFields(inputPath, outputPath) //
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		os.Exit(1)
//	}
//}
//
//type NewData struct {
//	NewName  string
//	Font     model.StdFontName
//	FontSize int
//	Flag     model.FieldFlag
//}
//
//// Here you Put All Changes you Want to Implement.
//var newNames = []NewData{
//	{"Name1", model.HelveticaBoldObliqueName, 16, model.FieldFlagDoNotSpellCheck},
//	{"Name2", model.TimesBoldItalicName, 12, model.FieldFlagDoNotSpellCheck},
//	{"Name3", model.TimesItalicName, 14, model.FieldFlagDoNotSpellCheck},
//	{"Name4", model.CourierBoldObliqueName, 16, model.FieldFlagDoNotSpellCheck},
//	{"Name5", model.HelveticaName, 18, model.FieldFlagMultiline},
//	{"Name6", model.CourierBoldObliqueName, 10, model.FieldFlagDoNotSpellCheck},
//	{"Name7", model.CourierBoldObliqueName, 16, model.FieldFlagMultiline},
//	{"Name8", model.CourierName, 16, model.FieldFlagMultiline},
//	{"Name9", model.TimesBoldName, 18, model.FieldFlagDoNotSpellCheck},
//	{"Name10", model.HelveticaBoldName, 12, model.FieldFlagDoNotSpellCheck},
//	{"Name11", model.CourierObliqueName, 8, model.FieldFlagDoNotSpellCheck},
//	{"Name12", model.TimesRomanName, 16, model.FieldFlagDoNotSpellCheck},
//	{"Name13", model.CourierBoldObliqueName, 16, model.FieldFlagDoNotSpellCheck},
//	{"Name14", model.TimesBoldItalicName, 16, model.FieldFlagReadOnly},
//}
//
//func updatePdfFields(inputPath, outputPath string) error { //
//	f, err := os.Open(inputPath)
//	if err != nil {
//		return err
//	}
//
//	defer f.Close()
//
//	mapData := make(map[string]NewData, 0)
//
//	mapNames := make(map[string]string, 0)
//
//	pdfReader, err := model.NewPdfReader(f)
//	if err != nil {
//		return err
//	}
//
//	acroForm := pdfReader.AcroForm
//
//	if acroForm == nil {
//		return nil
//	}
//
//	fdata, err := fjson.LoadFromPDFFile(inputPath)
//	if err != nil {
//		return err
//	}
//
//	fieldFallBacks := make(map[string]*annotator.AppearanceFont)
//	fieldAppearance := annotator.FieldAppearance{OnlyIfMissing: true, RegenerateTextFields: true}
//
//	fields := acroForm.AllFields()
//
//	if len(fields) != len(newNames) {
//		return fmt.Errorf("error -> the names provided are not complete")
//	}
//
//	for k, v := range fields {
//		mapNames[v.T.String()] = newNames[k].NewName
//		mapData[v.T.String()] = newNames[k]
//	}
//
//	newFData := fdata.SetNewFieldData(mapNames)
//	if newFData == nil {
//		return fmt.Errorf("error -> No Data")
//	}
//
//	for _, v := range mapData {
//		font, err := model.NewStandard14Font(v.Font)
//		if err != nil {
//			return err
//		}
//		fieldFallBacks[v.NewName] = &annotator.AppearanceFont{
//			Name: font.FontDescriptor().FontName.String(),
//			Font: font,
//			Size: float64(v.FontSize),
//		}
//	}
//
//	defaultFontReplacement, err := model.NewStandard14Font(model.TimesItalicName)
//
//	style := fieldAppearance.Style()
//	style.Fonts = &annotator.AppearanceFontStyle{
//		Fallback: &annotator.AppearanceFont{
//			Font: defaultFontReplacement,
//			Name: defaultFontReplacement.FontDescriptor().FontName.String(),
//			Size: 14,
//		},
//		FieldFallbacks: fieldFallBacks,
//		ForceReplace:   true,
//	}
//	fieldAppearance.SetStyle(style)
//
//	for _, field := range fields {
//		if v, ok := mapData[field.T.String()]; ok {
//			name := v.NewName
//			if field.T.String() == " Tx" {
//				field.SetFlag(mapData[field.T.String()].Flag)
//			}
//			objectString := core.MakeString(name)
//			field.T = objectString
//		}
//	}
//
//	err = acroForm.FillWithAppearance(newFData, fieldAppearance)
//	if err != nil {
//		return err
//	}
//
//	// You can comment to not Flatten if you don't need it.
//	err = pdfReader.FlattenFields(true, fieldAppearance)
//	if err != nil {
//		return err
//	}
//	// The document AcroForm field is no longer needed.
//	opt := &model.ReaderToWriterOpts{
//		SkipAcroForm: false,
//	}
//
//	pdfWriter, err := pdfReader.ToWriter(opt)
//	if err != nil {
//		return err
//	}
//
//	err = pdfWriter.WriteToFile(outputPath)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
