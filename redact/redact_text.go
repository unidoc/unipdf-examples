/*
 * Redact text: Redacts text that match given regexp patterns on a PDF document.
 *
 * Run as: go run redact_text.go input.pdf output.pdf
 */

package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/redactor"
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
		fmt.Printf("Usage: go run redact_text.go inputFile.pdf outputFile.pdf \n")
		os.Exit(1)
	}

	inputFile := os.Args[1]

	outputFile := os.Args[2]

	// List of regex patterns and replacement strings
	patterns := []string{
		// Regex for matching credit card number.
		`(^|\s+)(\d{4}[ -]\d{4}[ -]\d{4}[ -]\d{4})(?:\s+|$)`,
		// Regex for matching emails.
		`[a-zA-Z0-9\.\-+_]+@[a-zA-Z0-9\.\-+_]+\.[a-z]+`,
	}

	// Initialize the RectangleProps object.
	rectProps := &redactor.RectangleProps{
		FillColor:   creator.ColorBlack,
		BorderWidth: 0.0,
		FillOpacity: 1.0,
	}

	err := redactText(patterns, rectProps, inputFile, outputFile)
	if err != nil {
		panic(err)
	}
	fmt.Println("successfully redacted.")
}

// redactText redacts the text in `inputFile` according to given patterns and saves result at `outputFile`.
func redactText(patterns []string, rectProps *redactor.RectangleProps, inputFile, destFile string) error {

	// Initialize RedactionTerms with regex patterns.
	terms := []redactor.RedactionTerm{}
	for _, pattern := range patterns {
		regexp, err := regexp.Compile(pattern)
		if err != nil {
			panic(err)
		}
		redTerm := redactor.RedactionTerm{Pattern: regexp}
		terms = append(terms, redTerm)
	}

	pdfReader, f, err := model.NewPdfReaderFromFile(inputFile, nil)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// Define RedactionOptions.
	options := redactor.RedactionOptions{Terms: terms}
	red := redactor.New(pdfReader, &options, rectProps)
	if err != nil {
		return err
	}
	err = red.Redact()
	if err != nil {
		return err
	}
	// write the redacted document to file
	err = red.WriteToFile(destFile)
	if err != nil {
		return err
	}
	return nil
}
