/*
 * Applies outlines to a PDF file. The files are read from a JSON formatted file,
 * which can be created via pdf_get_outlines which outputs outlines for an input PDF file
 * in the JSON format.
 *
 * Run as: go run pdf_set_outlines.go input.pdf outlines.json output.pdf
 */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/model"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Usage: go run pdf_set_outlines.go input.pdf outlines.json output.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outlinesPath := os.Args[2]
	outPath := os.Args[3]

	// Enable debug-level logging.
	common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))

	fmt.Printf("Input file: %s\n", inputPath)
	fmt.Printf("Outlines file (JSON): %s\n", outlinesPath)
	fmt.Printf("Output file: %s\n", outPath)

	err := applyOutlines(inputPath, outlinesPath, outPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func applyOutlines(inputPath, outlinesPath, outPath string) error {
	data, err := ioutil.ReadFile(outlinesPath)
	if err != nil {
		return err
	}

	var newOutlines model.Outline
	err = json.Unmarshal(data, &newOutlines)
	if err != nil {
		return err
	}

	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return err
	}

	// Try decrypting with an empty one.
	if isEncrypted {
		auth, err := pdfReader.Decrypt([]byte(""))
		if err != nil {
			return err
		}
		if !auth {
			common.Log.Debug("Encrypted - unable to access - update code to specify pass")
			return nil
		}
	}

	fw, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer fw.Close()

	w := model.NewPdfWriter()
	for _, p := range pdfReader.PageList {
		err = w.AddPage(p)
		if err != nil {
			return err
		}
	}
	w.AddOutlineTree(newOutlines.ToOutlineTree())

	return w.Write(fw)
}
