/*
 * This example showcases how to validate a digital signature with UniDoc for PAdES compatibility.
 *
 * $ ./pdf_sign_validate_pades_b_b <INPUT_PDF_PATH>
 */
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/model/sighandler"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}
}

const usagef = "Usage: %s INPUT_PDF_PATH\n"

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf(usagef, os.Args[0])
		return
	}
	inputPath := args[1]

	// Create reader.
	file, err := os.Open(inputPath)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}
	defer file.Close()

	reader, err := model.NewPdfReader(file)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	// Create signature handler. It is possible to validate B-T level the same way if you create NewEtsiPAdESLevelT here
	padesHandler, err := sighandler.NewEtsiPAdESLevelB(nil, nil, nil)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	handlers := []model.SignatureHandler{
		padesHandler,
	}

	// Validate signatures.
	res, err := reader.ValidateSignatures(handlers)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}
	if len(res) == 0 {
		log.Fatal("Fail: no signature fields found")
	}

	if !res[0].IsSigned || !res[0].IsVerified {
		log.Fatal("Fail: validation failed")
	}

	for i, item := range res {
		fmt.Printf("--- Signature %d\n%s\n", i+1, item.String())
	}
}
