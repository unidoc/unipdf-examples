/*
 * This example showcases the usage of creator templates by creating a sample
 * boarding pass document.
 *
 * Run as: go run pdf_boarding_pass.go
 */

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"

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
	c.SetPageMargins(30, 30, 20, 20)

	// Read main content template.
	mainTpl, err := readTemplate("templates/main.tpl")
	if err != nil {
		log.Fatal(err)
	}

	// Read boarding pass JSON data.
	pass, err := readBoardingPassData("boarding_pass.json")
	if err != nil {
		log.Fatal(err)
	}

	// Create QR code image.
	qrCode, err := createQRCode("https://github.com/unidoc/unipdf-examples/tree/master/templates/boarding-pass", 100, 100)
	if err != nil {
		log.Fatal(err)
	}

	// Draw main content template.
	tplOpts := &creator.TemplateOptions{
		ImageMap: map[string]*model.Image{
			"qr-code-1": qrCode,
		},
	}

	if err := c.DrawTemplate(mainTpl, pass, tplOpts); err != nil {
		log.Fatal(err)
	}

	// Write output file.
	if err := c.WriteToFile("unipdf-boarding-pass.pdf"); err != nil {
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

// readBoardingPassData reads the boarding pass data from a specified JSON file.
func readBoardingPassData(jsonFile string) (*Pass, error) {
	file, err := os.Open(jsonFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	pass := &Pass{}
	if err := json.NewDecoder(file).Decode(pass); err != nil {
		return nil, err
	}

	return pass, nil
}

// Pass holds the boarding pass data.
type Pass struct {
	Etk           string   `json:"etk"`
	RegNumber     string   `json:"reg_num"`
	Name          string   `json:"passenger_name"`
	From          *Airport `json:"from"`
	Destination   *Airport `json:"destination"`
	FlightNumber  string   `json:"flight_number"`
	Gate          string   `json:"gate"`
	Class         string   `json:"class"`
	Seat          string   `json:"seat"`
	Date          string   `json:"date"`
	BoardingTime  string   `json:"boarding_time"`
	DepartureTime string   `json:"departure_time"`
	ArrivalTime   string   `json:"arrival_time"`
}

// Airport holds the airport information data.
type Airport struct {
	City string `json:"city"`
	Code string `json:"code"`
}
