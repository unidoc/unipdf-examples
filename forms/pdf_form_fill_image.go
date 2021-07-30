/*
 * Example of attaching an image to a push button field.
*
* Run as: go run pdf_form_fill_image.go input.pdf output.pdf.
*/

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/fjson"
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

// Example of filling PDF formdata with a form.
func main() {
	var (
		inputPath  string
		outputPath string
	)
	inputPath = os.Args[1]
	if len(os.Args) > 2 {
		outputPath = os.Args[2]
	}

	// Output path not specified: Export list of fields and data as JSON format.
	if len(outputPath) == 0 {
		fdata, err := fjson.LoadFromPDFFile(inputPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		if fdata == nil {
			fmt.Printf("No data\n")
			return
		}
		fjson, err := fdata.JSON()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", fjson)
		return
	}

	err := fillImageFields(inputPath, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Success, output written to %s\n", outputPath)
}

func fillImageFields(inputPath, outputPath string) error {
	fdata := &fjson.FieldData{}
	fdata.SetImageFromFile("01_Logo Emitente", "icon.png", nil)

	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return err
	}

	fieldAppearance := annotator.ImageFieldAppearance{OnlyIfMissing: true}

	// Populate the form data.
	err = pdfReader.AcroForm.FillWithAppearance(fdata, fieldAppearance)
	if err != nil {
		return err
	}

	// Generate a PdfWriter instance from existing PdfReader.
	pdfWriter, err := pdfReader.ToWriter(nil)
	if err != nil {
		return err
	}

	// Write to file.
	err = pdfWriter.WriteToFile(outputPath)
	return err
}
