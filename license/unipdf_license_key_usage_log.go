/*
 * This example shows how to get license key usage logs
 *
 * Run as: go run unipdf_license_usage_log.go input_dir output_dir
 */

package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/unidoc/unipdf/v3/common"
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
	license.SetMeteredKeyUsageLogVerboseMode(true)
	common.SetLogger(common.NewConsoleLogger(common.LogLevelInfo))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run unipdf_license_usage_log.go input_dir output_dir\n")
		os.Exit(1)
	}

	input_dir := os.Args[1]
	output_dir := os.Args[2]
	err := extractMultiple(input_dir, output_dir)
	if err != nil {
		panic(err)
	}
}

// extractPdfText extracts the text of pdf file provided by inputPath.
func extractPdfText(inputPath string) (string, error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return "", err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return "", err
	}

	text := ""
	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return "", err
		}

		ex, err := extractor.New(page)
		if err != nil {
			return "", err
		}

		res, err := ex.ExtractText()
		if err != nil {
			return "", err
		}
		pageMark := fmt.Sprintf("--- page %d ---- \n", pageNum)
		text += pageMark + res
	}

	return text, nil
}

// extractMultiple extracts text from multiple files in `inputDir`.
func extractMultiple(inputDir string, outputDir string) error {
	files, err := filepath.Glob(inputDir + "*.pdf")
	if err != nil {
		return err
	}

	for _, file := range files {
		txt, err := extractPdfText(file)
		if err != nil {
			return err
		}

		base := path.Base(file)
		name := strings.TrimSuffix(base, filepath.Ext(base))
		dest := path.Join(outputDir, name+".txt")

		f, err := os.Create(dest)
		if err != nil {
			return err
		}

		f.WriteString(txt)
	}
	return nil
}
