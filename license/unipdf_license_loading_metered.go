/*
 * unipdf_license_loading_metered.go:
 * Illustrates how to load a metered license API key.
 * Free api keys can be obtained at: https://cloud.unidoc.io
 *
 * Run as: go run unipdf_license_loading_metered.go
 */

package main

import (
	"fmt"

	"github.com/unidoc/unipdf/v4/common/license"
)

func init() {
	// To get your free API key for metered license, sign up on: https://cloud.unidoc.io
	// Make sure to be using UniPDF v3.19.1 or newer for Metered API key support.
	err := license.SetMeteredKey(`my metered api key goes here`)
	if err != nil {
		fmt.Printf("ERROR: Failed to set metered key: %v\n", err)
		fmt.Printf("Make sure to get a valid key from https://cloud.unidoc.io\n")
		panic(err)
	}
}

func main() {
	lk := license.GetLicenseKey()
	if lk == nil {
		fmt.Printf("Failed retrieving license key")
		return
	}
	fmt.Printf("License: %s\n", lk.ToString())

	// GetMeteredState freshly checks the state, contacting the licensing server.
	state, err := license.GetMeteredState()
	if err != nil {
		fmt.Printf("ERROR getting metered state: %+v\n", err)
		panic(err)
	}
	fmt.Printf("State: %+v\n", state)
	if state.OK {
		fmt.Printf("State is OK\n")
	} else {
		fmt.Printf("State is not OK\n")
	}
	fmt.Printf("Credits: %v\n", state.Credits)
	fmt.Printf("Used credits: %v\n", state.Used)
}
