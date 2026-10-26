# Signing

Examples of signing a PDF with keys loaded from different sources. See the [signatures README](../README.md) for the full list.

Each signature is added in an incremental revision on top of the file it
signs. Most examples sign the input as it is, so its bytes stay as they were.
`pdf_sign_new_page.go` and `pdf_sign_and_encrypt_pdf.go` first rewrite the
input with `PdfReader.ToWriter` (to add a page, or to encrypt it) and sign
the rewritten file. Examples that generate their own key pair use a
self-signed certificate, which proves integrity but not identity. To check
that a signer is trusted, see [../validation](../validation).

## pdf_sign_generate_keys.go

Signs with a generated RSA key pair and an `adbe.pkcs7.detached` signature.
The shortest signing example.

```bash
$ go run pdf_sign_generate_keys.go input.pdf input_signed.pdf
```

## pdf_sign_pkcs12.go

Signs with the RSA key and certificate in a PKCS12 (.p12/.pfx) file. Files
that also bundle CA certificates load, but only the signing certificate is
embedded. To embed the chain, see `pdf_sign_pem_multicert.go`.

```bash
$ go run pdf_sign_pkcs12.go signer.p12 <PASSWORD> input.pdf input_signed.pdf
```

## pdf_sign_pem_multicert.go

Signs with a private key and a certificate chain read from PEM files. The key
can be PKCS#1 or PKCS#8. The first certificate in the chain signs and the whole
chain is embedded. A file with only the signing certificate also works.

```bash
$ go run pdf_sign_pem_multicert.go input.pdf input_signed.pdf certs.pem key.pem
```

## pdf_sign_pkcs7_ecdsa.go

Signs with an ECDSA key from a PKCS12 file. ECDSA signatures are defined by
PDF 2.0, so viewers that support only earlier versions may not validate them.

```bash
$ go run pdf_sign_pkcs7_ecdsa.go signer.pfx <PASSWORD> input.pdf input_signed.pdf
```

## pdf_sign_timestamp.go

Signs with a PKCS12 file, then adds an RFC 3161 document timestamp in a second
revision. The timestamp comes from https://freetsa.org/tsr, so it needs network
access.

```bash
$ go run pdf_sign_timestamp.go signer.p12 <PASSWORD> input.pdf input_signed.pdf
```

## pdf_sign_new_page.go

Appends a new page to the document and places a visible signature on it.

```bash
$ go run pdf_sign_new_page.go input.pdf input_signed.pdf
```

## pdf_sign_twice_visible_annotation.go

Signs the last page twice, each time in its own revision, with its own key
pair and a visible annotation. The first signature still verifies after the
second is added because the second revision is appended.

```bash
$ go run pdf_sign_twice_visible_annotation.go input.pdf input_signed.pdf
```

## pdf_sign_and_encrypt_pdf.go

Encrypts the document with AES-128, then signs the encrypted document in a
revision that keeps the encryption. The output opens with the password
`password`.

```bash
$ go run pdf_sign_and_encrypt_pdf.go input.pdf input_signed.pdf
```

## pdf_sign_get_revision.go

Signs the input twice, then uses `PdfReader.GetRevision` to write the revision
the input had before signing. Use the same call to see what a signature
covered when later revisions were appended.

```bash
$ go run pdf_sign_get_revision.go input.pdf input_revision.pdf
```
