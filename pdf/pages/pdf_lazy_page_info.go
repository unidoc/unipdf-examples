/*
 * Prints PDF page info: Mediabox size and other parameters.
 * If [page num] is not specified prints out info for all pages.
 *
 * Run as: go run pdf_info.go input.pdf [page num]
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	common "github.com/unidoc/unidoc/common"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

const usage = "Usage:  go run pdf_lazy_info.go input.pdf [page num]\n"

func main() {
	var debug, trace, lazy bool
	flag.BoolVar(&lazy, "z", false, "Use lazy loading.")
	flag.BoolVar(&debug, "d", false, "Print debugging information.")
	flag.BoolVar(&trace, "e", false, "Print detailed debugging information.")
	makeUsage(usage)

	flag.Parse()
	args := flag.Args()

	if trace {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelTrace))
	} else if debug {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
	} else {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelInfo))
	}

	inputPath := args[0]

	num, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	pageNum := int(num)

	fmt.Printf("Input file: %s\n", inputPath)

	m := startMemoryMeasurement()
	page, mBox, err := pdfPageMediaBox(lazy, inputPath, pageNum)
	m.Stop()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("===================------------=======================")
	describPage(page, mBox)
	fmt.Println("======================================================")
	fmt.Printf("Lazy: %t\n", lazy)
	fmt.Printf("%s\n", m.Summary())
}

func pdfPageMediaBox(lazy bool, inputPath string, pageNum int) (*pdf.PdfPage, *pdf.PdfRectangle, error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	var pdfReader *pdf.PdfReader
	if lazy {
		pdfReader, err = pdf.NewPdfReaderLazy(f)
	} else {
		pdfReader, err = pdf.NewPdfReader(f)
	}
	if err != nil {
		return nil, nil, err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return nil, nil, err
	}

	// Try decrypting with an empty one.
	if isEncrypted {
		auth, err := pdfReader.Decrypt([]byte(""))
		if err != nil {
			return nil, nil, err
		}
		if !auth {
			common.Log.Debug("Encrypted - unable to access - update code to specify pass")
			return nil, nil, nil
		}
	}

	page, err := pdfReader.GetPage(pageNum)
	if err != nil {
		return nil, nil, err
	}

	mBox, err := page.GetMediaBox()
	return page, mBox, err
}

func describPage(page *pdf.PdfPage, mBox *pdf.PdfRectangle) {

	pageWidth := mBox.Urx - mBox.Llx
	pageHeight := mBox.Ury - mBox.Lly

	// fmt.Printf(" Page: %+v\n", page)
	if page.Rotate != nil {
		fmt.Printf(" Page rotation: %v\n", *page.Rotate)
	} else {
		fmt.Printf(" Page rotation: 0\n")
	}
	fmt.Printf(" Page mediabox: %+v\n", page.MediaBox)
	fmt.Printf(" Page height: %f\n", pageHeight)
	fmt.Printf(" Page width: %f\n", pageWidth)
}

type memoryMeasure struct {
	start     runtime.MemStats
	startTime time.Time
	end       runtime.MemStats
	duration  time.Duration
}

func startMemoryMeasurement() memoryMeasure {
	var m memoryMeasure

	runtime.ReadMemStats(&m.start)
	m.startTime = time.Now()
	return m
}

// Stops finishes the measurement.
func (m *memoryMeasure) Stop() {
	runtime.ReadMemStats(&m.end)
	m.duration = time.Since(m.startTime)
}

func (m memoryMeasure) Summary() string {
	alloc := float64(m.end.Alloc) - float64(m.start.Alloc)
	mallocs := int64(m.end.Mallocs) - int64(m.start.Mallocs)
	frees := int64(m.end.Frees) - int64(m.start.Frees)

	var b strings.Builder
	b.WriteString(fmt.Sprintf("Duration: %.2f seconds\n", m.duration.Seconds()))
	b.WriteString(fmt.Sprintf("Alloc: %.2f MB\n", alloc/1024.0/1024.0))
	b.WriteString(fmt.Sprintf("Mallocs: %d\n", mallocs))
	b.WriteString(fmt.Sprintf("Frees: %d\n", frees))
	return b.String()
}

// makeUsage updates flag.Usage to include usage message `msg`.
func makeUsage(msg string) {
	usage := flag.Usage
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, msg)
		usage()
	}
}
