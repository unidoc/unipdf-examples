/*
 * Retrieves outlines (bookmarks) from a PDF file and prints out in JSON format.
 * Note: The JSON output can be used with the related pdf_set_outlines.go example to
 * apply outlines to a PDF file.
 *
 * Run as: go run pdf_get_outlines.go input.pdf > outlines.json
 */

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/model"
)

const licenseKey = `
-----BEGIN UNIDOC LICENSE KEY-----
Free trial license keys are available at: https://unidoc.io/
-----END UNIDOC LICENSE KEY-----
`

func init() {
	// Enable debug-level logging.
	// unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	err := license.SetLicenseKey(licenseKey, `Company Name`)
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage:  go run pdf_get_outlines.go input.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]

	fmt.Printf("Input file: %s\n", inputPath)

	err := getOutlines(inputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func getOutlines(inputPath string) error {
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

	outlines, err := pdfReader.GetOutlines()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(outlines, "", "    ")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", data)

	return nil
}
