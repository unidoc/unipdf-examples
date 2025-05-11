/*
 * This example showcases how to digitally sign a PDF file using HSM via PKCS11.
 * with UniDoc.
 *
 * To create a key pair:
 * $ ./pdf_sign_hsm_pkcs11_cgo add test <PIN> <keypair_label>
 *
 * To sign a PDF:
 * $ ./pdf_sign_hsm_pkcs11_cgo sign test <PIN> <keypair_label> input.pdf input_signed.pdf
 *
 * See instructions for testing via SoftHSM in README.md.
 */
package main

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"hash"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ThalesIgnite/crypto11"

	"github.com/unidoc/unipdf/v4/annotator"
	"github.com/unidoc/unipdf/v4/common/license"
	"github.com/unidoc/unipdf/v4/core"
	"github.com/unidoc/unipdf/v4/model"
	"github.com/unidoc/unipdf/v4/model/sighandler"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}
}

// Library path might be different on different operating systems.
const PathSoftHSM = "/usr/local/lib/softhsm/libsofthsm2.so"

const (
	usage     = "Usage: %s add|sign PARAMETERS...\n"
	usageAdd  = "Usage: %s add TOKEN_LABEL TOKEN_PIN KEYPAIR_LABEL\n"
	usageSign = "Usage: %s sign TOKEN_LABEL TOKEN_PIN KEYPAIR_LABEL INPUT_PDF_PATH OUTPUT_PDF_PATH\n"
)

func main() {
	// Check specified action.
	args := os.Args

	lenArgs := len(args)
	if lenArgs < 2 {
		fmt.Printf(usage, os.Args[0])
		return
	}

	action := args[1]
	switch action {
	case "add":
		if lenArgs != 5 {
			fmt.Printf(usageAdd, os.Args[0])
			return
		}
	case "sign":
		if lenArgs != 7 {
			fmt.Printf(usageSign, os.Args[0])
			return
		}
	default:
		fmt.Printf(usage, os.Args[0])
		return
	}

	tokenLabel := args[2]
	tokenPin := args[3]
	keypairLabel := args[4]

	// Initialize PKCS11 session.
	// The PKCS11 store only exposes a crypto.Signer interface.
	// The signing process takes place inside the signer and it is only
	// possible while a session is open.
	ctx, err := initPKCS11Session(tokenLabel, tokenPin)
	if err != nil {
		log.Fatalf("Init Fail: %v\n", err)
	}
	defer ctx.Close()

	switch action {
	case "add":
		if _, err := addKeyPair(ctx, keypairLabel); err != nil {
			log.Fatalf("Add Key Pair Fail: %v\n", err)
		}

		log.Printf("Key pair successfully added for token %s\n", tokenLabel)
	case "sign":
		// Get a signer object associated with the keypair label.
		signer, err := getKeyPair(ctx, keypairLabel)
		if err != nil {
			log.Fatalf("Fail: %v\n", err)
		}

		if signer == nil {
			log.Printf("Fail: KeyPair with the specified label not found")
			return
		}

		cert, err := generateCertificate(signer)
		if err != nil {
			log.Fatalf("Fail: %v\n", err)
		}

		inputPath := args[5]
		outputPath := args[6]

		if err := sign(signer, cert, inputPath, outputPath); err != nil {
			log.Fatalf("Fail: %v\n", err)
		}

		log.Printf("PDF file successfully signed. Output path: %s\n", outputPath)
	}
}

// initPKCS11Session initializes a PKCS11 store and creates a new session.
func initPKCS11Session(tokenLabel, tokenPin string) (*crypto11.Context, error) {
	conf := crypto11.Config{
		Path:       PathSoftHSM,
		TokenLabel: tokenLabel,
		Pin:        tokenPin,
	}

	return crypto11.Configure(&conf)
}

// getKeyPair retrieves the key pair with the specified label.
func getKeyPair(ctx *crypto11.Context, keypairLabel string) (crypto11.Signer, error) {
	return ctx.FindKeyPair(nil, []byte(keypairLabel))
}

// addKeyPair adds a new public/private key pair with the specified label
// to the PKCS11 store.
func addKeyPair(ctx *crypto11.Context, keypairLabel string) (crypto11.SignerDecrypter, error) {
	// Generate random keypair id.
	resultId := make([]byte, 32)
	rand.Read(resultId)

	// Generate key pair.
	return ctx.GenerateRSAKeyPairWithLabel([]byte(resultId), []byte(keypairLabel), 2048)
}

// generateCertificate generates a X509 certificate based on the specified private key.
func generateCertificate(signer crypto.Signer) (*x509.Certificate, error) {
	// Initialize X509 certificate template.
	now := time.Now()
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test Company"},
		},
		NotBefore: now,
		NotAfter:  now.Add(time.Hour * 24 * 365),

		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Generate X509 certificate.
	certData, err := x509.CreateCertificate(rand.Reader, template, template, signer.Public(), signer)
	if err != nil {
		return nil, err
	}

	return x509.ParseCertificate(certData)
}

// sign signs the specified input PDF file using an adobeX509RSASHA1 signature handler
// and saves the result at the destination specified by the outputPath parameter.
func sign(signer crypto.Signer, certificate *x509.Certificate, inputPath, outputPath string) error {
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

	// Create signature function.
	signFunc := func(sig *model.PdfSignature, digest model.Hasher) ([]byte, error) {
		h, ok := digest.(hash.Hash)
		if !ok {
			return nil, errors.New("hash type error")
		}

		return signer.Sign(rand.Reader, h.Sum(nil), crypto.SHA1)
	}

	// Create custom signature handler.
	handler, err := sighandler.NewAdobeX509RSASHA1Custom(certificate, signFunc)
	if err != nil {
		return err
	}

	// Create signature.
	signature := model.NewPdfSignature(handler)
	signature.SetName("Test SoftHSM2 Signature")
	signature.SetReason("TestSoftHSM2Signature")
	signature.SetDate(time.Now(), "")

	if err := signature.Initialize(); err != nil {
		return err
	}

	// Create signature field and appearance.
	opts := annotator.NewSignatureFieldOpts()
	opts.FontSize = 10
	opts.Rect = []float64{10, 25, 75, 60}

	sigField, err := annotator.NewSignatureField(
		signature,
		[]*annotator.SignatureLine{
			annotator.NewSignatureLine("Name", "John Doe"),
			annotator.NewSignatureLine("Date", "2019.15.03"),
			annotator.NewSignatureLine("Reason", "SoftHSM2 Signature Test"),
		},
		opts,
	)
	sigField.T = core.MakeString("External signature")

	// Sign PDF.
	if err = appender.Sign(1, sigField); err != nil {
		return err
	}

	// Write output file.
	return appender.WriteToFile(outputPath)
}
