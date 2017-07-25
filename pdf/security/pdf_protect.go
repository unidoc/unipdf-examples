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

	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

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
	pdfWriter := pdf.NewPdfWriter()

	permissions := pdfcore.AccessPermissions{}
	// Allow printing with low quality
	permissions.Printing = true
	permissions.FullPrintQuality = false
	// Allow modifications.
	permissions.Modify = true
	// Allow annotations.
	permissions.Annotate = true
	permissions.FillForms = true
	// Allow modifying page order, rotating pages etc.
	permissions.RotateInsert = true
	// Allow extracting graphics.
	permissions.ExtractGraphics = true
	// Allow extracting graphics (accessibility)
	permissions.DisabilityExtract = true

	encryptOptions := &pdf.EncryptOptions{}
	encryptOptions.Permissions = permissions

	err := pdfWriter.Encrypt([]byte(userPassword), []byte(ownerPassword), encryptOptions)
	if err != nil {
		return err
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

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return err
		}

		err = pdfWriter.AddPage(page)
		if err != nil {
			return err
		}
	}

	fWrite, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer fWrite.Close()

	err = pdfWriter.Write(fWrite)
	if err != nil {
		return err
	}

	return nil
}
