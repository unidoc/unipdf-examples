# LTV

## Table of contents

- [Overview](#overview)
- [LTV enable workflows](#ltv-enable-workflows)
  - [Sign and LTV enable in the same revision](#1-sign-and-ltv-enable-in-the-same-revision)
  - [Sign and LTV enable in a separate revision](#2-sign-and-ltv-enable-in-a-separate-revision)
- [Usage](#usage)
  - [Running the examples](#running-the-examples)
  - [LTV client](#ltv-client)
  - [Sign and LTV enable in one revision](#workflow-1-sign-and-ltv-enable-in-one-revision)
  - [LTV enable signed file](#workflow-2-ltv-enable-signed-file)
  - [Protect validation data by adding a timestamp signature](#protect-validation-data-by-adding-a-timestamp-signature)
  - [Customize LTV client](#customize-ltv-client)

## Overview

LTV (Long-Term Validation) can be enabled by adding validation information
for the digital signatures applied to a PDF file.

The validation data required to LTV enable a signature consists of:
- The signing certificate chain, which contains the certificate used to apply
  the signature and all its issuers, up to a trusted root (CA) certificate.
- OCSP (Online Certificate Status Protocol) responses for each certificate in
  the added certificate chain.
- CRL (Certificate Revocation List) responses for each certificate in the
  added certificate chain.

`NOTE`: a signature can be LTV enabled by including only OCSP or only CRL
responses for the signing certificate chain. However, the more validation
data is included, the better.

OCSP and CRL responses provide revocation information for the certificates used
in the signing process. They are obtained by querying OCSP/CRL servers.
The responders (servers) used to obtain the OCSP responses are extracted from
the certificates in the signing chain
([x509.Certificate.OCSPServer](https://godoc.org/crypto/x509#Certificate)).
Similarly, the servers used to obtain CRL responses are also extracted from
each certificate in the signing chain
([x509.Certificate.CRLDistributionPoints](https://godoc.org/crypto/x509#Certificate)).

LTV enabling is done through the use of a DSS (Document Security Store) dictionary,
which contains the validation information described above.

```
DSS
    // Global validation data.
    // The validation data in these fields can be used to validate all
    // signatures in a PDF document.
    Certs [PDF stream array]
    OCSPs [PDF stream array]
    CRLs  [PDF stream array]

    // Signature specific validation data.
    // Each key in the VRI dictionary represents the hash of the Contents
    // field of a signature dictionary. Each signature entry in the VRI
    // dictionary has its own set of certificates, OCSP and CRL responses,
    // used to validate that particular signature exclusively.
    VRI:
        Hash signature 1: Cert [PDF stream array]
                          OCSP [PDF stream array]
                          CRL  [PDF stream array]
        ...
        Hash signature N: Cert [PDF stream array]
                          OCSP [PDF stream array]
                          CRL  [PDF stream array]
```

The Contents field of a signature is only known after signing and writing a
PDF file. In consequence, VRI signature entries cannot be added in the same
revision a signature is applied. However, global validation data can be added
in the signing revision, provided that at least the signing certificate is
known at signing time, which is not always the case (e.g.: when using external
services to apply the signature).

`NOTE`: Adobe's definition of LTV is not clearly defined. Also, they
have not disclosed their validation process process which results in the
`Signature is LTV enabled` label being displayed for a signature. However,
the general consensus is that both a signing certificate chain (which builds
up to a trusted root certificate) and revocation data for it (OCSP and CRL
responses) need to be included.

## LTV enable workflows

#### 1. Sign and LTV enable in the same revision

Useful when signing and LTV enabling must be done in one revision. At least
the signing certificate must be provided. The validation data is added to the
global scope of the DSS and can be used to validate any of the signatures in
the file.
```
    Revision 1 (signer 1):
      - Add validation information (certificates, OCSP and CRL information) to
        the global validation data of the DSS dictionary.
      - Apply signature 1.
      - Write signed file.
    ...
    Revision N (signer N):
      - Add validation information (certificates, OCSP and CRL information) to
        the global validation data of the DSS dictionary.
      - Apply signature N.
      - Write signed file.
```

#### 2. Sign and LTV enable in a separate revision
Useful when the validation data is not available at signing time. Furthermore,
both global and signature specific data can be added for all signatures applied
previous revisions.
```
    // Signing revisions.
    Revision 1 (signer 1):
      - Apply signature 1.
      - Write signed file.
    ...
    Revision N (signer N):
      - Apply signature N.
      - Write signed file.

    // LTV enable revision.
    Revision N+1:
      - Add both global and signature specific validation data (certificates,
        OCSP and CRL information) to the DSS dictionary.
      - Write file.
```

The workflows can be repeated as many times as needed, by adding extra
revisions to the document, either for additional signatures or for adding
validation data for the previous revision/revisions.

`NOTE`: Hybrid workflows (consisting of combinations of the two workflows above)
are also possible.

## Usage

#### Running the examples

```
go run pdf_sign_ltv_one_revision.go P12_FILE PASSWORD INPUT_PDF_PATH OUTPUT_PDF_PATH [EXTRA_CERTS.pem]
go run pdf_ltv_enable_signed_file.go INPUT_PDF_PATH OUTPUT_PDF_PATH [EXTRA_CERTS.pem]
go run pdf_sign_ltv_timestamp_revision.go P12_FILE PASSWORD INPUT_PDF_PATH OUTPUT_PDF_PATH [EXTRA_CERTS.pem]
```

`EXTRA_CERTS.pem` holds issuer certificates to add to the chain, in PEM
format. The signing examples also add the CA certificates bundled in the
PKCS12 file. All three fetch missing issuers, OCSP responses and CRLs over
the network. `pdf_sign_ltv_timestamp_revision.go` also requests a timestamp
from the freetsa.org TSA.

To check the output, run
[../validation/pdf_sign_validate_trust.go](../validation/pdf_sign_validate_trust.go)
in offline mode (it reads the revocation data from the DSS) or
[../validation/pdf_sign_validate_at_time.go](../validation/pdf_sign_validate_at_time.go).

#### LTV client

unipdf provides a configurable client for LTV enabling signed PDF documents,
which supports both workflows described above.

```go
type LTV
    //
    // Fields.
    //

    CertClient *sigutil.CertClient // Used to retrieve certificates.
    OCSPClient *sigutil.OCSPClient // Used to retrieve OCSP information.
    CRLClient  *sigutil.CRLClient  // Used to retrieve CRL information.

    // Specifies whether existing signature validations should be skipped.
    SkipExisting bool

    //
    // Workflow #1 methods: adding LTV data in the signing revision.
    //

    // EnableChain adds the specified certificate chain and validation data
    // for it to the global scope of the document DSS.
    EnableChain(chain []*x509.Certificate) error

    //
    // Workflow #2 methods: adding LTV data in a separate revision.
    //

    // LTV enables all signatures in a PDF document.
    EnableAll(extraCerts []*x509.Certificate) error

    // Enable LTV enables the specified signature.
    Enable(sig *model.PdfSignature, extraCerts []*x509.Certificate) error

    // LTV enables the signature dictionary of the PDF AcroForm field
    // identified the specified name.
    EnableByName(name string, extraCerts []*x509.Certificate) error
```

If `LTV.SkipExisting` is true (the default), existing signature validations
won't be reconstructed when calling `EnableAll`, `Enable` or `EnableByName`,
thus making the LTV process incremental.

Example:
```
    Revision 1: apply signature 1.
    Revision 2: LTV enable by using LTV.EnableAll (will enable signature 1).
    Revision 3: apply signature 2.
    Revision 4: apply signature 3.
    Revision 5: LTV enable by using LTV.EnableAll (will enable signatures 2 and 3).
```

#### Workflow #1: sign and LTV enable in one revision

The sample showcases applying a digital signature to a PDF document, and
LTV enabling the signing certificate chain, in a single revision. The example
uses a `adbe.pkcs7.detached` signature handler. However, the handler can be
swapped with a `adbe.x509.rsa_sha1` handler (no other changes are required).

The signing certificate chain must be provided (which needs to contain at
least the signing certificate). The LTV client completes the chain by
downloading missing issuer certificates. It stops at the first self-signed CA
certificate or when an issuer cannot be fetched, and it does not check that
the chain ends at a root you trust.

For more information, see [pdf_sign_ltv_one_revision.go](pdf_sign_ltv_one_revision.go).

```go
// Get private key and X509 certificate from a PKCS12 (.p12/.pfx) file.
pfxData, err := os.ReadFile("p12.pfx")
if err != nil {
    log.Fatal(err)
}

key, cert, caCerts, err := pkcs12.DecodeChain(pfxData, password)
if err != nil {
    log.Fatal(err)
}
priv, ok := key.(*rsa.PrivateKey)
if !ok {
    log.Fatal("PKCS12 private key is not an RSA key")
}

file, err := os.Open("unsigned-file.pdf")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

reader, err := model.NewPdfReader(file)
if err != nil {
    log.Fatal(err)
}

// Create appender.
appender, err := model.NewPdfAppender(reader)
if err != nil {
    log.Fatal(err)
}

// Create signature handler.
handler, err := sighandler.NewAdobePKCS7Detached(priv, cert)
if err != nil {
    log.Fatal(err)
}

// Create signature.
signature := model.NewPdfSignature(handler)
signature.SetName("Test Sign LTV enable")

if err := signature.Initialize(); err != nil {
    log.Fatal(err)
}

// Create signature appearance.
opts := annotator.NewSignatureFieldOpts()
opts.Rect = []float64{10, 25, 75, 60}

field, err := annotator.NewSignatureField(
    signature,
    []*annotator.SignatureLine{
        annotator.NewSignatureLine("Name", "John Doe"),
        annotator.NewSignatureLine("Reason", "Signature test"),
    },
    opts,
)
if err != nil {
    log.Fatal(err)
}
field.T = core.MakeString("Test Sign LTV enable")

if err = appender.Sign(1, field); err != nil {
    log.Fatal(err)
}

// LTV enable the certificate chain used to apply the signature.
ltv, err := model.NewLTV(appender)
if err != nil {
    log.Fatal(err)
}
if err := ltv.EnableChain(append([]*x509.Certificate{cert}, caCerts...)); err != nil {
    log.Fatal(err)
}

// Write output file.
if err = appender.WriteToFile("output.pdf"); err != nil {
    log.Fatal(err)
}
```

#### Workflow #2: LTV enable signed file

The sample showcases LTV enabling a signed PDF file through an additional
revision. This workflow is useful when the signing certificate is not
available at the time of signing (e.g.: when signing using external services).

The certificate chains used for the signatures in the file are
extracted from the signature dictionaries. The LTV client completes each
chain by downloading missing issuer certificates, with the same limits as in
workflow #1: it stops at a self-signed CA certificate or on a fetch error and
does not check trust.

`LTV.EnableAll` returns an error when the form has a signature field that is
not signed yet. The example walks the signature fields and calls `LTV.Enable`
only on signed ones.

For more information, see [pdf_ltv_enable_signed_file.go](pdf_ltv_enable_signed_file.go).

```go
// Create reader.
file, err := os.Open("signed-file.pdf")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

reader, err := model.NewPdfReader(file)
if err != nil {
    log.Fatal(err)
}

// Create appender.
appender, err := model.NewPdfAppender(reader)
if err != nil {
    log.Fatal(err)
}

// LTV enable.
ltv, err := model.NewLTV(appender)
if err != nil {
    log.Fatal(err)
}

// Enable only the fields that hold a signature; EnableAll fails on an
// empty signature field.
for _, field := range reader.AcroForm.AllFields() {
    sigField, ok := field.GetContext().(*model.PdfFieldSignature)
    if !ok || sigField.V == nil || sigField.V.Contents == nil {
        continue
    }
    if err := ltv.Enable(sigField.V, nil); err != nil {
        log.Fatal(err)
    }
}

// Write output file.
if err = appender.WriteToFile("output.pdf"); err != nil {
    log.Fatal(err)
}
```

#### Protect validation data by adding a timestamp signature.

Optionally, an extra revision consisting of a timestamp signature, in
order to protect the added validation information, can be added.

A timestamp signature to protect the document DSS can be added in both workflows:

Workflow #1
```
    Revision 1: apply signature and LTV enable the signing chain.
    Revision 2: apply timestamp signature.
```

Workflow #2:
```
    Revision 1: apply signature.
    Revision 2: LTV enable the added signature and apply a timestamp signature.
```

`NOTE`: the timestamp signature itself is not LTV enabled. In order to LTV
enable the timestamp signature, a new revision has to be added (in both
workflows), in order to add validation data for the applied timestamp
signature.

For more information, see [pdf_sign_ltv_timestamp_revision](pdf_sign_ltv_timestamp_revision.go).

#### Customize LTV client

The HTTP clients used by `LTV.CertClient`, `LTV.OCSPClient` and `LTV.CRLClient`
can be customized.

```go
ltv, err := model.NewLTV(appender)
if err != nil {
    log.Fatal(err)
}

// Set timeout for HTTP requests.
ltv.CertClient.HTTPClient.Timeout = 300 * time.Millisecond

// Set custom HTTP client.
ltv.OCSPClient.HTTPClient = &http.Client{
    Transport: &http.Transport{
        MaxIdleConns:      10,
        IdleConnTimeout:   30 * time.Second,
        ForceAttemptHTTP2: true,
    },
}

// Set HTTP client proxy.
proxyURL, err := url.Parse("https://proxy-addr:proxy-port")
if err != nil {
    log.Fatal(err)
}

ltv.CRLClient.HTTPClient = &http.Client{
    Transport: &http.Transport{
        Proxy: http.ProxyURL(proxyURL),
    },
    Timeout: 5 * time.Second,
}
```
