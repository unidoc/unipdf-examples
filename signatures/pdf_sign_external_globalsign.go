/*
 * This example showcases how to digitally sign a PDF file using the GlobalSign DSS API.
 *
 * $ ./pdf_sign_external_globalsign <INPUT_PDF_PATH> <OUTPUT_PDF_PATH> <API_KEY> <API_SECRET> <CERT_FILE_PATH> <KEY_FILE_PATH>
 */
package main

import (
	"bytes"
	"context"
	"crypto"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/unidoc/globalsign-dss"
	"github.com/unidoc/pkcs7"
	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"

	"golang.org/x/crypto/ocsp"
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
	sigLen = 8192

	apiKey, apiSecret, certFilepath, keyFilepath string
)

const usagef = "Usage: %s INPUT_PDF_PATH OUTPUT_PDF_PATH API_KEY API_SECRET CERT_FILE_PATH KEY_FILE_PATH\n"

func main() {
	args := os.Args
	if len(args) < 7 {
		fmt.Printf(usagef, os.Args[0])
		return
	}
	inputPath := args[1]
	outputPath := args[2]

	apiKey = args[3]
	apiSecret = args[4]
	certFilepath = args[5]
	keyFilepath = args[6]

	// We will simulate this by signing the PDF using GlobalSign DSS Service
	// and returning the signatured PDF data.
	pdfData, _, err := getExternalSignatureAndSign(inputPath)
	if err != nil {
		log.Fatalf("Fail signature: %v\n", err)
	}

	// Write output file.
	if err := ioutil.WriteFile(outputPath, pdfData, os.ModePerm); err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	log.Printf("PDF file successfully signed. Output path: %s\n", outputPath)
}

// generateSignedFile generates a signed version of the input PDF file using the
// specified signature handler.
func generateSignedFile(inputPath string, handler model.SignatureHandler, field *model.PdfFieldSignature) ([]byte, error) {
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

	if err = appender.Sign(1, field); err != nil {
		return nil, err
	}

	// Add LTV into signature.
	var certChain []*x509.Certificate
	if getter, ok := handler.(*externalSigner); ok {
		certChain = getter.getCertificateChain()
	}

	ltv, err := model.NewLTV(appender)
	if err != nil {
		return nil, err
	}
	ltv.CertClient.HTTPClient.Timeout = 30 * time.Second
	ltv.OCSPClient.HTTPClient.Timeout = 30 * time.Second
	ltv.CRLClient.HTTPClient.Timeout = 30 * time.Second

	err = ltv.EnableChain(certChain)
	if err != nil {
		return nil, err
	}

	// Write PDF file to buffer.
	pdfBuf := bytes.NewBuffer(nil)
	if err = appender.Write(pdfBuf); err != nil {
		return nil, err
	}

	return pdfBuf.Bytes(), nil
}

// getExternalSignatureAndSign get external signature, signs the specified
// PDF file and returns its signature.
func getExternalSignatureAndSign(inputPath string) ([]byte, *model.PdfSignature, error) {
	// Create signature handler.
	handler, err := NewGlobalSignPdfSignature()
	if err != nil {
		return nil, nil, err
	}

	option := &SignOption{
		Fullname: "GLOBALSIGN TEST ACCOUNT - FOR TESTING PURPOSE ONLY",
		Reason:   "Testing GlobalSign DSS",
		Location: "GLOBALSIGN TEST ACCOUNT - FOR TESTING PURPOSE ONLY",
		Annotate: true,
	}
	signatureField, signature, err := createSignatureField(option, handler)
	if err != nil {
		return nil, nil, err
	}

	pdfData, err := generateSignedFile(inputPath, handler, signatureField)
	if err != nil {
		return nil, nil, err
	}

	return pdfData, signature, nil
}

// externalSigner is wrapper for third-party signer,
// needs to implement function of model.SignatureHandler.
// in this example, we use GlobalSign as third-party signer.
type externalSigner struct {
	// ocsp retrieved during identity request.
	ocsp      []byte
	certChain []*x509.Certificate

	// client is third-party api client.
	// We make client as interface for supporting another third-party signer.
	client interface{}

	ctx context.Context
}

// NewGlobalSignPdfSignature wrap externalSigner to model.SignatureHandler.
// With this, you can implement any third-party client signer into signature handler.
func NewGlobalSignPdfSignature() (model.SignatureHandler, error) {
	ctx := context.Background()
	// Initiate globalsign client.
	c, err := globalsign.NewClient(apiKey, apiSecret, certFilepath, keyFilepath)
	if err != nil {
		return nil, err
	}

	return &externalSigner{
		client: c,
		ctx:    ctx,
	}, nil
}

