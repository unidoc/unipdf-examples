/**
 * This is a sample Go program that demonstrates how to use the UniPDF library
 * to extract text from within images in a PDF using an OCR service that returns
 * HOCR formatted output then writes the reconstructed text to a new PDF.
 *
 * This example uses https://github.com/unidoc/ocrserver as the OCR service.
 * However, UniPDF's OCR API is designed to support other OCR services that accept
 * image uploads via HTTP and return text or HOCR formatted results.
 *
 * Run as: go run reconstruct_pdf_from_hocr.go input.pdf
 */
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/unidoc/unipdf/v4/common/license"
	"github.com/unidoc/unipdf/v4/creator"
	"github.com/unidoc/unipdf/v4/extractor"
	"github.com/unidoc/unipdf/v4/model"
	"github.com/unidoc/unipdf/v4/ocr"
)

// BBox represents a bounding box with coordinates.
type BBox struct {
	X0 int
	Y0 int
	X1 int
	Y1 int
}

// Baseline represents baseline information for text.
type Baseline struct {
	Slope  float64
	Offset float64
}

// TitleAttributes contains parsed attributes from the title field.
type TitleAttributes struct {
	BBox        *BBox
	Baseline    *Baseline
	XSize       float64
	XDescenders float64
	XAscenders  float64
	XWConf      int    // Word confidence (0-100)
	PageNo      int    // Page number
	Image       string // Image path/name
}

// OCRWord represents a word element in hOCR.
type OCRWord struct {
	XMLName xml.Name `xml:"span"`
	Class   string   `xml:"class,attr"`
	ID      string   `xml:"id,attr"`
	Title   string   `xml:"title,attr"`
	Content string   `xml:",chardata"`
}

// OCRLine represents a line element in hOCR.
type OCRLine struct {
	XMLName xml.Name  `xml:"span"`
	Class   string    `xml:"class,attr"`
	ID      string    `xml:"id,attr"`
	Title   string    `xml:"title,attr"`
	Words   []OCRWord `xml:"span"`
}

// OCRPar represents a paragraph element in hOCR.
type OCRPar struct {
	XMLName xml.Name  `xml:"p"`
	Class   string    `xml:"class,attr"`
	ID      string    `xml:"id,attr"`
	Lang    string    `xml:"lang,attr"`
	Title   string    `xml:"title,attr"`
	Lines   []OCRLine `xml:"span"`
}

// OCRCArea represents a column area element in hOCR.
type OCRCArea struct {
	XMLName xml.Name `xml:"div"`
	Class   string   `xml:"class,attr"`
	ID      string   `xml:"id,attr"`
	Title   string   `xml:"title,attr"`
	Pars    []OCRPar `xml:"p"`
}

// OCRPage represents a page element in hOCR.
type OCRPage struct {
	XMLName xml.Name   `xml:"div"`
	Class   string     `xml:"class,attr"`
	ID      string     `xml:"id,attr"`
	Title   string     `xml:"title,attr"`
	CAreas  []OCRCArea `xml:"div"`
}

// HOCRDocument represents the root hOCR document structure.
type HOCRDocument struct {
	XMLName xml.Name  `xml:"div"`
	Pages   []OCRPage `xml:"div"`
}

