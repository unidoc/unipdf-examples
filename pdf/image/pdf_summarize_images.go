/*
 * Summarize images in a corpus of PDF files.  For each PDF file, passes through each page, goes
 * through the content stream and finds instances of both XObject Images and inline images.  Also
 * handles images referred within XObject Form content streams.
 * Outputs a summary of the images found.
 *
 * Run as: go run pdf_summarize_images.go ~/testdata/*.pdf
 */

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/unidoc/unidoc/common"
	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

const usage = "Usage: go run pdf_summarize_images.go testdata/*.pdf\n"

func main() {
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

	corpus := args[:]
	sort.Slice(corpus, func(i, j int) bool {
		fi, fj := corpus[i], corpus[j]
		si, sj := fileSizeMB(fi), fileSizeMB(fj)
		if si != sj {
			return si < sj
		}
		return fi < fj
	})

	corpusInfo := map[string][]imageInfo{}
	for i, inputPath := range corpus {
		fmt.Fprintf(os.Stderr, "%4d of %d %q", i, len(corpus), filepath.Base(inputPath))
		fileInfo, err := fileImages(inputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, " ERROR: %v\n", inputPath, err)
			continue
		}
		corpusInfo[inputPath] = fileInfo
		fmt.Fprintf(os.Stderr, "\n")
	}

	showSummary(corpus, corpusInfo)
	saveAsCsv("results.csv", corpus, corpusInfo)
}

// fileImages returns a list of imageInfo entries for the images in the PDF file `inputPath`.
func fileImages(inputPath string) ([]imageInfo, error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return nil, err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return nil, err
	}

	if isEncrypted {
		// Try decrypting with an empty one.
		auth, err := pdfReader.Decrypt([]byte(""))
		if err != nil {
			return nil, err
		}
		if !auth {
			fmt.Println("Need to decrypt with a specified user/owner password")
			return nil, nil
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(os.Stderr, " %d pages,", numPages)

	var fileInfo []imageInfo

	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page, err := pdfReader.GetPage(pageNum)
		showError(nil, err, "pdfReader.GetPage failed: page %d", pageNum)
		if err != nil {
			continue
		}

		// List images on the page.
		pageInfo, err := pageImages(page)
		if err != nil || len(pageInfo) == 0 {
			continue
		}
		for i := range pageInfo {
			pageInfo[i].path = inputPath
			pageInfo[i].page = pageNum
		}
		fileInfo = append(fileInfo, pageInfo...)
	}

	fmt.Fprintf(os.Stderr, " %d images", len(fileInfo))

	return fileInfo, nil
}

// pageImages returns a list of imageInfo entries for the images in the PDF page `page`.
func pageImages(page *pdf.PdfPage) ([]imageInfo, error) {
	contents, err := page.GetAllContentStreams()
	if err != nil {
		return nil, err
	}
	return contentStreamImages(contents, page.Resources)
}

// errors records the errors seen so far. It is used to display each error only once.
var errors = map[error]bool{nil: true}

