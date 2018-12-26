/*
 * Extract images from a PDF file. Passes through each page, goes through the content stream and
 * finds instances of both XObject Images and inline images. Also handles images referred within
 * XObject Form content streams.
 * The output files are saved as a zip archive.
 *
 * NOTE(peterwilliams97): Unlike pdf_extract_images.go,
 *       1) Losslessly compressed PDF images are saved in PNG format. (Lossily compressed PDF images
 *          still are saved in JPEG format)
 *       2) Images are saved in the same color space as they occur in PDF files.
 *
 * XXX(peterwilliams97): This file fixes an apparent problem in the UniDoc resampling code with
 *       handling 1 bit per component images. There is an additional problem that got6.DecodeBytes()
 *       returns CCITTFax images as 8 bits per pixel while PDF expects these images to be 1 bit per
 *       pixel. I tried modifying got.6 to fix this and found that ResampleBytes() didn't work with
 *       1 bit per pixel, so I just set imgData.img.BitsPerComponent = 8 for CCITTFaxEncoder images.
 *
 * TODO: Handle JBIG images.
 *       Handle CCITTFaxEncoder inline images?
 *       Handle CCITTFaxEncoder Group 3 images?
 *       Change got.6 to return 1 bit images.
 *       Save images in orientation they appear in the PDF file.
 *
 * Run as: go run pdf_extract_native_images.go input.pdf output.zip
 */

package main

import (
	"flag"
	"fmt"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/unidoc/unidoc/common"
	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

var xObjectImages = 0
var inlineImages = 0

const usage = "Usage: go run pdf_render_text.go testdata/*.pdf\n"

type details struct {
	inline bool
	filter string
	bpc    int
	colors int
	width  int
	height int
}

func main() {

	var showHelp, debug, trace, verbose bool
	var outputDir string
	var minDim int
	flag.BoolVar(&showHelp, "h", false, "Show this help message.")
	flag.BoolVar(&debug, "d", false, "Print debugging information.")
	flag.BoolVar(&trace, "e", false, "Print detailed debugging information.")
	flag.BoolVar(&verbose, "v", false, "Print extra page information.")
	flag.StringVar(&outputDir, "o", "", "Output directory.")
	flag.IntVar(&minDim, "m", 0, "Mininum image width and height to save.")
	makeUsage(usage)

	flag.Parse()
	args := flag.Args()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}
	if len(args) < 1 || outputDir == "" {
		flag.Usage()
		os.Exit(1)
	}

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.Mkdir(outputDir, 0777)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not create outputDir=%q err=%v\n", outputDir, err)
			os.Exit(1)
		}
	}

	if trace {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelTrace))
	} else if debug {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
	} else {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelError))
	}

	files := args[:]
	sort.Strings(files)
	sort.Slice(files, func(i, j int) bool {
		fi, fj := files[i], files[j]
		si, sj := fileSizeMB(fi), fileSizeMB(fj)
		if si != sj {
			return si < sj
		}
		return fi < fj
	})

	numFiles := 0

	for _, inputPath := range files {
		if !isWanted(inputPath) && len(files) > 1 {
			continue
		}
		sizeMB := fileSizeMB(inputPath)
		if !(minSizeMB <= sizeMB && sizeMB <= maxSizeMB) && len(files) > 1 {
			// fmt.Fprintf(os.Stderr, "%.3f MB, ", sizeMB)
			continue
		}
		if numFiles > maxFiles && len(files) > 1 {
			break
		}
		if verbose {
			fmt.Println("========================= ^^^ =========================")
		}
		t0 := time.Now()
		pdfReader, numPages, err := getReader(inputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\to====> Pdf File %3d of %d %q err=%v\n",
				numFiles, len(files), inputPath, err)
			continue
		}

		version := pdfReader.PdfVersion()

		fmt.Fprintf(os.Stderr, "Pdf File %3d of %d (%3s) %5.2f MB %3d pages %q ",
			numFiles, len(files), pdfReader.PdfVersion(), sizeMB, numPages, inputPath)
		if version.Minor < minVersionMinor && len(files) > 1 {
			fmt.Fprintln(os.Stderr, "")
			continue
		}

		docInfo, err := extractImagesToArchive(inputPath, pdfReader, outputDir, minDim)
		dt := time.Since(t0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "FAILED err=%v\n", err)
			continue
		}

		if summary := infoSummary(numPages, docInfo); len(summary) > 0 {
			fmt.Printf("%s\n", summary)
		}

		fmt.Fprintf(os.Stderr, "%3.1f sec\n", dt.Seconds())
		if err != nil {
			fmt.Fprintf(os.Stderr, "\tx====> Pdf File %3d of %d %q err=%v\n",
				numFiles, len(files), inputPath, err)
		}
		if verbose {
			fmt.Println("========================= ~~~ =========================")
		}

		numFiles++
	}
	fmt.Fprintf(os.Stderr, "Done %d files \n", numFiles)

	fmt.Fprintf(os.Stderr, "-- Summary\n")
	fmt.Fprintf(os.Stderr, "%d XObject images extracted\n", xObjectImages)
	fmt.Fprintf(os.Stderr, "%d inline images extracted\n", inlineImages)
	fmt.Fprintf(os.Stderr, "Total %d images\n", xObjectImages+inlineImages)
}

