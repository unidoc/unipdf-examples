package main

import (
	"errors"
	"fmt"

	// unilicense "github.com/unidoc/unidoc/license"
	common "github.com/unidoc/unidoc/common"
	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

// =================================================================================================
// Page transform code goes here
// =================================================================================================

// Replaces color objects on the page with grayscale ones.  Also references XObject Images and Forms
// to convert those to grayscale.
func convertPageToGrayscale(page *pdf.PdfPage, desc string) error {
	// For each page, we go through the resources and look for the images.
	resources, err := page.GetResources()
	if err != nil {
		panic(err)
		return err
	}

	contents, err := page.GetAllContentStreams()
	if err != nil {
		panic(err)
		return err
	}

	grayContent, err := transformContentStreamToGrayscale(contents, resources)
	if err != nil {
		// panic(err)
		return err
	}
	page.SetContentStreams([]string{string(grayContent)}, pdfcore.NewFlateEncoder())

	if gVerbose {
		fmt.Printf("Processed contents: %s\n", grayContent)
	}

	return nil
}

var peter3csData = []byte{
	0x00, 0x00, 0x02, // 0
	0x0C, 0x61, 0x70, // 1
	0x70, 0x6C, 0x02, // 2
	0x00, 0x00, 0x00, // 3
	0x6D, 0x6E, 0x74, // 4
	0x72, 0x52, 0x47, // 5
	0x42, 0x20, 0x58, // 6
	0x59, 0x5A, 0x20, // 7
	0x07, 0xCB, 0x00, // 8
	0x02, 0x00, 0x16, // 9
	0x00, 0x0E, 0x00, // 10
	0x22, 0x00, 0x2C, // 11
	0x61, 0x63, 0x73, // 12
	0x70, 0x41, 0x50, // 13
	0x50, 0x4C, 0x00, // 14
	0x00, 0x00, 0x00, // 15
	0x61, 0x70, 0x70, // 16
	0x6C, 0x00, 0x00, // 17
	0x04, 0x01, 0x00, // 18
	0x00, 0x00, 0x00, // 19
	0x00, 0x00, 0x00, // 20
	0x02, 0x00, 0x00, // 21
	0x00, 0x00, 0x00, // 22
	0x00, 0xF6, 0xD4, // 23
	0x00, 0x01, 0x00, // 24
	0x00, 0x00, 0x00, // 25
	0xD3, 0x2B, 0x00, // 26
	0x00, 0x00, 0x00, // 27
	0x00, 0x00, 0x00, // 28
	0x00, 0x00, 0x00, // 29
	0x00, 0x00, 0x00, // 30
	0x00, 0x00, 0x00, // 31
	0x00, 0x00, 0x00, // 32
	0x00, 0x00, 0x00, // 33
	0x00, 0x00, 0x00, // 34
	0x00, 0x00, 0x00, // 35
	0x00, 0x00, 0x00, // 36
	0x00, 0x00, 0x00, // 37
	0x00, 0x00, 0x00, // 38
	0x00, 0x00, 0x00, // 39
	0x00, 0x00, 0x00, // 40
	0x00, 0x00, 0x00, // 41
	0x00, 0x00, 0x00, // 42
	0x00, 0x00, 0x09, // 43
	0x64, 0x65, 0x73, // 44
	0x63, 0x00, 0x00, // 45
	0x00, 0xF0, 0x00, // 46
	0x00, 0x00, 0x71, // 47
	0x72, 0x58, 0x59, // 48
	0x5A, 0x00, 0x00, // 49
	0x01, 0x64, 0x00, // 50
	0x00, 0x00, 0x14, // 51
	0x67, 0x58, 0x59, // 52
	0x5A, 0x00, 0x00, // 53
	0x01, 0x78, 0x00, // 54
	0x00, 0x00, 0x14, // 55
	0x62, 0x58, 0x59, // 56
	0x5A, 0x00, 0x00, // 57
	0x01, 0x8C, 0x00, // 58
	0x00, 0x00, 0x14, // 59
	0x72, 0x54, 0x52, // 60
	0x43, 0x00, 0x00, // 61
	0x01, 0xA0, 0x00, // 62
	0x00, 0x00, 0x0E, // 63
	0x67, 0x54, 0x52, // 64
	0x43, 0x00, 0x00, // 65
	0x01, 0xB0, 0x00, // 66
	0x00, 0x00, 0x0E, // 67
	0x62, 0x54, 0x52, // 68
	0x43, 0x00, 0x00, // 69
	0x01, 0xC0, 0x00, // 70
	0x00, 0x00, 0x0E, // 71
	0x77, 0x74, 0x70, // 72
	0x74, 0x00, 0x00, // 73
	0x01, 0xD0, 0x00, // 74
	0x00, 0x00, 0x14, // 75
	0x63, 0x70, 0x72, // 76
	0x74, 0x00, 0x00, // 77
	0x01, 0xE4, 0x00, // 78
	0x00, 0x00, 0x27, // 79
	0x64, 0x65, 0x73, // 80
	0x63, 0x00, 0x00, // 81
	0x00, 0x00, 0x00, // 82
	0x00, 0x00, 0x17, // 83
	0x41, 0x70, 0x70, // 84
	0x6C, 0x65, 0x20, // 85
	0x31, 0x33, 0x22, // 86
	0x20, 0x52, 0x47, // 87
	0x42, 0x20, 0x53, // 88
	0x74, 0x61, 0x6E, // 89
	0x64, 0x61, 0x72, // 90
	0x64, 0x00, 0x00, // 91
	0x00, 0x00, 0x00, // 92
	0x00, 0x00, 0x00, // 93
	0x00, 0x00, 0x00, // 94
	0x17, 0x41, 0x70, // 95
	0x70, 0x6C, 0x65, // 96
	0x20, 0x31, 0x33, // 97
	0x22, 0x20, 0x52, // 98
	0x47, 0x42, 0x20, // 99
	0x53, 0x74, 0x61, // 100
	0x6E, 0x64, 0x61, // 101
	0x72, 0x64, 0x00, // 102
	0x00, 0x00, 0x00, // 103
	0x00, 0x00, 0x00, // 104
	0x00, 0x00, 0x00, // 105
	0x00, 0x00, 0x00, // 106
	0x00, 0x00, 0x00, // 107
	0x00, 0x00, 0x00, // 108
	0x00, 0x00, 0x00, // 109
	0x00, 0x00, 0x00, // 110
	0x00, 0x00, 0x00, // 111
	0x00, 0x00, 0x00, // 112
	0x00, 0x00, 0x00, // 113
	0x00, 0x00, 0x00, // 114
	0x00, 0x00, 0x00, // 115
	0x00, 0x00, 0x00, // 116
	0x00, 0x00, 0x58, // 117
	0x59, 0x5A, 0x58, // 118
	0x59, 0x5A, 0x20, // 119
	0x00, 0x00, 0x00, // 120
	0x00, 0x00, 0x00, // 121
	0x63, 0x0A, 0x00, // 122
	0x00, 0x35, 0x0F, // 123
	0x00, 0x00, 0x03, // 124
	0x30, 0x58, 0x59, // 125
	0x5A, 0x20, 0x00, // 126
	0x00, 0x00, 0x00, // 127
	0x00, 0x00, 0x53, // 128
	0x3D, 0x00, 0x00, // 129
	0xAE, 0x37, 0x00, // 130
	0x00, 0x15, 0x76, // 131
	0x58, 0x59, 0x5A, // 132
	0x20, 0x00, 0x00, // 133
	0x00, 0x00, 0x00, // 134
	0x00, 0x40, 0x89, // 135
	0x00, 0x00, 0x1C, // 136
	0xAF, 0x00, 0x00, // 137
	0xBA, 0x82, 0x63, // 138
	0x75, 0x72, 0x76, // 139
	0x00, 0x00, 0x00, // 140
	0x00, 0x00, 0x00, // 141
	0x00, 0x01, 0x01, // 142
	0xCC, 0x63, 0x75, // 143
	0x63, 0x75, 0x72, // 144
	0x76, 0x00, 0x00, // 145
	0x00, 0x00, 0x00, // 146
	0x00, 0x00, 0x01, // 147
	0x01, 0xCC, 0x63, // 148
	0x75, 0x63, 0x75, // 149
	0x72, 0x76, 0x00, // 150
	0x00, 0x00, 0x00, // 151
	0x00, 0x00, 0x00, // 152
	0x01, 0x01, 0xCC, // 153
	0x58, 0x59, 0x58, // 154
	0x59, 0x5A, 0x20, // 155
	0x00, 0x00, 0x00, // 156
	0x00, 0x00, 0x00, // 157
	0xF3, 0x1B, 0x00, // 158
	0x01, 0x00, 0x00, // 159
	0x00, 0x01, 0x67, // 160
	0xE7, 0x74, 0x65, // 161
	0x78, 0x74, 0x00, // 162
	0x00, 0x00, 0x00, // 163
	0x20, 0x43, 0x6F, // 164
	0x70, 0x79, 0x72, // 165
	0x69, 0x67, 0x68, // 166
	0x74, 0x20, 0x41, // 167
	0x70, 0x70, 0x6C, // 168
	0x65, 0x20, 0x43, // 169
	0x6F, 0x6D, 0x70, // 170
	0x75, 0x74, 0x65, // 171
	0x72, 0x73, 0x20, // 172
	0x31, 0x39, 0x39, // 173
	0x34, 0x00, 0x00, // 174
}

func addIccColorspaces(resources *pdf.PdfPageResources) error {
	peter3cs, err := pdf.NewPdfColorspaceICCBased(3)
	if err != nil {
		panic(err)
	}
	peter3cs.Range = []float64{0.0, 1.0, 0.0, 1.0, 0.0, 1.0}
	peter3cs.Data = peter3csData
	if resources.ColorSpace == nil {
		resources.ColorSpace = pdf.NewPdfPageResourcesColorspaces()
	}
	resources.ColorSpace.Add("peter3cs", peter3cs)
	return nil
}

func transformContentStreamToGrayscale(contents string, resources *pdf.PdfPageResources) ([]byte, error) {
	cstreamParser := pdfcontent.NewContentStreamParser(contents)
	operations, err := cstreamParser.Parse()
	if err != nil {
		return nil, err
	}
	processedOperations := &pdfcontent.ContentStreamOperations{}
	addIccColorspaces(resources)

	// The content stream processor keeps track of the graphics state and we can make our own
	// handlers to process certain commands, using the AddHandler method.  In this case, we hook up
	// to color related operands, and for image and form handling.
	processor := pdfcontent.NewContentStreamProcessor(operations)
	// Add handlers for colorspace related functionality.
	processor.AddHandler(pdfcontent.HandlerConditionEnumAllOperands, "",
		func(op *pdfcontent.ContentStreamOperation, gs pdfcontent.GraphicsState,
			resources *pdf.PdfPageResources) error {
			operand := op.Operand
			switch operand {
			case "CS", "cs": // Set colorspace operands.
				op := pdfcontent.ContentStreamOperation{}
				op.Operand = operand
				op.Params = []pdfcore.PdfObject{pdfcore.MakeName("DeviceGray")}
				*processedOperations = append(*processedOperations, &op)
				return nil
			case "SC", "SCN": // Set stroking color.
				color, err := gs.ColorspaceStroking.ColorToRGB(gs.ColorStroking)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					return err
				}
				rgbColor := color.(*pdf.PdfColorDeviceRGB)
				grayColor := rgbColor.ToGray()

				op := pdfcontent.ContentStreamOperation{}
				op.Operand = operand
				op.Params = []pdfcore.PdfObject{pdfcore.MakeFloat(grayColor.Val())}

				*processedOperations = append(*processedOperations, &op)
				return nil
			case "sc", "scn": // Set nonstroking color.
				color, err := gs.ColorspaceNonStroking.ColorToRGB(gs.ColorNonStroking)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					return err
				}
				rgbColor := color.(*pdf.PdfColorDeviceRGB)
				grayColor := rgbColor.ToGray()

				op := pdfcontent.ContentStreamOperation{}
				op.Operand = operand
				op.Params = []pdfcore.PdfObject{pdfcore.MakeFloat(grayColor.Val())}

				*processedOperations = append(*processedOperations, &op)
				return nil
			case "RG", "K": // Set RGB or CMYK stroking color.
				color, err := gs.ColorspaceStroking.ColorToRGB(gs.ColorStroking)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					return err
				}
				rgbColor := color.(*pdf.PdfColorDeviceRGB)
				grayColor := rgbColor.ToGray()

				op := pdfcontent.ContentStreamOperation{}
				op.Operand = "G"
				op.Params = []pdfcore.PdfObject{pdfcore.MakeFloat(grayColor.Val())}

				*processedOperations = append(*processedOperations, &op)
				return nil
			case "rg", "k": // Set RGB or CMYK as nonstroking color.
				color, err := gs.ColorspaceNonStroking.ColorToRGB(gs.ColorNonStroking)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					return err
				}
				rgbColor := color.(*pdf.PdfColorDeviceRGB)
				grayColor := rgbColor.ToGray()

				op := pdfcontent.ContentStreamOperation{}
				op.Operand = "g"
				op.Params = []pdfcore.PdfObject{pdfcore.MakeFloat(grayColor.Val())}

				*processedOperations = append(*processedOperations, &op)
				return nil
			}
			// default:  !@#$
			*processedOperations = append(*processedOperations, op)
			return nil
		})
	// Add handler for image related handling.  Note that inline images are completely stored with a
	// ContentStreamInlineImage object as the parameter for BI.
	processor.AddHandler(pdfcontent.HandlerConditionEnumOperand, "BI",
		func(op *pdfcontent.ContentStreamOperation, gs pdfcontent.GraphicsState, resources *pdf.PdfPageResources) error {
			if len(op.Params) != 1 {
				common.Log.Error("BI Error invalid number of params")
				return errors.New("invalid number of parameters")
			}
			// Inline image.
			iimg, ok := op.Params[0].(*pdfcontent.ContentStreamInlineImage)
			if !ok {
				common.Log.Error("Invalid handling for inline image")
				return errors.New("Invalid inline image parameter")
			}

			img, err := iimg.ToImage(resources)
			if err != nil {
				common.Log.Error("Error converting inline image to image: %v", err)
				return err
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
			grayImage, err := rgbColorSpace.ImageToGray(rgbImg)
			if err != nil {
				common.Log.Error("Error converting img to gray: %v", err)
				return err
			}
			grayInlineImg, err := pdfcontent.NewInlineImageFromImage(grayImage, nil)
			if err != nil {
				common.Log.Error("Error making a new inline image object: %v", err)
				return err
			}

			// Replace inline image data with the gray image.
			pOp := pdfcontent.ContentStreamOperation{}
			pOp.Operand = "BI"
			pOp.Params = []pdfcore.PdfObject{grayInlineImg}
			*processedOperations = append(*processedOperations, &pOp)

			return nil
		})

	// Handler for XObject Image and Forms.
	processedXObjects := map[string]bool{} // Keep track of processed XObjects to avoid repetition.

	processor.AddHandler(pdfcontent.HandlerConditionEnumOperand, "Do",
		func(op *pdfcontent.ContentStreamOperation, gs pdfcontent.GraphicsState, resources *pdf.PdfPageResources) error {
			// operand := op.Operand
			if gVerbose {
				fmt.Printf("Do handler: %s\n", op)
			}
			if len(op.Params) < 1 {
				fmt.Printf("ERROR: Invalid number of params for Do object.\n")
				return errors.New("Range check")
			}

			// XObject.
			name := op.Params[0].(*pdfcore.PdfObjectName)
			common.Log.Debug("Name=%#v=%#q", name, string(*name))

			// Only process each one once.
			_, has := processedXObjects[string(*name)]
			common.Log.Debug("has=%t %+v", has, processedXObjects)
			if has {
				return nil
			}
			processedXObjects[string(*name)] = true

			_, xtype := resources.GetXObjectByName(string(*name))
			common.Log.Debug("xtype=%+v pdf.XObjectTypeImage=%v", xtype, pdf.XObjectTypeImage)
			if xtype == pdf.XObjectTypeImage {
				if gVerbose {
					fmt.Printf(" XObject Image: %s\n", *name)
				}

				ximg, err := resources.GetXObjectImageByName(string(*name))
				if err != nil {
					fmt.Printf("Error : %v\n", err)
					return err
				}

				img, err := ximg.ToImage()
				if err != nil {
					fmt.Printf("Error : %v\n", err)
					return err
				}

				rgbImg, err := ximg.ColorSpace.ImageToRGB(*img)
				if err != nil {
					fmt.Printf("Error : %v\n", err)
					return err
				}

				rgbColorSpace := pdf.NewPdfColorspaceDeviceRGB()
				grayImage, err := rgbColorSpace.ImageToGray(rgbImg)
				if err != nil {
					fmt.Printf("Error : %v\n", err)
					return err
				}

				// Update the XObject image.
				err = ximg.SetImage(&grayImage, nil)
				if err != nil {
					fmt.Printf("Error : %v\n", err)
					return err
				}

				// Update the container.
				_ = ximg.ToPdfObject()
			} else if xtype == pdf.XObjectTypeForm {
				// Go through the XObject Form content stream.
				xform, err := resources.GetXObjectFormByName(string(*name))
				if err != nil {
					fmt.Printf("Error : %v\n", err)
					return err
				}

				formContent, err := xform.GetContentStream()
				if err != nil {
					fmt.Printf("Error : %v\n", err)
					return err
				}

				// Process the content stream in the Form object too:
				// XXX/TODO: Use either form resources (priority) and fall back to page resources
				// alternatively if not found.
				// formResources := xform.FormResources
				formResources := xform.Resources
				if formResources == nil {
					formResources = resources
				}

				// Process the content stream in the Form object too:
				grayContent, err := transformContentStreamToGrayscale(string(formContent), formResources)
				if err != nil {
					common.Log.Error("%v", err)
					return err
				}

				xform.SetContentStream(grayContent)
				// Update the container.
				_ = xform.ToPdfObject()
			}

			return nil
		})

	err = processor.Process(resources)
	if err != nil {
		common.Log.Error("Error processing: %v", err)
		return nil, err
	}

	if gVerbose {
		// For debug purposes: (high level logging).
		fmt.Printf("=== Unprocessed - Full list\n")
		for idx, op := range operations {
			fmt.Printf("U. Operation %d: %s - Params: %v\n", idx+1, op.Operand, op.Params)
		}
		fmt.Printf("=== Processed - Full list\n")
		for idx, op := range *processedOperations {
			fmt.Printf("P. Operation %d: %s - Params: %v\n", idx+1, op.Operand, op.Params)
		}
	}

	return processedOperations.Bytes(), nil
}
