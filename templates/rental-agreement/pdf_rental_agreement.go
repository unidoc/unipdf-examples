package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
)

type RentalAgreement struct {
	Date                  string   `json:"date"`
	CompanyName           string   `json:"company_name"`
	CompanyAddress        string   `json:"company_address"`
	Tenants               []string `json:"tenants"`
	ApartmentAddress      string   `json:"apartment_address"`
	BeginningDate         string   `json:"beginning_date"`
	EndingDate            string   `json:"ending_date"`
	MonthlyInstallment    string   `json:"monthly_installment"`
	InsufficientFundFee   string   `json:"ins_fee_amount"`
	LatePaymentFee        string   `json:"50.00"`
	SecurityDeposit       string   `json:"sec_deposit_amount"`
	PurchaseDepositAmount string   `json:"purchase_deposit_amount"`
	PetFee                string   `json:"pet_fee_amount"`
	MoveInCheckList       struct {
		LivingRoom  []string `json:"living_room"`
		DinningRoom []string `json:"dinning_room"`
		Kitchen     []string `json:"kitchen"`
		Bathroom    []string `json:"bathroom"`
		Other       []string `json:"other"`
	} `json:"move_in_check_list"`
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
	c.SetPageMargins(70, 50, 100, 100)
	// Read main content template.
	mainTpl, err := readTemplate("templates/main.tpl")
	if err != nil {
		log.Fatal(err)
	}
	// Read data from json.
	rentalAgreement, err := readRentalAgreement("rental_data.json")
	if err != nil {
		log.Fatal(err)
	}
	// Draw the main template.
	if err := c.DrawTemplate(mainTpl, rentalAgreement, nil); err != nil {
		log.Fatal(err)
	}
	// Draw header and footer.
	drawHeader := func(tplPath string, block *creator.Block, pageNum, totalPages int) {
		// Read template.
		tpl, err := readTemplate(tplPath)
		if err != nil {
			log.Fatal(err)
		}
		// Draw template.
		data := map[string]interface{}{
			"Date":       time.Now(),
			"Statement":  rentalAgreement,
			"PageNum":    pageNum,
			"TotalPages": totalPages,
		}
		if err := block.DrawTemplate(c, tpl, data, nil); err != nil {
			log.Fatal(err)
		}
	}

	c.DrawHeader(func(block *creator.Block, args creator.HeaderFunctionArgs) {
		drawHeader("templates/header.tpl", block, args.PageNum, args.TotalPages)
	})
	c.DrawFooter(func(block *creator.Block, args creator.FooterFunctionArgs) {
		drawHeader("templates/footer.tpl", block, args.PageNum, args.TotalPages)
	})
	// Write output file.
	if err := c.WriteToFile("unipdf-rental-agreement.pdf"); err != nil {
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

// readRentalAgreement reads the data for an rental agreement document from the
// specified JSON file.
func readRentalAgreement(jsonFile string) (*RentalAgreement, error) {
	file, err := os.Open(jsonFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	rentalAgreement := &RentalAgreement{}
	if err := json.NewDecoder(file).Decode(rentalAgreement); err != nil {
		return nil, err
	}

	return rentalAgreement, nil
}
