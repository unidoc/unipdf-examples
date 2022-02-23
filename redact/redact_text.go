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

	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/redactor"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: go run redact_text.go inputFile.pdf outputFile.pdf \n")
		os.Exit(1)
	}

	// Make sure to enter a valid license key.
	// Otherwise text is truncated and a watermark added to the text.
	// License keys are available via: https://unidoc.io
	/*
			license.SetLicenseKey(`
		-----BEGIN UNIDOC LICENSE KEY-----
		...key contents...
		-----END UNIDOC LICENSE KEY-----
		`)
	*/

	inputFile := os.Args[1]

	outputFile := os.Args[2]

	// List of regex patterns and replacement strings

	patterns := []string{
		// Regex for matching credit card number.
		`(^|\s+)(\d{4}[ -]\d{4}[ -]\d{4}[ -]\d{4})(?:\s+|$)`,
		// Regex for matching emails.
		`[a-zA-Z0-9\.\-+_]+@[a-zA-Z0-9\.\-+_]+\.[a-z]+`,
	}
	// Replace the first matches i.e the credit cards with `*` and the emails with `#`.
	replacements := []string{"*", "#"}
	err := redactText(patterns, replacements, inputFile, outputFile)
	if err != nil {
		panic(err)
	}
	fmt.Println("successfully redacted.")
}

// redactText redacts the text in `inputFile` according to given patterns and saves result at `outputFile`.
func redactText(patterns, replacements []string, inputFile, outputFile string) error {
	terms := []redactor.RedactionTerm{}
	for i, pattern := range patterns {
		regexp, err := regexp.Compile(pattern)
		if err != nil {
			panic(err)
		}
		replacement := replacements[i]
		redTerm := redactor.RedactionTerm{Pattern: regexp, Replacement: replacement}
		terms = append(terms, redTerm)
	}
	pdfReader, f, err := model.NewPdfReaderFromFile(inputFile, nil)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// Define RedactionOptions
	options := redactor.RedactionOptions{Terms: terms}
	red, err := redactor.New(pdfReader, &options)
	if err != nil {
		panic(err)
	}
	// Execute redaction on the file.
	err = red.Redact()
	if err != nil {
		panic(err)
	}
	// Write the redacted document to outputFile.
	err = red.WriteToFile(outputFile)
	if err != nil {
		return err
	}
	return nil
}
