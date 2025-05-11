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

	xmprights "github.com/trimmer-io/go-xmp/models/xmp_rights"
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
		fmt.Printf("Usage: %s INPUT_PDF_PATH\n", os.Args[0])
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

	// Unwrap github.com/trimmer-io/go-xmp/xmp.Document implementation on which base xmputil is implemented.
	goXmpDoc := xmpDoc.GetGoXmpDocument()

	// Getting direct access to go-xmp/xmp.Document allows extracting custom or undefined model for the XMP Metadata.
	// Multiple XMP Metadata models could be find in: https://github.com/trimmer-io/go-xmp/tree/master/models.
	xmpRightsModel := xmprights.FindModel(goXmpDoc)
	if xmpRightsModel == nil {
		fmt.Println("No XMP Media Management namespace defined within XMP document.")
		return
	}

	if xmpRightsModel.Certificate != "" {
		fmt.Printf("Certificate: %v\n", xmpRightsModel.Certificate)
	}
	if !xmpRightsModel.Owner.IsZero() {
		fmt.Println("Owners: ")
		for _, owner := range xmpRightsModel.Owner {
			fmt.Printf("- %v\n", owner)
		}
	}
	if !xmpRightsModel.UsageTerms.IsZero() {
		fmt.Println("Usage Terms: ")
		for _, usageTerms := range xmpRightsModel.UsageTerms {
			fmt.Printf("Lang: %s, Terms: %s\n", usageTerms.Lang, usageTerms.Value)
		}
	}
	if xmpRightsModel.WebStatement != "" {
		fmt.Printf("Web Statement: %v\n", xmpRightsModel.WebStatement)
	}
}
