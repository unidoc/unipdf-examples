/*
 * List all content streams for all pages in a pdf file.
 *
 * Run as: go run pdf_print_content_streams.go input.pdf [page]
 * The page number is optional (by default all pages are processed).
 */

package main

import (
	"fmt"
	"os"

	"strconv"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/contentstream"
	"github.com/unidoc/unipdf/v3/model"
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
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run pdf_list_content_streams.go input.pdf [page]\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	pageNum := -1

	if len(os.Args) >= 3 {
		val, err := strconv.ParseInt(os.Args[2], 10, 64)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		pageNum = int(val)
	}

	fmt.Println(inputPath)
	err := listContentStreams(inputPath, pageNum)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func listContentStreams(inputPath string, targetPageNum int) error {
	pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, nil)
	if err != nil {
		return err
	}
	defer f.Close()

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	fmt.Printf("--------------------\n")
	fmt.Printf("Content streams:\n")
	fmt.Printf("--------------------\n")
	for i := 0; i < numPages; i++ {
		pageNum := i + 1
		if pageNum != targetPageNum && targetPageNum != -1 {
			continue
		}
		fmt.Printf("Page %d\n", pageNum)

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return err
		}

		contentStreams, err := page.GetContentStreams()
		if err != nil {
			return err
		}
		fmt.Printf("Page %d has %d content streams:\n", pageNum, len(contentStreams))

		pageContentStr := ""

		// If the value is an array, the effect shall be as if all of the streams in the array were concatenated,
		// in order, to form a single stream.
		for _, cstream := range contentStreams {
			pageContentStr += cstream
		}
		fmt.Printf("%s\n", pageContentStr)

		cstreamParser := contentstream.NewContentStreamParser(pageContentStr)
		operations, err := cstreamParser.Parse()
		if err != nil {
			return err
		}

		fmt.Printf("=== Full list\n")
		for idx, op := range *operations {
			fmt.Printf("Operation %d: %s - Params: %v\n", idx+1, op.Operand, op.Params)
		}
	}

	return nil
}
