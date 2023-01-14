package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/unidoc/unipdf/v3/creator"
)

type Field struct {
	FieldName  string `json:"FieldName"`
	FieldValue string `json:"FieldValue"`
}
type Reciept struct {
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
	process(receipt)
}

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

func readReceipt(jsonFile string) (*Reciept, error) {
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
	receipt := Reciept{
		Fields: fields,
	}
	return &receipt, nil
}

func process(reciept *Reciept) {
	c := creator.New()
	tpl, err := readTemplate("./templates/receipt.tpl")
	if err != nil {
		log.Fatal(err)
	}

	// Draw front page template.
	if err := c.DrawTemplate(tpl, reciept, nil); err != nil {
		log.Fatal(err)
	}

	// Write output file.
	if err := c.WriteToFile("receipt.pdf"); err != nil {
		log.Fatal(err)
	}
}
