/*
 * PDF text segmentation: Find bounding boxes around all regions of PDF pages that contain text.
 *   Mark up the PDF file with rectangles for these regisions
 *
 * Run as: go run pdf_segment_text.go ~/testdata/pc-test/ocr/300-deu.pdf marked_up.pdf
 */

package main

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/unidoc/unidoc/common"
	"github.com/unidoc/unidoc/pdf/creator"
	"github.com/unidoc/unidoc/pdf/extractor"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

const (
	usage                 = "Usage: go run pdf_segment_text.go my_file.pdf marked_up.pdf\n"
	badFilesPath          = "bad.files"
	defaultNormalizeWidth = 60
)

const (
	// Tolerances for merging a column of text bounding boxes
	boxTolX0 = 50.0 // Left X tolerance for merging.
	boxTolX1 = 75.0 // Right X tolerance for merging.
	boxTolY  = 50.0 // Y (downwards) tolerance for merging
)

var (
	// Sanity checkes
	maxWidth      = 20.0
	maxHeight     = 100.0
	maxText       = 400
	validatePanic = false
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run pdf_text_locations.go input.pdf markedup.pdf\n")
		os.Exit(1)
	}

	// Make sure to enter a valid license key.
	// Otherwise text is truncated and a watermark added to the text.
	// License keys are available via: https://unidoc.io
	/*
			license.SetLicenseKey(`
		-----BEGIN UNIDOC LICENSE KEY-----
		...key contents...
		-----END UNIDOC LICENSE KEY-----
		`, "Key owner")
	*/

	var debug, trace, coalesce, merge, occlude bool
	firstPage := 1
	maxPages := 1
	maxText := 500
	maxLocations := 100

	flag.BoolVar(&debug, "d", false, "Print debugging information.")
	flag.BoolVar(&trace, "e", false, "Print detailed debugging information.")
	flag.BoolVar(&coalesce, "c", false, "coalesce text bounding boxes into word bounding boxed.")
	flag.BoolVar(&merge, "m", false, "Merge adjacent text bounding boxes.")
	flag.BoolVar(&occlude, "o", false, "Remove occluded text bounding boxes.")
	flag.IntVar(&firstPage, "p", firstPage, "First page to extract.")
	flag.IntVar(&maxPages, "n", maxPages, "Number of pages to extract. (-1 = all pages from first page")
	flag.IntVar(&maxText, "t", maxText, "Maximum number of characters of text to show per page.")
	flag.IntVar(&maxLocations, "l", maxLocations, "Maximum number of locations to show per page.")
	makeUsage(usage)
	flag.Parse()
	args := flag.Args()

	if len(args) < 2 {
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
	inPath := args[0]
	outPath := args[1]

	if coalesce {
		maxWidth = 500
	}

	// markupList is the list of text locations that will be marked up on the output PDF.
	var markupList []docTextLocations

	docLocations, err := readDocLocations(inPath, firstPage, maxPages, maxText, maxLocations)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	docLocations = docLocations.clean()
	// markupList = append(markupList, docLocations)

	var words, blobs, visible docTextLocations
	if coalesce {
		words = docLocations.coalesceToWords()
		words = words.clean()
		markupList = append(markupList, words)
	} else {
		words = docLocations
	}
	if merge {
		blobs = words.mergeAdjacent()
		blobs = blobs.clean()
		markupList = append(markupList, blobs)
	} else {
		blobs = words
	}
	if occlude {
		visible = blobs.removeHidden()
		visible = visible.clean()
		markupList = append(markupList, visible)
	} else {
		visible = blobs
	}

	fmt.Println("BEFORE (raw)")
	fmt.Printf("%s\n", docLocations.String())

	fmt.Println("AFTER (coalesced)")
	fmt.Printf("%s\n", words.String())

	fmt.Println("MERGED (merged)")
	fmt.Printf("%s\n", blobs.String())

	fmt.Println("VISIBLE (removed hidden)")
	fmt.Printf("%s\n", visible.String())

	common.Log.Info("outPath=%q pageLocations=%d", outPath, len(docLocations))
	if err := markupPdfResults(outPath, markupList); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// docTextLocations describes the locations of extracted text on the pages in a PDF document.
type docTextLocations []pageTextLocations

// pageTextLocations describes the locations of extracted text on a PDF page.
type pageTextLocations struct {
	pageNum   int                      // The page number.
	page      *pdf.PdfPage             // The page.
	locations []extractor.TextLocation // The locations of the text extracted from `page`.
	text      string                   // The extracted text.
}

func (dl docTextLocations) String() string {
	var sb strings.Builder

	fmt.Fprintln(&sb, separator)
	fmt.Printf("PDF text location extraction: %d pages\n", len(dl))
	fmt.Fprintln(&sb, separator)
	for _, pl := range dl {
		fmt.Fprintf(&sb, "%s\n", pl.String())
	}
	fmt.Fprintln(&sb, separator)

	return sb.String()
}

func (pl pageTextLocations) String() string {
	var sb strings.Builder

	fmt.Fprintln(&sb, separText)
	fmt.Fprintf(&sb, "Page %d: %d chars %d locations\n", pl.pageNum, len(pl.text), len(pl.locations))
	fmt.Fprintf(&sb, "\"%s\"\n", pl.text)
	fmt.Fprintln(&sb, separLocs)
	for i, loc := range pl.locations {
		fmt.Fprintf(&sb, "%6d: %s\n", i, loc)
	}
	fmt.Fprintf(&sb, "max=%s\n", pl.max().String())
	fmt.Fprintln(&sb, separator)

	return sb.String()
}

// readDocLocations extracts the text positions in PDF file `inPath` and returns them in a
// docTextLocations.
// Text is extracted from firstPage (1-offset). Up to `maxPages` page of text are extracted or up
// to the end of the PDF if `maxPages` < 0.
// Up to `maxText` characters of the extracted text are displayed.
// Up to `maxLocations` TextLocation's are extracted per page.
func readDocLocations(inPath string, firstPage, maxPages, maxText, maxLocations int) (docTextLocations, error) {
	f, err := os.Open(inPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return nil, err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return nil, err
	}

	if firstPage < 1 {
		firstPage = 1
	}
	if firstPage > numPages {
		firstPage = numPages
	}
	if maxPages > 0 && firstPage+maxPages-1 < numPages {
		numPages = firstPage + maxPages - 1
	}

	fmt.Println(separator)
	fmt.Printf("-PDF text location extraction: pages %d - %d\n", firstPage, numPages)
	fmt.Println(separator)

	var pageLocations docTextLocations

	for pageNum := firstPage; pageNum <= numPages; pageNum++ {
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return nil, err
		}

		ex, err := extractor.New(page)
		if err != nil {
			return nil, err
		}

		pageText, _, _, err := ex.ExtractPageText()
		if err != nil {
			return nil, err
		}
		text, locations := pageText.ToTextLocation()
		if maxLocations >= 0 && len(locations) > maxLocations {
			locations = locations[:maxLocations]
		}
		if maxText >= 0 && len(text) > maxText {
			text = text[:maxText]
		}
		pl := pageTextLocations{
			pageNum:   pageNum,
			page:      page,
			locations: locations,
			text:      text,
		}
		pageLocations = append(pageLocations, pl)
		fmt.Println(separator)
	}

	common.Log.Info("+PDF text location extraction: %d pages", numPages)

	return pageLocations, nil
}

