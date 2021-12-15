package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/model/xmputil"
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
	if len(args) < 2 {
		fmt.Printf(usage, os.Args[0])
		return
	}
	inputPath := args[1]

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

	// Extract XMP metadata from the PDF Catalog Metadata.
	metadata, ok := reader.GetCatalogMetadata()
	if !ok {
		log.Fatalf("No XMP metadata defined within Pdf Document")
	}

	// Get metadata stream bytes.
	stream, ok := core.GetStream(metadata)
	if !ok {
		log.Fatalf("Catalog metadata is expected to be a stream but is: %T", metadata)
	}

	// Load up the XMP document implementation from the input stream.
	xmpDoc, err := xmputil.LoadDocument(stream.Stream)
	if err != nil {
		log.Fatalf("Reading XMP metadata failed: %v", err)
	}

	// Extract PDF Namespace model, which contains information about PDF Document.
	// It mostly should match Pdf document Info dictionary.
	pdfInfoMetadata, ok := xmpDoc.GetPdfInfo()
	if !ok {
		fmt.Println("No PdfInfo namespace defined within XMP document.")
		return
	}

	if pdfInfoMetadata.InfoDict != nil {
		// Generally this value should match PDF document Info dictionary.
		infoDict, err := model.NewPdfInfoFromObject(pdfInfoMetadata.InfoDict)
		if err != nil {
			log.Fatalf("Err: %v", err)
		}

		fmt.Printf("InfoDict: %#v\n", infoDict)
	}

	if pdfInfoMetadata.PdfVersion != "" {
		fmt.Printf("Pdf Version: %s\n", pdfInfoMetadata.PdfVersion)
	}

	if pdfInfoMetadata.Copyright != "" {
		fmt.Printf("Copyright: %s\n", pdfInfoMetadata.Copyright)
	}
	fmt.Printf("Marked: %v\n", pdfInfoMetadata.Marked)
}
