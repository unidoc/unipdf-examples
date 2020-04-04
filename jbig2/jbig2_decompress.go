/*
 * This example showcases the decompression of the jbig2 encoded document and storing into
 * commonly used jpg format.
 *
 * Syntax: go run jbig2_decompress.go img.jb2
 */

package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"

	"github.com/unidoc/unipdf/v3/core"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run jbig2_decompress.go img.jb2 ...\n")
		os.Exit(1)
	}
	inputImage := os.Args[1]

	f, err := os.Open(inputImage)
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

	// Store all images in the jpg format.
	saveImage := func(i int, img image.Image) {
		imgFile, err := os.Create(fmt.Sprintf("jbig2_decoded_%d.jpg", i+1))
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
		defer imgFile.Close()

		err = jpeg.Encode(imgFile, img, &jpeg.Options{Quality: 100})
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
	}

	for i, img := range images {
		saveImage(i, img)
	}
}
