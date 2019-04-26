// Shows all fonts in a PDF file.
//
// Run as: go run pdf_fonts.go testdata/*.pdf testdata/**/*.pdf
//

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/bmatcuk/doublestar"
	"github.com/unidoc/unidoc/common"
	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

const usage = "Usage: go run pdf_fonts.go testdata/*.pdf\n"

func main() {
	var debug, trace bool
	flag.BoolVar(&debug, "d", false, "Print debugging information.")
	flag.BoolVar(&trace, "e", false, "Print detailed debugging information.")
	makeUsage(usage)

	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	if trace {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelTrace))
	} else if debug {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
	} else {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelInfo))
	}

	pathList, err := patternsToPaths(args)
	if err != nil {
		os.Exit(1)
	}

	occurrences, err := fontsInPdfList(pathList)
	if err != nil {
		os.Exit(1)
	}

	versionOccurrences := occurrences.byVersion()
	versions := mapKeysNumFiles(versionOccurrences)

	fmt.Printf("%d total font occurrences\n", len(occurrences))
	fmt.Printf("%4d files total.\n", occurrences.numFiles())
	for _, version := range versions {
		fmt.Printf("%4d files PDF version %s\n", versionOccurrences[version].numFiles(), version)
	}
	occurrences.showSubtypeCounts("All versions")
	for _, version := range versions {
		versionOccurrences[version].showSubtypeCounts(fmt.Sprintf("PDF version %s", version))
	}

	subtypeOccurrences := occurrences.bySubtype()
	subtypes := mapKeysNumFiles(subtypeOccurrences)
	occurrences.showEncodingCounts("All font subtypes")
	for _, subtype := range subtypes {
		subtypeOccurrences[subtype].showEncodingCounts(fmt.Sprintf("Font subtype %#q", subtype))
	}
	occurrences.showTounicodeCounts("All font subtypes")
	for _, subtype := range subtypes {
		subtypeOccurrences[subtype].showTounicodeCounts(fmt.Sprintf("Font subtype %#q", subtype))
	}

}

// fontOccurrence represents an occurrence of a font in a PDF file.
type fontOccurrence struct {
	filename  string
	version   string
	subtype   string
	encoding  string
	tounicode string
	basefont  string
}

// fontOccurrence represents a list of fontOccurences.
type occurrenceList []fontOccurrence

// showSubtypeCounts prints a summary of the font subtype counts in `occurrences`
func (occurrences occurrenceList) showSubtypeCounts(title string) {
	numFiles := occurrences.numFiles()
	keyOccurrences := occurrences.bySubtype()

	fmt.Println("===================================================== Subtype")
	fmt.Printf("%s\n", title)
	fmt.Printf("%d files tested\n", numFiles)
	fmt.Printf("%d font subtype occurrences\n", len(occurrences))

	keys := mapKeysNumFiles(keyOccurrences)

	fmt.Printf("%d subtypes\n", len(keys))
	for i, k := range keys {
		nOcc := len(keyOccurrences[k])
		nFiles := keyOccurrences[k].numFiles()
		percentOcc := float64(nOcc) / float64(len(occurrences)) * 100.0
		percentFiles := float64(nFiles) / float64(numFiles) * 100.0
		fmt.Printf("%2d: %#-18s %4d occurrences (%2.0f%%) Occurred in %3d (%2.0f%%) of files.\n",
			i, truncate(k, 20), nOcc, percentOcc, nFiles, percentFiles)
	}
}

// showEncodingCounts prints a summary of the font encoding counts in `occurrences`
func (occurrences occurrenceList) showEncodingCounts(title string) {
	numFiles := occurrences.numFiles()

	fmt.Println("===================================================== Encoding")
	fmt.Printf("%s\n", title)
	fmt.Printf("%d files tested\n", numFiles)
	fmt.Printf("%d font encoding occurrences\n", len(occurrences))

	keyOccurrences := occurrences.byEncoding()

	keys := mapKeysNumFiles(keyOccurrences)

	fmt.Printf("%d encodings\n", len(keys))
	for i, k := range keys {
		nOcc := len(keyOccurrences[k])
		nFiles := keyOccurrences[k].numFiles()
		percentOcc := float64(nOcc) / float64(len(occurrences)) * 100.0
		percentFiles := float64(nFiles) / float64(numFiles) * 100.0
		fmt.Printf("%4d: %-52s %4d occurrences (%2.0f%%) Occurred in %3d (%2.0f%%) of files.\n",
			i, truncate(k, 52), nOcc, percentOcc, nFiles, percentFiles)
	}
}