const (
	separator = "---------------------------------------------------------------"
	separText = "--------------------------- TEXT ------------------------------"
	separLocs = "------------------------- LOCATIONS ---------------------------"
)

// clean removes text locations containing only space. These locations mess up the down stream
// processing. !@#$
func (dl docTextLocations) clean() docTextLocations {
	out := make(docTextLocations, len(dl))
	for i, pl := range dl {
		out[i] = pl
		out[i].locations = cleanPage(pl.locations)
	}
	return out
}

// coalesceToWords combines individual letter text locations into word text locations for PDF document `dl`.
func (dl docTextLocations) coalesceToWords() docTextLocations {
	out := make(docTextLocations, len(dl))
	for i, pl := range dl {
		out[i] = pl.coalesceToWords()
	}
	return out
}

// coalesceToWords combines individual letter text locations into word text locations for PDF page `pl`.
func (pl pageTextLocations) coalesceToWords() pageTextLocations {
	return pageTextLocations{
		pl.pageNum,
		pl.page,
		coalescePage(pl.locations),
		pl.text,
	}
}

// mergeAdjacent merges adjacent text locations for PDF document `dl`.
func (dl docTextLocations) mergeAdjacent() docTextLocations {
	out := make(docTextLocations, len(dl))
	for i, pl := range dl {
		out[i] = pl.mergeAdjacent()
	}
	// panic("--merge")
	return out
}