// contentStreamImages returns a list of imageInfo entries for the images in the content stream `contents`.
func contentStreamImages(contents string, resources *pdf.PdfPageResources) ([]imageInfo, error) {
	cstreamParser := pdfcontent.NewContentStreamParser(contents)
	operations, err := cstreamParser.Parse()
	showError(errors, err, "cstreamParser.Parse failed")
	if err != nil {
		return nil, err
	}

	var infoList []imageInfo

	processedXObjects := map[string]bool{}

	for _, op := range *operations {
		if op.Operand == "BI" && len(op.Params) == 1 {
			// Inline image.
			iimg, ok := op.Params[0].(*pdfcontent.ContentStreamInlineImage)
			if !ok {
				continue
			}

			var width, height, cpts, bpc int
			img, err := iimg.ToImage(resources)
			showError(errors, err, "ToImage failed")
			if err == nil {
				width = int(img.Width)
				height = int(img.Height)
				cpts = img.ColorComponents
				bpc = int(img.BitsPerComponent)
			}

			var filter, colorspace string
			cs, err := iimg.GetColorSpace(resources)
			showError(errors, err, "GetColorSpace failed")
			if err == nil {
				colorspace = cs.String()
			}
			encoder, err := iimg.GetEncoder()
			showError(errors, err, "GetEncoder failed")
			if err == nil {
				filter = encoder.GetFilterName()
			}

			info := imageInfo{
				inline:     true,
				filter:     filter,
				width:      width,
				height:     height,
				cpts:       cpts,
				colorspace: colorspace,
				bpc:        bpc,
			}

			infoList = append(infoList, info)

		} else if op.Operand == "Do" && len(op.Params) == 1 {
			// XObject.
			name := op.Params[0].(*pdfcore.PdfObjectName)

			// Only process each one once.
			if _, has := processedXObjects[string(*name)]; has {
				continue
			}
			processedXObjects[string(*name)] = true

			_, xtype := resources.GetXObjectByName(*name)
			if xtype == pdf.XObjectTypeImage {

				ximg, err := resources.GetXObjectImageByName(*name)
				showError(errors, err, "GetXObjectImageByName failed: %q ", *name)
				if err != nil {
					continue
				}

				var width, height, cpts, bpc int
				img, err := ximg.ToImage()
				showError(errors, err, "ximg.ToImage failed: %q ", *name)
				if err == nil {
					cpts = img.ColorComponents
				}
				if ximg.Width != nil {
					width = int(*ximg.Width)
				}
				if ximg.Height != nil {
					height = int(*ximg.Height)
				}
				if ximg.BitsPerComponent != nil {
					bpc = int(*ximg.BitsPerComponent)
				}

				info := imageInfo{
					inline:     false,
					filter:     ximg.Filter.GetFilterName(),
					width:      width,
					height:     height,
					cpts:       cpts,
					colorspace: ximg.ColorSpace.String(),
					bpc:        bpc,
				}
				infoList = append(infoList, info)

			} else if xtype == pdf.XObjectTypeForm {
				// Go through the XObject Form content stream.
				xform, err := resources.GetXObjectFormByName(*name)
				showError(errors, err, "GetXObjectFormByName failed: %q", *name)
				if err != nil {
					continue
				}
				formContent, err := xform.GetContentStream()
				showError(errors, err, "GetContentStream failed: %q", *name)
				if err != nil {
					continue
				}

				// Process the content stream in the Form object too.
				formResources := xform.Resources
				if formResources == nil {
					formResources = resources
				}
				formDescs, err := contentStreamImages(string(formContent), formResources)
				showError(errors, err, "contentStreamImages failed: %q", *name)
				if err != nil {
					continue
				}
				infoList = append(infoList, formDescs...)
			}
		}
	}

	return infoList, nil
}

type imageInfo struct {
	path       string
	page       int
	inline     bool
	filter     string
	width      int
	height     int
	cpts       int
	colorspace string
	bpc        int
}

func (info imageInfo) String() string {
	name := "XObject"
	if info.inline {
		name = "Inline image"
	}
	return strings.Join([]string{
		fmt.Sprintf("%q:%d", filepath.Base(info.path), info.page),
		fmt.Sprintf("%s", name),
		fmt.Sprintf("  Filter: %s", info.filter),
		fmt.Sprintf("  Width: %d", info.width),
		fmt.Sprintf("  Height: %d", info.height),
		fmt.Sprintf("  Color components: %d", info.cpts),
		fmt.Sprintf("  ColorSpace: %s", info.colorspace),
		fmt.Sprintf("  BPC: %d", info.bpc),
	}, "\n")
}

func (info imageInfo) asStrings() []string {
	name := "XObject"
	if info.inline {
		name = "Inline image"
	}
	parts := []string{
		info.path,
		fmt.Sprintf("%d", info.page),
		name,
		info.filter,
		fmt.Sprintf("%d", info.width),
		fmt.Sprintf("%d", info.height),
		fmt.Sprintf("%d", info.cpts),
		info.colorspace,
		fmt.Sprintf("%d", info.bpc),
	}
	if len(parts) != len(header) {
		panic("csv")
	}
	return parts
}

