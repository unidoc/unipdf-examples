/*
 * Add bookmarks to pdf file.
 *
 * Run as: go run pdf_add_bookmarks.go input.pdf output.pdf
 */

package main

import (
	"fmt"
	"os"

	unicommon "github.com/unidoc/unidoc/common"
	unilicense "github.com/unidoc/unidoc/license"
	unipdf "github.com/unidoc/unidoc/pdf"
)

func initUniDoc(licenseKey string) error {
	if len(licenseKey) > 0 {
		err := unilicense.SetLicenseKey(licenseKey)
		if err != nil {
			return err
		}
	}

	// To make the library log we just have to initialise the logger which satisfies
	// the unicommon.Logger interface, unicommon.DummyLogger is the default and
	// does not do anything. Very easy to implement your own.
	// unicommon.SetLogger(unicommon.DummyLogger{})
	unicommon.SetLogger(unicommon.ConsoleLogger{})

	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Example adds bookmarks with page number referring to each page.\n")
		fmt.Printf("Usage: go run pdf_add_bookmarks.go input.pdf output.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	err := initUniDoc("")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err = addBookmarks(inputPath, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

func addBookmarks(inputPath string, outputPath string) error {
	pdfWriter := unipdf.NewPdfWriter()

	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	defer f.Close()

	pdfReader, err := unipdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return err
	}

	if isEncrypted {
		_, err = pdfReader.Decrypt([]byte(""))
		if err != nil {
			return err
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	outlineTree := unipdf.NewPdfOutlineTree()
	isFirst := true
	var node *unipdf.PdfOutlineItem

	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPageAsPdfPage(pageNum)
		if err != nil {
			return err
		}

		outputPage := page.GetPageAsIndirectObject()
		err = pdfWriter.AddPage(outputPage)
		if err != nil {
			return err
		}

		item := unipdf.NewOutlineBookmark(fmt.Sprintf("Page %d", i+1), outputPage)
		fmt.Printf("Item: %v\n", item)
		if isFirst {
			outlineTree.First = &item.PdfOutlineTreeNode
			node = item
			isFirst = false
		} else {
			node.Next = &item.PdfOutlineTreeNode
			item.Prev = &node.PdfOutlineTreeNode
			node = item
		}
	}
	outlineTree.Last = &node.PdfOutlineTreeNode
	pdfWriter.AddOutlineTree(&outlineTree.PdfOutlineTreeNode)

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

/*
pdfReader -> read a pdf file
pdfwriter -> create new
add all the pages

then generate the outlines...

p1, err := reader.GetPage(1)
p4, err := reader.GetPage(1)
p10, err := reader.GetPage(1)

// Make destination?

outlines := unipdf.NewOutlines()
ch1 := unipdf.NewOutline(p1, "Chapter 1")
subch := unipdf.NewOutline(p1, "Introduction")
outlines.AddOutline(&ch1)
ch1.AddOutline(&subch)
ch2 := unipdf.NewOutline(p4, "Chapter 2")
outlines.AddOutline(p10, ch2)
*/
