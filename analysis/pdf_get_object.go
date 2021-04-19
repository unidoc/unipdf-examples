/*
 * Get specific object from a PDF by number.  Prints the trailer dictionary if no number specified.
 * For streams, prints out the decoded stream.
 *
 * Run as: go run pdf_get_object.go input.pdf [num]
 */

package main

import (
	"fmt"
	"os"
	"strconv"

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
	if len(os.Args) < 2 {
		fmt.Println("Syntax: go run pdf_get_object.go input.pdf [num]")
		fmt.Println("If num is not specified, will display the trailer dictionary")
		os.Exit(1)
	}

	inputPath := os.Args[1]

	objNum := -1
	if len(os.Args) > 2 {
		num, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		objNum = num
	}

	fmt.Printf("Input file: %s\n", inputPath)
	err := inspectPdfObject(inputPath, objNum)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func inspectPdfObject(inputPath string, objNum int) error {
	readerOpts := model.NewReaderOpts()
	readerOpts.LazyLoad = false

	pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, readerOpts)
	if err != nil {
		return err
	}
	defer f.Close()

	// Print trailer
	if objNum == -1 {
		trailer, err := pdfReader.GetTrailer()
		if err != nil {
			return err
		}

		fmt.Printf("Trailer: %s\n", trailer.String())
		return nil
	}

	obj, err := pdfReader.GetIndirectObjectByNumber(objNum)
	if err != nil {
		return err
	}

	fmt.Printf("Object %d: %s\n", objNum, obj.String())

	if stream, is := obj.(*core.PdfObjectStream); is {
		decoded, err := core.DecodeStream(stream)
		if err != nil {
			return err
		}
		fmt.Printf("Decoded:\n%s", decoded)
	} else if indObj, is := obj.(*core.PdfIndirectObject); is {
		fmt.Printf("%T\n", indObj.PdfObject)
		fmt.Printf("%s\n", indObj.PdfObject.String())
	}

	return nil
}
