/*
 * PDF optimization (compression) example.
 *
 * Run as: go run pdf_apply_standard.go <input.pdf> <output.pdf>
 */

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/unidoc/unipdf/v4/common/license"
	"github.com/unidoc/unipdf/v4/core"
	"github.com/unidoc/unipdf/v4/model"
	"github.com/unidoc/unipdf/v4/model/xmputil"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}
}

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Printf("Usage: %s INPUT_PDF_PATH OUTPUT_PDF_PATH", os.Args[0])
		return
	}
	inputPath := args[1]
	outputPath := args[2]

	// Initialize starting time.
	start := time.Now()
	defer func() {
		duration := float64(time.Since(start)) / float64(time.Millisecond)
		fmt.Printf("Processing time: %.2f ms\n", duration)
	}()

	// Read some file to which you want to add XMP metadata.
	reader, file, err := model.NewPdfReaderFromFile(inputPath, nil)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
	defer file.Close()

	// Copy content of the reader into a writer.
	pdfWriter, err := reader.ToWriter(nil)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	// Check if the metadata is already defined within given catalog.
	var xmpDoc *xmputil.Document
	metadata, ok := reader.GetCatalogMetadata()
	if ok {
		stream, ok := core.GetStream(metadata)
		if !ok {
			log.Fatalf("Catalog metadata is expected to be a stream but is: %T", metadata)
		}

		xmpDoc, err = xmputil.LoadDocument(stream.Stream)
		if err != nil {
			log.Fatalf("Reading XMP metadata failed: %v", err)
		}
	} else {
		xmpDoc = xmputil.NewDocument()
	}

	// Read PdfInfo from the origin file.
	pdfInfo, err := reader.GetPdfInfo()
	if err != nil {
		log.Fatalf("Err: %v", err)
	}

	// Set up some custom fields in the document PdfInfo if it doesn't exist.
	var createdAt time.Time
	if pdfInfo.CreationDate != nil {
		createdAt = pdfInfo.CreationDate.ToGoTime()
	} else {
		createdAt = time.Now()
		creationDate, err := model.NewPdfDateFromTime(createdAt)
		if err != nil {
			log.Fatalf("Err: %v", err)
		}
		pdfInfo.CreationDate = &creationDate
	}

	modifiedAt := time.Now()
	pdfInfo.Author = core.MakeString("Example Author")
	modDate, err := model.NewPdfDateFromTime(modifiedAt)
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	pdfInfo.ModifiedDate = &modDate

	// Copy the content of the PdfInfo into XMP metadata.
	xmpPdfMetadata := &xmputil.PdfInfoOptions{
		InfoDict:   pdfInfo.ToPdfObject(),
		PdfVersion: reader.PdfVersion().String(),
		Copyright:  "Copyright Example",
		// Setting this value would clear current InfoDict within XMP Document (if exists) and overwrite by provided pdfInfo.
		// In some cases when the origin XMP Metadata contained different fields than the PDF Document, setting this value
		// to false would maintain values undefined within provided pdfInfo.
		Overwrite: true,
	}
	// Store PDF Metadata into xmp document.
	err = xmpDoc.SetPdfInfo(xmpPdfMetadata)
	if err != nil {
		log.Fatalf("Err: %v", err)
	}

	// Once we've defined all the XMP metadata we wanted we can extract raw bytes stream and store as catalog metadata PdfStream.
	// By doing:

	// 1. Marshal XMP Document into raw bytes.
	metadataBytes, err := xmpDoc.MarshalIndent("", "\t")
	if err != nil {
		log.Fatalf("Err: %v", err)
	}

	// 2. Create new PdfStream
	metadataStream, err := core.MakeStream(metadataBytes, nil)
	if err != nil {
		log.Fatalf("Err: %v", err)
	}

	// 3. Set the metadata stream as catalog metadata.
	if err = pdfWriter.SetCatalogMetadata(metadataStream); err != nil {
		log.Fatalf("Err: %v", err)
	}

	// Create output file.
	err = pdfWriter.WriteToFile(outputPath)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
}
