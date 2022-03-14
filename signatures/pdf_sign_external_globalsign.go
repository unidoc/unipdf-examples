/*
 * This example showcases how to sign a PDF document with the GlobalSign.
 * The file is signed using a private/public key pair.
 *
 * $ ./pdf_sign_external_globalsign <INPUT_PDF_PATH> <OUTPUT_PDF_PATH>
 */
package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/unidoc/globalsign-dss"
	"github.com/unidoc/unidoc-examples/signatures/sign_handler"
	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}

	// Set logger.
	common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))

	// Set signer factory.
	signerFactories["globalsign"] = NewGlobalSignDssSigner
}

var (
	inputFile  string
	outputFile string
	email      string
	fullname   string
	reason     string

	certPath string
	keyPath  string

	apiKey     = ""
	apiSecret  = ""
	apiBaseURL = "https://emea.api.dss.globalsign.com:8443"
)

func main() {
	flag.StringVar(&inputFile, "input-file", "", "file to be signed (required)")
	flag.StringVar(&outputFile, "output-file", "", "output result (required)")
	flag.StringVar(&email, "email", "", "email for signer identity (required)")
	flag.StringVar(&apiKey, "api-key", "", "API key (required)")
	flag.StringVar(&apiSecret, "api-secret", "", "API secret (required)")
	flag.StringVar(&certPath, "cert-file", "tls.cer", "certificate file for API (required)")
	flag.StringVar(&keyPath, "key-file", "key.pem", "key file for API (required)")
	flag.StringVar(&fullname, "name", "your n@me", "signer name")
	flag.StringVar(&reason, "reason", "enter your re@son", "signing reason")

	flag.Parse()

	if inputFile == "" || outputFile == "" || email == "" || apiKey == "" || apiSecret == "" || certPath == "" || keyPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	option := &SignOption{
		SignedBy: "UniDoc",
		Fullname: "Alip Sulistio",
		Reason:   "GlobalSign DSS Testing",
		Annotate: true,
	}

	sigGen := NewGlobalSignDssSigner(map[string]interface{}{
		"provider.globalsign.api_url":     apiBaseURL,
		"provider.globalsign.api_key":     apiKey,
		"provider.globalsign.api_secret":  apiSecret,
		"provider.globalsign.certificate": certPath,
		"provider.globalsign.private_key": keyPath,
	})

	if err := SignFile(context.Background(), inputFile, outputFile, option, sigGen); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("File signed successfully")
}

type globalsignDssSigner struct {
	apiBase      string
	apiKey       string
	apiSecret    string
	certFilepath string
	keyFilepath  string

	client *globalsign.Client
}

// Load .
func (s *globalsignDssSigner) Load() error {
	// Initiate globalsign client.
	c, err := globalsign.NewClient(s.apiKey, s.apiSecret, s.certFilepath, s.keyFilepath)
	if err != nil {
		return err
	}
	s.client = c

	return nil
}

// Sign .
func (s *globalsignDssSigner) Sign(ctx context.Context, rd *model.PdfReader, option *SignOption) (*model.PdfAppender, error) {
	// ensure pdf is decrypted
	isEncrypted, err := rd.IsEncrypted()
	if err != nil {
		return nil, err
	}

	if isEncrypted {
		log.Println("pdf is encrypted")
		auth, err := rd.Decrypt([]byte(option.Password))
		if err != nil {
			return nil, err
		}
		if !auth {
			return nil, errors.New("cannot open encrypted document, please specify password in option")
		}
	}

	isEncrypted, err = rd.IsEncrypted()
	if err != nil {
		return nil, err
	}

	log.Println("pdf is encrypted?", isEncrypted)

	ap, err := model.NewPdfAppender(rd)
	if err != nil {
		return nil, err
	}

	signerIdentity := map[string]interface{}{
		"common_name": option.Fullname,
	}

	// Create signature handler.
	handler, err := sign_handler.NewGlobalSignDSS(context.Background(), s.client, option.SignedBy, signerIdentity)
	if err != nil {
		return nil, err
	}

	field, err := createSignatureField(option, handler)
	if err != nil {
		return nil, err
	}

	// get cert chain
	var certChain []*x509.Certificate
	if getter, ok := handler.(sign_handler.CertificateChainGetter); ok {
		certChain = getter.GetCertificateChain()
	}

	// only sign first page
	if err = ap.Sign(1, field); err != nil {
		return nil, err
	}

	// add tlv
	ltv, err := model.NewLTV(ap)
	if err != nil {
		return nil, err
	}
	ltv.CertClient.HTTPClient.Timeout = 30 * time.Second
	ltv.OCSPClient.HTTPClient.Timeout = 30 * time.Second
	ltv.CRLClient.HTTPClient.Timeout = 1 * time.Microsecond // attempt to exclude crl

	err = ltv.EnableChain(certChain)
	if err != nil {
		return nil, err
	}

	return ap, nil
}

