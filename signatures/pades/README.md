# PAdES

Examples of PAdES (ETSI EN 319 142) signatures, one per baseline level. See the [signatures README](../README.md) for the full list.

Each level adds to the one before it:

- B-B is a basic CAdES signature (`ETSI.CAdES.detached`) with the signing
  certificate as a signed attribute.
- B-T adds a signature timestamp from a TSA, which proves the signature existed
  at that time.
- B-LT adds the validation data (certificates, CRLs and OCSP responses) to the
  document DSS, so it can be validated later without going online.
- B-LTA adds a document timestamp over the validation data, which can be
  renewed before the algorithms or certificates involved become weak.

All examples take the signing key and certificate from a PKCS12 file and the
issuing CA certificate from a PEM file. RSA and ECDSA keys both work. B-T and
above request timestamps from https://freetsa.org/tsr, and B-LT and above
fetch CRL and OCSP responses for the chain, so they need network access.

To validate the output, use
[../validation/pdf_sign_validate.go](../validation/pdf_sign_validate.go) for
integrity or
[../validation/pdf_sign_validate_trust.go](../validation/pdf_sign_validate_trust.go)
for trust.

## pdf_sign_pades_b_b.go

```bash
$ go run pdf_sign_pades_b_b.go signer.pfx <PASSWORD> ca.pem input.pdf input_signed.pdf
```

## pdf_sign_pades_b_t.go

```bash
$ go run pdf_sign_pades_b_t.go signer.pfx <PASSWORD> ca.pem input.pdf input_signed.pdf
```

## pdf_sign_pades_b_lt.go

Signs in one revision, then writes the collected validation data to the DSS
in a second revision.

```bash
$ go run pdf_sign_pades_b_lt.go signer.pfx <PASSWORD> ca.pem input.pdf input_signed.pdf
```

## pdf_sign_pades_b_lta.go

Same as B-LT, followed by a third revision with an `ETSI.RFC3161` document
timestamp.

```bash
$ go run pdf_sign_pades_b_lta.go signer.pfx <PASSWORD> ca.pem input.pdf input_signed.pdf
```

## pdf_sign_custom_client.go

B-LT signing through the setter-based `sighandler.NewEtsiPAdES(sighandler.LevelLT)`
API, with a custom timestamp client whose HTTP client has a 30 second timeout.
Use this form to set your own TSA, proxy or timeouts.

```bash
$ go run pdf_sign_custom_client.go signer.pfx <PASSWORD> ca.pem input.pdf input_signed.pdf
```
