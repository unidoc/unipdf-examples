/*
* Merge form data from FDF file to output PDF - flattened.
*
* Run as: go run pdf_form_fill_fdf_merge.go template.pdf input.fdf output.pdf
 */

package main

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/unidoc/unidoc/common"
	"github.com/unidoc/unidoc/pdf/annotator"
	"github.com/unidoc/unidoc/pdf/core"
	"github.com/unidoc/unidoc/pdf/fdf"
	"github.com/unidoc/unidoc/pdf/model"
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
	fieldMap, err := getFDFData(fdfPath)
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

	keys := []string{}
	for fieldName := range fieldMap {
		keys = append(keys, fieldName)
	}
	sort.Strings(keys)

	// Fill in the values.
	for _, fieldName := range keys {
		fieldDict := fieldMap[fieldName]
		val := core.TraceToDirectObject(fieldDict.Get("V"))
		fmt.Printf("Filling: %s: %+v %T\n", fieldName, val, val)
		found := false
		for _, f := range pdfReader.AcroForm.AllFields() {
			if f.PartialName() == fieldName {
				switch f.GetContext().(type) {
				case *model.PdfFieldText:
					switch t := val.(type) {
					case *core.PdfObjectName:
						name := t
						str := name.String()
						if len(str) > 1 && str[0] == '{' && str[len(str)-1] == '}' {
							fmt.Printf("%s - calculated field - skipping\n", fieldName)
							continue
						}
						fmt.Printf("ERROR got V as name -> converting to string: '%s'\n", name.String())
						f.V = core.MakeString(name.String())
					case *core.PdfObjectString:
						//string := t
						//f.V = t
						f.V = core.MakeEncodedString(t.String())
					default:
						fmt.Printf("Unsupported text field V: %T (%#v)\n", t, t)
					}
				case *model.PdfFieldButton, *model.PdfFieldChoice:
					switch val.(type) {
					case *core.PdfObjectName:
						for _, wa := range f.Annotations {
							wa.AS = val
						}
						f.V = val
					default:
						fmt.Printf("UNEXPECTED %s -> %v\n", fieldName, val)
						f.V = val
					}
				case *model.PdfFieldSignature:
					fmt.Printf("Signature not supported yet: %s/%v\n", fieldName, val)
				}

				found = true
				break
			}
		}
		if !found {
			fmt.Printf("'%s' NOT FOUND - not filled\n", fieldName)
		}
	}

	// Flatten.
	fieldAppearance := annotator.FieldAppearance{OnlyIfMissing: false}
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

func getFDFData(fdfPath string) (map[string]*core.PdfObjectDictionary, error) {
	f, err := os.Open(fdfPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fieldDataMap := map[string]*core.PdfObjectDictionary{}

	p, err := fdf.NewParser(f)
	if err != nil {
		return nil, err
	}

	fdfDict, err := p.Root()
	if err != nil {
		return nil, err
	}

	fields, found := core.GetArray(fdfDict.Get("Fields"))
	if !found {
		return nil, errors.New("Fields missing")
	}

	for i := 0; i < fields.Len(); i++ {
		fieldDict, has := core.GetDict(fields.Get(i))
		if has {
			// Key value field data.
			t, _ := core.GetString(fieldDict.Get("T"))
			if t != nil {
				fieldDataMap[t.Str()] = fieldDict
			}
		}
	}

	return fieldDataMap, nil
}
