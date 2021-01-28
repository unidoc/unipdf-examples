/*
 *
 * An implementation example of a custom encoder for handling JPEG2000 image format.
 *
 */

package jpeg2k

import (
	"errors"
	"fmt"

	"github.com/unidoc/unipdf/v3/core"
	"gopkg.in/gographics/imagick.v2/imagick"
)

// Custom encoder declaration that implements StreamEncoder interface.
type CustomJPXEncoder struct{}

// NewCustomJPXEncoder returns a new instance of CustomJPXEncoder.
func NewCustomJPXEncoder() *CustomJPXEncoder {
	return &CustomJPXEncoder{}
}

// GetFilterName returns the name of the encoding filter.
func (enc *CustomJPXEncoder) GetFilterName() string {
	return "CustomJPXEncoder"
}

// MakeDecodeParams makes a new instance of an encoding dictionary based on
// the current encoder settings.
func (enc *CustomJPXEncoder) MakeDecodeParams() core.PdfObject {
	return nil
}

// MakeStreamDict makes a new instance of an encoding dictionary for a stream object.
func (enc *CustomJPXEncoder) MakeStreamDict() *core.PdfObjectDictionary {
	dict := core.MakeDict()

	dict.Set("Filter", core.MakeName(enc.GetFilterName()))

	return dict
}

// UpdateParams updates the parameter values of the encoder.
func (enc *CustomJPXEncoder) UpdateParams(params *core.PdfObjectDictionary) {
}

// DecodeBytes decodes a slice of JPX encoded bytes and returns the result.
func (enc *CustomJPXEncoder) DecodeBytes(encoded []byte) ([]byte, error) {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	if err := mw.ReadImageBlob(encoded); err != nil {
		return nil, err
	}

	// fmt.Println(mw.IdentifyImage())

	imgWidth := mw.GetImageWidth()
	imgHeight := mw.GetImageHeight()
	imgDepth := mw.GetImageDepth()
	cs := mw.GetImageColorspace()

	fmt.Printf("Color space %v - %d bit\n", cs, imgDepth)

	imw := mw.NewPixelIterator()
	defer imw.Destroy()

	colorComponent := 3

	switch cs {
	case imagick.COLORSPACE_GRAY:
		colorComponent = 1
	case imagick.COLORSPACE_CMYK:
		colorComponent = 4
	}

	var decoded = make([]byte, imgWidth*imgHeight*uint(colorComponent)*imgDepth/8)
	index := 0

	for y := 0; y < int(imgHeight); y++ {
		pmw := imw.GetNextIteratorRow()

		for x := 0; x < int(imgWidth); x++ {
			if cs == imagick.COLORSPACE_GRAY {
				r := int(pmw[x].GetRed() * 255)
				if imgDepth == 16 {
					decoded[index] = byte((r >> 8) & 0xff)
					index++
					decoded[index] = byte(r & 0xff)
					index++
				} else {
					decoded[index] = uint8(r) & 0xff
					index++
				}
			} else if cs == imagick.COLORSPACE_RGB || cs == imagick.COLORSPACE_SRGB {
				// Get the RGB quanta from the source image
				r := int(pmw[x].GetRed() * 255)
				g := int(pmw[x].GetGreen() * 255)
				b := int(pmw[x].GetBlue() * 255)

				if imgDepth == 16 {
					decoded[index] = byte((r >> 8) & 0xff)
					index++
					decoded[index] = byte(r & 0xff)
					index++
					decoded[index] = byte((g >> 8) & 0xff)
					index++
					decoded[index] = byte(g & 0xff)
					index++
					decoded[index] = byte((b >> 8) & 0xff)
					index++
					decoded[index] = byte(b & 0xff)
					index++
				} else {
					decoded[index] = byte(r & 0xff)
					index++
					decoded[index] = byte(g & 0xff)
					index++
					decoded[index] = byte(b & 0xff)
					index++
				}
			} else {
				fmt.Printf("Unrecognized encoding: %v", cs)
			}
		}
	}

	return decoded, nil
}

// DecodeStream decodes a JPX encoded stream and returns the result as a
// slice of bytes.
func (enc *CustomJPXEncoder) DecodeStream(streamObj *core.PdfObjectStream) ([]byte, error) {
	return enc.DecodeBytes(streamObj.Stream)
}

// EncodeBytes JPX encodes the passed in slice of bytes.
func (enc *CustomJPXEncoder) EncodeBytes(data []byte) ([]byte, error) {
	return data, errors.New("Custom encode bytes not yet implemented")
}
