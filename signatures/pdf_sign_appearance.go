/*
 * This example showcases how to create appearance fields for digital signature.
 * The file is signed using a private/public key pair.
 *
 * $ ./pdf_sign_appearance <INPUT_PDF_PATH> <OUTPUT_PDF_PATH>
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

	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/model/sighandler"
)

var now = time.Now()

const usage = "Usage: %s INPUT_PDF_PATH OUTPUT_PDF_PATH\n"

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Printf(usage, os.Args[0])
		return
	}
	inputPath := args[1]
	outputPath := args[2]

	// Generate key pair.
	priv, cert, err := generateKeys()
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

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

	// Create appender.
	appender, err := model.NewPdfAppender(reader)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	// Create signature handler.
	handler, err := sighandler.NewAdobePKCS7Detached(priv, cert)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	// Create signature.
	signature := model.NewPdfSignature(handler)
	signature.SetName("Test Signature Appearance Name")
	signature.SetReason("TestSignatureAppearance Reason")
	signature.SetDate(time.Now(), "")

	if err := signature.Initialize(); err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	numPages, err := reader.GetNumPages()
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	// Create signature fields and add them on each page of the PDF file.
	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		// Annotation 1.
		opts := annotator.NewSignatureFieldOpts()
		opts.FontSize = 10
		opts.Rect = []float64{10, 25, 75, 60}

		sigField, err := annotator.NewSignatureField(
			signature,
			[]*annotator.SignatureLine{
				annotator.NewSignatureLine("Name", "Jane Doe"),
				annotator.NewSignatureLine("Date", "2019.01.03"),
				annotator.NewSignatureLine("Reason", "Some reason"),
				annotator.NewSignatureLine("Location", "New York"),
				annotator.NewSignatureLine("DN", "authority1:name1"),
			},
			opts,
		)
		sigField.T = core.MakeString(fmt.Sprintf("Signature %d", pageNum))

		if err = appender.Sign(pageNum, sigField); err != nil {
			log.Fatal("Fail: %v\n", err)
		}

		// Annotation 2.
		opts = annotator.NewSignatureFieldOpts()
		opts.FontSize = 8
		opts.Rect = []float64{250, 25, 325, 70}
		opts.TextColor = model.NewPdfColorDeviceRGB(255, 0, 0)

		sigField, err = annotator.NewSignatureField(
			signature,
			[]*annotator.SignatureLine{
				annotator.NewSignatureLine("Name", "John Doe"),
				annotator.NewSignatureLine("Date", "2019.03.14"),
				annotator.NewSignatureLine("Reason", "No reason"),
				annotator.NewSignatureLine("Location", "London"),
				annotator.NewSignatureLine("DN", "authority2:name2"),
			},
			opts,
		)
		sigField.T = core.MakeString(fmt.Sprintf("Signature2 %d", pageNum))

		if err = appender.Sign(pageNum, sigField); err != nil {
			log.Fatalf("Fail: %v\n", err)
		}

		// Annotation 3.
		opts = annotator.NewSignatureFieldOpts()
		opts.BorderSize = 1
		opts.FontSize = 10
		opts.Rect = []float64{475, 25, 590, 80}
		opts.FillColor = model.NewPdfColorDeviceRGB(255, 255, 0)
		opts.TextColor = model.NewPdfColorDeviceRGB(0, 0, 200)

		sigField, err = annotator.NewSignatureField(
			signature,
			[]*annotator.SignatureLine{
				annotator.NewSignatureLine("Name", "John Smith"),
				annotator.NewSignatureLine("Date", "2019.02.19"),
				annotator.NewSignatureLine("Reason", "Another reason"),
				annotator.NewSignatureLine("Location", "Paris"),
				annotator.NewSignatureLine("DN", "authority3:name3"),
			},
			opts,
		)
		sigField.T = core.MakeString(fmt.Sprintf("Signature3 %d", pageNum))

		if err = appender.Sign(pageNum, sigField); err != nil {
			log.Fatalf("Fail: %v\n", err)
		}
	}

	// Write output PDF file.
	err = appender.WriteToFile(outputPath)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	log.Printf("PDF file successfully signed. Output path: %s\n", outputPath)
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
			Organization: []string{"Test Company"},
		},
		NotBefore: now.Add(-time.Hour),
		NotAfter:  now.Add(time.Hour * 24 * 365),

		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
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
