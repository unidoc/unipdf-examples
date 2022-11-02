/*
 * An example code to extract bullet point list from pdf document.
 *
 * Run as: go run extract_bullet_list.go input.pdf 1 output.txt
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
	if len(os.Args) < 4 {
		fmt.Printf("Usage: go run extract_bullet_list.go inputFile.pdf 1 output.txt \n")
		os.Exit(1)
	}
	inputFile := os.Args[1]
	p := os.Args[2]
	output := os.Args[3]
	page, err := strconv.Atoi(p)
	if err != nil {
		fmt.Printf("wrong page value.\nUsage: go run extract_bullet_list.go inputFile.pdf 1 output.txt\n")
		os.Exit(1)
	}
	fout, err := os.Create(output)
	if err != nil {
		fmt.Printf("failed: %s", err)
		os.Exit(1)
	}
	defer fout.Close()
	txt := extractBulletList(inputFile, page)
	_, err = fout.WriteString(txt)
	if err != nil {
		fmt.Printf("failed to write: %s ", err)
		os.Exit(1)
	}
	fmt.Println(txt)
}

// extractBulletList extracts bullet point lists on `file` at `pageNum` and returns the string representation of the list.
func extractBulletList(file string, pageNum int) string {
	pdfReader, f, err := model.NewPdfReaderFromFile(file, nil)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	page, err := pdfReader.GetPage(pageNum)
	if err != nil {
		panic(err)
	}
	options := &extractor.Options{
		DisableDocumentTags: false,
	}
	ex, err := extractor.NewWithOptions(page, options)
	if err != nil {
		panic(err)
	}
	pageText, _, _, err := ex.ExtractPageText()
	if err != nil {
		panic(err)
	}
	lists := pageText.List()
	txt := lists.Text()
	return txt
}
