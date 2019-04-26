/*
 * PDF to text: Extract all text for each page of a pdf file.
 *
 * Run as
 go run extract_objects.go  ~/testdata/4e.pdf /Users/pcadmin/go-work/src/github.com/unidoc/unidoc/contrib/testdata/font/helminths.txt 19

 ~/testdata/pair-over-C.pdf
*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Syntax: go run extract_objects.go input.pdf output.pdf")
		os.Exit(1)
	}
	filename := os.Args[1]
	outname := os.Args[2]
	if outname == filename {
		panic("!!!")
	}
	keepNums := []int{}
	for _, a := range os.Args[3:] {
		n, err := strconv.Atoi(a)
		if err != nil {
			panic(err)
		}
		keepNums = append(keepNums, n)
	}
	fmt.Printf("keepNums=%+v\n", keepNums)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	fmt.Printf("data=%d\n", len(data))

	matchGroups := reObj.FindAllStringSubmatch(string(data), -1)
	fmt.Printf("matchGroups=%d\n", len(matchGroups))
	objNums := []int{}
	numObj := map[int]string{}
	for _, groups := range matchGroups {
		n, err := strconv.Atoi(string(groups[1]))
		if err != nil {
			panic(err)
		}
		objNums = append(objNums, n)
		numObj[n] = groups[0]
	}
	sort.Ints(objNums)
	fmt.Printf("objNums=%+v\n", objNums)

	replace := func(text string) string {
		groups := reLengthRef.FindStringSubmatch(text)
		n, err := strconv.Atoi(groups[1])
		if err != nil {
			panic(err)
		}
		fmt.Printf("numObj[%d]=%q\n", n, numObj[n])
		groups2 := reObj.FindStringSubmatch(numObj[n])
		out := fmt.Sprintf("/Length %s", groups2[2])
		fmt.Printf("groups2=%+v\n", groups2)
		fmt.Printf("%q->%q\n", text, out)
		return out
	}

	keepObjs := []string{}
	for _, n := range keepNums {
		obj := numObj[n]
		obj = reToUnicode.ReplaceAllString(obj, "\t")
		obj2 := reLength.ReplaceAllStringFunc(obj, replace)
		if obj2 != obj && len(obj) <= 200 {
			fmt.Printf("%q->%q\n", obj, obj2)
		}
		keepObjs = append(keepObjs, obj2)
	}
	keepObjs = append(keepObjs, " ")

	keepData := strings.Join(keepObjs, "\n")

	err = ioutil.WriteFile(outname, []byte(keepData), 0777)
	if err != nil {
		panic(err)
	}
}

var reObj = regexp.MustCompile(`(?ms)(\d+)\s+0\s+obj\s+(.+?)\s+endobj`)
var reLength = regexp.MustCompile(`(?ms)/Length\s+(\d+\s+0\s+R)`)
var reLengthRef = regexp.MustCompile(`(?ms)(\d+)\s+0\s+R`)
var reToUnicode = regexp.MustCompile(`(?ms)(/ToUnicode\s+\d+\s+0\s+R)`)

// var reObj = regexp.MustCompile(`(?ms)(\d+)\s+0\s+obj`)
