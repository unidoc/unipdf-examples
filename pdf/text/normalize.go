/*
 * PDF to text: Extract all text for each page of a pdf file.
 *
 * Run as: go run pdf_extract_text.go input.pdf
 */

package main

import (
	"fmt"

	"golang.org/x/text/unicode/norm"
)

func main() {
	// "ff":                          '\ufb00', // ﬀ
	// "ffi":                         '\ufb03', // ﬃ
	// "ffl":                         '\ufb04', // ﬄ
	// "fi":                          '\ufb01', // ﬁ
	diacritics := "őōṓṑỏόớờỡॐϖѽώ⅛⅕₁⒈"
	ligatures := "ﬀﬃﬄﬁ"
	e1 := "e\u0301"
	e2 := "é"
	e3 := "é́́"
	e4 := "e\u0301\u0301\u0301\u0301\u0301\u0301\u0301\u0301\u0301"
	e5 := "e\u0301\u0302\u0303\u0304\u0305\u0306\u0307\u0310\u0311" +
		"\u0321\u0322\u0331\u0332" +
		"\u0341\u0342\u0343" +
		"\u0351\u0352\u0353" +
		"\u0361\u0364\u0363"
	e6 := "e\u0301\u0302\u0332\u0341\u0342\u0343"
	same := e1 == e2
	runes := []rune(ligatures)
	fmt.Printf("ligatures=%q\n", ligatures)
	fmt.Printf("runes=%q\n", runes)
	fmt.Printf("e1=%q e2=%q same=%t\n", e1, e2, same)
	fmt.Printf("e3=%q=%#q=%+q\n", e3, e3, e3)
	fmt.Printf("\ne4=%q=%#q=%+q\n", e4, e4, e4)
	fmt.Printf("\n\n\n\ne5=%q=%#q\n\n\n", e5, e5)

	fmt.Printf("\n\n\n%q\n\n\n", "X\u0301\u0302\u0303\u0304\u0305\u0306\u0307\u0310\u0311"+
		"\u0321\u0322\u0331\u0332"+
		"\u0341\u0342\u0343"+
		"\u0351\u0352\u0353"+
		"\u0361\u0364\u0363")

	show("diacritics", diacritics)
	show("ligatures", ligatures)
	show("e1", e1)
	show("e2", e2)
	show("e6", e6)

}

func show(name, text string) {
	nc := norm.NFC.Bytes([]byte(text))
	nd := norm.NFD.Bytes([]byte(text))
	sd := string(nd)
	fmt.Printf("\n*** %s=%q=%+q\n", name, text, text)
	fmt.Printf("\t\tcomposed(%T)=%q=%+q\n", nc, nc, nc)
	fmt.Printf("\t\tseparated(%T)=%q=%+q\n", nd, nd, nd)
	fmt.Printf("\t\tstring(%T)=%q=%+q\n", sd, sd, sd)
	fmt.Printf("\n")
}
