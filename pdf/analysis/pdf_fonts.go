/*
 * Shows all fonts in a PDF file.
 *
 * Run as: go run pdf_fonts.go input.pdf
 */

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Syntax: go run pdf_fonts.go input.pdf")
		os.Exit(1)
	}

	// Enable debug-level logging.
	//unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	// filename := os.Args[1]
	// showOneFile(filename)

	pathList, err := patternsToPaths(os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	numFiles, occurrences, uniques, err := fontsInPdfList(pathList)
	if err != nil {
		os.Exit(1)
	}

	fmt.Printf("%d files tested (of %d tried)\n", numFiles, len(pathList))
	fmt.Printf("%d font occurences (not interesting)\n", len(occurrences))
	fmt.Printf("%d font subtype occurrences\n", len(uniques))

	subtypeCounts := uniques.typeCounts()
	subtypes := mapKeys(subtypeCounts)

	fmt.Printf("%d subtypes\n", len(subtypes))
	for i, subtype := range subtypes {
		count := subtypeCounts[subtype]
		percentSubtypes := float64(count) / float64(len(uniques)) * 100.0
		percentFiles := float64(count) / float64(numFiles) * 100.0
		fmt.Printf("%4d: %#-14q %4d files (%5.1f%%) Occurred in %.1f%% of files.\n",
			i, subtype, count, percentSubtypes, percentFiles)
	}

}

type fontOccurence struct {
	subtype  string
	basefont string
	filename string
	version  string
}

type compendium []fontOccurence

func (occurrences compendium) typeCounts() map[string]int {
	counts := map[string]int{}
	for _, o := range occurrences {
		counts[o.subtype]++
	}
	return counts
}

func mapKeys(m map[string]int) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return
}

func fontsInPdfList(pathList []string) (numFiles int, occurrences, uniques compendium, err error) {
	for i, filename := range pathList {
		version := ""
		fonts := []pdf.PdfFont{}
		version, fonts, err = fontsInPdf(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %d of %d %q err=%v\n", i, len(pathList), filename, err)
			continue
		}
		numFiles++
		typeCount := map[string]int{}
		for _, font := range fonts {
			subtype := font.Subtype()
			basefont := font.BaseFont()
			o := fontOccurence{subtype, basefont, filename, version}
			occurrences = append(occurrences, o)
			typeCount[subtype]++
		}
		for subtype := range typeCount {
			o := fontOccurence{subtype, "", filename, version}
			uniques = append(uniques, o)
		}
	}
	return
}

// fontsInPdf returns a list of the fonts in PDF file `filename`.
func fontsInPdf(filename string) (version string, fonts []pdf.PdfFont, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return
	}
	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return
	}

	// Try decrypting with an empty one.
	if isEncrypted {
		auth := false
		auth, err = pdfReader.Decrypt([]byte(""))
		if err != nil {
			return
		}
		if !auth {
			err = errors.New("Unable to decrypt password protected file - need to specify pass to Decrypt")
			return
		}
	}

	version = pdfReader.PdfVersion()

	for _, objNum := range pdfReader.GetObjectNums() {
		var obj pdfcore.PdfObject
		obj, err = pdfReader.GetIndirectObjectByNumber(objNum)
		if err != nil {
			return
		}
		font, err := pdf.NewPdfFontFromPdfObject(obj)
		if err != nil {
			continue
		}
		fonts = append(fonts, *font)
	}

	return
}

func showOneFile(filename string) {
	fmt.Printf("Input file: %s\n", filename)
	version, fonts, err := fontsInPdf(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	sort.Slice(fonts, func(i, j int) bool {
		t1, t2 := fonts[i].Subtype(), fonts[j].Subtype()
		if t1 != t2 {
			return t1 < t2
		}
		return fonts[i].BaseFont() < fonts[j].BaseFont()
	})

	fmt.Printf("%-40q (version %s) %d fonts\n", filename, version, len(fonts))
	for i, font := range fonts {
		fmt.Printf("%4d: %-8s %#q\n", i, font.Subtype(), font.BaseFont())
	}
}

// patternsToPaths returns a list of files matching the patterns in `patternList`
func patternsToPaths(patternList []string) ([]string, error) {
	pathList := []string{}
	for _, pattern := range patternList {
		files, err := filepath.Glob(pattern)
		if err != nil {
			fmt.Fprintf(os.Stderr, "patternsToPaths: Glob failed. pattern=%#q err=%v\n", pattern, err)
			return pathList, err
		}
		for _, path := range files {
			if !regularFile(path) {
				fmt.Fprintf(os.Stderr, "Not a regular file: %q\n", path)
				continue
			}
			pathList = append(pathList, path)
		}
	}
	return pathList, nil
}

// regularFile returns true if file `path` is a regular file
func regularFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not stat %q. err=%v", path, err)
		return false
	}
	return fi.Mode().IsRegular()
}
