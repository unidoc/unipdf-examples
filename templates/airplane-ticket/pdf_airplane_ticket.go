/*
 * This example showcases the usage of creator templates by creating a sample
 * airplane ticket.
 *
 * Run as: go run pdf_airplane_ticket.go
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
	c.SetPageMargins(50, 50, 25, 25)

	// Read main content template.
	mainTpl, err := readTemplate("templates/main.tpl")
	if err != nil {
		log.Fatal(err)
	}

	// Read ticket.
	ticket, err := readTicket("ticket.json")
	if err != nil {
		log.Fatal(err)
	}

	// Create QR code image.
	qrCode, err := createQRCode("https://github.com/unidoc/unipdf-examples/tree/master/templates/airplane-ticket", 500, 500)
	if err != nil {
		log.Fatal(err)
	}

	// Draw main content template.
	tplOpts := &creator.TemplateOptions{
		ImageMap: map[string]*model.Image{
			"qr-code-1": qrCode,
		},
		HelperFuncMap: template.FuncMap{
			"formatTime": func(val, format string) string {
				t, _ := time.Parse("2006-01-02T15:04:05", val)
				return t.Format(format)
			},
		},
	}

	data := map[string]interface{}{
		"company": "UniPDF Airlines",
		"ticket":  ticket,
	}

	if err := c.DrawTemplate(mainTpl, data, tplOpts); err != nil {
		log.Fatal(err)
	}

	// Write output file.
	if err := c.WriteToFile("unipdf-airplain-ticket.pdf"); err != nil {
		log.Fatal(err)
	}
}

// readTemplate reads the template at the specified file path and returns it
// as an io.Reader.
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

// Ticket represents a sample plane ticket.
type Ticket struct {
	Status    string `json:"status"`
	Passenger string `json:"passenger"`
	Document  string `json:"document"`
	Number    string `json:"number"`
	Order     string `json:"order"`
	Issued    string `json:"issued"`
	Routes    []struct {
		Flight           string `json:"flight"`
		FlightCompany    string `json:"flightCompany"`
		FlightPlaner     string `json:"flightPlaner"`
		Departure        string `json:"departure"`
		DepartureAirport string `json:"departureAirport"`
		Arrival          string `json:"arrival"`
		ArrivalAirport   string `json:"arrivalAirport"`
		Class            string `json:"class"`
		ClassAdd         string `json:"classAdd"`
		Baggage          string `json:"baggage"`
		BaggageAdd       string `json:"baggageAdd"`
		CheckIn          string `json:"checkIn"`
		CheckInAirport   string `json:"checkInAirport"`
	} `json:"routes"`
	Fares []struct {
		Name   string  `json:"name"`
		Charge float64 `json:"charge"`
	} `json:"fares"`
	PhoneNumbers []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"phoneNumbers"`
}

// readTicket reads the data for a plane ticket from a specified JSON file.
func readTicket(jsonFile string) (*Ticket, error) {
	file, err := os.Open(jsonFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ticket := &Ticket{}
	if err := json.NewDecoder(file).Decode(ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}

// createQRCode creates a new QR code image encoding the provided text, having
// the specified width and height.
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
