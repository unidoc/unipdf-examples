package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/core"
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
	c.SetPageMargins(30, 30, 20, 20)

	// Read main content template.
	mainTpl, err := readTemplate("templates/main.tpl")
	if err != nil {
		log.Fatal(err)
	}

	// Draw main content teplate.
	tplOpts := &creator.TemplateOptions{
		HelperFuncMap: template.FuncMap{
			"formatTime": func(val, format string) string {
				t, _ := time.Parse("2006-01-02T15:04:05", val)
				return t.Format(format)
			},
			"extendDict": func(m map[string]interface{}, params ...interface{}) (map[string]interface{}, error) {
				lenParams := len(params)
				if lenParams%2 != 0 {
					return nil, core.ErrRangeError
				}

				for i := 0; i < lenParams; i += 2 {
					key, ok := params[i].(string)
					if !ok {
						return nil, core.ErrTypeError
					}

					m[key] = params[i+1]
				}

				return m, nil
			},
		},
	}

	if err := c.DrawTemplate(mainTpl, nil, tplOpts); err != nil {
		log.Fatal(err)
	}

	// Write output file.
	if err := c.WriteToFile("unipdf-medical-bill.pdf"); err != nil {
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
