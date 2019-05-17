/*
 * Rotate pages in a PDF file.  Degrees needs to be a multiple of 90.
 * Example of how to manipulate pages with the pdf creator.
 *
 * Run as: go run pdf_rotate.go output.pdf <angle> input.pdf
 * The angle is specified in degrees.
 */

package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	unicommon "github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/creator"
	pdf "github.com/unidoc/unipdf/v3/model"
)

func init() {
	// Use debug-mode log level.
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))
}

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Usage: go run pdf_rotate.go input.pdf <angle> output.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[3]

	degrees, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		fmt.Printf("Invalid degrees: %v\n", err)
		os.Exit(1)
	}
	if degrees%90 != 0 {
		fmt.Printf("Degrees needs to be a multiple of 90\n")
		os.Exit(1)
	}

	err = rotatePdf(inputPath, degrees, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

// Rotate all pages by 90 degrees.
func rotatePdf(inputPath string, degrees int64, outputPath string) error {
	c := creator.New()

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

	// Try decrypting both with given password and an empty one if that fails.
	if isEncrypted {
		auth, err := pdfReader.Decrypt([]byte(""))
		if err != nil {
			return err
		}
		if !auth {
			return errors.New("Unable to decrypt pdf with empty pass")
		}
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

		err = c.AddPage(page)
		if err != nil {
			return err
		}

		_ = c.RotateDeg(degrees)
	}

	err = c.WriteToFile(outputPath)
	return err
}
