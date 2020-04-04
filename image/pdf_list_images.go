/*
 * List images in a PDF file.  Passes through each page, goes through the content stream and finds instances of both
 * XObject Images and inline images. Also handles images referred within XObject Form content streams.
 * Additionally outputs a summary of the filters and colorspaces used by the images found.
 *
 * Run as: go run pdf_list_images.go input.pdf
 */

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/contentstream"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
)

var colorspaces = map[string]int{}
var filters = map[string]int{}

func main() {
	// Enable console debug-level logging when debugging:.
	// common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))

	if len(os.Args) < 2 {
		fmt.Printf("Syntax: go run pdf_list_images.go input.pdf\n")
		os.Exit(1)
	}

	for _, inputPath := range os.Args[1:len(os.Args)] {
		fmt.Printf("Input file: %s\n", inputPath)

		err := listImages(inputPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("=======\nFilter summary:\n")
	for filter, instances := range filters {
		fmt.Printf(" %s: %d instance(s)\n", filter, instances)
	}
	fmt.Printf("=======\nColorspace summary:\n")
	for cs, instances := range colorspaces {
		fmt.Printf(" %s: %d instance(s)\n", cs, instances)
	}
}

// List images and properties of a PDF specified by inputPath.
func listImages(inputPath string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return err
	}

	if isEncrypted {
		// Try decrypting with an empty one.
		auth, err := pdfReader.Decrypt([]byte(""))
		if err != nil {
			return err
		}
		if !auth {
			fmt.Println("Need to decrypt with a specified user/owner password")
			return nil
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}
	fmt.Printf("PDF Num Pages: %d\n", numPages)

	for i := 0; i < numPages; i++ {
		fmt.Printf("-----\nPage %d:\n", i+1)

		page, err := pdfReader.GetPage(i + 1)
		if err != nil {
			return err
		}

		// List images on the page.
		err = listImagesOnPage(page)
		if err != nil {
			return err
		}
	}

	return nil
}

func listImagesOnPage(page *model.PdfPage) error {
	contents, err := page.GetAllContentStreams()
	if err != nil {
		return err
	}

	return listImagesInContentStream(contents, page.Resources)
}

func listImagesInContentStream(contents string, resources *model.PdfPageResources) error {
	cstreamParser := contentstream.NewContentStreamParser(contents)
	operations, err := cstreamParser.Parse()
	if err != nil {
		return err
	}

	processedXObjects := map[string]bool{}

	for _, op := range *operations {
		if op.Operand == "BI" && len(op.Params) == 1 {
			// Inline image.

			iimg, ok := op.Params[0].(*contentstream.ContentStreamInlineImage)
			if !ok {
				continue
			}

			img, err := iimg.ToImage(resources)
			if err != nil {
				return err
			}

			cs, err := iimg.GetColorSpace(resources)
			if err != nil {
				return err
			}

			encoder, err := iimg.GetEncoder()
			if err != nil {
				return err
			}

			fmt.Printf(" Inline image\n")
			fmt.Printf("  Filter: %s\n", encoder.GetFilterName())
			fmt.Printf("  Width: %d\n", img.Width)
			fmt.Printf("  Height: %d\n", img.Height)
			fmt.Printf("  Color components: %d\n", img.ColorComponents)
			fmt.Printf("  ColorSpace: %s\n", cs.String())
			//fmt.Printf("  ColorSpace: %+v\n", cs)
			fmt.Printf("  BPC: %d\n", img.BitsPerComponent)

			// Log filter use globally.
			filter := encoder.GetFilterName()
			if _, has := filters[filter]; has {
				filters[filter]++
			} else {
				filters[filter] = 1
			}
			// Log colorspace use globally.
			csName := "?"
			if cs != nil {
				csName = cs.String()
			}
			if _, has := colorspaces[csName]; has {
				colorspaces[csName]++
			} else {
				colorspaces[csName] = 1
			}
		} else if op.Operand == "Do" && len(op.Params) == 1 {
			// XObject.
			name := op.Params[0].(*core.PdfObjectName)

			// Only process each one once.
			_, has := processedXObjects[string(*name)]
			if has {
				continue
			}
			processedXObjects[string(*name)] = true

			_, xtype := resources.GetXObjectByName(*name)
			if xtype == model.XObjectTypeImage {
				fmt.Printf(" XObject Image: %s\n", *name)

				ximg, err := resources.GetXObjectImageByName(*name)
				if err != nil {
					return err
				}
				img, err := ximg.ToImage()
				if err != nil {
					return err
				}

				fmt.Printf("  Filter: %#v\n", ximg.Filter)
				fmt.Printf("  Width: %v\n", *ximg.Width)
				fmt.Printf("  Height: %d\n", *ximg.Height)
				fmt.Printf("  Color components: %d\n", img.ColorComponents)
				fmt.Printf("  ColorSpace: %s\n", ximg.ColorSpace.String())
				fmt.Printf("  ColorSpace: %#v\n", ximg.ColorSpace)
				fmt.Printf("  BPC: %v\n", *ximg.BitsPerComponent)

				// Log filter use globally.
				filter := ximg.Filter.GetFilterName()
				if _, has := filters[filter]; has {
					filters[filter]++
				} else {
					filters[filter] = 1
				}
				// Log colorspace use globally.
				cs := ximg.ColorSpace.String()
				if _, has := colorspaces[cs]; has {
					colorspaces[cs]++
				} else {
					colorspaces[cs] = 1
				}
			} else if xtype == model.XObjectTypeForm {
				// Go through the XObject Form content stream.
				fmt.Printf("--> XObject Form: %s\n", *name)
				xform, err := resources.GetXObjectFormByName(*name)
				if err != nil {
					return err
				}

				formContent, err := xform.GetContentStream()
				if err != nil {
					return err
				}
				fmt.Printf("xform: %#v\n", xform)
				fmt.Printf("xform res: %#v\n", xform.Resources)
				fmt.Printf("Content: %s\n", formContent)

				// Process the content stream in the Form object too:
				// XXX/TODO: Use either form resources (priority) and fall back to page resources alternatively if not found.
				if xform.Resources != nil {
					err = listImagesInContentStream(string(formContent), xform.Resources)
				} else {
					err = listImagesInContentStream(string(formContent), resources)
				}
				if err != nil {
					return err
				}
				fmt.Printf("<-- XObject Form: %s\n", *name)
			}
		}
	}

	return nil
}
