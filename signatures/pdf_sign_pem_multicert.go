/*
 * This example showcases signing a PDF file using a certificate chain and
 * a private key, both extracted from PEM files.
 * The first certificate in the chain is actually used for signing the input
 * PDF file, while the entire chain is embedded in the generated PDF signature
 * in order to validate its authenticity.
 * The example also works if the certificate file contains only the signing
 * certificate.
 *
 * $ ./pdf_sign_pem_multicert <INPUT_PDF_PATH> <OUTPUT_PDF_PATH> <CERTIFICATE_PATH> <PRIVATE_KEY_PATH>
 *
 * Example: ./pdf_sign_pem_multicert in.pdf out.pdf certs.pem key.pem
 */
package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"hash"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/model/sighandler"
)

const (
	usage = "Usage: %s INPUT_PDF_PATH OUTPUT_PDF_PATH CERTS_PATH PRIVATE_KEY_PATH\n"
)

func main() {
	args := os.Args
	if len(args) < 5 {
		fmt.Printf(usage, os.Args[0])
		return
	}
	inputPath, outputPath := args[1], args[2]
	certPath, privateKeyPath := args[3], args[4]

	// Load signing certificate and certificate chain as a PDF array object.
	// The signing certificate is used to sign the input PDF file.
	// The PDF certificate chain (which includes the signing certificate) will
	// be embedded in the generated PDF signature.
	signingCert, pdfCerts, err := loadCertificates(certPath)
	if err != nil {
		log.Fatal(err)
	}

	// Load private key.
	privateKey, err := loadPrivateKey(privateKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	// Sign input file and write output file.
	err = sign(inputPath, outputPath, signingCert, privateKey, pdfCerts)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("PDF file successfully signed. Output path: %s\n", outputPath)
}

func loadPrivateKey(privateKeyPath string) (*rsa.PrivateKey, error) {
	// Read private key file contents.
	privateKeyData, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	// Decode PEM block.
	block, _ := pem.Decode(privateKeyData)

	// Parse private key data.
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func loadCertificates(certPath string) (*x509.Certificate, *core.PdfObjectArray, error) {
	// Read certificate file contents.
	certData, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, nil, err
	}

	parseCert := func(data []byte) (*x509.Certificate, []byte, error) {
		// Decode PEM block.
		block, rest := pem.Decode(data)

		// Parse certificate.
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, nil, err
		}

		return cert, rest, nil
	}

	// Create PDF array object which will contain the certificate chain data,
	// loaded from the PEM file. The first element of the array must be the
	// signing certificate. The rest of the certificate chain is used for
	// validating the authenticity of the signing certificate.
	pdfCerts := core.MakeArray()

	// Parse signing certificate.
	signingCert, pemUnparsedData, err := parseCert(certData)
	if err != nil {
		return nil, nil, err
	}
	pdfCerts.Append(core.MakeString(string(signingCert.Raw)))

	// Parse the rest of the certificates contained in the PEM file,
	// if any, and add them to the PDF certificates array.
	for len(pemUnparsedData) != 0 {
		cert, rest, err := parseCert(pemUnparsedData)
		if err != nil {
			return nil, nil, err
		}

		pdfCerts.Append(core.MakeString(string(cert.Raw)))
		pemUnparsedData = rest
	}

	return signingCert, pdfCerts, nil
}

func sign(inputPath, outputPath string, signingCert *x509.Certificate,
	privateKey *rsa.PrivateKey, pdfCerts *core.PdfObjectArray) error {
	// Create signature function.
	signFunc := func(sig *model.PdfSignature, digest model.Hasher) ([]byte, error) {
		h, ok := digest.(hash.Hash)
		if !ok {
			return nil, errors.New("hash type error")
		}

		return privateKey.Sign(rand.Reader, h.Sum(nil), crypto.SHA1)
	}

	// Create custom signature handler.
	handler, err := sighandler.NewAdobeX509RSASHA1Custom(signingCert, signFunc)
	if err != nil {
		return err
	}

	// Create and initialize signature.
	signature := model.NewPdfSignature(handler)
	if err := signature.Initialize(); err != nil {
		return err
	}

	// Set signature fields.
	signature.SetName("Test PEM Multicert Signature")
	signature.SetReason("Test_PEM_Multicert_Signature")
	signature.SetDate(time.Now(), "")

	// Set signature certificate chain.
	signature.Cert = pdfCerts

	// Create signature field and appearance.
	opts := annotator.NewSignatureFieldOpts()
	opts.FontSize = 10
	opts.Rect = []float64{10, 25, 75, 60}

	sigField, err := annotator.NewSignatureField(
		signature,
		[]*annotator.SignatureLine{
			annotator.NewSignatureLine("Name", "John Doe"),
			annotator.NewSignatureLine("Date", "2020.03.27"),
			annotator.NewSignatureLine("Reason", "UniPDF Signature Test"),
		},
		opts,
	)
	sigField.T = core.MakeString("PEM Multicert signature")

	// Open input file.
	file, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
	defer file.Close()

	// Create reader.
	reader, err := model.NewPdfReader(file)
	if err != nil {
		return err
	}

	// Create appender.
	appender, err := model.NewPdfAppender(reader)
	if err != nil {
		return err
	}

	// Sign input PDF file.
	// The signature appearance is placed on the first page of the output PDF file.
	// See https://github.com/unidoc/unipdf-examples/blob/master/signatures/pdf_sign_appearance.go
	// for adding multiple signature appearance annotations, on multiple pages.
	if err = appender.Sign(1, sigField); err != nil {
		return err
	}

	// Write output file.
	return appender.WriteToFile(outputPath)
}
