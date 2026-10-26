# PDF Security Management

The example explains how you can use UniPDF to secure your PDF document. UniPDF allows you to check the security condition of existing PDF documents, set protection for new or existing PDF documents. You can also decrypt encrypted documents with a given password.  

## Examples

- [pdf_check_permissions.go](pdf_check_permissions.go) The example checks access permissions for a specified PDF.
- [pdf_check_permissions_certificate.go](pdf_check_permissions_certificate.go) The example checks, without decrypting, whether a recipient's certificate and private key can open a PDF encrypted for recipients and which permissions it grants, then decrypts it with `PdfReader.DecryptWithCertificate`. It shows the errors for a certificate that is not a recipient (`security.ErrNotRecipient`), a private key that does not open the recipient's entry (`security.ErrRecipientNotOpened`), recipient entries in an unsupported format (`security.ErrUnsupportedRecipient`) and a password protected PDF (`core.ErrPasswordRequired`).
- [pdf_decrypt_with_certificate.go](pdf_decrypt_with_certificate.go) The example opens a PDF encrypted for recipients with a recipient's certificate and private key passed in `model.ReaderOpts`, and writes an unencrypted copy.
- [pdf_encrypt_for_recipients.go](pdf_encrypt_for_recipients.go) The example encrypts a PDF for one or more recipient certificates with the public-key security handler (Adobe.PubSec), so that recipients open it with their certificate and private key instead of a password. It uses AES-256 by default or AES-128, and grants the recipients either all permissions (owner-equivalent) or printing only.
- [pdf_generate_recipient_keys.go](pdf_generate_recipient_keys.go) The example generates self-signed RSA test certificates and private keys for trying out the public-key encryption examples.
- [pdf_protect.go](pdf_protect.go) The example showcases how to protect PDF files by setting a password on it using UniPDF. This example both sets user and opening password and hard-codes the protection bits here, but easily adjusted in the code here although not on the command line.
- [pdf_security_info.go](pdf_security_info.go) The example outputs protection information about locked PDFs.
- [pdf_unlock.go](pdf_unlock.go) The example showcases how to unlock PDF files using UniPDF and it tries to decrypt encrypted documents with the given password, if that fails it tries an empty password as best effort.
