/*
 * Shows all fonts in a PDF file.
 *
 * Run as: go run pdf_fonts.go input.pdf
 */

package main

import (
	"errors"
	"fmt"
	"os"
	"sort"

	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Syntax: go run pdf_fonts.go input.pdf")
		os.Exit(1)
	}

	// Enable debug-level logging.
	//unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	inputPath := os.Args[1]

	fmt.Printf("Input file: %s\n", inputPath)
	fonts, err := fontsInPdf(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	sort.Slice(fonts, func(i, j int) bool {
		t1, t2 := fonts[i].Subtype(), fonts[j].Subtype()
		if t1 != t2 {
			return t1 < t2
		}
		return fonts[i].BaseFont() < fonts[j].BaseFont()
	})
	fmt.Printf("%d fonts\n", len(fonts))
	for i, font := range fonts {
		fmt.Printf("%3d: %-8s %#q\n", i, font.Subtype(), font.BaseFont())
	}
}

// fontsInPdf returns a list of the fonts in PDF file `inputPath`.
func fontsInPdf(inputPath string) (fonts []pdf.PdfFont, err error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return
	}
	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return
	}
	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return
	}

	// Try decrypting with an empty one.
	if isEncrypted {
		auth := false
		auth, err = pdfReader.Decrypt([]byte(""))
		if err != nil {
			return
		}
		if !auth {
			err = errors.New("Unable to decrypt password protected file - need to specify pass to Decrypt")
			return
		}
	}

	for _, objNum := range pdfReader.GetObjectNums() {
		var obj pdfcore.PdfObject
		obj, err = pdfReader.GetIndirectObjectByNumber(objNum)
		if err != nil {
			return
		}
		font, err := pdf.NewPdfFontFromPdfObject(obj)
		if err != nil {
			continue
		}
		fonts = append(fonts, *font)
	}

	return
}
