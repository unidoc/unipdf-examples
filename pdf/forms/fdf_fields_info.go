/*
 * Prints basic FDF field data in terms of key and values.
 *
 * Run as: go run fdf_fields_info.go input1.fdf [input2.fdf ...]
 */

package main

import (
	"errors"
	"fmt"
	"os"

	unicommon "github.com/unidoc/unidoc/common"
	"github.com/unidoc/unidoc/pdf/core"
	"github.com/unidoc/unidoc/pdf/fdf"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Print out basic properties of FDF files\n")
		fmt.Printf("Usage: go run fdf_fields_info.go input.fdf [input2.pdf] ...\n")
		os.Exit(1)
	}

	// Enable debug-level logging.
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

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
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	defer f.Close()

	p, err := fdf.NewParser(f)
	if err != nil {
		return err
	}

	fdfDict, err := p.Root()
	if err != nil {
		return err
	}

	fmt.Printf("FDF: %v\n", fdfDict)

	fields, found := core.GetArray(fdfDict.Get("Fields"))
	if !found {
		return errors.New("Fields missing")
	}

	for i := 0; i < fields.Len(); i++ {
		fieldDict, has := core.GetDict(fields.Get(i))
		if has {
			// Key value field data.
			t, _ := core.GetString(fieldDict.Get("T"))
			v := core.TraceToDirectObject(fieldDict.Get("V"))
			if t != nil && v != nil {
				fmt.Printf("Field T: %v, V: %v\n", t, v.String())
			}
		}
	}

	return nil
}
