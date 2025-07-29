/*
 * This example showcases how to create an adbe.pkcs7.detached compatible digital signature for a PDF file
 * using ECDSA (elliptic curve DSA) key. ECDSA keys can be used for signing PDF files starting from version 2.0.
 *
 * $ ./pdf_sign_pkcs7_ecdsa <FILE.PFX> <PASSWORD> <INPUT_PDF_PATH> <OUTPUT_PDF_PATH>
 */
package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/unidoc/unipdf/v4/annotator"
	"github.com/unidoc/unipdf/v4/common/license"
	"github.com/unidoc/unipdf/v4/core"
	"github.com/unidoc/unipdf/v4/model"
	"github.com/unidoc/unipdf/v4/model/sighandler"

	"software.sslmate.com/src/go-pkcs12"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}
}

const usagef = "Usage: %s PFX_FILE PASSWORD INPUT_PDF_PATH OUTPUT_PDF_PATH\n"

func main() {
	args := os.Args
	if len(args) < 5 {
		fmt.Printf(usagef, os.Args[0])
		return
	}
	pfxPath := args[1]
	password := args[2]
	inputPath := args[3]
	outputPath := args[4]

	// Get private key and X509 certificate from the ECDSA PFX file.
	pfxData, err := os.ReadFile(pfxPath)
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

	// Create signature handler with ECDSA key.
	handler, err := sighandler.NewAdobePKCS7DetachedEcdsa(priv.(*ecdsa.PrivateKey), cert)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	// Create signature.
	signature := model.NewPdfSignature(handler)
	signature.SetName("adbe.pkcs7.detached ECDSA PDF")
	signature.SetReason("Test ECDSA")
	signature.SetDate(time.Now(), "")

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
			annotator.NewSignatureLine("Date", "2025.07.28"),
			annotator.NewSignatureLine("Reason", "Test"),
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
