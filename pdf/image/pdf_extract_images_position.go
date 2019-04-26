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
 *       1 bit per pixel, so I just set imgMark.img.BitsPerComponent = 8 for CCITTFaxEncoder images.
 *
 * TODO: Handle JBIG images.
 *       Handle CCITTFaxEncoder inline images?
 *       Handle CCITTFaxEncoder Group 3 images?
 *       Change got.6 to return 1 bit images.
 *       Save images in orientation they appear in the PDF file.
 *
 * Run as: go run pdf_extract_native_images.go input.pdf output.folder


    ~/testdata/misc/rotate_2.pdf  Images with different orientations.

    Clipped image
    ~/testdata/other/pdf/itextpdf/xtra/src/test/resources/com/itextpdf/text/pdf/pdfcleanup/BigImage-tif.pdf

     ~/testdata/other/pdf/itextpdf/xtra/src/test/resources/com/itextpdf/text/pdf/pdfcleanup/cmp_rotatedImg.pdf"
  Width: 1000
  Height: 250
  Size 500.0x500.0

   Scales incorrectly. Unscaled image is correct.
   go run pdf_extract_images_position.go -p 5 -i 2 -s ~/testdata/programming/images/deskew/50682881b388d2a5c534caab6031db39e61d.pdf

   Looks correct
   go run pdf_extract_images_position.go -p 25 -i 3 ~/testdata/programming/digital_signatures/guide.pdf

   Color is reversed.
   go run pdf_extract_images_position.go -p 7 -i 1 ~/testdata/misc/slides_12.pdf

   Inline image
   ~/testdata/other/pdf/itextpdf/itext/src/test/resources/com/itextpdf/text/pdf/parser/PdfContentStreamProcessorTest/inlineImages01.pdf
   go run pdf_extract_images_position.go -p 443 -i 1 ~/testdata/other/unidoc-examples/pdf/testing/out.1/psrefman.pdf

   Rotated image?
   ~/testdata/science/climate/environment/20190117-105663-2.pdf

*/

package main

import (
	"flag"
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/unidoc/unidoc/common"
	pdfcore "github.com/unidoc/unidoc/pdf/core"
	"github.com/unidoc/unidoc/pdf/extractor"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

const usage = "Usage: go run df_extract_images.go input.pdf output.folder\n"

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
	var debug, trace bool
	var outputDir string
	var pageNum, imageNum int
	var unscaled bool
	flag.BoolVar(&debug, "d", false, "Print debugging information.")
	flag.BoolVar(&trace, "e", false, "Print detailed debugging information.")
	flag.StringVar(&outputDir, "o", "extracted.images", "Directory where extracted images are saved.")
	flag.IntVar(&pageNum, "p", 0, "Extract images only from this page number (1-offset).")
	flag.IntVar(&imageNum, "i", 0, "Extract images only with this image number (1-offset).")
	flag.BoolVar(&unscaled, "s", false, "Don't rescale.")
	makeUsage(usage)

	flag.Parse()
	args := flag.Args()

	if len(args) < 1 || len(outputDir) == 0 {
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

	fmt.Printf("Input file: %s\n", inputPath)
	err := extractImagesToFolder(inputPath, outputDir, pageNum, imageNum, unscaled)
	if err != nil {
		fmt.Printf("ERROR: Could not process inputPath=%q outputDir=%q err=%v\n",
			inputPath, outputDir, err)
		os.Exit(1)
	}
}

// extractImagesToFolder extracts images and properties of a PDF file specified by `inputPath`.
// The output images are stored into a directory whose path is given by `outputDir`.
func extractImagesToFolder(inputPath, outputDir string, pageN, imageN int, unscaled bool) error {
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

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.Mkdir(outputDir, 0777)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not create outputDir=%q err=%v\n", outputDir, err)
			return err
		}
	}

	name := filepath.Base(inputPath)
	{
		ext := filepath.Ext(name)
		name = name[:len(name)-len(ext)]
	}

	for pageNum := 1; pageNum <= numPages; pageNum++ {
		if pageN > 0 && pageNum != pageN {
			continue
		}

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return err
		}

		// Extract images on the page.
		images, err := extractImagesOnPage(page)
		if err != nil {
			return err
		}

		mbox, err := page.GetMediaBox()
		if err != nil {
			panic(err)
		}

		for idx, imgMark := range images.Images {
			if imageN > 0 && idx+1 != imageN {
				continue
			}
			fname := fmt.Sprintf("%s_page%d_img%d", name, pageNum, idx+1)

			gimg, err := imgMark.PageView(*mbox, !unscaled, !unscaled, false)
			if err != nil {
				return err
			}

			fname = addSuffix(fname, imgMark.Inline, ".i", ".x")
			fname = addSuffix(fname, unscaled, ".unscaled", "")
			fname = addSuffix(fname, imgMark.Lossy, ".jpg", ".png")
			outputPath := filepath.Join(outputDir, fname)

			fmt.Printf("  Converting to go image: page %d img %d ⇾ %q mbox=%+v\n  imgMark=%s\n",
				pageNum, idx+1, outputPath, *mbox, imgMark.String())

			imgf, err := os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE, 0755)
			if err != nil {
				return err
			}
			defer imgf.Close()

			if !imgMark.Lossy {
				err = png.Encode(imgf, gimg)
			} else {
				opt := jpeg.Options{Quality: 100}
				err = jpeg.Encode(imgf, gimg, &opt)
			}
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type imageData struct {
	img    *pdf.Image
	filter pdfcore.StreamEncoder
	inline bool
}

