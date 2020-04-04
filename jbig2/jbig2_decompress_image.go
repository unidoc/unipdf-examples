/*
 * This example showcases the decompression of the jbig2 encoded image and storing into
 * commonly used jpg format.
 *
 * As the input for this test the result of the compression example would be used (lossless image).
 * The result of this example is an image that has unchanged quality (in compare to the input of the compression example).
 */

package main

import (
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"

	"github.com/unidoc/unipdf/v3/core"
)

func main() {
	f, err := os.Open("checkerboard-squares-black-white.jb2")
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	// Create new JBIG2 Encoder/Decoder context.
	enc := &core.JBIG2Encoder{}

	// Decode all images from the 'data'.
	images, err := enc.DecodeImages(data)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	// The 'checkerboard-squares-black-white.jb2 file should have exactly one image stored.
	if len(images) != 1 {
		log.Fatalf("Error: Only a single image should be decoded\n")
	}
	// Create a new file for the decoded image.
	dec, err := os.Create("checkerboard-squares-black-white-decoded.jpg")
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	defer dec.Close()
	if err = jpeg.Encode(dec, images[0], &jpeg.Options{Quality: 100}); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}
