/*
 * pdf_search_replace_advanced.go - More advanced example of search and replace with UniPDF.
 * Replaces <text> with <replace text> in the output PDF. Takes font encoding
 * into account when finding occurrences of the search term. In addition, it
 * attempts to identify occurences of the search term which are contained within
 * multiple PDF text processing operators (Tj, TJ).
 *
 * Syntax: go run pdf_search_replace_advanced.go <input.pdf> <output.pdf> <text> <replace text>
 */

package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/contentstream"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/model/optimize"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}
}

type textChunk struct {
	font   *model.PdfFont
	strObj *core.PdfObjectString
	val    string
	idx    int
}

func (tc *textChunk) encode() {
	var encoded string
	if font := tc.font; font != nil {
		encodedBytes, numMisses := font.StringToCharcodeBytes(tc.val)
		if numMisses != 0 {
			common.Log.Debug("WARN: some runes could not be encoded.\n\t%s -> %v")
		}
		encoded = string(encodedBytes)
	}

	*tc.strObj = *core.MakeString(encoded)
}

type textChunks struct {
	text   string
	chunks []*textChunk
}

func (tc *textChunks) replace(search, replacement string) {
	text := tc.text
	chunks := tc.chunks

	// Steps:
	// 1. Search for the first index of the search term in the text.
	// 2. Use the found index to match the text chunk which contains
	//    (or partly contains) the search term.
	// 3. Replace the search term in the found text chunk. The search term
	//    will not always start at the beginning of the text chunk. Also,
	//    the search term could be split in multiple text chunks. If that's
	//    the case, replace the portion of the search term in the found
	//    chunk and continue removing characters from the following chunks
	//    until the search term has been completely erased.
	// 4. Offset the text chunks slice to the last processed text chunk from
	//    the previous step, if the text chunk was not completely erased, or
	//    to the next one otherwise. This is necessary so that the visited
	//    text chunks are skipped when searching for the next occurrence of the
	//    search term.
	// 5. Discard the part of the text up to (and including) the index found
	//    in step one.
	// 6. Move to step 1 in order to search for the search term in the remaining
	//    text.
	var chunkOffset int
	matchIdx := strings.Index(text, search)
	for currMatchIdx := matchIdx; matchIdx != -1; {
		for i, chunk := range chunks[chunkOffset:] {
			idx, lenChunk := chunk.idx, len(chunk.val)
			if currMatchIdx < idx || currMatchIdx > idx+lenChunk-1 {
				continue
			}
			chunkOffset += i + 1

			start := currMatchIdx - idx
			remaining := len(search) - (lenChunk - start)

			replaceVal := chunk.val[:start] + replacement
			if remaining < 0 {
				replaceVal += chunk.val[lenChunk+remaining:]
				chunkOffset--
			}

			chunk.val = replaceVal
			chunk.encode()

			for j := chunkOffset; remaining > 0; j++ {
				c := chunks[j]
				l := len(c.val)

				if l > remaining {
					c.val = c.val[remaining:]
				} else {
					c.val = ""
					chunkOffset++
				}

				c.encode()
				remaining -= l
			}

			break
		}

		text = text[matchIdx+1:]
		matchIdx = strings.Index(text, search)
		currMatchIdx += matchIdx + 1
	}

	tc.text = strings.Replace(tc.text, search, replacement, -1)
}

func main() {
	// Uncomment to enable debugging.
	//common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))

	if len(os.Args) < 5 {
		fmt.Printf("Usage: go run pdf_search_replace_advanced.go <input.pdf> <output.pdf> <text> <replace text>\n")
		os.Exit(0)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]
	searchText := os.Args[3]
	replaceText := os.Args[4]

	err := searchReplace(inputPath, outputPath, searchText, replaceText)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully created %s\n", outputPath)
}

