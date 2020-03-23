package main

import (
	"fmt"
	"image"
	// load jpeg decoder
	_ "image/jpeg"
	"os"

	"github.com/unidoc/unipdf/v3/core"
)

func main() {
	// let's read an jpeg rgba image from the file
	// convert it into JBIG2Image using auto threshold
	// and compress the black and white result into another file
	// read an image file 'my-image.jpg'
	f, err := os.Open("checkerboard-squares-black-white.jpg")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	// decode the jpeg image.
	img, _, err := image.Decode(f)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// JBIG2Encoder requires core.JBIG2Image as an input.
	// In order to convert golang image.Image into core.JBIG2Image use
	// a core.GoImageToJBIG2 function.
	// For the RGB and gray scale images there is a 'threshold' which states
	// at which level the values should be black and at which it should be white.
	// It is recommended to use core.JB2ImageAutoThreshold value which computes image histogram.
	// and on it's base gets proper value for the threshold.

	// convert the image into JBIG2Image using auto threshold.
	jb2Img, err := core.GoImageToJBIG2(img, core.JB2ImageAutoThreshold)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// create a jbig2 encoder context.
	enc := &core.JBIG2Encoder{
		// JBIG2 files could be stored as a separate files (mostly with .jb2 extension) or
		// as a part of PDF stream. In this case we want to store it as a file - thus set FileMode to true.
		FileMode: true,
	}
	settings := core.JBIG2EncoderSettings{
		// In order to have better compression, JBIG2 encoder allows to store
		// duplicated lines (image row bits) once, and then relates to it's value on the
		// subsequent rows.
		// In order to use it set DuplicatedLinesRemoval to true.
		DuplicatedLinesRemoval: true,
	}

	// Add JBIG2Image as a new page to the encoder context with the given settings.
	if err = enc.AddPageImage(jb2Img, settings); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Encode the data into jbig2 format.
	data, err := enc.Encode()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// write encoded data into a file with the extension '.jb2' - this is standard extension for the jbig2 files.
	encoded, err := os.Create("checkerboard-squares-black-white.jb2")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer encoded.Close()

	// Write encoded data into the file.
	_, err = encoded.Write(data)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
