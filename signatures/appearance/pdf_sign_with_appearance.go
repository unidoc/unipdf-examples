/*
 * This example showcases how to sign a PDF document with a visible signature
 * whose appearance is generated from the signature dictionary.
 *
 * PdfAppender.SignWithAppearance passes the signature field to a
 * FieldAppearanceGenerator before signing. annotator.FieldAppearance draws
 * the Name, M (signing date), Reason and Location entries of the signature
 * inside the widget rectangle, so the visible text always matches what was
 * signed. pdf_sign_appearance.go builds the same kind of appearance from
 * hand-written SignatureLines instead.
 * The file is signed using a generated private/public key pair.
 *
 * $ ./pdf_sign_with_appearance <INPUT_PDF_PATH> <OUTPUT_PDF_PATH>
 */
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/unidoc/unipdf/v5/annotator"
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

var now = time.Now()

const usagef = "Usage: %s INPUT_PDF_PATH OUTPUT_PDF_PATH\n"

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Printf(usagef, os.Args[0])
		return
	}
	inputPath := args[1]
	outputPath := args[2]

	// Generate key pair.
	priv, cert, err := generateKeys()
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

	// Create signature handler.
	handler, err := sighandler.NewAdobePKCS7Detached(priv, cert)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// The entries set here are the ones drawn in the appearance.
	signature := model.NewPdfSignature(handler)
	signature.SetName("Jane Doe")
	signature.SetDate(now, "")
	signature.SetReason("Contract approval")
	signature.SetLocation("New York")

	if err := signature.Initialize(); err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Create the signature field with an empty widget of the size and
	// position the appearance is drawn in. The field is numbered after the
	// Approval fields already in the document, so signing the output again
	// adds a new field above the previous one instead of reusing its name.
	n := 1
	for hasField(reader, fmt.Sprintf("Approval %d", n)) {
		n++
	}
	y := 50 + float64(n-1)*80
	field := model.NewPdfFieldSignature(signature)
	field.T = core.MakeString(fmt.Sprintf("Approval %d", n))
	field.Rect = core.MakeArrayFromFloats([]float64{50, y, 250, y + 70})

	// Generate the appearance and sign page 1.
	if err := appender.SignWithAppearance(1, field, annotator.FieldAppearance{}); err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Write output PDF file.
	err = appender.WriteToFile(outputPath)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	log.Printf("PDF file successfully signed. Output path: %s\n", outputPath)
}

// hasField reports whether the document form has a field named name.
func hasField(reader *model.PdfReader, name string) bool {
	if reader.AcroForm == nil {
		return false
	}
	for _, field := range reader.AcroForm.AllFields() {
		if field.PartialName() == name {
			return true
		}
	}
	return false
}

func generateKeys() (*rsa.PrivateKey, *x509.Certificate, error) {
	// Generate private key.
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	// Initialize X509 certificate template.
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "Jane Doe",
			Organization: []string{"Test Company"},
		},
		NotBefore: now.Add(-time.Hour),
		NotAfter:  now.Add(time.Hour * 24 * 365),

		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}

	// Generate X509 certificate.
	certData, err := x509.CreateCertificate(rand.Reader, &template, &template, priv.Public(), priv)
	if err != nil {
		return nil, nil, err
	}

	cert, err := x509.ParseCertificate(certData)
	if err != nil {
		return nil, nil, err
	}

	return priv, cert, nil
}
