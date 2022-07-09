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

	// Extract metadata if is already defined within given catalog.
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

	mm, ok := xmpDoc.GetMediaManagement()
	if !ok {
		mm = &xmputil.MediaManagement{
			// OriginalDocumentID is a persistent identifier of a document. It should persist no matter
			// what modification had been done on the document.
			// If the Media Management metadata is not defined within XMP document, this value would either be automatically
			// generated or set up to the one provided within MediaManagementOptions.
			// By setting this value here and copying it to the MediaManagementOptions we can control how this value
			// persists.
			OriginalDocumentID: "56119f84-a812-484a-bb4c-61c7e7cb3265",
		}
	}

	mmOptions := &xmputil.MediaManagementOptions{
		// OriginalDocumentID should maintain after any modification of provided document.
		OriginalDocumentID: string(mm.OriginalDocumentID),
		// Set this value if we want to create a new file (not overwrite current file).
		NewDocumentID: true,
		ModifyComment: "Added Media Management XMP Metadata",
		ModifyDate:    time.Now(),
		Modifier:      "Example User Modifier name",
	}
	if err = xmpDoc.SetMediaManagement(mmOptions); err != nil {
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
