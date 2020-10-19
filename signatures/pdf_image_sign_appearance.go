/*
 * This example showcases how to create appearance fields for digital signature.
 * The file is signed using a private/public key pair.
 *
 * $ ./pdf_image_sign_appearance <INPUT_PDF_PATH> <IMAGE_FILE> <WATERMARK_IMAGE_FILE> <OUTPUT_PDF_PATH>
 */
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"image"
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

var now = time.Now()

const usage = "Usage: %s INPUT_PDF_PATH IMAGE_FILE WATERMARK_IMAGE_FILE OUTPUT_PDF_PATH\n"

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
	if len(args) < 5 {
		fmt.Printf(usage, os.Args[0])
		return
	}
	inputPath := args[1]
	imageFile := args[2]
	watermarkImageFile := args[3]
	outputPath := args[4]

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

	// Create the image
	imgFile, err := os.Open(imageFile)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
	defer imgFile.Close()

	signatureImage, _, err := image.Decode(imgFile)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Create the watermark image
	wImgFile, err := os.Open(watermarkImageFile)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
	defer imgFile.Close()

	signatureWImage, _, err := image.Decode(wImgFile)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
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
	signature.SetReason("Test Signature Appearance Reason")
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

		// Image on the right
		opts := annotator.NewSignatureFieldOpts()
		opts.FontSize = 10
		opts.Rect = []float64{10, 25, 110, 75}
		opts.Image = signatureImage
		opts.ImagePosition = annotator.SignatureImageRight
		opts.WatermarkImage = signatureWImage

		sigField, err := annotator.NewSignatureField(
			signature,
			[]*annotator.SignatureLine{
				annotator.NewSignatureLine("Name", "Jane Doe"),
				annotator.NewSignatureLine("Date", "2019.01.03"),
				annotator.NewSignatureLine("Reason", "Image on right"),
				annotator.NewSignatureLine("Location", "New York"),
				annotator.NewSignatureLine("DN", "authority1:name1"),
			},
			opts,
		)
		sigField.T = core.MakeString(fmt.Sprintf("Signature %d", pageNum))

		if err = appender.Sign(pageNum, sigField); err != nil {
			log.Fatal("Fail: %v\n", err)
		}

		// Without image signature
		opts = annotator.NewSignatureFieldOpts()
		opts.FontSize = 10
		opts.Rect = []float64{170, 25, 270, 75}
		opts.WatermarkImage = signatureWImage

		sigField, err = annotator.NewSignatureField(
			signature,
			[]*annotator.SignatureLine{
				annotator.NewSignatureLine("Name", "Jane Doe"),
				annotator.NewSignatureLine("Date", "2019.01.03"),
				annotator.NewSignatureLine("Reason", "Without image sig"),
				annotator.NewSignatureLine("Location", "New York"),
				annotator.NewSignatureLine("DN", "authority1:name1"),
			},
			opts,
		)
		sigField.T = core.MakeString(fmt.Sprintf("Signature2 %d", pageNum))

		if err = appender.Sign(pageNum, sigField); err != nil {
			log.Fatal("Fail: %v\n", err)
		}

		// Image on the top
		opts = annotator.NewSignatureFieldOpts()
		opts.FontSize = 8
		opts.Rect = []float64{10, 90, 110, 140}
		opts.TextColor = model.NewPdfColorDeviceRGB(255, 0, 0)
		opts.Image = signatureImage
		opts.ImagePosition = annotator.SignatureImageTop
		opts.WatermarkImage = signatureWImage

		sigField, err = annotator.NewSignatureField(
			signature,
			[]*annotator.SignatureLine{
				annotator.NewSignatureLine("Name", "John Doe"),
				annotator.NewSignatureLine("Date", "2019.03.14"),
				annotator.NewSignatureLine("Reason", "Image on top"),
				annotator.NewSignatureLine("Location", "London"),
				annotator.NewSignatureLine("DN", "authority2:name2"),
			},
			opts,
		)
		sigField.T = core.MakeString(fmt.Sprintf("Signature3 %d", pageNum))

		if err = appender.Sign(pageNum, sigField); err != nil {
			log.Fatalf("Fail: %v\n", err)
		}

		// Image on bottom
		opts = annotator.NewSignatureFieldOpts()
		opts.BorderSize = 1
		opts.FontSize = 10
		opts.Rect = []float64{170, 90, 270, 140}
		opts.FillColor = model.NewPdfColorDeviceRGB(255, 255, 0)
		opts.TextColor = model.NewPdfColorDeviceRGB(0, 0, 200)
		opts.Image = signatureImage
		opts.ImagePosition = annotator.SignatureImageBottom
		opts.WatermarkImage = signatureWImage

		sigField, err = annotator.NewSignatureField(
			signature,
			[]*annotator.SignatureLine{
				annotator.NewSignatureLine("Name", "John Smith"),
				annotator.NewSignatureLine("Date", "2019.02.19"),
				annotator.NewSignatureLine("Reason", "Image on bottom"),
				annotator.NewSignatureLine("Location", "Paris"),
				annotator.NewSignatureLine("DN", "authority3:name3"),
			},
			opts,
		)
		sigField.T = core.MakeString(fmt.Sprintf("Signature4 %d", pageNum))

		if err = appender.Sign(pageNum, sigField); err != nil {
			log.Fatalf("Fail: %v\n", err)
		}

		// Only Image Signature
		opts = annotator.NewSignatureFieldOpts()
		opts.Rect = []float64{70, 50, 170, 100}
		opts.Image = signatureImage
		opts.WatermarkImage = signatureWImage

		sigField, err = annotator.NewSignatureField(signature, nil, opts)

		sigField.T = core.MakeString(fmt.Sprintf("Signature5 %d", pageNum))

		if err = appender.Sign(pageNum, sigField); err != nil {
			log.Fatalf("Fail: %v\n", err)
		}
	}

	// Write output PDF file.
	err = appender.WriteToFile(outputPath)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
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
