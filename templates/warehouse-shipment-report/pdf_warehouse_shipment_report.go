/*
 * This example showcases the usage of creator templates by creating a sample
 * warehouse shipment report.
 *
 * Run as: go run pdf_warehouse_shipment_report.go
 */

package main

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"io"
	"log"
	"os"
	"text/template"

	"github.com/boombuler/barcode/code128"
	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
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
	c.SetPageMargins(25, 25, 150, 25)

	// Read main content template.
	mainTpl, err := readTemplate("templates/main.tpl")
	if err != nil {
		log.Fatal(err)
	}

	// Read shipment report JSON data.
	report, err := readReportData("shipment_report.json")
	if err != nil {
		log.Fatal(err)
	}

	// Draw main content template.
	imageMap := map[string]*model.Image{}
	tplOpts := &creator.TemplateOptions{
		ImageMap: imageMap,
		HelperFuncMap: template.FuncMap{
			"createBarcode": func(text string) (string, error) {
				barcode, err := code128.Encode(text)
				if err != nil {
					return "", err
				}

				tmpImg := image.NewRGBA(image.Rect(0, 0, 100, 30))

				draw.NearestNeighbor.Scale(tmpImg, image.Rect(0, 0, 100, 20), image.Image(barcode), barcode.Bounds(), draw.Over, nil)

				col := color.RGBA{0, 0, 0, 255}
				point := fixed.Point26_6{X: fixed.I(10), Y: fixed.I(30)}

				d := &font.Drawer{
					Dst:  tmpImg,
					Src:  image.NewUniform(col),
					Face: basicfont.Face7x13,
					Dot:  point,
				}
				d.DrawString(text)

				img, err := model.ImageHandling.NewImageFromGoImage(tmpImg)
				if err != nil {
					return "", err
				}

				imageMap[text] = img

				return text, nil
			},
		},
	}

	if err := c.DrawTemplate(mainTpl, report.Shipments, tplOpts); err != nil {
		log.Fatal(err)
	}

	// Draw header.
	drawHeader := func(tplPath string, block *creator.Block, company *Company) {
		// Read template.
		tpl, err := readTemplate(tplPath)
		if err != nil {
			log.Fatal(err)
		}

		// Draw template.
		if err := block.DrawTemplate(c, tpl, company, tplOpts); err != nil {
			log.Fatal(err)
		}
	}

	c.DrawHeader(func(block *creator.Block, args creator.HeaderFunctionArgs) {
		drawHeader("templates/header.tpl", block, report.Company)
	})

	// Write output file.
	if err := c.WriteToFile("unipdf-warehouse-shipment-report.pdf"); err != nil {
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

// readReportData reads the shipment report data from a specified JSON file.
func readReportData(jsonFile string) (*Report, error) {
	file, err := os.Open(jsonFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	report := &Report{}
	if err := json.NewDecoder(file).Decode(report); err != nil {
		return nil, err
	}

	return report, nil
}

// Report represents the report data.
type Report struct {
	Company   *Company    `json:"company"`
	Shipments []*Shipment `json:"shipments"`
}

// Company represents company data.
type Company struct {
	Brand   string `json:"brand"`
	Name    string `json:"name"`
	POBox   string `json:"pobox"`
	Address string `json:"address"`
	Area    string `json:"area"`
	City    string `json:"city"`
	Country string `json:"country"`
}

// Shipment represent the shipment data.
type Shipment struct {
	Time   string   `json:"time"`
	Orders []*Order `json:"orders"`
}

// Order represents the order data.
type Order struct {
	Time     string     `json:"time"`
	PIC      string     `json:"pic"`
	Products []*Product `json:"products"`
}

// Product represent the product data.
type Product struct {
	Barcode string `json:"barcode"`
	Code    string `json:"code"`
	Name    string `json:"name"`
}