// mergeAdjacent merges adjacent text locations for PDF page `pl`.
func (pl pageTextLocations) mergeAdjacent() pageTextLocations {
	return pageTextLocations{
		pl.pageNum,
		pl.page,
		mergeLocations(pl.locations),
		pl.text,
	}
}

func (dl docTextLocations) removeHidden() docTextLocations {
	out := make(docTextLocations, len(dl))
	for i, pl := range dl {
		out[i] = pl.removeHidden()
	}
	// panic("--merge")
	return out
}

func (pl pageTextLocations) removeHidden() pageTextLocations {
	return pageTextLocations{
		pl.pageNum,
		pl.page,
		removeHidden(pl.locations),
		pl.text,
	}
}

// cleanPage removes pure space text locations from `locations`.
// FIXME !@#$ Don't remove these. Mark them.
func cleanPage(locations []extractor.TextLocation) []extractor.TextLocation {
	// return locations
	var out []extractor.TextLocation
	for i, loc := range locations {
		if loc.BBox.Llx == 0.0 {
			if len(loc.Text) > 0 && loc.Text != " " {
				panic(fmt.Errorf("Llx == 0.0: i=%d loc=%s\n\t%#v", i, loc, loc))
			}
			continue
		}
		if len(loc.Text) == 0 || loc.Text == " " {
			continue
		}
		if !validate(loc) {
			ppanic("cleanPage: too high")
		}
		out = append(out, loc)
	}
	return out
}

// coalescePage combines individual letter text locations into word text locations.
func coalescePage(locations []extractor.TextLocation) []extractor.TextLocation {
	common.Log.Info("===================== coalescePage %d", len(locations))

	var out []extractor.TextLocation
	if len(locations) == 0.0 {
		return out
	}

	tolWPage := aveWidth(locations)
	tolHPage := aveHeight(locations) * 1.0
	common.Log.Info("tolWPage=%.2f tolHPage=%.2f", tolWPage, tolHPage)

	lines := toLines(locations, tolHPage)

	common.Log.Info("===================== %d lines", len(lines))
	for i, ln := range lines {
		var parts []string
		for _, loc := range ln {
			parts = append(parts, loc.Text)
		}
		text := strings.Join(parts, "")
		common.Log.Info("line %d: %d %q", i, len(ln), text)
	}

	for _, ln := range lines {
		tolWLine := aveWidth(ln)
		tol := 0.5 * (tolWPage + tolWLine)
		ln = coalesceLine(ln, tol)
		out = append(out, ln...)
	}

	return out
}

// toLines splits `locations` into slices of elements that are in the same horizontal line.
// NOTE: This doesn't work with columns and raised lines the "energy transferred" condition below.
func toLines(locations []extractor.TextLocation, tol float64) [][]extractor.TextLocation {
	common.Log.Info("===================== -toLines locations=%d tol=%.2f", len(locations), tol)

	var lines [][]extractor.TextLocation
	if len(locations) == 0 {
		return lines
	}

	y0 := locations[0].BBox.Ury
	ln := []extractor.TextLocation{locations[0]}
	for i := 1; i < len(locations); i++ {
		loc := locations[i]
		y := loc.BBox.Ury
		common.Log.Debug("y0=%.2f y=%.2f y0-tol=%2.f %t", y0, y, y0-tol, y < y0-tol)

		if y < y0-tol && y != 0.0 {
			common.Log.Info("y=%.2f->%.2f %4d %s", y0, y, i, locString(loc))
			if len(ln) > 0 {
				lines = append(lines, ln)
			}
			ln = []extractor.TextLocation{loc}
			y0 = y
			continue
		}
		ln = append(ln, loc)

	}
	if len(ln) > 0 {
		lines = append(lines, ln)
	}
	common.Log.Info("===================== +toLines %d->%d", len(locations), len(lines))
	return lines
}

