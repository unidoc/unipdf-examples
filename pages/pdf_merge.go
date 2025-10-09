/*
 * Basic merging of PDF files.
 * Simply loads all pages for each file and writes to the output file.
 * See pdf_merge_advanced.go for a more advanced version which handles merging document forms (acro forms) also.
 *
 * Run as: go run pdf_merge.go output.pdf input1.pdf input2.pdf input3.pdf ...
 */

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v4/common/license"
	"github.com/unidoc/unipdf/v4/pdfutil"
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
	if len(os.Args) < 4 {
		fmt.Printf("Requires at least 3 arguments: output_path and 2 input paths\n")
		fmt.Printf("Usage: go run pdf_merge.go output.pdf input1.pdf input2.pdf input3.pdf ...\n")
		os.Exit(0)
	}

	outputPath := ""
	inputPaths := []string{}

	// Sanity check the input arguments.
	for i, arg := range os.Args {
		if i == 0 {
			continue
		} else if i == 1 {
			outputPath = arg
			continue
		}

		inputPaths = append(inputPaths, arg)
	}

	// Merge pdfs without forms.
	err := pdfutil.MergePdf(inputPaths, outputPath, false)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}
