/*
 * This example showcases how to create a PAdES B-LTA compatible digital signature for a PDF file.
 *
 * $ ./pdf_sign_pades_b_lta <FILE.PFX> <PASSWORD> <FILE.PEM> <INPUT_PDF_PATH> <OUTPUT_PDF_PATH>
 */
package main

import (
	"bytes"
	"crypto"
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

	// Get CA Certificate from the PEM file.
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
		handler, err = sighandler.NewEtsiPAdESLevelLT(key, cert, cacert, timestampServerURL, appender)
	case *ecdsa.PrivateKey:
		handler, err = sighandler.NewEtsiPAdESLevelLTEcdsa(key, cert, cacert, timestampServerURL, appender)
	default:
		log.Fatalf("Fail: unsupported private key type %T\n", priv)
	}
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Create signature.
	signature := model.NewPdfSignature(handler)
	signature.SetName("PAdES B-LTA Signature PDF")
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
	field.T = core.MakeString("PAdES B-LTA Signature")

	if err = appender.Sign(1, field); err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Write output to buffer.
	buffer := bytes.NewBuffer(nil)
	err = appender.Write(buffer)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// We need the second pass to correctly save DSS/VRI information.
	pdf2, err := model.NewPdfReader(bytes.NewReader(buffer.Bytes()))
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	appender2, err := model.NewPdfAppender(pdf2)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	appender2.SetDSS(appender.GetDSS())

	buf2 := bytes.NewBuffer(nil)
	err = appender2.Write(buf2)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Finally, we add document timestamp for B-LTA compatibility.
	pdf3, err := model.NewPdfReader(bytes.NewReader(buf2.Bytes()))
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	appender3, err := model.NewPdfAppender(pdf3)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// A fixed Contents size makes an oversized token fail with
	// model.ErrSignNotEnoughSpace; without it a token a few bytes longer
	// than the size probe is silently truncated.
	handler, err = sighandler.NewDocTimeStampWithOpts(timestampServerURL, crypto.SHA512, &sighandler.DocTimeStampOpts{
		SignatureSize: 8192,
	})
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	signature = model.NewPdfSignature(handler)
	signature.SetName("Test Signature")
	signature.SetDate(time.Now(), "")

	err = signature.Initialize()
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Create signature appearance.
	opts = annotator.NewSignatureFieldOpts()
	opts.Rect = []float64{0, 0, 0, 0}

	sigField, err := annotator.NewSignatureField(
		signature,
		[]*annotator.SignatureLine{
			annotator.NewSignatureLine("Name", "Jane Doe"),
			annotator.NewSignatureLine("Reason", "Test Document TimeStamp"),
		},
		opts,
	)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	if err = appender3.Sign(1, sigField); err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	err = appender3.WriteToFile(outputPath)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	log.Printf("PDF file successfully signed. Output path: %s\n", outputPath)
}
