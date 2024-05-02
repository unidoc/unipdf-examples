/*
 * This example showcases the usage of creator templates by creating a sample
 * log book report.
 *
 * Run as: go run pdf_log_book.go
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

// Item represents a log book item.
type Item struct {
	Source        string `json:"Source"`
	Manufacturer  string `json:"Manufacturer"`
	Model         string `json:"Model"`
	VIN           string `json:"VIN"`
	Received      string `json:"Received"`
	Sent          string `json:"Sent"`
	BuyerName     string `json:"Buyer_Name"`
	BuyerAddress  string `json:"Buyer_Address"`
	BuyerState    string `json:"Buyer_State"`
	BuyerZip      string `json:"Buyer_Zip"`
	Discarded     string `json:"Discarded"`
	DiscardReason string `json:"DiscardReason"`
}

// LogBookData represents data used for log book document.
type LogBookData struct {
	Items       []Item `json:"Items"`
	DateOfPrint string `json:"DateOfPrint"`
	DateRange   string `json:"DateRange"`
}

func main() {
	c := creator.New()
	c.SetPageSize(creator.PageSize{842, 595})
	c.SetPageMargins(10, 10, 65, 55)

	// Read main content template.
	mainTpl, err := readTemplate("templates/main.tpl")
	if err != nil {
		log.Fatal(err)
	}

	exoBold, err := model.NewPdfFontFromTTFFile("./templates/res/Exo-Bold.ttf")
	if err != nil {
		log.Fatal(err)
	}

	exoItalic, err := model.NewPdfFontFromTTFFile("./templates/res/Exo-Italic.ttf")
	if err != nil {
		log.Fatal(err)
	}

	exoRegular, err := model.NewPdfFontFromTTFFile("./templates/res/Exo-Regular.ttf")
	if err != nil {
		log.Fatal(err)
	}

	// Read data from JSON file.
	data, err := readData("contents/operations_log.json")
	if err != nil {
		log.Fatal(err)
	}
	pageContent := splitData(data.Items)

	// Create template options.
	tplOpts := &creator.TemplateOptions{
		FontMap: map[string]*model.PdfFont{
			"exo-bold":    exoBold,
			"exo-italic":  exoItalic,
			"exo-regular": exoRegular,
		},
		HelperFuncMap: template.FuncMap{
			"isEven": func(num int) bool {
				return num%2 == 0
			},
			"getSlice": func(s string) []string {
				return strings.Split(s, ",")
			},
			"htmlescaper": func(value string) string {
				return template.HTMLEscaper(value)
			},
		},
	}

	// Draw template.
	logBookData := map[string]interface{}{
		"DateOfPrint":  data.DateOfPrint,
		"DateRange":    data.DateRange,
		"NumOfRecords": len(data.Items),
		"PageToItems":  pageContent,
	}

	if err = c.DrawTemplate(mainTpl, logBookData, tplOpts); err != nil {
		log.Fatal(err)
	}

	// Draw front page.
	c.CreateFrontPage(func(args creator.FrontpageFunctionArgs) {
		// Read front page template.
		frontPageTpl, err := readTemplate("templates/front-page.tpl")
		if err != nil {
			log.Fatal(err)
		}

		// Draw front page template.
		if err := c.DrawTemplate(frontPageTpl, c, tplOpts); err != nil {
			log.Fatal(err)
		}
	})

	draw := func(tplPath string, block *creator.Block, pageNum int) {
		// Read template.
		tpl, err := readTemplate(tplPath)
		if err != nil {
			log.Fatal(err)
		}

		// Draw template.
		data := map[string]interface{}{
			"PageNum":     pageNum,
			"DateOfPrint": data.DateOfPrint,
			"DateRange":   data.DateRange,
		}
		if err := block.DrawTemplate(c, tpl, data, tplOpts); err != nil {
			log.Fatal(err)
		}
	}

	c.DrawHeader(func(block *creator.Block, args creator.HeaderFunctionArgs) {
		draw("templates/header.tpl", block, args.PageNum)
	})

	c.DrawFooter(func(block *creator.Block, args creator.FooterFunctionArgs) {
		draw("templates/footer.tpl", block, args.PageNum)
	})

	// Write output file.
	if err := c.WriteToFile("unipdf-log-book.pdf"); err != nil {
		log.Fatal(err)
	}
}

// readTemplate reads the template at the specified file path and returns an io.Reader.
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

// readData reads data from the json file and decodes it to LogBookData.
func readData(jsonFile string) (*LogBookData, error) {
	file, err := os.Open(jsonFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data LogBookData
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// splitData splits `items` data and returns a map of page to items.
func splitData(items []Item) map[int][]Item {
	start := 0
	end := 0
	page := 2
	size := 0
	pageContent := map[int][]Item{}
	for end < len(items) {
		if page == 2 {
			size = 21
		} else if page == 4 {
			size = 27
		} else {
			size = 28
		}
		end = end + size
		current := items[start:end]
		start = end
		pageContent[page] = current
		page += 2
	}

	return pageContent
}