// coalesceLine combines the characters in `locations` into words.
func coalesceLine(locations []extractor.TextLocation, tol float64) []extractor.TextLocation {
	common.Log.Debug("===================== -coalesceLine %d", len(locations))

	var out []extractor.TextLocation
	if len(locations) == 0 {
		return out
	}
	loc0 := locations[0]
	validate(loc0)
	filling := true
	for i := 1; i < len(locations); i++ {
		loc := locations[i]
		validate(loc)
		common.Log.Debug("%4d %s", i, locString(loc))
		// loc.BBox.Llx > loc0.BBox.Urx+tol is a new word
		// loc.BBox.Urx < loc0.BBox.Urx is a case where toLines combined lines incorrectly. e.g.
		//           energy transferred
		// Voltage = ------------------
		//           Coulomb of charge
		// as {"energy transferred", "Voltage"} in a single line.
		if loc.BBox.Llx > loc0.BBox.Urx+tol || (loc.BBox.Urx < loc0.BBox.Urx && loc.BBox.Lly < loc0.BBox.Lly) {
			out = append(out, loc0)
			loc0 = loc
			filling = false
			continue
		}
		if loc.BBox.Urx < loc0.BBox.Urx {
			common.Log.Error("i=%d\n\tloc0=%v\n\tloc =%v", i, loc0, loc)
			panic("Urx")
		}
		loc0.BBox.Urx = loc.BBox.Urx
		loc0.Text += loc.Text
		validate(loc0)
		filling = true
	}
	if filling {
		out = append(out, loc0)
	}
	common.Log.Debug("===================== +coalesceLine %d->%d", len(locations), len(out))
	return out
}

// mergeLocations merges adjacent text locations.
// sort boxes top to bottom, left to right
// for each box:
//     Look for boxes below that are contiguous in y and
//     have left and right sides within tolerance
func mergeLocations(locations []extractor.TextLocation) []extractor.TextLocation {
	common.Log.Info("===================== -mergeLocations %d", len(locations))
	for i00, loc00 := range locations {
		fmt.Printf("\t%4d: %+v\n", i00, loc00)
	}
	common.Log.Info("===================== *mergeLocations %d", len(locations))

	// panic("merge")
	var out []extractor.TextLocation
	if len(locations) == 0.0 {
		return out
	}
	merged := make([]bool, len(locations))
	var mergeList []mergeItem

	count1 := 0

	// locations is sorted top to bottom, left to right (!@#$ Ll or Ur)
	for i00, loc00 := range locations {
		mergeResult := mergeItem{i00: i00}
		count1++
		if count1 > len(locations)*2 {
			panic("count1")
		}
		if merged[i00] {
			common.Log.Info("i00=%d is merged", i00)
			// panic("merged")
			continue
		}
		common.Log.Info("mergeLocations i00=%d of %d count=%d", i00, len(locations), count1)

		// Find all locations that can be merged to loc00.
		common.Log.Info("--------------666---------------")
		for i0 := i00; i0 < len(locations)-1; {
			adjacent, i1 := findConnected(locations, merged, loc00, i0)
			if i1 < 0 {
				break
			}
			if i1 < i0 {
				panic("stalled")
			}
			for _, j := range adjacent {
				expand(&loc00, locations[j])
				mergeResult.parts = append(mergeResult.parts, j)

				merged[j] = true
				common.Log.Info("  i0=%d j=%d locations=%d \n\tloc00=%s",
					i0, j, len(locations), loc00)
			}
			i0 = i1 + 1
		}

		common.Log.Info("**i00=%d loc00=%+v (%.3f x %.3f)", i00, loc00,
			loc00.BBox.Width(), loc00.BBox.Height())
		out = append(out, loc00)
		mergeResult.loc00 = loc00
		if len(mergeResult.parts) > 0 {
			mergeList = append(mergeList, mergeResult)
		}
	}

	common.Log.Info("===================== +mergeLocations %d->%d", len(locations), len(out))
	common.Log.Info("mergeList=%s", showMergeList(mergeList, locations))
	// panic("done")
	return out
}