type pageInfo struct {
	pageNum    int
	imageDescs []string
}

func infoSummary(numPages int, docInfo []pageInfo) string {
	if len(docInfo) == 0 {
		return ""
	}
	parts := []string{fmt.Sprintf("%d of %d pages contain images", len(docInfo), numPages)}
	for _, info := range docInfo {
		parts = append(parts, fmt.Sprintf("\tPage %d: %d images", info.pageNum, len(info.imageDescs)))
		for i, desc := range info.imageDescs {
			parts = append(parts, fmt.Sprintf("\t\t %2d: %s", i+1, desc))
		}
	}
	return strings.Join(parts, "\n")
}

func makePath(inputPath, outputDir, fname string) string {
	base := filepath.Base(inputPath)
	base = strings.TrimSuffix(base, filepath.Ext(base))
	name := strings.Join([]string{base, fname}, "-")
	out := filepath.Join(outputDir, name)
	return out
}

// extractImagesToArchive extracts images and properties of a PDF specified by `inputPath`.
// The output images are stored into a zip archive whose path is given by `outputDir`.
func extractImagesToArchive(inputPath string, pdfReader *pdf.PdfReader, outputDir string,
	minDim int) ([]pageInfo, error) {
	var docInfo []pageInfo

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return docInfo, err
	}

	for i := 0; i < numPages; i++ {
		info := pageInfo{pageNum: i + 1}

		page, err := pdfReader.GetPage(i + 1)
		if err != nil {
			return docInfo, err
		}

		// List images on the page.
		images, err := extractImagesOnPage(page, minDim)
		if err != nil {
			return docInfo, err
		}

		for idx, imgData := range images {
			img := imgData.img
			fname := fmt.Sprintf("p%d_%d", i+1, idx)

			var lossless bool // Is image compressed losslessly?

			switch imgData.filter.(type) {
			case *pdfcore.FlateEncoder:
				lossless = true
			case *pdfcore.CCITTFaxEncoder:
				lossless = true
				// XXX(peterwilliams97) Hack to work around got6.DecodeBytes() returning an 8 bits
				// per component raster and sampling.ResampleBytes() not working for 1 bits per
				// pixel
				imgData.img.BitsPerComponent = 8
			}

			gimg, err := img.ToGoImage()
			if err != nil {
				return docInfo, err
			}

			if lossless {
				fname += ".png"
			} else {
				fname += ".jpg"
			}

			fname = makePath(inputPath, outputDir, fname)

			info.imageDescs = append(info.imageDescs, fmt.Sprintf("%q %s", fname, img))

			imgf, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0666)
			if err != nil {
				return docInfo, err
			}
			defer imgf.Close()

			if lossless {
				err = png.Encode(imgf, gimg)
			} else {
				opt := jpeg.Options{Quality: 100}
				err = jpeg.Encode(imgf, gimg, &opt)
			}
			if err != nil {
				return docInfo, err
			}
		}
		if len(info.imageDescs) > 0 {
			docInfo = append(docInfo, info)
		}
	}

	return docInfo, nil
}

type imageData struct {
	img    *pdf.Image
	filter pdfcore.StreamEncoder
}

// extractImagesOnPage returns a slice of all images on page `page`.
func extractImagesOnPage(page *pdf.PdfPage, minDim int) ([]imageData, error) {
	contents, err := page.GetAllContentStreams()
	if err != nil {
		return nil, err
	}
	return extractImagesInContentStream(contents, page.Resources, minDim)
}

