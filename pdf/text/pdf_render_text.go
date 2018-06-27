/*
 * PDF to text: Extract all text for each page of a pdf file.
 *
 * Run as: go run pdf_extract_text.go input.pdf
 */

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"github.com/unidoc/unidoc/common"
	pdfcore "github.com/unidoc/unidoc/pdf/core"
	"github.com/unidoc/unidoc/pdf/extractor"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run pdf_render_text.go input.pdf\n")
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

	// For debugging.
	common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
	files := os.Args[1:]
	sort.Strings(files)

	exclusions := map[string]bool{
		`The-Byzantine-Generals-Problem.pdf`: true,
		`endosymbiotictheory_marguli.pdf`:    true,
		`iverson.pdf`:                        true,
		`p253-porter.pdf`:                    true,
		`shamirturing.pdf`:                   true,
		`warnock_camelot.pdf`:                true,
		`B02.pdf`:                            true,

		// [DEBUG]  text.go:617 getFont: NewPdfFontFromPdfObject failed. name=`C0_0` err=Bad state
		// [DEBUG]  processor.go:279 Processor handler error: Bad state
		// [ERROR]  text.go:212 Error processing: Bad state
		`Data Classification For Dummies_Identity Finder_Special Edition_Todd_Feinman.pdf`: true,

		// 0x29 ')' missing from encoding
		`000018.pdf`: true,
		// 0x3c '<' missing from encoding
		`03-block-v2-annotated.pdf`: true,
	}
	files2 := []string{}
	for _, inputPath := range files {
		if _, ok := exclusions[filepath.Base(inputPath)]; ok {
			continue
		}
		if strings.Contains(inputPath, "xxx.hard") {
			continue
		}
		files2 = append(files2, inputPath)
	}
	files = files2

	for i, inputPath := range files {
		fmt.Println("======================== ^^^ ========================")
		fmt.Printf("Pdf File %3d of %d %q\n", i+1, len(files), inputPath)
		err := outputPdfTextRecover(i, len(files), inputPath)
		if err != nil {
			marker := ""
			if err != pdf.ErrEncrypted && err != pdfcore.ErrNoPdfVersion {
				marker = "******"
			}
			fmt.Fprintf(os.Stderr, "Pdf File %3d of %d %q err=%v %s\n",
				i+1, len(files), inputPath, err, marker)
			if err == pdf.ErrEncrypted || err == pdfcore.ErrNoPdfVersion || true {
				continue
			}
			os.Exit(1)
		}
		fmt.Println("======================== ||| ========================")
	}
	fmt.Fprintf(os.Stderr, "Done %d files\n", len(files))
}

// outputPdfTextRecover prints out contents of PDF file to stdout and recovers from panics
func outputPdfTextRecover(i, n int, inputPath string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Recovered: %d of %d: %#q r=%#v\n", i, n, inputPath, r)
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
			panic(err)
		}
	}()
	err = outputPdfText(inputPath)
	return
}

// outputPdfText prints out contents of PDF file to stdout.
func outputPdfText(inputPath string) error {
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
	if err != nil {
		return err
	}

	fmt.Println("---------------------------------------")
	fmt.Printf("PDF text rendering: %q\n", inputPath)
	fmt.Println("---------------------------------------")
	for i := 0; i < numPages; i++ {
		pageNum := i + 1

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

		fmt.Printf("Page %d: %q\n", pageNum, inputPath)
		fmt.Printf("%s\n", text)
		fmt.Println("---------------------------------------")
	}

	return nil
}
