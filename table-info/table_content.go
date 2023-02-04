/*
* This example shows how to calculate the percentage of page content that is inside a table
* using the TextMark.TableInfo() method.
* Run as: go run table_content.go inputFile.pdf pageNum
*/

package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/extractor"
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
		fmt.Printf("Usage: go run table_content.go inputFile.pdf pageNum\n")
		os.Exit(1)
	}
	inputFile := os.Args[1]
	page := os.Args[2]
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		panic(err)
	}
	checkContentDistribution(inputFile, pageNum)
}

func checkContentDistribution(inputFile string, pageNum int) {
	fmt.Printf("running with %s\n", inputFile)
	pdfReader, f, err := model.NewPdfReaderFromFile(inputFile, nil)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	page, err := pdfReader.GetPage(pageNum)
	if err != nil {
		panic(err)
	}
	ex, err := extractor.New(page)
	if err != nil {
		panic(err)
	}
	pageText, _, _, err := ex.ExtractPageText()
	if err != nil {
		panic(err)
	}

	text := pageText.Text()
	length := len(text)
	count := 0
	for _, textMark := range pageText.Marks().Elements() {
		table, _ := textMark.TableInfo()
		if table != nil {
			count += len(textMark.Text)
		}
	}
	var distribution float64
	if length != 0 {
		distribution = float64(count) / float64(length)
	}
	fmt.Printf("\n %.3f percent of the page content is inside a table \n", distribution)
}
