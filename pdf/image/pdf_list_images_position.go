/*
 * List images in a PDF file.  Passes through each page, goes through the content stream and finds instances of both
 * XObject Images and inline images. Also handles images referred within XObject Form content streams.
 * Additionally outputs a summary of the filters and colorspaces used by the images found.
 *
 * Run as: go run pdf_list_images.go input.pdf
 */

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/user"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar"
	"github.com/unidoc/unidoc/common"
	"github.com/unidoc/unidoc/pdf/extractor"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

var colorspaces = map[string]int{}

const usage = "Usage: go run pdf_list_images_position.go testdata/*.pdf\n"

// isInteresting returns true for the images we are interested in.
// This is a user-supplied function.
// It currently returns true for images that are scaled anisotropically.
func isInteresting(page *pdf.PdfPage, imgMark extractor.ImageMark) bool {
	img := imgMark.Image
	big := img.Width >= 20 && img.Height >= 20
	if !big {
		return false
	}
	return true

	sx, sy := imgMark.CTM.ScalingFactorX(), imgMark.CTM.ScalingFactorY()
	w, h := math.Abs(sx), math.Abs(sy)
	W, H := math.Abs(float64(img.Width)), math.Abs(float64(img.Height))

	method := 5
	switch method {
	case 0: // Rotated
		return math.Abs(imgMark.CTM.Angle()) >= 1.0
	case 1:
		if w < 0.001 {
			return false
		}
		r := h / w
		return r <= 0.5 || r >= 2.0
	case 2: // Very anistropic
		if w < 0.001 || H == 0 {
			return false
		}
		rr := (h / w) * (W / H)
		return rr <= 0.5 || rr >= 2.0
	case 3: // Anistropic
		if w < 0.001 || H == 0 {
			return false
		}
		rr := (h / w) * (W / H)
		return rr <= 0.9 || rr >= 1.1
	case 4: // Inline
		return imgMark.Inline
	case 5: // Clipped
		d := 100.0
		mbox, err := page.GetMediaBox()
		if err != nil {
			panic(err)
		}
		llx, lly := imgMark.CTM.Translation()
		urx, ury := llx+sx, lly+sy
		ok := llx < mbox.Llx-d || lly < mbox.Lly-d || urx > mbox.Urx+d || ury > mbox.Ury+d
		if ok {
			fmt.Fprintf(os.Stderr, "*** mbox=%+v ctm=%s\n", *mbox, imgMark.CTM.String())
		}
		return ok
	}
	return false
}

func main() {
	// Make sure to enter a valid license key.
	// Otherwise text is truncated and a watermark added to the text.
	// License keys are available via: https://unidoc.io
	/*
			license.SetLicenseKey(`
		-----BEGIN UNIDOC LICENSE KEY-----
		...key contents...
		-----END UNIDOC LICENSE KEY-----
		`)
	*/
	var debug, trace bool
	var maxInteresting int
	var maxSizeMB float64
	flag.BoolVar(&debug, "d", false, "Print debugging information.")
	flag.BoolVar(&trace, "e", false, "Print detailed debugging information.")
	flag.IntVar(&maxInteresting, "i", 0, "Stop when this interesting files are found.")
	flag.Float64Var(&maxSizeMB, "s", 0, "Search files up this size in MBytes.")
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

	corpus, err := patternsToPaths(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "patternsToPaths failed: : %v\n", err)
		os.Exit(2)
	}
	corpus = removeDuplicates(corpus)
	corpus = filterCorpus(corpus, maxSizeMB)
	corpus = cleanCorpus(corpus)
	// Process files from smallest to largest.
	sort.Slice(corpus, func(i, j int) bool {
		fi, fj := corpus[i], corpus[j]
		si, sj := fileSizeMB(fi), fileSizeMB(fj)
		if si != sj {
			return si < sj
		}
		return fi < fj
	})

	type result struct {
		numPages    int
		numImages   int
		size        float64
		headline    string
		interesting map[int]map[int]bool
	}
	var results []result

	for i, inputPath := range corpus {
		if (maxCorpus >= 0 && maxCorpus < len(corpus)) && i+1 < maxCorpus {
			continue
		}
		fmt.Fprintf(os.Stderr, "%d of %d: inputPath=%q\n", i+1, len(corpus), inputPath)
		// defer func() {
		// 	if r := recover(); r != nil {
		// 		fmt.Fprintf(os.Stderr, "%d of %d: inputPath=%q Error: %v\n", i+1, len(corpus), inputPath, r)
		// 		os.Exit(33)
		// 	}
		// }()
		if !isWanted(inputPath) {
			continue
		}
		numPages, numImages, report, interesting, err := listImages(inputPath)
		if err != nil {
			fmt.Printf("%d of %d: inputPath=%q Error: %v\n", i+1, len(corpus), inputPath, err)
			continue
		}
		headline := fmt.Sprintf("%d of %d: %.1f MB %d pages %d images %q",
			i+1, len(corpus), fileSizeMB(inputPath), numPages, numImages, inputPath)
		if interesting != nil {
			fmt.Fprintf(os.Stderr, "%d:: %s %s\n", len(results)+1, headline, showInteresting(interesting))
			fmt.Println("===========================================")
			fmt.Printf("%s\n", headline)
			fmt.Printf("Interesting: %s\n", headline)
			fmt.Printf("%s\n", strings.Join(report, ""))
			results = append(results, result{numPages, numImages, fileSizeMB(inputPath), headline, interesting})
			if len(interesting) == 0 {
				common.Log.Error("interesting=%#v", interesting)
				panic("A")
			}
			if sum(interesting) == 0 {
				common.Log.Error("interesting=%#v", interesting)
				panic("B")
			}
		}
		if maxInteresting > 0 && len(results) >= maxInteresting {
			break
		}
	}

	// Sort results by most interesting first. This is my judgement of what is interesting.
	sort.Slice(results, func(i, j int) bool {
		fi, fj := results[i], results[j]
		pi, pj := fi.numPages, fj.numPages
		gi, gj := fi.numImages, fj.numImages
		si, sj := fi.size, fj.size
		ii, ij := sum(fi.interesting), sum(fj.interesting)
		ti := pi*pi + gi*gi
		tj := pj*pj + gj*gj
		if ti != tj {
			return ti < tj
		}
		if gi != gj {
			return gi < gj
		}
		if si != sj {
			return si < sj
		}
		if ii != ij {
			return ii > ij
		}
		return fi.headline < fj.headline
	})

	for _, f := range []*os.File{os.Stdout, os.Stderr} {
		fmt.Fprintln(f, "========================= |||| ========================= ")
		fmt.Fprintf(f, "%d interesting results\n", len(results))
		for i, r := range results {
			fmt.Fprintf(f, "%3d:: %s %s\n", i, r.headline, showInteresting(r.interesting))
		}
	}

	fmt.Printf("=======\nColorspace summary:\n")
	for cs, instances := range colorspaces {
		fmt.Printf(" %s: %d instance(s)\n", cs, instances)
	}
}

