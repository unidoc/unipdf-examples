/*
 * This example showcases the usage of creator templates to create a receipt document
 *
 * Run as: go run pdf_receipt.go
 */
package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/unidoc/unipdf/v3/creator"
)

// func init() {
// 	// Make sure to load your metered License API key prior to using the library.
// 	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io.
// 	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
// 	if err != nil {
// 		panic(err)
// 	}

// 	common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
// }

// Field represents field item.
type Field struct {
	FieldName  string `json:"FieldName"`
	FieldValue string `json:"FieldValue"`
}

// Receipt represent Receipt object.
type Receipt struct {
	Title  string
	Fields []Field
}

func main() {
	filePath := "./contents/receipt.json"
	receipt, err := readReceipt(filePath)
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}
	render(receipt)
}

// readTemplate reads template file.
func readTemplate(tplFile string) (io.Reader, error) {
	file, err := os.Open(tplFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, file); err != nil {
		return nil, err
	}

	return buf, nil
}

// readReceipt reads the receipt json file and decods it to `Receipt` object.
func readReceipt(jsonFile string) (*Receipt, error) {
	file, err := os.Open(jsonFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var fields []Field
	decoder := json.NewDecoder(file)

	err = decoder.Decode(&fields)
	if err != nil {
		return nil, err
	}
	receipt := Receipt{
		Title:  "Receipt",
		Fields: fields,
	}
	return &receipt, nil
}

// render reads template file and draws the template content to output file.
func render(receipt *Receipt) {
	c := creator.New()
	c.SetPageMargins(15, 15, 20, 20)
	c.SetPageSize(creator.PageSizeA5)
	tpl, err := readTemplate("./templates/receipt.tpl")
	if err != nil {
		log.Fatal(err)
	}
	// Draw front page template.
	if err := c.DrawTemplate(tpl, receipt, nil); err != nil {
		log.Fatal(err)
	}

	// Write output file.
	if err := c.WriteToFile("unipdf-receipt.pdf"); err != nil {
		log.Fatal(err)
	}
}