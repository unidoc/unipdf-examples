/*
 * Convert a PDF to grayscale in a vectorized fashion, including images and all content.
 *
 * This advanced example demonstrates some of the more complex capabilities of UniPDF, showing the capability to process
 * and transform objects and contents.
 *
 * Run as: go run pdf_grayscale_transform.go color.pdf output.pdf
 */

package main

import (
	"errors"
	"fmt"

	"os"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/contentstream"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/ps"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Syntax: go run pdf_grayscale_transform.go input.pdf output.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	err := convertPdfToGrayscale(inputPath, outputPath)
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Completed, see output %s\n", outputPath)
}

func convertPdfToGrayscale(inputPath, outputPath string) error {
	pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, nil)
	if err != nil {
		return err
	}
	defer f.Close()

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}
	fmt.Printf("PDF Num Pages: %d\n", numPages)

	// Process each page using the following callback
	// when generating PdfWriter.
	opts := &model.ReaderToWriterOpts{
		PageProcessCallback: func(pageNum int, page *model.PdfPage) error {
			fmt.Printf("Processing page %d/%d\n", pageNum, numPages)

			err = convertPageToGrayscale(page)
			if err != nil {
				return err
			}

			return nil
		},
	}

	// Generate a PdfWriter instance from existing PdfReader.
	pdfWriter, err := pdfReader.ToWriter(opts)
	if err != nil {
		return err
	}

	// Write to file.
	err = pdfWriter.WriteToFile(outputPath)
	if err != nil {
		return err
	}

	return nil
}

// Replaces color objects on the page with grayscale ones.  Also references XObject Images and Forms
// to convert those to grayscale.
func convertPageToGrayscale(page *model.PdfPage) error {
	// For each page, we go through the resources and look for the images.
	contents, err := page.GetAllContentStreams()
	if err != nil {
		return err
	}

	grayContent, err := transformContentStreamToGrayscale(contents, page.Resources)
	if err != nil {
		return err
	}
	page.SetContentStreams([]string{string(grayContent)}, core.NewFlateEncoder())

	//fmt.Printf("Processed contents: %s\n", grayContent)

	return nil
}

// Check if colorspace represents a Pattern colorspace.
func isPatternCS(cs model.PdfColorspace) bool {
	_, isPattern := cs.(*model.PdfColorspaceSpecialPattern)
	return isPattern
}

