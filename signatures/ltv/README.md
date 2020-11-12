# LTV

## Table of contents

- [Overview](#overview)
- [LTV enable workflows](#ltv-enable-workflows)
  - [Sign and LTV enable in the same revision](#1-sign-and-ltv-enable-in-the-same-revision)
  - [Sign and LTV enable in a separate revision](#2-sign-and-ltv-enable-in-a-separate-revision)
- [Usage](#usage)
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
    CLRs  [PDF stream array]

    // Signature specific validation data.
    // Each key in the VRI dictionary represents the hash of the Contents
    // field of a signature dictionary. Each signature entry in the VRI
    // dictionary has its own set of certificates, OCSP and CRL responses,
    // used to validate that particular signature exclusively.
    VRI:
        Hash signature 1: Cert [PDF stream array]
                          OCSP [PDF stream array]
                          CLR  [PDF stream array]
        ...
        Hash signature N: Cert [PDF stream array]
                          OCSP [PDF stream array]
                          CLR  [PDF stream array]
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
least the signing certificate). The LTV client builds the certificate chain
up to a trusted root certificate (by downloading any missing certificates).

For more information, see [pdf_sign_ltv_one_revision.go](pdf_sign_ltv_one_revision.go).

```go
// Get private key and X509 certificate from a PKCS12 (.p12/.pfx) file.
pfxData, err := ioutil.ReadFile("p12.pfx")
if err != nil {
    log.Fatal(err)
}

priv, cert, err := pkcs12.Decode(pfxData, password)
if err != nil {
    log.Fatal(err)
}

file, err := os.Open("unsigned-file.pdf")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

// Create appender.
appender, err := model.NewPdfAppender(reader)
if err != nil {
    log.Fatal(err)
}

// Create signature handler.
handler, err := sighandler.NewAdobePKCS7Detached(priv.(*rsa.PrivateKey), cert)
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
field.T = core.MakeString("Test Sign LTV enable")

if err = appender.Sign(1, field); err != nil {
    log.Fatal(err)
}

// LTV enable the certificate chain used to apply the signature.
ltv, err := model.NewLTV(appender)
if err != nil {
    log.Fatal(err)
}
if err := ltv.EnableChain([]*x509.Certificate{cert}); err != nil {
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
extracted from the signature dictionaries.  The LTV client builds the
certificate chains up to a trusted root certificate, by downloading any
missing certificates.

For more information, see [pdf_ltv_enable_signed_file.go](pdf_ltv_enable_signed_file.go)
and [pdf_sign_ltv_extra_revision.go](pdf_sign_ltv_extra_revision.go).

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

if err := ltv.EnableAll(nil); err != nil {
    log.Fatal(err)
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
