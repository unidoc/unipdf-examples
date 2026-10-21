/*
 * This example showcases how to validate the digital signatures of a PDF with
 * UniDoc. It covers adbe.x509.rsa_sha1, adbe.pkcs7.detached,
 * ETSI.CAdES.detached (PAdES B-B to B-LTA) and ETSI.RFC3161 document
 * timestamps, and exits non-zero if any signature fails.
 *
 * ValidateSignatures checks integrity only: IsVerified means the signature
 * matches the document under the certificate embedded in the PDF, which any
 * self-signed certificate satisfies. To check that the signer is trusted, see
 * pdf_sign_validate_trust.go.
 *
 * $ ./pdf_sign_validate <INPUT_PDF_PATH>
 */
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/unidoc/unipdf/v5/common/license"
	"github.com/unidoc/unipdf/v5/model"
	"github.com/unidoc/unipdf/v5/model/sighandler"
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
		log.Fatalf("Fail: %v\n", err)
	}
	defer file.Close()

	reader, err := model.NewPdfReader(file)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Create signature handlers.
	handlerX509RSASHA1, err := sighandler.NewAdobeX509RSASHA1(nil, nil)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	handlerPKCS7Detached, err := sighandler.NewAdobePKCS7Detached(nil, nil)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// The PAdES handler validates every PAdES level, not only B-B.
	handlerPAdES, err := sighandler.NewEtsiPAdESLevelB(nil, nil, nil)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	handlerDocTimeStamp, err := sighandler.NewDocTimeStamp("", 0)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	handlers := []model.SignatureHandler{
		handlerX509RSASHA1,
		handlerPKCS7Detached,
		handlerPAdES,
		handlerDocTimeStamp,
	}

	// Validate signatures.
	res, err := reader.ValidateSignatures(handlers)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
	if len(res) == 0 {
		log.Fatal("Fail: no signature fields found")
	}

	var failed []string
	for i, item := range res {
		fmt.Printf("--- Signature %d\n%s\n", i+1, item.String())
		if !item.IsSigned || !item.IsVerified {
			failed = append(failed, strconv.Itoa(i+1))
		}
	}
	if len(failed) > 0 {
		log.Fatalf("Fail: validation failed for signature(s) %s\n", strings.Join(failed, ", "))
	}
}
