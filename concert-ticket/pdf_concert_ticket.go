package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"text/template"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
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
	qrCode, err := createQRCode("https://github.com/unidoc/unipdf-examples/tree/master/concert-ticket/", 500, 500)
	if err != nil {
		panic(err)
	}
	// Read ticket data.
	ticket, err := readTemplateData("./templates/concert-ticket.json")
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
			"extendDict": func(m map[string]interface{}, params ...interface{}) (map[string]interface{}, error) {
				lenParams := len(params)
				if lenParams%2 != 0 {
					return nil, core.ErrRangeError
				}

				out := make(map[string]interface{}, len(m))
				for key, val := range m {
					out[key] = val
				}

				for i := 0; i < lenParams; i += 2 {
					key, ok := params[i].(string)
					if !ok {
						return nil, core.ErrTypeError
					}

					out[key] = params[i+1]
				}

				return out, nil
			},
		},
	}
	// Draw main content template.
	if err := c.DrawTemplate(tpl, ticket, tplOpts); err != nil {
		log.Fatal(err)
	}

	// Write to output file.
	if err := c.WriteToFile("unipdf_ticket.pdf"); err != nil {
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
