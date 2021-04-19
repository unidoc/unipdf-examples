/*
 * Basic example for text searching: Retrieving position of a signature line in PDF where the signature line is given by
 * "__________________" text. And positioned with a Tm operation above.
 *
 * Run as: go run pdf_detect_signature.go input.pdf
 */

package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/contentstream"
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
		fmt.Printf("Usage: go run pdf_detect_signature.go input.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]

	err := detectSignatureInput(inputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func detectSignatureInput(inputPath string) error {
	readerOpts := model.NewReaderOpts()
	readerOpts.LazyLoad = false

	pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, readerOpts)
	if err != nil {
		return err
	}
	defer f.Close()

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return err
		}

		if pageNum == 1 {
			found, x, y, err := locateSignatureLine(page)
			if err != nil {
				return err
			}
			// Happens if did not find the "___.." line or if there was no Tm position marker before.
			if !found || (x == 0 && y == 0) {
				return errors.New("Unable to find the signature line")
			}

			fmt.Printf("Position: x: %f, y: %f\n", x, y)
		}

	}

	return nil
}

func locateSignatureLine(page *model.PdfPage) (bool, float64, float64, error) {
	found := false
	x := float64(0)
	y := float64(0)

	pageContentStr, err := page.GetAllContentStreams()
	if err != nil {
		return found, x, y, err
	}

	cstreamParser := contentstream.NewContentStreamParser(pageContentStr)
	if err != nil {
		return found, x, y, err
	}

	operations, err := cstreamParser.Parse()
	if err != nil {
		return found, x, y, err
	}

	for _, op := range *operations {
		switch {
		case op.Operand == "Tm" && len(op.Params) == 6:
			if val, has := core.GetFloatVal(op.Params[4]); has {
				x = val
			}

			if val, has := core.GetFloatVal(op.Params[5]); has {
				y = val
			}

		case op.Operand == "Tj" && len(op.Params) == 1:
			str, isStr := core.GetStringVal(op.Params[0])
			if isStr {
				if strings.Contains(str, "________________") {
					fmt.Printf("Tj: %s\n", str)
					found = true
					break
				}
			}
		}
	}

	return found, x, y, nil
}
