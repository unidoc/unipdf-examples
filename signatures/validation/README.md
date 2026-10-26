# Validation

Examples of validating the signatures in a PDF. See the [signatures README](../README.md) for the full list.

Validation answers three separate questions, and each example covers a
different mix of them:

- Integrity. The signature matches the bytes it covers, under the certificate
  embedded in the PDF. `PdfReader.ValidateSignatures` checks this. Any
  self-signed certificate passes.
- Trust. The signer certificate chains to a root you chose, was valid at the
  validation time and, if requested, was not revoked.
  `PdfReader.ValidateSignaturesWithOpts` with `SignatureValidationOpts.Roots`
  checks this. Roots always come from the caller, never from the document.
- Coverage. The signed bytes are the document the reader sees. Content
  appended in a later revision does not break the signature, so it has to be
  checked separately.

Trust needs a PEM bundle of the root CA certificates you accept as
`ROOTS_PEM_PATH`.

## pdf_sign_validate.go

Checks integrity for every signature: `adbe.x509.rsa_sha1`,
`adbe.pkcs7.detached`, PAdES (`ETSI.CAdES.detached`, B-B to B-LTA) and
`ETSI.RFC3161` document timestamps. Exits non-zero if any signature fails.

```bash
$ go run pdf_sign_validate.go input_signed.pdf
```

## pdf_sign_validate_trust.go

Checks trust against `ROOTS_PEM_PATH`. Revocation data is read from the
document (DSS, PAdES revocation attribute) by default, or also fetched online
with `online`. An unknown revocation status counts as untrusted.

When the document embeds only the signer certificate, the chain can't be built
and trust fails with "certificate signed by unknown authority". Pass the
intermediate CA certificates as a fourth argument. The `Certificate issuer:`
line in the output, which comes from `SignerCertificate()`, names the
certificate that is missing. When a signer is not trusted, the example also
prints its serial number.

```bash
$ go run pdf_sign_validate_trust.go input_signed.pdf roots.pem
$ go run pdf_sign_validate_trust.go input_signed.pdf roots.pem offline intermediates.pem
```

## pdf_sign_validate_online.go

Checks trust with online CRL and OCSP lookups through configured HTTP clients.
Revocation data in the document is used first. OCSP responders are tried
before CRLs, and a failed lookup is reported in `TrustWarnings`. A revocation
status that is still unknown after the lookups makes the signature untrusted.

```bash
$ go run pdf_sign_validate_online.go input_signed.pdf roots.pem
```

## pdf_sign_validate_at_time.go

Checks trust at a point in time other than now, for archived documents whose
certificates have since expired. By default the current time is used, or the
time of a trusted timestamp token. Pass an RFC 3339 time, or `signing-time` to
use the signer's own /M date for signatures without a trusted timestamp. The
signer chooses /M, so use it only when you accept that.

```bash
$ go run pdf_sign_validate_at_time.go input_signed.pdf roots.pem
$ go run pdf_sign_validate_at_time.go input_signed.pdf roots.pem 2025-06-01T12:00:00Z
```

## pdf_sign_validate_coverage.go

Checks coverage for every signature. It reports whether the signature covers
the whole file or only its own revision, how many revisions were appended
after it, and whether the bytes left out of the signed range are only the
signature's /Contents value. Appended revisions are reported but don't fail,
since a second signature or LTV data also adds one. Use `PdfReader.GetRevision`
(see [../signing/pdf_sign_get_revision.go](../signing/pdf_sign_get_revision.go))
to see what they changed.

```bash
$ go run pdf_sign_validate_coverage.go input_signed.pdf
```
