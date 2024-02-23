/*
* This example shows how to divide the extracted text into inside table content and
* outside table content.
* Run as : go run partition_text.go inputFile.pdf pageNum
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
		fmt.Printf("Usage: go run partition_text.go inputFile.pdf pageNum\n")
		os.Exit(1)
	}
	inputFile := os.Args[1]
	page := os.Args[2]
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		panic(err)
	}
	partitionText(inputFile, pageNum)
}

// partitionText divides the extracted text into inside table and outside table sections.
func partitionText(inputFile string, pageNum int) {
	filePath := inputFile
	fmt.Printf("working with file %s\n", filePath)
	pdfReader, f, err := model.NewPdfReaderFromFile(filePath, nil)
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
	var prevTable *extractor.TextTable
	count := 0
	indexes := [][]int{}
	index := []int{}
	for _, tMark := range pageText.Marks().Elements() {
		currentTable, _ := tMark.TableInfo()
		if prevTable == nil && currentTable != nil {
			// new table here
			index = append(index, count)
		} else if prevTable != nil && currentTable == nil {
			// end of table here
			index = append(index, count)
			indexes = append(indexes, index)
			index = []int{}
		}
		prevTable = currentTable
		count += len(tMark.Text)
	}
	beg := 0
	for i, idx := range indexes {
		idx1 := idx[0]
		idx2 := idx[1]
		fmt.Printf("\n------------------- outside table begins ---------------\n")
		fmt.Print(text[beg:idx1])
		fmt.Printf("\n------------------- outside table ends ------------------\n")
		fmt.Printf("\n------------------- inside table begins -----------------\n")
		fmt.Print(text[idx1:idx2])
		fmt.Printf("\n------------------- inside table ends -------------------\n")
		beg = idx2
		if i == len(indexes)-1 {
			fmt.Printf("\n--------------- out side table  begins --------------\n")
			fmt.Print(text[beg:])
			fmt.Printf("\n--------------- outside table ends ------------------\n")
		}
	}
}
