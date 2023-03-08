/*
 * This example showcases the usage of creator templates by creating a sample
 * log book report.
 *
 * Run as: go pdf_log_book.go
 */
package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"os"
	"strings"

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

type Item struct {
	Source        string `json:"Source"`
	Manufacturer  string `json:"Manufacturer"`
	Model         string `json:"Model"`
	VIN           string `json"VIN"`
	Received      string `json:"Received"`
	Sent          string `json:"Sent"`
	Buyer_Name    string `json:"Buyer_Name"`
	Buyer_Address string `json:"Buyer_Address"`
	Buyer_State   string `json:"Buyer_State"`
	Buyer_Zip     string `json:"Buyer_Zip"`
	Discarded     string `json:"Discarded"`
	DiscardReason string `json:"DiscardReason"`
}

func main() {
	c := creator.New()
	size := creator.PageSize{842, 595}
	c.SetPageSize(size)
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

	// Read data from JSON.
	items, err := readData("contents/operations_log.json")
	if err != nil {
		log.Fatal(err)
	}
	pageContent := splitData(items)
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
			"add": func(num1, num2 int) int {
				return num1 + num2
			},
			"getSlice": func(s string) []string {
				return strings.Split(s, ",")
			},
			"htmlescaper": func(value string) string {
				return template.HTMLEscaper(value)
			},
		},
	}

	if err = c.DrawTemplate(mainTpl, pageContent, tplOpts); err != nil {
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
			"DateOfPrint": "12/28/2020",
			"DateRange":   "06/10/2018 - 06/09/2019",
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

// readData reads data from the json file and decodes it to `MedicalData` object.
func readData(jsonFile string) ([]Item, error) {
	file, err := os.Open(jsonFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []Item
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// splitData splits `items` data to per page content.
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
