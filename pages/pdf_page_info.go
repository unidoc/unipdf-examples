/*
 * Prints PDF page info: Mediabox size and other parameters.
 * If [page num] is not specified prints out info for all pages.
 *
 * Run as: go run pdf_info.go input.pdf [page num]
 */

package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/model"
	pdf "github.com/unidoc/unipdf/v3/model"
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
		fmt.Printf("Usage:  go run pdf_info.go input.pdf [page num]\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]

	pageNum := 0
	if len(os.Args) > 2 {
		num, err := strconv.ParseInt(os.Args[2], 10, 64)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		pageNum = int(num)
	}

	fmt.Printf("Input file: %s\n", inputPath)

	err := printPdfPageProperties(inputPath, pageNum)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func printPdfPageProperties(inputPath string, pageNum int) error {
	pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, nil)
	if err != nil {
		return err
	}
	defer f.Close()

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	// If invalid pagenum, print all pages.
	if pageNum <= 0 || pageNum > numPages {
		for i := 0; i < numPages; i++ {
			page, err := pdfReader.GetPage(i + 1)
			if err != nil {
				return err
			}
			fmt.Printf("-- Page %d\n", i+1)
			err = processPage(page)
			if err != nil {
				return err
			}
		}
	} else {
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return err
		}
		fmt.Printf("-- Page %d\n", pageNum)
		err = processPage(page)
		if err != nil {
			return err
		}
	}

	return nil
}

func processPage(page *pdf.PdfPage) error {
	pageWidth, pageHeight, err := page.Size()
	if err != nil {
		return err
	}

	fmt.Printf(" Page: %+v\n", page)
	if page.Rotate != nil {
		fmt.Printf(" Page rotation: %v\n", *page.Rotate)
	} else {
		fmt.Printf(" Page rotation: 0\n")
	}
	fmt.Printf(" Page mediabox: %+v\n", page.MediaBox)
	fmt.Printf(" Page height: %f\n", pageHeight)
	fmt.Printf(" Page width: %f\n", pageWidth)

	return nil
}
