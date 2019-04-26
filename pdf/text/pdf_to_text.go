/*
 * PDF to text: Extract all text for each page of a pdf file.
 *
 * Run as: go run pdf_to_text.go in.pdf out.tx
 */

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/unidoc/unidoc/common"
	"github.com/unidoc/unidoc/pdf/extractor"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

const usage = "Usage: go run pdf_to_text.go in.pdf out.txt\n"

func main() {
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
	var debug, trace bool
	flag.BoolVar(&debug, "d", false, "Print debugging information.")
	flag.BoolVar(&trace, "e", false, "Print detailed debugging information.")
	makeUsage(usage)

	flag.Parse()
	args := flag.Args()

	if len(args) < 2 {
		flag.Usage()
		os.Exit(1)
	}
	if trace {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelTrace))
	} else if debug {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
	} else {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelInfo))
	}

	inputPath := args[0]
	outputPath := args[1]

	err := outputPdfText(inputPath, outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't process %q err=%v\n", inputPath, err)
		os.Exit(1)
	}
}

// outputPdfText prints out text of PDF file `inputPath` to `outputPath`.
func outputPdfText(inputPath, outputPath string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return err
	}
	numPages, err := pdfReader.GetNumPages()

	texts := []string{}
	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return err
		}
		ex, err := extractor.New(page)
		if err != nil {
			return err
		}
		text, err := ex.ExtractText()
		if err != nil {
			return err
		}
		texts = append(texts, text)
	}

	ioutil.WriteFile(outputPath, []byte(strings.Join(texts, "")), 0644)
	return nil
}

// makeUsage updates flag.Usage to include usage message `msg`.
func makeUsage(msg string) {
	usage := flag.Usage
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, msg)
		usage()
	}
}
