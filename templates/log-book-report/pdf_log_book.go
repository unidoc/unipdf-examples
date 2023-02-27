/*
 * This example showcases the usage of creator templates by creating a sample
 * log book report.
 *
 * Run as: go pdf_log_book.go
 */
package main

import (
	"bytes"
	"io"
	"log"
	"os"

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

func main() {
	c := creator.New()
	size := creator.PageSize{842, 595}
	c.SetPageSize(size)
	c.SetPageMargins(10, 10, 35, 35)

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

	tplOpts := &creator.TemplateOptions{
		FontMap: map[string]*model.PdfFont{
			"exo-bold":    exoBold,
			"exo-italic":  exoItalic,
			"exo-regular": exoRegular,
		}}

	if err = c.DrawTemplate(mainTpl, tplOpts, nil); err != nil {
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
		if err := c.DrawTemplate(frontPageTpl, nil, tplOpts); err != nil {
			log.Fatal(err)
		}
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
