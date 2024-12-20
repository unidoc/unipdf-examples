package main

import (
	"fmt"
	"os"

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
	pattern := "Australia"
	pages := []int{1}
	replacement := "America"
	filePath := "./test-data/file1.pdf"

	outputPath := "./test-data/result.pdf"
	reader, _, err := model.NewPdfReaderFromFile(filePath, nil)
	if err != nil {
		fmt.Printf("Failed to create PDF reader: %v", err)
		os.Exit(1)
	}
	editor := extractor.NewEditor(reader)

	err = editor.Replace(pattern, replacement, pages)
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
