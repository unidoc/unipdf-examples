/*
 * This example demonstrates how to extract text concurrently
 * with each page extraction running in its own go routine.
 * This can be useful for large documents processing.
 *
 * Run as: go run concurrent_extraction_page_level.go <input.pdf> <output_dir>
 */

package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/unidoc/unipdf/v4/common/license"
	"github.com/unidoc/unipdf/v4/extractor"
	"github.com/unidoc/unipdf/v4/model"
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
		fmt.Printf("Usage: go run concurrent_extraction_page_level.go input.pdf output_dir\n")
		os.Exit(1)
	}
	inputPdf := os.Args[1]
	outputDir := os.Args[2]

	start := time.Now()

	runPageConcurrent(inputPdf, outputDir)

	duration := time.Since(start)
	fmt.Println("time taken for concurrent extraction", duration)
}

// runPageConcurrent takes the input PDF and destination output directory and runs the extraction concurrently on page level.
func runPageConcurrent(filename string, outputDir string) {
	// Create reader.
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open input file")
		return
	}
	defer file.Close()

	if _, err := os.Stat(outputDir); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(outputDir, fs.ModePerm)
			if err != nil {
				fmt.Printf("Error: failed to create directory %s\n", outputDir)
				return
			}
		}
	}

	reader, err := model.NewPdfReader(file)
	if err != nil {
		fmt.Printf("Could not create reader")
		return
	}

	// Get total number of pages.
	numPages, err := reader.GetNumPages()
	if err != nil {
		fmt.Printf("Could not retrieve number of pages")
		return
	}

	testGoroutineExtract(reader, numPages, outputDir)
}

func testGoroutineExtract(reader *model.PdfReader, numPages int, outputDir string) {
	// Extract text from pages simultaneously.
	var wg sync.WaitGroup

	extractPage := func(pageNum int) error {
		page, err := reader.GetPage(pageNum)
		if err != nil {
			return err
		}

		ex, err := extractor.New(page)
		if err != nil {
			return err
		}

		text, err := ex.ExtractText()
		if err != nil {
			return err
		}
		// for test purposes, save each page text to its own file
		filePath := filepath.Join(outputDir, strconv.Itoa(pageNum)+".txt")
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("Error. failed to create file. %v\n", err)
		}
		_, err = file.WriteString(text)
		if err != nil {
			fmt.Printf("Error. failed to write content. %v\n", err)
		}
		return nil
	}

	for i := 1; i <= numPages; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			extractErr := extractPage(j)
			if extractErr != nil {
				fmt.Printf("extractPage error")
				return
			}
		}(i)
	}

	wg.Wait() // Wait for all goroutines to finish
}