var header = []string{
	"Path",
	"Page number",
	"Type",
	"Filter",
	"Width",
	"Height",
	"Cpts",
	"Colors Space",
	"BPC",
}

// saveAsCsv saves `fileInfo` as a CSV file.
func saveAsCsv(filename string, corpus []string, corpusInfo map[string][]imageInfo) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	err = w.Write(header)
	if err != nil {
		return err
	}

	for _, fn := range corpus {
		infoList, ok := corpusInfo[fn]
		if !ok {
			continue
		}
		for _, info := range infoList {
			err := w.Write(info.asStrings())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func showSummary(corpus []string, corpusInfo map[string][]imageInfo) {
	numFiles := len(corpusInfo)
	numImages := sumVals(corpusInfo)
	fmt.Println("=================================================")
	fmt.Printf("Totals:%d of files contain images. %6d images\n", numFiles, len(corpus), numImages)
	boolSummary("inline", corpusInfo, func(info imageInfo) bool { return info.inline })
	stringSummary("filter", corpusInfo, func(info imageInfo) string { return info.filter })
	stringSummary("color", corpusInfo, func(info imageInfo) string { return info.colorspace })
	intSummary("cpts", corpusInfo, func(info imageInfo) int { return info.cpts })
	intSummary("bpc", corpusInfo, func(info imageInfo) int { return info.bpc })
	// intSummary("width", corpusInfo, func(info imageInfo) int { return info.width })
	// intSummary("height", corpusInfo, func(info imageInfo) int { return info.height })
}

func boolSummary(title string, corpusInfo map[string][]imageInfo, selector func(imageInfo) bool) {
	numFiles := len(corpusInfo)
	numImages := sumVals(corpusInfo)
	byImage, byFile := boolCounts(corpusInfo, selector)
	imageKeys, fileKeys := boolKeys(byImage), boolKeys(byFile)
	fmt.Println("-----------------------------------------")
	fmt.Printf("%s\n", title)
	fmt.Printf("By image: %d\n", len(byImage))
	for _, k := range imageKeys {
		fmt.Printf("\t%+15v\t%s\n", k, percentage(byImage[k], numImages))
	}
	fmt.Printf("By file: %d\n", len(byFile))
	for _, k := range fileKeys {
		fmt.Printf("\t%+15v\t%s\n", k, percentage(byFile[k], numFiles))
	}
}

func intSummary(title string, corpusInfo map[string][]imageInfo, selector func(imageInfo) int) {
	numFiles := len(corpusInfo)
	numImages := sumVals(corpusInfo)
	byImage, byFile := intCounts(corpusInfo, selector)
	imageKeys, fileKeys := intKeys(byImage), intKeys(byFile)
	fmt.Println("-----------------------------------------")
	fmt.Printf("%s\n", title)
	fmt.Printf("By image: %d\n", len(byImage))
	for _, k := range imageKeys {
		fmt.Printf("\t%+15v\t%s\n", k, percentage(byImage[k], numImages))
	}
	fmt.Printf("By file: %d\n", len(byFile))
	for _, k := range fileKeys {
		fmt.Printf("\t%+15v\t%s\n", k, percentage(byFile[k], numFiles))
	}
}

func stringSummary(title string, corpusInfo map[string][]imageInfo, selector func(imageInfo) string) {
	numFiles := len(corpusInfo)
	numImages := sumVals(corpusInfo)
	byImage, byFile := stringCounts(corpusInfo, selector)
	imageKeys, fileKeys := stringKeys(byImage), stringKeys(byFile)
	fmt.Println("-----------------------------------------")
	fmt.Printf("%s\n", title)
	fmt.Printf("By image: %d\n", len(byImage))
	for _, k := range imageKeys {
		fmt.Printf("\t%+15v\t%s\n", k, percentage(byImage[k], numImages))
	}
	fmt.Printf("By file: %d\n", len(byFile))
	for _, k := range fileKeys {
		fmt.Printf("\t%+15v\t%s\n", k, percentage(byFile[k], numFiles))
	}
}

func boolKeys(counts map[bool]int) []bool {
	var keys []bool
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		ki, kj := keys[i], keys[j]
		ni, nj := counts[ki], counts[kj]
		if ni != nj {
			return ni > nj
		}
		return kj
	})
	return keys
}

