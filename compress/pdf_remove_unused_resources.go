/*
 * This example shows how to use remove unused resources optimization example.
 * The compression accomplished using this filter is lossless.
 *
 * Run as: go run pdf_remove_unused_resources.go <input.pdf> <output.pdf>
 */

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/model/optimize"
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
	args := os.Args
	if len(args) < 3 {
		fmt.Printf("Usage: %s <INPUT_PDF_PATH> <OUTPUT_PDF_PATH>\n", os.Args[0])
		return
	}
	inputPath := args[1]
	outputPath := args[2]

	// Initialize starting time
	start := time.Now()
	
	inputFile, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	defer inputFile.Close()
	
	pdfReader, err := model.NewPdfReader(inputFile)
	if err != nil {
		log.Fatalf("Failed to create reader: %v", err)
	}

	pdfWriter, err := pdfReader.ToWriter(nil)
	if err != nil {
		fmt.Printf("Failed to get writer from reader %v", err)
	}

	// Set optimizer
	pdfWriter.SetOptimizer(optimize.New(optimize.Options{
		CleanUnusedResources: true,
	}))
	
	// Write document to file.
	if err := pdfWriter.WriteToFile(outputPath); err != nil {
		log.Fatalf("Failed to write to file %v", err)

	}

	// Get processing time.
	duration := float64(time.Since(start)) / float64(time.Millisecond)

	// Get input file stat.
	inputFileInfo, err := os.Stat(inputPath)
	if err != nil {
		log.Fatalf("Failed to get input file stat: %v\n", err)
	}

	// Get output file stat
	outputFileInfo, err := os.Stat(outputPath)
	if err != nil {
		log.Fatalf("Failed to get output file path: %v\n", err)
	}

	// Print information
	fmt.Printf("Input file size = %d bytes \n",inputFileInfo.Size())
	fmt.Printf("Optimized file Size = %d bytes\n", outputFileInfo.Size())
	fmt.Printf("Time taken to process %s = %.2fms\n", inputPath, duration)
}
