/*
 * This example showcases the conversion of the jpg encoded image
 * into jbig2 encoding format.
 *
 * The result jbig2 file is compressed with a lossless method stored with
 * the standard  jbig2 - .jb2 extension. The file compressed in this example
 * is stored in just 377 bytes - compared to the 142 794 bytes for original image.
 * This gives compression ratio (uncompressed size/compressed size) of 378.76, which leads to 99.735%
 * space savings for given example.
 */

package main

import (
	"fmt"
	"image"
	"log"

	// load jpeg decoder
	_ "image/jpeg"
	"os"

	"github.com/unidoc/unipdf/v3/core"
)

func main() {
	// Let's read an jpeg rgba image from the file, convert it into JBIG2Image
	// using auto threshold and compress the black and white result into another file.
	f, err := os.Open("checkerboard-squares-black-white.jpg")
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	// JBIG2Encoder requires core.JBIG2Image as an input.
	// In order to convert golang image.Image into core.JBIG2Image use
	// the core.GoImageToJBIG2 function.
	// For the RGB and gray scale images there is a 'threshold' which states
	// at which level the values should be black and at which it should be white.
	// It is recommended to use core.JB2ImageAutoThreshold value which computes image histogram.
	// and on it's base gets proper value for the threshold.
	// Convert the image into JBIG2Image using auto threshold.
	jb2Img, err := core.GoImageToJBIG2(img, core.JB2ImageAutoThreshold)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	// Create a JBIG2 Encoder/Decoder context. In this example we're setting page settings
	// used for the encoding process.
	enc := &core.JBIG2Encoder{
		DefaultPageSettings: core.JBIG2EncoderSettings{
			// JBIG2 files could be stored as a separate files (mostly with .jb2 extension) or
			// as a part of PDF stream. In this case we want to store it as a file - thus set FileMode to true.
			FileMode: true,
			// In order to have better compression, JBIG2 encoder allows storing
			// duplicated lines (image row bits) once, and then relates to their value on the
			// subsequent rows.
			// In order to use it set DuplicatedLinesRemoval to true.
			DuplicatedLinesRemoval: true,
		},
	}
	// Add JBIG2Image as a new page to the encoder context with the default page settings.
	if err = enc.AddPageImage(jb2Img, nil); err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	// Encode the data into jbig2 format and return as the slice of bytes.
	data, err := enc.Encode()
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	// Write encoded data into a file with the extension '.jb2' - this is standard extension for the jbig2 files.
	encodedFile, err := os.Create("checkerboard-squares-black-white.jb2")
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	defer encodedFile.Close()

	_, err = encodedFile.Write(data)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	fmt.Printf("Created JBIG2 Encoded file successfully")
}
