/*
 * PDF optimization (compression) example.
 *
 * Run as: go run pdf_optimize.go <input.pdf> <output.pdf>
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
		fmt.Printf("Usage: %s INPUT_PDF_PATH OUTPUT_PDF_PATH\n", os.Args[0])
		return
	}
	inputPath := args[1]
	outputPath := args[2]

	// Initialize starting time.
	start := time.Now()

	// Get input file stat.
	inputFileInfo, err := os.Stat(inputPath)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Create reader.
	inputFile, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
	defer inputFile.Close()

	reader, err := model.NewPdfReader(inputFile)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Generate a PDFWriter from PDFReader.
	pdfWriter, err := reader.ToWriter(nil)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Set optimizer.
	pdfWriter.SetOptimizer(optimize.New(optimize.Options{
		CombineDuplicateDirectObjects:   true,
		CombineIdenticalIndirectObjects: true,
		CombineDuplicateStreams:         true,
		CompressStreams:                 true,
		UseObjectStreams:                true,
		ImageQuality:                    80,
		ImageUpperPPI:                   100,
		CleanUnusedResources:            true,
	}))

	// Create output file.
	err = pdfWriter.WriteToFile(outputPath)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Get output file stat.
	outputFileInfo, err := os.Stat(outputPath)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Print basic optimization statistics.
	inputSize := inputFileInfo.Size()
	outputSize := outputFileInfo.Size()
	ratio := 100.0 - (float64(outputSize) / float64(inputSize) * 100.0)
	duration := float64(time.Since(start)) / float64(time.Millisecond)

	fmt.Printf("Original file: %s\n", inputPath)
	fmt.Printf("Original size: %d bytes\n", inputSize)
	fmt.Printf("Optimized file: %s\n", outputPath)
	fmt.Printf("Optimized size: %d bytes\n", outputSize)
	fmt.Printf("Compression ratio: %.2f%%\n", ratio)
	fmt.Printf("Processing time: %.2f ms\n", duration)
}