func transformContentStreamToGrayscale(contents string, resources *model.PdfPageResources) ([]byte, error) {
	cstreamParser := contentstream.NewContentStreamParser(contents)
	operations, err := cstreamParser.Parse()
	if err != nil {
		return nil, err
	}
	processedOperations := &contentstream.ContentStreamOperations{}

	transformedPatterns := map[core.PdfObjectName]bool{} // List of already transformed patterns. Avoid multiple conversions.
	transformedShadings := map[core.PdfObjectName]bool{} // List of already transformed shadings. Avoid multiple conversions.

	// The content stream processor keeps track of the graphics state and we can make our own handlers to process certain commands,
	// using the AddHandler method.  In this case, we hook up to color related operands, and for image and form handling.
	processor := contentstream.NewContentStreamProcessor(*operations)
	// Add handlers for colorspace related functionality.
	processor.AddHandler(contentstream.HandlerConditionEnumAllOperands, "",
		func(op *contentstream.ContentStreamOperation, gs contentstream.GraphicsState, resources *model.PdfPageResources) error {
			operand := op.Operand
			switch operand {
			case "CS": // Set colorspace operands (stroking).
				if isPatternCS(gs.ColorspaceStroking) {
					// If referring to a pattern colorspace with an external definition, need to update the definition.
					// If has an underlying colorspace, then go and change it to DeviceGray.
					// Needs to be specified externally in the colorspace resources.

					csname := op.Params[0].(*core.PdfObjectName)
					if *csname != "Pattern" {
						// Update if referring to an external colorspace in resources.
						cs, ok := resources.GetColorspaceByName(*csname)
						if !ok {
							common.Log.Debug("Undefined colorspace for pattern (%s)", csname)
							return errors.New("Colorspace not defined")
						}

						patternCS, ok := cs.(*model.PdfColorspaceSpecialPattern)
						if !ok {
							return errors.New("Type error")
						}

						if patternCS.UnderlyingCS != nil {
							// Swap out for a gray colorspace.
							patternCS.UnderlyingCS = model.NewPdfColorspaceDeviceGray()
						}

						err = resources.SetColorspaceByName(*csname, patternCS)
						if err != nil {
							return err
						}
					}
					*processedOperations = append(*processedOperations, op)
					return nil
				}

				op := contentstream.ContentStreamOperation{}
				op.Operand = operand
				op.Params = []core.PdfObject{core.MakeName("DeviceGray")}
				*processedOperations = append(*processedOperations, &op)
				return nil
			case "cs": // Set colorspace operands (non-stroking).
				if isPatternCS(gs.ColorspaceNonStroking) {
					// If referring to a pattern colorspace with an external definition, need to update the definition.
					// If has an underlying colorspace, then go and change it to DeviceGray.
					// Needs to be specified externally in the colorspace resources.

					csname := op.Params[0].(*core.PdfObjectName)
					if *csname != "Pattern" {
						// Update if referring to an external colorspace in resources.
						cs, ok := resources.GetColorspaceByName(*csname)
						if !ok {
							common.Log.Debug("Undefined colorspace for pattern (%s)", csname)
							return errors.New("Colorspace not defined")
						}

						patternCS, ok := cs.(*model.PdfColorspaceSpecialPattern)
						if !ok {
							return errors.New("Type error")
						}

						if patternCS.UnderlyingCS != nil {
							// Swap out for a gray colorspace.
							patternCS.UnderlyingCS = model.NewPdfColorspaceDeviceGray()
						}

						resources.SetColorspaceByName(*csname, patternCS)
					}
					*processedOperations = append(*processedOperations, op)
					return nil
				}

				op := contentstream.ContentStreamOperation{}
				op.Operand = operand
				op.Params = []core.PdfObject{core.MakeName("DeviceGray")}
				*processedOperations = append(*processedOperations, &op)
				return nil

			case "SC", "SCN": // Set stroking color.  Includes pattern colors.
				if isPatternCS(gs.ColorspaceStroking) {
					op := contentstream.ContentStreamOperation{}
					op.Operand = operand
					op.Params = []core.PdfObject{}

					patternColor, ok := gs.ColorStroking.(*model.PdfColorPattern)
					if !ok {
						return errors.New("Invalid stroking color type")
					}

					if patternColor.Color != nil {
						color, err := gs.ColorspaceStroking.ColorToRGB(patternColor.Color)
						if err != nil {
							fmt.Printf("Error: %v\n", err)
							return err
						}
						rgbColor := color.(*model.PdfColorDeviceRGB)
						grayColor := rgbColor.ToGray()

						op.Params = append(op.Params, core.MakeFloat(grayColor.Val()))
					}

					if _, has := transformedPatterns[patternColor.PatternName]; has {
						// Already processed, need not change anything, except underlying color if used.
						op.Params = append(op.Params, core.MakeName(string(patternColor.PatternName)))
						*processedOperations = append(*processedOperations, &op)
						return nil
					}
					transformedPatterns[patternColor.PatternName] = true

					// Look up the pattern name and convert it.
					pattern, found := resources.GetPatternByName(patternColor.PatternName)
					if !found {
						return errors.New("Undefined pattern name")
					}

					grayPattern, err := convertPatternToGray(pattern)
					if err != nil {
						common.Log.Debug("Unable to convert pattern to grayscale: %v", err)
						return err
					}
					resources.SetPatternByName(patternColor.PatternName, grayPattern.ToPdfObject())

					op.Params = append(op.Params, core.MakeName(string(patternColor.PatternName)))
					*processedOperations = append(*processedOperations, &op)
				} else {
					color, err := gs.ColorspaceStroking.ColorToRGB(gs.ColorStroking)
					if err != nil {
						fmt.Printf("Error with ColorToRGB: %v\n", err)
						return err
					}
					rgbColor := color.(*model.PdfColorDeviceRGB)
					grayColor := rgbColor.ToGray()

					op := contentstream.ContentStreamOperation{}
					op.Operand = operand
					op.Params = []core.PdfObject{core.MakeFloat(grayColor.Val())}
					*processedOperations = append(*processedOperations, &op)
				}

				return nil
			case "sc", "scn": // Set nonstroking color.
				if isPatternCS(gs.ColorspaceNonStroking) {
					op := contentstream.ContentStreamOperation{}
					op.Operand = operand
					op.Params = []core.PdfObject{}

					patternColor, ok := gs.ColorNonStroking.(*model.PdfColorPattern)
					if !ok {
						return errors.New("Invalid stroking color type")
					}

					if patternColor.Color != nil {
						color, err := gs.ColorspaceNonStroking.ColorToRGB(patternColor.Color)
						if err != nil {
							fmt.Printf("Error : %v\n", err)
							return err
						}
						rgbColor := color.(*model.PdfColorDeviceRGB)
						grayColor := rgbColor.ToGray()

						op.Params = append(op.Params, core.MakeFloat(grayColor.Val()))
					}

					if _, has := transformedPatterns[patternColor.PatternName]; has {
						// Already processed, need not change anything, except underlying color if used.
						op.Params = append(op.Params, core.MakeName(string(patternColor.PatternName)))
						*processedOperations = append(*processedOperations, &op)
						return nil
					}
					transformedPatterns[patternColor.PatternName] = true

					// Look up the pattern name and convert it.
					pattern, found := resources.GetPatternByName(patternColor.PatternName)
					if !found {
						return errors.New("Undefined pattern name")
					}

					grayPattern, err := convertPatternToGray(pattern)
					if err != nil {
						common.Log.Debug("Unable to convert pattern to grayscale: %v", err)
						return err
					}
					resources.SetPatternByName(patternColor.PatternName, grayPattern.ToPdfObject())

					op.Params = append(op.Params, core.MakeName(string(patternColor.PatternName)))
					*processedOperations = append(*processedOperations, &op)
				} else {
					color, err := gs.ColorspaceNonStroking.ColorToRGB(gs.ColorNonStroking)
					if err != nil {
						fmt.Printf("Error: %v\n", err)
						return err
					}
					rgbColor := color.(*model.PdfColorDeviceRGB)
					grayColor := rgbColor.ToGray()

					op := contentstream.ContentStreamOperation{}
					op.Operand = operand
					op.Params = []core.PdfObject{core.MakeFloat(grayColor.Val())}

					*processedOperations = append(*processedOperations, &op)
				}
				return nil
			case "RG", "K": // Set RGB or CMYK stroking color.
				color, err := gs.ColorspaceStroking.ColorToRGB(gs.ColorStroking)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					return err
				}
				rgbColor := color.(*model.PdfColorDeviceRGB)
				grayColor := rgbColor.ToGray()

				op := contentstream.ContentStreamOperation{}
				op.Operand = "G"
				op.Params = []core.PdfObject{core.MakeFloat(grayColor.Val())}

				*processedOperations = append(*processedOperations, &op)
				return nil
			case "rg", "k": // Set RGB or CMYK as nonstroking color.
				color, err := gs.ColorspaceNonStroking.ColorToRGB(gs.ColorNonStroking)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					return err
				}
				rgbColor := color.(*model.PdfColorDeviceRGB)
				grayColor := rgbColor.ToGray()

				op := contentstream.ContentStreamOperation{}
				op.Operand = "g"
				op.Params = []core.PdfObject{core.MakeFloat(grayColor.Val())}

				*processedOperations = append(*processedOperations, &op)
				return nil
			case "sh": // Paints the shape and color defined by shading dict.
				if len(op.Params) != 1 {
					return errors.New("Params to sh operator should be 1")
				}
				shname, ok := op.Params[0].(*core.PdfObjectName)
				if !ok {
					return errors.New("sh parameter should be a name")
				}
				if _, has := transformedShadings[*shname]; has {
					// Already processed, no need to do anything.
					*processedOperations = append(*processedOperations, op)
					return nil
				}
				transformedShadings[*shname] = true

				shading, found := resources.GetShadingByName(*shname)
				if !found {
					return errors.New("Shading not defined in resources")
				}

				grayShading, err := convertShadingToGray(shading)
				if err != nil {
					return err
				}

				resources.SetShadingByName(*shname, grayShading.GetContext().ToPdfObject())
			}
			*processedOperations = append(*processedOperations, op)

			return nil
		})
	// Add handler for image related handling.  Note that inline images are completely stored with a ContentStreamInlineImage
	// object as the parameter for BI.
	processor.AddHandler(contentstream.HandlerConditionEnumOperand, "BI",
		func(op *contentstream.ContentStreamOperation, gs contentstream.GraphicsState, resources *model.PdfPageResources) error {
			if len(op.Params) != 1 {
				fmt.Printf("BI Error invalid number of params\n")
				return errors.New("invalid number of parameters")
			}
			// Inline image.
			iimg, ok := op.Params[0].(*contentstream.ContentStreamInlineImage)
			if !ok {
				fmt.Printf("Error: Invalid handling for inline image\n")
				return errors.New("Invalid inline image parameter")
			}

			img, err := iimg.ToImage(resources)
			if err != nil {
				fmt.Printf("Error converting inline image to image: %v\n", err)
				return err
			}

			cs, err := iimg.GetColorSpace(resources)
			if err != nil {
				fmt.Printf("Error getting color space for inline image: %v\n", err)
				return err
			}
			rgbImg, err := cs.ImageToRGB(*img)
			if err != nil {
				fmt.Printf("Error converting image to rgb: %v\n", err)
				return err
			}
			rgbColorSpace := model.NewPdfColorspaceDeviceRGB()
			grayImage, err := rgbColorSpace.ImageToGray(rgbImg)
			if err != nil {
				fmt.Printf("Error converting img to gray: %v\n", err)
				return err
			}

			// Update the XObject image.
			// Use same encoder as input data.  Make sure for DCT filter it is updated to 1 color component.
			encoder, err := iimg.GetEncoder()
			if err != nil {
				fmt.Printf("Error getting encoder for inline image: %v\n", err)
				return err
			}
			if dctEncoder, is := encoder.(*core.DCTEncoder); is {
				dctEncoder.ColorComponents = 1
			}

			grayInlineImg, err := contentstream.NewInlineImageFromImage(grayImage, encoder)
			if err != nil {
				if err == core.ErrUnsupportedEncodingParameters {
					// Unsupported encoding parameters, revert to a basic flate encoder without predictor.
					encoder = core.NewFlateEncoder()
				}
				// Try again, fail on error.
				grayInlineImg, err = contentstream.NewInlineImageFromImage(grayImage, encoder)
				if err != nil {
					fmt.Printf("Error making a new inline image object: %v\n", err)
					return err
				}
			}

			// Replace inline image data with the gray image.
			pOp := contentstream.ContentStreamOperation{}
			pOp.Operand = "BI"
			pOp.Params = []core.PdfObject{grayInlineImg}
			*processedOperations = append(*processedOperations, &pOp)

			return nil
		})

	// Handler for XObject Image and Forms.
	processedXObjects := map[string]bool{} // Keep track of processed XObjects to avoid repetition.

	processor.AddHandler(contentstream.HandlerConditionEnumOperand, "Do",
		func(op *contentstream.ContentStreamOperation, gs contentstream.GraphicsState, resources *model.PdfPageResources) error {
			if len(op.Params) < 1 {
				fmt.Printf("ERROR: Invalid number of params for Do object.\n")
				return errors.New("Range check")
			}

			// XObject.
			name := op.Params[0].(*core.PdfObjectName)

			// Only process each one once.
			_, has := processedXObjects[string(*name)]
			if has {
				return nil
			}
			processedXObjects[string(*name)] = true

			_, xtype := resources.GetXObjectByName(*name)
			if xtype == model.XObjectTypeImage {
				//fmt.Printf(" XObject Image: %s\n", *name)

				ximg, err := resources.GetXObjectImageByName(*name)
				if err != nil {
					fmt.Printf("Error w/GetXObjectImageByName : %v\n", err)
					return err
				}

				img, err := ximg.ToImage()
				if err != nil {
					fmt.Printf("Error w/ToImage: %v\n", err)
					return err
				}

				rgbImg, err := ximg.ColorSpace.ImageToRGB(*img)
				if err != nil {
					fmt.Printf("Error ImageToRGB: %v\n", err)
					return err
				}

				rgbColorSpace := model.NewPdfColorspaceDeviceRGB()
				grayImage, err := rgbColorSpace.ImageToGray(rgbImg)
				if err != nil {
					fmt.Printf("Error ImageToGray: %v\n", err)
					return err
				}

				// Update the XObject image.
				// Use same encoder as input data.  Make sure for DCT filter it is updated to 1 color component.
				encoder := ximg.Filter
				if dctEncoder, is := encoder.(*core.DCTEncoder); is {
					dctEncoder.ColorComponents = 1
				}

				ximgGray, err := model.NewXObjectImageFromImage(&grayImage, nil, encoder)
				if err != nil {
					if err == core.ErrUnsupportedEncodingParameters {
						// Unsupported encoding parameters, revert to a basic flate encoder without predictor.
						encoder = core.NewFlateEncoder()
					}

					// Try again, fail if error.
					ximgGray, err = model.NewXObjectImageFromImage(&grayImage, nil, encoder)
					if err != nil {
						fmt.Printf("Error creating image: %v\n", err)
						return err
					}
				}

				// Update the entry.
				err = resources.SetXObjectImageByName(*name, ximgGray)
				if err != nil {
					fmt.Printf("Failed setting x object: %v (%s)\n", err, string(*name))
					return err
				}
			} else if xtype == model.XObjectTypeForm {
				//fmt.Printf(" XObject Form: %s\n", *name)

				// Go through the XObject Form content stream.
				xform, err := resources.GetXObjectFormByName(*name)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					return err
				}

				formContent, err := xform.GetContentStream()
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					return err
				}

				// Process the content stream in the Form object too:
				// XXX/TODO/Consider: Use either form resources (priority) and fall back to page resources alternatively if not found.
				// Have not come into cases where needed yet.
				formResources := xform.Resources
				if formResources == nil {
					formResources = resources
				}

				// Process the content stream in the Form object too:
				grayContent, err := transformContentStreamToGrayscale(string(formContent), formResources)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					return err
				}

				xform.SetContentStream(grayContent, nil)

				// Update the resource entry.
				resources.SetXObjectFormByName(*name, xform)
			}

			return nil
		})

	err = processor.Process(resources)
	if err != nil {
		fmt.Printf("Error processing: %v\n", err)
		return nil, err
	}

	// For debug purposes: (high level logging).
	//
	//fmt.Printf("=== Unprocessed - Full list\n")
	//for idx, op := range operations {
	//	fmt.Printf("U. Operation %d: %s - Params: %v\n", idx+1, op.Operand, op.Params)
	//}
	//fmt.Printf("=== Processed - Full list\n")
	//for idx, op := range *processedOperations {
	//	fmt.Printf("P. Operation %d: %s - Params: %v\n", idx+1, op.Operand, op.Params)
	//}

	return processedOperations.Bytes(), nil
}

