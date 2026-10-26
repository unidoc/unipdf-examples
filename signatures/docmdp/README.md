# DocMDP

Examples of certification signatures with DocMDP permissions. See the [signatures README](../README.md) for the full list.

A DocMDP signature certifies the document and states which changes are
allowed in later revisions. The permission is one of:

- `mdp.NoChanges` (P=1). No changes are allowed.
- `mdp.FillForms` (P=2). Filling in forms, instantiating page templates and
  signing are allowed.
- `mdp.FillFormsAndAnnots` (P=3). As P=2, plus creating, deleting and changing
  annotations.

A document can have only one DocMDP signature, and it must be the first
signature. To check later revisions against the permission, read it from the
document with `reader.GetPerms().DocMDP.GetDocMDPPermission()` and pass it to
`sighandler.NewDocMDPHandler`, wrapping the inner handler. The handler
enforces the permission it is given, not the document's /P, so a looser value
would accept changes the document forbids. A disallowed change
leaves the signed bytes intact, but the handler reports the signature as
invalid and lists the changes in `SignatureValidationResult.DiffResults`. All
examples sign with a generated key pair, so they run offline.

## pdf_sign_docmdp.go

Signs with `mdp.FillForms` (P=2).

```bash
$ go run pdf_sign_docmdp.go input.pdf input_signed.pdf
```

## pdf_sign_docmdp_valid_changes.go

Signs with `mdp.FillFormsAndAnnots` (P=3), adds a square annotation to page 1
in a new revision, then validates. It prints that no DocMDP violation was
detected.

```bash
$ go run pdf_sign_docmdp_valid_changes.go input.pdf input_changed.pdf
```

## pdf_sign_docmdp_invalid_changes.go

Signs with `mdp.FillFormsAndAnnots` (P=3), adds a page in a new revision, then
validates. Adding pages is not allowed at P=3, so it prints that a DocMDP
violation was detected.

```bash
$ go run pdf_sign_docmdp_invalid_changes.go input.pdf input_changed.pdf
```