type mergeItem struct {
	i00   int
	parts []int
	loc00 extractor.TextLocation
}

func showMergeList(mergeList []mergeItem, locations []extractor.TextLocation) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d merge items\n", len(mergeList))
	for i, m := range mergeList {
		fmt.Fprintf(&sb, "\t%6d: %+v\n", i, m)
		fmt.Fprintf(&sb, "\t\t%6d: %+v\n", m.i00, locations[m.i00])
		for _, j := range m.parts {
			fmt.Fprintf(&sb, "\t\t%6d: %q\n", j, locations[j])
		}
		fmt.Fprintf(&sb, "\t%s: %q\n", "loc00", m.loc00)
	}

	return sb.String()
}

// findConnected returns the indexes into `locations` of the elements after `i0` for which
// `compY` and `compX` return true for pairs on intermediate elements.
func findConnected(locations []extractor.TextLocation, visited []bool,
	loc00 extractor.TextLocation, i0 int) ([]int, int) {

	i1 := -1
	imax := -1
	var adjacent []int
outer:
	for ; i0 < len(locations); i0 = i1 + 1 {
		loc0 := locations[i0]
		i1 = findAdjacentY(locations, loc0, i0)
		if i1 < 0 {
			break outer
		}
		// !@#$ texts is for debugging only
		var texts []string
		for i := i0 + 1; i <= i1; i++ {
			texts = append(texts, locations[i].Text)
		}
		// indexes up to and including i1 are within merging distance of loc00 in the Y direction
		common.Log.Info("findConnected: i0=%d i1=%d loc00=%q\n\t%q", i0, i1, loc00.Text, texts)

		for j := i0 + 1; j <= i1; j++ {
			match := compX(loc0, locations[j])
			common.Log.Info("findConnected: i0=%d j=%d visited=%t comp=%t\n\t%+v\n\t%+v",
				i0, j, visited[j], match,
				loc0, locations[j])
			if visited[j] {
				continue
			}
			if match {
				adjacent = append(adjacent, j)
				imax = j
			}
		}
		if len(adjacent) == 0 {
			break
		}

		common.Log.Info("adjacent=%d %v", len(adjacent), adjacent)
		if len(adjacent) == 0 {
			panic("Can't happen")
		}
	}
	return adjacent, imax
}

// findAdjacentY searches `locations` for elements after `i0`for which `compY` returns true.
func findAdjacentY(locations []extractor.TextLocation, loc00 extractor.TextLocation, i0 int) int {
	imax := -1
	for i, loc := range locations[i0+1:] {
		if !compY(loc00, loc) {
			break
		}
		imax = i0 + 1 + i
	}
	return imax
}

// compY returns true if `loc` is adjacent to `loc0` in the Y direction.
// loc.Lly is guaranteed to be below `loc0.Lly`.
func compY(loc0, loc extractor.TextLocation) bool {
	dy := loc0.BBox.Lly - loc.BBox.Ury
	ok := dy < boxTolY
	common.Log.Debug("\t\tcompY: dy=%.2f %t", dy, ok)
	return ok
}

// compX returns true if `loc` is adjacent to `loc0` in the X direction.
func compX(loc0, loc extractor.TextLocation) bool {
	inside := loc0.BBox.Llx <= loc.BBox.Llx && loc.BBox.Urx <= loc0.BBox.Urx
	common.Log.Info("compX: inside=%t", inside)
	if inside {
		return true
	}
	return compX0(loc0, loc) && compX1(loc0, loc)
}

