/*
 * This example showcases how to retrieve the specific or previous revision of a PDF document
 * The file is signed using a private/public key pair.
 *
 * $ ./pdf_sign_docmdp <INPUT_PDF_PATH> <OUTPUT_PDF_PATH>
 */
package main

import (
	"bytes"
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
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/model/sighandler"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}
}

func main() {
	args := os.Args
	if len(args) < 3 {
		log.Fatalln("Usage: go run pdf_sign_get_revision INPUT_PDF_PATH OUTPUT_PDF_PATH")
	}

	inputPath := args[1]
	outputPath := args[2]

	// Read the original file.
	f, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Preparing first revision
	// Add signature
	buf, err := addSignature(pdfReader, 0)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Open reader for the signed document
	pdfReader, err = model.NewPdfReader(bytes.NewReader(buf))
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Preparing second revision
	// Add signature
	buf, err = addSignature(pdfReader, 1)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	pdfReader2, err := model.NewPdfReader(bytes.NewReader(buf))
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	err = retrieveSpecificRevision(pdfReader2, 0, outputPath)
	// or
	//err = retrievePreviousRevision(pdfReader2, outputPath)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	fmt.Println("Done")
}

func retrievePreviousRevision(reader *model.PdfReader, outputPath string) error {
	prevReader, err := reader.GetPreviousRevision()
	if err != nil {
		return err
	}

	appender, err := model.NewPdfAppender(prevReader)
	if err != nil {
		return err
	}
	// Write output PDF file.
	if err = appender.WriteToFile(outputPath); err != nil {
		return err
	}

	log.Printf("PDF file successfully signed. Output path: %s\n", outputPath)
	return nil
}

func retrieveSpecificRevision(reader *model.PdfReader, revisionNumber int, outputPath string) error {
	prevReader, err := reader.GetRevision(revisionNumber)
	if err != nil {
		return err
	}

	appender, err := model.NewPdfAppender(prevReader)
	if err != nil {
		return err
	}
	// Write output PDF file.
	if err = appender.WriteToFile(outputPath); err != nil {
		return err
	}

	log.Printf("PDF file successfully signed. Output path: %s\n", outputPath)
	return nil
}

func addSignature(pdfReader *model.PdfReader, signNumber int) ([]byte, error) {
	totalPage, err := pdfReader.GetNumPages()
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Create appender.
	appender, err := model.NewPdfAppender(pdfReader)
	if err != nil {
		return nil, err
	}

	// Generate key pair.
	priv, cert, err := generateSigKeys()
	if err != nil {
		return nil, err
	}

	// Create signature handler.
	handler, err := sighandler.NewAdobePKCS7Detached(priv, cert)
	if err != nil {
		return nil, err
	}

	// Create signature.
	signature := model.NewPdfSignature(handler)
	signature.SetName("Test Signature Appearance Name")
	signature.SetReason("TestSignatureAppearance Reason")
	signature.SetDate(time.Now(), "")

	// Initialize signature.
	if err := signature.Initialize(); err != nil {
		return nil, err
	}

	opts := annotator.NewSignatureFieldOpts()
	opts.FontSize = 8
	opts.Rect = []float64{float64(50 + signNumber*100), 250, float64(150 + signNumber*100), 300}
	opts.TextColor = model.NewPdfColorDeviceRGB(255, 0, 0)

	sigField, err := annotator.NewSignatureField(
		signature,
		[]*annotator.SignatureLine{
			annotator.NewSignatureLine("Name", "John Doe"),
			annotator.NewSignatureLine("Date", "2019.03.14"),
			annotator.NewSignatureLine("Reason", fmt.Sprintf("Test sign #%d", signNumber)),
			annotator.NewSignatureLine("Location", "London"),
			annotator.NewSignatureLine("DN", "authority2:name2"),
		},
		opts,
	)
	if err != nil {
		return nil, err
	}

	sigField.T = core.MakeString(fmt.Sprintf("New Page Signature %d", signNumber))

	if err = appender.Sign(totalPage, sigField); err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	buf := &bytes.Buffer{}
	// Write output PDF file.
	if err = appender.Write(buf); err != nil {
		return nil, err
	}

	log.Println("PDF file successfully signed")

	return buf.Bytes(), nil
}

func generateSigKeys() (*rsa.PrivateKey, *x509.Certificate, error) {
	var now = time.Now()

	// Generate private key.
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	// Initialize X509 certificate template.
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "any",
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
