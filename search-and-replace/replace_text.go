/*
 * This example code shows how to do search and replace operation in PDF using unipdf
 *
 * Run as: go run replace_text.go <pattern> <replacement> <pages> <input> <output>
 *
 * example: go run replace_text.go "Australia" "America" "1,2" ./test-data/file1.pdf ./test-data/result.pdf
 */
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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
	// Ensure enough arguments are provided
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run replace_text.go <pattern> <replacement> <pages> <input> <output>")
		os.Exit(1)
	}

	// Parse positional arguments
	pattern := os.Args[1]
	replacement := os.Args[2]
	pagesArg := os.Args[3]
	filePath := os.Args[4]
	outputPath := os.Args[5]

	// Convert pages string to a slice of integers
	pageStrings := strings.Split(pagesArg, ",")
	pageList := []int{}
	for _, pageStr := range pageStrings {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			fmt.Printf("Invalid page number: %s\n", pageStr)
			os.Exit(1)
		}
		pageList = append(pageList, page)
	}

	reader, _, err := model.NewPdfReaderFromFile(filePath, nil)
	if err != nil {
		fmt.Printf("Failed to create PDF reader: %v", err)
		os.Exit(1)
	}
	editor := extractor.NewEditor(reader)

	err = editor.Replace(pattern, replacement, pageList)
	if err != nil {
		fmt.Printf("Failed to search pattern: %v\n", err)
		os.Exit(1)
	}

	err = editor.WriteToFile(outputPath)
	if err != nil {
		fmt.Printf("Failed to write to file: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Finished replacing %s by %s and saved the output file at %s\n", pattern, replacement, filePath)
}