// compX0 returns true if `loc`, which is below `loc0` matches on the left side.
func compX0(loc0, loc extractor.TextLocation) bool {
	dx := loc.BBox.Llx - loc0.BBox.Llx
	ok := math.Abs(dx) < boxTolX0
	common.Log.Info("compX0: dx=%.2f %t", dx, ok)
	return ok
}

// compX1 returns true if `loc`, which is below `loc0` matches on the right side.
func compX1(loc0, loc extractor.TextLocation) bool {
	dx := loc.BBox.Urx - loc0.BBox.Urx
	ok := math.Abs(dx) < boxTolX1
	common.Log.Info("compX1: dx=%.2f %t", dx, ok)
	return ok
}

func expand(loc0 *extractor.TextLocation, loc extractor.TextLocation) {
	if loc.BBox.Lly < loc0.BBox.Lly {
		loc0.BBox.Lly = loc.BBox.Lly
	}
	if loc.BBox.Ury > loc0.BBox.Ury {
		loc0.BBox.Ury = loc.BBox.Ury
	}
	if loc.BBox.Llx < loc0.BBox.Llx {
		loc0.BBox.Llx = loc.BBox.Llx
	}
	if loc.BBox.Urx > loc0.BBox.Urx {
		loc0.BBox.Urx = loc.BBox.Urx
	}

	loc0.Text += "\n" + loc.Text
	// if !validate(*loc0) {
	// 	panic("yikes")
	// }
	if len(loc.Text) > 200 {
		panic("longggggg")
	}
}

// removeHidden removes text locations that are contained within other text locations.
// !@#$ Need to do better and find smallest number of rectangle that cover contiguous text locations.
func removeHidden(locations []extractor.TextLocation) []extractor.TextLocation {
	common.Log.Info("===================== -removeHidden %d", len(locations))

	llx := orderedBy(locations, func(bi, bj pdf.PdfRectangle) bool { return bi.Llx < bj.Llx })
	lly := orderedBy(locations, func(bi, bj pdf.PdfRectangle) bool { return bi.Lly < bj.Lly })
	urx := orderedBy(locations, func(bi, bj pdf.PdfRectangle) bool { return bi.Urx > bj.Urx })
	ury := orderedBy(locations, func(bi, bj pdf.PdfRectangle) bool { return bi.Ury > bj.Ury })

	common.Log.Info("removeHidden")
	common.Log.Info("llx=%+v", llx)
	common.Log.Info("lly=%+v", lly)
	common.Log.Info("urx=%+v", urx)
	common.Log.Info("ury=%+v", ury)
	for i, loc := range locations {
		common.Log.Info("%6d: %2d,%2d,%2d,%2d %v", i,
			llx.inverse[i], lly.inverse[i], urx.inverse[i], ury.inverse[i], loc)
	}

	var out []extractor.TextLocation
	for i, loc := range locations {
		ellx := llx.below(i)
		elly := lly.below(i)
		eurx := urx.below(i)
		eury := ury.below(i)
		lists := [][]int{ellx, elly, eurx, eury}
		cover := intersection(lists)
		common.Log.Info("%4d: %+v cover=%+v", i, lists, cover)
		if len(cover) == 0 {
			out = append(out, loc)
		} else {
			common.Log.Info("%2d %v is blocked by %v", i, loc, cover)
		}
	}

	common.Log.Info("===================== +removeHidden %d->%d", len(locations), len(out))
	return out
}

type ordering struct {
	order   []int
	inverse []int
}

// orderedBy returns ordering `o` such that
//   locations[o.order[i]] is sorted by `comp`,
//   o.inverse[i] is the index of locations[i] in o.index.
func orderedBy(locations []extractor.TextLocation, comp func(bi, bj pdf.PdfRectangle) bool) ordering {

	order := make([]int, len(locations))
	for i := range locations {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool {
		ki, kj := order[i], order[j]
		bi, bj := locations[ki].BBox, locations[kj].BBox
		return comp(bi, bj)
	})

	inverse := make([]int, len(order))
	for i, k := range order {
		inverse[k] = i
	}
	return ordering{order: order, inverse: inverse}
}

