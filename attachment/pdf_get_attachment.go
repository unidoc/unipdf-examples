/*
* Retrieve list of attachment file and save it locally.
 */

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/model"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}
}

func main() {
	inputPath := "output.pdf"

	err := listAttachments(inputPath)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	fmt.Println("Done")
}

func listAttachments(inputPath string) error {
	// Read the input pdf file.
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return err
	}

	files, err := pdfReader.GetAttachedFiles()
	if err != nil {
		return err
	}

	_, err = os.Stat("output")
	if os.IsNotExist(err) {
		err = os.Mkdir("output", 0777)
		if err != nil {
			return err
		}
	}

	for _, v := range files {
		err := os.WriteFile(filepath.Join("output", fmt.Sprintf("%s.xml", v.Name)), v.Content, 0655)
		if err != nil {
			return err
		}

		fmt.Printf("Name: %s; Hash: %s\n", v.Name, v.Hash)
	}

	return nil
}
