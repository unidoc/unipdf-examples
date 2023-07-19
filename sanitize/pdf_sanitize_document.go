package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/sanitize"
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
		fmt.Printf("Usage: pdf_sanitize_document.go <INPUT_PDF_PATH> <OUTPUT_PDF_PATH>\n", os.Args[0])
		return
	}
	inputPath := args[1]
	outputPath := args[2]

	start := time.Now()
	pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, nil)
	if err != nil {
		log.Fatalf("Failed to create reader: %v", err)
	}
	defer f.Close()

	pdfWriter, err := pdfReader.ToWriter(nil)
	if err != nil {
		log.Fatalf("Failed to get writer object: %v", err)
	}

	opts := sanitize.SanitizationOpts{
		JavaScript:  true,
		URI:         true,
		GoToR:       true,
		GoTo:        true,
		RenditionJS: true,
		OpenAction:  true,
		Launch:      true,
	}

	pdfWriter.SetOptimizer(sanitize.New(opts))
	pdfWriter.WriteToFile(outputPath)

	duration := float64(time.Since(start)) / float64(time.Millisecond)
	inputFileInfo, err := os.Stat(inputPath)
	if err != nil {
		fmt.Printf("Failed to get inputFile info %v", err)
	}

	outputFileInfo, err := os.Stat(outputPath)
	if err != nil {
		fmt.Printf("Failed to get outputFile info %v", err)
	}

	fmt.Printf("Input file size %d bytes\n", inputFileInfo.Size())
	fmt.Printf("Output file size %d bytest\n",outputFileInfo.Size())
	fmt.Printf("Processing time %.2f ms",duration)

}
