# External signing

Examples of signing with external services, cloud KMS and HSMs. See the [signatures README](../README.md) for the full list.

## pdf_sign_external.go

Writes the PDF with an empty signature, takes the bytes covered by its
ByteRange, gets a detached PKCS7 signature for them and places it into the
Contents of that same PDF. The external service is simulated locally with a
generated key and self-signed certificate, so it runs offline.

```bash
$ go run pdf_sign_external.go input.pdf input_signed.pdf
```

## pdf_sign_external_aws_kms.go

Signs with an asymmetric AWS KMS key. The key needs key usage `SIGN_VERIFY`
and an RSA key spec (`RSA_2048`, `RSA_3072` or `RSA_4096`); the example signs
with `RSASSA_PKCS1_V1_5_SHA_256` and fails on any other key type. Credentials
and region come from the standard AWS configuration (`AWS_PROFILE` or
`AWS_ACCESS_KEY_ID`/`AWS_SECRET_ACCESS_KEY`, and `AWS_REGION`, or
`~/.aws/credentials` and `~/.aws/config`). The identity needs `kms:GetPublicKey`
and `kms:Sign` on the key.

```bash
$ export AWS_REGION=us-west-1
$ go run pdf_sign_external_aws_kms.go input.pdf input_signed.pdf <KEY_ID>
```

## pdf_sign_external_google_cloud_kms.go

Signs with a Google Cloud KMS key version. The key needs purpose
`ASYMMETRIC_SIGN` and algorithm `RSA_SIGN_PKCS1_2048_SHA256`,
`RSA_SIGN_PKCS1_3072_SHA256` or `RSA_SIGN_PKCS1_4096_SHA256`. Credentials are
a service account JSON key file passed on the command line; the account needs
the `cloudkms.signerVerifier` role on the key.

```bash
$ go run pdf_sign_external_google_cloud_kms.go input.pdf input_signed.pdf \
    service-account.json \
    projects/<PROJECT>/locations/<LOCATION>/keyRings/<KEY_RING>/cryptoKeys/<KEY>/cryptoKeyVersions/<VERSION>
```

## pdf_sign_external_globalsign.go

Signs through the GlobalSign Digital Signing Service API
(`https://emea.api.dss.globalsign.com:8443`). It needs the API key and secret
of a DSS account, plus the mTLS client certificate and private key files
GlobalSign issues with that account. The signature embeds the GlobalSign
certificate chain, the OCSP response and a GlobalSign timestamp, and the
example adds LTV data to the DSS.

```bash
$ go run pdf_sign_external_globalsign.go input.pdf input_signed.pdf \
    <API_KEY> <API_SECRET> tls_cert.pem tls_key.pem
```

## pdf_sign_hsm_pkcs11_cgo.go

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

Creating a token labeled "test", selecting the PIN numbers as prompted

```bash
$ softhsm2-util --init-token --slot 0 --label "test"
```

#### Usage

The example loads the PKCS11 library from `/usr/local/lib/softhsm/libsofthsm2.so`.
Set `SOFTHSM2_LIB` if it is installed elsewhere, for example
`/opt/homebrew/lib/softhsm/libsofthsm2.so` with Homebrew or
`/usr/lib/softhsm/libsofthsm2.so` on Debian/Ubuntu packages.

```bash
$ export SOFTHSM2_LIB=/opt/homebrew/lib/softhsm/libsofthsm2.so
```

`test` in the commands below is the token label used with `--label` above.
Create a key pair:
```bash
$ go run pdf_sign_hsm_pkcs11_cgo.go add test <PIN> <KEYPAIR_LABEL>
```

Sign PDF file:
```bash
$ go run pdf_sign_hsm_pkcs11_cgo.go sign test <PIN> <KEYPAIR_LABEL> input.pdf input_signed.pdf
```

Signed output is in `input_signed.pdf`.
