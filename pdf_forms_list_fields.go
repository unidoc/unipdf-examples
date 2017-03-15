/*
 * Lists form fields in a PDF file.
 *
 * Run as: go run pdf_forms_list_fields.go input.pdf
 */

package main

import (
	"fmt"
	"os"

	unicommon "github.com/unidoc/unidoc/common"
	unipdf "github.com/unidoc/unidoc/pdf/model"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run pdf_forms_list_fields.go input.pdf [input2.pdf] ...\n")
		os.Exit(1)
	}

	// Enable debug-level logging.
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))
	//unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelTrace))

	for _, inputPath := range os.Args[1:len(os.Args)] {
		err := listFormFields(inputPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
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

	pdfReader, err := unipdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return err
	}

	// Try decrypting with an empty one.
	if isEncrypted {
		auth, err := pdfReader.Decrypt([]byte(""))
		if err != nil {
			return err
		}
		if !auth {
			fmt.Printf(" Encrypted! Need to modify code to decrypt with your password.\n")
			return nil
		}
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
	fmt.Printf(" #Fields: %d\n", len(*acroForm.Fields))
	fmt.Printf(" =====\n")

	for idx, field := range *acroForm.Fields {
		fmt.Printf(" -Field %d (%p): %+v\n", idx+1, field, *field)
		for _, child := range field.KidsF {
			switch c := child.(type) {
			case *unipdf.PdfField:
				fmt.Printf(" --Field: %+v\n", *c)
			case *unipdf.PdfAnnotationWidget:
				fmt.Printf(" --Widget: %+v\n", *c)
			default:
				fmt.Printf(" f--UNKNOWN\n")
			}
		}
	}

	return nil
}
