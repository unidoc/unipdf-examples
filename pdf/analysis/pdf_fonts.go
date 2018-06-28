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
	"strings"

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

	occurrences, err := fontsInPdfList(pathList)
	if err != nil {
		os.Exit(1)
	}
	fmt.Printf("%d total font occurrences\n", len(occurrences))

	uniques := occurrences.collapseBasefont()
	versionOccurrences := uniques.byVersion()
	versions := mapKeys(versionOccurrences)

	fmt.Printf("%4d files total.\n", uniques.numFiles())
	for _, version := range versions {
		fmt.Printf("%4d files PDF version %s\n", versionOccurrences[version].numFiles(), version)
	}
	uniques.showSubtypeCounts("All versions")
	for _, version := range versions {
		versionOccurrences[version].showSubtypeCounts(fmt.Sprintf("PDF version %s", version))
	}

	uniques.showEncodingCounts("All versions")
	for _, version := range versions {
		versionOccurrences[version].showEncodingCounts(fmt.Sprintf("PDF version %s", version))
	}
}

// fontOccurence represents an occurrence of a font in a PDF file.
type fontOccurence struct {
	filename string
	version  string
	subtype  string
	encoding string
	basefont string
}

// fontOccurence represents a list of fontOccurences.
type occurrenceList []fontOccurence

// showSubtypeCounts prints a summary of the font subtype counts in `occurrences`
func (occurrences occurrenceList) showSubtypeCounts(title string) {
	numFiles := occurrences.numFiles()

	fmt.Println("===================================================== Subtypes")
	fmt.Printf("%s\n", title)
	fmt.Printf("%d files tested\n", numFiles)
	fmt.Printf("%d font subtype occurrences\n", len(occurrences))

	subtypeOccurrences := occurrences.bySubtype()

	// Reduce the Type0`,`CIDFontType0` and `CIDFontType2` counts to
	// to `CIDFontType0` and `CIDFontType2` subsets of Type0
	numCIDs := len(subtypeOccurrences[`CIDFontType0`]) + len(subtypeOccurrences[`CIDFontType2`])
	if numCIDs < len(subtypeOccurrences[`Type0`]) {
		fmt.Fprintf(os.Stderr, "This is impossible\n")
		return
	}
	delete(subtypeOccurrences, `Type0`)

	subtypes := mapKeys(subtypeOccurrences)
	sort.SliceStable(subtypes, func(i, j int) bool {
		return len(subtypeOccurrences[subtypes[i]]) > len(subtypeOccurrences[subtypes[j]])
	})

	fmt.Printf("%d subtypes\n", len(subtypes))
	for i, subtype := range subtypes {
		count := len(subtypeOccurrences[subtype])
		percentSubtypes := float64(count) / float64(len(occurrences)-numCIDs) * 100.0
		percentFiles := float64(count) / float64(numFiles) * 100.0
		fmt.Printf("%4d: %#-14q %4d files (%4.1f%%) Occurred in %4.1f%% of files.\n",
			i, subtype, count, percentSubtypes, percentFiles)
	}
}

