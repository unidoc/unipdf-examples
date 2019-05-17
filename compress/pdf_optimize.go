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

	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/model/optimize"
)

const usage = "Usage: %s INPUT_PDF_PATH OUTPUT_PDF_PATH\n"

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Printf(usage, os.Args[0])
		return
	}
	inputPath := args[1]
	outputPath := args[2]

	// Initialize starting time.
	start := time.Now()

	// Get input file stat.
	inputFileInfo, err := os.Stat(inputPath)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	// Create reader.
	inputFile, err := os.Open(inputPath)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}
	defer inputFile.Close()

	reader, err := model.NewPdfReader(inputFile)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	// Get number of pages in the input file.
	pages, err := reader.GetNumPages()
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	// Add input file pages to the writer.
	writer := model.NewPdfWriter()
	for i := 1; i <= pages; i++ {
		page, err := reader.GetPage(i)
		if err != nil {
			log.Fatal("Fail: %v\n", err)
		}

		if err = writer.AddPage(page); err != nil {
			log.Fatal("Fail: %v\n", err)
		}
	}

	// Add reader AcroForm to the writer.
	if reader.AcroForm != nil {
		writer.SetForms(reader.AcroForm)
	}

	// Set optimizer.
	writer.SetOptimizer(optimize.New(optimize.Options{
		CombineDuplicateDirectObjects:   true,
		CombineIdenticalIndirectObjects: true,
		CombineDuplicateStreams:         true,
		CompressStreams:                 true,
		UseObjectStreams:                true,
		ImageQuality:                    80,
		ImageUpperPPI:                   100,
	}))

	// Create output file.
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}
	defer outputFile.Close()

	// Write output file.
	err = writer.Write(outputFile)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	// Get output file stat.
	outputFileInfo, err := os.Stat(outputPath)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
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
