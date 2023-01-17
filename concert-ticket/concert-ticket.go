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
	FieldName  string `json:"field_name"`
	FieldValue string `json:"field_value"`
}
type Ticket struct {
	Detail            []Field  `json:"ticket_detail"`
	RulesOfAttendance []string `json:"rules_of_attendance"`
	RulesOfPurchase   []string `json:"rules_of_purchase"`
}

func main() {

	ticket, err := readTemplateData("./templates/concert-ticket.json")
	if err != nil {
		panic(err)
	}
	process(ticket)
	// fmt.Printf("type %T  and value %v\n", ticket.RulesOfAttendance, ticket.RulesOfAttendance)
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
func process(ticket *Ticket) {
	c := creator.New()
	c.SetPageMargins(20, 20, 20, 20)
	tpl, err := readTemplate("./templates/ticket.tpl")
	if err != nil {
		log.Fatal(err)
	}
	// opts := creator.TemplateOptions{}
	// Draw front page template.
	if err := c.DrawTemplate(tpl, ticket, nil); err != nil {
		log.Fatal(err)
	}

	// Write output file.
	if err := c.WriteToFile("ticket.pdf"); err != nil {
		log.Fatal(err)
	}
}

func readTemplateData(filePath string) (*Ticket, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var ticket Ticket
	decoder := json.NewDecoder(file)

	err = decoder.Decode(&ticket)
	if err != nil {
		return nil, err
	}

	return &ticket, nil
}