// showEncodingCounts prints a summary of the font encoding counts in `occurrences`
func (occurrences occurrenceList) showEncodingCounts(title string) {
	numFiles := occurrences.numFiles()

	fmt.Println("===================================================== Encodings")
	fmt.Printf("%s\n", title)
	fmt.Printf("%d files tested\n", numFiles)
	fmt.Printf("%d font encoding occurrences\n", len(occurrences))

	encodingOccurrences := occurrences.byEncoding()

	encodings := mapKeys(encodingOccurrences)
	sort.SliceStable(encodings, func(i, j int) bool {
		return len(encodingOccurrences[encodings[i]]) > len(encodingOccurrences[encodings[j]])
	})

	fmt.Printf("%d encodings\n", len(encodings))
	for i, encoding := range encodings {
		count := len(encodingOccurrences[encoding])
		percentEncodings := float64(count) / float64(len(occurrences)) * 100.0
		percentFiles := float64(count) / float64(numFiles) * 100.0
		fmt.Printf("%4d: %#-21q %4d files (%4.1f%%) Occurred in %4.1f%% of files.\n",
			i, encoding, count, percentEncodings, percentFiles)
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

// byVersion returns `occurrences` as a map of occurrenceLists keyed by PDF version.
func (occurrences occurrenceList) byVersion() map[string]occurrenceList {
	keyedOccurrences := map[string]occurrenceList{}
	for _, o := range occurrences {
		keyedOccurrences[o.version] = append(keyedOccurrences[o.version], o)
	}
	return keyedOccurrences
}

// bySubtype returns `occurrences` as a map of occurrenceLists keyed by font subtype.
func (occurrences occurrenceList) bySubtype() map[string]occurrenceList {
	keyedOccurrences := map[string]occurrenceList{}
	for _, o := range occurrences {
		keyedOccurrences[o.subtype] = append(keyedOccurrences[o.subtype], o)
	}
	return keyedOccurrences
}

// byEncoding returns `occurrences` as a map of occurrenceLists keyed by font encoding.
func (occurrences occurrenceList) byEncoding() map[string]occurrenceList {
	keyedOccurrences := map[string]occurrenceList{}
	for _, o := range occurrences {
		keyedOccurrences[o.encoding] = append(keyedOccurrences[o.encoding], o)
	}
	return keyedOccurrences
}

// collapseSubtype returns `occurrences` with all entries that differ only by subtype replaced by
// a single entry.
func (occurrences occurrenceList) collapseSubtype() occurrenceList {
	fields := fontOccurence{subtype: "remove"}
	return occurrences.reduce(fields)
}

// collapseEncoding returns `occurrences` with all entries that differ only by encoding replaced by
// a single entry.
func (occurrences occurrenceList) collapseEncoding() occurrenceList {
	fields := fontOccurence{encoding: "remove"}
	return occurrences.reduce(fields)
}

// collapseEncoding returns `occurrences` with all entries that differ only by encoding replaced by
// a single entry.
func (occurrences occurrenceList) collapseBasefont() occurrenceList {
	fields := fontOccurence{basefont: "remove"}
	return occurrences.reduce(fields)
}

// reduce returns `occurrences` with all entries that differ only by the fields with non-empty
// fields in `fields` replaced by a single entry.
func (occurrences occurrenceList) reduce(fields fontOccurence) occurrenceList {
	keyOccurrence := map[string]fontOccurence{}
	for _, o := range occurrences {
		key := o.toKey(fields)
		keyOccurrence[key] = o
	}
	collapsed := occurrenceList{}
	for _, o := range keyOccurrence {
		collapsed = append(collapsed, o)
	}
	return collapsed
}

// toKey returns a string containing the fields in `o` that aren't set in `fields`.
func (o fontOccurence) toKey(fields fontOccurence) string {
	parts := []string{"", "", "", "", ""}
	if fields.filename == "" {
		parts[0] = o.filename
	}
	if fields.version == "" {
		parts[1] = o.version
	}
	if fields.subtype == "" {
		parts[2] = o.subtype
	}
	if fields.encoding == "" {
		parts[3] = o.encoding
	}
	if fields.basefont == "" {
		parts[4] = o.basefont
	}
	return strings.Join(parts, ":")
}

// fontsInPdfList returns the fonts used in the PDF files in `pathList`.
func fontsInPdfList(pathList []string) (occurrences occurrenceList, err error) {
	for i, filename := range pathList {
		version := ""
		fonts := []pdf.PdfFont{}
		version, fonts, err = fontsInPdf(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %d of %d %q err=%v\n", i, len(pathList), filename, err)
			continue
		}
		for _, font := range fonts {
			subtype := font.Subtype()
			basefont := font.BaseFont()
			encoding := ""
			if font.Encoder() != nil {
				encoding = font.Encoder().String()
			}
			o := fontOccurence{filename, version, subtype, encoding, basefont}
			occurrences = append(occurrences, o)
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

// showOneFile prints the fonts used in PDF file `filename`.
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

// patternsToPaths returns a list of files matching the patterns in `patternList`.
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

// regularFile returns true if file `path` is a regular file.
func regularFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not stat %q. err=%v", path, err)
		return false
	}
	return fi.Mode().IsRegular()
}

// mapKeys returns the keys in `m` sorted alphabetically.
func mapKeys(m map[string]occurrenceList) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return
}
