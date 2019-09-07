/*
 * pdf_search_replace.go - Basic example of find and replace with UniDoc.
 * Replaces <text> with <replace text> in the output PDF.
 *
 * Syntax: go run pdf_search_replace.go <input.pdf> <output.pdf> <text> <replace text>
 */

package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/contentstream"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/model/optimize"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Printf("Usage: go run pdf_search_replace.go <input.pdf> <output.pdf> <text> <replace text>\n")
		os.Exit(0)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]
	searchText := os.Args[3]
	replaceText := os.Args[4]

	err := searchReplace(inputPath, outputPath, searchText, replaceText)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully created %s\n", outputPath)
}

func searchReplace(inputPath, outputPath, searchText, replaceText string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()
	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return err
	}

	pdfWriter := model.NewPdfWriter()

	encrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return err
	}
	if encrypted {
		ok, err := pdfReader.Decrypt([]byte(""))
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("Encrypted")
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	for n := 1; n <= numPages; n++ {
		page, err := pdfReader.GetPage(n)
		if err != nil {
			return err
		}

		err = searchReplacePageText(page, searchText, replaceText)
		if err != nil {
			return err
		}

		err = pdfWriter.AddPage(page)
		if err != nil {
			return err
		}
	}

	fw, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer fw.Close()

	opt := optimize.Options{
		CombineDuplicateStreams:         true,
		CombineIdenticalIndirectObjects: true,
		UseObjectStreams:                true,
		CompressStreams:                 true,
	}
	pdfWriter.SetOptimizer(optimize.New(opt))

	return pdfWriter.Write(fw)
}

func searchReplacePageText(page *model.PdfPage, searchText, replaceText string) error {
	contents, err := page.GetAllContentStreams()
	if err != nil {
		return err
	}

	csParser := contentstream.NewContentStreamParser(contents)
	ops, err := csParser.Parse()
	if err != nil {
		return err
	}

	replaceFunc := func(objptr *core.PdfObject) {
		strobj, ok := core.GetString(*objptr)
		if !ok {
			common.Log.Debug("Invalid parameter, skipping")
			return
		}

		if strings.Contains(strobj.String(), searchText) {
			s := strings.Replace(strobj.String(), searchText, replaceText, -1)
			*strobj = *core.MakeString(s)
		}
	}

	processor := contentstream.NewContentStreamProcessor(*ops)
	processor.AddHandler(contentstream.HandlerConditionEnumAllOperands, "",
		func(op *contentstream.ContentStreamOperation, gs contentstream.GraphicsState, resources *model.PdfPageResources) error {
			switch op.Operand {
			case `Tj`, `'`:
				if len(op.Params) != 1 {
					common.Log.Debug("Invalid: Tj/' with invalid set of parameters - skip")
					return nil
				}
				replaceFunc(&op.Params[0])
			case `''`:
				if len(op.Params) != 3 {
					common.Log.Debug("Invalid: '' with invalid set of parameters - skip")
					return nil
				}
				replaceFunc(&op.Params[3])
			case `TJ`:
				if len(op.Params) != 1 {
					common.Log.Debug("Invalid: TJ with invalid set of parameters - skip")
					return nil
				}
				arr, _ := core.GetArray(op.Params[0])
				for i := range arr.Elements() {
					obj := arr.Get(i)
					replaceFunc(&obj)
					arr.Set(i, obj)
				}
			}

			return nil
		})

	err = processor.Process(page.Resources)
	if err != nil {
		return err
	}

	return page.SetContentStreams([]string{ops.String()}, core.NewFlateEncoder())
}