// below returns the indexes in `o` that are below locations[`i`] for the ordering given by `comp`
// over `locations` in orderedBy().
func (o ordering) below(i int) []int {
	k := o.inverse[i]
	return o.order[:k]
}

// intersection returns the sorted intersection of the sets in `lists`.
func intersection(lists [][]int) []int {
	m := slice2map(lists[0])
	for _, lst := range lists[1:] {
		m = intersect(m, lst)
	}
	return map2slice2map(m)
}

// intersect returns the intersection of set `m` and slice `lst` as a set
func intersect(m map[int]struct{}, lst []int) map[int]struct{} {
	c := map[int]struct{}{}
	for _, k := range lst {
		if _, ok := m[k]; ok {
			c[k] = struct{}{}
		}
	}
	return c
}

func slice2map(lst []int) map[int]struct{} {
	m := map[int]struct{}{}
	for _, k := range lst {
		m[k] = struct{}{}
	}
	return m
}

func map2slice2map(m map[int]struct{}) []int {
	lst := make([]int, 0, len(m))
	for k := range m {
		lst = append(lst, k)
	}
	sort.Ints(lst)
	return lst
}

// const granularity = 0.1

// func lowerBound(indexes []int, locations []extractor.TextLocation, i0 int) int {
// 	return indexes[i0]
// 	// bound := locations[i0] - granularity
// 	// for i := i0 - 1; i >= 0; i-- {
// 	// 	if locations[i] < bound {
// 	// 		return i + 1
// 	// 	}
// 	// }
// 	return i0
// }
// func upperBound(indexes []int, locations []extractor.TextLocation, i0 int) int {
// 	return indexes[i0]
// 	// bound := locations[i0] + granularity
// 	// for i := i0 - 1; i >= 0; i-- {
// 	// 	if locations[i] > bound {
// 	// 		return i + 1
// 	// 	}
// 	// }
// 	// return i0
// }

const borderWidth = 1.5

var colors = []string{
	"#0000ff", // Blue border.
	"#ff0000", // Red border.
	"#ffff00", // Yellow border.
}

// markupPdfResults marks up the PDF file described by the page fields in `pageLocations` with the
// text locations in `pageLocations` and saves the marked up PDF file as `outPath`.
func markupPdfResults(outPath string, markupList []docTextLocations) error {
	if len(markupList) == 0 {
		common.Log.Info("markupPdfResults: Nothing to do.")
		return nil
	}
	common.Log.Debug("markupPdfResults: markupList=%d %d outPath=%q", len(markupList),
		len(markupList[0]), outPath)

	// Make a new PDF creator.
	c := creator.New()

	outPages := 0
	for i, pl0 := range markupList[0] {
		common.Log.Info("-markupPdfResults: page %d locations=%d", i+1, len(pl0.locations))
		page := pl0.page
		if page.MediaBox == nil {
			common.Log.Error("page %d: No MediaBox.", i+1)
			continue
		}
		if err := c.AddPage(page); err != nil {
			common.Log.Error("page %d: AddPage failed", i)
			return err
		}
		outPages++

		h := page.MediaBox.Ury
		for k := range markupList {
			pl := markupList[k][i]
			col := colors[k]
			r := float64(len(markupList)) / float64(k+1)
			wid := borderWidth * r
			for j, loc := range pl.locations {
				r := loc.BBox
				if j%1000 == 500 {
					common.Log.Info("markupPdfResults: %4d of %d %s",
						j+1, len(pl.locations), rectString(r))
				}

				rect := c.NewRectangle(r.Llx, h-r.Lly, r.Urx-r.Llx, -(r.Ury - r.Lly))
				rect.SetBorderColor(creator.ColorRGBFromHex(col))
				rect.SetBorderWidth(wid)
				if err := c.Draw(rect); err != nil {
					return err
				}
			}
			common.Log.Info("+markupPdfResults: list=%d page=%d locations=%d col=%q wid=%.3g\n\tmx=%s",
				k+1, i+1, len(pl.locations), col, wid, pl.max())
		}
	}
	if outPages == 0 {
		return errors.New("no pages in marked up PDF")
	}
	return c.WriteToFile(outPath)
}

