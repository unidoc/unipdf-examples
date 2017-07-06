/*
 * Check a directory of PNG files and check if any contain color
 */

package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	debug := false // Write debug level info to stdout?
	flag.BoolVar(&debug, "d", false, "Enable debug logging")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-d -k] <PDF file>\n", os.Args[0])
		os.Exit(1)
	}

	pdf := args[0]
	dir := fmt.Sprintf("%s.rasters", pdf)

	if err := os.MkdirAll(dir, 0777); err != nil {
		fmt.Fprintf(os.Stderr, "MkdirAll failed. dir=%q err=%v\n", dir, err)
		os.Exit(1)
	}

	err := runGhostscript(pdf, dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runGhostscript failed. err=%v\n", err)
		os.Exit(1)
	}

}

var (
	gsImageFormat = "doc-%03d.png"
)

// runGhostscript runs Ghostscript on file `pdf` to create file one png file per page in directory
// `outputDir`
func runGhostscript(pdf, outputDir string) error {
	fmt.Printf("runGhostscript: pdf=%#q outputDir=%#q\n", pdf, outputDir)
	outputPath := filepath.Join(outputDir, gsImageFormat)
	output := fmt.Sprintf("-sOutputFile=%s", outputPath)
	cmd := exec.Command(
		ghostscriptName(),
		"-dSAFER",
		"-dBATCH",
		"-dNOPAUSE",
		"-r150",
		"-sDEVICE=png16m",
		"-dTextAlphaBits=1",
		"-dGraphicsAlphaBits=1",
		output,
		pdf)
	fmt.Printf("runGhostscript: cmd=%#q\n", cmd.Args)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "runGhostscript: Could not process pdf=%q err=%v\nstdout=\n%s\nstderr=\n%s\n",
			pdf, err, stdout, stderr)
	}
	return err
}

// ghostscriptName returns the name of the Ghostscript binary on this OS
func ghostscriptName() string {
	if runtime.GOOS == "windows" {
		return "gswin64c.exe"
	}
	return "gs"
}
