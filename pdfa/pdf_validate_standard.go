/*
 * PDF optimization (compression) example.
 *
 * Run as: go run pdfa_validate_standard.go <input.pdf>
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

const licenseKey = `
-----BEGIN UNIDOC LICENSE KEY-----
Free trial license keys are available at: https://unidoc.io/
-----END UNIDOC LICENSE KEY-----
`

const usage = "Usage: %s INPUT_PDF_PATH\n"

func init() {
	err := license.SetLicenseKey(licenseKey, `Company Name`)
	if err != nil {
		panic(err)
	}
}

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
		log.Fatalf("Fail: %v\n", err)
	}

	// Create reader.
	inputFile, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
	defer inputFile.Close()

	detailedReader, err := model.NewDetailedPdfReader(inputFile)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}


	// Apply standard PDF/A-1B.
	standards := []model.Standard{
		pdfa.NewProfile1A(nil),
		pdfa.NewProfile1B(nil),
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