// locString returns a string describing rectangle `r`.
func locString(loc extractor.TextLocation) string {
	return fmt.Sprintf("<%s %q>", rectString(loc.BBox), loc.Text)
}

func validate(loc extractor.TextLocation) bool {
	b := loc.BBox
	if b.Urx < b.Llx {
		common.Log.Error("@!@ Negative width  loc=%s\n\t%#v", loc, loc)
		panic("validate:width")
		return false
	}
	if b.Ury < b.Lly {
		common.Log.Error("@!@ Negative height  loc=%s\n\t%#v", loc, loc)
		panic("validate:height")
		return false
	}
	// if b.Width() > maxWidth {
	// 	common.Log.Error("@!@ Too wide: w=%.2f  loc=%s\n\t%#v", loc.BBox.Width(), loc, loc)
	// 	ppanic("validate:width")
	// 	return false
	// }
	// if b.Height() > maxHeight {
	// 	common.Log.Error("@!@ Too high: h=%.2f loc=%s\n\t%#v", loc.BBox.Height(), loc, loc)
	// 	ppanic("validate:height")
	// 	return false
	// }
	// if len(loc.Text) > maxText {
	// 	common.Log.Error("@!@ Too long: h=%.2f loc=%s\n\t%#v", loc.BBox.Height(), loc, loc)
	// 	ppanic("validate:text")
	// 	return false
	// }
	return true
}

func ppanic(msg string) {
	if validatePanic {
		panic(msg)
	}
}

// rectString returns a string describing rectangle `r`.
func rectString(r pdf.PdfRectangle) string {
	return fmt.Sprintf("{llx: %4.1f lly: %4.1f urx: %4.1f ury: %4.1f (%.1f x %.1f)}",
		r.Llx, r.Lly, r.Urx, r.Ury, r.Urx-r.Llx, r.Ury-r.Lly)
}

func aveWidth(locations []extractor.TextLocation) float64 {
	if len(locations) == 0.0 {
		return 0.0
	}
	var total float64
	for _, loc := range locations {
		total += loc.BBox.Width()
	}
	return total / float64(len(locations))
}

func aveHeight(locations []extractor.TextLocation) float64 {
	if len(locations) == 0.0 {
		return 0.0
	}
	var total float64
	for _, loc := range locations {
		total += loc.BBox.Height()
	}
	return total / float64(len(locations))
}

// makeUsage updates flag.Usage to include usage message `msg`.
func makeUsage(msg string) {
	usage := flag.Usage
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, msg)
		usage()
	}
}

// truncate returns `text` truncated to `n` characters.
func truncate(text string, n int) string {
	if len(text) < n {
		return text
	}
	return text[:n]
}

// maxTextLocations tracks the biggest bounding boxes in a pageTextLocations
// !@#$ Do we ned this
type maxTextLocations struct {
	w, h       float64 // Largest width and height
	text       string  // Longest text
	iw, ih, it int     // Indexes of locations of w, h and text
	tw, th     string  // Text of w and h
}

func (s maxTextLocations) String() string {
	return fmt.Sprintf("{maxTextLocations: %.1f x %.1f %d  %d,%d,%d\n\t%q\n\t%q\n\t%q}",
		s.w, s.h, len(s.text), s.iw, s.ih, s.it, s.tw, s.th, s.text)
}

func (pl pageTextLocations) max() maxTextLocations {
	var s maxTextLocations
	for i, loc := range pl.locations {
		b := loc.BBox
		if b.Width() > s.w {
			s.w = b.Width()
			s.iw = i
			s.tw = loc.Text
		}
		if b.Height() > s.h {
			s.h = b.Height()
			s.ih = i
			s.th = loc.Text
		}
		if len(loc.Text) > len(s.text) {
			s.text = loc.Text
			s.it = i
		}
	}
	return s
}
