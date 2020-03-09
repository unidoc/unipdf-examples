# PDF Security Management

The example explains how you can use UniPDF to secure your PDF document. UniPDF allows you to check the security condition of existing PDF documents, set protection for new or existing PDF documents. You can also decrypt encrypted documents with a given password.  

## Examples

- [pdf_check_permissions.go](pdf_check_permissions.go) The example checks access permissions for a specified PDF.
- [pdf_protect.go](pdf_protect.go) The example showcases how to protect PDF files by setting a password on it using UniPDF. This example both sets user and opening password and hard-codes the protection bits here, but easily adjusted in the code here although not on the command line.
- [pdf_security_info.go](pdf_security_info.go) The example outputs protection information about locked PDFs.
- [pdf_unlock.go](pdf_unlock.go) The example showcases how to unlock PDF files using UniPDF and it tries to decrypt encrypted documents with the given password, if that fails it tries an empty password as best effort.
