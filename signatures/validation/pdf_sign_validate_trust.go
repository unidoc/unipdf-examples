/*
 * This example showcases how to validate digital signatures against a set of
 * trusted root certificates, including certificate revocation checks.
 *
 * Validating with ValidateSignatures alone only proves that the signature
 * matches the document bytes under the certificate embedded in the PDF. Any
 * self-signed certificate passes that check. ValidateSignaturesWithOpts also
 * verifies that the signer certificate chains to a root you trust, that it
 * is valid at the validation time (now, or the time of a trusted timestamp
 * token), and, when requested, that it was not revoked.
 *
 * $ ./pdf_sign_validate_trust <INPUT_PDF_PATH> <ROOTS_PEM_PATH> [offline|online]
 */
package main

import (
	"crypto/x509"
	"fmt"
	"log"
	"os"

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

const usagef = "Usage: %s INPUT_PDF_PATH ROOTS_PEM_PATH [offline|online]\n"

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Printf(usagef, os.Args[0])
		return
	}
	inputPath := args[1]
	rootsPath := args[2]

	// Revocation data is read from the document (DSS, PAdES revocation
	// attribute) by default. "online" additionally fetches CRLs and OCSP
	// responses from the URLs named in the certificates.
	revocation := model.RevocationCheckOffline
	if len(args) > 3 {
		switch args[3] {
		case "offline":
		case "online":
			revocation = model.RevocationCheckOnline
		default:
			fmt.Printf(usagef, os.Args[0])
			return
		}
	}

	// The trust roots come from the caller, never from the document. Load
	// them from a PEM bundle of the CA certificates you accept.
	rootsPEM, err := os.ReadFile(rootsPath)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(rootsPEM) {
		log.Fatalf("Fail: no certificates found in %s\n", rootsPath)
	}

	file, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
	defer file.Close()

	reader, err := model.NewPdfReader(file)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// One handler per signature type; each verifies integrity only.
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
	handlerTimestamp, err := sighandler.NewDocTimeStamp("", 0)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
	handlers := []model.SignatureHandler{
		handlerX509RSASHA1,
		handlerPKCS7Detached,
		handlerPAdES,
		handlerTimestamp,
	}

	opts := &model.SignatureValidationOpts{
		Roots:      roots,
		Revocation: revocation,
		// Treat an unknown revocation status as untrusted. Leave this unset
		// to accept documents without revocation data and inspect
		// RevocationStatus yourself.
		RequireRevocationStatus: true,
	}

	res, err := reader.ValidateSignaturesWithOpts(handlers, opts)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
	if len(res) == 0 {
		log.Fatal("Fail: no signature fields found")
	}

	for i, item := range res {
		fmt.Printf("--- Signature %d\n%s\n", i+1, item.String())
		switch {
		case !item.IsVerified:
			fmt.Println("Result: signature does not match the document (integrity failure)")
		case !item.IsTrusted:
			fmt.Println("Result: signature is intact but the signer is not trusted")
			for _, reason := range item.TrustErrors {
				fmt.Printf("  - %s\n", reason)
			}
		default:
			fmt.Printf("Result: signed by %s, trusted via %s\n",
				item.TrustChain[0].Subject.CommonName,
				item.TrustChain[len(item.TrustChain)-1].Subject.CommonName)
			for _, warning := range item.TrustWarnings {
				fmt.Printf("  warning: %s\n", warning)
			}
		}
	}
}