// Convert a pattern to grayscale (tiling or shading pattern).
func convertPatternToGray(pattern *model.PdfPattern) (*model.PdfPattern, error) {
	// Case 1: Colored tiling patterns.  Need to process the content stream and replace.
	if pattern.IsTiling() {
		tilingPattern := pattern.GetAsTilingPattern()

		if tilingPattern.IsColored() {
			// A colored tiling pattern can use color operators in its stream, need to process the stream.

			content, err := tilingPattern.GetContentStream()
			if err != nil {
				return nil, err
			}

			grayContents, err := transformContentStreamToGrayscale(string(content), tilingPattern.Resources)
			if err != nil {
				return nil, err
			}

			tilingPattern.SetContentStream(grayContents, nil)

			// Update in-memory pdf objects.
			_ = tilingPattern.ToPdfObject()
		}
	} else if pattern.IsShading() {
		// Case 2: Shading patterns.  Need to create a new colorspace that can map from N=3,4 colorspaces to grayscale.
		shadingPattern := pattern.GetAsShadingPattern()

		grayShading, err := convertShadingToGray(shadingPattern.Shading)
		if err != nil {
			return nil, err
		}
		shadingPattern.Shading = grayShading

		// Update in-memory pdf objects.
		_ = shadingPattern.ToPdfObject()
	}

	return pattern, nil
}

