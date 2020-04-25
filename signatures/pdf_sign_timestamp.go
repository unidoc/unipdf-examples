/*
 * This example showcases how to digitally sign and timestamp a PDF file using a
 * PKCS12 (.p12/.pfx) file.
 *
 * $ ./pdf_sign_timestamp <FILE.p12> <PASSWORD> <INPUT_PDF_PATH> <OUTPUT_PDF_PATH>
 */
package main

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/model/sighandler"
	"golang.org/x/crypto/pkcs12"
)

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
			annotator.NewSignatureLine("Date", "2019.16.04"),
			annotator.NewSignatureLine("Reason", "External signature test"),
		},
		opts,
	)
	field.T = core.MakeString("Self signed PDF")

	if err = appender.Sign(1, field); err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	outDoc := bytes.NewBuffer(nil)
	if err = appender.Write(outDoc); err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	pdf1, err := model.NewPdfReader(bytes.NewReader(outDoc.Bytes()))
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	appender, err = model.NewPdfAppender(pdf1)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	handler, err = sighandler.NewDocTimeStamp("https://freetsa.org/tsr", crypto.SHA512)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	// Create signature field and appearance.
	signature = model.NewPdfSignature(handler)
	signature.SetName("Test Appender")
	signature.SetReason("TestAppenderSignPage4")
	signature.SetDate(time.Now(), "")

	if err := signature.Initialize(); err != nil {
		return
	}

	sigField := model.NewPdfFieldSignature(signature)
	sigField.T = core.MakeString("Signature1")
	sigField.Rect = core.MakeArray(
		core.MakeInteger(0),
		core.MakeInteger(0),
		core.MakeInteger(0),
		core.MakeInteger(0),
	)

	if err = appender.Sign(1, sigField); err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	err = appender.WriteToFile(outputPath)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	log.Printf("PDF file successfully signed. Output path: %s\n", outputPath)
}