// showTounicodeCounts prints a summary of the  ToUnicode cmap name counts in `occurrences`
func (occurrences occurrenceList) showTounicodeCounts(title string) {
	numFiles := occurrences.numFiles()

	fmt.Println("===================================================== ToUnicode")
	fmt.Printf("%s\n", title)
	fmt.Printf("%d files tested\n", numFiles)
	fmt.Printf("%d ToUnicode cmap name occurrences\n", len(occurrences))

	keyOccurrences := occurrences.byTounicode()

	allKeys := mapKeysNumFiles(keyOccurrences)
	numNames := len(allKeys)
	keyOccurrences = topElements(keyOccurrences, 10)

	keys := mapKeysNumFiles(keyOccurrences)

	namesString := ""
	if numNames > len(keys) {
		namesString = fmt.Sprintf("%#q", truncateSlice(allKeys, 10, 20))
	}
	fmt.Printf("%d ToUnicode cmap names %s\n", numNames, namesString)

	fmt.Printf("%d ToUnicode cmap names\n", numNames)
	for i, k := range keys {
		nOcc := len(keyOccurrences[k])
		nFiles := keyOccurrences[k].numFiles()
		percentOcc := float64(nOcc) / float64(len(occurrences)) * 100.0
		percentFiles := float64(nFiles) / float64(numFiles) * 100.0
		fmt.Printf("%2d: %-22s %4d occurrences (%2.0f%%) Occurred in %3d (%2.0f%%) of files.\n",
			i, truncate(k, 22), nOcc, percentOcc, nFiles, percentFiles)
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

// topElements returns the `n` elements of `keyOccurrences` with the occurrences that appear in the
// most PDF files.
// The remaining elements are consolidated in a new "[other]" element.
func topElements(keyOccurrences map[string]occurrenceList, n int) map[string]occurrenceList {
	if len(keyOccurrences) <= n {
		return keyOccurrences
	}
	keys := mapKeysNumFiles(keyOccurrences)
	other := occurrenceList{}
	for _, k := range keys[n:] {
		other = append(other, keyOccurrences[k]...)
	}
	top := map[string]occurrenceList{"[other]": other}
	for _, k := range keys[:n] {
		top[k] = keyOccurrences[k]
	}
	return top
}

// byFilename returns `occurrences` as a map of occurrenceLists keyed by PDF file name.
func (occurrences occurrenceList) byFilename() map[string]occurrenceList {
	keyOccurrences := map[string]occurrenceList{}
	for _, o := range occurrences {
		keyOccurrences[o.filename] = append(keyOccurrences[o.filename], o)
	}
	return keyOccurrences
}

// byVersion returns `occurrences` as a map of occurrenceLists keyed by PDF version.
func (occurrences occurrenceList) byVersion() map[string]occurrenceList {
	keyOccurrences := map[string]occurrenceList{}
	for _, o := range occurrences {
		keyOccurrences[o.version] = append(keyOccurrences[o.version], o)
	}
	return keyOccurrences
}

// bySubtype returns `occurrences` as a map of occurrenceLists keyed by font subtype.
func (occurrences occurrenceList) bySubtype() map[string]occurrenceList {
	keyOccurrences := map[string]occurrenceList{}
	for _, o := range occurrences {
		keyOccurrences[o.subtype] = append(keyOccurrences[o.subtype], o)
	}
	return keyOccurrences
}

// byEncoding returns `occurrences` as a map of occurrenceLists keyed by font encoding.
func (occurrences occurrenceList) byEncoding() map[string]occurrenceList {
	keyOccurrences := map[string]occurrenceList{}
	for _, o := range occurrences {
		keyOccurrences[o.encoding] = append(keyOccurrences[o.encoding], o)
	}
	return keyOccurrences
}

// byTounicode returns `occurrences` as a map of occurrenceLists keyed by ToUnicode cmap name.
func (occurrences occurrenceList) byTounicode() map[string]occurrenceList {
	keyOccurrences := map[string]occurrenceList{}
	for _, o := range occurrences {
		keyOccurrences[o.tounicode] = append(keyOccurrences[o.tounicode], o)
	}
	return keyOccurrences
}

// mergeOccurrences returns the elements of the keyed `keyOccurrences` merged into a single
// occurrenceList.
func mergeOccurrences(keyOccurrences map[string]occurrenceList) occurrenceList {
	occurrences := occurrenceList{}
	for _, occ := range keyOccurrences {
		occurrences = append(occurrences, occ...)
	}
	return occurrences
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
			// CIDFontType? font objects aren't fonts.
			if subtype == `CIDFontType0` || subtype == `CIDFontType2` {
				continue
			}
			basefont := font.BaseFont()
			encoding := "[none]"
			if font.Encoder() != nil {
				encoding = font.Encoder().String()
			}
			tounicode := font.ToUnicode()
			if tounicode == "" {
				tounicode = "[none]"
			}
			o := fontOccurrence{filename, version, subtype, encoding, tounicode, basefont}
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

	version = pdfReader.PdfVersion().String()

	// FIXME(peterwilliams97)  GetIndirectObjectByNumber doesn't resolve object references in objects
	// referenced by the returned object. e.g. Stueckelberg6LV.pdf
	for _, objNum := range pdfReader.GetObjectNums() {
		var obj pdfcore.PdfObject
		obj, err = pdfReader.GetIndirectObjectByNumber(objNum)
		if err != nil {
			continue
		}

		obj = pdfcore.TraceToDirectObject(obj)
		if stream, is := obj.(*pdfcore.PdfObjectStream); is {
			obj = stream.PdfObjectDictionary
		}
		if dict, is := obj.(*pdfcore.PdfObjectDictionary); !is {
			continue
		} else {
			typ := dict.Get("Type")
			if typ == nil {
				continue
			}
			if n, ok := pdfcore.GetName(typ); !ok || string(*n) != "Font" {
				continue
			}
		}

		var font *pdf.PdfFont
		font, err = pdf.NewPdfFontFromPdfObject(obj)
		if err != nil && err != pdf.ErrType1CFontNotSupported && err != pdf.ErrType3FontNotSupported {
			return
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
		files, err := doublestar.Glob(pattern)
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

// mapKeysNumFiles returns the keys in `keyOccurrences` sorted by the number of files they occur in.
func mapKeysNumFiles(keyOccurrences map[string]occurrenceList) (keys []string) {
	keys = mapKeys(keyOccurrences)
	sort.SliceStable(keys, func(i, j int) bool {
		return keyOccurrences[keys[i]].numFiles() > keyOccurrences[keys[j]].numFiles()
	})
	return
}

// mapKeys returns the keys in `keyOccurrences` sorted alphabetically.
func mapKeys(keyOccurrences map[string]occurrenceList) (keys []string) {
	for k := range keyOccurrences {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return
}

// truncate returns a string with the first `n` characters in `s`.
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}

// truncateSlice returns up to `m` of the unique `n`-truncated strings in `arr`.
func truncateSlice(arr []string, n, m int) []string {
	set := map[string]bool{}
	for _, s := range arr {
		set[truncate(s, n)] = true
	}
	truncated := []string{}
	for s := range set {
		truncated = append(truncated, s)
	}
	sort.Strings(truncated)
	if len(truncated) <= m {
		return truncated
	}
	return append(truncated[:m], "...")
}

// makeUsage updates flag.Usage to include usage message `msg`.
func makeUsage(msg string) {
	usage := flag.Usage
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, msg)
		usage()
	}
}
