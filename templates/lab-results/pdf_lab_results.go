/*
 * This example showcases the usage of creator templates by creating a sample
 * lab results document.
 *
 * Run as: go run pdf_lab_results.go
 */

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"text/template"

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
	c.SetPageMargins(15, 15, 180, 100)

	// Read main content template.
	mainTpl, err := readTemplate("templates/main.tpl")
	if err != nil {
		log.Fatal(err)
	}

	// Read lab results JSON data.
	data, err := readResultData("lab_results.json")
	if err != nil {
		log.Fatal(err)
	}

	// Draw main content template.
	tplOpts := &creator.TemplateOptions{
		HelperFuncMap: template.FuncMap{
			"infoColor": func(level Level) string {
				if color, ok := levelColor[level]; ok {
					return color
				}

				return "#000000"
			},
		},
	}

	if err := c.DrawTemplate(mainTpl, data, tplOpts); err != nil {
		log.Fatal(err)
	}

	// Draw header and footer.
	drawHeader := func(tplPath string, block *creator.Block, labResults *LabResults, pageNum, totalPages int) {
		// Read template.
		tpl, err := readTemplate(tplPath)
		if err != nil {
			log.Fatal(err)
		}

		// Draw template.
		data := map[string]interface{}{
			"Date":         "02/03/2023 0945",
			"SpecimenID":   labResults.SpecimenID,
			"ControlID":    labResults.ControlID,
			"AcctNum":      labResults.AcctNum,
			"Phone":        labResults.Phone,
			"Rte":          labResults.Rte,
			"SupportPhone": labResults.SupportPhone,
			"PageNum":      pageNum,
			"TotalPages":   totalPages,
		}

		if err := block.DrawTemplate(c, tpl, data, tplOpts); err != nil {
			log.Fatal(err)
		}
	}

	c.DrawHeader(func(block *creator.Block, args creator.HeaderFunctionArgs) {
		drawHeader("templates/header.tpl", block, data, args.PageNum, args.TotalPages)
	})
	c.DrawFooter(func(block *creator.Block, args creator.FooterFunctionArgs) {
		drawHeader("templates/footer.tpl", block, data, args.PageNum, args.TotalPages)
	})

	// Write output file.
	if err := c.WriteToFile("unipdf-lab-results.pdf"); err != nil {
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

// readResultData reads the lab results data from a specified JSON file.
func readResultData(jsonFile string) (*LabResults, error) {
	file, err := os.Open(jsonFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := &LabResults{}
	if err := json.NewDecoder(file).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

// LabResults hold lab results data.
type LabResults struct {
	SpecimenID   string        `json:"specimen_id"`
	ControlID    string        `json:"control_id"`
	AcctNum      string        `json:"account_number"`
	Phone        string        `json:"phone"`
	SupportPhone string        `json:"support_phone"`
	Rte          string        `json:"rte"`
	Patient      Patient       `json:"patient"`
	Specimen     Specimen      `json:"specimen"`
	Physician    Physician     `json:"physician"`
	Results      []*TestResult `json:"tests"`
}

// Patient holds patient related data.
type Patient struct {
	Birthdate string `json:"birthdate"`
	Age       string `json:"age"`
	Gender    string `json:"gender"`
	Id        string `json:"id"`
}

// Specimen holds test specimen data.
type Specimen struct {
	Collected string `json:"collected"`
	Received  string `json:"received"`
	Entered   string `json:"entered"`
	Reported  string `json:"reported"`
}

// Physician holds physician data.
type Physician struct {
	Ordering  string `json:"ordering"`
	Referring string `json:"referring"`
	Id        string `json:"id"`
	Npi       string `json:"npi"`
}

// Level represents test results info level.
type Level string

// Acceptable level values.
const (
	Green  Level = "green"
	Yellow       = "yellow"
	Red          = "red"
)

// Map of result level and text color.
var levelColor = map[Level]string{
	Green:  "#407505",
	Yellow: "#F5A623",
	Red:    "#FF0000",
}

// TestResult represents each test result.
type TestResult struct {
	Info        string   `json:"info"`
	Level       Level    `json:"level"`
	Test        string   `json:"test"`
	Result      string   `json:"result"`
	Flag        string   `json:"flag"`
	Units       string   `json:"units"`
	RefInterval string   `json:"ref_interval"`
	Lab         string   `json:"lab"`
	Description []string `json:"description"`
}
