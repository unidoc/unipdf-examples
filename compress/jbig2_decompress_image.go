package main

import (
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"os"

	"github.com/unidoc/unipdf/core"
)

func main() {
	// read jbig2 encoded file with the checkerboard-squares-black-white image.
	f, err := os.Open("checkerboard-squares-black-white.jb2")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	// read all bytes from the file
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Create new JBIG2Encoder/Decoder
	enc := &core.JBIG2Encoder{}

	// Decode all images from the 'data'.
	images, err := enc.DecodeImages(data)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	// there should be exactly one image.
	if len(images) != 1 {
		fmt.Printf("Error: Only a single image should be decoded\n")
		os.Exit(1)
	}
	// Create a new decoded file.
	dec, err := os.Create("checkerboard-squares-black-white-decoded.jpg")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer dec.Close()

	// encode using commonly used jpeg.
	if err = jpeg.Encode(dec, images[0], &jpeg.Options{Quality: 100}); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
