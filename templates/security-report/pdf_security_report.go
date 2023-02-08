/*
 * This example showcases the usage of creator templates by creating a sample
 * security report.
 *
 * Run as: go run pdf_security_report.go
 */

package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"sort"
	"text/template"

	"github.com/unidoc/unichart"
	"github.com/unidoc/unichart/dataset"
	"github.com/unidoc/unichart/dataset/sequence"
	"github.com/unidoc/unichart/render"
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
	c.AddTOC = true
	c.SetPageMargins(30, 30, 70, 40)

	// Read main content template.
	mainTpl, err := readTemplate("templates/main.tpl")
	if err != nil {
		log.Fatal(err)
	}

	// Draw main content template.
	chartMap := map[string]render.ChartRenderable{}
	tplOpts := &creator.TemplateOptions{
		HelperFuncMap: template.FuncMap{
			"CreateLineChart": func(name string, points int, min, max float64) string {
				chartMap[name] = createLineChart(points, min, max)
				return name
			},
			"CreateBarChart": func(name string, valMap map[string]interface{}) string {
				chartMap[name] = createBarChart(valMap)
				return name
			},
			"CreateStackedBar": createStackedBar,
			"CreateStackedBarChart": func(name string, bars ...unichart.StackedBar) string {
				chartMap[name] = createStackedBarChart(bars...)
				return name
			},
			"CreatePieChart": func(name string, valMap map[string]interface{}) string {
				chartMap[name] = createPieChart(valMap)
				return name
			},
		},
		ChartMap: chartMap,
	}

	if err := c.DrawTemplate(mainTpl, nil, tplOpts); err != nil {
		log.Fatal(err)
	}

	// Draw header.
	c.DrawHeader(func(header *creator.Block, args creator.HeaderFunctionArgs) {
		// Skip front page.
		if args.PageNum == 1 {
			return
		}

		// Read header template.
		headerTpl, err := readTemplate("templates/header.tpl")
		if err != nil {
			log.Fatal(err)
		}

		// Draw header template.
		if err := header.DrawTemplate(c, headerTpl, nil, nil); err != nil {
			log.Fatal(err)
		}
	})

	// Draw footer.
	c.DrawFooter(func(footer *creator.Block, args creator.FooterFunctionArgs) {
		// Skip front page.
		if args.PageNum == 1 {
			return
		}

		// Read footer template.
		footerTpl, err := readTemplate("templates/footer.tpl")
		if err != nil {
			log.Fatal(err)
		}

		// Draw footer template.
		data := map[string]interface{}{
			"PageNum":    args.PageNum,
			"TotalPages": args.TotalPages,
		}

		if err := footer.DrawTemplate(c, footerTpl, data, nil); err != nil {
			log.Fatal(err)
		}
	})

	// Draw front page.
	c.CreateFrontPage(func(args creator.FrontpageFunctionArgs) {
		// Read front page template.
		frontPageTpl, err := readTemplate("templates/front-page.tpl")
		if err != nil {
			log.Fatal(err)
		}

		// Draw front page template.
		if err := c.DrawTemplate(frontPageTpl, nil, nil); err != nil {
			log.Fatal(err)
		}
	})

	// Write output file.
	if err := c.WriteToFile("unipdf-security-report.pdf"); err != nil {
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

// createPieChart creates a sample pie chart based on the provided values map.
func createPieChart(valMap map[string]interface{}) render.ChartRenderable {
	vals := make([]dataset.Value, 0, len(valMap))
	for key, val := range valMap {
		fVal, _ := val.(float64)
		vals = append(vals, dataset.Value{
			Label: key,
			Value: fVal,
			Style: render.Style{
				FontSize:    6,
				StrokeWidth: 1,
			},
		})
	}

	sort.Slice(vals, func(i, j int) bool {
		return vals[i].Value < vals[j].Value
	})

	return &unichart.PieChart{
		Values: vals,
	}
}

// createPieChart creates a sample line chart based on the provided configuration.
func createLineChart(points int, min, max float64) render.ChartRenderable {
	mainSeries := dataset.ContinuousSeries{
		XValues: sequence.Wrapper{Sequence: sequence.NewLinearSequence().WithStart(1.0).WithEnd(float64(points)).WithStep(1)}.Values(),
		YValues: sequence.Wrapper{Sequence: sequence.NewRandomSequence().WithLen(points).WithMin(min).WithMax(max)}.Values(),
	}
	linRegSeries := &dataset.LinearRegressionSeries{
		InnerSeries: mainSeries,
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
		BarWidth:   20,
		BarSpacing: 20,
		Bars:       parseChartValMap(valMap),
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

// createPieChart creates a sample stacked bar chart based on the provided bar values.
func createStackedBarChart(bars ...unichart.StackedBar) render.ChartRenderable {
	return &unichart.StackedBarChart{
		YAxis: render.Style{
			FontSize: 6,
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
	vals := make([]dataset.Value, 0, len(valMap))
	for key, val := range valMap {
		fVal, _ := val.(float64)
		vals = append(vals, dataset.Value{
			Label: key,
			Value: fVal,
			Style: render.Style{
				FontSize: 6,
			},
		})
	}

	sort.Slice(vals, func(i, j int) bool {
		return vals[i].Value < vals[j].Value
	})

	return vals
}
