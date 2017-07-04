/*
 * Check a directory of PNG files and check if any contain color
 */

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/unidoc/unidoc/common"
)

func main() {
	debug := false // Write debug level info to stdout?
	keep := false  // Keep the rasters used for PDF comparison"
	flag.BoolVar(&debug, "d", false, "Enable debug logging")
	flag.BoolVar(&keep, "k", false, "Keep the rasters used for PDF comparison")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-d -k] <PNG directory> ...\n", os.Args[0])
		os.Exit(1)
	}

	dir := args[0]

	numColor, err := colorDirectoryPages(".png", dir, keep)

	if err != nil {
		fmt.Fprintf(os.Stderr, "colorDirectoryPages failed. err=%v", err)
	}

	fmt.Printf("%d color pages\n", numColor)

}

var (
	gsImagePattern = `doc-(\d+).png$`
	gsImageRegex   = regexp.MustCompile(gsImagePattern)
)

// colorDirectoryPages returns a lists of the page numbers of the image files that match `mask` in
// directories `dir` that have any color pixels.
func colorDirectoryPages(mask, dir string, keep bool) ([]int, error) {
	pattern := filepath.Join(dir, mask)
	files, err := filepath.Glob(pattern)
	if err != nil {
		common.Log.Error("isColorDirectory: Glob failed. pattern=%#q err=%v", pattern, err)
		return nil, err
	}

	colorPages := []int{}
	for _, path := range files {
		matches := gsImageRegex.FindStringSubmatch(path)
		if len(matches) == 0 {
			continue
		}
		pageNum, err := strconv.Atoi(matches[1])
		if err != nil {
			panic(err)
			return colorPages, err
		}
		// common.Log.Error("isColorDirectory:  path=%#q", path)
		isColor, err := isColorImage(path, keep)
		// common.Log.Error("isColorDirectory: isColor=%t path=%#q", isColor, path)
		if err != nil {
			panic(err)
			return colorPages, err
		}
		if isColor {
			colorPages = append(colorPages, pageNum)
			// common.Log.Error("isColorDirectory: colorPages=%d %d", len(colorPages), colorPages)
		}
	}
	return colorPages, nil
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
			rgb := float64(abs(rg)+abs(gb)) / float64(0xFFFF) * float64(0xFF)
			if rgb > colorThreshold {
				return true
			}
		}
	}
	return false
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
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
			if rg < 0 {
				rg = -rg
			}
			if gb < 0 {
				gb = -gb
			}
			rgb := float64(rg+gb) / float64(0xFFFF) * float64(0xFF)
			if rgb > colorThreshold {
				img.Set(x, y, black)
				data = append(data, rgb)
			}
		}
	}
	return img, summarizeSeries(data)
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

// makeUniqueDir creates a new directory inside `baseDir`
func makeUniqueDir(baseDir string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 1000; i++ {
		dir := filepath.Join(baseDir, fmt.Sprintf("dir.%03d", i))
		if _, err := os.Stat(dir); err != nil {
			if os.IsNotExist(err) {
				if err := os.MkdirAll(dir, 0777); err != nil {
					panic(err)
				}
				return dir
			}
		}
		time.Sleep(time.Duration(r.Float64() * float64(time.Second)))
	}
	panic("Cannot create new directory")
}

// removeDir removes directory `dir` and its contents
func removeDir(dir string) error {
	err1 := os.RemoveAll(dir)
	err2 := os.Remove(dir)
	if err1 != nil {
		return err1
	}
	return err2
}

// patternsToPaths returns a list of files matching the patterns in `patternList`
func patternsToPaths(patternList []string) ([]string, error) {
	pathList := []string{}
	for _, pattern := range patternList {
		files, err := filepath.Glob(pattern)
		if err != nil {
			common.Log.Error("patternsToPaths: Glob failed. pattern=%#q err=%v", pattern, err)
			return pathList, err
		}
		for _, path := range files {
			if !regularFile(path) {
				fmt.Fprintf(os.Stderr, "Not a regular file. %#q\n", path)
				continue
			}
			pathList = append(pathList, path)
		}
	}
	return pathList, nil
}

// regularFile returns true if file `path` is a regular file
func regularFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	return fi.Mode().IsRegular()
}

// fileSize returns the size of file `path` in bytes
func fileSize(path string) int64 {
	fi, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	return fi.Size()
}

type statistics struct {
	enabled        bool
	testResultPath string
	imageInfoPath  string
	testResultList []testResult
	testResultMap  map[string]int
}

func (s *statistics) load() error {
	if !s.enabled {
		return nil
	}
	s.testResultList = []testResult{}
	s.testResultMap = map[string]int{}

	testResultList, err := testResultRead(s.testResultPath)
	if err != nil {
		return err
	}
	for _, e := range testResultList {
		s.addTestResult(e, true)
	}

	return nil
}

func (s *statistics) _save() error {
	if !s.enabled {
		return nil
	}
	return testResultWrite(s.testResultPath, s.testResultList)
}

func (s *statistics) addTestResult(e testResult, force bool) {
	if !s.enabled {
		return
	}
	i, ok := s.testResultMap[e.name]
	if !ok {
		s.testResultList = append(s.testResultList, e)
		s.testResultMap[e.name] = len(s.testResultList) - 1
	} else {
		s.testResultList[i] = e
	}
	if force {
		s._save()
	}
}

type testResult struct {
	name     string
	colorIn  bool
	colorOut bool
	numPages int
	duration float64
	xobjImg  int
	xobjForm int
}

func testResultRead(path string) ([]testResult, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []testResult{}, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)

	results := []testResult{}
	for i := 0; ; i++ {
		row, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			common.Log.Error("testResultRead: i=%d err=%v", i, err)
			return results, err
		}
		if i == 0 {
			continue
		}
		e := testResult{
			name:     row[0],
			colorIn:  toBool(row[1]),
			colorOut: toBool(row[2]),
			numPages: toInt(row[3]),
			duration: toFloat(row[4]),
			xobjImg:  toInt(row[5]),
			xobjForm: toInt(row[6]),
		}
		results = append(results, e)
	}
	return results, nil
}

func testResultWrite(path string, results []testResult) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)

	if err := w.Write([]string{"name", "colorIn", "colorOut", "numPages", "duration",
		"imageXobj", "formXobj"}); err != nil {
		return err
	}
	for i, e := range results {
		row := []string{
			e.name,
			fmt.Sprintf("%t", e.colorIn),
			fmt.Sprintf("%t", e.colorOut),
			fmt.Sprintf("%d", e.numPages),
			fmt.Sprintf("%.3f", e.duration),
			fmt.Sprintf("%d", e.xobjImg),
			fmt.Sprintf("%d", e.xobjForm),
		}
		if err := w.Write(row); err != nil {
			common.Log.Error("testResultWrite: Error writing record. i=%d path=%#q err=%v",
				i, path, err)
		}
	}

	w.Flush()
	return w.Error()
}

func toBool(s string) bool {
	return strings.ToLower(strings.TrimSpace(s)) == "true"
}

func toInt(s string) int {
	s = strings.TrimSpace(s)
	x, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return x
}

func toFloat(s string) float64 {
	s = strings.TrimSpace(s)
	x, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return x
}

func changeDir(path, dir string) string {
	_, name := filepath.Split(path)
	return filepath.Join(dir, name)
}
