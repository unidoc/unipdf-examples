// Find all the duplicate files in a test corpus.
//
// go run corpus_duplicates.go testdata/*.pdf testdata/**/*.pdf

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/bmatcuk/doublestar"
)

const usage = "Usage: corpus_duplicates  <file1> <file2> ..."

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}

	pathList, err := patternsToPaths(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "patternsToPaths failed. args=%#q err=%v\n", args, err)
		os.Exit(1)
	}

	fmt.Printf("%d files\n", len(pathList))

	hashPaths := map[string][]string{}
	for _, path := range pathList {
		if !regularFile(path) {
			continue
		}
		hash, err := fileHash(path)
		if err != nil {
			os.Exit(1)
		}
		hashPaths[hash] = append(hashPaths[hash], path)
	}
	duplicates := []string{}
	for hash, paths := range hashPaths {
		if len(paths) > 1 {
			sort.Strings(paths)
			duplicates = append(duplicates, hash)
		}
	}
	sort.Strings(duplicates)
	fmt.Printf("%d duplicates\n", len(duplicates))
	for i, hash := range duplicates {
		paths := hashPaths[hash]
		fmt.Printf("%3d: %2d: %+v\n", i, len(paths), paths)
	}
}

// patternsToPaths returns a list of files matching the patterns in `patternList`
func patternsToPaths(patternList []string) ([]string, error) {
	pathList := []string{}
	for _, pattern := range patternList {
		files, err := doublestar.Glob(pattern)
		if err != nil {
			fmt.Fprintf(os.Stderr, "patternsToPaths: Glob failed. pattern=%#q err=%v\n", pattern, err)
			return pathList, err
		}
		for _, path := range files {
			if !regularFile(path) {
				fmt.Fprintf(os.Stderr, "Not a regular file: %q\n", path)
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
		fmt.Fprintf(os.Stderr, "Could not stat %q. err=%v", path, err)
		return false
	}
	return fi.Mode().IsRegular()
}

// fileHash returns a string with the sha256 hash of the file `path`.
func fileHash(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read %q. err=%v", path, err)
		return "", err
	}
	hasher := sha256.New()
	hasher.Write(b)
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
