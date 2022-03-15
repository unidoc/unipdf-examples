package sign_handler

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

	"github.com/unidoc/globalsign-dss"
	"github.com/unidoc/pkcs7"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
	"golang.org/x/crypto/ocsp"
)

const sigLen = 8192

// SignerCallback .
type SignerCallback func(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error)

// Signer implements custom crypto.Signer which utilize globalsign DSS API
// to sign signature digest
type Signer struct {
	// sign callback
	callback SignerCallback
}

func (s *Signer) EncryptionAlgorithmOID() asn1.ObjectIdentifier {
	return pkcs7.OIDEncryptionAlgorithmRSA
}

// Public .
func (s *Signer) Public() crypto.PublicKey {
	return nil
}

// Sign request
func (s *Signer) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	if s.callback == nil {
		return nil, errors.New("signer func not implemented")
	}

	return s.callback(rand, digest, opts)
}

// NewSigner create crypto.Signer implementation
func NewSigner(cb SignerCallback) crypto.Signer {
	return &Signer{
		callback: cb,
	}
}

type CertificateChainGetter interface {
	GetCertificateChain() []*x509.Certificate
}

// GlobalsignDSS is custom unidoc sighandler which leverage
// globalsign DSS service
type GlobalsignDSS struct {
	signer   string
	identity map[string]interface{}

	// ocsp retrieved during identity request
	ocsp      []byte
	certChain []*x509.Certificate

	client *globalsign.Client

	// a caller context
	ctx context.Context
}

func (h *GlobalsignDSS) GetCertificateChain() []*x509.Certificate {
	return h.certChain
}

func (h *GlobalsignDSS) getCertificates(sig *model.PdfSignature) ([]*x509.Certificate, error) {
	var certData []byte
	switch certObj := sig.Cert.(type) {
	case *core.PdfObjectString:
		certData = certObj.Bytes()
	case *core.PdfObjectArray:
		if certObj.Len() == 0 {
			return nil, errors.New("no signature certificates found")
		}
		for _, obj := range certObj.Elements() {
			certStr, ok := core.GetString(obj)
			if !ok {
				return nil, fmt.Errorf("invalid certificate object type in signature certificate chain: %T", obj)
			}
			certData = append(certData, certStr.Bytes()...)
		}
	default:
		return nil, fmt.Errorf("invalid signature certificate object type: %T", certObj)
	}

	certs, err := x509.ParseCertificates(certData)
	if err != nil {
		return nil, err
	}

	return certs, nil
}

// IsApplicable .
func (h *GlobalsignDSS) IsApplicable(sig *model.PdfSignature) bool {
	if sig == nil || sig.Filter == nil || sig.SubFilter == nil {
		return false
	}
	return (*sig.Filter == "Adobe.PPKMS" || *sig.Filter == "Adobe.PPKLite") && *sig.SubFilter == "adbe.pkcs7.detached"
}

// Validate .
func (h *GlobalsignDSS) Validate(sig *model.PdfSignature, digest model.Hasher) (model.SignatureValidationResult, error) {
	return model.SignatureValidationResult{
		IsSigned:   true,
		IsVerified: true,
	}, nil
}

// InitSignature sets the PdfSignature parameters.
func (h *GlobalsignDSS) InitSignature(sig *model.PdfSignature) error {
	// request new identification based on signer
	identity, err := h.client.DSSService.DSSGetIdentity(h.ctx, h.signer, &globalsign.IdentityRequest{
		SubjectDn: globalsign.SubjectDn{},
	})
	if err != nil {
		return err
	}

	// OCSP Response in base64 format
	h.ocsp, err = base64.StdEncoding.DecodeString(identity.OCSP)
	if err != nil {
		return fmt.Errorf("invalid ocsp response, err: %v", err)
	}

	// create certificate chain
	// from signing and ca cert
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
	h.certChain = certChain

	// Create PDF array object which will contain the certificate chain data
	pdfCerts := core.MakeArray()
	for _, cert := range certChain {
		pdfCerts.Append(core.MakeString(string(cert.Raw)))
	}

	// append cert to signature
	sig.Cert = pdfCerts

	handler := *h
	sig.Handler = &handler
	sig.Filter = core.MakeName("Adobe.PPKLite")
	sig.SubFilter = core.MakeName("adbe.pkcs7.detached")
	sig.Reference = nil

	// reserve initial size
	return handler.Sign(sig, nil)
}

// NewDigest .
func (h *GlobalsignDSS) NewDigest(sig *model.PdfSignature) (model.Hasher, error) {
	return bytes.NewBuffer(nil), nil
}

