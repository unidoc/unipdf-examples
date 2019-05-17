/*
 * Check access permissions for a specified PDF.
 *
 * Run as: go run pdf_check_permissions.go input.pdf [password]
 */

package main

import (
	"fmt"
	"os"

	unicommon "github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/core/security"
	pdf "github.com/unidoc/unipdf/v3/model"
)

func init() {
	// Set debug-level logging via console.
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run pdf_check_permissions.go input.pdf [password]\n")
		os.Exit(0)
	}
	filepath := os.Args[1]

	password := "" // Default try empty pass if not specified.
	if len(os.Args) >= 3 {
		password = os.Args[2]
	}

	err := printAccessInfo(filepath, password)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func printAccessInfo(inputPath string, password string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	canView, perms, err := pdfReader.CheckAccessRights([]byte(password))
	if err != nil {
		return err
	}

	if !canView {
		fmt.Printf("%s - Cannot view - No access with the specified password\n", inputPath)
		return nil
	}

	fmt.Printf("Input file %s\n", inputPath)
	fmt.Printf("Access Permissions: %+v\n", perms)
	fmt.Printf("--------\n")

	// Print a text summary of the flags.
	booltext := map[bool]string{false: "No", true: "Yes"}
	allowed := func(p security.Permissions) string {
		return booltext[perms.Allowed(p)]
	}

	fmt.Printf("Printing allowed? - %s\n", allowed(security.PermPrinting))
	if perms.Allowed(security.PermPrinting) {
		fmt.Printf("Full print quality (otherwise print in low res)? - %s\n", allowed(security.PermFullPrintQuality))
	}
	fmt.Printf("Modifications allowed? - %s\n", allowed(security.PermModify))
	fmt.Printf("Allow extracting graphics? %s\n", allowed(security.PermExtractGraphics))
	fmt.Printf("Can annotate? - %s\n", allowed(security.PermAnnotate))
	if perms.Allowed(security.PermAnnotate) {
		fmt.Printf("Can fill forms? - Yes\n")
	} else {
		fmt.Printf("Can fill forms? - %s\n", allowed(security.PermFillForms))
	}
	fmt.Printf("Extract text, graphics for users with disabilities? - %s\n", allowed(security.PermDisabilityExtract))

	return nil
}
