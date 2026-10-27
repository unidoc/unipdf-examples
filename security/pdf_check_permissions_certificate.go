/*
 * Checks whether a recipient's certificate and private key can open a PDF file
 * encrypted with the public-key security handler (Adobe.PubSec), see
 * pdf_encrypt_for_recipients.go, and which permissions the file grants the
 * recipient. The check does not decrypt the file. If the recipient can open it,
 * the example then decrypts it with PdfReader.DecryptWithCertificate and prints
 * the number of pages.
 *
 * The certificate is read from a PEM or DER file, the private key from a PEM file
 * (PKCS #1 or PKCS #8). The private key must be an RSA key.
 *
 * Run as: go run pdf_check_permissions_certificate.go input.pdf recipient.crt recipient.key
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
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run pdf_check_permissions_certificate.go input.pdf recipient.crt recipient.key")
		fmt.Println("Checks the access permissions a PDF encrypted for recipients grants a recipient.")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	certPath := os.Args[2]
	keyPath := os.Args[3]

	err := printAccessInfo(inputPath, certPath, keyPath)
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
			fmt.Println("The PDF is password protected, see pdf_check_permissions.go.")
		}
		os.Exit(1)
	}
}

func printAccessInfo(inputPath, certPath, keyPath string) error {
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
	if !isEncrypted {
		fmt.Printf("%s - is not encrypted\n", inputPath)
		return nil
	}

	// Check the access rights without decrypting the document. A certificate
	// that cannot open it is reported as an error, see main.
	_, perms, err := pdfReader.CheckAccessRightsWithCertificate(cert, key)
	if err != nil {
		return err
	}

	fmt.Printf("Input file %s\n", inputPath)
	fmt.Printf("Access Permissions: %+v\n", perms)
	fmt.Printf("--------\n")

	// Print a text summary of the flags.
	booltext := map[bool]string{false: "No", true: "Yes"}
	allowed := func(p security.Permissions) string {
		return booltext[perms.Allowed(p)]
	}

	fmt.Printf("Owner-equivalent (all permissions)? - %s\n", booltext[perms == security.PermOwner])
	fmt.Printf("Printing allowed? - %s\n", allowed(security.PermPrinting))
	if perms.Allowed(security.PermPrinting) {
		fmt.Printf("Full print quality (otherwise print in low res)? - %s\n", allowed(security.PermFullPrintQuality))
	}
	fmt.Printf("Modifications allowed? - %s\n", allowed(security.PermModify))
	fmt.Printf("Assemble (insert, rotate, delete pages) allowed? - %s\n", allowed(security.PermRotateInsert))
	fmt.Printf("Allow extracting graphics? %s\n", allowed(security.PermExtractGraphics))
	fmt.Printf("Can annotate? - %s\n", allowed(security.PermAnnotate))
	if perms.Allowed(security.PermAnnotate) {
		fmt.Printf("Can fill forms? - Yes\n")
	} else {
		fmt.Printf("Can fill forms? - %s\n", allowed(security.PermFillForms))
	}
	fmt.Printf("Extract text, graphics for users with disabilities? - %s\n", allowed(security.PermDisabilityExtract))
	fmt.Printf("--------\n")

	// The same certificate and key now decrypt the document on this reader.
	_, err = pdfReader.DecryptWithCertificate(cert, key)
	if err != nil {
		return err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}
	fmt.Printf("Decrypted, the document has %d page(s)\n", numPages)

	return nil
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
