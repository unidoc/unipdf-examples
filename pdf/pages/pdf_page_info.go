/*
 * Prints PDF page info: Mediabox size and other parameters.
 * If [page num] is not specified prints out info for all pages.
 *
 * Run as: go run pdf_info.go input.pdf [page num]
 */

package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	unicommon "github.com/unidoc/unidoc/common"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage:  go run pdf_info.go input.pdf [page num]\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]

	pageNum := 0
	if len(os.Args) > 2 {
		num, err := strconv.ParseInt(os.Args[2], 10, 64)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		pageNum = int(num)
	}

	// Enable debug-level logging.
	// unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	fmt.Printf("Input file: %s\n", inputPath)

	m := startMemoryMeasurement()

	err := printPdfPageProperties(inputPath, pageNum)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	m.Stop()
	fmt.Println("======================================")
	fmt.Printf("%s\n", m.Summary())
}

func printPdfPageProperties(inputPath string, pageNum int) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return err
	}

	// Try decrypting with an empty one.
	if isEncrypted {
		auth, err := pdfReader.Decrypt([]byte(""))
		if err != nil {
			return err
		}
		if !auth {
			unicommon.Log.Debug("Encrypted - unable to access - update code to specify pass")
			return nil
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	// If invalid pagenum, print all pages.
	if pageNum <= 0 || pageNum > numPages {
		for i := 0; i < numPages; i++ {
			page, err := pdfReader.GetPage(i + 1)
			if err != nil {
				return err
			}
			fmt.Printf("-- Page %d\n", i+1)
			err = processPage(page)
			if err != nil {
				return err
			}
		}
	} else {
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return err
		}
		fmt.Printf("-- Page %d\n", pageNum)
		err = processPage(page)
		if err != nil {
			return err
		}
	}

	return nil
}

func processPage(page *pdf.PdfPage) error {
	mBox, err := page.GetMediaBox()
	if err != nil {
		return err
	}
	pageWidth := mBox.Urx - mBox.Llx
	pageHeight := mBox.Ury - mBox.Lly

	fmt.Printf(" Page: %+v\n", page)
	if page.Rotate != nil {
		fmt.Printf(" Page rotation: %v\n", *page.Rotate)
	} else {
		fmt.Printf(" Page rotation: 0\n")
	}
	fmt.Printf(" Page mediabox: %+v\n", page.MediaBox)
	fmt.Printf(" Page height: %f\n", pageHeight)
	fmt.Printf(" Page width: %f\n", pageWidth)

	return nil
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
