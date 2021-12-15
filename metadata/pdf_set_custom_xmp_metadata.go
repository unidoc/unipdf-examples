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

	xmprights "github.com/trimmer-io/go-xmp/models/xmp_rights"
	"github.com/trimmer-io/go-xmp/xmp"
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
		fmt.Printf( "Usage: %s INPUT_PDF_PATH OUTPUT_PDF_PATH\n", os.Args[0])
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

	// Extract XMP metadata from the PDF Catalog Metadata.
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
		// Otherwise, simply create a new XMP document,
		xmpDoc = xmputil.NewDocument()
	}

	// Unwrap github.com/trimmer-io/go-xmp/xmp.Document implementation on which base xmputil is implemented.
	goXmpDoc := xmpDoc.GetGoXmpDocument()

	// Getting direct access to go-xmp/xmp.Document allows extracting custom or undefined model for the XMP Metadata.
	// Multiple XMP Metadata models could be find in: https://github.com/trimmer-io/go-xmp/tree/master/models.
	xmpRightsModel, err := xmprights.MakeModel(goXmpDoc)
	if err != nil {
		log.Fatalf("Err: %v\n", err)
	}

	// Set up the fields in the xmp model.
	xmpRightsModel.Certificate = "56c69b3eaf10cff5c3bd8932f0169b"
	xmpRightsModel.UsageTerms = xmp.NewAltString("My Custom Usage Terms")
	xmpRightsModel.Owner = xmp.NewStringArray("Custom Owner")
	xmpRightsModel.WebStatement = "My Custom Web Statement"

	// Sync the model with the XMP Document.
	if err = xmpRightsModel.SyncToXMP(goXmpDoc); err != nil {
		log.Fatalf("Err: %v\n", err)
	}

	// The xmputil.Document is a wrapper over go-xmp/xmp.Document and will keep all the changes done on top of it.
	// We can safely marshal and store Metadata in the Pdf document.
	data, err := xmpDoc.MarshalIndent("", "\t")
	if err != nil {
		log.Fatalf("Err: %v\n", err)
	}

	// Create a simple PdfStream from the metadata bytes.
	metadataStream, err := core.MakeStream(data, nil)
	if err != nil {
		log.Fatalf("Err: %v\n", err)
	}

	// Set up catalog metadata as PdfStream.
	if err = pdfWriter.SetCatalogMetadata(metadataStream); err != nil {
		log.Fatalf("Err: %v\n", err)
	}

	// Create output file.
	err = pdfWriter.WriteToFile(outputPath)
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
}
