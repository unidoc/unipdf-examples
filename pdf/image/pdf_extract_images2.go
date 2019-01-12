/*
 * Extract images from a PDF file. Passes through each page, goes through the content stream and finds instances of both
 * XObject Images and inline images. Also handles images referred within XObject Form content streams.
 * The output files are saved as a zip archive.
 *
 * Also extracts display position and dimensions from the PDF based on the GraphicsState and CTM.
 * The coordinates extracted are in the PDF coordinate system (origin in lower left corner).
 *
 * Run as: go run pdf_extract_images2.go input.pdf output.zip
 */

package main

import (
	"archive/zip"
	"fmt"
	"image/jpeg"
	"os"

	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

var xObjectImages = 0
var inlineImages = 0

func main() {
	// Enable debug-level console logging, when debuggingn:
	//unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	if len(os.Args) < 3 {
		fmt.Printf("Syntax: go run pdf_extract_images2.go input.pdf output.zip\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	fmt.Printf("Input file: %s\n", inputPath)
	err := extractImagesToArchive(inputPath, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("-- Summary\n")
	fmt.Printf("%d XObject images extracted\n", xObjectImages)
	fmt.Printf("%d inline images extracted\n", inlineImages)
	fmt.Printf("Total %d images\n", xObjectImages+inlineImages)
}

// Extracts images and properties of a PDF specified by inputPath.
// The output images are stored into a zip archive whose path is given by outputPath.
func extractImagesToArchive(inputPath, outputPath string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return err
	}

	// Try decrypting with an empty one.
	if isEncrypted {
		auth, err := pdfReader.Decrypt([]byte(""))
		if err != nil {
			// Encrypted and we cannot do anything about it.
			return err
		}
		if !auth {
			fmt.Println("Need to decrypt with password")
			return nil
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}
	fmt.Printf("PDF Num Pages: %d\n", numPages)

	// Prepare output archive.
	zipf, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer zipf.Close()
	zipw := zip.NewWriter(zipf)

	for i := 0; i < numPages; i++ {
		fmt.Printf("-----\nPage %d:\n", i+1)

		page, err := pdfReader.GetPage(i + 1)
		if err != nil {
			return err
		}

		// List images on the page.
		images, err := extractImagesOnPage(page)
		if err != nil {
			return err
		}

		for idx, img := range images {
			fname := fmt.Sprintf("p%d_%d.jpg", i+1, idx)
			fmt.Println(fname)
			fmt.Println(img)

			gimg, err := img.image.ToGoImage()
			if err != nil {
				return err
			}

			imgf, err := zipw.Create(fname)
			if err != nil {
				return err
			}
			opt := jpeg.Options{Quality: 100}
			err = jpeg.Encode(imgf, gimg, &opt)
			if err != nil {
				return err
			}
		}
	}

	// Make sure to check the error on Close.
	err = zipw.Close()
	if err != nil {
		return err
	}

	return nil
}

// imageObject represents an extracted image with image data, position and dimensions.
type imageObject struct {
	image *pdf.Image

	// Dimensions of the image as displayed in the PDF.
	width  float64
	height float64

	// Position of the image in PDF coordinates (lower left corner).
	x float64
	y float64

	// Angle if rotated.
	angle float64
}

func (i imageObject) String() string {
	return fmt.Sprintf("[X=%.2f Y=%.2f] [W=%.2f H=%2f] [angle=%.2f]", i.x, i.y, i.width, i.height, i.angle)
}

func extractImagesOnPage(page *pdf.PdfPage) ([]imageObject, error) {
	contents, err := page.GetAllContentStreams()
	if err != nil {
		return nil, err
	}

	return extractImagesInContentStream(contents, page.Resources)
}

func extractImagesInContentStream(contents string, resources *pdf.PdfPageResources) ([]imageObject, error) {
	var extractedImages []imageObject
	cstreamParser := pdfcontent.NewContentStreamParser(contents)
	operations, err := cstreamParser.Parse()
	if err != nil {
		return nil, err
	}

	processedXObjects := map[string]bool{}

	processor := pdfcontent.NewContentStreamProcessor(*operations)
	processor.AddHandler(pdfcontent.HandlerConditionEnumAllOperands, "",
		func(op *pdfcontent.ContentStreamOperation, gs pdfcontent.GraphicsState, resources *pdf.PdfPageResources) error {
			if op.Operand == "BI" && len(op.Params) == 1 {
				// BI: Inline image.

				iimg, ok := op.Params[0].(*pdfcontent.ContentStreamInlineImage)
				if !ok {
					return nil
				}

				img, err := iimg.ToImage(resources)
				if err != nil {
					return err
				}

				cs, err := iimg.GetColorSpace(resources)
				if err != nil {
					return err
				}
				if cs == nil {
					// Default if not specified?
					cs = pdf.NewPdfColorspaceDeviceGray()
				}
				fmt.Printf("Cs: %T\n", cs)

				rgbImg, err := cs.ImageToRGB(*img)
				if err != nil {
					return err
				}
				fmt.Printf("@BI CTM: %s\n", gs.CTM.String())
				xDim := gs.CTM.ScalingFactorX()
				yDim := gs.CTM.ScalingFactorY()
				xPos, yPos := gs.CTM.Translation()
				angle := gs.CTM.Angle()
				fmt.Printf("Size: (%v, %v) at (%v, %v), angle: %v\n", xDim, yDim, xPos, yPos, angle)
				imgObj := imageObject{
					image:  &rgbImg,
					x:      xPos,
					y:      yPos,
					width:  xDim,
					height: yDim,
				}

				extractedImages = append(extractedImages, imgObj)
				inlineImages++
			} else if op.Operand == "Do" && len(op.Params) == 1 {
				// Do: XObject.
				name := op.Params[0].(*pdfcore.PdfObjectName)

				// Only process each one once.
				_, has := processedXObjects[string(*name)]
				if has {
					// Already processed.
					return nil
				}
				processedXObjects[string(*name)] = true

				_, xtype := resources.GetXObjectByName(*name)
				if xtype == pdf.XObjectTypeImage {
					fmt.Printf(" XObject Image: %s\n", *name)

					ximg, err := resources.GetXObjectImageByName(*name)
					if err != nil {
						return err
					}

					img, err := ximg.ToImage()
					if err != nil {
						return err
					}

					fmt.Printf("@Do CTM: %s\n", gs.CTM.String())
					xDim := gs.CTM.ScalingFactorX()
					yDim := gs.CTM.ScalingFactorY()
					xPos, yPos := gs.CTM.Translation()
					angle := gs.CTM.Angle()
					fmt.Printf("Size: (%v, %v) at (%v, %v), angle: %v\n", xDim, yDim, xPos, yPos, angle)

					rgbImg, err := ximg.ColorSpace.ImageToRGB(*img)
					if err != nil {
						return err
					}
					imgObj := imageObject{
						image:  &rgbImg,
						x:      xPos,
						y:      yPos,
						width:  xDim,
						height: yDim,
					}

					extractedImages = append(extractedImages, imgObj)
					xObjectImages++
				} else if xtype == pdf.XObjectTypeForm {
					// Go through the XObject Form content stream.
					xform, err := resources.GetXObjectFormByName(*name)
					if err != nil {
						return err
					}

					formContent, err := xform.GetContentStream()
					if err != nil {
						return err
					}

					// Process the content stream in the Form object too:
					formResources := xform.Resources
					if formResources == nil {
						formResources = resources
					}

					// Process the content stream in the Form object too:
					formImages, err := extractImagesInContentStream(string(formContent), formResources)
					if err != nil {
						return err
					}
					extractedImages = append(extractedImages, formImages...)
				}
			}
			return nil
		})

	err = processor.Process(resources)
	if err != nil {
		return nil, err
	}

	return extractedImages, nil
}
