/*
 * Prints basic FDF field data in terms of key and values.
 *
 * Run as: go run fdf_fields_info.go input1.fdf [input2.fdf ...]
 */

package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/fdf"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Print out basic properties of FDF files\n")
		fmt.Printf("Usage: go run fdf_fields_info.go input.fdf [input2.pdf] ...\n")
		os.Exit(1)
	}

	// Enable debug-level logging.
	common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))

	for _, inputPath := range os.Args[1:len(os.Args)] {
		fmt.Printf("Input file: %s\n", inputPath)

		err := printFdfInfo(inputPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}
}

func printFdfInfo(inputPath string) error {
	fdf, err := fdf.LoadFromPath(inputPath)
	if err != nil {
		return err
	}

	fieldMap, err := fdf.FieldDictionaries()
	if err != nil {
		return err
	}

	// Sort field names alphabetically.
	keys := []string{}
	for key, _ := range fieldMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		fieldDict := fieldMap[key]
		// Key value field data.
		t, _ := core.GetString(fieldDict.Get("T"))
		v := core.TraceToDirectObject(fieldDict.Get("V"))
		if t != nil && v != nil {
			fmt.Printf("Field T: %v, V: %v\n", t, v.String())
		}

	}

	return nil
}
