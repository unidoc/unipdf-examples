package main

import (
	"archive/zip"
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

// func init() {
// 	// Make sure to load your metered License API key prior to using the library.
// 	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
// 	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
// 	if err != nil {
// 		panic(err)
// 	}
// }
func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Syntax: go run pdf_extract_fonts.go input.pdf output.zip\n")
		os.Exit(1)
	}
	inputPath := os.Args[1]
	outputPath := os.Args[2]

	err := extractFontToArchive(inputPath, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
func extractFontToArchive(inputPath string, outputPath string) error {

	pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, nil)
	if err != nil {
		return err
	}
	defer f.Close()
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Printf("PDF Num Pages: %d\n", numPages)
	zipFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	zipw := zip.NewWriter(zipFile)
	afont := []string{} //Font Bag

	for i := 0; i < numPages; i++ {
		page, err := pdfReader.GetPage(i + 1)
		if err != nil {
			return err
		}
		pdfExtractor, err := extractor.New(page)
		lFont, err := pdfExtractor.ExtractPageFonts()
		for _, f := range lFont.Fonts {
			if InsideBag(afont, f.FontName) { //Check Duplicate Font from other Page
				continue
			}
			if f.FontFileName != "" {
				zipFile, err := zipw.Create(f.FontFileName)
				if err != nil {
					return err
				}
				_, err = zipFile.Write(f.FontData)
				if err != nil {
					return err
				}
			}
			afont = append(afont, f.FontName)
			fmt.Printf("Font Name \t: %s\nType \t\t: %s\nEncoding \t: %v\nIsCID\t\t: %t\nIsSimple\t: %t\nToUnicode\t: %t\n  \n\n", f.FontName, f.FontType, f.PdfFont.Encoder().String(), f.IsCID, f.IsSimple, f.ToUnicode)
		}

		if err != nil {
			return err
		}
	}
	err = zipw.Close()
	if err != nil {
		return err
	}
	return nil
}

//InsideBag is for check duplicate font in every page.
func InsideBag(afont []string, s string) bool {
	for _, v := range afont {
		if v == s {
			return true
		}
	}
	return false
}
