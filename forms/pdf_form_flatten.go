/*
 * Flatten form data in PDF files, moving to content stream from annotations, so cannot be edited.
 * Note: Works for forms that have been filled in an editor and have the appearance streams generated.
 *
 * Run as: go run pdf_form_flatten.go <outputdir> <pdf files...>
 */

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/model"
)

func main() {
	// When debugging, enable debug-level logging via console:
	common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))

	if len(os.Args) < 3 {
		fmt.Printf("Usage: go run pdf_form_flatten.go <outputdir> <input1.pdf> [input2.pdf] ...\n")
		os.Exit(1)
	}

	outputDir := os.Args[1]

	fails := map[string]string{}
	failKeys := []string{}
	processed := 0

	for i := 2; i < len(os.Args); i++ {
		inputPath := os.Args[i]
		name := filepath.Base(inputPath)
		outputPath := filepath.Join(outputDir, fmt.Sprintf("flattened_%s", name))
		err := flattenPdf(inputPath, outputPath)
		if err != nil {
			fmt.Printf("%s - Error: %v\n", inputPath, err)
			fails[inputPath] = err.Error()
			failKeys = append(failKeys, inputPath)
		}
		processed++
	}

	fmt.Printf("Total %d processed / %d failures\n", processed, len(failKeys))
	sort.Strings(failKeys)
	for _, k := range failKeys {
		fmt.Printf("%s: %v\n", k, fails[k])
	}
}

// flattenPdf flattens annotations and forms moving the appearance stream to the page contents so cannot be
// modified.
func flattenPdf(inputPath, outputPath string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return err
	}

	fieldAppearance := annotator.FieldAppearance{OnlyIfMissing: true}
	err = pdfReader.FlattenFields(false, fieldAppearance)
	if err != nil {
		return err
	}

	pdfWriter := model.NewPdfWriter()
	pdfWriter.SetForms(nil)

	for _, p := range pdfReader.PageList {
		err := pdfWriter.AddPage(p)
		if err != nil {
			return err
		}
	}

	fout, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer fout.Close()

	err = pdfWriter.Write(fout)
	return err
}