// NewGlobalSignDssSigner create and return instance signature
// generator backed by global sign
func NewGlobalSignDssSigner(param map[string]interface{}) Signer {
	mReader := NewMapReader(param)
	apiURL := mReader.String("provider.globalsign.api_url", "")
	apiKey := mReader.String("provider.globalsign.api_key", "")
	apiSecret := mReader.String("provider.globalsign.api_secret", "")
	apiCertFile := mReader.String("provider.globalsign.certificate", "")
	keyFile := mReader.String("provider.globalsign.private_key", "")

	return &globalsignDssSigner{
		apiBase:      apiURL,
		apiKey:       apiKey,
		apiSecret:    apiSecret,
		certFilepath: apiCertFile,
		keyFilepath:  keyFile,
	}
}

// SignOption contains both digital signing
// and annotation properties
type SignOption struct {
	SignedBy string
	Fullname string
	Reason   string
	Location string

	// Annonate signature?
	Annotate bool

	// position of annotation
	Position []float64

	// Annotation font size
	FontSize int

	// extra signature annotation fields
	Extra map[string]string

	FilePath string

	// just in case source file is protected
	// and defalt password is not empty
	Password string
}

func defaultSignOption() *SignOption {
	return &SignOption{
		FontSize: 11,
	}
}

// generateSignedFile generates a signed version of the input PDF file using the
// specified signature handler.
func generateSignedFile(inputPath string, handler model.SignatureHandler, option *SignOption) ([]byte, *model.PdfSignature, error) {
	if option == nil {
		option = defaultSignOption()
	}

	// generate timestamp
	now := time.Now()

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

	// Create pdf appender.
	appender, err := model.NewPdfAppender(reader)
	if err != nil {
		return nil, nil, err
	}

	// Create signature.
	signature := model.NewPdfSignature(handler)
	signature.SetName(option.SignedBy)
	signature.SetReason(option.Reason)
	signature.SetDate(now, time.RFC1123)
	signature.SetLocation(option.Location)

	if err := signature.Initialize(); err != nil {
		return nil, nil, err
	}

	// Create signature field and appearance.
	var field *model.PdfFieldSignature

	// onyl when annotate option is enabled
	if option.Annotate {
		opts := annotator.NewSignatureFieldOpts()
		opts.FontSize = 10

		// set default position
		opts.Rect = []float64{10, 25, 75, 60}
		if option.Position != nil && len(option.Position) == 4 {
			opts.Rect = option.Position
		}

		signatureFields := []*annotator.SignatureLine{
			annotator.NewSignatureLine("Signed By", option.SignedBy),
			annotator.NewSignatureLine("Date", now.Format(time.RFC1123)),
			annotator.NewSignatureLine("Reason", option.Reason),
			annotator.NewSignatureLine("Location", option.Location),
		}

		for k, v := range option.Extra {
			signatureFields = append(signatureFields, annotator.NewSignatureLine(k, v))
		}

		field, err = annotator.NewSignatureField(
			signature,
			signatureFields,
			opts,
		)
		field.T = core.MakeString("Signature")
	}

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

// SignerFactory .
type SignerFactory func(map[string]interface{}) Signer

var signerFactories = make(map[string]SignerFactory)

// CreateSigner .
func CreateSigner(signerType string, param map[string]interface{}) Signer {
	factory, ok := signerFactories[signerType]
	if !ok {
		return nil
	}

	return factory(param)
}

// Sign apply digital signing from inputFile to outputFile
// with signature generator callback
func SignFile(ctx context.Context, inputFile, outputFile string, option *SignOption, signer Signer) error {

	fin, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer fin.Close()

	rd, err := model.NewPdfReader(fin)
	if err != nil {
		return err
	}

	ap, err := Sign(ctx, rd, option, signer)
	if err != nil {
		return err
	}

	fout, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer fout.Close()

	return ap.Write(fout)
}

// Sign apply digital signing to given pdf reader which return pdf appender
func Sign(ctx context.Context, rd *model.PdfReader, option *SignOption, signer Signer) (*model.PdfAppender, error) {
	if signer == nil {
		return nil, errors.New("signer not provided")
	}

	if err := signer.Load(); err != nil {
		return nil, err
	}

	if option == nil {
		option = defaultSignOption()
	}

	return signer.Sign(ctx, rd, option)
}

// Signer abstract pdf signer implementation
type Signer interface {
	// Load init and prepare signer
	// it may fail on bad configuration
	Load() error

	// Sign .
	Sign(context.Context, *model.PdfReader, *SignOption) (*model.PdfAppender, error)
}

// UpdateInfo set tool author info for created pdf
func UpdateInfo(author, creator string) {
	model.SetPdfAuthor(author)
	model.SetPdfCreator(creator)
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

// GenerateChecksum returns checksum of a reader,
// reader seek head will be returned to 0 (beginning of file)
func GenerateChecksum(r io.ReadSeeker) []byte {
	bufferedReader := bufio.NewReader(r)
	computedChecksum := sha256.New()
	_, err := bufferedReader.WriteTo(computedChecksum)
	defer r.Seek(0, io.SeekStart)
	if err != nil {
		return make([]byte, 0)
	}

	return computedChecksum.Sum(nil)
}

func loadPrivateKey(privateKeyData string) (*rsa.PrivateKey, error) {
	// Decode PEM block.
	block, _ := pem.Decode([]byte(privateKeyData))

	// Parse private key data.
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func loadCertificates(certData string) (*x509.Certificate, *core.PdfObjectArray, error) {
	parseCert := func(data []byte) (*x509.Certificate, []byte, error) {
		// Decode PEM block.
		block, rest := pem.Decode(data)

		// Parse certificate.
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, nil, err
		}

		return cert, rest, nil
	}

	// Create PDF array object which will contain the certificate chain data,
	// loaded from the PEM file. The first element of the array must be the
	// signing certificate. The rest of the certificate chain is used for
	// validating the authenticity of the signing certificate.
	pdfCerts := core.MakeArray()

	// Parse signing certificate.
	signingCert, pemUnparsedData, err := parseCert([]byte(certData))
	if err != nil {
		return nil, nil, err
	}
	pdfCerts.Append(core.MakeString(string(signingCert.Raw)))

	// Parse the rest of the certificates contained in the PEM file,
	// if any, and add them to the PDF certificates array.
	for len(pemUnparsedData) != 0 {
		cert, rest, err := parseCert(pemUnparsedData)
		if err != nil {
			return nil, nil, err
		}

		pdfCerts.Append(core.MakeString(string(cert.Raw)))
		pemUnparsedData = rest
	}

	return signingCert, pdfCerts, nil
}

func createSignatureField(option *SignOption, handler model.SignatureHandler, certChain ...*x509.Certificate) (*model.PdfFieldSignature, error) {
	// Create signature.
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
		return nil, err
	}

	now := time.Now()
	signature.SetName(option.Fullname)
	signature.SetReason(option.Reason)
	signature.SetDate(now, "D:20060102150405-07'00'")
	signature.SetLocation(option.Location)

	// Create signature field and appearance.
	signatureFields := make([]*annotator.SignatureLine, 0)
	opts := annotator.NewSignatureFieldOpts()

	// onyl when annotate option is enabled
	if option.Annotate {
		if option.FontSize > 0 {
			opts.FontSize = float64(option.FontSize)
		}

		// set default position
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
		return nil, err
	}

	field.T = core.MakeString("Signature")

	return field, nil
}

// MapReader .
type MapReader struct {
	m map[string]interface{}
}

func (m *MapReader) String(key string, def ...string) string {
	if len(def) == 0 {
		def = []string{""}
	}

	v, ok := m.m[key]
	if !ok {
		return def[0]
	}

	vv, ok := v.(string)
	if !ok {
		return def[0]
	}

	return vv
}

// NewMapReader .
func NewMapReader(m map[string]interface{}) *MapReader {
	return &MapReader{m: flatten(m)}
}

func flatten(value interface{}) map[string]interface{} {
	return flattenPrefixed(value, "")
}

func flattenPrefixed(value interface{}, prefix string) map[string]interface{} {
	m := make(map[string]interface{})
	flattenPrefixedToResult(value, prefix, m)
	return m
}

func flattenPrefixedToResult(value interface{}, prefix string, m map[string]interface{}) {
	base := ""
	if prefix != "" {
		base = prefix + "."
	}

	cm, ok := value.(map[string]interface{})
	if ok {
		for k, v := range cm {
			flattenPrefixedToResult(v, base+k, m)
		}
	} else {
		if prefix != "" {
			m[prefix] = value
		}
	}
}

// CopyFile the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
