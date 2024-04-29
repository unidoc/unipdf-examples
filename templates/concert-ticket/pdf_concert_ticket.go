/*
 * This example demonstrates creating a concert ticket using the creator templates.
 *
 * Run as: go run pdf_concert_ticket.go
 */

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

// Field represents a single field with name and value.
type Field struct {
	FieldName  string `json:"field_name"`
	FieldValue string `json:"field_value"`
}

// Ticket holds all data related to a ticket.
type Ticket struct {
	EventTime         string   `json:"event_time"`
	TicketNumber      string   `json:"ticket_number"`
	Detail            []Field  `json:"ticket_detail"`
	RulesOfAttendance []string `json:"rules_of_attendance"`
	RulesOfPurchase   []string `json:"rules_of_purchase"`
}

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io.
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}

	common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
}
func main() {
	c := creator.New()
	c.SetPageMargins(20, 20, 20, 20)

	// Create qrCode.
	qrCode, err := createQRCode("https://github.com/unidoc/unipdf-examples/tree/master/templates/concert-ticket/", 500, 500)
	if err != nil {
		panic(err)
	}

	// Read ticket data.
	ticket, err := readTemplateData("./concert-ticket.json")
	if err != nil {
		panic(err)
	}

	// Read template file.
	tpl, err := readTemplate("./templates/main.tpl")
	if err != nil {
		log.Fatal(err)
	}

	// Draw main content template.
	tplOpts := &creator.TemplateOptions{
		ImageMap: map[string]*model.Image{
			"qr-code": qrCode,
		},
		HelperFuncMap: template.FuncMap{
			"formatTime": func(val, format string) string {
				t, _ := time.Parse("2006-01-02T15:04:05Z", val)
				return t.Format(format)
			},
		},
	}

	if err := c.DrawTemplate(tpl, ticket, tplOpts); err != nil {
		log.Fatal(err)
	}

	// Write to output file.
	if err := c.WriteToFile("unipdf-ticket.pdf"); err != nil {
		log.Fatal(err)
	}
}

// createQRCode creates a new QR code image encoding the provided text with the specified width and height.
func createQRCode(text string, width, height int) (*model.Image, error) {
	qrCode, err := qr.Encode(text, qr.M, qr.Auto)
	if err != nil {
		return nil, err
	}

	qrCode, err = barcode.Scale(qrCode, width, height)
	if err != nil {
		return nil, err
	}

	img, err := model.ImageHandling.NewImageFromGoImage(qrCode)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// readTemplate reads template file and returns an io.Reader.
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

// readTemplateData reads the ticket data from the json file provided by `filePath`.
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
