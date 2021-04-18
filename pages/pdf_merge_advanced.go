/*
 * Merge PDF files, including form field data (AcroForms).
 * For a more basic merging of PDF page contents, see pdf_merge.go.
 *
 * Run as: go run pdf_merge_advanced.go output.pdf input1.pdf input2.pdf input3.pdf ...
 */

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Requires at least 3 arguments: output_path and 2 input paths\n")
		fmt.Printf("Usage: go run pdf_merge.go output.pdf input1.pdf input2.pdf input3.pdf ...\n")
		os.Exit(0)
	}

	outputPath := ""
	inputPaths := []string{}

	// Sanity check the input arguments.
	for i, arg := range os.Args {
		if i == 0 {
			continue
		} else if i == 1 {
			outputPath = arg
			continue
		}

		inputPaths = append(inputPaths, arg)
	}

	err := mergePdf(inputPaths, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

// Merge form resources.
// TODO: Add handling for cases where same resource name is used with different values.  In that case, need to rename
// the resource and change all references to that value with the new value.
func mergeResources(r, r2 *model.PdfPageResources) (*model.PdfPageResources, error) {
	// Merge XObject resources.
	if r.XObject == nil {
		r.XObject = r2.XObject
	} else {
		xobjs, _ := core.GetDict(r.XObject)
		if r2.XObject != nil {
			xobjs2, _ := core.GetDict(r2.XObject)
			for _, key := range xobjs2.Keys() {
				val := xobjs2.Get(key)
				// Add XObjects from r2.  Overwrite if existing...
				// TODO: Handle overwrites properly.
				xobjs.Set(key, val)
			}
		}
	}

	// Merge Colorspace resources.
	colorspaces, err := r.GetColorspaces()
	if err != nil {
		return nil, err
	}
	colorspaces2, err := r2.GetColorspaces()
	if err != nil {
		return nil, err
	}
	if colorspaces == nil {
		r.SetColorSpace(colorspaces2)
	} else {
		if colorspaces2 != nil {
			for key, val := range colorspaces2.Colorspaces {
				// Add the r2 colorspaces to r. Overwrite if duplicate.  Ensure only present once in Names.
				if _, has := colorspaces.Colorspaces[key]; !has {
					colorspaces.Names = append(colorspaces.Names, key)
				}
				r.SetColorspaceByName(core.PdfObjectName(key), val)
			}
		}
	}

	// Merge ExtGState resources.
	if r.ExtGState == nil {
		r.ExtGState = r2.ExtGState
	} else {
		extgstates, _ := core.GetDict(r.ExtGState)

		if r2.ExtGState != nil {
			extgstates2, _ := core.GetDict(r2.ExtGState)
			for _, key := range extgstates2.Keys() {
				// TODO: Handle overwrites properly.
				val := extgstates2.Get(key)
				extgstates.Set(key, val)
			}
		}
	}

	if r.Shading == nil {
		r.Shading = r2.Shading
	} else {
		shadings, _ := core.GetDict(r.Shading)
		if r2.Shading != nil {
			shadings2, _ := core.GetDict(r2.Shading)
			for _, key := range shadings2.Keys() {
				val := shadings2.Get(key)
				shadings.Set(key, val)
			}
		}
	}

	if r.Pattern == nil {
		r.Pattern = r2.Pattern
	} else {
		shadings, _ := core.GetDict(r.Pattern)
		if r2.Pattern != nil {
			patterns2, _ := core.GetDict(r2.Pattern)
			for _, key := range patterns2.Keys() {
				val := patterns2.Get(key)
				shadings.Set(key, val)
			}
		}
	}

	if r.Font == nil {
		r.Font = r2.Font
	} else {
		fonts, _ := core.GetDict(r.Font)
		if r2.Font != nil {
			fonts2, _ := core.GetDict(r2.Font)
			for _, key := range fonts2.Keys() {
				val := fonts2.Get(key)
				fonts.Set(key, val)
			}
		}
	}

	if r.ProcSet == nil {
		r.ProcSet = r2.ProcSet
	} else {
		procsets, _ := core.GetDict(r.ProcSet)
		if r2.ProcSet != nil {
			procsets2, _ := core.GetDict(r2.ProcSet)
			for _, key := range procsets2.Keys() {
				val := procsets2.Get(key)
				procsets.Set(key, val)
			}
		}
	}

	if r.Properties == nil {
		r.Properties = r2.Properties
	} else {
		props, _ := core.GetDict(r.Properties)
		if r2.Properties != nil {
			props2, _ := core.GetDict(r2.Properties)
			for _, key := range props2.Keys() {
				val := props2.Get(key)
				props.Set(key, val)
			}
		}
	}

	return r, nil
}

// Merge two interactive forms.
func mergeForms(form, form2 *model.PdfAcroForm, docNum int) (*model.PdfAcroForm, error) {
	// Use whatever value comes first..
	// TODO: Consider adding a more intelligent, preferential handling based on actual values.  If needed.

	if form.NeedAppearances == nil {
		form.NeedAppearances = form2.NeedAppearances
	}

	if form.SigFlags == nil {
		form.SigFlags = form2.SigFlags
	}

	if form.CO == nil {
		form.CO = form2.CO
	}

	if form.DR == nil {
		form.DR = form2.DR
	} else if form2.DR != nil {
		dr, err := mergeResources(form.DR, form2.DR)
		if err != nil {
			return nil, err
		}
		form.DR = dr
	}

	if form.DA == nil {
		form.DA = form2.DA
	}

	if form.Q == nil {
		form.Q = form2.Q
	}

	if form.XFA == nil {
		form.XFA = form2.XFA
	} else {
		if form2.XFA != nil {
			// TODO: Handle merging XFA.
			common.Log.Debug("TODO: Handle XFA merging - Currently just using first one that is encountered")
		}
	}

	// Fields.
	if form.Fields == nil {
		form.Fields = form2.Fields
	} else {
		// Make a top-level field for the doc (non-terminal field).
		docfield := model.NewPdfField()
		docfield.T = core.MakeString(fmt.Sprintf("doc%d", docNum))
		docfield.Kids = []*model.PdfField{}
		if form2.Fields != nil {
			for _, subfield := range *form2.Fields {
				subfield.Parent = docfield // Update parent.
				docfield.Kids = append(docfield.Kids, subfield)
			}
		}
		*form.Fields = append(*form.Fields, docfield)
	}

	return form, nil
}

func mergePdf(inputPaths []string, outputPath string) error {
	pdfWriter := model.NewPdfWriter()

	var forms *model.PdfAcroForm

	for docIdx, inputPath := range inputPaths {
		f, err := os.Open(inputPath)
		if err != nil {
			return err
		}

		defer f.Close()

		pdfReader, err := model.NewPdfReader(f)
		if err != nil {
			return err
		}

		isEncrypted, err := pdfReader.IsEncrypted()
		if err != nil {
			return err
		}

		if isEncrypted {
			_, err = pdfReader.Decrypt([]byte(""))
			if err != nil {
				return err
			}
		}

		numPages, err := pdfReader.GetNumPages()
		if err != nil {
			return err
		}

		for i := 0; i < numPages; i++ {
			pageNum := i + 1

			page, err := pdfReader.GetPage(pageNum)
			if err != nil {
				return err
			}

			err = pdfWriter.AddPage(page)
			if err != nil {
				return err
			}
		}

		// Handle forms.
		if pdfReader.AcroForm != nil {
			if forms == nil {
				forms = pdfReader.AcroForm
			} else {
				forms, err = mergeForms(forms, pdfReader.AcroForm, docIdx+1)
				if err != nil {
					return err
				}
			}
		}
	}

	fWrite, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer fWrite.Close()

	// Set the merged forms object.
	if forms != nil {
		pdfWriter.SetForms(forms)
	}

	err = pdfWriter.Write(fWrite)
	if err != nil {
		return err
	}

	return nil
}
