/*
 * Generates self-signed RSA test identities for the public-key encryption
 * examples (pdf_encrypt_for_recipients.go, pdf_decrypt_with_certificate.go and
 * pdf_check_permissions_certificate.go). For every NAME it writes two PEM files:
 *   NAME.crt - the certificate, which you give to whoever encrypts documents for you;
 *   NAME.key - the private key, which only the recipient keeps.
 *
 * The identities are for trying out the examples only: the private key is
 * written unencrypted and the certificate is not issued by a CA.
 *
 * Run as: go run pdf_generate_recipient_keys.go NAME [NAME ...]
 * E.g.:   go run pdf_generate_recipient_keys.go alice bob
 */

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run pdf_generate_recipient_keys.go NAME [NAME ...]")
		fmt.Println("Writes a test certificate NAME.crt and private key NAME.key for every NAME.")
		os.Exit(1)
	}

	for _, name := range os.Args[1:] {
		err := generateIdentity(name)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Wrote %s.crt and %s.key\n", name, name)
	}
}

func generateIdentity(name string) error {
	certPath := name + ".crt"
	keyPath := name + ".key"

	// Never replace an existing identity: PDFs already encrypted for it could
	// no longer be opened.
	for _, path := range []string{certPath, keyPath} {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("%s already exists", path)
		}
	}

	// Recipients must have RSA keys.
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// A recipient is identified in the encrypted file by the issuer and serial
	// number of its certificate, so every identity gets its own name and a
	// random serial number.
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 127))
	if err != nil {
		return err
	}

	now := time.Now()
	template := x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   name,
			Organization: []string{"Test Company"},
		},
		NotBefore: now.Add(-time.Hour).UTC(),
		NotAfter:  now.AddDate(1, 0, 0).UTC(),
		KeyUsage: x509.KeyUsageKeyEncipherment |
			x509.KeyUsageDataEncipherment |
			x509.KeyUsageDigitalSignature,
	}

	// Self-signed X509 certificate.
	certData, err := x509.CreateCertificate(rand.Reader, &template, &template, priv.Public(), priv)
	if err != nil {
		return err
	}

	keyData, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return err
	}

	// Write the private key first, so that a certificate is never left
	// without its key.
	err = writeNewFile(keyPath, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyData}), 0600)
	if err != nil {
		return err
	}
	return writeNewFile(certPath, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certData}), 0644)
}

// writeNewFile writes data to a file that must not exist yet.
func writeNewFile(path string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, perm)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	return err
}
