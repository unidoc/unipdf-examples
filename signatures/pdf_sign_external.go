/*
 * This example showcases how to digitally sign a PDF file using an external
 * signing service which returns PKCS7 package. The external service is
 * is simulated by signing the file with UniDoc.
 *
 * $ ./pdf_sign_external <INPUT_PDF_PATH> <OUTPUT_PDF_PATH>
 */
package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"io/ioutil"
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

	// Generate PDF file signed with empty signature.
	handler, err := sighandler.NewEmptyAdobePKCS7Detached(8192)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	pdfData, signature, err := generateSignedFile(inputPath, handler)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	// Parse signature byte range.
	byteRange, err := parseByteRange(signature.ByteRange)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	// This would be the time to send the PDF buffer to a signing device or
	// signing web service and get back the signature. We will simulate this by
	// signing the PDF using UniDoc and returning the signature data.
	signatureData, err := getExternalSignature(inputPath)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	// Apply external signature to the PDF data buffer.
	// Overwrite the generated empty signature with the signature
	// bytes retrieved from the external service.
	sigBytes := make([]byte, 8192)
	copy(sigBytes, signatureData)

	sig := core.MakeHexString(string(sigBytes)).WriteString()
	copy(pdfData[byteRange[1]:byteRange[2]], []byte(sig))

	// Write output file.
	if err := ioutil.WriteFile(outputPath, pdfData, os.ModePerm); err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	log.Printf("PDF file successfully signed. Output path: %s\n", outputPath)
}

// generateSignedFile generates a signed version of the input PDF file using the
// specified signature handler.
func generateSignedFile(inputPath string, handler model.SignatureHandler) ([]byte, *model.PdfSignature, error) {
	// Create reader.
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader, err := model.NewPdfReader(file)
	if err != nil {
		return nil, nil, err
	}

	// Create appender.
	appender, err := model.NewPdfAppender(reader)
	if err != nil {
		return nil, nil, err
	}

	// Create signature.
	signature := model.NewPdfSignature(handler)
	signature.SetName("Test External Signature")
	signature.SetReason("TestAppenderExternalSignature")
	signature.SetDate(now, "")

	if err := signature.Initialize(); err != nil {
		return nil, nil, err
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
	field.T = core.MakeString("External signature")

	if err = appender.Sign(1, field); err != nil {
		return nil, nil, err
	}

	// Write PDF file to buffer.
	pdfBuf := bytes.NewBuffer(nil)
	if err = appender.Write(pdfBuf); err != nil {
		return nil, nil, err
	}

	return pdfBuf.Bytes(), signature, nil
}

// getExternalSignature simulates an external service which signs the specified
// PDF file and returns its signature.
func getExternalSignature(inputPath string) ([]byte, error) {
	// Generate private key.
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	cert, err := x509.ParseCertificate(certData)
	if err != nil {
		return nil, err
	}

	// Sign input file.
	handler, err := sighandler.NewAdobePKCS7Detached(priv, cert)
	if err != nil {
		return nil, err
	}

	_, signature, err := generateSignedFile(inputPath, handler)
	if err != nil {
		return nil, err
	}

	return signature.Contents.Bytes(), nil
}

// parseByteRange parses the ByteRange value of the signature field.
func parseByteRange(byteRange *core.PdfObjectArray) ([]int64, error) {
	if byteRange == nil {
		return nil, errors.New("byte range cannot be nil")
	}
	if byteRange.Len() != 4 {
		return nil, errors.New("invalid byte range length")
	}

	s1, err := core.GetNumberAsInt64(byteRange.Get(0))
	if err != nil {
		return nil, errors.New("invalid byte range value")
	}
	l1, err := core.GetNumberAsInt64(byteRange.Get(1))
	if err != nil {
		return nil, errors.New("invalid byte range value")
	}

	s2, err := core.GetNumberAsInt64(byteRange.Get(2))
	if err != nil {
		return nil, errors.New("invalid byte range value")
	}
	l2, err := core.GetNumberAsInt64(byteRange.Get(3))
	if err != nil {
		return nil, errors.New("invalid byte range value")
	}

	return []int64{s1, s1 + l1, s2, s2 + l2}, nil
}
