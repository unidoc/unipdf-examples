/*
 * Outputs information XML metadata of root catalog for PDF files.
 *
 * Note: Each component within a PDF can have associated metadata stream. This example showcases
 * how to retrieve the metadata for the root catalog (typically document information). Similar methodology
 * can be applied to XML metadata for inner components.
 *
 * Run as: go run pdf_metadata_get_xml.go input1.pdf [input2.pdf] ...
 */

package main

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"sort"

	unicommon "github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Outputs XML metadata for PDF files\n")
		fmt.Printf("Usage: go run pdf_metadata_get_xml.go input1.pdf [input2.pdf] ...\n")
		os.Exit(1)
	}

	// Enable debug-level logging.
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	fmt.Printf("XML metadata for root catalog\n")
	for _, inputPath := range os.Args[1:len(os.Args)] {
		fmt.Printf("Input file: %s\n", inputPath)

		err := printXMLMetadataForPdf(inputPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}
}

func printXMLMetadataForPdf(inputPath string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return err
	}
	/*
		numPages, err := pdfReader.GetNumPages()
		if err != nil {
			return err
		}
	*/

	// XXX/FIXME: Much of the clunky type casting and tracing is being improved in v3.
	catalogDict, err := getRootCatalog(pdfReader)
	if err != nil {
		return err
	}

	// Get metadata.
	metadataObj, err := resolve(pdfReader, catalogDict.Get("Metadata"))
	if err != nil {
		return err
	}
	metadataStream, has := metadataObj.(*core.PdfObjectStream)
	if !has {
		fmt.Printf("Metadata for root catalog not present")
		return nil
	}
	xmlMetadata, err := core.DecodeStream(metadataStream)
	if err != nil {
		return err
	}
	//fmt.Printf("Raw - XML Metadata: %s\n", xmlMetadata)

	xmlDecoder := xml.NewDecoder(bytes.NewReader(xmlMetadata))
	var xmp xmpMetadata
	err = xmlDecoder.Decode(&xmp)
	if err != nil {
		return err
	}

	// Print xmp metadata.
	xmp.print()

	return nil
}

// getRootCatalog returns the root catalog for the PDF.
// Note: a lot of the syntax here will be simplified in v3.
func getRootCatalog(r *model.PdfReader) (*core.PdfObjectDictionary, error) {
	// Trailer dictionary points to catalog.
	trailerDict, err := r.GetTrailer()
	if err != nil {
		return nil, err
	}
	if trailerDict == nil {
		return nil, errors.New("missing trailer dict")
	}

	// Get root catalog
	catalogObj, err := resolve(r, trailerDict.Get("Root"))
	if err != nil {
		return nil, err
	}
	catalogDict, has := core.TraceToDirectObject(catalogObj).(*core.PdfObjectDictionary)
	if !has {
		return nil, errors.New("catalog dict missing")
	}

	return catalogDict, nil
}

// resolve resolves references and returns only resolved objects within context of a specific reader `r`.
// Note: This function will not be needed in v3 which simplifies working with PdfObjects.
func resolve(r *model.PdfReader, obj core.PdfObject) (core.PdfObject, error) {
	switch t := obj.(type) {
	case *core.PdfObjectReference:
		ref := t
		var err error
		obj, err = r.GetIndirectObjectByNumber(int(ref.ObjectNumber))
		if err != nil {
			return nil, err
		}
	}
	return obj, nil
}

// xmpMetadata is a simple representation of XMP metadata. It can be used for unmarshalling of basic
// XMP field data under XMP -> Description.
// This serves for example use and can be extended for custom applications.
type xmpMetadata struct {
	Descriptions []struct {
		Tags []struct {
			XMLName xml.Name
			//Value xml.CharData `xml:",chardata"`
			Value xml.CharData `xml:",innerxml"`
		} `xml:",any"`
	} `xml:"RDF>Description"`
}

// keyValMap returns a key -> value map of Description data for XMP metadata.
// Mostly
func (xmp xmpMetadata) keyValMap() map[string]string {
	keyVal := map[string]string{}

	for _, desc := range xmp.Descriptions {
		for _, tag := range desc.Tags {
			keyVal[tag.XMLName.Local] = string(tag.Value)
		}

	}
	return keyVal
}

func (xmp xmpMetadata) print() {
	keyVal := xmp.keyValMap()

	keys := []string{}
	for key, _ := range keyVal {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	fmt.Printf("XMP RDF Description data:\n")
	for _, key := range keys {
		fmt.Printf("%s - %s\n", key, keyVal[key])
	}
}
