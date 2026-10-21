/*
 * This example showcases how to digitally sign a PDF file using an external
 * signing service Google Cloud KMS. The key must have purpose
 * ASYMMETRIC_SIGN and one of the algorithms RSA_SIGN_PKCS1_2048_SHA256,
 * RSA_SIGN_PKCS1_3072_SHA256 or RSA_SIGN_PKCS1_4096_SHA256.
 *
 * $ ./pdf_sign_external_google_cloud_kms <INPUT_PDF_PATH> <OUTPUT_PDF_PATH> <GOOGLE_CLOUD_CREDENTIALS_PATH> <KEY_NAME>
 *
 * <GOOGLE_CLOUD_CREDENTIALS_PATH> is a service account JSON key file.
 * <KEY_NAME> format example: projects/my-project/locations/us-east1/keyRings/my-key-ring/cryptoKeys/my-key/cryptoKeyVersions/123
 */
package main

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"math/big"
	"os"
	"time"

	kms "cloud.google.com/go/kms/apiv1"
	kmspb "cloud.google.com/go/kms/apiv1/kmspb"
	gcOption "google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/unidoc/pkcs7"
	"github.com/unidoc/unipdf/v5/annotator"
	"github.com/unidoc/unipdf/v5/common/license"
	"github.com/unidoc/unipdf/v5/core"
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

var (
	now = time.Now()

	// sigLen signature length.
	sigLen = 8192
)

const usagef = "Usage: %s INPUT_PDF_PATH OUTPUT_PDF_PATH GOOGLE_CLOUD_CREDENTIALS_PATH KEY_NAME\n"

func main() {
	args := os.Args
	if len(args) < 5 {
		fmt.Printf(usagef, os.Args[0])
		return
	}
	inputPath := args[1]
	outputPath := args[2]
	gcCredentialsPath := args[3]
	keyName := args[4]

	// Initiate Google Cloud KMS client.
	client, err := kms.NewKeyManagementClient(context.Background(), gcOption.WithCredentialsFile(gcCredentialsPath))
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
	defer func() { _ = client.Close() }()

	handler, err := GcKmsExternalSigner(client, keyName)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	pdfData, err := generateSignedFile(inputPath, handler)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Write output file.
	if err := os.WriteFile(outputPath, pdfData, 0644); err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	log.Printf("PDF file successfully signed. Output path: %s\n", outputPath)
}

// generateSignedFile signs the input PDF file with the specified signature
// handler. The appender reserves the Contents placeholder, computes the
// ByteRange content and passes it to the handler's Sign method.
func generateSignedFile(inputPath string, handler model.SignatureHandler) ([]byte, error) {
	// Create reader.
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	reader, err := model.NewPdfReader(file)
	if err != nil {
		return nil, err
	}

	// Create appender.
	appender, err := model.NewPdfAppender(reader)
	if err != nil {
		return nil, err
	}

	// Create signature.
	signature := model.NewPdfSignature(handler)
	signature.SetName("GOOGLECLOUD KMS Sign")
	signature.SetReason("TestGCKMSSignature")
	signature.SetDate(now, "")

	if err := signature.Initialize(); err != nil {
		return nil, err
	}

	// Create signature field and appearance.
	opts := annotator.NewSignatureFieldOpts()
	opts.FontSize = 10
	opts.Rect = []float64{10, 25, 75, 60}

	field, err := annotator.NewSignatureField(
		signature,
		[]*annotator.SignatureLine{
			annotator.NewSignatureLine("Name", "John Doe"),
			annotator.NewSignatureLine("Date", now.Format("2006.01.02")),
			annotator.NewSignatureLine("Reason", "GC KMS Sign Test"),
		},
		opts,
	)
	if err != nil {
		return nil, err
	}

	field.T = core.MakeString("External signature")

	if err = appender.Sign(1, field); err != nil {
		return nil, err
	}

	// Write PDF file to buffer.
	pdfBuf := bytes.NewBuffer(nil)
	if err = appender.Write(pdfBuf); err != nil {
		return nil, err
	}

	return pdfBuf.Bytes(), nil
}