func addSuffix(prefix string, condition bool, is, isnt string) string {
	if condition {
		return prefix + is
	}
	return prefix + isnt
}

func count(data []imageData, inline bool) int {
	n := 0
	for _, d := range data {
		if d.inline == inline {
			n++
		}
	}
	return n
}

// extractImagesOnPage returns a slice of all images on page `page`.
func extractImagesOnPage(page *pdf.PdfPage) (*extractor.PageImages, error) {
	pageExtractor, err := extractor.New(page)
	if err != nil {
		return nil, err
	}
	return pageExtractor.ExtractPageImages(nil)
}

// func describeImage(img *pdf.Image) string {
// 	desc := fmt.Sprintf("%dx%d cpts=%d bpp=%d",
// 		img.Width, img.Height, img.ColorComponents, img.BitsPerComponent)
// 	if len(desc) > 100 {
// 		panic("ddd")
// 	}
// 	return desc
// }

// func describeImageMark(mark extractor.ImageMark) string {
// 	return mark.String()
// 	// desc := fmt.Sprintf("%.1fx%.1f (%.1f,%.1f) ϴ=%.1f img=[%s]",
// 	// 	mark.Width, mark.Height, mark.X, mark.Y, mark.Angle, describeImage(mark.Image))
// 	// // if math.Abs(img.Angle) > 0.1 {
// 	// // 	panic("ppppp")
// 	// // }
// 	// return desc
// }

// makeUsage updates flag.Usage to include usage message `msg`.
func makeUsage(msg string) {
	usage := flag.Usage
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, msg)
		usage()
	}
}

// func inverseTransform(imgMark extractor.ImageMark) core.Matrix {
// 	tx, ty, w, h, θ := imgMark.X, imgMark.Y, imgMark.Width, imgMark.Height, imgMark.Angle
// 	c, s := math.Cos(θ), math.Sin(θ)
// 	core.NewMatrix(c/w, s/w, -c/n, c/h, tx, ty float64)
// 	cos(θ)
// 	// ⎡⎣⎢cos(α)−sin(α)0sin(α)cos(α)0−Txcos(α)−Tysin(α)−Tycos(α)+Txsin(α)1⎤⎦⎥
// }

// func crop(imgMark extractor.ImageMark, mbox core.PdfRectangl) *pdf.Image {
// 	imgBox := asRect(imgMark)
// 	cropBox := intersection(imgBox, mbox)
// 	imgW, imgH := dimensions(imgBox)
// 	cropW, cropH := dimensions(cropBox)

// 	img := imgMark.Image
// 	if cropW >= imgW && cropH >= imgH {
// 		return img
// 	}
// }

// func dimensions(r core.PdfRectangle) (float64, float64) {
// 	return math.Abs(r.Urx - r.Llx), math.Abs(r.Ury - r.Lly)
// }

// func asRect(imgMark extractor.ImageMark) core.PdfRectangle {
// 	r := core.PdfRectangle{
// 		Llx: 0,
// 		Lly: 0,
// 		Urx: 1,
// 		Ury: 1,
// 	}
// 	r.Llx *= imgMark.Width
// 	r.Lly *= imgMark.Height
// 	r.Urx *= imgMark.Width
// 	r.Ury *= imgMark.Height

// 	r.Llx, r.Lly = rotateXY(r.Llx, r.Lly, imgMark.Angle)
// 	r.Urx, r.Ury = rotateXY(r.Urx, r.Ury, imgMark.Angle)

// 	r.Llx += imgMark.X
// 	r.Lly += imgMark.Y
// 	r.Urx += imgMark.X
// 	r.Ury += imgMark.Y

// 	return y
// }

// func rotateXY(x, y, θ float64) (float64, float64) {
// 	c, s := math.Cos(θ), math.Sin(θ)
// 	return x*c + y*s, -x*s + y*c
// }

// func intersection(r1, r2 core.PdfRectangle) core.PdfRectangle {
// 	r := core.PdfRectangle{
// 		Llx: maxFloat(r1.Llx, r2.Llx),
// 		Lly: maxFloat(r1.Lly, r2.Lly),
// 		Urx: minFloat(r1.Urx, r2.Urx),
// 		Ury: minFloat(r1.Ury, r2.Ury),
// 	}
// 	return r
// }

// // minFloat returns the lesser of `a` and `b`.
// func minFloat(a, b float64) float64 {
// 	if a < b {
// 		return a
// 	}
// 	return b
// }

// // maxFloat returns the greater of `a` and `b`.
// func maxFloat(a, b float64) float64 {
// 	if a > b {
// 		return a
// 	}
// 	return b
// }
