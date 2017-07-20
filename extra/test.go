package main

import (
	"encoding/json"
	"fmt"
)

/*
	text = I'm a text.
	number = 1234
	floats = [1.1 2.2 3.3]
	innermap = map[bar:2 foo:1]
	innermap.foo = 1
	innermap.bar = 2
	The whole map: map[floats:[1.1 2.2 3.3] innermap:map[foo:1 bar:2] text:I'm a text. number:1234]
*/

func main() {
	s := `[{"text":"I'm a text.",
	       "number":1234,"floats":[1.1, 2.2, 3.3],
		   "innermap":{"foo":1,"bar":2}}]`

	var dataList []map[string]interface{}
	err := json.Unmarshal([]byte(s), &dataList)
	if err != nil {
		panic(err)
	}

	data := dataList[0]
	fmt.Println("text =", data["text"])
	fmt.Println("number =", data["number"])
	fmt.Println("floats =", data["floats"])
	fmt.Println("innermap =", data["innermap"])

	innermap, ok := data["innermap"].(map[string]interface{})
	if !ok {
		panic("inner map is not a map!")
	}
	fmt.Println("innermap.foo =", innermap["foo"])
	fmt.Println("innermap.bar =", innermap["bar"])

	fmt.Println("The whole map:", data)
}