// crc32c computes the CRC32C checksum Cloud KMS uses for in-transit
// integrity checks, see https://cloud.google.com/kms/docs/data-integrity-guidelines
func crc32c(data []byte) uint32 {
	return crc32.Checksum(data, crc32.MakeTable(crc32.Castagnoli))
}

// CryptoSigner implements `crypto.Signer` on top of a Google Cloud KMS key.
type CryptoSigner struct {
	keyName   string
	client    *kms.KeyManagementClient
	publicKey crypto.PublicKey
}

// NewCryptoSigner fetches the public key of the KMS key version and returns a
// signer for it. Only RSA PKCS#1 v1.5 SHA-256 keys are supported.
func NewCryptoSigner(c *kms.KeyManagementClient, keyName string) (*CryptoSigner, error) {
	result, err := c.GetPublicKey(context.Background(), &kmspb.GetPublicKeyRequest{
		Name: keyName,
	})
	if err != nil {
		return nil, fmt.Errorf("get KMS public key: %w", err)
	}

	switch result.Algorithm {
	case kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_2048_SHA256,
		kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_3072_SHA256,
		kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_4096_SHA256:
	default:
		return nil, fmt.Errorf("KMS key %s uses %s, this example requires RSA_SIGN_PKCS1_2048_SHA256, RSA_SIGN_PKCS1_3072_SHA256 or RSA_SIGN_PKCS1_4096_SHA256", keyName, result.Algorithm)
	}

	key := []byte(result.Pem)
	if int64(crc32c(key)) != result.PemCrc32C.GetValue() {
		return nil, errors.New("GetPublicKey: response corrupted in-transit")
	}

	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("GetPublicKey: response contains no PEM block")
	}
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &CryptoSigner{
		keyName:   keyName,
		client:    c,
		publicKey: publicKey,
	}, nil
}

// EncryptionAlgorithmOID returns encryption OID algorithm supported by external services.
func (cs *CryptoSigner) EncryptionAlgorithmOID() asn1.ObjectIdentifier {
	return pkcs7.OIDEncryptionAlgorithmRSASHA256
}

// Public returns the signer public key.
func (cs *CryptoSigner) Public() crypto.PublicKey {
	return cs.publicKey
}

// Sign signs the SHA-256 digest with the KMS key.
func (cs *CryptoSigner) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	req := &kmspb.AsymmetricSignRequest{
		Name: cs.keyName,
		Digest: &kmspb.Digest{
			Digest: &kmspb.Digest_Sha256{
				Sha256: digest,
			},
		},
		DigestCrc32C: wrapperspb.Int64(int64(crc32c(digest))),
	}

	result, err := cs.client.AsymmetricSign(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to sign digest: %v", err)
	}

	if !result.VerifiedDigestCrc32C {
		return nil, errors.New("AsymmetricSign: request corrupted in-transit")
	}

	if int64(crc32c(result.Signature)) != result.SignatureCrc32C.GetValue() {
		return nil, errors.New("AsymmetricSign: response corrupted in-transit")
	}

	return result.Signature, nil
}

// externalSigner is wrapper for third-party signer,
// needs to implement function of model.SignatureHandler.
// in this example, we use Google Cloud KMS as third-party signer.
type externalSigner struct {
	certChain []*x509.Certificate

	// signer customized `crypto.Signer`
	signer crypto.Signer
}

