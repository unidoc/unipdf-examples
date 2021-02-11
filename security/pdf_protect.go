/*
 * Protects PDF files by setting a password on it. This example both sets user
 * and opening password and hard-codes the protection bits here, but easily adjusted
 * in the code here although not on the command line.
 *
 * The user-pass is a password required to view the file with the access specified by certain permission flags (specified
 * in the code example below), whereas the owner pass is needed to have full access to the file.
 * See pdf_check_permissions.go for an example about checking the permissions for a given PDF file.
 *
 * If anyone is supposed to be able to read the PDF under the given access restrictions, then the user password should
 * be left empty ("").
 *
 * Run as: go run pdf_protect.go input.pdf <user-pass> <owner-pass> output.pdf
 * Sets a user and owner password for the PDF.
 */

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/core/security"
	pdf "github.com/unidoc/unipdf/v3/model"
)

const licenseKey = `
-----BEGIN UNIDOC LICENSE KEY-----
Free trial license keys are available at: https://unidoc.io/
-----END UNIDOC LICENSE KEY-----
`

func init() {
	// Enable debug-level logging.
	// unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	err := license.SetLicenseKey(licenseKey, `Company Name`)
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run pdf_protect.go input.pdf <user-pass> <owner-pass> output.pdf")
		fmt.Println("Sets a user and owner password for the PDF.")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	userPassword := os.Args[2]
	ownerPassword := os.Args[3]
	outputPath := os.Args[4]

	err := protectPdf(inputPath, outputPath, userPassword, ownerPassword)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

func protectPdf(inputPath string, outputPath string, userPassword, ownerPassword string) error {
	permissions := security.PermPrinting | // Allow printing with low quality
		security.PermFullPrintQuality |
		security.PermModify | // Allow modifications.
		security.PermAnnotate | // Allow annotations.
		security.PermFillForms |
		security.PermRotateInsert | // Allow modifying page order, rotating pages etc.
		security.PermExtractGraphics | // Allow extracting graphics.
		security.PermDisabilityExtract // Allow extracting graphics (accessibility)

	encryptOptions := &pdf.EncryptOptions{
		Permissions: permissions,
	}

	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return err
	}
	if isEncrypted {
		return fmt.Errorf("The PDF is already locked (need to unlock first)")
	}

	// Generate a PdfWriter instance from existing PdfReader.
	pdfWriter, err := pdfReader.ToWriter(nil)
	if err != nil {
		return err
	}

	// Encrypt document before writing to file.
	err = pdfWriter.Encrypt([]byte(userPassword), []byte(ownerPassword), encryptOptions)
	if err != nil {
		return err
	}

	// Write to file.
	err = pdfWriter.WriteToFile(outputPath)
	return err
}
