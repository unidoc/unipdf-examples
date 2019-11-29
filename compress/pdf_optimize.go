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
	"time"

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
	flag.BoolVar(&debug, "d", false, "Print debugging information.")
	flag.BoolVar(&trace, "e", false, "Print detailed debugging information.")
	flag.BoolVar(&simpleColor, "c", false, "Convert ICC color images to grayscale or RGB.")
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

	inputPath := args[0]
	outputPath := args[1]
	err := optimizePDF(inputPath, outputPath, simpleColor)
	if err != nil {
		log.Fatal("Fail: %v\n", err)
	}
}

// optimizePDF reduces the size of PDF `inputPath` and writes the result to `outputPath`. If
// simpleColor is true, ICC color images are converted to DeviceGray or DeviceRGB.
func optimizePDF(inputPath, outputPath string, simpleColor bool) error {
	// Initialize starting time.
	start := time.Now()

	// Get input file stat.
	inputFileInfo, err := os.Stat(inputPath)
	if err != nil {
		return err
	}

	// Create reader.
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	reader, err := model.NewPdfReader(inputFile)
	if err != nil {
		return err
	}

	// Get number of pages in the input file.
	numPages, err := reader.GetNumPages()
	if err != nil {
		return err
	}

	// Add input file pages to the writer.
	writer := model.NewPdfWriter()
	for pageNum := 1; pageNum <= numPages; pageNum++ {
		common.Log.Debug("page %d", pageNum)
		page, err := reader.GetPage(pageNum)
		if err != nil {
			return err
		}
		if simpleColor {
			err = convertPageToSimpleColor(page)
			if err != nil {
				return err
			}
		}
		if err = writer.AddPage(page); err != nil {
			return err
		}
	}

	// Add reader AcroForm to the writer.
	if reader.AcroForm != nil {
		writer.SetForms(reader.AcroForm)
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
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Write output file.
	err = writer.Write(outputFile)
	if err != nil {
		return err
	}

	// Get output file stat.
	outputFileInfo, err := os.Stat(outputPath)
	if err != nil {
		return err
	}

	// Print basic optimization statistics.
	inputSize := inputFileInfo.Size()
	outputSize := outputFileInfo.Size()
	percentage := float64(outputSize) / float64(inputSize) * 100.0
	ratio := 100.0 - percentage
	duration := time.Since(start).Seconds()

	fmt.Printf(" Original file: %s\n", inputPath)
	fmt.Printf("Optimized file: %s\n", outputPath)
	fmt.Printf(" Original size: %d bytes = %5.2f MB\n", inputSize, mb(inputSize))
	fmt.Printf("Optimized size: %d bytes = %5.2f MB\n", outputSize, mb(outputSize))
	fmt.Printf("Compression ratio: %.2f%% = %.2f%% of original\n", ratio, percentage)
	fmt.Printf("Processing time: %.3f sec\n", duration)

	return nil
}

// convertPageToSimpleColor goes through all XObject images referenced by `page` and converts 1 and
// 3 component ICC images to RGB to grayscale and RGB respectively.
func convertPageToSimpleColor(page *model.PdfPage) error {
	contents, err := page.GetAllContentStreams()
	if err != nil {
		return err
	}
	return transformImagesToSimpleColor(contents, page.Resources)
}

// convertPageToSimpleColor goes through all XObject images in `contents` and converts 1 and 3
// component ICC images to grayscale and RGB respectively.
func transformImagesToSimpleColor(contents string, resources *model.PdfPageResources) error {
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

			// We currently only process 1 and 3 component ICC images.
			cs := ximg.ColorSpace
			common.Log.Debug("ColorSpace=%+q cpts=%d", cs.String(), cs.GetNumComponents())
			if !(cs.String() == "ICCBased" && (cs.GetNumComponents() == 3 || cs.GetNumComponents() == 1)) {
				return nil
			}

			common.Log.Debug("Converting 3 cpt ICC to RGB")
			img, err := ximg.ToImage()
			if err != nil {
				return fmt.Errorf("transformImagesToSimpleColor: err=%v", err)
			}
			rgbImg, err := ximg.ColorSpace.ImageToRGB(*img)
			if err != nil {
				return fmt.Errorf("transformImagesToSimpleColor: err=%v", err)
			}
			if cs.GetNumComponents() == 1 {
				rgbColorSpace := model.NewPdfColorspaceDeviceRGB()
				rgbImg, err = rgbColorSpace.ImageToGray(rgbImg)
				if err != nil {
					return fmt.Errorf("transformImagesToSimpleColor: err=%v", err)
				}
			}
			// Update the XObject image.
			// Use same encoder as input data.  Make sure for DCT filter it is updated to 1 color component.
			encoder := ximg.Filter
			if dctEncoder, is := encoder.(*core.DCTEncoder); is {
				dctEncoder.ColorComponents = cs.GetNumComponents()
			}
			ximgRGB, err := model.NewXObjectImageFromImage(&rgbImg, nil, encoder)
			if err != nil {
				if err == core.ErrUnsupportedEncodingParameters {
					// Unsupported encoding parameters, revert to a basic flate encoder without predictor.
					encoder = core.NewFlateEncoder()
					panic("will this work")
				}
				// Try again, fail if error.
				ximgRGB, err = model.NewXObjectImageFromImage(&rgbImg, nil, encoder)
				if err != nil {
					return fmt.Errorf("transformImagesToSimpleColor: Error creating image err=%v", err)
				}
			}
			// Update the entry.
			err = resources.SetXObjectImageByName(*nameObj, ximgRGB)
			if err != nil {
				return fmt.Errorf("transformImagesToSimpleColor: Failed setting XObject %+q. err=%v",
					name, err)
			}

			return nil
		})

	return processor.Process(resources)
}

// mb returns the number of megabytes in `size` bytes.
func mb(size int64) float64 {
	return float64(size) / 1024.0 / 1024.0
}

// makeUsage updates flag.Usage to include usage message `msg`.
func makeUsage(msg string) {
	usage := flag.Usage
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, msg)
		usage()
	}
}