// ParseTitleAttributes parses the title attribute string and extracts structured data.
func ParseTitleAttributes(title string) *TitleAttributes {
	attrs := &TitleAttributes{}

	// Regular expressions for parsing different attributes
	bboxRe := regexp.MustCompile(`bbox\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)`)
	baselineRe := regexp.MustCompile(`baseline\s+([-\d.]+)\s+([-\d.]+)`)
	xSizeRe := regexp.MustCompile(`x_size\s+([\d.]+)`)
	xDescendersRe := regexp.MustCompile(`x_descenders\s+([\d.]+)`)
	xAscendersRe := regexp.MustCompile(`x_ascenders\s+([\d.]+)`)
	xWConfRe := regexp.MustCompile(`x_wconf\s+(\d+)`)
	pagenoRe := regexp.MustCompile(`ppageno\s+(\d+)`)
	imageRe := regexp.MustCompile(`image\s+"([^"]*)"`)

	// Parse bbox
	if matches := bboxRe.FindStringSubmatch(title); matches != nil {
		x0, _ := strconv.Atoi(matches[1])
		y0, _ := strconv.Atoi(matches[2])
		x1, _ := strconv.Atoi(matches[3])
		y1, _ := strconv.Atoi(matches[4])
		attrs.BBox = &BBox{X0: x0, Y0: y0, X1: x1, Y1: y1}
	}

	// Parse baseline
	if matches := baselineRe.FindStringSubmatch(title); matches != nil {
		slope, _ := strconv.ParseFloat(matches[1], 64)
		offset, _ := strconv.ParseFloat(matches[2], 64)
		attrs.Baseline = &Baseline{Slope: slope, Offset: offset}
	}

	// Parse x_size
	if matches := xSizeRe.FindStringSubmatch(title); matches != nil {
		attrs.XSize, _ = strconv.ParseFloat(matches[1], 64)
	}

	// Parse x_descenders
	if matches := xDescendersRe.FindStringSubmatch(title); matches != nil {
		attrs.XDescenders, _ = strconv.ParseFloat(matches[1], 64)
	}

	// Parse x_ascenders
	if matches := xAscendersRe.FindStringSubmatch(title); matches != nil {
		attrs.XAscenders, _ = strconv.ParseFloat(matches[1], 64)
	}

	// Parse x_wconf (word confidence)
	if matches := xWConfRe.FindStringSubmatch(title); matches != nil {
		attrs.XWConf, _ = strconv.Atoi(matches[1])
	}

	// Parse ppageno
	if matches := pagenoRe.FindStringSubmatch(title); matches != nil {
		attrs.PageNo, _ = strconv.Atoi(matches[1])
	}

	// Parse image
	if matches := imageRe.FindStringSubmatch(title); matches != nil {
		attrs.Image = matches[1]
	}

	return attrs
}

// GetText returns the text content from a word, trimming whitespace.
func (w *OCRWord) GetText() string {
	return strings.TrimSpace(w.Content)
}

// GetAttributes returns parsed title attributes for the word.
func (w *OCRWord) GetAttributes() *TitleAttributes {
	return ParseTitleAttributes(w.Title)
}

// GetText returns the concatenated text content from all words in the line.
func (l *OCRLine) GetText() string {
	var text strings.Builder
	for i, word := range l.Words {
		if i > 0 {
			text.WriteString(" ")
		}
		text.WriteString(word.GetText())
	}
	return text.String()
}

// GetAttributes returns parsed title attributes for the line.
func (l *OCRLine) GetAttributes() *TitleAttributes {
	return ParseTitleAttributes(l.Title)
}

// GetText returns the concatenated text content from all lines in the paragraph.
func (p *OCRPar) GetText() string {
	var text strings.Builder
	for i, line := range p.Lines {
		if i > 0 {
			text.WriteString("\n")
		}
		text.WriteString(line.GetText())
	}
	return text.String()
}

// GetAttributes returns parsed title attributes for the paragraph.
func (p *OCRPar) GetAttributes() *TitleAttributes {
	return ParseTitleAttributes(p.Title)
}

// GetText returns the concatenated text content from all paragraphs in the area.
func (c *OCRCArea) GetText() string {
	var text strings.Builder
	for i, par := range c.Pars {
		if i > 0 {
			text.WriteString("\n\n")
		}
		text.WriteString(par.GetText())
	}
	return text.String()
}

// GetAttributes returns parsed title attributes for the column area.
func (c *OCRCArea) GetAttributes() *TitleAttributes {
	return ParseTitleAttributes(c.Title)
}

// GetText returns the concatenated text content from all areas in the page.
func (p *OCRPage) GetText() string {
	var text strings.Builder
	for i, carea := range p.CAreas {
		if i > 0 {
			text.WriteString("\n\n")
		}
		text.WriteString(carea.GetText())
	}
	return text.String()
}

