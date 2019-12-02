/*
 * PDF optimization (compression) example.
 *
 * Run as: go run pdf_optimize.go <input.pdf> <output.pdf>
 */
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/bmatcuk/doublestar"
	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/contentstream"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/model/optimize"
)

const usage = "Usage: %s INPUT_PDF_PATH OUTPUT_PDF_PATH\n"

func main() {
	var debug, trace bool
	var simpleColor bool
	outDir := "output"
	fileSort := ""
	flag.BoolVar(&debug, "d", false, "Print debugging information.")
	flag.BoolVar(&trace, "e", false, "Print detailed debugging information.")
	flag.BoolVar(&simpleColor, "c", false, "Convert ICC color images to grayscale or RGB.")
	flag.StringVar(&outDir, "o", outDir, "Output directory.")
	flag.StringVar(&fileSort, "s", fileSort, "File sort. a for ascending size, d for descending size")
	makeUsage(usage)

	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}
	if len(fileSort) > 0 {
		fileSort = strings.ToLower(fileSort[:1])
		if fileSort != "a" && fileSort != "d" {
			flag.Usage()
			os.Exit(1)
		}
	}

	if trace {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelTrace))
	} else if debug {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
	} else {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelInfo))
	}

	err := os.MkdirAll(outDir, 0777)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}

	pathList, err := patternsToPaths(args)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}
	if fileSort != "" {
		decending := fileSort == "d"
		sort.Slice(pathList, func(i, j int) bool {
			pi, pj := pathList[i], pathList[j]
			si, sj := fileSizeMB(pi), fileSizeMB(pj)
			if si != sj {
				if decending {
					return si > sj
				} else {
					return si < sj
				}
			}
			return pi < pj
		})
	}

	pathList = cleanCorpus(pathList)

	var results []optResult
	iStart := 5090
	for i, inPath := range pathList {
		if len(pathList) >= iStart && i+1 < iStart && fileSort != "d" {
			continue
		}
		fmt.Printf("-Original file: %s\n", inPath)
		name := filepath.Base(inPath)
		outPath := filepath.Join(outDir, name)
		start := time.Now()
		skip, err := optimizePDF(inPath, outPath, simpleColor)
		if err != nil {
			if skip {
				fmt.Fprintf(os.Stderr, "%4d of %d: skipping %s err=%v\n",
					i+1, len(pathList), inPath, err)
				continue
			}
			log.Fatalf("Fail: %d %q err=%v\n", i, inPath, err)
		}
		duration := time.Since(start).Seconds()

		r := optResult{
			inPath:    inPath,
			outPath:   outPath,
			inSizeMB:  fileSizeMB(inPath),
			outSizeMB: fileSizeMB(outPath),
			duration:  duration,
		}
		results = append(results, r)

		fmt.Printf(" Original file: %s\n", inPath)
		fmt.Printf("Optimized file: %s\n", outPath)
		fmt.Fprintf(os.Stderr, "%4d of %d: %s\n", i+1, len(pathList), r)
	}

	fmt.Fprintln(os.Stderr, "=================================================")
	fmt.Fprintf(os.Stderr, "%d of %d files tested\n", len(results), len(pathList))
	if len(results) == 0 {
		return
	}
	r := arithMean(results)
	fmt.Fprintf(os.Stderr, "    Average: %s\n", r) 
}

type optResult struct {
	inPath    string
	outPath   string
	inSizeMB  float64
	outSizeMB float64
	duration  float64
}

func (r optResult) String() string {
	ratio := r.outSizeMB / r.inSizeMB
	return fmt.Sprintf("%6.3f->%6.3f MB %5.2f%% %5.2f sec %s",
		r.inSizeMB, r.outSizeMB, ratio, r.duration, r.inPath)
}

func arithMean(results []optResult) optResult {
	if len(results) == 0 {
		return optResult{}
	}
	totalInSize := 0.0
	totalOutSize := 0.0
	totalDuration := 0.0
	for _, r := range results {
		totalInSize += r.inSizeMB
		totalOutSize += r.outSizeMB
		totalDuration += r.duration
	}
	return optResult{
		inSizeMB:  totalInSize / float64(len(results)),
		outSizeMB: totalOutSize / float64(len(results)),
		duration:  totalDuration / float64(len(results)),
	}
}

func optimizePDF(inPath, outPath string, simpleColor bool) (bool, error) {
	tmpPath := outPath + ".tmp.pdf"
	skip, err := optimizePDF_(inPath, tmpPath, simpleColor)
	if err != nil {
		common.Log.Info("err=%v deleting %q", err, tmpPath)
		err2 := os.Remove(tmpPath)
		if err2 != nil {
			common.Log.Error("os.Remove(%s) failed. err=%v", tmpPath, err)
		}
		return true, err
	}
	err = os.Rename(tmpPath, outPath)
	if err != nil {
		common.Log.Error("os.Rename(%s, %s) failed.derr=%v", tmpPath, outPath, err)
		return false, err
	}
	return skip, err
}

