/*
 * Initiate empty PDF Writer, clean the Annots from pages and show analysis for newly created PDF.
 *
 * Run as: go run pdf_modify_contents.go input.pdf output.pdf
 */

package main

import (
	"bytes"
	"flag"
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

type cmdOptions struct {
	pdfPassword string
}

func main() {
	var opt cmdOptions
	flag.StringVar(&opt.pdfPassword, "password", "", "PDF Password (empty default)")
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Syntax: go run main.go [options] input.pdf output.pdf")
		os.Exit(1)
	}

	inputPath := args[0]
	outputPath := args[1]
	fmt.Printf("Input file: %s\n", inputPath)

	readerOpts := model.NewReaderOpts()
	readerOpts.Password = opt.pdfPassword

	pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, readerOpts)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	pdfWriter := model.NewPdfWriter()

	// Removes annotations from each pages.
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		panic(err)
	}
	for i := 0; i < numPages; i++ {
		page, err := pdfReader.GetPage(i + 1)
		if err != nil {
			panic(err)
		}

		if page.Annots != nil {
			fmt.Printf("-- Removing Annots on Page %d\n", i+1)
			page.Annots = nil
		}

		err = pdfWriter.AddPage(page)
		if err != nil {
			panic(err)
		}
	}

	outDoc := bytes.NewBuffer(nil)
	err = pdfWriter.Write(outDoc)
	if err != nil {
		panic(err)
	}

	pdfReader, err = model.NewPdfReader(bytes.NewReader(outDoc.Bytes()))
	if err != nil {
		panic(err)
	}

	err = inspectPdf(pdfReader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(outputPath, outDoc.Bytes(), 0777)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Output file: %s\n", outputPath)
}

func inspectPdf(pdfReader *model.PdfReader) error {
	objNums := pdfReader.GetObjectNums()

	// Output.
	fmt.Printf("%d PDF objects:\n", len(objNums))
	for i, objNum := range objNums {
		obj, err := pdfReader.GetIndirectObjectByNumber(objNum)
		if err != nil {
			return err
		}
		fmt.Println("=========================================================")
		fmt.Printf("%3d: %d 0 %T\n", i, objNum, obj)
		if stream, is := obj.(*core.PdfObjectStream); is {
			decoded, err := core.DecodeStream(stream)
			if err != nil {
				return err
			}
			fmt.Printf("Decoded:\n%s\n", decoded)
		} else if indObj, is := obj.(*core.PdfIndirectObject); is {
			fmt.Printf("%T\n", indObj.PdfObject)
			fmt.Printf("%s\n", indObj.PdfObject.String())
		}
	}

	return nil
}
