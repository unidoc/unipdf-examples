/*
 * This example showcases the usage of creator templates by creating a sample
 * medical bill.
 *
 * Run as: go run pdf_medical_bill.go
 */

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
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

	// Read medical bill JSON data.
	bill, err := readBillData("medical_bill.json")
	if err != nil {
		log.Fatal(err)
	}

	data := map[string]interface{}{
		"institution": map[string]interface{}{
			"name":     "UniDoc Medial Center",
			"address1": "123 Main Street",
			"address2": "Anywhere, NY 12345 - 6789",
		},
		"bill": bill,
	}

	if err := c.DrawTemplate(mainTpl, data, nil); err != nil {
		log.Fatal(err)
	}

	// Write output file.
	if err := c.WriteToFile("unipdf-medical-bill.pdf"); err != nil {
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

// readBillData reads the medical bill data from a specified JSON file.
func readBillData(jsonFile string) (*Bill, error) {
	file, err := os.Open(jsonFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bill := &Bill{}
	if err := json.NewDecoder(file).Decode(bill); err != nil {
		return nil, err
	}

	return bill, nil
}

// Bill holds the medical bill data.
type Bill struct {
	Guarantor     Guarantor         `json:"guarantor"`
	StatementDate string            `json:"statement_date"`
	DueDate       string            `json:"due_date"`
	Services      []*MedicalService `json:"services"`
	Total         string            `json:"total"`
}

// Guarantor holds guarantor data.
type Guarantor struct {
	Number   string `json:"number"`
	Name     string `json:"name"`
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
}

// MedicalService holds medical service list data.
type MedicalService struct {
	Items            []*ServiceItem `json:"items"`
	TotalPayments    string         `json:"total_payments"`
	TotalAdjustments string         `json:"total_adjustments"`
	PatientDue       string         `json:"patient_due"`
}

// ServiceItem holds medical service item data.
type ServiceItem struct {
	Date        string `json:"date,omitempty"`
	Description string `json:"desc"`
	Charges     string `json:"charges,omitempty"`
	Payment     string `json:"payment,omitempty"`
}
