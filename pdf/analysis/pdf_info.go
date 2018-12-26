/*
 * Prints basic PDF info: number of pages and encryption status.
 *
 * Run as: go run pdf_info.go input1.pdf [input2.pdf] ...
 */

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/unidoc/unidoc/common"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

const usage = "Print out basic properties of PDF files\n" +
	"Usage: go run pdf_info.go input.pdf [input2.pdf] ...\n"

type PdfProperties struct {
	Version     string
	IsEncrypted bool
	CanView     bool // Is the document viewable without password?
	NumPages    int
}

func main() {
	var showHelp, debug, trace bool
	flag.BoolVar(&showHelp, "h", false, "Show this help message.")
	flag.BoolVar(&debug, "d", false, "Print debugging information.")
	flag.BoolVar(&trace, "e", false, "Print detailed debugging information.")

	if len(os.Args) < 2 {
		fmt.Printf("Print out basic properties of PDF files\n")
		fmt.Printf("Usage: go run pdf_info.go input.pdf [input2.pdf] ...\n")
		os.Exit(1)
	}
	makeUsage(usage)

	flag.Parse()
	args := flag.Args()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	if trace {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelTrace))
	} else if debug {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
	} else {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelError))
	}

	for _, inputPath := range args {
		fmt.Printf("Input file: %s\n", inputPath)

		ret, err := getPdfProperties(inputPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf(" PDF Version: %s\n", ret.Version)
		fmt.Printf(" Num Pages: %d\n", ret.NumPages)
		fmt.Printf(" Is Encrypted: %t\n", ret.IsEncrypted)
		fmt.Printf(" Is Viewable (without pass): %t\n", ret.CanView)
	}
}

func getPdfProperties(inputPath string) (*PdfProperties, error) {
	ret := PdfProperties{}

	f, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return nil, err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return nil, err
	}

	ret.IsEncrypted = isEncrypted
	ret.CanView = true

	// Try decrypting with an empty one.
	if isEncrypted {
		auth, err := pdfReader.Decrypt([]byte(""))
		if err != nil {
			return nil, err
		}
		ret.CanView = auth
		return &ret, nil
	}

	ret.Version = pdfReader.PdfVersion().String()

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return nil, err
	}
	ret.NumPages = numPages

	return &ret, nil
}

// makeUsage updates flag.Usage to include usage message `msg`.
func makeUsage(msg string) {
	usage := flag.Usage
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, msg)
		usage()
	}
}
