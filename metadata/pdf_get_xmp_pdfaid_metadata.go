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

	// Extract PDF/A ID Namespace model, which contains information about conformance of given document.
	// I.e.:
	// A document conformant with the PDF/A-1B would have its values defined as:
	// xmputil.PdfAID{
	//   Part: 1,
	//   Conformance: "B",
	// }
	pdfaId, ok := xmpDoc.GetPdfAID()
	if !ok {
		fmt.Println("No PDF/A ID namespace defined within XMP document.")
	} else {
		fmt.Printf("Document was marked as conformant with PDF/A-%d%s\n", pdfaId.Part, pdfaId.Conformance)
	}

}
