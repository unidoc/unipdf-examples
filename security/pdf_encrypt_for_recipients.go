/*
 * Encrypts a PDF file for one or more recipients with the public-key security
 * handler (Adobe.PubSec). There is no password: every recipient opens the file
 * with their own certificate and private key, see pdf_decrypt_with_certificate.go
 * and pdf_check_permissions_certificate.go.
 *
 * Recipient certificates are read from PEM or DER files (.crt, .cer, .pem) and
 * must have RSA public keys. Test identities can be created with
 * pdf_generate_recipient_keys.go.
 *
 * By default the file is encrypted with AES-256 (PDF 2.0) and the recipients get
 * all permissions (security.PermOwner), which makes them owner-equivalent: they
 * may also change the security settings of the file.
 *   -aes128      encrypts with AES-128 instead, for viewers without PDF 2.0 support.
 *   -print-only  lets the recipients view and print the file, but not modify it or
 *                copy its content. Extraction for accessibility is always granted
 *                by the public-key handler.
 * The permissions are the same for all recipients of the file.
 *
 * Run as: go run pdf_encrypt_for_recipients.go [-aes128] [-print-only] input.pdf output.pdf recipient.crt [recipient2.crt ...]
 */

package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v5/common/license"
	"github.com/unidoc/unipdf/v5/core/security"
	"github.com/unidoc/unipdf/v5/model"
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
	aes128 := flag.Bool("aes128", false, "encrypt with AES-128 instead of AES-256 (PDF 2.0)")
	printOnly := flag.Bool("print-only", false, "allow the recipients to view and print only")
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintln(out, "Usage: go run pdf_encrypt_for_recipients.go [-aes128] [-print-only] input.pdf output.pdf recipient.crt [recipient2.crt ...]")
		fmt.Fprintln(out, "Encrypts the PDF for the recipient certificates. Options go before the files.")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 3 {
		flag.Usage()
		os.Exit(1)
	}

	inputPath := args[0]
	outputPath := args[1]
	certPaths := args[2:]

	err := encryptForRecipients(inputPath, outputPath, certPaths, *aes128, *printOnly)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

func encryptForRecipients(inputPath, outputPath string, certPaths []string, aes128, printOnly bool) error {
	var recipients []*x509.Certificate
	for _, path := range certPaths {
		cert, err := loadCertificate(path)
		if err != nil {
			return err
		}
		recipients = append(recipients, cert)
	}

	// Passing nil options has the same effect as these defaults.
	encryptOptions := &model.EncryptOptions{
		Algorithm:   model.AES_256bit,
		Permissions: security.PermOwner, // All permissions: recipients are owner-equivalent.
	}
	if aes128 {
		encryptOptions.Algorithm = model.AES_128bit
	}
	if printOnly {
		encryptOptions.Permissions = security.PermPrinting | // Allow printing with low quality
			security.PermFullPrintQuality // and with full quality.
	}

	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return err
	}
	if isEncrypted {
		return fmt.Errorf("The PDF is already encrypted (need to decrypt first)")
	}

	// Generate a PdfWriter instance from existing PdfReader.
	pdfWriter, err := pdfReader.ToWriter(nil)
	if err != nil {
		return err
	}

	// Encrypt document for the recipients before writing to file.
	err = pdfWriter.EncryptForRecipients(recipients, encryptOptions)
	if err != nil {
		return err
	}

	// Write to file.
	return pdfWriter.WriteToFile(outputPath)
}

// loadCertificate reads an X509 certificate from a PEM or DER file and checks
// that it has an RSA public key, the only kind a recipient can have.
func loadCertificate(path string) (*x509.Certificate, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if block, _ := pem.Decode(data); block != nil {
		data = block.Bytes
	}

	cert, err := x509.ParseCertificate(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if _, ok := cert.PublicKey.(*rsa.PublicKey); !ok {
		return nil, fmt.Errorf("%s: public key is %T, recipients must have RSA keys", path, cert.PublicKey)
	}
	return cert, nil
}