// GcKmsExternalSigner wrap externalSigner to model.SignatureHandler.
// With this, you can implement any third-party client signer into signature handler.
func GcKmsExternalSigner(client *kms.KeyManagementClient, keyName string) (*externalSigner, error) {
	signer, err := NewCryptoSigner(client, keyName)
	if err != nil {
		return nil, err
	}

	// Initialize X509 certificate template.
	template := x509.Certificate{
		SerialNumber: new(big.Int),
		Subject: pkix.Name{
			CommonName:   "UniDOC",
			Organization: []string{"Test Company"},
		},
		NotBefore:          now.Add(-time.Hour),
		NotAfter:           now.Add(time.Hour * 24 * 365),
		PublicKeyAlgorithm: x509.RSA,
		SignatureAlgorithm: x509.SHA256WithRSA,
		KeyUsage:           x509.KeyUsageDigitalSignature,
	}

	// Self-sign the certificate with the KMS key.
	certData, err := x509.CreateCertificate(rand.Reader, &template, &template, signer.Public(), signer)
	if err != nil {
		return nil, err
	}

	cert, err := x509.ParseCertificate(certData)
	if err != nil {
		return nil, err
	}

	return &externalSigner{
		certChain: []*x509.Certificate{cert},
		signer:    signer,
	}, nil
}

// InitSignature sets the `model.PdfSignature` parameters and reserves the
// Contents placeholder. No KMS call is made here.
func (es *externalSigner) InitSignature(sig *model.PdfSignature) error {
	// Create PDF array object which will contain the certificate chain data.
	pdfCerts := core.MakeArray()
	for _, cert := range es.certChain {
		pdfCerts.Append(core.MakeString(string(cert.Raw)))
	}

	// Append cert to signature.
	sig.Cert = pdfCerts

	handler := *es
	sig.Handler = &handler
	sig.Filter = core.MakeName("Adobe.PPKLite")
	sig.SubFilter = core.MakeName("adbe.pkcs7.detached")
	sig.Reference = nil

	return es.Sign(sig, nil)
}

// Sign signs the ByteRange content passed in digest with the KMS key. A nil
// digest only reserves the Contents placeholder.
func (es *externalSigner) Sign(sig *model.PdfSignature, digest model.Hasher) error {
	if digest == nil {
		sig.Contents = core.MakeHexString(string(make([]byte, sigLen)))
		return nil
	}

	buffer := digest.(*bytes.Buffer)
	signedData, err := pkcs7.NewSignedData(buffer.Bytes())
	if err != nil {
		return err
	}

	// Set digest algorithm which supported by Google Cloud KMS.
	signedData.SetDigestAlgorithm(pkcs7.OIDDigestAlgorithmSHA256)

	// Add certificate to `signedData`.
	err = signedData.AddSigner(es.certChain[0], es.signer, pkcs7.SignerInfoConfig{})
	if err != nil {
		return err
	}

	// Call Detach() is you want to remove content from the signature
	// and generate an S/MIME detached signature.
	signedData.Detach()
	// Finish() to obtain the signature bytes.
	detachedSignature, err := signedData.Finish()
	if err != nil {
		return err
	}
	if len(detachedSignature) > sigLen {
		return fmt.Errorf("signature is %d bytes, but only %d bytes are reserved for Contents", len(detachedSignature), sigLen)
	}

	data := make([]byte, sigLen)
	copy(data, detachedSignature)

	sig.Contents = core.MakeHexString(string(data))
	return nil
}

// IsApplicable returns true if the signature handler is applicable for the PdfSignature.
func (es *externalSigner) IsApplicable(sig *model.PdfSignature) bool {
	if sig == nil || sig.Filter == nil || sig.SubFilter == nil {
		return false
	}

	return (*sig.Filter == "Adobe.PPKMS" || *sig.Filter == "Adobe.PPKLite") && *sig.SubFilter == "adbe.pkcs7.detached"
}

// NewDigest creates a new digest.
func (es *externalSigner) NewDigest(sig *model.PdfSignature) (model.Hasher, error) {
	return bytes.NewBuffer(nil), nil
}

// Validate validates PdfSignature.
func (es *externalSigner) Validate(sig *model.PdfSignature, digest model.Hasher) (model.SignatureValidationResult, error) {
	return model.SignatureValidationResult{
		IsSigned:   true,
		IsVerified: true,
	}, nil
}
