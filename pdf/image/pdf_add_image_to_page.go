/*
 * Insert an image into a specific page, location in a PDF file.
 * If unsure about position, try getting the dimensions of a PDF with pdf_page_info.go first.
 *
 * Add image to a specific page of a PDF.  xPos and yPos define the lower left corner of the image location, and iwidth
 * is the width of the image in PDF coordinates (height/width ratio is maintained).
 *
 * Example go run pdf_add_image_to_page.go /tmp/input.pdf/tmp/output.pdf 1 /tmp/image.jpg 0 0 100
 * adds the image to the lower left corner of the page (0,0).  The width is 100 (typical page width 612).
 *
 * Run as: go run pdf_add_image_to_page.go input.pdf output.pdf page image.jpg xpos ypos width
 */

package main

import (
	"fmt"
	"os"
	"strconv"

	unicommon "github.com/unidoc/unidoc/common"
	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func main() {
	if len(os.Args) < 8 {
		fmt.Printf("Usage: go run pdf_add_image_to_page.go input.pdf output.pdf page image.jpg xpos ypos width\n")
		os.Exit(1)
	}

	// Use debug logging.
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))
	//unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelTrace))

	inputPath := os.Args[1]
	outputPath := os.Args[2]
	pageNumStr := os.Args[3]
	imagePath := os.Args[4]

	xPos, err := strconv.ParseFloat(os.Args[5], 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	yPos, err := strconv.ParseFloat(os.Args[6], 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	iwidth, err := strconv.ParseFloat(os.Args[7], 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("xPos: %d, yPos: %d\n", xPos, yPos)
	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err = addImageToPdf(inputPath, outputPath, imagePath, pageNum, xPos, yPos, iwidth)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

// Add image to a specific page of a PDF.  xPos and yPos define the lower left corner of the image location, and iwidth
// is the width of the image in PDF coordinates (height/width ratio is maintained).
func addImageToPdf(inputPath string, outputPath string, imagePath string, pageNum int, xPos float64, yPos float64, iwidth float64) error {
	// Open the image file.
	imgReader, err := os.Open(imagePath)
	if err != nil {
		unicommon.Log.Error("Error opening file: %s", err)
		return err
	}
	defer imgReader.Close()
	// Load the image with default handler.
	img, err := pdf.ImageHandling.Read(imgReader)
	if err != nil {
		unicommon.Log.Error("Error loading image: %s", err)
		return err
	}

	/*
		// When dealing with transparent images, sometimes need the following trick to make all the pixels
		// transparent:
		img.AlphaMap(func(alpha byte) byte {
			if alpha > 50 {
				return alpha
			} else {
				return 0 // Transparent
			}
		})
	*/

	// Read the input pdf file.
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}
	if pageNum <= 0 || pageNum > numPages {
		return fmt.Errorf("Page number out of range (%d/%d)", pageNum, numPages)
	}

	// Load the pages.
	pages := []*pdf.PdfPage{}
	for i := 0; i < numPages; i++ {
		page, err := pdfReader.GetPage(i + 1)
		if err != nil {
			return err
		}

		pages = append(pages, page)
	}

	// Add the image to the selected page.
	selPage := pages[pageNum-1]

	fmt.Printf("Page: %+v\n", selPage)
	mediabox, err := selPage.GetMediaBox()
	if err != nil {
		return err
	}
	fmt.Printf("Page mediabox: %+v\n", mediabox)

	// Find a free name for the image.
	num := 1
	imgName := pdfcore.PdfObjectName(fmt.Sprintf("Img%d", num))
	for selPage.Resources.HasXObjectByName(imgName) {
		num++
		imgName = pdfcore.PdfObjectName(fmt.Sprintf("Img%d", num))
	}

	encoder := pdfcore.NewFlateEncoder()

	// Create the XObject image.
	ximg, err := pdf.NewXObjectImageFromImage(img, nil, encoder)
	if err != nil {
		unicommon.Log.Error("Failed to create xobject image: %s", err)
		return err
	}

	// Add to the page resources.
	err = selPage.AddImageResource(imgName, ximg)
	if err != nil {
		return err
	}

	// Find an available GS name.
	i := 0
	gsName := pdfcore.PdfObjectName(fmt.Sprintf("GS%d", i))
	for selPage.HasExtGState(gsName) {
		i++
		gsName = pdfcore.PdfObjectName(fmt.Sprintf("GS%d", i))
	}

	// Create a normal graphics state.
	gs0 := pdfcore.MakeDict()
	gs0.Set("BM", pdfcore.MakeName("Normal"))
	selPage.AddExtGState(gsName, gs0)

	imWidth := iwidth
	imHeight := float64(img.Height) / float64(img.Width) * iwidth

	// Create content stream to add to the page contents.
	creator := pdfcontent.NewContentCreator()
	creator.
		Add_q().                               // Wrap.
		Add_gs(gsName).                        // Set graphics state.
		Add_cm(1, 0, 0, 1, xPos, yPos).        // Set position.
		Add_cm(imWidth, 0, 0, imHeight, 0, 0). // Scale to desired image size.
		Add_Do(imgName).                       // Draw the image.
		Add_Q()                                // Unwrap (go to prior graphics stack).

	selPage.AddContentStreamByString(creator.String())

	// Write output.
	pdfWriter := pdf.NewPdfWriter()
	for _, page := range pages {
		err = pdfWriter.AddPage(page)
		if err != nil {
			unicommon.Log.Error("Failed to add page: %s", err)
			return err
		}
	}

	fWrite, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer fWrite.Close()

	err = pdfWriter.Write(fWrite)
	if err != nil {
		return err
	}

	return nil
}
