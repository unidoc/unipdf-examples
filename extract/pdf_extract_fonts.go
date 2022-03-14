package main

import (
	"archive/zip"
	"fmt"
	"os"
	"strconv"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
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
	var (
		outputPath string
		pageNumber int
	)
	const usage = "Usage: go run pdf_extract_fonts.go input.pdf <pagenumber> output.zip\n<pagenumber> is int, set to 0 to extract all page fonts.\n"
	if len(os.Args) < 2 {
		fmt.Printf(usage)
		os.Exit(1)
	}
	inputPath := os.Args[1]
	if len(os.Args) > 2 {
		if i, err := strconv.Atoi(os.Args[2]); err == nil {
			pageNumber = i
		} else {
			fmt.Printf("Error: %v\n", err)
			fmt.Printf(usage)
			os.Exit(1)
		}
	}
	if len(os.Args) > 3 {
		outputPath = os.Args[3]
	}

	err := extractFontToArchive(inputPath, outputPath, pageNumber)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Printf(usage)
		os.Exit(1)
	}
}

// extractFontToArchive prints out font information from pdf to stdout, can add to archive if has output path. it can be used for specific page.
func extractFontToArchive(inputPath string, outputPath string, page int) error {
	var (
		pdfFont           *extractor.PageFonts
		zipFile           *os.File
		zipw              *zip.Writer
		embeddedFontTotal = 0
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

	if page > 0 { // for specific page
		pdfFont, err = extractSpecificPageFont(pdfFont, pdfReader, page)
		if err != nil {
			return err
		}
	} else { // for range of PDF pages
		for i := 0; i < pdfPage; i++ {
			pdfFont, err = extractSpecificPageFont(pdfFont, pdfReader, (i + 1))
			if err != nil {
				return err
			}
		}
	}

	if len(outputPath) > 0 {
		zipFile, err = os.Create(outputPath)
		if err != nil {
			return err
		}
		zipw = zip.NewWriter(zipFile)
	}

	for _, font := range pdfFont.Fonts {
		var (
			hasEmbeddedFont = false
			extracted       = false
		)
		if len(font.FontFileName) > 0 {
			hasEmbeddedFont = true
			if len(outputPath) > 0 {
				zipFile, err := zipw.Create(font.FontFileName)
				if err != nil {
					return err
				}
				_, err = zipFile.Write(font.FontData)
				if err != nil {
					return err
				}
				extracted = true
			}
			embeddedFontTotal++

		}
		fmt.Println("------------------------------")
		fmt.Printf("Font Name \t: %s\nType \t\t: %s\nEncoding \t: %v\nIsCID\t\t: %t\nIsSimple\t: %t\nToUnicode\t: %t\nEmbedded\t: %v\nExtracted\t: %v", font.FontName, font.FontType, font.PdfFont.Encoder().String(), font.IsCID, font.IsSimple, font.ToUnicode, hasEmbeddedFont, extracted)
		fmt.Println("\n------------------------------\n")
	}
	if len(outputPath) > 0 {
		err = zipw.Close()
		if embeddedFontTotal == 0 {
			os.Remove(outputPath)
		}
	}

	if err != nil {
		return err
	}
	return nil
}

func extractSpecificPageFont(pdfFont *extractor.PageFonts, pdfReader *model.PdfReader, page int) (*extractor.PageFonts, error) {
	currentPage, err := pdfReader.GetPage(page)
	if err != nil {
		return nil, err
	}
	pdfExtractor, err := extractor.New(currentPage)
	if err != nil {
		return nil, err
	}
	pdfFont, err = pdfExtractor.ExtractFonts(pdfFont)
	return pdfFont, nil
}
