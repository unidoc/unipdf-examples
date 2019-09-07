/*
 * Lists form fields in a PDF file.
 *
 * Run as: go run pdf_form_list_fields.go input.pdf
 */

package main

import (
	"errors"
	"fmt"
	"os"

	unicommon "github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run pdf_forms_list_fields.go input.pdf [input2.pdf] ...\n")
		os.Exit(1)
	}

	// When debugging, enable debug-level logging via console:
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	for _, inputPath := range os.Args[1:len(os.Args)] {
		err := listFormFields(inputPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

func listFormFields(inputPath string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	fmt.Printf("Input file: %s\n", inputPath)

	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return err
	}

	acroForm := pdfReader.AcroForm
	if acroForm == nil {
		fmt.Printf(" No formdata present\n")
		return nil
	}

	fmt.Printf(" AcroForm (%p)\n", acroForm)
	fmt.Printf(" NeedAppearances: %v\n", acroForm.NeedAppearances)
	fmt.Printf(" SigFlags: %v\n", acroForm.SigFlags)
	fmt.Printf(" CO: %v\n", acroForm.CO)
	fmt.Printf(" DR: %v\n", acroForm.DR)
	fmt.Printf(" DA: %v\n", acroForm.DA)
	fmt.Printf(" Q: %v\n", acroForm.Q)
	fmt.Printf(" XFA: %v\n", acroForm.XFA)
	if acroForm.Fields != nil {
		fmt.Printf(" #Fields: %d\n", len(*acroForm.Fields))
	} else {
		fmt.Printf("No fields set\n")
	}
	fmt.Printf(" =====\n")

	fields := acroForm.AllFields()

	for idx, field := range fields {
		fmt.Printf("=====\n")
		fmt.Printf("Field %d\n", idx+1)
		if !field.IsTerminal() {
			fmt.Printf("- Skipping over non-terminal field\n")
			continue
		}

		fullname, err := field.FullName()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		fflags := field.Flags()
		fmt.Printf("Name: %v\n", fullname)
		fmt.Printf(" Flags: %s (%d)\n", fflags, fflags)

		ctx := field.GetContext()
		switch t := ctx.(type) {
		case *model.PdfFieldButton:
			fmt.Printf(" Button\n")
			if t.IsCheckbox() {
				fmt.Printf(" - Checkbox\n")
			}
			if t.IsPush() {
				fmt.Printf(" - Push\n")
			}
			if t.IsRadio() {
				fmt.Printf(" - Radio\n")
			}
			fmt.Printf(" - '%v'\n", t.V)
		case *model.PdfFieldText:
			fmt.Printf(" Text\n")
			fmt.Printf(" - '%v'\n", t.V)
			if str, ok := core.GetString(t.V); ok {
				fmt.Printf(" - Decoded: '%s'\n", str.Decoded())
			}

		case *model.PdfFieldChoice:
			fmt.Printf(" Choice\n")
			fmt.Printf(" - '%v'\n", t.V)
		case *model.PdfFieldSignature:
			fmt.Printf(" Signature\n")
			fmt.Printf(" - '%v'\n", t.V)
		default:
			fmt.Printf(" Unknown\n")
			continue
		}

		fmt.Printf(" Annotations: %d\n", len(field.Annotations))
		for j, wa := range field.Annotations {
			fmt.Printf(" - Annotation %d \n", j+1)

			// Note: The P field is optional
			if wa.P != nil {
				pageind, ok := wa.P.(*core.PdfIndirectObject)
				if !ok {
					fmt.Printf("not indirect, got: %v\n", pageind)
					fmt.Printf("%")
					return errors.New("Type check error")
				}
				_, pagenum, err := pdfReader.PageFromIndirectObject(pageind)
				if err != nil {
					return err
				}
				fmt.Printf(" - Page number: %d\n", pagenum)
			} else {
				// If P not set, go through pages and look for match.
				// TODO: Make a map of widget annotations to page numbers a priori.
				for pageidx, page := range pdfReader.PageList {
					annotations, err := page.GetAnnotations()
					if err != nil {
						return err
					}
					for _, annot := range annotations {
						switch t := annot.GetContext().(type) {
						case *model.PdfAnnotationWidget:
							if wa == t {
								fmt.Printf(" - Page number: %d\n", pageidx+1)
							}
						}
					}
				}
			}
			fmt.Printf(" - Rect: %+v\n", wa.Rect)
			fmt.Printf(" - wa.AS: %v\n", wa.AS)
			fmt.Printf(" - wa.AP: %v\n", wa.AP)
			fmt.Printf(" - wa.F:  %v\n", wa.F)

			// Example of how to fetch the appearance stream data.
			if apDict, has := core.GetDict(wa.AP); has {
				n, has := core.GetStream(apDict.Get("N"))
				if has {
					decoded, err := core.DecodeStream(n)
					if err != nil {
						fmt.Printf("Decoding error: %v\n", err)
						return err
					}
					fmt.Printf("   - N: '%s'\n", decoded)
				} else {
					fmt.Printf("   - N not set\n")
				}

				if d, has := core.GetDict(apDict.Get("D")); has {
					appKey := core.MakeName("Off")
					if appname, has := core.GetName(wa.AS); has {
						appKey = appname
					}

					fmt.Printf("   - D dict: % s\n", d.Keys())
					fmt.Printf("   - App Key: '%s'\n", *appKey)
					if radioApp, has := core.GetStream(d.Get(*appKey)); has {
						decoded, err := core.DecodeStream(radioApp)
						if err != nil {
							fmt.Printf("  - Decoding error: %v\n", err)
							return err
						}

						fmt.Printf("   - Radio appearance: '%s'\n", decoded)
					}

				}

			} else {
				fmt.Printf("   - Appearance dict not present: %s\n", apDict)
			}
		}
	}

	return nil
}
