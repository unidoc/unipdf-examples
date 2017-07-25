/*
 * Get specific object from a PDF by number.  Prints the trailer dictionary if no number specified.
 * For streams, prints out the decoded stream.
 *
 * Run as: go run pdf_get_object.go input.pdf [num]
 */

package main

import (
	"fmt"
	"os"
	"strconv"

	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Syntax: go run pdf_get_object.go input.pdf [num]")
		fmt.Println("If num is not specified, will display the trailer dictionary")
		os.Exit(1)
	}

	// Load a console logger in debug mode.
	//unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	inputPath := os.Args[1]

	objNum := -1
	if len(os.Args) > 2 {
		num, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		objNum = num
	}

	fmt.Printf("Input file: %s\n", inputPath)
	err := inspectPdfObject(inputPath, objNum)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func inspectPdfObject(inputPath string, objNum int) error {
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
		// If encrypted, try decrypting with an empty one.
		// Can also specify a user/owner password here by modifying the line below.
		auth, err := pdfReader.Decrypt([]byte(""))
		if err != nil {
			fmt.Printf("Decryption error: %v\n", err)
			return err
		}
		if !auth {
			fmt.Println(" This file is encrypted with opening password. Modify the code to specify the password.")
			return nil
		}
	}

	// Print trailer
	if objNum == -1 {
		trailer, err := pdfReader.GetTrailer()
		if err != nil {
			return err
		}

		fmt.Printf("Trailer: %s\n", trailer.String())
		return nil
	}

	obj, err := pdfReader.GetIndirectObjectByNumber(objNum)
	if err != nil {
		return err
	}

	fmt.Printf("Object %d: %s\n", objNum, obj.String())

	if stream, is := obj.(*pdfcore.PdfObjectStream); is {
		decoded, err := pdfcore.DecodeStream(stream)
		if err != nil {
			return err
		}
		fmt.Printf("Decoded:\n%s", decoded)
	} else if indObj, is := obj.(*pdfcore.PdfIndirectObject); is {
		fmt.Printf("%T\n", indObj.PdfObject)
		fmt.Printf("%s\n", indObj.PdfObject.String())
	}

	return nil
}
