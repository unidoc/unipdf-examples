/*
 * This example showcases the usage of creator templates by creating a sample medication schedule report.
 *
 * Run as: go run pdf-medication-schedule.go
 */

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/ean"
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

type MedicalData struct {
	Patient struct {
		Name                 string `json:"name"`
		SocialSecurityNumber string `json:"social_security_number"`
		Dob                  string `json:"dob"`
	} `json:"patient"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
	EmergencyLine   string `json:"emergency_line"`
	InformationLine string `json:"information_line"`
	Website         string `json:"website"`
}

func main() {
	c := creator.New()
	size := creator.PageSize{279.4 * creator.PPMM, 215.9 * creator.PPMM}
	c.SetPageSize(size)
	c.SetPageMargins(20, 20, 44, 44)
	// Read main content template.
	mainTpl, err := readTemplate("templates/main.tpl")
	if err != nil {
		log.Fatal(err)
	}
	barCode, err := makeBarcode("0123456789012", 50, 72)
	if err != nil {
		log.Fatal(err)
	}
	_ = barCode
	arialBold, err := model.NewPdfFontFromTTFFile("./templates/res/arialbd.ttf")
	if err != nil {
		log.Fatal(err)
	}
	arial, err := model.NewPdfFontFromTTFFile("./templates/res/arial.ttf")
	if err != nil {
		log.Fatal(err)
	}
	// Create template options.
	tplOpts := &creator.TemplateOptions{
		ImageMap: map[string]*model.Image{
			"barCode": barCode,
		},
		FontMap: map[string]*model.PdfFont{
			"arial-bold": arialBold,
			"arial":      arial,
		},
	}
	// Read data from json.
	medicationData, err := readData("data.json")
	if err != nil {
		log.Fatal(err)
	}
	// Draw content template.
	if err := c.DrawTemplate(mainTpl, medicationData, tplOpts); err != nil {
		log.Fatal(err)
	}

	c.DrawFooter(func(block *creator.Block, args creator.FooterFunctionArgs) {
		// Read template.
		tpl, err := readTemplate("templates/footer.tpl")
		if err != nil {
			log.Fatal(err)
		}
		// Draw template.
		data := map[string]interface{}{
			"PageNum": args.PageNum,
		}
		if err := block.DrawTemplate(c, tpl, data, nil); err != nil {
			log.Fatal(err)
		}
	})

	// Write to output file.
	if err := c.WriteToFile("unipdf-medication-schedule.pdf"); err != nil {
		log.Fatal(err)
	}
}

func makeBarcode(codeStr string, width float64, oversampling int) (*model.Image, error) {
	bcode, err := ean.Encode(codeStr)
	if err != nil {
		log.Fatal(err)
	}

	pixelWidth := oversampling * int(math.Ceil(width))
	bcodeImg, err := barcode.Scale(bcode, pixelWidth, pixelWidth)
	if err != nil {
		log.Fatal(err)
	}

	image, err := model.ImageHandling.NewImageFromGoImage(bcodeImg)

	return image, err
}

// readTemplate reads the template at the specified file path and returns it as an io.Reader.
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

// readReceipt reads the receipt json file and decodes it to `Receipt` object.
func readData(jsonFile string) (*MedicalData, error) {
	file, err := os.Open(jsonFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data MedicalData
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
