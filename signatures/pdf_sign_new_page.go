/*
 * This example showcases how to append a new page with signature to a PDF document.
 * The file is signed using a private/public key pair.
 *
 * $ ./pdf_sign_new_page <INPUT_PDF_PATH> <OUTPUT_PDF_PATH>
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

const licenseKey = `
-----BEGIN UNIDOC LICENSE KEY-----
Free trial license keys are available at: https://unidoc.io/
-----END UNIDOC LICENSE KEY-----
`

func init() {
	// Enable debug-level logging.
	// unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	err := license.SetLicenseKey(licenseKey, `Company Name`)
	if err != nil {
		panic(err)
	}
}

func main() {
	args := os.Args
	if len(args) < 3 {
		log.Fatalln("Usage: go run pdf_sign_new_page INPUT_PDF_PATH OUTPUT_PDF_PATH")
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

	// Add new page and write it into a buffer.
	buf, err := addPage(pdfReader)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	pdfReader, err = model.NewPdfReader(bytes.NewReader(buf.Bytes()))
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	totalPage, err := pdfReader.GetNumPages()
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Add signature and write it to output.pdf file.
	err = addSignature(pdfReader, totalPage, outputPath)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	fmt.Println("Done")
}

func addPage(reader *model.PdfReader) (*bytes.Buffer, error) {
	writer, err := reader.ToWriter(&model.ReaderToWriterOpts{})
	if err != nil {
		return nil, err
	}

	if err = writer.AddPage(model.NewPdfPage()); err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	if err = writer.Write(buf); err != nil {
		return nil, err
	}

	return buf, nil
}

func addSignature(reader *model.PdfReader, pageNum int, outputPath string) error {
	// Create appender.
	appender, err := model.NewPdfAppender(reader)
	if err != nil {
		return err
	}

	// Generate key pair.
	priv, cert, err := generateSigKeys()
	if err != nil {
		return err
	}

	// Create signature handler.
	handler, err := sighandler.NewAdobePKCS7Detached(priv, cert)
	if err != nil {
		return err
	}

	// Create signature.
	signature := model.NewPdfSignature(handler)
	signature.SetName("Test Signature Appearance Name")
	signature.SetReason("TestSignatureAppearance Reason")
	signature.SetDate(time.Now(), "")

	// Initialize signature.
	if err := signature.Initialize(); err != nil {
		return err
	}

	opts := annotator.NewSignatureFieldOpts()
	opts.FontSize = 8
	opts.Rect = []float64{250, 250, 325, 300}
	opts.TextColor = model.NewPdfColorDeviceRGB(255, 0, 0)

	sigField, err := annotator.NewSignatureField(
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
	if err != nil {
		return err
	}

	sigField.T = core.MakeString("New Page Signature")

	if err = appender.Sign(pageNum, sigField); err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Write output PDF file.
	if err = appender.WriteToFile(outputPath); err != nil {
		return err
	}

	log.Printf("PDF file successfully signed. Output path: %s\n", outputPath)

	return nil
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