func searchReplace(inputPath, outputPath, searchText, replaceText string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()
	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return err
	}

	pdfWriter := model.NewPdfWriter()

	encrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return err
	}
	if encrypted {
		ok, err := pdfReader.Decrypt([]byte(""))
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("Encrypted")
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	for n := 1; n <= numPages; n++ {
		page, err := pdfReader.GetPage(n)
		if err != nil {
			return err
		}

		err = searchReplacePageText(page, searchText, replaceText)
		if err != nil {
			return err
		}

		err = pdfWriter.AddPage(page)
		if err != nil {
			return err
		}
	}

	fw, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer fw.Close()

	opt := optimize.Options{
		CombineDuplicateStreams:         true,
		CombineIdenticalIndirectObjects: true,
		UseObjectStreams:                true,
		CompressStreams:                 true,
	}
	pdfWriter.SetOptimizer(optimize.New(opt))

	return pdfWriter.Write(fw)
}

func searchReplacePageText(page *model.PdfPage, searchText, replaceText string) error {
	contents, err := page.GetAllContentStreams()
	if err != nil {
		return err
	}

	ops, err := contentstream.NewContentStreamParser(contents).Parse()
	if err != nil {
		return err
	}

	// Generate text chunks.
	var currFont *model.PdfFont
	tc := textChunks{}

	textProcFunc := func(objptr *core.PdfObject) {
		strObj, ok := core.GetString(*objptr)
		if !ok {
			common.Log.Debug("Invalid parameter, skipping")
			return
		}

		str := strObj.String()
		if currFont != nil {
			decoded, _, numMisses := currFont.CharcodeBytesToUnicode(strObj.Bytes())
			if numMisses != 0 {
				common.Log.Debug("WARN: some charcodes could not be decoded.\n\t%v -> %s", strObj.Bytes(), decoded)
			}
			str = decoded
		}

		tc.chunks = append(tc.chunks, &textChunk{
			font:   currFont,
			strObj: strObj,
			val:    str,
			idx:    len(tc.text),
		})
		tc.text += str
	}

	processor := contentstream.NewContentStreamProcessor(*ops)
	processor.AddHandler(contentstream.HandlerConditionEnumAllOperands, "",
		func(op *contentstream.ContentStreamOperation, gs contentstream.GraphicsState, resources *model.PdfPageResources) error {
			switch op.Operand {
			case `Tj`, `'`:
				if len(op.Params) != 1 {
					common.Log.Debug("Invalid: Tj/' with invalid set of parameters - skip")
					return nil
				}
				textProcFunc(&op.Params[0])
			case `''`:
				if len(op.Params) != 3 {
					common.Log.Debug("Invalid: '' with invalid set of parameters - skip")
					return nil
				}
				textProcFunc(&op.Params[3])
			case `TJ`:
				if len(op.Params) != 1 {
					common.Log.Debug("Invalid: TJ with invalid set of parameters - skip")
					return nil
				}
				arr, _ := core.GetArray(op.Params[0])
				for i := range arr.Elements() {
					obj := arr.Get(i)
					textProcFunc(&obj)
					arr.Set(i, obj)
				}
			case "Tf":
				if len(op.Params) != 2 {
					common.Log.Debug("Invalid: Tf with invalid set of parameters - skip")
					return nil
				}

				fname, ok := core.GetName(op.Params[0])
				if !ok || fname == nil {
					common.Log.Debug("ERROR: could not get font name")
					return nil
				}

				fObj, has := resources.GetFontByName(*fname)
				if !has {
					common.Log.Debug("ERROR: font %s not found", fname.String())
					return nil
				}

				pdfFont, err := model.NewPdfFontFromPdfObject(fObj)
				if err != nil {
					common.Log.Debug("ERROR: loading font")
					return nil
				}
				currFont = pdfFont
			}

			return nil
		})

	if err = processor.Process(page.Resources); err != nil {
		return err
	}

	tc.replace(searchText, replaceText)
	return page.SetContentStreams([]string{ops.String()}, core.NewFlateEncoder())
}
