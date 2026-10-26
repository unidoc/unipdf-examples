/*
 * This example showcases how to check how much of a PDF file each digital
 * signature covers.
 *
 * A signature only protects the bytes listed in its ByteRange. Content can be
 * appended to a signed file in an incremental update without breaking the
 * signature, so a signature can be intact, and even trusted, while the
 * document a reader opens is not the one that was signed. Appended revisions
 * are normal (a second signature, LTV data, form filling), but each one needs
 * to be accounted for. The coverage fields of SignatureValidationResult
 * report, per signature:
 *   - IsCoveringEntireFile: no bytes were appended after the signature.
 *   - IsCoveringEntireRevision: the signature covers its whole revision.
 *   - RevisionsAfterSignature: how many revisions were appended after it.
 *   - ContentsGapOnly: the only bytes left out of the signed range are the
 *     signature's own /Contents value.
 *
 * The example exits non-zero if a signature does not verify, does not cover
 * its revision, or leaves out more than its /Contents value. To see what the
 * appended revisions changed, extract the signed revision with
 * PdfReader.GetRevision (see signing/pdf_sign_get_revision.go) and compare it
 * with the full document.
 *
 * $ ./pdf_sign_validate_coverage <INPUT_PDF_PATH>
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

	file, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
	defer file.Close()

	reader, err := model.NewPdfReader(file)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	handlerX509RSASHA1, err := sighandler.NewAdobeX509RSASHA1(nil, nil)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
	handlerPKCS7Detached, err := sighandler.NewAdobePKCS7Detached(nil, nil)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
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

	// The coverage fields are derived from the file bytes, so they are set
	// by ValidateSignatures and ValidateSignaturesWithOpts alike.
	res, err := reader.ValidateSignatures(handlers)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
	if len(res) == 0 {
		log.Fatal("Fail: no signature fields found")
	}

	var failed []string
	for i, item := range res {
		name := ""
		if len(item.Fields) > 0 {
			name = item.Fields[0].PartialName()
		}
		fmt.Printf("--- Signature %d: field %q\n", i+1, name)

		// IsVerified is also false when no handler supports the signature or
		// its ByteRange is missing or invalid; Errors says which.
		if !item.IsVerified {
			fmt.Println("Result: signature not verified")
			for _, e := range item.Errors {
				fmt.Printf("  - %s\n", e)
			}
			failed = append(failed, strconv.Itoa(i+1))
			continue
		}

		ok := true
		switch {
		case item.IsCoveringEntireFile:
			fmt.Println("Result: signature covers the entire file")
		case item.IsCoveringEntireRevision:
			fmt.Printf("Result: signature covers its revision; %d revision(s) appended after it\n",
				item.RevisionsAfterSignature)
		default:
			fmt.Println("Result: signed range stops before the end of its revision")
			ok = false
		}
		// Signers write /Contents as a hex string and start the ByteRange at
		// offset 0. Anything else excludes bytes the signature does not
		// protect, so it is treated as a failure here.
		if !item.ContentsGapOnly {
			fmt.Println("Warning: ByteRange leaves out more than the signature /Contents")
			ok = false
		}

		if !ok {
			failed = append(failed, strconv.Itoa(i+1))
		}
	}
	if len(failed) > 0 {
		log.Fatalf("Fail: coverage check failed for signature(s) %s\n", strings.Join(failed, ", "))
	}
}
