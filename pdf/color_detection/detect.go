package main

import (
	"errors"
	"fmt"

	common "github.com/unidoc/unidoc/common"
	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

// isPageColored returns true if `page` contains color. It also references
// XObject Images and Forms to _possibly_ record if they contain color
func isPageColored(page *pdf.PdfPage, desc string, debug bool) (bool, error) {
	// For each page, we go through the resources and look for the images.
	resources, err := page.GetResources()
	if err != nil {
		common.Log.Error("GetResources failed. err=%v", err)
		return false, err
	}

	contents, err := page.GetAllContentStreams()
	if err != nil {
		common.Log.Error("GetAllContentStreams failed. err=%v", err)
		return false, err
	}

	if debug {
		fmt.Println("\n===============***================")
		fmt.Printf("%s\n", desc)
		fmt.Println("===============+++================")
		fmt.Printf("%s\n", contents)
		fmt.Println("==================================")
	}

	colored, err := isContentStreamColored(contents, resources, debug)
	if debug {
		common.Log.Info("colored=%t err=%v", colored, err)
	}
	if err != nil {
		common.Log.Error("isContentStreamColored failed. err=%v", err)
		return false, err
	}
	return colored, nil
}

// isPatternCS returns true if `colorspace` represents a Pattern colorspace.
func isPatternCS(cs pdf.PdfColorspace) bool {
	_, isPattern := cs.(*pdf.PdfColorspaceSpecialPattern)
	return isPattern
}

// isContentStreamColored returns true if `contents` contains any color object
func isContentStreamColored(contents string, resources *pdf.PdfPageResources, debug bool) (bool, error) {
	cstreamParser := pdfcontent.NewContentStreamParser(contents)
	operations, err := cstreamParser.Parse()
	if err != nil {
		return false, err
	}

	colored := false                                    // Has a colored mark been detected in the stream?
	coloredPatterns := map[pdfcore.PdfObjectName]bool{} // List of already detected patterns. Re-use for subsequent detections.
	coloredShadings := map[string]bool{}                // List of already detected shadings. Re-use for subsequent detections.

	// The content stream processor keeps track of the graphics state and we can make our own handlers to process
	// certain commands using the AddHandler method. In this case, we hook up to color related operands, and for image
	// and form handling.
	processor := pdfcontent.NewContentStreamProcessor(*operations)
	// Add handlers for colorspace related functionality.
	processor.AddHandler(pdfcontent.HandlerConditionEnumAllOperands, "",
		func(op *pdfcontent.ContentStreamOperation, gs pdfcontent.GraphicsState,
			resources *pdf.PdfPageResources) error {
			operand := op.Operand
			switch operand {
			case "SC", "SCN": // Set stroking color.  Includes pattern colors.
				// if gs.ColorspaceStroking != nil {
				// 	if gs.ColorspaceStroking.GetNumComponents() != len((op.Params)) {
				// 		common.Log.Error("Mismatch cs=%#v op=%#v", gs.ColorspaceStroking, op)
				// 		panic("wtf")
				// 	}
				// }
				if isPatternCS(gs.ColorspaceStroking) {
					op := pdfcontent.ContentStreamOperation{}
					op.Operand = operand
					op.Params = []pdfcore.PdfObject{}

					patternColor, ok := gs.ColorStroking.(*pdf.PdfColorPattern)
					if !ok {
						return errors.New("Invalid stroking color type")
					}

					if patternColor.Color != nil {
						color, err := gs.ColorspaceStroking.ColorToRGB(patternColor.Color)
						if err != nil {
							common.Log.Error("err=%v", err)
							return err
						}
						rgbColor := color.(*pdf.PdfColorDeviceRGB)
						if rgbColor.IsColored() {
							if debug {
								common.Log.Info("op=%s col=%t", op, true)
							}
							colored = true
							return nil
						}
					}

					if col, ok := coloredPatterns[patternColor.PatternName]; ok {
						// Already processed, need not change anything, except underlying color if used.
						if col {
							if debug {
								common.Log.Info("op=%s col=%t", op, col)
							}
							colored = true
						}
						return nil
					}

					// Look up the pattern name and convert it.
					pattern, found := resources.GetPatternByName(patternColor.PatternName)
					if !found {
						return errors.New("Undefined pattern name")
					}
					col, err := isPatternColored(pattern, debug)
					if err != nil {
						common.Log.Error("isPatternColored failed. err=%v", err)
						return err
					}
					coloredPatterns[patternColor.PatternName] = col
					colored = colored || col
					if debug {
						common.Log.Info("op=%s col=%t", op, col)
					}

				} else {
					color, err := gs.ColorspaceStroking.ColorToRGB(gs.ColorStroking)
					if err != nil {
						common.Log.Error("Error with ColorToRGB: %v", err)
						return err
					}
					rgbColor := color.(*pdf.PdfColorDeviceRGB)
					col := rgbColor.IsColored()
					colored = colored || col
					if debug {
						common.Log.Info("op=%s ColorspaceStroking=%T ColorStroking=%#v col=%t",
							op, gs.ColorspaceStroking, gs.ColorStroking, col)
						if col {
							panic("Done")
						}
					}
				}
				return nil
			case "sc", "scn": // Set non-stroking color.
				// common.Log.Info("!!!*op=%s ColorspaceNonStroking=%s ColorNonStroking=%+v",
				// 	op, gs.ColorspaceNonStroking, gs.ColorNonStroking)
				// if gs.ColorspaceNonStroking.GetNumComponents() != len((op.Params)) {
				// 	common.Log.Error("Mismatch cs=%s op=%s", gs.ColorspaceNonStroking, op)
				// 	panic("wtf")
				// }
				if isPatternCS(gs.ColorspaceNonStroking) {
					op := pdfcontent.ContentStreamOperation{}
					op.Operand = operand
					op.Params = []pdfcore.PdfObject{}
					patternColor, ok := gs.ColorNonStroking.(*pdf.PdfColorPattern)
					if !ok {
						return errors.New("Invalid stroking color type")
					}
					if patternColor.Color != nil {
						color, err := gs.ColorspaceNonStroking.ColorToRGB(patternColor.Color)
						if err != nil {
							common.Log.Error("err=%v", err)
							return err
						}
						rgbColor := color.(*pdf.PdfColorDeviceRGB)
						col := rgbColor.IsColored()
						colored = colored || col
						if debug {
							common.Log.Info("op=%#v col=%t", op, col)
						}
					}
					if col, ok := coloredPatterns[patternColor.PatternName]; ok {
						// Already processed, need not change anything, except underlying color if used.
						colored = colored || col
						if debug {
							common.Log.Info("op=%#v col=%t", op, col)
						}
						return nil
					}

					// Look up the pattern name and convert it.
					pattern, found := resources.GetPatternByName(patternColor.PatternName)
					if !found {
						return errors.New("Undefined pattern name")
					}
					col, err := isPatternColored(pattern, debug)
					if err != nil {
						common.Log.Debug("Unable to convert pattern to grayscale: %v", err)
						return err
					}
					coloredPatterns[patternColor.PatternName] = col
				} else {
					// common.Log.Info("!!!!op=%s ColorspaceNonStroking=%s ColorNonStroking=%+v",
					// 	op, gs.ColorspaceNonStroking, gs.ColorNonStroking)
					color, err := gs.ColorspaceNonStroking.ColorToRGB(gs.ColorNonStroking)
					if err != nil {
						common.Log.Error("err=%v", err)
						return err
					}
					rgbColor := color.(*pdf.PdfColorDeviceRGB)
					col := rgbColor.IsColored()
					if debug {
						// common.Log.Info("op=%s col=%t ColorNonStroking=%#v rgbColor=%+v",
						// 	op, col, gs.ColorNonStroking, rgbColor)
						// common.Log.Info("ColorspaceNonStroking=\n%#v\n%s",
						// 	gs.ColorspaceNonStroking, gs.ColorspaceNonStroking)
					}
					colored = colored || col
					if debug {
						common.Log.Info("op=%s ColorspaceNonStroking=%T ColorNonStroking=%#v col=%t",
							op, gs.ColorspaceNonStroking, gs.ColorNonStroking, col)
						if col {
							panic("Done")
						}
					}

				}
				return nil
			case "RG", "K": // Set RGB or CMYK stroking color.
				color, err := gs.ColorspaceStroking.ColorToRGB(gs.ColorStroking)
				if err != nil {
					common.Log.Error("err=%v", err)
					return err
				}
				rgbColor := color.(*pdf.PdfColorDeviceRGB)
				col := rgbColor.IsColored()
				if debug {
					common.Log.Info("op=%s ColorspaceNonStroking=%T ColorNonStroking=%#v col=%t",
						op, gs.ColorspaceNonStroking, gs.ColorNonStroking, col)
					if col {
						panic("Done")
					}
				}
				colored = colored || col
				return nil
			case "rg", "k": // Set RGB or CMYK as non-stroking color.
				color, err := gs.ColorspaceNonStroking.ColorToRGB(gs.ColorNonStroking)
				if err != nil {
					common.Log.Error("err=%v", err)
					return err
				}
				rgbColor := color.(*pdf.PdfColorDeviceRGB)
				// fmt.Printf("rgbColor=%v\n", rgbColor)
				col := rgbColor.IsColored()
				colored = colored || col
				if debug {
					common.Log.Info("op=%s ColorspaceStroking=%T ColorStroking=%#v col=%t",
						op, gs.ColorspaceStroking, gs.ColorStroking, col)
					if col {
						panic("Done")
					}
				}
				return nil
			case "sh": // Paints the shape and color defined by shading dict.
				if len(op.Params) != 1 {
					return errors.New("Params to sh operator should be 1")
				}
				shname, ok := op.Params[0].(*pdfcore.PdfObjectName)
				if !ok {
					return errors.New("sh parameter should be a name")
				}
				if col, has := coloredShadings[string(*shname)]; has {
					// Already processed, no need to do anything.
					colored = colored || col
					if debug {
						common.Log.Info("col=%t", col)
					}
					return nil
				}

				shading, found := resources.GetShadingByName(*shname)
				if !found {
					common.Log.Error("Shading not defined in resources. shname=%#q", string(*shname))
					return errors.New("Shading not defined in resources")
				}
				col, err := isShadingColored(shading)
				if err != nil {
					return err
				}
				coloredShadings[string(*shname)] = col
			}
			return nil
		})

	// Add handler for image related handling.  Note that inline images are completely stored with a ContentStreamInlineImage
	// object as the parameter for BI.
	processor.AddHandler(pdfcontent.HandlerConditionEnumOperand, "BI",
		func(op *pdfcontent.ContentStreamOperation, gs pdfcontent.GraphicsState, resources *pdf.PdfPageResources) error {
			if len(op.Params) != 1 {
				err := errors.New("invalid number of parameters")
				common.Log.Error("BI error. err=%v")
				return err
			}
			// Inline image.
			iimg, ok := op.Params[0].(*pdfcontent.ContentStreamInlineImage)
			if !ok {
				common.Log.Error("Invalid handling for inline image")
				return errors.New("Invalid inline image parameter")
			}
			if debug {
				common.Log.Info("iimg=%s", iimg)
			}
			img, err := iimg.ToImage(resources)
			if err != nil {
				common.Log.Error("Error converting inline image to image: %v", err)
				return err
			}

			if debug {
				common.Log.Info("img=%v %d", img.ColorComponents, img.BitsPerComponent)
			}

			if img.ColorComponents <= 1 {
				return nil
			}

			cs, err := iimg.GetColorSpace(resources)
			if err != nil {
				common.Log.Error("Error getting color space for inline image: %v", err)
				return err
			}
			rgbImg, err := cs.ImageToRGB(*img)
			if err != nil {
				common.Log.Error("Error converting image to rgb: %v", err)
				return err
			}
			rgbColorSpace := pdf.NewPdfColorspaceDeviceRGB()
			col := rgbColorSpace.IsImageColored(rgbImg)
			colored = colored || col
			if debug {
				common.Log.Info("col=%t", col)
			}

			return nil
		})

	// !@#$% Black background is here
	// Handler for XObject Image and Forms.
	processedXObjects := map[string]bool{} // Keep track of processed XObjects to avoid repetition.

	processor.AddHandler(pdfcontent.HandlerConditionEnumOperand, "Do",
		func(op *pdfcontent.ContentStreamOperation, gs pdfcontent.GraphicsState, resources *pdf.PdfPageResources) error {

			if len(op.Params) < 1 {
				common.Log.Error("Invalid number of params for Do object")
				return errors.New("Range check")
			}

			// XObject.
			name := op.Params[0].(*pdfcore.PdfObjectName)
			common.Log.Debug("Name=%#v=%#q", name, string(*name))

			// Only process each one once.
			_, has := processedXObjects[string(*name)]
			common.Log.Debug("name=%q has=%t processedXObjects=%+v", *name, has, processedXObjects)
			if has {
				return nil
			}
			processedXObjects[string(*name)] = true

			_, xtype := resources.GetXObjectByName(*name)
			common.Log.Debug("xtype=%+v pdf.XObjectTypeImage=%v", xtype, pdf.XObjectTypeImage)
			if xtype == pdf.XObjectTypeImage {

				ximg, err := resources.GetXObjectImageByName(*name)
				if err != nil {
					common.Log.Error("Error w/GetXObjectImageByName : %v", err)
					return err
				}
				if debug {
					common.Log.Info("!!Filter=%s ColorSpace=%s ImageMask=%v wxd=%dx%d",
						ximg.Filter.GetFilterName(), ximg.ColorSpace,
						ximg.ImageMask, *ximg.Width, *ximg.Height)
				}
				switch ximg.Filter.GetFilterName() {
				case "JPXDecode":
					colored = true
					return nil
				case "CCITTDecode", "JBIG2Decode", "RunLengthDecode":
					return nil
				}

				img, err := ximg.ToImage()
				if err != nil {
					common.Log.Error("Error w/ToImage: %v", err)
					return err
				}

				rgbImg, err := ximg.ColorSpace.ImageToRGB(*img)
				if err != nil {
					common.Log.Error("Error ImageToRGB: %v", err)
					return err
				}

				if debug {
					common.Log.Info("img: ColorComponents=%d wxh=%dx%d", img.ColorComponents, img.Width, img.Height)
					common.Log.Info("ximg: ColorSpace=%T=%s mask=%v", ximg.ColorSpace, ximg.ColorSpace, ximg.Mask)
					common.Log.Info("rgbImg: ColorComponents=%d wxh=%dx%d", rgbImg.ColorComponents, rgbImg.Width, rgbImg.Height)
				}

				rgbColorSpace := pdf.NewPdfColorspaceDeviceRGB()
				col := rgbColorSpace.IsImageColored(rgbImg)
				colored = colored || col
				// !@#$ Update XObj colored map
				if debug {
					common.Log.Info("col=%t", col)
				}

			} else if xtype == pdf.XObjectTypeForm {
				common.Log.Debug(" XObject Form: %s")

				// Go through the XObject Form content stream.
				xform, err := resources.GetXObjectFormByName(*name)
				if err != nil {
					fmt.Printf("Error : %v\n", err)
					return err
				}

				formContent, err := xform.GetContentStream()
				if err != nil {
					common.Log.Error("err=%v")
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
				col, err := isContentStreamColored(string(formContent), formResources, debug)
				if err != nil {
					common.Log.Error("err=%v", err)
					return err
				}
				colored = colored || col
				// !@#$ Update colored XObj map
				if debug {
					common.Log.Info("col=%t", col)
				}

			}

			return nil
		})

	err = processor.Process(resources)
	if err != nil {
		common.Log.Error("processor.Process returned: err=%v", err)
		return false, err
	}

	// if common.LogLevel >= common.LogLevelDebug {
	// 	// For debug purposes: (high level logging).
	// 	fmt.Printf("=== Unprocessed - Full list\n")
	// 	for idx, op := range *operations {
	// 		fmt.Printf("U. Operation %d: %s - Params: %v\n", idx+1, op.Operand, op.Params)
	// 	}
	// 	fmt.Printf("=== Processed - Full list\n")
	// 	for idx, op := range *processedOperations {
	// 		fmt.Printf("P. Operation %d: %s - Params: %v\n", idx+1, op.Operand, op.Params)
	// 	}
	// }
	return colored, nil
}

// isPatternColored returns true if `pattern` contains color (tiling or shading pattern).
func isPatternColored(pattern *pdf.PdfPattern, debug bool) (bool, error) {
	// Case 1: Colored tiling patterns.  Need to process the content stream and replace.
	if pattern.IsTiling() {
		tilingPattern := pattern.GetAsTilingPattern()
		if tilingPattern.IsColored() {
			// A colored tiling pattern can use color operators in its stream, need to process the stream.
			content, err := tilingPattern.GetContentStream()
			if err != nil {
				return false, err
			}
			colored, err := isContentStreamColored(string(content), tilingPattern.Resources, debug)
			return colored, err
		}
	} else if pattern.IsShading() {
		// Case 2: Shading patterns.  Need to create a new colorspace that can map from N=3,4 colorspaces to grayscale.
		shadingPattern := pattern.GetAsShadingPattern()
		colored, err := isShadingColored(shadingPattern.Shading)
		return colored, err
	}
	common.Log.Error("isPatternColored. pattern is neither tiling nor shading")
	panic("wtf")
	return false, nil
}

// isShadingColored returns true is  `shading` is a colored colorspace
func isShadingColored(shading *pdf.PdfShading) (bool, error) {
	cs := shading.ColorSpace
	if cs.GetNumComponents() == 1 {
		// Grayscale colorspace
		return false, nil
	} else if cs.GetNumComponents() == 3 {
		// RGB colorspace
		return true, nil
	} else if cs.GetNumComponents() == 4 {
		// CMYK colorspace
		return true, nil
	} else {
		err := errors.New("Unsupported pattern colorspace for color detection")
		common.Log.Error("isShadingColored: colorpace N=%d err=%v", cs.GetNumComponents(), err)
		return false, err
	}
}
