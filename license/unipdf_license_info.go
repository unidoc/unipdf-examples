/*
 * Check UniPDF License key info.
 *
 * Run as: go run unipdf_license_info.go
 */

package main

import (
	"fmt"

	"github.com/unidoc/unipdf/v3/common/license"
)

const licenseKey = `
-----BEGIN UNIDOC LICENSE KEY-----
Free trial license keys are available at: https://unidoc.io/
-----END UNIDOC LICENSE KEY-----
`

func init() {
	err := license.SetLicenseKey(licenseKey, `Company Name`)
	if err != nil {
		panic(err)
	}
}

func main() {
	lk := license.GetLicenseKey()
	if lk == nil {
		fmt.Printf("Failed retrieving license key")

		return
	}

	fmt.Printf("%s\n", lk.ToString())
}
