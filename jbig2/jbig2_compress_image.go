/*
 * This example showcases the conversion of the jpg encoded image
 * into jbig2 encoding format.
 *
 * The input file would be firstly converted into bi-level image
 * and then stored using jbig2 encoding format.
 *
 * Syntax: go run jbig2_compress_image.go img.jpg
 */

package main

import (
	"fmt"
	"image"
	// load jpeg image decoder
	_ "image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/unidoc/unipdf/v3/core"
)

func main() {
	// Let's read an jpeg rgba image from the file, convert it into JBIG2Image
	// using auto threshold and compress the black and white result into another file.
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run jbig2_compress_image.go img.jpg ...\n")
		os.Exit(1)
	}

	inputImage := os.Args[1]
	_, fileName := filepath.Split(inputImage)

	f, err := os.Open(inputImage)
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
	encodedFile, err := os.Create(strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".jb2")
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	defer encodedFile.Close()

	_, err = encodedFile.Write(data)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	fmt.Println("Created JBIG2 Encoded file successfully")
}
