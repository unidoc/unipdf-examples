/*
 * Flatten form data in a PDF file, move to content stream, so cannot be edited.
 * Note: Currently only works for field data with an apperance stream.
 * TODO: Add support for generating appearance streams (default for fields missing AP or force generation on all).
 *
 * Run as: go run pdf_form_flatten.go <input.pdf> <flattened.pdf>
 */

package main

import (
	"fmt"
	"os"

	"path/filepath"
	"sort"

	unicommon "github.com/unidoc/unidoc/common"
	"github.com/unidoc/unidoc/pdf/model"
)

func main() {
	// When debugging, enable debug-level logging via console:
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	if len(os.Args) < 3 {
		//fmt.Printf("Usage: go run pdf_form_flatten.go <input.pdf> <output.pdf>\n")
		fmt.Printf("Usage: go run pdf_form_flatten.go <outputdir> <input1.pdf> ...\n")

		os.Exit(1)
	}

	outputDir := os.Args[1]

	fails := map[string]string{}
	failKeys := []string{}
	processed := 0

	for i := 2; i < len(os.Args); i++ {
		inputPath := os.Args[i]
		name := filepath.Base(inputPath)
		outputPath := filepath.Join(outputDir, fmt.Sprintf("flattened_%s", name))
		err := flattenPdf(inputPath, outputPath)
		if err != nil {
			fmt.Printf("%s - Error: %v\n", inputPath, err)
			fails[inputPath] = err.Error()
			failKeys = append(failKeys, inputPath)
		}
		processed++
	}

	fmt.Printf("Total %d processed / %d failures\n", processed, len(failKeys))
	sort.Strings(failKeys)
	for _, k := range failKeys {
		fmt.Printf("%s: %v\n", k, fails[k])
	}
}

// flattenPdf flattens annotations and forms moving the appearance stream to the page contents so cannot be
// modified.
// TODO: Cross-reference between annotations and form fields, should give an error if a form field does not
// have an appearance stream (widget annotation) present for drawing it.
func flattenPdf(inputPath, outputPath string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return err
	}

	err = pdfReader.FlattenFields(false)
	//err = flatten(pdfReader, true, false)
	if err != nil {
		return err
	}

	pdfWriter := model.NewPdfWriter()
	pdfWriter.SetForms(nil)

	for _, p := range pdfReader.PageList {
		err := pdfWriter.AddPage(p)
		if err != nil {
			return err
		}
	}

	fout, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer fout.Close()

	err = pdfWriter.Write(fout)
	return err
}

/*
// flatten flattens the PDF loaded in `pdf`. If `doannots` is true, all annotations will be flattened.
// If `dofields` is true will lookup all widget annotations corresponding to form fields and flatten them.
//
// References to flattened annotations will be removed from Page Annots array. For fields the AcroForm entry
// will be emptied.
//
// TODO: Plan is to add Flatten(dofields, doannots bool) error to PdfReader.
func flatten(pdf *model.PdfReader, dofields, doannots bool) error {

}

// Get the active XObject Form for an appearance dictionary.
// Default gets the N entry, and if it is a dictionary, picks the entry referred to by AS.
// If returned XObject Form is nil (and no errors) it indicates that the annotation has no appearance.
func getAnnotationActiveAppearance(annot *model.PdfAnnotation) (*model.XObjectForm, *model.PdfRectangle, error) {

	fmt.Printf("----\n")
	fmt.Printf("annot: %#v\n", annot)
	fmt.Printf("context: %#v\n", annot.GetContext())
	fmt.Printf("obj: %v\n", annot.GetContainingPdfObject())

	// Appearance dictionary entries (Table 168 p. 397).
	apDict, has := core.GetDict(annot.AP)
	if !has {
		return nil, nil, errors.New("field missing AP dictionary")
	}

	// Get the Rect specifying the display rectangle.
	rectArr, has := core.GetArray(annot.Rect)
	if !has || rectArr.Len() != 4 {
		return nil, nil, errors.New("rect invalid")
	}
	rect, err := model.NewPdfRectangle(*rectArr)
	if err != nil {
		return nil, nil, err
	}

	nobj := core.TraceToDirectObject(apDict.Get("N"))
	switch t := nobj.(type) {
	case *core.PdfObjectStream:
		stream := t
		xform, err := model.NewXObjectFormFromStream(stream)
		return xform, rect, err
	case *core.PdfObjectDictionary:
		// An annotation representing multiple fields may have many appearances.
		// As an example checkbox may have two appearance states On and Off.
		// Its appearance dictionary would contain /N << /On Ref /Off Ref >>, the choice is
		// determines by the AS entry in the annotation dictionary.
		nDict := t

		fmt.Printf("AS: %v\n", annot.AS)
		state, has := core.GetName(annot.AS)
		if !has {
			//state = core.MakeName("Off") // Default.
			// No appearance (nil).
			return nil, nil, nil
		}

		// If only 1 element, use it.
		// XXX/FIXME: Not always good behavior? What about V?
		/
			if len(nDict.Keys()) == 1 {
				state = &nDict.Keys()[0]
			}
			fmt.Printf("State: '%s'\n", state.String())
		/

		if nDict.Get(*state) == nil {
			fmt.Printf("Error: AS state not specified in AP dict\n")
			return nil, nil, nil
		}

		stream, has := core.GetStream(nDict.Get(*state))
		if !has {
			fmt.Printf("Unable to access stream for %v\n", state)
			return nil, nil, errors.New("stream missing")
		}
		xform, err := model.NewXObjectFormFromStream(stream)
		return xform, rect, err
	}

	fmt.Printf("Invalid type for N: %T\n", nobj)
	return nil, nil, errors.New("type check error")
}

*/
