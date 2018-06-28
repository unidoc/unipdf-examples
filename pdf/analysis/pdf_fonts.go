// Shows all fonts in a PDF file.
//
// Run as: go run pdf_fonts.go o testdata/*.pdf testdata/**/*.pdf
//

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
		fmt.Fprintln(os.Stderr, "Usage: go run pdf_fonts.go testdata/*.pdf")
		os.Exit(1)
	}

	// Enable debug-level logging.
	//unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	// showOneFile(os.Args[1])

	pathList, err := patternsToPaths(os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	_, occurrences, uniques, err := fontsInPdfList(pathList)
	if err != nil {
		os.Exit(1)
	}
	versionUniques := uniques.byVersion()
	versions := versionKeys(versionUniques)
	fmt.Printf("%d total font occurrences\n", len(occurrences))

	fmt.Printf("%4d files total.\n", uniques.numFiles())
	for _, version := range versions {
		fmt.Printf("%4d files PDF version %s\n", versionUniques[version].numFiles(), version)
	}
	uniques.showCounts("All versions")
	for _, version := range versions {
		versionUniques[version].showCounts(fmt.Sprintf("PDF version %s", version))
	}
}

// fontOccurence represents an occurrence of a font in a PDF file.
type fontOccurence struct {
	subtype  string
	basefont string
	filename string
	version  string
}

// fontOccurence represents a list of fontOccurences.
type occurrenceList []fontOccurence

// showCounts prints a summary of the font counts in `occurrences`
func (occurrences occurrenceList) showCounts(title string) {
	numFiles := occurrences.numFiles()

	fmt.Println("=====================================================")
	fmt.Printf("%s\n", title)
	fmt.Printf("%d files tested\n", numFiles)
	fmt.Printf("%d font subtype occurrences\n", len(occurrences))

	subtypeCounts := occurrences.subtypeCounts()

	// Reduce the Type0`,`CIDFontType0` and `CIDFontType2` counts to
	// to `CIDFontType0` and `CIDFontType2` subsets of Type0
	cids := subtypeCounts[`CIDFontType0`] + subtypeCounts[`CIDFontType2`]
	if cids < subtypeCounts[`Type0`] {
		fmt.Fprintf(os.Stderr, "This is impossible: subtypeCounts=%+v\n", subtypeCounts)
		return
	}
	delete(subtypeCounts, `Type0`)

	subtypes := mapKeys(subtypeCounts)
	sort.SliceStable(subtypes, func(i, j int) bool {
		return subtypeCounts[subtypes[i]] > subtypeCounts[subtypes[j]]
	})

	fmt.Printf("%d subtypes\n", len(subtypes))
	for i, subtype := range subtypes {
		count := subtypeCounts[subtype]
		percentSubtypes := float64(count) / float64(len(occurrences)-cids) * 100.0
		percentFiles := float64(count) / float64(numFiles) * 100.0
		fmt.Printf("%4d: %#-14q %4d files (%5.1f%%) Occurred in %.1f%% of files.\n",
			i, subtype, count, percentSubtypes, percentFiles)
	}
}

// numFiles returns the number of unique PDF files in `occurrences`.
func (occurrences occurrenceList) numFiles() int {
	counts := map[string]int{}
	for _, o := range occurrences {
		counts[o.filename]++
	}
	return len(counts)
}

// subtypeCounts returns the number of occurrences of each font subtype in `occurrences`.
func (occurrences occurrenceList) subtypeCounts() map[string]int {
	counts := map[string]int{}
	for _, o := range occurrences {
		counts[o.subtype]++
	}
	return counts
}

// versionCounts returns the number of occurrence of each PDF version in `occurrences`.
func (occurrences occurrenceList) byVersion() map[string]occurrenceList {
	versionOccurrences := map[string]occurrenceList{}
	for _, o := range occurrences {
		versionOccurrences[o.version] = append(versionOccurrences[o.version], o)
	}
	return versionOccurrences
}

// fontsInPdfList returns the fonts used in the PDF files in `pathList`
func fontsInPdfList(pathList []string) (numFiles int, occurrences, uniques occurrenceList, err error) {
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

// showOneFile prints the fonts used in PDF file `filename`
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

// mapKeys returns the keys in `m`
func mapKeys(m map[string]int) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return
}

func versionKeys(m map[string]occurrenceList) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return
}
