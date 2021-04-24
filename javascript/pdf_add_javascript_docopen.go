/*
 * Add a Document Open Javascript Action to a PDF.
 */

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/common/license"
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
	if len(os.Args) < 3 {
		fmt.Printf("Usage:  go run pdf_add_js_docopen.go input.pdf output.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	// Most basic example - print out "Hello World" dialog.
	// For a more information on what can be done, see:
	// - Acrobat Javascript Scripting Guide
	//   https://www.adobe.com/content/dam/acom/en/devnet/acrobat/pdfs/acrojsguide.pdf
	// - Javscript for Acrobat API Reference
	//   https://www.adobe.com/content/dam/acom/en/devnet/acrobat/pdfs/js_api_reference.pdf
	docOpenJS := `app.alert('Hello world!');`

	err := addPDFDocOpenJSAction(inputPath, outputPath, docOpenJS)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func addPDFDocOpenJSAction(inputPath, outputPath string, docOpenJS string) error {
	pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, nil)
	if err != nil {
		return err
	}
	defer f.Close()

	w, err := pdfReader.ToWriter(nil)
	if err != nil {
		return err
	}
	if err := addJSOpenAction(w, docOpenJS); err != nil {
		return err
	}

	return w.WriteToFile(outputPath)
}

func addJSOpenAction(w *model.PdfWriter, docOpenJS string) error {
	action := model.NewPdfActionJavaScript()
	action.JS = core.MakeEncodedString(docOpenJS, true)

	name := "javascript-action-name"
	nameTree := core.MakeDictMap(map[string]core.PdfObject{
		"JavaScript": core.MakeDictMap(map[string]core.PdfObject{
			"Names": core.MakeArray(core.MakeString(name), action.ToPdfObject()),
		}),
	})

	if err := w.SetNamedDestinations(nameTree); err != nil {
		return err
	}
	if err := w.SetOpenAction(action.ToPdfObject()); err != nil {
		return err
	}

	return nil
}
