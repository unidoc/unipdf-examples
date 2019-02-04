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
	var showHelp, debug, trace bool
	flag.BoolVar(&showHelp, "h", false, "Show this help message.")
	flag.BoolVar(&debug, "d", false, "Print debugging information.")
	flag.BoolVar(&trace, "e", false, "Print detailed debugging information.")
	makeUsage(usage)

	flag.Parse()
	args := flag.Args()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}
	if trace {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelTrace))
	} else if debug {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
	} else {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelError))
	}

	corpus, err := patternsToPaths(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "patternsToPaths failed: : %v\n", err)
		os.Exit(2)
	}
	corpus = removeDuplicates(corpus)
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
		numPages  int
		numImages int
		size      float64
		summary   string
	}
	var results []result

	for i, inputPath := range corpus {
		fmt.Printf("^^^ %d: %s\n", i, inputPath)
		// defer func() {
		// 	if r := recover(); r != nil {
		// 		fmt.Fprintf(os.Stderr, "%d of %d: inputPath=%q Error: %v\n", i+1, len(corpus), inputPath, r)
		// 		os.Exit(33)
		// 	}
		// }()
		numPages, numImages, report, interesting, err := listImages(inputPath)
		if err != nil {
			fmt.Printf("%d of %d: inputPath=%q Error: %v\n", i+1, len(corpus), inputPath, err)
			continue
		}
		headline := fmt.Sprintf("%d of %d: %.1f MB %d pages %d images %q",
			i+1, len(corpus), fileSizeMB(inputPath), numPages, numImages, inputPath)
		if interesting {
			fmt.Fprintf(os.Stderr, "%s\n", headline)
			fmt.Println("===========================================")
			fmt.Printf("%s\n", headline)
			fmt.Printf("%s\n", strings.Join(report, ""))
			results = append(results, result{numPages, numImages, fileSizeMB(inputPath), headline})
		}
	}

	// Sort results by most interesting first. This is my judgement of what is interesting.
	sort.Slice(results, func(i, j int) bool {
		fi, fj := results[i], results[j]
		pi, pj := fi.numPages, fj.numPages
		gi, gj := fi.numImages, fj.numImages
		si, sj := fi.size, fj.size
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
		return fi.summary < fi.summary
	})

	for _, f := range []*os.File{os.Stdout, os.Stderr} {
		fmt.Fprintln(f, "========================= |||| ========================= ")
		fmt.Fprintf(f, "%d interesting results\n", len(results))
		for i, r := range results {
			fmt.Fprintf(f, "%3d:: %s\n", i, r.summary)
		}
	}

	fmt.Printf("=======\nColorspace summary:\n")
	for cs, instances := range colorspaces {
		fmt.Printf(" %s: %d instance(s)\n", cs, instances)
	}
}

// listImages returns a report on and other information about the images in PDF file `inputPath`.
func listImages(inputPath string) (int, int, []string, bool, error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return 0, 0, nil, false, err
	}
	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return 0, 0, nil, false, err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return 0, 0, nil, false, err
	}

	if isEncrypted {
		// Try decrypting with an empty one.
		auth, err := pdfReader.Decrypt([]byte(""))
		if err != nil {
			return 0, 0, nil, false, err
		}
		if !auth {
			fmt.Println("Need to decrypt with a specified user/owner password")
			return 0, 0, nil, false, nil
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return numPages, 0, nil, false, err
	}

	numImages := 0
	report := []string{fmt.Sprintf("PDF Num Pages: %d\n", numPages)}
	interesting := true
	for pageNum := 1; pageNum <= numPages; pageNum++ {
		fmt.Printf("^^$ page %d\n", pageNum)
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return numPages, numImages, report, interesting, err
		}

		// List images on the page.
		nImages, pageReport, interestingPage, err := listImagesOnPage(page)
		if err != nil {
			return numPages, numImages, report, interesting, err
		}
		if !interestingPage {
			continue
		}
		numImages += nImages
		report = append(report, fmt.Sprintln("--------------------"))
		report = append(report, fmt.Sprintf("Page %d of %d:\n", pageNum, numPages))
		report = append(report, pageReport...)
		interesting = interesting || interestingPage
	}

	return numPages, numImages, report, interesting, nil
}

// listPageImages returns a report on the images in PDF page `page`.
func listImagesOnPage(page *pdf.PdfPage) (int, []string, bool, error) {
	pageExtractor, err := extractor.New(page)
	if err != nil {
		return 0, nil, false, err
	}
	images, err := pageExtractor.ExtractPageImages()
	if err != nil {
		return 0, nil, false, err
	}
	fmt.Printf("&&& %d images\n", len(images.Images))
	report, interesting := listPageImages(images)
	fmt.Printf("&&& %d images\n", len(images.Images))
	return len(images.Images), report, interesting, nil
}

// listPageImages returns a report on the images in `images`.
func listPageImages(images *extractor.PageImages) ([]string, bool) {
	var report []string
	interesting := true

	for i, imgData := range images.Images {
		img := imgData.Image

		report = append(report, fmt.Sprintf(" image %d\n", i))
		report = append(report, fmt.Sprintf("  Width: %d\n", img.Width))
		report = append(report, fmt.Sprintf("  Height: %d\n", img.Height))
		report = append(report, fmt.Sprintf("  Color components: %d\n", img.ColorComponents))
		report = append(report, fmt.Sprintf("  BPC: %d\n", img.BitsPerComponent))
		report = append(report, fmt.Sprintf("  Size %.1fx%.1f\n", imgData.Width, imgData.Height))
		report = append(report, fmt.Sprintf("  CTM (%.1f,%.1f) ϴ=%.1f\n", imgData.X, imgData.Y, imgData.Angle))
		if math.Abs(imgData.Angle) >= 1.0 && img.Width >= 100.0 && img.Height >= 100.0 {
			interesting = true
		}
		// Log colorspace use globally.
		csName := fmt.Sprintf("%d", img.ColorComponents)

		if _, has := colorspaces[csName]; has {
			colorspaces[csName]++
		} else {
			colorspaces[csName] = 1
		}
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
}

// makeUsage updates flag.Usage to include usage message `msg`.
func makeUsage(msg string) {
	usage := flag.Usage
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, msg)
		usage()
	}
}
