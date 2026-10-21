/*
 * This example showcases how to create a PAdES B-T compatible digital signature for a PDF file.
 *
 * $ ./pdf_sign_pades_b_t <FILE.PFX> <PASSWORD> <FILE.PEM> <INPUT_PDF_PATH> <OUTPUT_PDF_PATH>
 */
package main

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/unidoc/unipdf/v5/annotator"
	"github.com/unidoc/unipdf/v5/common/license"
	"github.com/unidoc/unipdf/v5/core"
	"github.com/unidoc/unipdf/v5/model"
	"github.com/unidoc/unipdf/v5/model/sighandler"

	"software.sslmate.com/src/go-pkcs12"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}
}

const usagef = "Usage: %s PFX_FILE PASSWORD PEM_FILE INPUT_PDF_PATH OUTPUT_PDF_PATH\n"

func main() {
	args := os.Args
	if len(args) < 6 {
		fmt.Printf(usagef, os.Args[0])
		return
	}
	pfxPath := args[1]
	password := args[2]
	pemPath := args[3]
	inputPath := args[4]
	outputPath := args[5]

	// Get private key and X509 certificate from the PFX file.
	pfxData, err := os.ReadFile(pfxPath)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	priv, cert, _, err := pkcs12.DecodeChain(pfxData, password)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Get cacert certificate from the PEM file.
	caCertF, err := os.ReadFile(pemPath)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	certDERBlock, _ := pem.Decode(caCertF)
	if certDERBlock == nil {
		log.Fatalf("Fail: no PEM certificate found in %s\n", pemPath)
	}

	cacert, err := x509.ParseCertificate(certDERBlock.Bytes)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

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

	// Create appender.
	appender, err := model.NewPdfAppender(reader)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Set timestamp server.
	timestampServerURL := "https://freetsa.org/tsr"

	// Create signature handler.
	var handler model.SignatureHandler
	switch key := priv.(type) {
	case *rsa.PrivateKey:
		handler, err = sighandler.NewEtsiPAdESLevelT(key, cert, cacert, timestampServerURL)
	case *ecdsa.PrivateKey:
		handler, err = sighandler.NewEtsiPAdESLevelTEcdsa(key, cert, cacert, timestampServerURL)
	default:
		log.Fatalf("Fail: unsupported private key type %T\n", priv)
	}
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Create signature.
	signature := model.NewPdfSignature(handler)
	signature.SetName("PAdES B-T Signature PDF")
	signature.SetReason("TestPAdESPDF")
	signature.SetDate(time.Now(), "")

	if err := signature.Initialize(); err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Create signature field and appearance.
	opts := annotator.NewSignatureFieldOpts()
	opts.FontSize = 10
	opts.Rect = []float64{10, 25, 75, 60}

	field, err := annotator.NewSignatureField(
		signature,
		[]*annotator.SignatureLine{
			annotator.NewSignatureLine("Name", "John Doe"),
			annotator.NewSignatureLine("Date", "2023.05.08"),
			annotator.NewSignatureLine("Reason", "PAdES signature test"),
		},
		opts,
	)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
	field.T = core.MakeString("PAdES B-T Signature")

	if err = appender.Sign(1, field); err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Write output PDF file.
	err = appender.WriteToFile(outputPath)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	log.Printf("PDF file successfully signed. Output path: %s\n", outputPath)
}
