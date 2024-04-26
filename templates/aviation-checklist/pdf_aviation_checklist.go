/*
 * This example showcases the usage of creator templates by creating a sample
 * aviation checklist.
 *
 * Run as: go run pdf_aviation_checklist.go
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

func main() {
	c := creator.New()
	c.SetPageMargins(30, 30, 92, 50)

	// Read main content template.
	mainTpl, err := readTemplate("templates/main.tpl")
	if err != nil {
		log.Fatal(err)
	}

	// Read checklist JSON data.
	checks, err := readChecklistData("checklist.json")
	if err != nil {
		log.Fatal(err)
	}

	// Draw main content template.
	font, err := model.NewStandard14Font(model.HelveticaName)
	if err != nil {
		log.Fatal(err)
	}

	fontBold, err := model.NewStandard14Font(model.HelveticaBoldName)
	if err != nil {
		log.Fatal(err)
	}

	tplOpts := &creator.TemplateOptions{
		HelperFuncMap: template.FuncMap{
			"calcTableHeight": func(check *Checks) float64 {
				titleSp := c.NewStyledParagraph()
				chunk := titleSp.SetText(check.Title)
				chunk.Style.Font = font

				height := titleSp.Height()
				for _, v := range check.Items {
					var chunk *creator.TextChunk
					itemSp := c.NewStyledParagraph()

					if v.Action != nil && *v.Action != "" {
						chunk = itemSp.SetText(*v.Action)
						chunk.Style.Font = fontBold
					} else {
						label := ""
						if v.Check != nil && *v.Check != "" {
							label = *v.Check
						}

						val := ""
						if v.Value != nil && *v.Value != "" {
							label = *v.Value
						}

						chunk = itemSp.SetText(label + val)
						chunk.Style.Font = font
					}

					height += itemSp.Height()
				}

				return height
			},
			"isFitInPageHeight": func(height float64) bool {
				return height <= c.Height()-320
			},
		},
	}

	if err := c.DrawTemplate(mainTpl, checks, tplOpts); err != nil {
		log.Fatal(err)
	}

	// Draw header and footer.
	drawHeader := func(tplPath string, block *creator.Block, version, releaseDate string, pageNum, totalPages int) {
		// Read template.
		tpl, err := readTemplate(tplPath)
		if err != nil {
			log.Fatal(err)
		}

		// Draw template.
		data := map[string]interface{}{
			"Version":     version,
			"ReleaseDate": releaseDate,
			"PageNum":     pageNum,
			"TotalPages":  totalPages,
		}

		if err := block.DrawTemplate(c, tpl, data, tplOpts); err != nil {
			log.Fatal(err)
		}
	}

	c.DrawHeader(func(block *creator.Block, args creator.HeaderFunctionArgs) {
		drawHeader("templates/header.tpl", block, checks.Version, checks.ReleaseDate, args.PageNum, args.TotalPages)
	})
	c.DrawFooter(func(block *creator.Block, args creator.FooterFunctionArgs) {
		drawHeader("templates/footer.tpl", block, checks.Version, checks.ReleaseDate, args.PageNum, args.TotalPages)
	})

	// Write output file.
	if err := c.WriteToFile("unipdf-aviation-checklist.pdf"); err != nil {
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

// readChecklistData reads the checklist data from a specified JSON file.
func readChecklistData(jsonFile string) (*Checklist, error) {
	file, err := os.Open(jsonFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	check := &Checklist{}
	if err := json.NewDecoder(file).Decode(check); err != nil {
		return nil, err
	}

	return check, nil
}

// Checklist represent checklist data.
type Checklist struct {
	Version     string    `json:"version"`
	ReleaseDate string    `json:"release"`
	Checks      []*Checks `json:"checklist"`
}

// Checks represent each group of checklist.
type Checks struct {
	Title string       `json:"title"`
	Items []*CheckItem `json:"items"`
}

// CheckItem represents each item in checklist groups.
type CheckItem struct {
	Action *string `json:"action"`
	Check  *string `json:"check,omitempty"`
	Value  *string `json:"value,omitempty"`
}

// DisplayText returns combination text of checklist label, separator, and value.
func (c *CheckItem) DisplayText() string {
	if c.Value == nil {
		return *c.Check
	}

	// Max line width.
	maxWidth := 258.0
	// Font size use for rendering checklist items.
	fontSize := 10.0

	font, err := model.NewStandard14Font(model.HelveticaName)
	if err != nil {
		log.Fatal(err)
	}

	textWidth := 0.0
	for _, r := range *c.Check + *c.Value {
		textWidth += getRuneWidth(font, r) * fontSize
	}

	if textWidth < maxWidth {
		availWidth := maxWidth - textWidth
		separator := '.'
		sepWidth := getRuneWidth(font, separator) * fontSize
		sepCount := int(availWidth / sepWidth)

		// When the distance between check label and it's value is too narrow
		// to use dots as separator, replace dots with spaces.
		if sepCount <= 5 {
			separator = ' '
			sepWidth := getRuneWidth(font, separator)
			sepCount = int(availWidth / sepWidth)
		}
		sepText := strings.Repeat(string(separator), sepCount)

		return *c.Check + sepText + *c.Value
	}

	return *c.Check + "  " + *c.Value
}

// getRuneWidth calculates rune width based on provided font.
func getRuneWidth(font *model.PdfFont, r rune) float64 {
	metrics, bool := font.GetRuneMetrics(r)
	if !bool {
		log.Fatal("failed to get width")
	}

	return metrics.Wx / 1000
}
