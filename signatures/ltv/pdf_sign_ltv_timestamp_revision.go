/*
 * This example showcases how to digitally sign a PDF file using a
 * PKCS12 (.p12/.pfx) file and LTV enable the signature by adding a second
 * revision to the document (containing the validation data) and a timestamp
 * signature (in order to protect the validation information).
 *
 * It needs network access: the timestamp is requested from the freetsa.org
 * TSA, and the validation data is fetched from the OCSP responders and CRL
 * distribution points named in the certificates.
 *
 * $ ./pdf_sign_ltv_timestamp_revision <FILE.p12> <P12_PASS> <INPUT_PDF_PATH> <OUTPUT_PDF_PATH> [<EXTRA_CERTS.pem>]
 */

package main

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"slices"
	"time"

	"software.sslmate.com/src/go-pkcs12"

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

const usagef = "Usage: %s P12_FILE PASSWORD INPUT_PDF_PATH OUTPUT_PDF_PATH [EXTRA_CERTS.pem]\n"

func main() {
	args := os.Args
	if len(args) < 5 || len(args) > 6 {
		fmt.Printf(usagef, os.Args[0])
		return
	}
	p12Path := args[1]
	password := args[2]
	inputPath := args[3]
	outputPath := args[4]

	// Load private key and X509 certificate from the PKCS12 file.
	pfxData, err := os.ReadFile(p12Path)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	key, cert, caCerts, err := pkcs12.DecodeChain(pfxData, password)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
	priv, ok := key.(*rsa.PrivateKey)
	if !ok {
		log.Fatalf("Fail: the PKCS12 private key is %T; this example needs an RSA key\n", key)
	}

	// Load certificate chain: the signer, the CA certificates bundled in the
	// PKCS12 file, then any certificates from EXTRA_CERTS.pem.
	certChain := append([]*x509.Certificate{cert}, caCerts...)
	if len(args) == 6 {
		issuerCertData, err := os.ReadFile(args[5])
		if err != nil {
			log.Fatalf("Fail: %v\n", err)
		}

		for len(issuerCertData) != 0 {
			var block *pem.Block
			block, issuerCertData = pem.Decode(issuerCertData)
			if block == nil {
				break
			}

			issuer, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				log.Fatalf("Fail: %v\n", err)
			}
			if !slices.ContainsFunc(certChain, issuer.Equal) {
				certChain = append(certChain, issuer)
			}
		}
	}

	// Sign and write file to buffer.
	signedBytes, err := signFile(inputPath, priv, cert)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// LTV enable, timestamp and write file to buffer.
	signedBytes, err = ltvEnableAndTimestamp(bytes.NewReader(signedBytes), certChain)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Optionally, LTV enable the applied timestamp signature.
	signedBytes, err = ltvEnableTimestampSig(bytes.NewReader(signedBytes))
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Write output file to disk.
	err = os.WriteFile(outputPath, signedBytes, 0644)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	log.Printf("PDF file successfully signed. Output path: %s\n", outputPath)
}

func signFile(inputPath string, priv *rsa.PrivateKey, cert *x509.Certificate) ([]byte, error) {
	// Create reader.
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader, err := model.NewPdfReader(file)
	if err != nil {
		return nil, err
	}

	// Create appender.
	appender, err := model.NewPdfAppender(reader)
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
	signature.SetName("Test Sign LTV enable")
	signature.SetReason("TestSignLTV")
	signature.SetDate(time.Now(), "")

	if err := signature.Initialize(); err != nil {
		return nil, err
	}

	// Create signature field and appearance.
	opts := annotator.NewSignatureFieldOpts()
	opts.FontSize = 10
	opts.Rect = []float64{10, 25, 75, 60}

	field, err := annotator.NewSignatureField(
		signature,
		[]*annotator.SignatureLine{
			annotator.NewSignatureLine("Name", "John Doe"),
			annotator.NewSignatureLine("Date", time.Now().Format("2006.01.02")),
			annotator.NewSignatureLine("Reason", "Signature test"),
		},
		opts,
	)
	if err != nil {
		return nil, err
	}
	field.T = core.MakeString("Test Sign LTV enable")

	if err = appender.Sign(1, field); err != nil {
		return nil, err
	}

	// Write output PDF file to buffer.
	buf := bytes.NewBuffer(nil)
	if err = appender.Write(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ltvEnableAndTimestamp takes a *bytes.Reader rather than an io.ReadSeeker
// because the reader and appender require an io.ReaderAt source.
func ltvEnableAndTimestamp(r *bytes.Reader, certChain []*x509.Certificate) ([]byte, error) {
	// Create reader.
	reader, err := model.NewPdfReader(r)
	if err != nil {
		return nil, err
	}

	// Create appender.
	appender, err := model.NewPdfAppender(reader)
	if err != nil {
		return nil, err
	}

	// LTV enable the certificate chain used to apply the signature.
	ltv, err := model.NewLTV(appender)
	if err != nil {
		return nil, err
	}

	if err := enableSignedFields(reader, ltv, certChain); err != nil {
		return nil, err
	}

	// Create timestamp handler.
	handler, err := sighandler.NewDocTimeStamp("https://freetsa.org/tsr", crypto.SHA512)
	if err != nil {
		return nil, err
	}

	// Create signature field and appearance.
	signature := model.NewPdfSignature(handler)
	signature.SetName("Test Sign Timestamp")
	signature.SetReason("TestSignTimestamp")

	if err := signature.Initialize(); err != nil {
		return nil, err
	}

	sigField := model.NewPdfFieldSignature(signature)
	sigField.T = core.MakeString("Test Sign Timestamp")
	sigField.Rect = core.MakeArray(
		core.MakeInteger(0),
		core.MakeInteger(0),
		core.MakeInteger(0),
		core.MakeInteger(0),
	)

	if err = appender.Sign(1, sigField); err != nil {
		return nil, err
	}

	// Write output PDF file to buffer.
	buf := bytes.NewBuffer(nil)
	if err = appender.Write(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func ltvEnableTimestampSig(r *bytes.Reader) ([]byte, error) {
	// Create reader.
	reader, err := model.NewPdfReader(r)
	if err != nil {
		return nil, err
	}

	// Create appender.
	appender, err := model.NewPdfAppender(reader)
	if err != nil {
		return nil, err
	}

	// LTV enable the certificate chain used to apply the signature.
	ltv, err := model.NewLTV(appender)
	if err != nil {
		return nil, err
	}

	if err := enableSignedFields(reader, ltv, nil); err != nil {
		return nil, err
	}

	// Write output PDF file to buffer.
	buf := bytes.NewBuffer(nil)
	if err = appender.Write(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// enableSignedFields LTV enables every signature field that holds a
// signature. LTV.EnableAll fails on an empty signature field.
func enableSignedFields(reader *model.PdfReader, ltv *model.LTV, extraCerts []*x509.Certificate) error {
	if reader.AcroForm == nil {
		return nil
	}
	for _, field := range reader.AcroForm.AllFields() {
		sigField, ok := field.GetContext().(*model.PdfFieldSignature)
		if !ok || sigField.V == nil || sigField.V.Contents == nil {
			continue
		}
		if err := ltv.Enable(sigField.V, extraCerts); err != nil {
			return err
		}
	}
	return nil
}
