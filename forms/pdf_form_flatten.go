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

	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/common/license"
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
	err = pdfReader.FlattenFields(true, fieldAppearance)
	if err != nil {
		return err
	}

	// AcroForm field is no longer needed.
	opt := &model.ReaderToWriterOpts{
		SkipAcroForm: true,
	}

	// Generate a PdfWriter instance from existing PdfReader.
	pdfWriter, err := pdfReader.ToWriter(opt)
	if err != nil {
		return err
	}

	// Write to file.
	err = pdfWriter.WriteToFile(outputPath)
	return err
}
