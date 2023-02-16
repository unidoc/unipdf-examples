/*
 * This example showcases the usage of creator templates by creating a sample medication schedule report.
 *
 * Run as: go run pdf-medication-schedule.go
 */

package main

import (
	"bytes"
	"io"
	"log"
	"os"

	"github.com/unidoc/unipdf/v3/creator"
)

// func init() {
// 	// Make sure to load your metered License API key prior to using the library.
// 	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io.
// 	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
// 	if err != nil {
// 		panic(err)
// 	}
// 	common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
// }

func main() {
	c := creator.New()
	c.SetPageMargins(90, 60, 95, 300)
	// pageSize 279.4 x 215.9
	size := creator.PageSize{279.4 * creator.PPMM, 215.9 * creator.PPMM}
	// c.SetPageSize(creator.PageSizeA4)
	c.SetPageSize(size)
	// Read main content template.
	mainTpl, err := readTemplate("templates/main.tpl")
	if err != nil {
		log.Fatal(err)
	}

	// Draw content template.
	if err := c.DrawTemplate(mainTpl, nil, nil); err != nil {
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
