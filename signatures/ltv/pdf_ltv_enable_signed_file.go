/*
 * This example showcases how to LTV enable the signatures in a signed PDF file,
 * by adding a second revision to the document, containing the validation data.
 *
 * $ ./pdf_ltv_enable_signed_file <INPUT_PDF_PATH> <OUTPUT_PDF_PATH> [<EXTRA_CERTS.pem>]
 */

package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/model"
)

const licenseKey = `
-----BEGIN UNIDOC LICENSE KEY-----
Free trial license keys are available at: https://unidoc.io/
-----END UNIDOC LICENSE KEY-----
`

const usage = "Usage: %s INPUT_PDF_PATH OUTPUT_PDF_PATH [EXTRA_CERTS]\n"

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
		fmt.Printf(usage, os.Args[0])
		return
	}
	inputPath := args[1]
	outputPath := args[2]

	// Load certificate chain.
	var certChain []*x509.Certificate
	if len(args) == 4 {
		issuerCertData, err := ioutil.ReadFile(args[3])
		if err != nil {
			log.Fatal("Fail: %v\n", err)
		}

		for len(issuerCertData) != 0 {
			var block *pem.Block
			block, issuerCertData = pem.Decode(issuerCertData)
			if block == nil {
				break
			}

			issuer, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				log.Fatal("Fail: %v\n", err)
			}
			certChain = append(certChain, issuer)
		}
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

	// LTV enable the signed file.
	ltv, err := model.NewLTV(appender)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	if err := ltv.EnableAll(certChain); err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	// Write output PDF file.
	err = appender.WriteToFile(outputPath)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	log.Printf("PDF file successfully LTV enabled. Output path: %s\n", outputPath)
}
