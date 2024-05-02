/*
 * PDF/A-3 optimization (compression) example.
 *
 * Run as: go run pdfa3_validate_standard.go <input.pdf>
 */

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/model/pdfa"
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
	if len(args) < 2 {
		fmt.Printf("Usage: %s INPUT_PDF_PATH", os.Args[0])
		return
	}
	inputPath := args[1]

	// Initialize starting time.
	start := time.Now()

	// Create reader.
	inputFile, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
	defer inputFile.Close()

	detailedReader, err := model.NewCompliancePdfReader(inputFile)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Apply standard PDF/A-3.
	standards := []model.StandardImplementer{
		pdfa.NewProfile3A(nil),
		pdfa.NewProfile3B(nil),
		pdfa.NewProfile3U(nil),
	}

	// Iterate over input standards and check if the document passes its requirements.
	for _, standard := range standards {
		if err = standard.ValidateStandard(detailedReader); err != nil {
			fmt.Printf("Input document didn't pass the standard: %s - %v\n", standard.StandardName(), err)
		}
	}

	duration := float64(time.Since(start)) / float64(time.Millisecond)
	fmt.Printf("Processing time: %.2f ms\n", duration)
}