// Sign .
func (h *GlobalsignDSS) Sign(sig *model.PdfSignature, digest model.Hasher) error {
	if digest == nil {
		sig.Contents = core.MakeHexString(string(make([]byte, sigLen)))
		return nil
	}

	buffer := digest.(*bytes.Buffer)
	signedData, err := pkcs7.NewSignedData(buffer.Bytes())
	if err != nil {
		return err
	}

	// set digest algorithm which supported by globalsign
	signedData.SetDigestAlgorithm(pkcs7.OIDDigestAlgorithmSHA256)

	// get certificate chain
	certs := h.GetCertificateChain()

	// callback
	cb := func(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
		// sign digest
		signature, err := h.client.DSSService.DSSIdentitySign(h.ctx, h.signer, &globalsign.IdentityRequest{SubjectDn: globalsign.SubjectDn{}}, digest)
		if err != nil {
			return nil, err
		}

		return signature, nil
	}

	siConfig := pkcs7.SignerInfoConfig{}
	if len(h.ocsp) != 0 {

		// verify ocsp response
		_, err := ocsp.ParseResponseForCert(h.ocsp, certs[0], certs[1])
		if err != nil {
			return err
		}

		siConfig.ExtraSignedAttributes = []pkcs7.Attribute{
			{
				Type: pkcs7.OIDAttributeAdobeRevocation,
				Value: pkcs7.RevocationInfoArchival{
					Crl: []asn1.RawValue{},
					Ocsp: []asn1.RawValue{
						{FullBytes: h.ocsp},
					},
					OtherRevInfo: []asn1.RawValue{},
				},
			},
		}
	}

	// if contains certificate chains
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

	// add timestamp token to first signer
	err = signedData.RequestSignerTimestampToken(0, func(digest []byte) ([]byte, error) {
		hasher := sha256.New()
		hasher.Write(digest)
		hasher.Sum(nil)

		// request timestamp token
		t, err := h.client.DSSService.DSSIdentityTimestamp(h.ctx, h.signer, &globalsign.IdentityRequest{SubjectDn: globalsign.SubjectDn{}}, hasher.Sum(nil))
		if err != nil {
			return nil, err
		}

		return t, nil
	})
	if err != nil {
		return err
	}

	// Call Detach() is you want to remove content from the signature
	// and generate an S/MIME detached signature
	signedData.Detach()
	// Finish() to obtain the signature bytes
	detachedSignature, err := signedData.Finish()
	if err != nil {
		return err
	}

	data := make([]byte, sigLen)
	copy(data, detachedSignature)

	sig.Contents = core.MakeHexString(string(data))
	return nil
}

// NewGlobalSignDSS create custom unidoc sighandler which leverage globalsign DSS service
// this handler assume that globalsign credential already have been set
// please see globalsign.SetupCredential
func NewGlobalSignDSS(ctx context.Context, c *globalsign.Client, signer string, identity map[string]interface{}) (model.SignatureHandler, error) {
	return &GlobalsignDSS{
		ctx:      ctx,
		client:   c,
		signer:   signer,
		identity: identity,
	}, nil
}

// InitFunc allow customize pdf signature initialization
type InitFunc func(model.SignatureHandler, *model.PdfSignature) error

// SignFunc allow customize signing implementation
type SignFunc func(model.SignatureHandler, *model.PdfSignature, model.Hasher) error

type CustomHandler struct {
	initFunc InitFunc
	signFunc SignFunc
}

// NewDigest .
func (h *CustomHandler) NewDigest(sig *model.PdfSignature) (model.Hasher, error) {
	return bytes.NewBuffer(nil), nil
}

// IsApplicable .
func (h *CustomHandler) IsApplicable(sig *model.PdfSignature) bool {
	if sig == nil || sig.Filter == nil || sig.SubFilter == nil {
		return false
	}

	return (*sig.Filter == "Adobe.PPKMS" || *sig.Filter == "Adobe.PPKLite") && *sig.SubFilter == "adbe.pkcs7.detached"
}

// InitSignature .
func (h *CustomHandler) InitSignature(sig *model.PdfSignature) error {
	if h.initFunc != nil {
		return h.initFunc(h, sig)
	}

	return nil
}

// Sign .
func (h *CustomHandler) Sign(sig *model.PdfSignature, digest model.Hasher) error {
	if h.signFunc != nil {
		return h.signFunc(h, sig, digest)
	}

	return nil
}

// Validate .
func (h *CustomHandler) Validate(sig *model.PdfSignature, digest model.Hasher) (model.SignatureValidationResult, error) {

	return model.SignatureValidationResult{
		IsSigned:   true,
		IsVerified: true,
	}, nil
}

// NewCustomHandler allow client to implement their own init and sign function
func NewCustomHandler(initFunc InitFunc, signFunc SignFunc) model.SignatureHandler {
	return &CustomHandler{
		initFunc: initFunc,
		signFunc: signFunc,
	}
}