func intKeys(counts map[int]int) []int {
	var keys []int
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		ki, kj := keys[i], keys[j]
		ni, nj := counts[ki], counts[kj]
		if ni != nj {
			return ni > nj
		}
		return ki < kj
	})
	return keys
}

func stringKeys(counts map[string]int) []string {
	var keys []string
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		ki, kj := keys[i], keys[j]
		ni, nj := counts[ki], counts[kj]
		if ni != nj {
			return ni > nj
		}
		return ki < kj
	})
	return keys
}

func boolCounts(corpusInfo map[string][]imageInfo, selector func(imageInfo) bool) (
	map[bool]int, map[bool]int) {
	byImage := map[bool]int{}
	byFile := map[bool]int{}
	for _, infoList := range corpusInfo {
		vals := map[bool]bool{}
		for _, info := range infoList {
			byImage[selector(info)] += 1
			vals[selector(info)] = true
		}
		for v := range vals {
			byFile[v] += 1
		}
	}
	return byImage, byFile
}

func intCounts(corpusInfo map[string][]imageInfo, selector func(imageInfo) int) (
	map[int]int, map[int]int) {
	byImage := map[int]int{}
	byFile := map[int]int{}
	for _, infoList := range corpusInfo {
		vals := map[int]bool{}
		for _, info := range infoList {
			byImage[selector(info)] += 1
			vals[selector(info)] = true
		}
		for v := range vals {
			byFile[v] += 1
		}
	}
	return byImage, byFile
}

func stringCounts(corpusInfo map[string][]imageInfo, selector func(imageInfo) string) (
	map[string]int, map[string]int) {
	byImage := map[string]int{}
	byFile := map[string]int{}
	for _, infoList := range corpusInfo {
		vals := map[string]bool{}
		for _, info := range infoList {
			byImage[selector(info)] += 1
			vals[selector(info)] = true
		}
		for v := range vals {
			byFile[v] += 1
		}
	}
	return byImage, byFile
}

// showError prints an error message `format` for error `err` if `err` has not been reported before.
// `errors` tracks errors seen so far. The caller can make `errors` per-page, per-file or global.
func showError(errors map[error]bool, err error, format string, args ...interface{}) bool {
	seen := false
	if errors != nil {
		_, ok := errors[err]
		seen = seen || ok
	}
	if seen && err != nil {
		msg := fmt.Sprintf(format, args...)
		fmt.Printf("%s. err=%v\n", msg, err)
	}
	if errors != nil {
		errors[err] = true
	}
	return err != nil
}

func sumVals(corpusInfo map[string][]imageInfo) int {
	n := 0
	for _, info := range corpusInfo {
		n += len(info)
	}
	return n
}

func percentage(n, total int) string {
	perc := 0.0
	if total > 0 {
		perc = 100.0 * float64(n) / float64(total)
	}
	return fmt.Sprintf("%6d of %d (%4.1f%%)", n, total, perc)
}

// fileSizeMB returns the size of file `path` in megabytes.
func fileSizeMB(path string) float64 {
	fi, err := os.Stat(path)
	if err != nil {
		return -1.0
	}
	return float64(fi.Size()) / 1024.0 / 1024.0
}

// makeUsage updates flag.Usage to include usage message `msg`.
func makeUsage(msg string) {
	usage := flag.Usage
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, msg)
		usage()
	}
}