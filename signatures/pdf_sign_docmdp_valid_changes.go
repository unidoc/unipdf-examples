/*
 * This example showcases how to sign a PDF document with the DocMDP restriction and add some valid changes.
 * The file is signed using a private/public key pair.
 *
 * $ ./pdf_sign_docmdp_valid_changes <INPUT_PDF_PATH> <OUTPUT_PDF_PATH>
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

	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/model/diffpolicy"
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
		log.Fatalln("Usage: go run pdf_sign_docmdp_valid_changes INPUT_PDF_PATH OUTPUT_PDF_PATH")
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

	totalPage, err := pdfReader.GetNumPages()
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Add signature and write it to output.pdf file.
	buf, err := addSignature(pdfReader, totalPage, outputPath)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	pdfReader2, err := model.NewPdfReader(bytes.NewReader(buf))
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	err = addSomeValidChanges(pdfReader2, outputPath)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	err = ValidateFile(outputPath)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	fmt.Println("Done")
}

func addSomeValidChanges(reader *model.PdfReader, outputPath string) error {
	// Create appender.
	appender, err := model.NewPdfAppender(reader)
	if err != nil {
		return err
	}

	page, err := reader.GetPage(1)
	if err != nil {
		return err
	}
	annotation := model.NewPdfAnnotationSquare()
	rect := model.PdfRectangle{Ury: 250.0, Urx: 150.0, Lly: 50.0, Llx: 50.0}
	annotation.Rect = rect.ToPdfObject()
	annotation.IC = core.MakeArrayFromFloats([]float64{4.0, 0.0, 0.3})
	annotation.CA = core.MakeFloat(0.5)

	page.AddAnnotation(annotation.PdfAnnotation)

	appender.UpdatePage(page)

	err = appender.WriteToFile(outputPath)
	if err != nil {
		return err
	}

	return nil
}

func ValidateFile(inputPath string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
	defer f.Close()

	reader, err := model.NewPdfReader(f)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	perms := reader.GetPerms()
	if perms == nil || perms.DocMDP == nil {
		return errors.New("unexpected perms object")
	}
	docMDPPerm, ok := perms.DocMDP.GetDocMDPPermission()
	if !ok {
		return errors.New("unexpected docMDP object")
	}

	inner_handler, _ := sighandler.NewAdobePKCS7Detached(nil, nil)
	handlerDocMdp, _ := sighandler.NewDocMDPHandler(inner_handler, docMDPPerm)

	handlers := []model.SignatureHandler{handlerDocMdp}

	res, err := reader.ValidateSignatures(handlers)
	if err != nil {
		return err
	}
	for i, validateResult := range res {
		log.Printf("== Signature %d", i+1)
		log.Printf("%s", validateResult.String())
	}
	return nil
}

func addSignature(reader *model.PdfReader, pageNum int, outputPath string) ([]byte, error) {
	// Create appender.
	appender, err := model.NewPdfAppender(reader)
	if err != nil {
		return nil, err
	}

	// Generate key pair.
	priv, cert, err := generateSigKeys()
	if err != nil {
		return nil, err
	}

	// Create signature handler.
	inner_handler, err := sighandler.NewAdobePKCS7Detached(priv, cert)
	if err != nil {
		return nil, err
	}

	handler, err := sighandler.NewDocMDPHandler(inner_handler, diffpolicy.FillingFormsAndAnnotations)
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
		return nil, err
	}

	sigField.T = core.MakeString("New Page Signature")

	if err = appender.Sign(pageNum, sigField); err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Write output PDF file.
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