// Convert shading to grayscale.
// This one is slightly involved as a shading defines a color as function of position, i.e. color(x,y) = F(x,y).
// Since the function can be challenging to change, we define new DeviceN colorspace with a color conversion
// function.
func convertShadingToGray(shading *model.PdfShading) (*model.PdfShading, error) {
	cs := shading.ColorSpace

	if cs.GetNumComponents() == 1 {
		// Already grayscale, should be fine. No action taken.
		return shading, nil
	} else if cs.GetNumComponents() == 3 {
		// Create a new DeviceN colorspace that converts R,G,B -> Grayscale
		// Use: gray := 0.3*R + 0.59G + 0.11B
		// PS program: { 0.11 mul exch 0.59 mul add exch 0.3 mul add }.
		transformFunc := &model.PdfFunctionType4{}
		transformFunc.Domain = []float64{0, 1, 0, 1, 0, 1}
		transformFunc.Range = []float64{0, 1}
		rgbToGrayPsProgram := ps.NewPSProgram()
		rgbToGrayPsProgram.Append(ps.MakeReal(0.11))
		rgbToGrayPsProgram.Append(ps.MakeOperand("mul"))
		rgbToGrayPsProgram.Append(ps.MakeOperand("exch"))
		rgbToGrayPsProgram.Append(ps.MakeReal(0.59))
		rgbToGrayPsProgram.Append(ps.MakeOperand("mul"))
		rgbToGrayPsProgram.Append(ps.MakeOperand("add"))
		rgbToGrayPsProgram.Append(ps.MakeOperand("exch"))
		rgbToGrayPsProgram.Append(ps.MakeReal(0.3))
		rgbToGrayPsProgram.Append(ps.MakeOperand("mul"))
		rgbToGrayPsProgram.Append(ps.MakeOperand("add"))
		transformFunc.Program = rgbToGrayPsProgram

		// Define the DeviceN colorspace that performs the R,G,B -> Gray conversion for us.
		transformcs := model.NewPdfColorspaceDeviceN()
		transformcs.AlternateSpace = model.NewPdfColorspaceDeviceGray()
		transformcs.ColorantNames = core.MakeArray(core.MakeName("R"), core.MakeName("G"), core.MakeName("B"))
		transformcs.TintTransform = transformFunc

		// Replace the old colorspace with the new.
		shading.ColorSpace = transformcs

		return shading, nil
	} else if cs.GetNumComponents() == 4 {
		// Create a new DeviceN colorspace that converts C,M,Y,K -> Grayscale.
		// Use: gray = 1.0 - min(1.0, 0.3*C + 0.59*M + 0.11*Y + K)  ; where BG(k) = k simply.
		// PS program: {exch 0.11 mul add exch 0.59 mul add exch 0.3 mul add dup 1.0 ge { pop 1.0 } if}
		transformFunc := &model.PdfFunctionType4{}
		transformFunc.Domain = []float64{0, 1, 0, 1, 0, 1, 0, 1}
		transformFunc.Range = []float64{0, 1}

		cmykToGrayPsProgram := ps.NewPSProgram()
		cmykToGrayPsProgram.Append(ps.MakeOperand("exch"))
		cmykToGrayPsProgram.Append(ps.MakeReal(0.11))
		cmykToGrayPsProgram.Append(ps.MakeOperand("mul"))
		cmykToGrayPsProgram.Append(ps.MakeOperand("add"))
		cmykToGrayPsProgram.Append(ps.MakeOperand("exch"))
		cmykToGrayPsProgram.Append(ps.MakeReal(0.59))
		cmykToGrayPsProgram.Append(ps.MakeOperand("mul"))
		cmykToGrayPsProgram.Append(ps.MakeOperand("add"))
		cmykToGrayPsProgram.Append(ps.MakeOperand("exch"))
		cmykToGrayPsProgram.Append(ps.MakeReal(0.30))
		cmykToGrayPsProgram.Append(ps.MakeOperand("mul"))
		cmykToGrayPsProgram.Append(ps.MakeOperand("add"))
		cmykToGrayPsProgram.Append(ps.MakeOperand("dup"))
		cmykToGrayPsProgram.Append(ps.MakeReal(1.0))
		cmykToGrayPsProgram.Append(ps.MakeOperand("ge"))
		// Add sub procedure.
		subProc := ps.NewPSProgram()
		subProc.Append(ps.MakeOperand("pop"))
		subProc.Append(ps.MakeReal(1.0))
		cmykToGrayPsProgram.Append(subProc)
		cmykToGrayPsProgram.Append(ps.MakeOperand("if"))
		transformFunc.Program = cmykToGrayPsProgram

		// Define the DeviceN colorspace that performs the R,G,B -> Gray conversion for us.
		transformcs := model.NewPdfColorspaceDeviceN()
		transformcs.AlternateSpace = model.NewPdfColorspaceDeviceGray()
		transformcs.ColorantNames = core.MakeArray(core.MakeName("C"), core.MakeName("M"), core.MakeName("Y"), core.MakeName("K"))
		transformcs.TintTransform = transformFunc

		// Replace the old colorspace with the new.
		shading.ColorSpace = transformcs

		return shading, nil
	} else {
		common.Log.Debug("Cannot convert to shading pattern grayscale, color space N = %d", cs.GetNumComponents())
		return nil, errors.New("Unsupported pattern colorspace for grayscale conversion")
	}
}
