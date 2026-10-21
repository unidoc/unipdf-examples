/*
 * This example showcases how to digitally sign a PDF file using an external
 * signing service AWS KMS. The key must be an asymmetric KMS key with key
 * usage SIGN_VERIFY and an RSA key spec (RSA_2048, RSA_3072 or RSA_4096);
 * it is used with the RSASSA_PKCS1_V1_5_SHA_256 signing algorithm.
 *
 * Credentials and region come from the standard AWS configuration: the
 * AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY/AWS_PROFILE and AWS_REGION
 * environment variables, or ~/.aws/credentials and ~/.aws/config.
 *
 * $ ./pdf_sign_external_aws_kms <INPUT_PDF_PATH> <OUTPUT_PDF_PATH> <KEY_ID>
 */
package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"

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

const usagef = "Usage: %s INPUT_PDF_PATH OUTPUT_PDF_PATH KEY_ID\n" +
	"AWS credentials and region (AWS_REGION) are read from the standard AWS configuration.\n"

func main() {
	args := os.Args
	if len(args) < 4 {
		fmt.Printf(usagef, os.Args[0])
		return
	}
	inputPath := args[1]
	outputPath := args[2]
	keyId := args[3]

	handler, err := AwsKmsExternalSigner(keyId)
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
	defer file.Close()

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
	signature.SetName("AWS KMS Sign")
	signature.SetReason("TestAwsKMSSignature")
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
			annotator.NewSignatureLine("Reason", "AWS KMS Sign Test"),
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

// CryptoSigner implements `crypto.Signer` on top of an AWS KMS key.
type CryptoSigner struct {
	keyId     string
	client    *kms.KMS
	publicKey *rsa.PublicKey
}

// NewCryptoSigner fetches the public key of the KMS key and returns a signer
// for it. Only RSA keys are supported.
func NewCryptoSigner(c *kms.KMS, keyId string) (*CryptoSigner, error) {
	publicKeyResp, err := c.GetPublicKey(&kms.GetPublicKeyInput{
		KeyId: aws.String(keyId),
	})
	if err != nil {
		return nil, fmt.Errorf("get KMS public key: %w", err)
	}

	publicKey, err := x509.ParsePKIXPublicKey(publicKeyResp.PublicKey)
	if err != nil {
		return nil, err
	}

	rsaKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("KMS key %s has a %T public key, this example requires an RSA key spec (RSA_2048, RSA_3072 or RSA_4096)", keyId, publicKey)
	}

	return &CryptoSigner{
		keyId:     keyId,
		client:    c,
		publicKey: rsaKey,
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
	result, err := cs.client.Sign(&kms.SignInput{
		KeyId:   aws.String(cs.keyId),
		Message: digest,
		// Set AWS `MessageType` to digest as we send the digest instead raw message.
		MessageType: aws.String("DIGEST"),
		// Must match the SHA-256 digest algorithm set on the PKCS7 package.
		SigningAlgorithm: aws.String("RSASSA_PKCS1_V1_5_SHA_256"),
	})
	if err != nil {
		return nil, err
	}

	return result.Signature, nil
}

// externalSigner is wrapper for third-party signer,
// needs to implement function of model.SignatureHandler.
// in this example, we use AWS KMS as third-party signer.
type externalSigner struct {
	certChain []*x509.Certificate

	// signer customized `crypto.Signer`
	signer crypto.Signer
}

// AwsKmsExternalSigner wrap externalSigner to model.SignatureHandler.
// With this, you can implement any third-party client signer into signature handler.
func AwsKmsExternalSigner(keyId string) (*externalSigner, error) {
	// Initiate AWS session from the shared config files and environment.
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, err
	}
	if aws.StringValue(sess.Config.Region) == "" {
		return nil, errors.New("AWS region is not set, set AWS_REGION or a region in ~/.aws/config")
	}

	signer, err := NewCryptoSigner(kms.New(sess), keyId)
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

	// Set digest algorithm which supported by AWS KMS.
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
