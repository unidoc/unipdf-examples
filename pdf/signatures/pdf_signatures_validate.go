/*
 * pdf_signatures_validate.go - Example for validating digital signatures in PDF.
 * Run as: go run pdf_signatures_validate.go input.pdf [input2.pdf ...]
 */

package main

import (
	"fmt"
	"os"

	unicommon "github.com/unidoc/unidoc/common"
	"github.com/unidoc/unidoc/pdf/model"
)

func init() {
	// Set debug-level logging.
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run pdf_signatures_validate.go input1.pdf [input2.pdf ...]\n")
		os.Exit(0)
	}

	processed := 0
	fails := 0
	for _, inputPath := range os.Args[1:] {
		fmt.Printf("Starting '%s'\n", inputPath)

		err := processPdf(inputPath)
		if err != nil {
			fmt.Printf("Error: %v - skipping file\n", err)
			fails++
		}
		processed++
	}
	fmt.Printf("Completed. %d failed /%d processed\n", fails, processed)
}

// processPdf processes a single PDF.
func processPdf(inputPath string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return err
	}
	defer f.Close()

	pdf, err := model.NewPdfReader(f)
	if err != nil {
		return err
	}

	sh := model.DefaultSignatureHandler{}
	validation, err := sh.Validate(pdf)
	if err != nil {
		fmt.Printf("Validation error: %v\n", err)
		return err
	}

	if len(validation) == 0 {
		fmt.Printf("Not signed\n")
	} else {
		allVerified := true
		allTrusted := true

		for vi, v := range validation {
			fmt.Printf("Signature %d\n", vi+1)
			fmt.Printf("  Name: %s\n", v.Name)
			fmt.Printf("  Date: %v\n", v.Date)
			fmt.Printf("  ContactInfo: %s\n", v.ContactInfo)
			fmt.Printf("  Reason: %s\n", v.Reason)
			fmt.Printf("  Location: %v\n", v.Location)
			fmt.Printf("  Num fields: %d\n", len(v.Fields))

			if v.IsVerified {
				fmt.Printf("  - Verified / ")
			} else {
				fmt.Printf("  - Not Verified / ")
				allVerified = false
			}

			if v.IsTrusted {
				fmt.Printf("Is trusted\n")
			} else {
				fmt.Printf("Not trusted\n")
				allTrusted = false
			}
		}
		fmt.Printf("All verified: %v\n", allVerified)
		fmt.Printf("All trusted: %v\n", allTrusted)
	}

	return nil
}
