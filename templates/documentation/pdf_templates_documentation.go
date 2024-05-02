/*
 * This example generates the documentation of the creator templates feature.
 *
 * Run as: go run pdf_templates_documentation.go
 */

package main

import (
	"bytes"
	"encoding/xml"
	"image/color"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/unidoc/unichart"
	"github.com/unidoc/unichart/dataset"
	"github.com/unidoc/unichart/dataset/sequence"
	"github.com/unidoc/unichart/render"
	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/model/optimize"
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
	c.SetPageMargins(25, 25, 50, 30)

	// Create common template options.
	tplOpts := createTplOpts(c)

	// Draw main content.
	drawContent(c, tplOpts)

	// Draw header and footer.
	drawHeaders(c, tplOpts)

	// Draw front page.
	drawFrontPage(c, tplOpts)

	// Draw TOC.
	drawTOC(c, tplOpts)

	// Set optimizer.
	c.SetOptimizer(optimize.New(optimize.Options{
		CombineDuplicateDirectObjects:   true,
		CombineIdenticalIndirectObjects: true,
		CombineDuplicateStreams:         true,
		CompressStreams:                 true,
		UseObjectStreams:                true,
	}))

	// Write output file.
	if err := c.WriteToFile("unipdf-templates-documentation.pdf"); err != nil {
		log.Fatal(err)
	}
}

// drawContent renders the main content templates.
func drawContent(c *creator.Creator, tplOpts *creator.TemplateOptions) {
	// Read main content template.
	mainTpl, err := readTemplate("templates/main.tpl")
	if err != nil {
		log.Fatal(err)
	}

	// Read chapter templates.
	files, err := os.ReadDir("templates/chapters")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		filename := file.Name()

		chapterTpl, err := readTemplate(filepath.Join("templates/chapters", filename))
		if err != nil {
			log.Fatal(err)
		}

		tplOpts.SubtemplateMap[strings.TrimSuffix(filename, filepath.Ext(filename))] = chapterTpl
	}

	// Draw main content template.
	data := map[string]interface{}{
		"newline": "&#xA;",
	}

	if err := c.DrawTemplate(mainTpl, data, tplOpts); err != nil {
		log.Fatal(err)
	}
}

// drawHeaders renders the header and footer templates.
func drawHeaders(c *creator.Creator, tplOpts *creator.TemplateOptions) {
	drawHeader := func(tplPath string, block *creator.Block, pageNum, totalPages int) {
		// Skip front page.
		if pageNum == 1 {
			return
		}

		// Read template.
		tpl, err := readTemplate(tplPath)
		if err != nil {
			log.Fatal(err)
		}

		// Draw template.
		data := map[string]interface{}{
			"Date":       time.Now(),
			"PageNum":    pageNum,
			"TotalPages": totalPages,
			"newline":    "&#xA;",
		}

		if err := block.DrawTemplate(c, tpl, data, tplOpts); err != nil {
			log.Fatal(err)
		}
	}

	// Draw header.
	c.DrawHeader(func(block *creator.Block, args creator.HeaderFunctionArgs) {
		drawHeader("templates/header.tpl", block, args.PageNum, args.TotalPages)
	})

	// Draw footer.
	c.DrawFooter(func(block *creator.Block, args creator.FooterFunctionArgs) {
		drawHeader("templates/footer.tpl", block, args.PageNum, args.TotalPages)
	})
}

// drawFrontPage renders the front page template.
func drawFrontPage(c *creator.Creator, tplOpts *creator.TemplateOptions) {
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
}

// drawTOC renders the table of contents of the document.
func drawTOC(c *creator.Creator, tplOpts *creator.TemplateOptions) {
	c.AddTOC = true
	c.CreateTableOfContents(func(toc *creator.TOC) error {
		primaryColor := tplOpts.ColorMap["primary"]

		heading := toc.Heading()
		heading.SetTextAlignment(creator.TextAlignmentCenter)
		heading.SetMargins(0, 0, 25, 25)
		headingChunk := heading.SetText("Table of Contents")
		headingChunk.Style.FontSize = 14
		headingChunk.Style.Color = primaryColor

		helvBold, err := model.NewStandard14Font(model.HelveticaBoldName)
		if err != nil {
			log.Fatal(err)
		}

		for _, line := range toc.Lines() {
			line.Number.Style.Color = primaryColor
			line.Number.Style.Font = helvBold
			line.Title.Style.Color = primaryColor
			line.Title.Style.Font = helvBold
			line.Separator.Style.Color = primaryColor
			line.Page.Style.Color = primaryColor
			line.Page.Style.Font = helvBold
		}

		return nil
	})
}