func showInteresting(interesting map[int]map[int]bool) string {
	var parts []string
	for _, pageNum := range keysIIB(interesting) {
		parts = append(parts, fmt.Sprintf("Page %d: %+v", pageNum, keysIB(interesting[pageNum])))
	}
	return fmt.Sprintf("Interesting: %d pages %d images: %s",
		len(interesting), sum(interesting), strings.Join(parts, ", "))
}

func keysIIB(m map[int]map[int]bool) []int {
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

func keysIB(m map[int]bool) []int {
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

func sum(interesting map[int]map[int]bool) int {
	n := 0
	for _, p := range interesting {
		n += len(p)
	}
	return n
}

// listImages returns a report on and other information about the images in PDF file `inputPath`.
func listImages(inputPath string) (int, int, []string, map[int]map[int]bool, error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return 0, 0, nil, nil, err
	}
	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return 0, 0, nil, nil, err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return 0, 0, nil, nil, err
	}

	if isEncrypted {
		// Try decrypting with an empty one.
		auth, err := pdfReader.Decrypt([]byte(""))
		if err != nil {
			return 0, 0, nil, nil, err
		}
		if !auth {
			fmt.Println("Need to decrypt with a specified user/owner password")
			return 0, 0, nil, nil, nil
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return numPages, 0, nil, nil, err
	}

	numImages := 0
	report := []string{fmt.Sprintf("PDF Num Pages: %d\n", numPages)}
	interesting := map[int]map[int]bool{}
	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return numPages, numImages, report, interesting, err
		}

		// List images on the page.
		nImages, pageReport, interestingPage, err := listImagesOnPage(page)
		if err != nil {
			return numPages, numImages, report, interesting, err
		}
		if interestingPage == nil {
			continue
		}
		if len(interestingPage) == 0 {
			common.Log.Error("interestingPage=%#v", interestingPage)
			panic("1")
		}

		numImages += nImages
		report = append(report, fmt.Sprintln("--------------------"))
		report = append(report, fmt.Sprintf("Page %d of %d:\n", pageNum, numPages))
		report = append(report, pageReport...)
		interesting[pageNum] = interestingPage
	}

	if len(interesting) == 0 {
		interesting = nil
	}

	return numPages, numImages, report, interesting, nil
}

// listPageImages returns a report on the images in PDF page `page`.
func listImagesOnPage(page *pdf.PdfPage) (int, []string, map[int]bool, error) {
	pageExtractor, err := extractor.New(page)
	if err != nil {
		return 0, nil, nil, err
	}
	images, err := pageExtractor.ExtractPageImages(nil)
	if err != nil {
		return 0, nil, nil, err
	}
	// fmt.Printf("&&& %d images\n", len(images.Images))
	report, interesting := listPageImages(page, images)
	// fmt.Printf("&&& %d images\n", len(images.Images))
	return len(images.Images), report, interesting, nil
}

