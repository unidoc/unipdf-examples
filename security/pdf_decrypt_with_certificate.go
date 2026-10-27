/*
 * Decrypts a PDF file encrypted for recipients with the public-key security
 * handler (Adobe.PubSec), see pdf_encrypt_for_recipients.go, and writes an
 * unencrypted copy. The file is decrypted with the recipient's certificate and
 * private key using PdfReader.DecryptWithCertificate. Only a recipient with all
 * permissions (owner-equivalent) may remove the encryption.
 *
 * The certificate is read from a PEM or DER file, the private key from a PEM file
 * (PKCS #1 or PKCS #8). The private key must be an RSA key. For a .p12/.pfx digital
 * ID, get both with pkcs12.DecodeChain as in signatures/signing/pdf_sign_pkcs12.go.
 *
 * Run as: go run pdf_decrypt_with_certificate.go input.pdf recipient.crt recipient.key output.pdf
 */

package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v5/common/license"
	"github.com/unidoc/unipdf/v5/core"
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
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run pdf_decrypt_with_certificate.go input.pdf recipient.crt recipient.key output.pdf")
		fmt.Println("Decrypts a PDF encrypted for recipients, using a recipient's certificate and private key.")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	certPath := os.Args[2]
	keyPath := os.Args[3]
	outputPath := os.Args[4]

	err := decryptPdf(inputPath, certPath, keyPath, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		switch {
		case errors.Is(err, security.ErrNotRecipient):
			fmt.Println("The PDF was not encrypted for this certificate.")
		case errors.Is(err, security.ErrRecipientNotOpened):
			fmt.Println("The PDF was encrypted for this certificate, but the private key does not open it. Does the key belong to the certificate?")
		case errors.Is(err, security.ErrUnsupportedRecipient):
			fmt.Println("The PDF has recipient entries in a format that is not supported (e.g. RSAES-OAEP), so it is unknown whether this certificate is among them.")
		case errors.Is(err, core.ErrPasswordRequired):
			fmt.Println("The PDF is password protected, see pdf_unlock.go.")
		}
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

func decryptPdf(inputPath, certPath, keyPath, outputPath string) error {
	cert, err := loadCertificate(certPath)
	if err != nil {
		return err
	}
	key, err := loadPrivateKey(keyPath)
	if err != nil {
		return err
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
		// The unencrypted copy carries no permissions, so only an owner-equivalent
		// recipient may make it. A password protected PDF is rejected with
		// core.ErrPasswordRequired.
		_, perms, err := pdfReader.CheckAccessRightsWithCertificate(cert, key)
		if err != nil {
			return err
		}
		if perms != security.PermOwner {
			return errors.New("recipient is not allowed to remove the encryption")
		}

		// Decrypt with the recipient's certificate and private key.
		_, err = pdfReader.DecryptWithCertificate(cert, key)
		if err != nil {
			return err
		}
	}

	// Generate a PdfWriter instance from existing PdfReader.
	pdfWriter, err := pdfReader.ToWriter(nil)
	if err != nil {
		return err
	}

	// Write to file.
	return pdfWriter.WriteToFile(outputPath)
}

// loadCertificate reads an X509 certificate from a PEM or DER file.
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
	return cert, nil
}

// loadPrivateKey reads an RSA private key from a PKCS #1 or PKCS #8 PEM file.
func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("%s: no PEM data found", path)
	}

	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("%s: private key is %T, only RSA keys are supported", path, key)
	}
	return rsaKey, nil
}
