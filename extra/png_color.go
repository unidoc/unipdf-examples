/*
 * Check a directory of PNG files and check if any contain color
 */

package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/unidoc/unidoc/common"
)

func main() {
	debug := false // Write debug level info to stdout?
	keep := false  // Keep the rasters used for PDF comparison"
	flag.BoolVar(&debug, "d", false, "Enable debug logging")
	flag.BoolVar(&keep, "k", false, "Keep the difference rasters")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-d -k] <PNG directory> ...\n", os.Args[0])
		os.Exit(1)
	}

	dir := args[0]

	numPages, colorPages, err := colorDirectoryPages("*.png", dir, keep)

	if err != nil {
		fmt.Fprintf(os.Stderr, "colorDirectoryPages failed. err=%v", err)
	}

	fmt.Printf("%d color of %d pages %+v\n", len(colorPages), numPages, colorPages)
}

var (
	gsImagePattern = `(\d+)\.png$`
	gsImageRegex   = regexp.MustCompile(gsImagePattern)
)

// colorDirectoryPages returns a lists of the page numbers of the image files that match `mask` in
// directories `dir` that have any color pixels.
func colorDirectoryPages(mask, dir string, keep bool) (int, []int, error) {
	pattern := filepath.Join(dir, mask)
	fmt.Printf("pattern=%q\n", pattern)
	files, err := filepath.Glob(pattern)
	if err != nil {
		common.Log.Error("colorDirectoryPages: Glob failed. pattern=%#q err=%v", pattern, err)
		return 0, nil, err
	}

	numPages := 0
	colorPages := []int{}
	for _, path := range files {
		fmt.Printf("%s\n", path)
		matches := gsImageRegex.FindStringSubmatch(path)
		if len(matches) == 0 {
			continue
		}
		pageNum, err := strconv.Atoi(matches[1])
		if err != nil {
			panic(err)
			return numPages, colorPages, err
		}
		numPages++
		// common.Log.Error("isColorDirectory:  path=%#q", path)
		isColor, err := isColorImage(path, keep)
		// common.Log.Error("isColorDirectory: isColor=%t path=%#q", isColor, path)
		if err != nil {
			panic(err)
			return numPages, colorPages, err
		}
		if isColor {
			colorPages = append(colorPages, pageNum)
			// common.Log.Error("isColorDirectory: colorPages=%d %d", len(colorPages), colorPages)
		}
	}
	return numPages, colorPages, nil
}

// isColorImage returns true if image file `path` contains color
func isColorImage(path string, keep bool) (bool, error) {
	img, err := readImage(path)
	if err != nil {
		return false, err
	}
	isColor := imgIsColor(img)
	if isColor && keep {
		markedPath := fmt.Sprintf("%s.marked.png", path)
		markedImg, summary := imgMarkColor(img)
		common.Log.Error("markedPath=%#q %s", markedPath, summary)
		err = writeImage(markedPath, markedImg)
	}
	return isColor, err
}

const colorThreshold = 5.0

// imgIsColor returns true if image `img` contains color
func imgIsColor(img image.Image) bool {
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rg := int(r) - int(g)
			gb := int(g) - int(b)
			rgb := normalize(abs(rg) + abs(gb))
			if rgb > colorThreshold {
				return true
			}
		}
	}
	return false
}

func imgMarkColor(imgIn image.Image) (image.Image, string) {
	img := image.NewNRGBA(imgIn.Bounds())
	black := color.RGBA{0, 0, 0, 255}
	// white := color.RGBA{255, 255, 255, 255}
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	data := []float64{}
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b, _ := imgIn.At(x, y).RGBA()
			rg := int(r) - int(g)
			gb := int(g) - int(b)
			rgb := normalize(abs(rg) + abs(gb))
			if rgb > colorThreshold {
				img.Set(x, y, black)
				data = append(data, rgb)
			}
		}
	}
	return img, summarizeSeries(data)
}

func normalize(v int) float64 {
	return float64(v) / float64(0xFFFF) * float64(0xFF)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func summarizeSeries(data []float64) string {
	n := len(data)
	total := 0.0
	min := +1e20
	max := -1e20
	for _, x := range data {
		total += x
		if x < min {
			min = x
		}
		if x > max {
			max = x
		}
	}
	mean := total / float64(n)
	return fmt.Sprintf("n=%d min=%.3f mean=%.3f max=%.3f", n, min, mean, max)
}

// readImage reads image file `path` and returns its contents as an Image.
func readImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		common.Log.Error("readImage: Could not open file. path=%#q err=%v", path, err)
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	return img, err
}

// readImage reads image file `path` and returns its contents as an Image.
func writeImage(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		common.Log.Error("writeImage: Could not create file. path=%#q err=%v", path, err)
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}
