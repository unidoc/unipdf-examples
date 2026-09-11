/*
 * Redact text: Redacts text that match given regexp patterns on a PDF document.
 *
 * Run as: go run redact_text.go input.pdf output.pdf
 */

package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/unidoc/unipdf/v5/common/license"
	"github.com/unidoc/unipdf/v5/creator"
	"github.com/unidoc/unipdf/v5/model"
	"github.com/unidoc/unipdf/v5/redactor"
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
	// Verification is enabled by default: after redaction, the redactor checks that every
	// matched term was physically removed from the page and form XObject content streams.
	// Set DisableVerification to true to skip this check.
	options := redactor.RedactionOptions{Terms: terms, DisableVerification: false}
	red := redactor.New(pdfReader, &options, rectProps)
	if err != nil {
		return err
	}
	err = red.Redact()
	if err != nil {
		// An IncompleteRedactionError means the document was processed, but some matches
		// could not be verified as removed. The output must not be treated as safely redacted.
		var incomplete *redactor.IncompleteRedactionError
		if errors.As(err, &incomplete) {
			for _, failure := range incomplete.Failures {
				fmt.Printf("Unverified redaction on page %d, term %q: %s\n", failure.Page, failure.Term, failure.Reason)
			}
		}
		return err
	}
	// write the redacted document to destFile.
	err = red.WriteToFile(destFile)
	if err != nil {
		return err
	}
	return nil
}