// extractImagesInContentStream returns a slice of all images in content stream `contents`.
func extractImagesInContentStream(contents string, resources *pdf.PdfPageResources,
	minDim int) ([]imageData, error) {
	images := []imageData{}
	cstreamParser := pdfcontent.NewContentStreamParser(contents)
	operations, err := cstreamParser.Parse()
	if err != nil {
		return nil, err
	}

	processedXObjects := map[string]bool{}

	// Range through all the content stream operations.
	for _, op := range *operations {
		if op.Operand == "BI" && len(op.Params) == 1 {
			// BI: Inline image.

			iimg, ok := op.Params[0].(*pdfcontent.ContentStreamInlineImage)
			if !ok {
				continue
			}

			img, err := iimg.ToImage(resources)
			if err != nil {
				return nil, err
			}

			cs, err := iimg.GetColorSpace(resources)
			if err != nil {
				return nil, err
			}
			if cs == nil {
				// Default if not specified?
				cs = pdf.NewPdfColorspaceDeviceGray()
			}

			images = append(images, imageData{img: img})
			inlineImages++
		} else if op.Operand == "Do" && len(op.Params) == 1 {
			// Do: XObject.
			name := op.Params[0].(*pdfcore.PdfObjectName)

			// Only process each one once.
			_, has := processedXObjects[string(*name)]
			if has {
				continue
			}
			processedXObjects[string(*name)] = true

			_, xtype := resources.GetXObjectByName(*name)
			if xtype == pdf.XObjectTypeImage {
				fmt.Printf(" XObject Image: %s\n", *name)

				ximg, err := resources.GetXObjectImageByName(*name)
				if err != nil {
					return nil, err
				}

				fmt.Printf("ximg=%s\n", ximg)
				fmt.Printf("ximg.Filter=%T=%s\n", ximg.Filter, ximg.Filter.GetFilterName())

				img, err := ximg.ToImage()
				if err != nil {
					return nil, err
				}
				if int(img.Width) >= minDim && int(img.Height) >= minDim {
					images = append(images, imageData{img, ximg.Filter})
				}
				xObjectImages++
			} else if xtype == pdf.XObjectTypeForm {
				// Go through the XObject Form content stream.
				xform, err := resources.GetXObjectFormByName(*name)
				if err != nil {
					return nil, err
				}

				formContent, err := xform.GetContentStream()
				if err != nil {
					return nil, err
				}

				// Process the content stream in the Form object too:
				formResources := xform.Resources
				if formResources == nil {
					formResources = resources
				}

				// Process the content stream in the Form object too:
				formImages, err := extractImagesInContentStream(string(formContent), formResources, minDim)
				if err != nil {
					return nil, err
				}
				images = append(images, formImages...)
			}
		}
	}

	return images, nil
}

// makeUsage updates flag.Usage to include usage message `msg`.
func makeUsage(msg string) {
	usage := flag.Usage
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, msg)
		usage()
	}
}

// getReader returns a PdfReader and the number of pages for PDF file `inputPath`.
func getReader(inputPath string) (pdfReader *pdf.PdfReader, numPages int, err error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return nil, 0, err
	}
	defer f.Close()

	pdfReader, err = pdf.NewPdfReader(f)
	if err != nil {
		return nil, 0, err
	}
	numPages, err = pdfReader.GetNumPages()
	return pdfReader, numPages, err
}

// fileSizeMB returns the size of file `path` in megabytes
func fileSizeMB(path string) float64 {
	fi, err := os.Stat(path)
	if err != nil {
		return -1.0
	}
	return float64(fi.Size()) / 1024.0 / 1024.0
}

// isWanted is for customising test runs to include desired files.
// It should return true for the files you want to process.
// e.g.The commented core returns true for files containing Type0 font dicts in clear text.
func isWanted(filename string) bool {
	for _, s := range exclusions {
		if strings.Contains(filename, s) {
			return false
		}
	}
	return true
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return strings.Contains(string(data), "/TrueType")
	return (strings.Contains(string(data), "/Type1") &&
		!strings.Contains(string(data), "/Type0") &&
		!strings.Contains(string(data), "/Type1C"))
}

var (
	maxFiles        = 1000000
	minSizeMB       = 0.60
	maxSizeMB       = 1.0e20
	minVersionMinor = 3
	exclusions      = []string{
		// Report to UniDoc
		"Presto_UserGuide.pdf",
		"constrained_decoding.pdf",
		"2_LarryTalk.pdf",
	}
)
