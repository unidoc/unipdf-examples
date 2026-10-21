/*
 * This example showcases how to validate digital signatures with online
 * certificate revocation checks (CRL and OCSP).
 *
 * With RevocationCheckOnline, unipdf first looks for revocation evidence
 * embedded in the document (DSS, VRI entries, PAdES revocation attribute).
 * Only for chain certificates that evidence does not settle does it contact
 * the CRL distribution points and OCSP responders named in the certificate,
 * over HTTP or HTTPS, through the clients configured here. A request that
 * fails leaves the status unknown. With RequireRevocationStatus, as set
 * here, an unknown status makes the signature untrusted; without it the
 * failure is only reported as a warning.
 *
 * $ ./pdf_sign_validate_online <INPUT_PDF_PATH> <ROOTS_PEM_PATH>
 */
package main

import (
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/unidoc/unipdf/v5/common/license"
	"github.com/unidoc/unipdf/v5/model"
	"github.com/unidoc/unipdf/v5/model/sighandler"
	"github.com/unidoc/unipdf/v5/model/sigutil"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}
}

const usagef = "Usage: %s INPUT_PDF_PATH ROOTS_PEM_PATH\n"

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Printf(usagef, os.Args[0])
		return
	}
	inputPath := args[1]
	rootsPath := args[2]

	rootsPEM, err := os.ReadFile(rootsPath)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(rootsPEM) {
		log.Fatalf("Fail: no certificates found in %s\n", rootsPath)
	}

	// The HTTP client used for CRL and OCSP requests. Set the timeout,
	// proxy and TLS configuration your environment needs here; the same
	// client can be shared across many validations and wrapped with a cache
	// if CRLs from the same issuer are fetched repeatedly.
	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: http.DefaultTransport.(*http.Transport).Clone(),
	}
	crlClient := sigutil.NewCRLClient()
	crlClient.HTTPClient = httpClient
	ocspClient := sigutil.NewOCSPClient()
	ocspClient.HTTPClient = httpClient

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
		Revocation: model.RevocationCheckOnline,
		CRLClient:  crlClient,
		OCSPClient: ocspClient,
		// An unreachable responder leaves the status unknown. Requiring a
		// known status turns that into an untrusted signature (fail closed),
		// which is the right choice when the outcome drives a decision.
		// Without it the signature stays trusted and the problem is listed
		// in TrustWarnings.
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
		fmt.Printf("--- Signature %d (%s)\n", i+1, item.Name)
		fmt.Printf("Integrity:  %v\n", item.IsVerified)
		fmt.Printf("Trusted:    %v\n", item.IsTrusted)
		fmt.Printf("Revocation: %s (signer CRL found: %v, signer OCSP found: %v)\n",
			item.RevocationStatus, item.IsCrlFound, item.IsOcspFound)
		if item.RevocationStatus == model.RevocationStatusRevoked {
			fmt.Printf("Revoked at: %s\n", item.RevocationTime.Format(time.RFC3339))
		}
		for _, reason := range item.TrustErrors {
			fmt.Printf("  error: %s\n", reason)
		}
		for _, warning := range item.TrustWarnings {
			fmt.Printf("  warning: %s\n", warning)
		}
	}
}
