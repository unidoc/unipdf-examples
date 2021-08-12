/*
 * Check UniPDF License key info for an offline license key.
 *
 * Run as: go run unipdf_offline_license_info.go
 */

package main

import (
	"fmt"

	"github.com/unidoc/unipdf/v3/common/license"
)

const offlineLicenseKey = `
-----BEGIN UNIDOC LICENSE KEY-----
Key contents here.
-----END UNIDOC LICENSE KEY-----
`

func init() {
	// Load the offline license key.
	err := license.SetLicenseKey(offlineLicenseKey, `Company Name`)
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
