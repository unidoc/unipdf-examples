/**
 * This example demonstrates how to copy pages from an existing PDF file to a new PDF file
 * and preserve the structure tree information.
 *
 * Usage:
 * go run pdf_copy_page_with_accessibility.go INPUT_PDF_PATH OUTPUT_PDF_PATH SPACE_SEPARATED_LIST_OF_PAGE_NUMBER_TO_COPY
 */

package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/unidoc/unipdf/v4/common/license"
	"github.com/unidoc/unipdf/v4/core"
	"github.com/unidoc/unipdf/v4/creator"
	"github.com/unidoc/unipdf/v4/model"
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
		fmt.Printf("Usage: %s INPUT_PDF_PATH OUTPUT_PDF_PATH SPACE_SEPARATED_LIST_OF_PAGE_NUMBER_TO_COPY\n", os.Args[0])
		return
	}

	inputPath := args[1]
	outputPath := args[2]
	pageNumbers := args[3:]

	c := creator.New()
	c.SetLanguage("en-US")
	c.SetPdfWriterAccessFunc(func(w *model.PdfWriter) error {
		w.SetCatalogMarkInfo(core.MakeDictMap(map[string]core.PdfObject{
			"Marked": core.MakeBool(true),
		}))

		w.SetVersion(1, 5)

		return nil
	})

	// Set the viewer preferences.
	vp := model.NewViewerPreferences()
	vp.SetDisplayDocTitle(true)

	c.SetViewerPreferences(vp)

	// Set PDF title.
	model.SetPdfTitle("Output Sample PDF")

	// Construct the StructTreeRoot.
	newStr := model.NewStructTreeRoot()

	// Construct base K dictionary.
	docK := model.NewKDictionary()
	docK.S = core.MakeName(string(model.StructureTypeDocument))

	newStr.AddKDict(docK)

	readerOpt := model.ReaderOpts{
		ComplianceMode: true,
	}

	r, _, err := model.NewPdfReaderFromFile(inputPath, &readerOpt)
	if err != nil {
		fmt.Printf("Error loading input file: %v\n", err)
		return
	}

	// Get the original `StructTreeRoot` object.
	strObj, found := r.GetCatalogStructTreeRoot()
	if !found {
		fmt.Printf("No StructTreeRoot found in the input PDF\n")
		return
	}

	orgStr, err := model.NewStructTreeRootFromPdfObject(strObj)
	if err != nil {
		fmt.Printf("Error parsing StructTreeRoot object: %v\n", err)
		return
	}

	orgK := orgStr.K[0].GetChildren()

	newPageIdx := 0

	for _, i := range pageNumbers {
		n, err := strconv.Atoi(i)
		if err != nil {
			fmt.Printf("Invalid page number: %s\n", i)
			continue
		}

		orgPage, err := r.GetPage(n)
		if err != nil {
			fmt.Printf("Error getting page %d: %v\n", n, err)
			continue
		}

		dupPage := orgPage.Duplicate()
		dupPage.SetStructParentsKey(newPageIdx)
		c.AddPage(dupPage)

		// Add section K object to store all original K objects from template page.
		sectK := model.NewKDictionary()
		sectK.S = core.MakeName(string(model.StructureTypeSection))
		sectK.T = core.MakeString(fmt.Sprintf("Page %d", n))
		sectK.GenerateRandomID()

		sectKv := model.KValue{}
		sectKv.SetKDict(sectK)

		// Copy all K objects from the original page to the new page.
		for _, k := range orgK {
			copyK := deepCopyKObject(k)

			if copyK != nil {
				sectK.AddChild(copyK)
			}
		}

		setKPageNumber(&sectKv, int64(n))

		docK.AddChild(&sectKv)

		newPageIdx++
	}

	c.SetStructTreeRoot(newStr)

	err = c.WriteToFile(outputPath)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
	}
}

func setKPageNumber(kv *model.KValue, page int64) {
	if k := kv.GetKDict(); k != nil {
		k.SetPageNumber(page)

		for _, child := range k.GetChildren() {
			setKPageNumber(child, page)
		}
	}
}

func deepCopyKObject(origin *model.KValue) *model.KValue {
	copyKVal := model.NewKValue()

	if oKDict := origin.GetKDict(); oKDict != nil {
		copy := model.NewKDictionary()
		copy.S = oKDict.S
		copy.ID = oKDict.ID
		copy.Lang = oKDict.Lang
		copy.Alt = oKDict.Alt
		copy.T = oKDict.T
		copy.Pg = oKDict.Pg
		copy.C = oKDict.C

		for _, child := range oKDict.GetChildren() {
			copy.AddChild(deepCopyKObject(child))
		}

		copyKVal.SetKDict(copy)

		return copyKVal
	} else if refObj := origin.GetRefObject(); refObj != nil {
		copyKVal.SetRefObject(refObj)

		return copyKVal
	} else if mcid := origin.GetMCID(); mcid != nil {
		copyKVal.SetMCID(*mcid)

		return copyKVal
	}

	return nil
}