// createTplOpts returns the common options used by all rendered templates.
func createTplOpts(c *creator.Creator) *creator.TemplateOptions {
	// Create image map.
	logoImg, err := newImageFromFile("templates/res/images/logo.png")
	if err != nil {
		log.Fatal(err)
	}

	sampleImg, err := newImageFromFile("templates/res/images/sample-image-2.jpg")
	if err != nil {
		log.Fatal(err)
	}

	imageMap := map[string]*model.Image{
		"logo":   logoImg,
		"sample": sampleImg,
	}

	// Create color map.
	primaryBgGradient := c.NewLinearGradientColor([]*creator.ColorPoint{
		creator.NewColorPoint(creator.ColorRGBFromHex("#54b7de"), 0.0),
		creator.NewColorPoint(creator.ColorRGBFromHex("#0772cd"), 0.3),
		creator.NewColorPoint(creator.ColorRGBFromHex("#3e5ede"), 1.0),
	})
	primaryBgGradient.SetAngle(180)

	secondaryBgGradient := c.NewLinearGradientColor([]*creator.ColorPoint{
		creator.NewColorPoint(creator.ColorRGBFromHex("#f00c27"), 0.0),
		creator.NewColorPoint(creator.ColorRGBFromHex("#0772cd"), 0.2),
		creator.NewColorPoint(creator.ColorRGBFromHex("#0772cd"), 0.8),
		creator.NewColorPoint(creator.ColorRGBFromHex("#f00c27"), 1.0),
	})
	secondaryBgGradient.SetAngle(180)

	primaryLightBgGradient := c.NewLinearGradientColor([]*creator.ColorPoint{
		creator.NewColorPoint(creator.ColorRGBFromHex("#3e5ede"), 0.0),
		creator.NewColorPoint(creator.ColorRGBFromHex("#0772cd"), 0.3),
		creator.NewColorPoint(creator.ColorRGBFromHex("#54b7de"), 1.0),
	})
	primaryLightBgGradient.SetAngle(180)

	colorMap := map[string]creator.Color{
		"primary":                   creator.ColorRGBFromHex("#0772cd"),
		"secondary":                 creator.ColorRGBFromHex("#f00c27"),
		"text":                      creator.ColorRGBFromHex("#333333"),
		"light-gray":                creator.ColorRGBFromHex("#f9fafe"),
		"medium-gray":               creator.ColorRGBFromHex("#dce0ef"),
		"white":                     creator.ColorWhite,
		"black":                     creator.ColorBlack,
		"primary-bg-gradient":       primaryBgGradient,
		"secondary-bg-gradient":     secondaryBgGradient,
		"primary-light-bg-gradient": primaryLightBgGradient,
	}

	// Create font map.
	dejaVuSansMono, err := model.NewPdfFontFromTTFFile("templates/res/fonts/DejaVuSansMono.ttf")
	if err != nil {
		log.Fatal(err)
	}

	fontMap := map[string]*model.PdfFont{
		"deja-vu-sans-mono": dejaVuSansMono,
	}

	// Create chart map.
	chartMap := map[string]render.ChartRenderable{}

	// Create function map.
	funcMap := template.FuncMap{
		"loop": func(size uint64) []struct{} {
			return make([]struct{}, size)
		},
		"sum": func(elements ...interface{}) int {
			var sum int
			for _, element := range elements {
				val, _ := element.(int)
				sum += val
			}
			return sum
		},
		"now":        time.Now,
		"strToUpper": strings.ToUpper,
		"xmlEscape":  xmlEscapeText,
		"createLineChart": func(name string, points int, min, max float64) string {
			chartMap[name] = createLineChart(points, min, max)
			return name
		},
		"createBarChart": func(name string, valMap map[string]interface{}) string {
			chartMap[name] = createBarChart(valMap)
			return name
		},
		"createStackedBar": createStackedBar,
		"createStackedBarChart": func(name string, bars ...unichart.StackedBar) string {
			chartMap[name] = createStackedBarChart(bars...)
			return name
		},
		"createPieChart": func(name string, valMap map[string]interface{}, isDonut bool) string {
			chartMap[name] = createPieChart(valMap, isDonut)
			return name
		},
	}

	// Load subtemplates.
	helpersTpl, err := readTemplate("templates/helpers.tpl")
	if err != nil {
		log.Fatal(err)
	}

	return &creator.TemplateOptions{
		HelperFuncMap: funcMap,
		FontMap:       fontMap,
		ImageMap:      imageMap,
		ColorMap:      colorMap,
		SubtemplateMap: map[string]io.Reader{
			"helpers": helpersTpl,
		},
		ChartMap: chartMap,
	}
}

