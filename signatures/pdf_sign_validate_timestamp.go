/*
 * This example showcases how to validate a digital signature containing timestamp with UniDoc.
 *
 * $ ./pdf_sign_validate_timestamp <INPUT_PDF_PATH>
 */
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/model/sighandler"
)

const usage = "Usage: %s INPUT_PDF_PATH\n"

func main() {
	args := os.Args
	if len(args) < 1 {
		fmt.Printf(usage, os.Args[0])
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

	// Create signature handlers.
	handler, err := sighandler.NewDocTimeStamp("", 0)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	handlers := []model.SignatureHandler{
		handler,
	}

	// Validate signatures.
	res, err := reader.ValidateSignatures(handlers)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}
	if len(res) == 0 {
		log.Fatal("Fail: no signature fields found")
	}

	timestampRes := make(map[int]model.SignatureValidationResult)

	for i, item := range res {
		// find only timestamp signatures
		if d, ok := core.GetDict(item.Fields[0].V); ok && d.Get("SubFilter").String() == "ETSI.RFC3161" {
			timestampRes[i] = item
		}
	}

	if len(timestampRes) == 0 {
		log.Fatal("Fail: validation failed")
	}

	for i, item := range timestampRes {
		fmt.Printf("--- Signature %d\n%s\n", i+1, item.String())
	}
}
