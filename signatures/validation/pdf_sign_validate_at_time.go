/*
 * This example showcases how to validate the signatures of an archived
 * document at an explicit point in time.
 *
 * By default ValidateSignaturesWithOpts evaluates certificate validity and
 * revocation at the current time, or at the time of a trusted timestamp
 * token when the signature carries one. A signature without a trusted
 * timestamp whose certificate has since expired or been revoked is therefore
 * reported as untrusted today, even if it was fine when it was made.
 *
 * An archive that has other evidence of when a document was signed (a
 * registration log, a trusted receipt) can pass that instant as
 * ValidationTime. The signature dictionary's own /M date is chosen by the
 * signer, so accepting it is a decision the caller makes knowingly; this
 * example does so only when the user asks for it, and only for signatures
 * without a trusted timestamp token. An explicit ValidationTime takes
 * precedence over the token, so passing /M for a timestamped signature would
 * replace the time the TSA vouched for with the one the signer chose.
 * Signatures without /M and document timestamps keep the default evaluation.
 *
 * $ ./pdf_sign_validate_at_time <INPUT_PDF_PATH> <ROOTS_PEM_PATH> [RFC3339_TIME|signing-time]
 */
package main

import (
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/unidoc/unipdf/v5/common/license"
	"github.com/unidoc/unipdf/v5/core"
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

const usagef = "Usage: %s INPUT_PDF_PATH ROOTS_PEM_PATH [RFC3339_TIME|signing-time]\n"

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Printf(usagef, os.Args[0])
		return
	}
	inputPath := args[1]
	rootsPath := args[2]
	timeArg := ""
	if len(args) > 3 {
		timeArg = args[3]
	}

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
		Roots: roots,
		// Archived documents carry their revocation evidence in the DSS
		// (see the LTV examples). Reading it back needs no network and is
		// the only evidence that describes the state at signing time.
		Revocation:              model.RevocationCheckOffline,
		RequireRevocationStatus: true,
	}
	if timeArg != "" && timeArg != "signing-time" {
		at, err := time.Parse(time.RFC3339, timeArg)
		if err != nil {
			log.Fatalf("Fail: invalid time %q: %v\n", timeArg, err)
		}
		opts.ValidationTime = at
	}

	// With no time given each signature is evaluated now, or at the time of
	// its trusted timestamp token.
	start := time.Now()
	res, err := reader.ValidateSignaturesWithOpts(handlers, opts)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
	if len(res) == 0 {
		log.Fatal("Fail: no signature fields found")
	}

	if timeArg == "signing-time" {
		// The caller chooses to believe the signer's own /M date, but only
		// where nothing better fixes the time: an explicit ValidationTime
		// would override a trusted timestamp token, so signatures evaluated
		// at a token time (before start) keep that result. Document
		// timestamps carry no signer date to trust.
		byDate := map[time.Time][]int{}
		for i, item := range res {
			if item.Date == (model.PdfDate{}) || isDocTimeStamp(item) || item.ValidatedAt.Before(start) {
				continue
			}
			at := item.Date.ToGoTime().UTC()
			byDate[at] = append(byDate[at], i)
		}
		for at, indexes := range byDate {
			dateOpts := *opts
			dateOpts.ValidationTime = at
			atRes, err := reader.ValidateSignaturesWithOpts(handlers, &dateOpts)
			if err != nil {
				log.Fatalf("Fail: %v\n", err)
			}
			for _, i := range indexes {
				res[i] = atRes[i]
			}
		}
	}

	for i, item := range res {
		signerDate := "none"
		if item.Date != (model.PdfDate{}) {
			signerDate = item.Date.ToGoTime().Format(time.RFC3339)
		}
		validatedAt := "not evaluated"
		if !item.ValidatedAt.IsZero() {
			validatedAt = item.ValidatedAt.Format(time.RFC3339)
		}

		fmt.Printf("--- Signature %d (%s)\n", i+1, item.Name)
		fmt.Printf("Signer's date:    %s\n", signerDate)
		fmt.Printf("Validated at:     %s\n", validatedAt)
		fmt.Printf("Integrity:        %v\n", item.IsVerified)
		fmt.Printf("Trusted:          %v\n", item.IsTrusted)
		fmt.Printf("Revocation:       %s\n", item.RevocationStatus)
		if !item.RevocationTime.IsZero() {
			fmt.Printf("Revoked at:       %s\n", item.RevocationTime.Format(time.RFC3339))
		}
		for _, reason := range item.TrustErrors {
			fmt.Printf("  error: %s\n", reason)
		}
		for _, warning := range item.TrustWarnings {
			fmt.Printf("  warning: %s\n", warning)
		}
	}
}

func isDocTimeStamp(item model.SignatureValidationResult) bool {
	if len(item.Fields) == 0 {
		return false
	}
	d, ok := core.GetDict(item.Fields[0].V)
	if !ok {
		return false
	}
	name, ok := core.GetNameVal(d.Get("SubFilter"))
	return ok && name == "ETSI.RFC3161"
}