// InitSignature sets the PdfSignature parameters.
func (es *externalSigner) InitSignature(sig *model.PdfSignature) error {
	gsClient, ok := es.client.(*globalsign.Client)
	if !ok {
		return fmt.Errorf("Not GlobalSign client.")
	}

	// Request new identification based on signer.
	identity, err := globalsign.DSSService.DSSGetIdentity(gsClient.DSSService, es.ctx, "GLOBALSIGN TEST ACCOUNT - FOR TESTING PURPOSE ONLY", &globalsign.IdentityRequest{
		SubjectDn: globalsign.SubjectDn{},
	})
	if err != nil {
		return err
	}

	// OCSP Response in base64 format.
	ocsp, err := base64.StdEncoding.DecodeString(identity.OCSP)
	if err != nil {
		return fmt.Errorf("invalid ocsp response, err: %v", err)
	}
	es.ocsp = ocsp

	// Create certificate chain from signing and CA cert.
	var certChain []*x509.Certificate
	issuerCertData := []byte(identity.SigningCert)
	for len(issuerCertData) != 0 {
		var block *pem.Block
		block, issuerCertData = pem.Decode(issuerCertData)
		if block == nil {
			break
		}

		issuer, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return err
		}

		certChain = append(certChain, issuer)
	}

	caCertData := []byte(identity.CA)
	for len(caCertData) != 0 {
		var block *pem.Block
		block, caCertData = pem.Decode(caCertData)
		if block == nil {
			break
		}

		issuer, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return err
		}

		certChain = append(certChain, issuer)
	}
	es.certChain = certChain

	// Create PDF array object which will contain the certificate chain data.
	pdfCerts := core.MakeArray()
	for _, cert := range certChain {
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

// Sign return error on failed sign.
func (es *externalSigner) Sign(sig *model.PdfSignature, digest model.Hasher) error {
	if digest == nil {
		sig.Contents = core.MakeHexString(string(make([]byte, sigLen)))
		return nil
	}

	gsClient, ok := es.client.(*globalsign.Client)
	if !ok {
		return fmt.Errorf("Not GlobalSign client.")
	}

	buffer := digest.(*bytes.Buffer)
	signedData, err := pkcs7.NewSignedData(buffer.Bytes())
	if err != nil {
		return err
	}

	// Set digest algorithm which supported by globalsign.
	signedData.SetDigestAlgorithm(pkcs7.OIDDigestAlgorithmSHA256)

	// Get certificate chain.
	certs := es.certChain

	// Callback.
	cb := func(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
		// Sign digest.
		signature, err := gsClient.DSSService.DSSIdentitySign(es.ctx, "GLOBALSIGN TEST ACCOUNT - FOR TESTING PURPOSE ONLY", &globalsign.IdentityRequest{SubjectDn: globalsign.SubjectDn{}}, digest)
		if err != nil {
			return nil, err
		}
		return signature, nil
	}

	siConfig := pkcs7.SignerInfoConfig{}
	if len(certs) > 1 {
		// Verify OCSP response.
		_, err := ocsp.ParseResponseForCert(es.ocsp, certs[0], certs[1])
		if err != nil {
			return err
		}

		siConfig.ExtraSignedAttributes = []pkcs7.Attribute{
			{
				Type: pkcs7.OIDAttributeAdobeRevocation,
				Value: RevocationInfoArchival{
					Crl: []asn1.RawValue{},
					Ocsp: []asn1.RawValue{
						{FullBytes: es.ocsp},
					},
					OtherRevInfo: []asn1.RawValue{},
				},
			},
		}
	}

	// If contains certificate chains.
	if len(certs) > 1 {
		err = signedData.AddSignerChain(certs[0], NewSigner(cb), certs[1:], siConfig)
		if err != nil {
			return err
		}
	} else if len(certs) == 1 {
		err = signedData.AddSigner(certs[0], NewSigner(cb), siConfig)
		if err != nil {
			return err
		}
	}

	// Add timestamp token to first signer.
	err = signedData.RequestSignerTimestampToken(0, func(digest []byte) ([]byte, error) {
		hasher := sha256.New()
		hasher.Write(digest)
		hasher.Sum(nil)

		// Request timestamp token.
		t, err := gsClient.DSSService.DSSIdentityTimestamp(es.ctx, "UniDoc", &globalsign.IdentityRequest{SubjectDn: globalsign.SubjectDn{}}, hasher.Sum(nil))
		if err != nil {
			return nil, err
		}

		return t, nil
	})
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

// getCertificateChain returns certificate chain.
func (es *externalSigner) getCertificateChain() []*x509.Certificate {
	return es.certChain
}

// RevocationInfoArchival is OIDAttributeAdobeRevocation attribute.
type RevocationInfoArchival struct {
	Crl          []asn1.RawValue `asn1:"explicit,tag:0,optional"`
	Ocsp         []asn1.RawValue `asn1:"explicit,tag:1,optional"`
	OtherRevInfo []asn1.RawValue `asn1:"explicit,tag:2,optional"`
}

// SignerCallback callback function of `Signer`.
type SignerCallback func(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error)

// Signer implements custom crypto.Signer which utilize globalsign DSS API
// to sign signature digest.
type Signer struct {
	// sign callback
	callback SignerCallback
}

// EncryptionAlgorithmOID returns asn1.ObjectIdentifier with value of pcks7.OIDEncryptionAlgorithmRSA.
func (s *Signer) EncryptionAlgorithmOID() asn1.ObjectIdentifier {
	return pkcs7.OIDEncryptionAlgorithmRSA
}

// Public returns PublicKey.
func (s *Signer) Public() crypto.PublicKey {
	return nil
}

// Sign with custom crypto.Signer which utilize globalsign DSS API.
func (s *Signer) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	if s.callback == nil {
		return nil, errors.New("signer func not implemented")
	}
	return s.callback(rand, digest, opts)
}

// NewSigner returns a new crypto.Signer implementation.
func NewSigner(cb SignerCallback) crypto.Signer {
	return &Signer{
		callback: cb,
	}
}

// SignOption contains both digital signing
// and annotation properties.
type SignOption struct {
	SignedBy string
	Fullname string
	Reason   string
	Location string

	// Annonate signature.
	Annotate bool

	// Position of annotation.
	Position []float64

	// Annotation font size.
	FontSize int

	// Extra signature annotation fields.
	Extra map[string]string

	FilePath string

	// Just in case source file is protected
	// and default password is not empty.
	Password string
}

// createSignatureField creates a signature field and an associated signature,
// based on the specified options.
func createSignatureField(option *SignOption, handler model.SignatureHandler, certChain ...*x509.Certificate) (*model.PdfFieldSignature, *model.PdfSignature, error) {
	// Create new signature.
	signature := model.NewPdfSignature(handler)

	if len(certChain) > 0 {
		// Create PDF array object which will contain the certificate chain data,
		// The first element of the array must be the signing certificate.
		// The rest of the certificate chain is used for validating the authenticity
		// of the signing certificate.
		pdfCerts := core.MakeArray()
		for _, cert := range certChain {
			pdfCerts.Append(core.MakeString(string(cert.Raw)))
		}

		signature.Cert = pdfCerts
	}

	if err := signature.Initialize(); err != nil {
		return nil, nil, err
	}

	now := time.Now()
	signature.SetName(option.Fullname)
	signature.SetReason(option.Reason)
	signature.SetDate(now, "D:20060102150405-07'00'")
	signature.SetLocation(option.Location)

	// Create signature field and appearance.
	signatureFields := make([]*annotator.SignatureLine, 0)
	opts := annotator.NewSignatureFieldOpts()

	// Only when annotate option is enabled.
	if option.Annotate {
		if option.FontSize > 0 {
			opts.FontSize = float64(option.FontSize)
		}

		// Set default position.
		opts.Rect = []float64{10, 25, 75, 60}
		if option.Position != nil && len(option.Position) == 4 {
			opts.Rect = option.Position
		}

		signatureFields = append(signatureFields,
			annotator.NewSignatureLine("Signed By", option.Fullname),
			annotator.NewSignatureLine("Date", now.Format(time.RFC1123)),
			annotator.NewSignatureLine("Reason", option.Reason),
			annotator.NewSignatureLine("Location", option.Location),
		)

		for k, v := range option.Extra {
			signatureFields = append(signatureFields, annotator.NewSignatureLine(k, v))
		}
	}

	field, err := annotator.NewSignatureField(
		signature,
		signatureFields,
		opts,
	)
	if err != nil {
		return nil, nil, err
	}

	field.T = core.MakeString("External Signature")

	return field, signature, nil
}