// optimizePDF reduces the size of PDF `inPath` and writes the result to `outPath`. If
// simpleColor is true, ICC color images are converted to DeviceGray or DeviceRGB.
func optimizePDF_(inPath, outPath string, simpleColor bool) (bool, error) {
	// Create reader.
	inputFile, err := os.Open(inPath)
	if err != nil {
		return false, err
	}
	defer inputFile.Close()

	pdfReader, err := model.NewPdfReader(inputFile)
	if err != nil {
		return true, err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		panic(err)
		return false, err
	}

	if isEncrypted {
		// Decrypt if needed.  Put your password in the empty string below.
		auth, err := pdfReader.Decrypt([]byte(""))
		if err != nil {
			return false, err
		}
		if !auth {
			fmt.Fprintf(os.Stderr, "%s - Unable to access (encrypted)\n", inPath)
			return true, nil
		}
	}

	// Get number of pages in the input file.
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return true, err
	}

	// Add input file pages to the writer.
	writer := model.NewPdfWriter()
	for pageNum := 1; pageNum <= numPages; pageNum++ {
		// common.Log.Info("page %d", pageNum)
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return false, err
		}
		if simpleColor {
			err := pageToSimpleColor(page)
			if err != nil {
				common.Log.Info("Could not convert %q: page %d err=%v. Reverting to original",
					inPath, pageNum, err)
				page, _ = pdfReader.GetPage(pageNum)
			}
		}
		if err = writer.AddPage(page); err != nil {
			return false, err
		}
	}

	// Add pdfReader AcroForm to the writer.
	if pdfReader.AcroForm != nil {
		writer.SetForms(pdfReader.AcroForm)
	}

	// Set optimizer.
	writer.SetOptimizer(optimize.New(optimize.Options{
		CombineDuplicateDirectObjects:   true,
		CombineIdenticalIndirectObjects: true,
		CombineDuplicateStreams:         true,
		CompressStreams:                 true,
		UseObjectStreams:                true,
		ImageQuality:                    80,
		ImageUpperPPI:                   100,
	}))

	// Create output file.
	outputFile, err := os.Create(outPath)
	if err != nil {
		common.Log.Error("Create failed. outPath=%q err=%v", outPath, err)
		return false, err
	}
	defer outputFile.Close()

	common.Log.Info("Write: inPath=%q", inPath)

	// Write output file.
	err = writer.Write(outputFile)
	if err != nil {
		common.Log.Error("Write failed. inPath=%q err=%v", inPath, err)
		return false, err
	}
	return false, nil
}

// pageToSimpleColor goes through all XObject images referenced by `page` and converts 1 and
// 3 component ICC images to RGB to grayscale and RGB respectively.
func pageToSimpleColor(page *model.PdfPage) error {
	contents, err := page.GetAllContentStreams()
	if err != nil {
		return err
	}
	// fmt.Printf("--------------------\n%s\n--------------------\n", contents)
	return imagesToSimpleColor(contents, page.Resources)
}

// pageToSimpleColor goes through all XObject images in `contents` and converts 1 and 3
// component ICC images to grayscale and RGB respectively.
func imagesToSimpleColor(contents string, resources *model.PdfPageResources) error {
	cstreamParser := contentstream.NewContentStreamParser(contents)
	operations, err := cstreamParser.Parse()
	if err != nil {
		return err
	}

	// Keep track of processed XObjects to avoid repetition.
	processedXObjects := map[string]struct{}{}

	processor := contentstream.NewContentStreamProcessor(*operations)
	// Handler for XObject Images.
	processor.AddHandler(contentstream.HandlerConditionEnumOperand, "Do",
		func(op *contentstream.ContentStreamOperation, gs contentstream.GraphicsState,
			resources *model.PdfPageResources) error {
			if len(op.Params) < 1 {
				return fmt.Errorf("Invalid number of params for Do object. op=%v", op)
			}
			// XObject.
			nameObj, ok := core.GetName(op.Params[0])
			if !ok {
				return fmt.Errorf("Invalid type for Do object name. op=%v", op)
			}
			name := string(*nameObj)

			// Only process each one once.
			_, has := processedXObjects[name]
			if has {
				return nil
			}
			processedXObjects[name] = struct{}{}

			// We only process images
			if _, xtype := resources.GetXObjectByName(*nameObj); xtype != model.XObjectTypeImage {
				return nil
			}

			ximg, err := resources.GetXObjectImageByName(*nameObj)
			if err != nil {
				return err
			}

			// if *ximg.BitsPerComponent != 8 {
			// 	return nil
			// }
			// We currently only process 1 and 3 component ICC images.
			cs := ximg.ColorSpace
			// common.Log.Info("ColorSpace=%+q cpts=%d", cs.String(), cs.GetNumComponents())
			convertable := cs.String() == "ICCBased" && (cs.GetNumComponents() == 3 || cs.GetNumComponents() == 1)
			convertable = convertable || cs.String() == "CalGray"
			if !convertable {
				return nil
			}

			common.Log.Debug("Converting 3 cpt ICC to RGB")
			img, err := ximg.ToImage()
			if err != nil {
				return fmt.Errorf("imagesToSimpleColor: err=%v", err)
			}
			var rgbImg model.Image
			if cs.String() == "CalGray" {
				rgbImg = *img
			} else {
				rgbImg, err = ximg.ColorSpace.ImageToRGB(*img)
				if err != nil {
					return fmt.Errorf("imagesToSimpleColor: err=%v", err)
				}
				if cs.GetNumComponents() == 1 {
					rgbColorSpace := model.NewPdfColorspaceDeviceRGB()
					rgbImg, err = rgbColorSpace.ImageToGray(rgbImg)
					if err != nil {
						return fmt.Errorf("imagesToSimpleColor: err=%v", err)
					}
				}
			}
			// Update the XObject image.
			// Use same encoder as input data.  Make sure for DCT filter it is updated to 1 color component.
			encoder := ximg.Filter
			// common.Log.Info("encoder=%T %+q", encoder, encoder.GetFilterName())

			if dctEncoder, is := encoder.(*core.DCTEncoder); is {
				common.Log.Info("DCTEncoder: %d->%d", dctEncoder.ColorComponents, cs.GetNumComponents())
				dctEncoder.ColorComponents = cs.GetNumComponents()
			}
			ximgRGB, err := model.NewXObjectImageFromImage(&rgbImg, nil, encoder)
			if err != nil {
				// if err == core.ErrUnsupportedEncodingParameters {
				// 	// Unsupported encoding parameters, revert to a basic flate encoder without predictor.
				// 	err2 := fmt.Errorf("imagesToSimpleColor: Error creating image encoder=%#v err=%v",
				// 		encoder, err)
				// 	return err2
				// 	encoder = core.NewFlateEncoder()
				// }
				// // Try again, fail if error.
				// ximgRGB, err = model.NewXObjectImageFromImage(&rgbImg, nil, encoder)
				// if err != nil {
				return fmt.Errorf("imagesToSimpleColor: Error creating image err=%v", err)
				// }
			}
			// Update the entry.
			err = resources.SetXObjectImageByName(*nameObj, ximgRGB)
			if err != nil {
				return fmt.Errorf("imagesToSimpleColor: Failed setting XObject %+q. err=%v",
					name, err)
			}

			return nil
		})

	return processor.Process(resources)
}

