/*
 * This example showcases how to digitally sign a PDF file using an external
 * signing service which returns PKCS7 package. The PDF is first written with
 * an empty signature, the bytes covered by its ByteRange are sent to the
 * service, and the returned signature is placed into the Contents of that
 * same PDF. The external service is simulated locally with a self-signed
 * certificate.
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
	"log"
	"math/big"
	"os"
	"time"

	"github.com/unidoc/pkcs7"
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

// sigLen is the number of bytes reserved for the signature Contents.
const sigLen = 8192

const usagef = "Usage: %s INPUT_PDF_PATH OUTPUT_PDF_PATH\n"

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Printf(usagef, os.Args[0])
		return
	}
	inputPath := args[1]
	outputPath := args[2]

	// Generate PDF file signed with empty signature.
	handler, err := sighandler.NewEmptyAdobePKCS7Detached(sigLen)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	pdfData, signature, err := generateSignedFile(inputPath, handler)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Parse signature byte range.
	byteRange, err := parseByteRange(signature.ByteRange)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// The signed content is the whole file except the Contents hex string,
	// which sits between the two byte ranges.
	content := append(append([]byte{}, pdfData[byteRange[0]:byteRange[1]]...), pdfData[byteRange[2]:byteRange[3]]...)

	// This would be the time to send the content (or its SHA-256 digest) to a
	// signing device or signing web service and get back the PKCS7 signature.
	signatureData, err := getExternalSignature(content)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
	if len(signatureData) > sigLen {
		log.Fatalf("Fail: signature is %d bytes, but only %d bytes are reserved for Contents\n", len(signatureData), sigLen)
	}

	// Apply external signature to the PDF data buffer.
	// Overwrite the generated empty signature with the signature
	// bytes retrieved from the external service.
	sigBytes := make([]byte, sigLen)
	copy(sigBytes, signatureData)

	sig := core.MakeHexString(string(sigBytes)).Write()
	copy(pdfData[byteRange[1]:byteRange[2]], sig)

	// Write output file.
	if err := os.WriteFile(outputPath, pdfData, 0644); err != nil {
		log.Fatalf("Fail: %v\n", err)
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
			annotator.NewSignatureLine("Date", now.Format("2006.01.02")),
			annotator.NewSignatureLine("Reason", "External signature test"),
		},
		opts,
	)
	if err != nil {
		return nil, nil, err
	}
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

// getExternalSignature simulates an external service which signs the
// specified content with its own key and returns a detached PKCS7 signature.
func getExternalSignature(content []byte) ([]byte, error) {
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

	// Create a detached PKCS7 signature over the content.
	signedData, err := pkcs7.NewSignedData(content)
	if err != nil {
		return nil, err
	}
	signedData.SetDigestAlgorithm(pkcs7.OIDDigestAlgorithmSHA256)
	if err := signedData.AddSigner(cert, priv, pkcs7.SignerInfoConfig{}); err != nil {
		return nil, err
	}
	signedData.Detach()

	return signedData.Finish()
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
