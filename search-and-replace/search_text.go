/*
 * This example code shows how to do text searching on pdf using unipdf
 *
 * Run as: go run search_text.go <pattern> <pages> <input>
 *
 * Example: go run search_text.go "copyright law" "1,2" ./test-data/file1.pdf
 */

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/extractor"
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
	// Ensure enough arguments are provided
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run main.go <pattern> <pages> <input>")
		os.Exit(1)
	}

	// Parse positional arguments
	pattern := os.Args[1]
	pagesArg := os.Args[2]
	filePath := os.Args[3]

	// Convert pages string to a slice of integers
	pageStrings := strings.Split(pagesArg, ",")
	pageList := []int{}
	for _, pageStr := range pageStrings {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			fmt.Printf("Invalid page number: %s\n", pageStr)
			os.Exit(1)
		}
		pageList = append(pageList, page)
	}

	// Create a new PDF reader
	reader, _, err := model.NewPdfReaderFromFile(filePath, nil)
	if err != nil {
		fmt.Printf("Failed to create PDF reader: %v\n", err)
		os.Exit(1)
	}

	// Create an Editor object for searching
	editor := extractor.NewEditor(reader)

	// Perform the search for the specified pattern on the given pages
	matchesPerPage, err := editor.Search(pattern, pageList)
	if err != nil {
		fmt.Printf("Failed to search pattern: %v\n", err)
		os.Exit(1)
	}

	// Print formatted search results
	printSearchResults(matchesPerPage, pageList, pattern)
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
