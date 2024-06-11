/*
 * This example demonstrates how to extract content form multiple documents concurrently
 * with each document extraction running in its own go routine.
 * N.B. currently concurrency is supported on a document level which means we can only extract one document per go routine.
 *
 * Run as: go run concurrent_extraction.go <input1.pdf> <input2.pdf> <input3.pdf>... <output_dir>
 */

package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

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
		fmt.Printf("Usage: go run concurrent_extraction.go input1.pdf input2.pdf input3.pdf... output_dir\n")
		os.Exit(1)
	}
	inputDocuments := []string{}
	args := os.Args
	var outputDir string
	for i := 1; i < len(args); i++ {
		if i == len(args)-1 {
			outputDir = args[i]
		} else {
			inputDocuments = append(inputDocuments, args[i])
		}
	}

	start := time.Now()
	runConcurrent(inputDocuments, outputDir)
	duration := time.Since(start)
	fmt.Println("time taken for concurrent extraction", duration)
}

// runConcurrent takes the list of input documents and destination output directory and runs the extraction concurrently.
func runConcurrent(documents []string, outputDir string) {
	res := make(chan map[string]string, len(documents))

	err := concurrentExtraction(documents, res)
	if err != nil {
		fmt.Printf("Error. extraction failed %v.\n", err)
	}
	outputPath := outputDir
	if _, err := os.Stat(outputPath); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(outputPath, fs.ModePerm)
			if err != nil {
				fmt.Printf("Error. failed to create directory %s\n", outputPath)
			}
		}
	}

	for i := 0; i < len(documents); i++ {
		result := <-res
		for path, content := range result {
			basename := filepath.Base(path)
			fileName := strings.TrimSuffix(basename, filepath.Ext(basename)) + ".txt"
			filePath := filepath.Join(outputPath, fileName)
			file, err := os.Create(filePath)
			if err != nil {
				fmt.Printf("Error. failed to create file. %v\n", err)
			}
			_, err = file.WriteString(content)
			if err != nil {
				fmt.Printf("Error. failed to write content. %v\n", err)
			}
		}
	}
}

// concurrentExtraction launches a go routine for each document in `documents` and writes the result of extraction to
// the channel `res`.
func concurrentExtraction(documents []string, res chan map[string]string) error {
	for _, docPath := range documents {
		filePath := docPath
		go func(path string, res chan map[string]string) {
			result, err := extractSingleDoc(path)
			if err != nil {
				fmt.Printf("Error. Failed to extract file %v. %v\n", filePath, err)
			}
			temp := map[string]string{
				filePath: result,
			}
			res <- temp
		}(filePath, res)
	}

	return nil
}

// extractSingleDoc takes a single file specified by the `filePath` and returns the extracted text.
func extractSingleDoc(filePath string) (string, error) {
	pdfReader, _, err := model.NewPdfReaderFromFile(filePath, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create pdf reader: %w", err)
	}
	pages, err := pdfReader.GetNumPages()
	if err != nil {
		return "", err
	}
	var result string
	for i := 0; i < pages; i++ {
		pageNum := i + 1
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return "", fmt.Errorf("failed to get page %d: %w", pageNum, err)
		}

		ex, err := extractor.New(page)
		if err != nil {
			return "", fmt.Errorf("failed to create extractor: %w", err)
		}

		text, err := ex.ExtractText()
		if err != nil {
			return "", fmt.Errorf("failed to extract text: %w", err)
		}
		result += text
	}

	return result, nil
}
