/*
 * This example showcases how to prepare report with charts from csv data.
 *
 * Run as: go run pdf_report_from_csv.go
 */
package main

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/unidoc/unichart"
	"github.com/unidoc/unichart/dataset"
	"github.com/unidoc/unichart/dataset/sequence"
	"github.com/unidoc/unichart/render"
	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

func init() {
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}
}

// chartColors list some colors to be used.
var chartColors = []color.Color{
	color.RGBA{R: 38, G: 198, B: 218},
	color.RGBA{R: 255, G: 167, B: 38},
	color.RGBA{R: 67, G: 160, B: 71},
	color.RGBA{R: 186, G: 104, B: 200},
	color.RGBA{R: 255, G: 64, B: 129},
	color.RGBA{R: 255, G: 109, B: 0},
}

// salesData represents a single row of sales data.
type salesData struct {
	DateOfSale      time.Time
	QuantitySold    float64
	SalePrice       float64
	SalespersonName string
}

func main() {
	robotoFontRegular, err := model.NewPdfFontFromTTFFile("./Roboto-Regular.ttf")
	if err != nil {
		common.Log.Info("Failed to load font %s", err)
		return
	}

	robotoFontPro, err := model.NewPdfFontFromTTFFile("./Roboto-Bold.ttf")
	if err != nil {
		common.Log.Info("Failed to load font %s", err)
		return
	}

	c := creator.New()
	filePath := "./test-data.csv"
	data, _, err := loadCsv(filePath)
	if err != nil {
		common.Log.Info("failed to load data %s", err)
		return
	}

	cumulativeSums := make(map[string]float64) // Cumulative Sales by person
	sum := 0.0
	dates := []time.Time{}    // for the x-axis
	totalSums := []float64{}  // commutative sum over time
	dailySales := []float64{} // daily sales
	for _, dataRow := range data {
		sum += dataRow.SalePrice
		cumulativeSums[dataRow.SalespersonName] += dataRow.SalePrice
		dates = append(dates, dataRow.DateOfSale)
		totalSums = append(totalSums, sum)
		dailySales = append(dailySales, dataRow.SalePrice)
	}

	Contributions := map[string]float64{}
	for person, totalSale := range cumulativeSums {
		Contribution := totalSale / sum
		Contributions[person] = Contribution
	}

	// create front page
	c.CreateFrontPage(func(args creator.FrontpageFunctionArgs) {
		doFirstPage(c, robotoFontRegular, robotoFontPro)
	})

	doFooter(c, robotoFontPro)

	// Create the chapter on page 1
	ch := c.NewChapter("Sales Report")

	chapterFont := robotoFontRegular
	chapterFontColor := creator.ColorRGBFrom8bit(72, 86, 95)
	chapterFontSize := 18.0

	normalFont := robotoFontRegular
	normalFontColor := creator.ColorRGBFrom8bit(72, 86, 95)
	normalFontSize := 10.0

	ch.GetHeading().SetFont(chapterFont)
	ch.GetHeading().SetFontSize(chapterFontSize)
	ch.GetHeading().SetColor(chapterFontColor)

	p := c.NewParagraph("Here we provide sales report for the month of January." +
		"The data reflects a diverse range of sales activities, showcasing the efforts of various salespeople. " +
		"The sales prices ranged from $40 to $250, highlighting the varying value of the products sold." +
		"Overall, this month demonstrated strong performance across the board," +
		" indicating effective sales strategies and engagement by the team throughout January.")
	p.SetFont(normalFont)
	p.SetFontSize(normalFontSize)
	p.SetColor(normalFontColor)
	p.SetMargins(0, 0, 5, 0)
	ch.Add(p)
	c.Draw(ch)

	// Create Bar chart
	p = c.NewParagraph("Total Sales By SalesPerson")
	p.SetFont(robotoFontRegular)
	p.SetFontSize(10)
	p.SetPos(60, 150)
	c.Draw(p)

	bChart := createBarChart(cumulativeSums)
	bChart.SetHeight(200)
	bChart.SetWidth(300)
	barChart := creator.NewChart(bChart)
	barChart.SetPos(50, 170)
	err = c.Draw(barChart)

	if err != nil {
		common.Log.Info("Failed to draw chart. %s.", err)
		return
	}

	// Create pie chart
	p = c.NewParagraph("Sales Contributions by Sales Persons")
	p.SetFont(robotoFontRegular)
	p.SetFontSize(10)
	p.SetPos(370, 150)
	c.Draw(p)

	pChart := createPieChart(Contributions, false)
	pChart.SetHeight(200)
	pChart.SetWidth(200)
	pieChart := creator.NewChart(pChart)
	pieChart.SetPos(360, 170)
	err = c.Draw(pieChart)
	if err != nil {
		common.Log.Info("Failed to draw pie chart. %s", err)
		return
	}

	// Create Line Chart
	title := "Total sales trend over time"
	lChart := createLineChart(dates, title, totalSums, dailySales)
	lChart.SetHeight(300)
	lChart.SetWidth(550)
	lineChart := creator.NewChart(lChart)
	lineChart.SetPos(40, 400)

	err = c.Draw(lineChart)

	if err != nil {
		common.Log.Info("Failed to draw line chart. %s", err)
		return
	}

	err = c.WriteToFile("report_from_csv.pdf")
	if err != nil {
		common.Log.Info("Failed to write to file. %s", err)
		return
	}
}

