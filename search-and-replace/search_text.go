package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

func main() {
	// Input parameters
	filePath := "./test-data/file1.pdf" // Path to the PDF file
	pattern := "Australia"              // Text pattern to search for
	pages := []int{1}                   // Page numbers to search on

	// Create a new PDF reader
	reader, _, err := model.NewPdfReaderFromFile(filePath, nil)
	if err != nil {
		fmt.Printf("Failed to create PDF reader: %v\n", err)
		os.Exit(1)
	}

	// Create an Editor object for searching
	editor := extractor.NewEditor(reader)

	// Perform the search for the specified pattern on the given pages
	matchesPerPage, err := editor.Search(pattern, pages)
	if err != nil {
		fmt.Printf("Failed to search pattern: %v\n", err)
		os.Exit(1)
	}

	// Print formatted search results
	printSearchResults(matchesPerPage, pages, pattern)
}

// printSearchResults formats and prints the search results.
// It displays indexes as [beg:end] and locations as {Llx Lly Urx Ury}.
// If no matches are found for a page, it prints a not found message.
func printSearchResults(matchesPerPage map[int]extractor.Match, pages []int, pattern string) {
	foundAny := false // Flag to check if any match is found across all pages

	for _, page := range pages {
		result, exists := matchesPerPage[page]
		if exists && len(result.Indexes) > 0 {
			foundAny = true
			fmt.Printf("Page %d:\n", page)

			// Prepare index strings
			var indexStrings []string
			for _, idx := range result.Indexes {
				indexStrings = append(indexStrings, fmt.Sprintf("[%d:%d]", idx[0], idx[1]))
			}
			fmt.Printf("indexes: %s\n", strings.Join(indexStrings, ", "))

			// Prepare location strings
			var locationStrings []string
			for _, box := range result.Locations {
				locationStrings = append(locationStrings, fmt.Sprintf("{%.2f %.2f %.2f %.2f}", box.BBox.Llx, box.BBox.Lly, box.BBox.Urx, box.BBox.Ury))
			}
			fmt.Printf("locations: %s\n\n", strings.Join(locationStrings, ", "))
		} else {
			// If no matches found for the current page
			fmt.Printf("Page %d:\n", page)
			fmt.Println("pattern didn't match any text\n")
		}
	}

	if !foundAny {
		// If no matches found in any of the pages
		fmt.Println("pattern didn't match any text in the specified pages.")
	}
}
