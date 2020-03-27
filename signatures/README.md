# Digital signatures.

Examples for digital signing of PDF files with UniDoc:
- [pdf_sign_generate_keys.go](pdf_sign_generate_keys.go) Example of signing using generated private/public key pair.
- [pdf_sign_pkcs12.go](pdf_sign_pkcs12.go) Example of signing using PKCS12 (.p12/.pfx) file.
- [pdf_sign_external.go](pdf_sign_external.go) Example of PKCS7 signing with an external service with an interim step, creating a PDF with a blank signature and then replacing the blank signature with the actual signature from the signing service.
- [pdf_sign_hsm_pkcs11.go](pdf_sign_hsm_pkcs11.go) Example of signing with a PKCS11 service using SoftHSM and the crypto11 package.
- [pdf_sign_appearance.go](pdf_sign_appearance.go) Example of creating signature appearance fields.
- [pdf_sign_validate.go](pdf_sign_validate.go) Example of signature validation.
- [pdf_sign_pem_multicert.go](pdf_sign_pem_multicert.go) Example of signing using a certificate chain and a private key, extracted from PEM files.

## pkcs_sign_hsm_pkcs11.go

The code example shows how to sign with a HSM via PKCS11 as supported by the
crypto11 library.  
The example uses SoftHSM which is great for testing digital signatures via
PKCS11 without any hardware requirements.

#### Prerequisites

Ubuntu/Debian
```bash
$ sudo apt-get install libssl-dev
$ sudo apt-get install autotools-dev
$ sudo apt-get install autoconf
$ sudo apt-get install libtool
```

CentOS/RHEL
```bash
$ sudo yum group install "Development Tools"
$ sudo yum install openssl-devel
```

#### Installation

```bash
$ git clone https://github.com/opendnssec/SoftHSMv2.git
$ cd SoftHSMv2
$ sh autogen.sh
$ ./configure
$ make
$ sudo make install
```

#### Configuration

```bash
$ mkdir -p /home/user/.config/softhsm2/tokens
$ cd /home/user/.config/softhsm2
$ touch softhsm2.conf
$ export SOFTHSM2_CONF=/home/user/.config/softhsm2/softhsm2.conf
```

#### Contents of softhsm2.conf

```
directories.tokendir = /home/user/.config/softhsm2/tokens
objectstore.backend = file
log.level = DEBUG
slots.removable = true
```

#### Create token

Creating a token "test", selecting the PIN numbers as prompted

```bash
$ softhsm2-util --init-token --slot 0 --label "test"
```

#### Usage

Create a key pair:
```bash
$ go run pdf_sign_hsm_pkcs11.go add test <PIN> <KEYPAIR_LABEL>
```

Sign PDF file:
```bash
$ go run pdf_sign_hsm_pkcs11.go sign test <PIN> <KEYPAIR_LABEL> input.pdf input_signed.pdf
```

Signed output is in `input_signed.pdf`