// readTemplate reads the template at the specified file path
// and returns it as an io.Reader.
func readTemplate(tplFile string) (*bytes.Buffer, error) {
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

// newImageFromFile loads and returns the image at the specified path.
func newImageFromFile(path string) (*model.Image, error) {
	// Open image file.
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Load the image using the default handler.
	img, err := model.ImageHandling.Read(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// xmlEscapeText XML escapes the specified text.
func xmlEscapeText(text string) (string, error) {
	b := bytes.NewBuffer(nil)
	if err := xml.EscapeText(b, []byte(text)); err != nil {
		return "", err
	}

	return b.String(), nil
}

// chartColors represents the collection of colors used by the chart components.
var chartColors = []color.Color{
	color.RGBA{R: 38, G: 198, B: 218},
	color.RGBA{R: 255, G: 167, B: 38},
	color.RGBA{R: 67, G: 160, B: 71},
	color.RGBA{R: 186, G: 104, B: 200},
	color.RGBA{R: 255, G: 64, B: 129},
	color.RGBA{R: 255, G: 109, B: 0},
}

// createPieChart creates a sample pie or donut chart based on
// the provided values map.
func createPieChart(valMap map[string]interface{}, isDonut bool) render.ChartRenderable {
	var (
		vals = make([]dataset.Value, 0, len(valMap))
		idx  = 0
	)

	for key, val := range valMap {
		fVal, _ := val.(float64)
		vals = append(vals, dataset.Value{
			Label: key,
			Value: fVal,
			Style: render.Style{
				FontSize:    8,
				FillColor:   chartColors[idx],
				StrokeWidth: 1,
			},
		})

		idx++
		if idx >= len(chartColors) {
			idx = 0
		}
	}

	sort.Slice(vals, func(i, j int) bool {
		return vals[i].Value < vals[j].Value
	})

	if isDonut {
		return &unichart.DonutChart{
			Values: vals,
		}
	}

	return &unichart.PieChart{
		Values: vals,
	}
}

// createPieChart creates a sample line chart based on
// the provided configuration.
func createLineChart(points int, min, max float64) render.ChartRenderable {
	mainSeries := dataset.ContinuousSeries{
		XValues: sequence.Wrapper{Sequence: sequence.NewLinearSequence().WithStart(1.0).WithEnd(float64(points)).WithStep(1)}.Values(),
		YValues: sequence.Wrapper{Sequence: sequence.NewRandomSequence().WithLen(points).WithMin(min).WithMax(max)}.Values(),
		Style: render.Style{
			FillColor: chartColors[0],
		},
	}
	linRegSeries := &dataset.LinearRegressionSeries{
		InnerSeries: mainSeries,
		Style: render.Style{
			StrokeColor: chartColors[1],
		},
	}

	return &unichart.Chart{
		Series: []dataset.Series{
			mainSeries,
			linRegSeries,
		},
	}
}

// createPieChart creates a sample bar chart based on the provided values map.
func createBarChart(valMap map[string]interface{}) render.ChartRenderable {
	chart := &unichart.BarChart{
		Bars: parseChartValMap(valMap),
	}

	// Set Y-axis custom range.
	var max float64
	for _, bar := range chart.Bars {
		if max < bar.Value {
			max = bar.Value
		}
	}

	rng := &sequence.ContinuousRange{}
	rng.SetMin(0)
	rng.SetMax(max)
	chart.YAxis.Range = rng

	return chart
}

// createPieChart creates a sample stacked bar chart based on
// the provided bar values.
func createStackedBarChart(bars ...unichart.StackedBar) render.ChartRenderable {
	return &unichart.StackedBarChart{
		YAxis: render.Style{
			FontSize: 8,
		},
		BarSpacing:   5,
		IsHorizontal: true,
		Bars:         bars,
	}
}

// createStackedBar creates a stacked bar based on the provided values map.
// This function is exposed via a helper function so that stacked bar can
// created directly inside template files.
func createStackedBar(name string, valMap map[string]interface{}) unichart.StackedBar {
	return unichart.StackedBar{
		Name:   name,
		Width:  20,
		Values: parseChartValMap(valMap),
	}
}

// parseChartValMap parses a value map and returns it as a dataset slice.
func parseChartValMap(valMap map[string]interface{}) []dataset.Value {
	var (
		vals = make([]dataset.Value, 0, len(valMap))
		idx  = 0
	)

	for key, val := range valMap {
		fVal, _ := val.(float64)
		vals = append(vals, dataset.Value{
			Label: key,
			Value: fVal,
			Style: render.Style{
				FontSize:    8,
				FillColor:   chartColors[idx],
				StrokeColor: chartColors[idx],
			},
		})

		idx++
		if idx >= len(chartColors) {
			idx = 0
		}
	}

	sort.Slice(vals, func(i, j int) bool {
		return vals[i].Value < vals[j].Value
	})

	return vals
}
