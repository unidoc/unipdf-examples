# Digital signatures.

Examples for digital signing of PDF files with UniDoc.

## Signing

- [pdf_sign_generate_keys.go](signing/pdf_sign_generate_keys.go) Example of signing using generated private/public key pair.
- [pdf_sign_pkcs12.go](signing/pdf_sign_pkcs12.go) Example of signing using PKCS12 (.p12/.pfx) file.
- [pdf_sign_pem_multicert.go](signing/pdf_sign_pem_multicert.go) Example of signing using a certificate chain and a private key, extracted from PEM files.
- [pdf_sign_pkcs7_ecdsa.go](signing/pdf_sign_pkcs7_ecdsa.go) Example of adbe.pkcs7.detached signing with an ECDSA key.
- [pdf_sign_timestamp.go](signing/pdf_sign_timestamp.go) Example of signing and timestamping using PKCS12 (.p12/.pfx) file.
- [pdf_sign_custom_client.go](signing/pdf_sign_custom_client.go) Example of signing with a custom timestamp client.
- [pdf_sign_new_page.go](signing/pdf_sign_new_page.go) Example of appending a new page with signature to a PDF document.
- [pdf_sign_twice_visible_annotation.go](signing/pdf_sign_twice_visible_annotation.go) Example of signing a PDF document twice with visible signature annotations.
- [pdf_sign_and_encrypt_pdf.go](signing/pdf_sign_and_encrypt_pdf.go) Example of signing an encrypted PDF document.

## Signature appearance

- [pdf_sign_appearance.go](appearance/pdf_sign_appearance.go) Example of creating signature appearance fields.
- [pdf_image_sign_appearance.go](appearance/pdf_image_sign_appearance.go) Example of creating signature appearance fields with images.

## External signing

- [pdf_sign_external.go](external/pdf_sign_external.go) Example of PKCS7 signing with an external service with an interim step, creating a PDF with a blank signature and then replacing the blank signature with the actual signature from the signing service.
- [pdf_sign_external_aws_kms.go](external/pdf_sign_external_aws_kms.go) Example of signing with AWS KMS.
- [pdf_sign_external_google_cloud_kms.go](external/pdf_sign_external_google_cloud_kms.go) Example of signing with Google Cloud KMS.
- [pdf_sign_external_globalsign.go](external/pdf_sign_external_globalsign.go) Example of signing with the GlobalSign DSS API.
- [pdf_sign_hsm_pkcs11_cgo.go](external/pdf_sign_hsm_pkcs11_cgo.go) Example of signing with a PKCS11 service using SoftHSM and the crypto11 package. See the [external signing guide](external) for setup.

## PAdES

- [pdf_sign_pades_b_b.go](pades/pdf_sign_pades_b_b.go) Example of signing with a PAdES B-B compatible digital signature.
- [pdf_sign_pades_b_t.go](pades/pdf_sign_pades_b_t.go) Example of signing with a PAdES B-T compatible digital signature.
- [pdf_sign_pades_b_lt.go](pades/pdf_sign_pades_b_lt.go) Example of signing with a PAdES B-LT compatible digital signature.
- [pdf_sign_pades_b_lta.go](pades/pdf_sign_pades_b_lta.go) Example of signing with a PAdES B-LTA compatible digital signature.

## DocMDP

- [pdf_sign_docmdp.go](docmdp/pdf_sign_docmdp.go) Example of signing with the DocMDP restriction.
- [pdf_sign_docmdp_valid_changes.go](docmdp/pdf_sign_docmdp_valid_changes.go) Example of signing with the DocMDP restriction and adding valid changes.
- [pdf_sign_docmdp_invalid_changes.go](docmdp/pdf_sign_docmdp_invalid_changes.go) Example of signing with the DocMDP restriction and adding invalid changes.

## Validation

- [pdf_sign_validate.go](validation/pdf_sign_validate.go) Example of signature integrity validation.
- [pdf_sign_validate_pades_b_b.go](validation/pdf_sign_validate_pades_b_b.go) Example of PAdES signature integrity validation.
- [pdf_sign_validate_timestamp.go](validation/pdf_sign_validate_timestamp.go) Example of validating a signature that contains a timestamp.
- [pdf_sign_validate_trust.go](validation/pdf_sign_validate_trust.go) Example of signature validation against trusted root certificates with CRL/OCSP revocation checks.
- [pdf_sign_validate_at_time.go](validation/pdf_sign_validate_at_time.go) Example of validating archived signatures at an explicit point in time.
- [pdf_sign_validate_online.go](validation/pdf_sign_validate_online.go) Example of signature validation with online CRL and OCSP revocation checks through configured HTTP clients.
- [pdf_sign_get_revision.go](validation/pdf_sign_get_revision.go) Example of retrieving a specific or previous revision of a signed PDF document.

## LTV

For LTV enabling digital signatures, see the [LTV](ltv) guide and samples.
