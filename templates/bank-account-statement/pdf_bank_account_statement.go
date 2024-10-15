/*
 * This example showcases the usage of creator templates by creating a sample
 * bank account statement.
 *
 * Run as: go run pdf_bank_account_statement.go
 */

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"
	"text/template"
	"time"

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
	c.SetPageMargins(50, 50, 80, 25)

	// Read main content template.
	mainTpl, err := readTemplate("templates/main.tpl")
	if err != nil {
		log.Fatal(err)
	}

	// Read account statement data.
	statement, err := readAccountStatement("account_statement.json")
	if err != nil {
		log.Fatal(err)
	}

	// Draw main content template.
	tplOpts := &creator.TemplateOptions{
		HelperFuncMap: template.FuncMap{
			"strRepeat": strings.Repeat,
			"loop": func(size uint64) []struct{} {
				return make([]struct{}, size)
			},
			"formatTime": func(val, format string) string {
				t, _ := time.Parse("2006-01-02T15:04:05", val)
				return t.Format(format)
			},
		},
	}

	if err := c.DrawTemplate(mainTpl, statement, tplOpts); err != nil {
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
			"Statement":  statement,
			"PageNum":    pageNum,
			"TotalPages": totalPages,
		}

		if err := block.DrawTemplate(c, tpl, data, tplOpts); err != nil {
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
	if err := c.WriteToFile("unipdf-bank-account-statement.pdf"); err != nil {
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

// AccountStatement represents a sample account statement.
type AccountStatement struct {
	BankName                 string  `json:"bankName"`
	BankNameState            string  `json:"bankNameState"`
	AccountNumber            string  `json:"accountNumber"`
	DateBegin                string  `json:"dateBegin"`
	DateEnd                  string  `json:"dateEnd"`
	CompanyName              string  `json:"companyName"`
	CompanyAddress           string  `json:"companyAddress"`
	ReportAddress            string  `json:"reportAddress"`
	PhoneFree                string  `json:"phoneFree"`
	Phone                    string  `json:"phone"`
	Tty                      string  `json:"tTY"`
	Online                   string  `json:"online"`
	White                    string  `json:"white"`
	BusinessPlanURL          string  `json:"businessPlanURL"`
	AccountOptionsURL        string  `json:"accountOptionsURL"`
	Advt                     string  `json:"advt"`
	BeginningBalance         float64 `json:"beginningBalance"`
	Withdrawals              float64 `json:"withdrawals"`
	Deposits                 float64 `json:"deposits"`
	EndingBalance            float64 `json:"endingBalance"`
	AverageBalance           float64 `json:"averageBalance"`
	DepositRTN               string  `json:"depositRTN"`
	WireRTN                  string  `json:"wireRTN"`
	StandardServiceFee       float64 `json:"standardServiceFee"`
	MinimumRequired          float64 `json:"minimumRequired"`
	ServiceFee               float64 `json:"serviceFee"`
	ServiceDiscount          float64 `json:"serviceDiscount"`
	TransactionUnits         float64 `json:"transactionUnits"`
	TransactionUnitsIncluded float64 `json:"transactionUnitsIncluded"`
	TransactionExcessUnits   float64 `json:"transactionExcessUnits"`
	ServiceCharge            float64 `json:"serviceCharge"`
	TotalServiceCharge       float64 `json:"totalServiceCharge"`
	FeedbackPhone            string  `json:"feedbackPhone"`
	TransactionDeposits      float64 `json:"transactionDeposits"`
	TransactionWithdrawals   float64 `json:"transactionWithdrawals"`
	Transactions             []struct {
		Date               string      `json:"date"`
		Check              interface{} `json:"check"`
		Details            string      `json:"details"`
		Deposits           float64     `json:"deposits"`
		Withdrawals        float64     `json:"withdrawals"`
		EndingDailyBalance float64     `json:"endingDailyBalance"`
	} `json:"Transactions"`
	AccountOptionsLabels []string `json:"accountOptionsLabels"`
}

// readAccountStatement reads the data for an account statement from a
// specified JSON file.
func readAccountStatement(jsonFile string) (*AccountStatement, error) {
	file, err := os.Open(jsonFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	statement := &AccountStatement{}
	if err := json.NewDecoder(file).Decode(statement); err != nil {
		return nil, err
	}

	return statement, nil
}