// listPageImages returns a report on the images in `images`.
func listPageImages(page *pdf.PdfPage, images *extractor.PageImages) ([]string, map[int]bool) {
	var report []string
	interesting := map[int]bool{}

	for i, imgMark := range images.Images {
		img := imgMark.Image

		report = append(report, fmt.Sprintf(" image %d\n", i))
		report = append(report, fmt.Sprintf("  Width: %d\n", img.Width))
		report = append(report, fmt.Sprintf("  Height: %d\n", img.Height))
		report = append(report, fmt.Sprintf("  Color components: %d\n", img.ColorComponents))
		report = append(report, fmt.Sprintf("  BPC: %d\n", img.BitsPerComponent))
		report = append(report, fmt.Sprintf("  Size %.1fx%.1f\n",
			imgMark.CTM.ScalingFactorX(), imgMark.CTM.ScalingFactorY()))
		tx, ty := imgMark.CTM.Translation()
		report = append(report, fmt.Sprintf("  CTM (%.1f,%.1f) ϴ=%.1f\n", tx, ty, imgMark.CTM.Angle()))
		if isInteresting(page, imgMark) {
			interesting[i+1] = true
		}
		// Log colorspace use globally.
		csName := fmt.Sprintf("%d", img.ColorComponents)

		if _, has := colorspaces[csName]; has {
			colorspaces[csName]++
		} else {
			colorspaces[csName] = 1
		}
	}

	if len(interesting) == 0 {
		interesting = nil
	}

	return report, interesting
}

// patternsToPaths returns a list of files matching the patterns in `patternList`
func patternsToPaths(patternList []string) ([]string, error) {
	usr, _ := user.Current()
	dir := usr.HomeDir

	pathList := []string{}
	for _, pattern := range patternList {
		pattern = strings.Replace(pattern, "~", dir, -1)
		fmt.Printf("patternsToPaths: pattern=%q\n", pattern)
		files, err := doublestar.Glob(pattern)
		if err != nil {
			fmt.Fprintf(os.Stderr, "patternsToPaths: Glob failed. pattern=%#q err=%v\n", pattern, err)
			return pathList, err
		}
		for _, path := range files {
			if regularFile(path) {
				pathList = append(pathList, path)
			}
		}
	}
	fmt.Printf("patternsToPaths: %d -> %d\n", len(patternList), len(pathList))
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

// fileSizeMB returns the size of file `path` in megabytes.
func fileSizeMB(path string) float64 {
	fi, err := os.Stat(path)
	if err != nil {
		return -1.0
	}
	return float64(fi.Size()) / 1024.0 / 1024.0
}

// removeDuplicates returns `corpus` with duplicate strings removed.
func removeDuplicates(corpus []string) []string {
	seen := map[string]bool{}
	var cleaned []string
	for _, path := range corpus {
		if _, ok := seen[path]; !ok {
			cleaned = append(cleaned, path)
			seen[path] = true
		}
	}
	return cleaned
}

func filterCorpus(corpus []string, maxSizeMB float64) []string {
	if maxSizeMB <= 0.0 {
		return corpus
	}
	var cleaned []string
	for _, path := range corpus {
		if fileSizeMB(path) <= maxSizeMB {
			cleaned = append(cleaned, path)
		}
	}
	return cleaned
}

// cleanCorpus returns `corpus` with known bad files removed.
func cleanCorpus(corpus []string) []string {
	var cleaned []string
	for _, path := range corpus {
		keep := true
		for _, bad := range badFiles {
			if strings.Contains(path, bad) {
				keep = false
			}
		}
		if keep {
			cleaned = append(cleaned, path)
		}
	}
	return cleaned
}

var badFiles = []string{
	"bookmarks_circular.pdf",            // Stack overflow in reader
	"4865ab395ed664c3ee17.pdf",          // Stack overflow in image forms
	"circularReferencesInResources.pdf", // Stack overflow in image forms
	"mrm-icdar.pdf",                     // !@#$
	"ghmt.pdf",                          // !@#$
	"SA_implementations.pdf",            // !@#$
	"naacl06-shinyama.pdf",              // !@#$
	"a_im_",                             // !@#$
	"CaiHof-CIKM2004.pdf",
	"blurhmt.pdf",
	"ESCP-R reference_151008.pdf",
	"a_imagemask.pdf",
	"sample_chapter_verilab_aop_cookbook.pdf",
	"TWISCKeyDist.pdf",
	"ergodicity/1607.04968.pdf",
	"1812.09449.pdf", // hangs
	"INF586.pdf",     // hangs
}

const maxCorpus = 10920

// makeUsage updates flag.Usage to include usage message `msg`.
func makeUsage(msg string) {
	usage := flag.Usage
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, msg)
		usage()
	}
}

// isWanted is for customising test runs to include desired files.
// It should return true for the files you want to process.
// e.g.The commented core returns true for files containing Type0 font dicts in clear text.
func isWanted(filename string) bool {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return !strings.Contains(string(data), "Linearized")
}
