/*
 * PDF to text: Extract all text for each page of a pdf file.
 *
 * Run as: go run pdf_extract_text.go testdata/*.pdf
 */

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"unicode"
)

const (
	usage                 = "Usage: go run normalize_text.go file1 file2 ...\n"
	defaultNormalizeWidth = 60
)

func main() {
	makeUsage(usage)

	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	for i, inFile := range args {
		outFile := fmt.Sprintf("normal.%d", i)
		fmt.Printf("%d: %q -> %q\n", i, inFile, outFile)
		normaliseFile(inFile, outFile, 60)
	}

}

func normaliseFile(inFile, outFile string, width int) {
	text := fileToString(inFile)
	text = normalizeText(text, width)
	stringToFile(outFile, text)
}

// getReader returns a PdfReader and the number of pages for PDF file `inputPath`.
func fileToString(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(fmt.Errorf("Couldn't open %q. err=%v", filename, err))
	}
	return string(data)
}

func stringToFile(filename, s string) {
	data := []byte(s)
	err := ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		panic(fmt.Errorf("Couldn't open %q. err=%v", filename, err))
	}
}

// makeUsage updates flag.Usage to include usage message `msg`.
func makeUsage(msg string) {
	usage := flag.Usage
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, msg)
		usage()
	}
}

// normalizeText returns `text` with runs of spaces of any kind (spaces, tabs, line breaks, etc)
// reduced to a single space. `width` is the target line width.
func normalizeText(text string, width int) string {
	if width < 0 {
		width = defaultNormalizeWidth
	}
	return splitLines(reduceSpaces(text), width)
}

// reduceSpaces returns `text` with runs of spaces of any kind (spaces, tabs, line breaks, etc)
// reduced to a single space.
func reduceSpaces(text string) string {
	text = reSpace.ReplaceAllString(text, " ")
	return strings.Trim(text, " \t\n\r\v")
}

var reSpace = regexp.MustCompile(`(?m)\s+`)

// splitLines inserts line breaks in string `text`. `width` is the target line width.
func splitLines(text string, width int) string {
	runes := []rune(text)
	if len(runes) < 2 {
		return text
	}
	lines := []string{}
	chars := []rune{}
	for i := 0; i < len(runes)-1; i++ {
		r, r1 := runes[i], runes[i+1]
		chars = append(chars, r)
		if (len(chars) >= width && unicode.IsSpace(r)) || (r == '.' && unicode.IsSpace(r1)) {
			lines = append(lines, string(chars))
			chars = []rune{}
		}
	}
	chars = append(chars, runes[len(runes)-1])
	if len(chars) > 0 {
		lines = append(lines, string(chars))
	}
	for i, ln := range lines {
		lines[i] = strings.Trim(ln, " \t\n\r\v")
	}
	return strings.Join(lines, "\n")
}
