package main

import (
	"bytes"
	"io"
	"log"
	"os"

	"github.com/unidoc/unipdf/v3/creator"
)

func main() {
	process()
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
func process() {
	c := creator.New()
	tpl, err := readTemplate("./templates/ticket.tpl")
	if err != nil {
		log.Fatal(err)
	}
	// opts := creator.TemplateOptions{}
	// Draw front page template.
	if err := c.DrawTemplate(tpl, nil, nil); err != nil {
		log.Fatal(err)
	}

	// Write output file.
	if err := c.WriteToFile("ticket.pdf"); err != nil {
		log.Fatal(err)
	}
}
