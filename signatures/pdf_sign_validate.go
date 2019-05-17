/*
 * This example showcases how to validate a digital signature with UniDoc.
 *
 * $ ./pdf_sign_validate <INPUT_PDF_PATH>
 */
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/model/sighandler"
)

const usage = "Usage: %s P12_FILE PASSWORD INPUT_PDF_PATH OUTPUT_PDF_PATH\n"

func main() {
	args := os.Args
	if len(args) < 2 {
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
	handlerX509RSASHA1, err := sighandler.NewAdobeX509RSASHA1(nil, nil)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	handlerPKCS7Detached, err := sighandler.NewAdobePKCS7Detached(nil, nil)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	handlers := []model.SignatureHandler{
		handlerX509RSASHA1,
		handlerPKCS7Detached,
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
