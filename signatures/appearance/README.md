# Signature appearance

Examples of visible signatures. See the [signatures README](../README.md) for the full list.

A signature appearance is what a viewer draws in the signature widget. The
signer writes it, so a name or date shown there is not evidence of who signed
or when. Validation reads those from the signature and its certificate. Every
example here signs with a generated key and a self-signed certificate, so it
runs offline.

## pdf_sign_with_appearance.go

Signs with `PdfAppender.SignWithAppearance` and `annotator.FieldAppearance`.
The appearance is generated from the Name, M (signing date), Reason and
Location entries of the signature dictionary, so the visible text matches what
was signed. Set those entries before calling `SignWithAppearance`. The field
needs a widget rectangle on page 1. Running the example on its own output adds
a new field, `Approval 2` and so on, above the previous one.

```bash
$ go run pdf_sign_with_appearance.go input.pdf input_signed.pdf
```

## pdf_sign_appearance.go

Builds text appearances by hand with `annotator.NewSignatureField` and a list
of `SignatureLine`s, and sets the font size, colours and border through
`SignatureFieldOpts`. It adds three fields on every page. They share one
signature dictionary, so the document has a single signature shown in several
places. The lines are free text, so they can differ from the signature
dictionary.

```bash
$ go run pdf_sign_appearance.go input.pdf input_signed.pdf
```

## pdf_image_sign_appearance.go

Adds image appearances: an image beside, above or below the text lines, an
image-only appearance, and a watermark image behind each field. Five fields go
on every page, sharing one signature dictionary. The signature and watermark
images can be PNG or JPEG.

```bash
$ go run pdf_image_sign_appearance.go input.pdf signature.png watermark.png input_signed.pdf
```
