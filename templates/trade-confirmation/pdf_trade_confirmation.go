/*
 * This example showcases the usage of creator templates by creating a sample
 * trade confirmation document.
 *
 * Run as: go run pdf_trade_confirmation.go
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
	c.SetPageMargins(50, 50, 25, 25)

	// Read main content template.
	mainTpl, err := readTemplate("templates/main.tpl")
	if err != nil {
		log.Fatal(err)
	}

	// Read trading data json.
	trade, err := readTradeData("trade.json")
	if err != nil {
		log.Fatal(err)
	}

	// Draw main content template.

	data := map[string]interface{}{
		"firmName":    "UniDoc Financial Firm",
		"firmAddress": "123 Main Street\nPortland, ME 12345\n(123) 456-789",
		"trade":       trade,
	}

	if err := c.DrawTemplate(mainTpl, data, nil); err != nil {
		log.Fatal(err)
	}

	// Write output file.
	if err := c.WriteToFile("unipdf-trade-confirmation.pdf"); err != nil {
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

// Trade represents the sample trading data.
type Trade struct {
	Date          string           `json:"date"`
	AccountNumber string           `json:"accountNumber"`
	Name          string           `json:"userName"`
	Address       string           `json:"userAddress"`
	Action        string           `json:"actionName"`
	ProductDesc   string           `json:"productDesc"`
	BoughtUnit    string           `json:"boughtUnit"`
	BoughtPrice   string           `json:"boughtPrice"`
	OrderNumber   string           `json:"orderNumber"`
	Calculation   TradeCalculation `json:"data"`
	MarkupValue   string           `json:"markupValue"`
	InfoUrl       string           `json:"infoUrl"`
}

// TradeCalculation represents trade calculation data.
type TradeCalculation struct {
	PrincipalAmount string `json:"principalAmount"`
	AccruedInterest string `json:"accruedInterest"`
	TransactionFee  string `json:"transactionFee"`
	Total           string `json:"total"`
	BankQualified   string `json:"bankQualified"`
	State           string `json:"state"`
	DatedToDate     string `json:"datedToDate"`
	YieldToMaturity string `json:"yieldToMaturity"`
	YieldToCall     string `json:"yieldToCall"`
	Callable        string `json:"callable"`
	TaxExcempt      string `json:"taxExcempt"`
	Capacity        string `json:"capacity"`
	BondForm        string `json:"bondForm"`
	TradeDate       string `json:"tradeDate"`
	TradeTime       string `json:"tradeTime"`
	SettlementDate  string `json:"settlementDate"`
}

// readTradeData reads the trading data from a specified JSON file.
func readTradeData(jsonFile string) (*Trade, error) {
	file, err := os.Open(jsonFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	trade := &Trade{}
	if err := json.NewDecoder(file).Decode(trade); err != nil {
		return nil, err
	}

	return trade, nil
}
