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
	"github.com/unidoc/unioffice/common/license"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

func init() {
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully set the Unidoc License Key")
}

var chartColors = []color.Color{
	color.RGBA{R: 38, G: 198, B: 218},
	color.RGBA{R: 255, G: 167, B: 38},
	color.RGBA{R: 67, G: 160, B: 71},
	color.RGBA{R: 186, G: 104, B: 200},
	color.RGBA{R: 255, G: 64, B: 129},
	color.RGBA{R: 255, G: 109, B: 0},
}

type SalesData struct {
	DateOfSale      time.Time
	QuantitySold    float64
	SalePrice       float64
	SalespersonName string
}

func main() {
	robotoFontRegular, err := model.NewPdfFontFromTTFFile("./Roboto-Regular.ttf")
	if err != nil {
		panic(err)
	}

	robotoFontPro, err := model.NewPdfFontFromTTFFile("./Roboto-Bold.ttf")
	if err != nil {
		panic(err)
	}

	c := creator.New()
	filePath := "./test-data.csv"
	data, _, err := load_csv(filePath)
	if err != nil {
		panic("failed to load data ")
	}

	logoImg, err := c.NewImageFromFile("./unidoc-logo.png")
	if err != nil {
		panic("filed to load image %s")
	}

	cumulativeSums := make(map[string]float64)
	sum := 0.0
	dates := []time.Time{}
	totalSums := []float64{}
	dailySales := []float64{}
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
		DoFirstPage(c, robotoFontRegular, robotoFontPro)
	})

	doFooter(c, robotoFontPro)
	c.DrawHeader(func(block *creator.Block, args creator.HeaderFunctionArgs) {
		// Draw the header on a block. The block size is the size of the page's top margins.
		block.Draw(logoImg)
	})

	chartTile := "Sales Contributions by Sales Persons"
	pChart := createPieChart(Contributions, false, chartTile)
	pChart.SetHeight(200)
	pChart.SetWidth(200)
	pieChart := creator.NewChart(pChart)
	pieChart.SetPos(350, 30)
	err = c.Draw(pieChart)
	if err != nil {
		panic(err)
	}

	barTitle := "Total Sales By Sales Person"
	bChart := createBarChart(cumulativeSums, barTitle)
	bChart.SetHeight(200)
	bChart.SetWidth(300)
	barChart := creator.NewChart(bChart)
	barChart.SetPos(10, 30)
	err = c.Draw(barChart)

	if err != nil {
		panic(err)
	}

	lChart := createLineChart(dates, totalSums, dailySales)
	lChart.SetHeight(350)
	lChart.SetWidth(600)
	lineChart := creator.NewChart(lChart)
	lineChart.SetPos(5, 300)

	err = c.Draw(lineChart)

	if err != nil {
		panic(err)
	}

	err = c.WriteToFile("bar-chart.pdf")
	if err != nil {
		panic(err)
	}
}

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

func DoFirstPage(c *creator.Creator, fontRegular *model.PdfFont, fontBold *model.PdfFont) {
	helvetica, _ := model.NewStandard14Font("Helvetica")
	helveticaBold, _ := model.NewStandard14Font("Helvetica-Bold")

	p := c.NewParagraph("UniDoc")
	p.SetFont(helvetica)
	p.SetFontSize(48)
	p.SetMargins(85, 0, 150, 0)
	p.SetColor(creator.ColorRGBFrom8bit(56, 68, 77))
	c.Draw(p)

	p = c.NewParagraph("Example Report")
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

func createBarChart(valMap map[string]float64, title string) render.ChartRenderable {
	chart := &unichart.BarChart{
		Bars:  parseChartValMap(valMap),
		Title: title,
		TitleStyle: render.Style{
			FontSize:            12,
			TextHorizontalAlign: render.TextHorizontalAlignLeft,
			Padding: render.Box{
				Top:    10,
				Bottom: 10,
				IsSet:  true,
			},
		},
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

func createLineChart(xValues1 []time.Time, yValues1, yValues2 []float64) render.ChartRenderable {
	mainSeries := dataset.TimeSeries{
		XValues: xValues1,
		YValues: yValues1,
		Name:    "Total Sale",
	}
	secondSeries := dataset.TimeSeries{
		XValues: xValues1,
		YValues: yValues2,
		Name:    "Sales By Person",
	}
	ch := &unichart.Chart{
		Series: []dataset.Series{
			mainSeries,
			secondSeries,
		},
	}
	ch.Elements = []render.Renderable{
		unichart.Legend(ch),
	}
	return ch
}

func createPieChart(valMap map[string]float64, isDonut bool, title string) render.ChartRenderable {
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
		Title:  title,

		TitleStyle: render.Style{
			FontSize:            12,
			TextHorizontalAlign: render.TextHorizontalAlignCenter,
			Padding: render.Box{
				Top:    10,
				Bottom: 10,
				IsSet:  true,
			},
		},
	}
}

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

func load_csv(filePath string) ([]SalesData, []string, error) {
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

	var sales []SalesData
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

		// Create a new Sale object and populate it
		sale := SalesData{
			DateOfSale:      dateOfSale,
			QuantitySold:    quantitySold,
			SalePrice:       salePrice,
			SalespersonName: row[3],
		}

		sales = append(sales, sale)
	}
	return sales, header, nil
}