// patternsToPaths returns a list of files matching the patterns in `patternList`.
func patternsToPaths(patternList []string) ([]string, error) {
	common.Log.Info("patternList=%d", len(patternList))
	var pathList []string

	for i, pattern := range patternList {
		pattern = expandUser(pattern)
		files, err := doublestar.Glob(pattern)
		if err != nil {
			common.Log.Error("patternsToPaths: Glob failed. pattern=%#q err=%v", pattern, err)
			return pathList, err
		}
		common.Log.Debug("patternList[%d]=%q %d matches", i, pattern, len(files))
		for _, filename := range files {
			ok, err := regularFile(filename)
			if err != nil {
				return pathList, fmt.Errorf("patternsToPaths: regularFile failed. pattern=%#q err=%v",
					pattern, err)
			}
			if !ok {
				continue
			}
			pathList = append(pathList, filename)
		}
	}
	pathList = stringUniques(pathList)
	return pathList, nil
}

// homeDir is the current user's home directory.
var homeDir = getHomeDir()

// getHomeDir returns the current user's home directory.
func getHomeDir() string {
	usr, _ := user.Current()
	return usr.HomeDir
}

// expandUser returns `filename` with "~"" replaced with user's home directory.
func expandUser(filename string) string {
	return strings.Replace(filename, "~", homeDir, -1)
}

// regularFile returns true if file `filename` is a regular file.
func regularFile(filename string) (bool, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return false, err
	}
	return fi.Mode().IsRegular(), nil
}

func fileSizeMB(filename string) float64 {
	fi, err := os.Stat(filename)
	if err != nil {
		common.Log.Error("stat failed. filename=%q err=%v", filename, err)
		return -1
	}
	return float64(fi.Size()) / 1024.0 / 1024.0
}

// stringUniques returns the unique strings in `arr`.
func stringUniques(arr []string) []string {
	set := map[string]struct{}{}
	var uniques []string
	for _, s := range arr {
		if _, exists := set[s]; !exists {
			uniques = append(uniques, s)
		}
		set[s] = struct{}{}
	}
	return uniques
}

// makeUsage updates flag.Usage to include usage message `msg`.
func makeUsage(msg string) {
	usage := flag.Usage
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, msg)
		usage()
	}
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
	"1812.09449.pdf",                         // hangs
	"INF586.pdf",                             // hangs
	"commercial-invoice-template-230kb.pdf",  // r=invalid pad length
	"CGU_Motor_Vehicle_Insurance_Policy.pdf", // r=invalid pad length
	"Forerunner_230_OM_EN.pdf} r=invalid",    // r=invalid pad length
	"transitions_test.pdf",                   //required attribute missing (No /Type/Font )
	"page_tree_multiple_levels.pdf",          //required attribute missing
	// "book.pdf",                               //version not found

	// // !@#$
	"/Symbolics_Common_Lis",
	// "CAM_Low Back Pain",
	// "yangbio",
}