// GetAttributes returns parsed title attributes for the page.
func (p *OCRPage) GetAttributes() *TitleAttributes {
	return ParseTitleAttributes(p.Title)
}

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}

	// common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run reconstruct_pdf_from_hocr.go input.pdf\n")
		os.Exit(1)
	}

	// Load images from the PDF.
	images, err := loadImages(os.Args[1])
	if err != nil {
		fmt.Printf("Error loading images: %v\n", err)
		os.Exit(1)
	}

	outDir := "output"
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err := os.Mkdir(outDir, 0755)
		if err != nil {
			fmt.Printf("Error creating output directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Process each image with OCR.
	for pageNum, imgList := range images {
		fmt.Printf("Processing images on page %d\n", pageNum+1)
		for _, img := range imgList {
			ocrPage, err := processImage(img)
			if err != nil {
				fmt.Printf("Error processing image on page %d: %s\n", pageNum+1, err)
				continue
			}

			// Successfully processed image, ocrPage contains the parsed data
			writeContentAsPDF(ocrPage, fmt.Sprintf("output/page_%d.pdf", pageNum+1))
		}
	}
}

// loadImages loads images from the specified PDF file.
func loadImages(inputPath string) ([][]image.Image, error) {
	result := make([][]image.Image, 0)

	// Load images from the PDF and return an error if any occurs.
	pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, nil)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return nil, err
	}

	fmt.Print("Loading images from PDF document")

	totalImages := 0
	for i := 0; i < numPages; i++ {
		page, err := pdfReader.GetPage(i + 1)
		if err != nil {
			return nil, err
		}

		pextract, err := extractor.New(page)
		if err != nil {
			return nil, err
		}

		pimages, err := pextract.ExtractPageImages(nil)
		if err != nil {
			return nil, err
		}

		result = append(result, make([]image.Image, 0))
		for _, img := range pimages.Images {
			goImg, err := img.Image.ToGoImage()
			if err != nil {
				return nil, err
			}

			rotatedImg := imaging.Rotate(goImg, 360-float64(*page.Rotate), color.Transparent)

			result[i] = append(result[i], rotatedImg)
		}

		totalImages += len(pimages.Images)

		fmt.Print(".")
	}

	fmt.Println(" Done")

	fmt.Printf("Total: %d images\n", totalImages)

	return result, nil
}

// processImage sends the image to the OCR service and processes the HOCR response.
func processImage(img image.Image) (*OCRPage, error) {
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		return nil, err
	}

	imgReader := bytes.NewReader(buf.Bytes())

	// Configure OCR service options.
	opts := ocr.OCROptions{
		Url:           "http://localhost:8080/file",
		Method:        "POST",
		FileFieldName: "file",
		Headers: map[string]string{
			"Accept": "application/json",
		},
		FormFields: map[string]string{
			"format": "hocr",
		},
		TimeoutSeconds: 30,
	}

	// Create OCR client.
	client := ocr.NewHTTPOCRService(opts)

	result, err := client.ExtractText(context.Background(), imgReader, "image.jpg")
	if err != nil {
		return nil, fmt.Errorf("error extracting text: %w", err)
	}

	// Parse JSON response to extract the "result" field.
	var jsonObj map[string]interface{}
	if err := json.Unmarshal(result, &jsonObj); err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %w", err)
	}

	content, ok := jsonObj["result"].(string)
	if !ok {
		return nil, fmt.Errorf("result field is not a string")
	}

	// Parse hOCR HTML content
	var ocrPage OCRPage
	if err := xml.Unmarshal([]byte(content), &ocrPage); err != nil {
		return nil, fmt.Errorf("error unmarshalling HOCR data: %w", err)
	}

	return &ocrPage, nil
}

func writeContentAsPDF(p *OCRPage, outputPath string) error {
	if p.CAreas == nil || len(p.CAreas) == 0 {
		fmt.Println("CAreas empty")
		return nil
	}

	pt := p.GetAttributes()

	c := creator.New()
	c.SetPageSize(creator.PageSize{float64(pt.BBox.X1), float64(pt.BBox.Y1)})
	c.NewPage()

	for _, a := range p.CAreas {
		// Process each column area in the page
		at := a.GetAttributes()

		adiv := c.NewDivision()
		adiv.SetMargins(float64(at.BBox.X0)-float64(pt.BBox.X0),
			float64(pt.BBox.Y1)-float64(at.BBox.Y1),
			float64(pt.BBox.X1)-float64(at.BBox.X1),
			float64(at.BBox.Y0)-float64(pt.BBox.Y0),
		)

		for _, par := range a.Pars {
			// Process each paragraph in the area
			for _, l := range par.Lines {
				// Process each line in the paragraph
				lt := l.GetAttributes()

				for _, w := range l.Words {
					// Process each word in the line
					wt := w.GetAttributes()

					sp := c.NewStyledParagraph()
					sp.SetPos(float64(wt.BBox.X0), float64(wt.BBox.Y1))
					sp.SetMargins(0, float64(pt.BBox.X1)-float64(wt.BBox.X1), 0, float64(pt.BBox.Y1)-float64(wt.BBox.Y1))
					sp.SetFontSize(lt.XSize)
					sp.SetText(w.GetText())

					adiv.Add(sp)
				}
			}
		}

		c.Draw(adiv)
	}

	err := c.WriteToFile(outputPath)
	if err != nil {
		return err
	}

	fmt.Println("Saved reconstructed PDF to", outputPath)

	return nil
}
