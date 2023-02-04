package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type RentalAgreement struct {
	Date                           string   `json:"date"`
	CompanyName                    string   `json:"company_name"`
	CompanyAddress                 string   `json:"company_address"`
	Tenants                        []string `json:"tenants"`
	ApartmentAddress               string   `json:"apartment_address"`
	UnitSize                       int64    `json:"unit_size"`
	BeginningDate                  string   `json:"beginning_date"`
	EndingDate                     string   `json:"ending_date"`
	MonthlyInstallment             string   `json:"monthly_installment"`
	MinimumAbandonmentDays         int      `json:"minimum_abandonment_days"`
	InsufficientFundFee            string   `json:"ins_fee_amount"`
	LatePaymentFee                 string   `json:"late_fee_amount"`
	SecurityDeposit                string   `json:"sec_deposit_amount"`
	SecurityDepositReturnTime      int      `json:"security_deposit_return_time"`
	PurchaseDepositAmount          string   `json:"purchase_deposit_amount"`
	PurchaseAmount                 string   `json:"purchase_amount"`
	ParkingSpacesDesc              string   `json:"parking_spaces_description"`
	PetFee                         string   `json:"pet_fee_amount"`
	TerminationFee                 string   `json:"termination_fee"`
	TerminationNoticePeriod        int      `json:"early_termination_notice_period"`
	NumberOfBedrooms               int      `json:"number_of_bedrooms"`
	NumberOfBathRooms              float64  `json:"number_of_bathrooms"`
	NumberOfAllowedPets            int      `json:"number_of_allowed_pets"`
	NumberOfParkingSpaces          int      `json:"number_of_parking_spaces"`
	CancellationNotificationPeriod int      `json:"cancellation_notification_period"`
	ContinuationNotificationPeriod int      `json:"continuation_notification_period"`
	MoveInCheckList                struct {
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
	c.SetPageMargins(90, 60, 90, 120)
	// Read main content template.
	mainTpl, err := readTemplate("templates/main.tpl")
	if err != nil {
		log.Fatal(err)
	}

	arialBold, err := model.NewPdfFontFromTTFFile("./templates/res/arialbd.ttf")
	if err != nil {
		log.Fatal(err)
	}
	tplOpts := &creator.TemplateOptions{
		FontMap: map[string]*model.PdfFont{
			"arial-bold": arialBold,
		},
		HelperFuncMap: template.FuncMap{
			"formatTime": func(val, format string) string {
				t, _ := time.Parse("2006-01-02T00:00:00", val)
				return t.Format(format)
			},
			"listNames": func(tenants []string) string {
				nameList := ""
				for i, t := range tenants {
					if i < (len(tenants) - 2) {
						nameList = nameList + t + ", "
					} else if i == (len(tenants) - 2) {
						nameList = nameList + t + " and "
					} else {
						nameList = nameList + t
					}
				}
				return nameList
			},
			"computeMargin": func(text string) float64 {
				// An arbitrary margin calculation to position the lines.
				// Basically the number of the characters multiplied by 5 happens to position the line right next to the text.
				// TODO maybe calculating the width of the text based on the font would make this accurate.
				return float64(len(text)+10) * 5
			},
			"numberToWord": func(number int, capitalize bool) string {
				w := ""
				if number < 20 {
					w = NumberToWord[number]
					if capitalize {
						w = cases.Title(language.English).String(w)
					}
					return w
				}

				r := number % 10
				if r == 0 {
					w = NumberToWord[number]
				} else {
					w = NumberToWord[number-r] + " " + NumberToWord[r]
				}
				if capitalize {
					w = cases.Title(language.English).String(w)
				}
				return w
			},
		},
	}

	// Read data from json.
	rentalAgreement, err := readRentalAgreement("rental_data.json")
	if err != nil {
		log.Fatal(err)
	}
	// Draw the main template.
	if err := c.DrawTemplate(mainTpl, rentalAgreement, tplOpts); err != nil {
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
			"Agreement":  rentalAgreement,
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

var NumberToWord = map[int]string{
	1:  "one",
	2:  "two",
	3:  "three",
	4:  "four",
	5:  "five",
	6:  "six",
	7:  "seven",
	8:  "eight",
	9:  "nine",
	10: "ten",
	11: "eleven",
	12: "twelve",
	13: "thirteen",
	14: "fourteen",
	15: "fifteen",
	16: "sixteen",
	17: "seventeen",
	18: "eighteen",
	19: "nineteen",
	20: "twenty",
	30: "thirty",
	40: "forty",
	50: "fifty",
	60: "sixty",
	70: "seventy",
	80: "eighty",
	90: "ninety",
}
