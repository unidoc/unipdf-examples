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

const licenseKey = `
-----BEGIN UNIDOC LICENSE KEY-----
Free trial license keys are available at: https://unidoc.io/
-----END UNIDOC LICENSE KEY-----
`

const usage = "Usage: %s INPUT_PDF_PATH\n"

func init() {
	err := license.SetLicenseKey(licenseKey, `Company Name`)
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

	// Extract XMP MediaManagement Namespace model, which contains information about document history and identifiers.
	mm, ok := xmpDoc.GetMediaManagement()
	if !ok {
		fmt.Println("No XMP Media Management namespace defined within XMP document.")
		return
	}

	// Print out the media management.
	fmt.Println("Document MediaManagement: ")

	// OriginalDocumentID must be created to identify a new document. This identifies a document as a conceptual entity.
	fmt.Printf("OriginalDocumentID: %v\n", mm.OriginalDocumentID)

	// DocumentID when a document is copied to a new file path or converted to a new format with
	// Save As, another new document ID should usually be assigned. This identifies a general version or
	// branch of a document. You can use it to track different versions or extracted portions of a document
	// with the same original-document ID.
	fmt.Printf("DocumentID: %v\n", mm.DocumentID)

	// InstanceID to track a document’s editing history, you must assign a new instance ID
	// whenever a document is saved after any changes. This uniquely identifies an exact version of a
	// document. It is used in resource references (to identify both the document or part itself and the
	// referenced or referencing documents), and in document-history resource events (to identify the
	// document instance that resulted from the change).
	fmt.Printf("InstanceID: %v\n", mm.InstanceID)

	// DerivedFrom references the source document from which this one is derived,
	// typically through a Save As operation that changes the file name or format. It is a minimal reference;
	// missing components can be assumed to be unchanged. For example, a new version might only need
	// to specify the instance ID and version number of the previous version, or a rendition might only need
	// to specify the instance ID and rendition class of the original.
	if mm.DerivedFrom != nil {
		fmt.Printf("DerivedFrom: %v\n", mm.DerivedFrom)
	}
	// VersionID are meant to associate the document with a product version that is part of a release process. They can be useful in tracking the
	// document history, but should not be used to identify a document uniquely in any context.
	// Usually it simply works by incrementing integers 1,2,3...
	if mm.VersionID != "" {
		fmt.Printf("VersionID: %v\n", mm.VersionID)
	}
	// Versions is the history of the document versions along with the comments, timestamps and issuers.
	if len(mm.Versions) > 0 {
		fmt.Println("Versions:")
		for _, v := range mm.Versions {
			fmt.Printf("\t%+v,\n", v)
		}
	}
}