// doFooter adds the page footers.
func doFooter(c *creator.Creator, font *model.PdfFont) {
	c.DrawFooter(func(block *creator.Block, args creator.FooterFunctionArgs) {
		// Draw the on a block for each page.
		p := c.NewParagraph("unidoc.io")
		p.SetFont(font)
		p.SetFontSize(8)
		p.SetPos(50, 20)
		p.SetColor(creator.ColorRGBFrom8bit(63, 68, 76))
		block.Draw(p)

		strPage := fmt.Sprintf("Page %d of %d", args.PageNum, args.TotalPages)
		p = c.NewParagraph(strPage)
		p.SetFont(font)
		p.SetFontSize(8)
		p.SetPos(300, 20)
		p.SetColor(creator.ColorRGBFrom8bit(63, 68, 76))
		block.Draw(p)
	})
}

// DoFirstPage creates the front page of the document.
func doFirstPage(c *creator.Creator, fontRegular *model.PdfFont, fontBold *model.PdfFont) {
	helvetica, _ := model.NewStandard14Font("Helvetica")
	helveticaBold, _ := model.NewStandard14Font("Helvetica-Bold")

	p := c.NewParagraph("UniDoc")
	p.SetFont(helvetica)
	p.SetFontSize(48)
	p.SetMargins(85, 0, 150, 0)
	p.SetColor(creator.ColorRGBFrom8bit(56, 68, 77))
	c.Draw(p)

	p = c.NewParagraph("Sample Report From CSV")
	p.SetFont(helveticaBold)
	p.SetFontSize(30)
	p.SetMargins(85, 0, 0, 0)
	p.SetColor(creator.ColorRGBFrom8bit(45, 148, 215))
	c.Draw(p)

	t := time.Now().UTC()
	dateStr := t.Format("1 Jan, 2006 15:04")

	p = c.NewParagraph(dateStr)
	p.SetFont(helveticaBold)
	p.SetFontSize(12)
	p.SetMargins(90, 0, 5, 0)
	p.SetColor(creator.ColorRGBFrom8bit(56, 68, 77))
	c.Draw(p)
}

// createBarChart creates a bar chart given `valMap`.
func createBarChart(valMap map[string]float64) render.ChartRenderable {
	chart := &unichart.BarChart{
		Bars:     parseChartValMap(valMap),
		BarWidth: 35,
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

// createLineChart creates line chart using xValues time series  and yValues1, yValues2 series values.
func createLineChart(xValues []time.Time, title string, yValues1, yValues2 []float64) render.ChartRenderable {
	mainSeries := dataset.TimeSeries{
		XValues: xValues,
		YValues: yValues1,
		Name:    "Total Sales Price",
	}

	secondSeries := dataset.TimeSeries{
		XValues: xValues,
		YValues: yValues2,
		Name:    "Sales By Person",
	}
	ch := &unichart.Chart{
		Series: []dataset.Series{
			mainSeries,
			secondSeries,
		},
		Title: title,
		XAxis: unichart.XAxis{
			Name: "Time",
		},
		YAxis: unichart.YAxis{
			Name: "Sales Price",
		},
	}
	ch.Elements = []render.Renderable{
		unichart.Legend(ch),
	}
	return ch
}

// createPieChart creates pie chart based on the valMap values.
func createPieChart(valMap map[string]float64, isDonut bool) render.ChartRenderable {
	var (
		vals = make([]dataset.Value, 0, len(valMap))
		idx  = 0
	)

	for key, val := range valMap {
		vals = append(vals, dataset.Value{
			Label: fmt.Sprintf("%s \n %.2f%% ", key, val*100),
			Value: val,
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

// parseChartValMap parses the valMap and returns an array of dataset.Value.
func parseChartValMap(valMap map[string]float64) []dataset.Value {
	var (
		vals = make([]dataset.Value, 0, len(valMap))
		idx  = 0
	)

	for key, val := range valMap {

		vals = append(vals, dataset.Value{
			Label: key,
			Value: val,
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

// loadCsv loads values from csv files and returns SalesData object and the list of csv headers.
func loadCsv(filePath string) ([]salesData, []string, error) {
	file, err := os.Open(filePath) // Ensure the filename matches your CSV file
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	var sales []salesData
	var header []string

	for i, row := range rows {
		if i == 0 {
			header = append(header, row...)
			continue
		}

		dateOfSale, err := time.Parse("2006-01-02", row[0])
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing sale price: %v", err)
		}

		quantitySold, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing sale price: %v", err)
		}

		salePrice, err := strconv.ParseFloat(row[2], 64)
		if err != nil {

			return nil, nil, fmt.Errorf("error parsing sale price: %v", err)
		}

		sale := salesData{
			DateOfSale:      dateOfSale,
			QuantitySold:    quantitySold,
			SalePrice:       salePrice,
			SalespersonName: row[3],
		}
		sales = append(sales, sale)
	}

	return sales, header, nil
}
