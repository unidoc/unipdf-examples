/*
 * This example shows how to remove unused resource from pdf pages while splitting.
 *
 * Run as: go run pdf-split-with-unused-resources.go <input.pdf> <output-dir>
 * In this example the splitting is done one page per document but it can be easily extended to any kind of splitting.
 */

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/contentstream"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
)

// func init() {
// 	// Make sure to load your metered License API key prior to using the library.
// 	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
// 	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
// 	if err != nil {
// 		panic(err)
// 	}
// }

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: go run pdf-split-with-unused-resources.go <input.pdf> <output-dir>\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]

	outputDir := os.Args[2]
	splitPdfFile(inputPath, outputDir)
}
func splitPdfFile(inputFile string, outputDir string) {
	pdfBytes, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	r := bytes.NewReader(pdfBytes)
	pdfReader, err := model.NewPdfReader(r)
	if err != nil {
		log.Printf("Failed to create reader: %v", err)
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		log.Printf("Failed to get number of pages: %v", err)
	}
	for pageIdx := 0; pageIdx < numPages; pageIdx++ {
		pdfPage, _ := pdfReader.GetPage(pageIdx + 1)
		cleanUnusedXobjects(pdfPage)
		w := model.NewPdfWriter()
		if err := w.AddPage(pdfPage); err != nil {
			log.Fatalf("Failed to convert pageIdx = %d, %v", pageIdx, err)
		}
		if err := w.WriteToFile(fmt.Sprintf("./%s/output/doc_page_%d_new.pdf", outputDir, pageIdx)); err != nil {
			log.Fatalf("Failed to write to file %d: %v", pageIdx, err)
		}
	}
}

func cleanUnusedXobjects(page *model.PdfPage) {
	contents, err := page.GetAllContentStreams()
	if err != nil {
		common.Log.Debug("failed to get page content stream")
	}
	parser := contentstream.NewContentStreamParser(contents)
	operations, err := parser.Parse()
	if err != nil {
		common.Log.Debug("failed to parse content stream")
	}
	usedObjectsNames := []string{}
	for _, op := range *operations {
		operand := op.Operand
		// Draw xObject operator.
		if operand == "Do" {
			params := op.Params
			imageName := params[0].String()
			usedObjectsNames = append(usedObjectsNames, imageName)
		}
	}

	xObject := page.Resources.XObject
	dict, ok := xObject.(*core.PdfObjectDictionary)
	if ok {
		keys := getKeys(dict)
		for _, k := range keys {
			if exists(k, usedObjectsNames) {
				continue
			}
			name := *core.MakeName(k)
			dict.Remove(name)
		}
	}

}

func getKeys(dict *core.PdfObjectDictionary) []string {
	keys := []string{}
	for _, k := range dict.Keys() {
		keys = append(keys, k.String())
	}
	return keys
}
func exists(element string, elements []string) bool {
	for _, el := range elements {
		if element == el {
			return true
		}
	}
	return false
}
