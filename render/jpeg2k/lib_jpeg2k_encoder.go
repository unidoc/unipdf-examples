/*
 *
 * An implementation example of a custom encoder for handling JPEG2000 image format.
 *
 */

package jpeg2k

import (
	"errors"
	"fmt"
	"math"

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

	imgWidth := mw.GetImageWidth()
	imgHeight := mw.GetImageHeight()
	imgDepth := mw.GetImageDepth()
	cs := mw.GetImageColorspace()

	imw := mw.NewPixelIterator()
	defer imw.Destroy()

	var colorComponent uint

	switch cs {
	case imagick.COLORSPACE_GRAY:
		colorComponent = 1
	case imagick.COLORSPACE_RGB, imagick.COLORSPACE_SRGB:
		colorComponent = 3
	case imagick.COLORSPACE_CMYK:
		colorComponent = 4
	default:
		return nil, fmt.Errorf("Provided image with unsupported colorspace: %v", cs)
	}

	var depthMultiplier float64
	switch imgDepth {
	case 8:
		depthMultiplier = float64(math.MaxUint8)
	case 16:
		depthMultiplier = float64(math.MaxUint16)
	default:
		return nil, fmt.Errorf("Unsupported image bit depth: %v", imgDepth)
	}

	var decoded = make([]byte, imgWidth*imgHeight*colorComponent*imgDepth/8)
	var index int

	for y := 0; y < int(imgHeight); y++ {
		pmw := imw.GetNextIteratorRow()

		for x := 0; x < int(imgWidth); x++ {
			switch cs {
			case imagick.COLORSPACE_GRAY:
				g := uint16(pmw[x].GetRed() * depthMultiplier)
				if imgDepth == 16 {
					decoded[index] = byte((g >> 8) & 0xff)
					index++
					decoded[index] = byte(g & 0xff)
					index++
				} else {
					decoded[index] = uint8(g) & 0xff
					index++
				}

			case imagick.COLORSPACE_RGB, imagick.COLORSPACE_SRGB:
				// Get the RGB from the source image
				r := uint16(pmw[x].GetRed() * depthMultiplier)
				g := uint16(pmw[x].GetGreen() * depthMultiplier)
				b := uint16(pmw[x].GetBlue() * depthMultiplier)

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

			case imagick.COLORSPACE_CMYK:
				c := uint16(pmw[x].GetCyan() * depthMultiplier)
				m := uint16(pmw[x].GetMagenta() * depthMultiplier)
				y := uint16(pmw[x].GetYellow() * depthMultiplier)
				k := uint16(pmw[x].GetBlack() * depthMultiplier)

				if imgDepth == 16 {
					decoded[index] = byte((c >> 8) & 0xff)
					index++
					decoded[index] = byte(c & 0xff)
					index++
					decoded[index] = byte((m >> 8) & 0xff)
					index++
					decoded[index] = byte(m & 0xff)
					index++
					decoded[index] = byte((y >> 8) & 0xff)
					index++
					decoded[index] = byte(y & 0xff)
					index++
					decoded[index] = byte((k >> 8) & 0xff)
					index++
					decoded[index] = byte(k & 0xff)
					index++
				} else {
					decoded[index] = byte(c & 0xff)
					index++
					decoded[index] = byte(m & 0xff)
					index++
					decoded[index] = byte(y & 0xff)
					index++
					decoded[index] = byte(k & 0xff)
					index++
				}
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
