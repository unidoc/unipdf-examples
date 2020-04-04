/*
 * This example showcases the decompression of the jbig2 encoded document
 * with additional 'globals' jbig2 file.
 *
 * JBIG2 standard allows to store common segment definitions called 'Globals' which may be stored
 * on separate byte stream or file. In order to decode this kind of documents, we need to
 * firstly decode the globals file (jbig2_example_globals.jb2), and then apply it to the JBIG2Encoder context that would be used
 * to decode the main file (jbig2_example_globals.jb2).
 */

package main

import (
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"

	"github.com/unidoc/unipdf/v3/core"
)

func main() {
	// JBIG2 Files could also be stored in two files/ data streams.
	// The Globals which contains globally defined jbig2 segments
	// and the main jb2 file.
	// In the PDF stream this is done automatically by the UniPDF library.
	// In order to decode it manually a user is responsible for decoding JBIG2 Globals.
	// At first we need to read and decode Globals file.
	globalsFile, err := os.Open("jbig2_example_globals.jb2")
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	defer globalsFile.Close()

	// Read all data from the globals file.
	globalsData, err := ioutil.ReadAll(globalsFile)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	// Create JBIG2 Encoder/Decoder context for the globals file.
	globalsDecoder := &core.JBIG2Encoder{}

	// Decode the globals using 'DecodeGlobals' method.
	globals, err := globalsDecoder.DecodeGlobals(globalsData)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	// Now read the main JBIG2 file and decode it with the the use of provided 'globals'.
	jbig2File, err := os.Open("jbig2_example.jb2")
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	defer jbig2File.Close()

	exampleFileData, err := ioutil.ReadAll(jbig2File)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	// Create new JBIG2 Decoder context with previously decoded 'globals' and decode the images.
	enc := &core.JBIG2Encoder{
		Globals: globals,
	}

	images, err := enc.DecodeImages(exampleFileData)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	for i, img := range images {
		imgFile, err := os.Create(fmt.Sprintf("jbig2_example_decoded_%d", i+1))
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
		defer imgFile.Close()
		err = jpeg.Encode(imgFile, img, &jpeg.Options{Quality: 100})
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
	}
	fmt.Printf("Decoded %d images.", len(images))
}
