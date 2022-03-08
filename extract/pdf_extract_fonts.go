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
	var outputPath string
	if len(os.Args) < 2 {
		fmt.Printf("Syntax: go run extract/pdf_extract_fonts.go input.pdf output.zip\n")
		os.Exit(1)
	}
	inputPath := os.Args[1]
	if len(os.Args) == 3 {
		outputPath = os.Args[2]
	}

	err := extractFontToArchive(inputPath, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// extractFontToArchive prints out font information from pdf to stdout, can add to archive if has output path.
func extractFontToArchive(inputPath string, outputPath string) error {
	var (
		pdfFont *extractor.PageFonts
		zipFile *os.File
		zipw    *zip.Writer
	)

	pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, nil)
	if err != nil {
		return err
	}
	defer f.Close()
	pdfPage, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}
	for i := 0; i < pdfPage; i++ {
		currentPage, err := pdfReader.GetPage(i + 1)
		if err != nil {
			return err
		}
		pdfExtractor, err := extractor.New(currentPage)
		if err != nil {
			return err
		}
		pdfFont, err = pdfExtractor.ExtractFonts(pdfFont)
	}

	if len(outputPath) > 0 {
		zipFile, err = os.Create(outputPath)
		if err != nil {
			return err
		}
		zipw = zip.NewWriter(zipFile)
	}

	for _, font := range pdfFont.Fonts {
		if len(font.FontFileName) > 0 && len(outputPath) > 0 {
			zipFile, err := zipw.Create(font.FontFileName)
			if err != nil {
				return err
			}
			_, err = zipFile.Write(font.FontData)
			if err != nil {
				return err
			}
		}
		fmt.Println("------------------------------")
		fmt.Printf("Font Name \t: %s\nType \t\t: %s\nEncoding \t: %v\nIsCID\t\t: %t\nIsSimple\t: %t\nToUnicode\t: %t", font.FontFileName, font.FontType, font.PdfFont.Encoder().String(), font.IsCID, font.IsSimple, font.ToUnicode)
		fmt.Println("\n------------------------------\n")
	}
	if len(outputPath) > 0 {
		err = zipw.Close()
	}
	if err != nil {
		return err
	}
	return nil
}
