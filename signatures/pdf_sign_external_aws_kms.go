/*
 * This example showcases how to digitally sign a PDF file using an external
 * signing service AWS KMS, you will need AWS KMS ASYMMETRIC KEY with SIGN AND VERIFY Key Usage.
 *
 * $ ./pdf_sign_external_aws_kms <INPUT_PDF_PATH> <OUTPUT_PDF_PATH> <KEY_ID>
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
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"

	"github.com/unidoc/pkcs7"
	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/model/sighandler"
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

const usagef = "Usage: %s INPUT_PDF_PATH OUTPUT_PDF_PATH KEY_ID\n"

func main() {
	args := os.Args
	if len(args) < 4 {
		fmt.Printf(usagef, os.Args[0])
		return
	}
	inputPath := args[1]
	outputPath := args[2]
	keyId := args[3]

	// Generate PDF file signed with empty signature.
	handler, err := sighandler.NewEmptyAdobePKCS7Detached(sigLen)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	_, signature, err := generateSignedFile(inputPath, handler)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	log.Println("Do sign contents")

	// Parse signature byte range.
	byteRange, err := parseByteRange(signature.ByteRange)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// This would be the time to send the PDF buffer to a signing device or
	// signing web service and get back the signature. We will simulate this by
	// signing the PDF using UniDoc and returning the signature data.
	signatureData, pdfData, err := getExternalSignatureAndSign(inputPath, keyId)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Apply external signature to the PDF data buffer.
	// Overwrite the generated empty signature with the signature
	// bytes retrieved from the external service.
	sigBytes := make([]byte, sigLen)
	copy(sigBytes, signatureData)

	sig := core.MakeHexString(string(sigBytes)).WriteString()
	copy(pdfData[byteRange[1]:byteRange[2]], []byte(sig))

	// Write output file.
	if err := ioutil.WriteFile(outputPath, pdfData, os.ModePerm); err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	log.Printf("PDF file successfully signed. Output path: %s\n", outputPath)
}

// generateSignedFile generates a signed version of the input PDF file using the
// specified signature handler.
func generateSignedFile(inputPath string, handler model.SignatureHandler) ([]byte, *model.PdfSignature, error) {
	// Create reader.
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader, err := model.NewPdfReader(file)
	if err != nil {
		return nil, nil, err
	}

	// Create appender.
	appender, err := model.NewPdfAppender(reader)
	if err != nil {
		return nil, nil, err
	}

	// Create signature.
	signature := model.NewPdfSignature(handler)
	signature.SetName("AWS KMS Sign")
	signature.SetReason("TestAwsKMSSignature")
	signature.SetDate(now, "")

	if err := signature.Initialize(); err != nil {
		return nil, nil, err
	}

	// Create signature field and appearance.
	opts := annotator.NewSignatureFieldOpts()
	opts.FontSize = 10
	opts.Rect = []float64{10, 25, 75, 60}

	field, err := annotator.NewSignatureField(
		signature,
		[]*annotator.SignatureLine{
			annotator.NewSignatureLine("Name", "John Doe"),
			annotator.NewSignatureLine("Date", "2023.01.01"),
			annotator.NewSignatureLine("Reason", "AWS KMS Sing Test"),
		},
		opts,
	)
	if err != nil {
		return nil, nil, err
	}

	field.T = core.MakeString("External signature")

	if err = appender.Sign(1, field); err != nil {
		return nil, nil, err
	}

	// Write PDF file to buffer.
	pdfBuf := bytes.NewBuffer(nil)
	if err = appender.Write(pdfBuf); err != nil {
		return nil, nil, err
	}

	return pdfBuf.Bytes(), signature, nil
}

// getExternalSignatureAndSign simulates an external service which signs the specified
// PDF file and returns its pdf signed bytes and signature.
func getExternalSignatureAndSign(inputPath string, keyId string) ([]byte, []byte, error) {
	awsKmsSign, err := AwsKmsExternalSigner(keyId)
	if err != nil {
		return nil, nil, err
	}

	now := time.Now()

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

	// We sign certificate using external signer that implement `crypto.Signer`.
	certData, err := x509.CreateCertificate(rand.Reader, &template, &template, awsKmsSign.getPublicKey(), awsKmsSign.signer)
	if err != nil {
		return nil, nil, err
	}

	cert, err := x509.ParseCertificate(certData)
	if err != nil {
		return nil, nil, err
	}

	certChain := []*x509.Certificate{cert}

	awsKmsSign.certChain = certChain

	// Sign input file.
	handler := awsKmsSign.toSigHandler()

	pdfBytes, signature, err := generateSignedFile(inputPath, handler)
	if err != nil {
		return nil, nil, err
	}

	return signature.Contents.Bytes(), pdfBytes, nil
}

// parseByteRange parses the ByteRange value of the signature field.
func parseByteRange(byteRange *core.PdfObjectArray) ([]int64, error) {
	if byteRange == nil {
		return nil, errors.New("byte range cannot be nil")
	}
	if byteRange.Len() != 4 {
		return nil, errors.New("invalid byte range length")
	}

	s1, err := core.GetNumberAsInt64(byteRange.Get(0))
	if err != nil {
		return nil, errors.New("invalid byte range value")
	}
	l1, err := core.GetNumberAsInt64(byteRange.Get(1))
	if err != nil {
		return nil, errors.New("invalid byte range value")
	}

	s2, err := core.GetNumberAsInt64(byteRange.Get(2))
	if err != nil {
		return nil, errors.New("invalid byte range value")
	}
	l2, err := core.GetNumberAsInt64(byteRange.Get(3))
	if err != nil {
		return nil, errors.New("invalid byte range value")
	}

	return []int64{s1, s1 + l1, s2, s2 + l2}, nil
}

// SignerCallback callback function of `crypto.Signer`
type SignerCallback func(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error)

// CryptoSigner Customize crypto signer wrapper to implement `crypto.Signer`.
type CryptoSigner struct {
	keyId  string
	client *kms.KMS

	callback SignerCallback
}

// EncryptionAlgorithmOID returns encryption OID algorithm supported by external services.
func (cs *CryptoSigner) EncryptionAlgorithmOID() asn1.ObjectIdentifier {
	return pkcs7.OIDEncryptionAlgorithmRSASHA256
}

// Public returns signer public key by calling to client services.
func (cs *CryptoSigner) Public() crypto.PublicKey {
	if cs.client == nil {
		return nil
	}

	publicKeyResp, err := cs.client.GetPublicKey(&kms.GetPublicKeyInput{
		KeyId: aws.String(cs.keyId),
	})
	if err != nil {
		return nil
	}

	spki, err := x509.ParsePKIXPublicKey(publicKeyResp.PublicKey)
	if err != nil {
		return nil
	}

	return spki
}

// Sign with customized `crypto.Signer`.
func (cs *CryptoSigner) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	if cs.callback != nil {
		return cs.callback(rand, digest, opts)
	}
	// Encrypt the data.
	result, err := cs.client.Sign(&kms.SignInput{
		KeyId:   aws.String(cs.keyId),
		Message: digest,
		// Set AWS `MessageType` to digest as we send the digest instead raw message.
		MessageType: aws.String("DIGEST"),
		// SigningAlgorithm must be match with digest algorithm.
		// Use:
		// - `crypto.SHA256` for `RSASSA_PKCS1_V1_5_SHA_256` and `RSASSA_PSS_SHA_256`
		// - `crypto.SHA384` for `RSASSA_PKCS1_V1_5_SHA_384` and `RSASSA_PSS_SHA_384`
		// - `crypto.SHA512` for `RSASSA_PKCS1_V1_5_SHA_512` and `RSASSA_PSS_SHA_512`
		SigningAlgorithm: aws.String("RSASSA_PKCS1_V1_5_SHA_256"),
	})
	if err != nil {
		return nil, err
	}

	return result.Signature, nil
}

func NewCryptoSigner(c *kms.KMS, keyId string, cb SignerCallback) crypto.Signer {
	return &CryptoSigner{
		keyId:    keyId,
		client:   c,
		callback: cb,
	}
}

// externalSigner is wrapper for third-party signer,
// needs to implement function of model.SignatureHandler.
// in this example, we use AWS KMS as third-party signer.
type externalSigner struct {
	certChain []*x509.Certificate

	// keyId AWS KMS ASYMMETRIC KEY_ID.
	keyId string
	// client is third-party api client.
	// We make client as interface for supporting another third-party signer.
	client *kms.KMS

	// signer customized `crypto.Signer`
	signer crypto.Signer

	ctx context.Context
}

// AwsKmsExternalSigner wrap externalSigner to model.SignatureHandler.
// With this, you can implement any third-party client signer into signature handler.
func AwsKmsExternalSigner(keyId string) (*externalSigner, error) {
	ctx := context.Background()

	// Initiate AWS session.
	sess, err := session.NewSession(&aws.Config{
		CredentialsChainVerboseErrors: aws.Bool(true),
		Region:                        aws.String("us-west-1"),
	})
	if err != nil {
		return nil, err
	}
	svc := kms.New(sess)

	// Generate external `crypto.Singer`.
	extSigner := NewCryptoSigner(svc, keyId, nil)

	return &externalSigner{
		keyId:  keyId,
		client: svc,
		ctx:    ctx,
		signer: extSigner,
	}, nil
}

// InitSignature sets the `model.PdfSignature` parameters.
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

	digest, err := handler.NewDigest(sig)
	if err != nil {
		return err
	}

	return es.Sign(sig, digest)
}

// Sign return error on failed sign.
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

	// Use customized callback if you want use different signature algorithm.
	/* cb := func(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
		* Sign digest.
		*	signResult, err := es.client.Sign(&kms.SignInput{
		*		KeyId:            aws.String(es.keyId),
		*		Message:          digest,
		*		MessageType:      aws.String("DIGEST"),
		*		SigningAlgorithm: aws.String("RSASSA_PKCS1_V1_5_SHA_256"),
		*	})
		*	if err != nil {
		*		return nil, err
		*	}
	  *
		*	return signResult.Signature, nil
		*}
	*/

	siConfig := pkcs7.SignerInfoConfig{}

	signer := es.signer
	// If want to use customized callback for signature.
	// signer = NewCryptoSigner(nil, "", cb)

	// Add certificate to `signedData`.
	err = signedData.AddSigner(es.certChain[0], signer, siConfig)
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

// Returns signer public key by calling to singer services.
func (es *externalSigner) getPublicKey() interface{} {
	return es.signer.Public()
}

// toSigHandler cast `externalSigner` to `model.SignatureHandler`.
func (es *externalSigner) toSigHandler() model.SignatureHandler {
	return es
}
