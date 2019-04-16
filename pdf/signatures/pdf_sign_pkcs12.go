/*
 * This example showcases how to digitally sign a PDF file using a
 * PKCS12 (.p12/.pfx) file.
 *
 * $ ./pdf_sign_external <FILE.p12> <PASSWORD> <INPUT_PDF_PATH> <OUTPUT_PDF_PATH>
 */
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/unidoc/unidoc/pdf/annotator"
	"github.com/unidoc/unidoc/pdf/core"
	"github.com/unidoc/unidoc/pdf/model"
	"github.com/unidoc/unidoc/pdf/model/sighandler"
	"golang.org/x/crypto/pkcs12"
)

var now = time.Now()

const usage = "Usage: %s P12_FILE PASSWORD INPUT_PDF_PATH OUTPUT_PDF_PATH\n"

func main() {
	args := os.Args
	if len(args) < 4 {
		fmt.Printf(usage, os.Args[0])
		return
	}
	p12Path := args[1]
	password := args[2]
	inputPath := args[3]
	outputPath := args[4]

	// Get private key and X509 certificate from the P12 file.
	pfxData, err := ioutil.ReadFile(p12Path)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	priv, cert, err := pkcs12.Decode(pfxData, password)
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
	handler, err := sighandler.NewAdobePKCS7Detached(priv.(*rsa.PrivateKey), cert)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	// Create signature.
	signature := model.NewPdfSignature(handler)
	signature.SetName("Test Self Signed PDF")
	signature.SetReason("TestSelfSignedPDF")
	signature.SetDate(now, "")

	if err := signature.Initialize(); err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	// Create signature field and appearance.
	opts := annotator.NewSignatureFieldOpts()
	opts.FontSize = 10
	opts.Rect = []float64{10, 25, 75, 60}

	field, err := annotator.NewSignatureField(
		signature,
		[]*annotator.SignatureLine{
			annotator.NewSignatureLine("Name", "John Doe"),
			annotator.NewSignatureLine("Date", "2019.16.04"),
			annotator.NewSignatureLine("Reason", "External signature test"),
		},
		opts,
	)
	field.T = core.MakeString("Self signed PDF")

	if err = appender.Sign(1, field); err != nil {
		log.Fatal("Fail: %v\n", err)
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
