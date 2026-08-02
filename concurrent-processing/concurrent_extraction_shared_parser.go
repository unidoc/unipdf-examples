/*
 * This example demonstrates page level concurrent text extraction backed by a
 * single shared parser.
 *
 * The parser and reader I/O layer takes an io.ReaderAt, which makes the parser
 * internally thread-safe. model.NewPdfReaderFromParser resolves the document
 * structure eagerly, so the resulting reader can be used from several goroutines
 * without further synchronization: load the document once, then read pages in
 * parallel.
 *
 * NOTE: readers built with NewPdfReaderFromParser do not support signature
 * verification or encrypted documents, and cannot be used with model.NewPdfAppender.
 * For those cases use model.NewPdfReader as usual.
 *
 * Run as: go run concurrent_extraction_shared_parser.go <input.pdf> <output_dir>
 */

package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/unidoc/unipdf/v5/common/license"
	"github.com/unidoc/unipdf/v5/core"
	"github.com/unidoc/unipdf/v5/extractor"
	"github.com/unidoc/unipdf/v5/model"
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
	if len(os.Args) < 3 {
		fmt.Printf("Usage: go run concurrent_extraction_shared_parser.go input.pdf output_dir\n")
		os.Exit(1)
	}
	inputPath := os.Args[1]
	outputDir := os.Args[2]

	start := time.Now()

	if err := extractConcurrently(inputPath, outputDir); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("time taken for concurrent extraction: %v\n", time.Since(start))
}

// extractConcurrently extracts the text of every page in its own goroutine,
// with all goroutines sharing one parser, and writes each page to its own file.
func extractConcurrently(inputPath, outputDir string) error {
	if err := os.MkdirAll(outputDir, fs.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", outputDir, err)
	}

	// Buffer the document so the parser is backed by an io.ReaderAt that can
	// serve many goroutines at once. A *os.File works here too.
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	parser, err := core.NewParser(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create parser: %w", err)
	}

	// Building the reader from the parser resolves the document structure up
	// front, which is what makes the concurrent page access below safe.
	reader, err := model.NewPdfReaderFromParser(parser)
	if err != nil {
		return fmt.Errorf("failed to create reader: %w", err)
	}

	numPages, err := reader.GetNumPages()
	if err != nil {
		return err
	}

	// Collect per page results in a pre-sized slice so no locking is needed:
	// each goroutine owns exactly one element.
	texts := make([]string, numPages+1)
	errs := make([]error, numPages+1)

	var wg sync.WaitGroup
	for i := 1; i <= numPages; i++ {
		wg.Add(1)
		go func(pageNum int) {
			defer wg.Done()
			texts[pageNum], errs[pageNum] = extractPage(reader, pageNum)
		}(i)
	}
	wg.Wait()

	for pageNum := 1; pageNum <= numPages; pageNum++ {
		if errs[pageNum] != nil {
			fmt.Printf("page %d: %v\n", pageNum, errs[pageNum])
			continue
		}

		filePath := filepath.Join(outputDir, strconv.Itoa(pageNum)+".txt")
		if err := os.WriteFile(filePath, []byte(texts[pageNum]), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", filePath, err)
		}
	}

	fmt.Printf("extracted %d pages to %s\n", numPages, outputDir)

	return nil
}

// extractPage extracts the text of a single page. It is called concurrently for
// every page, sharing one reader.
func extractPage(reader *model.PdfReader, pageNum int) (string, error) {
	page, err := reader.GetPage(pageNum)
	if err != nil {
		return "", err
	}

	ex, err := extractor.New(page)
	if err != nil {
		return "", err
	}

	return ex.ExtractText()
}
