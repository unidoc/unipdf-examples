package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	// test0()
	// jobs := decodeJson(text)
	// fmt.Printf("jobs=%+v\n", jobs)

	url := "https://us-central1-training-173612.cloudfunctions.net/printableJobs?printerId=0xdec3e6b8"
	// res, err := http.Get(url)
	// if err != nil {
	// 	panic(err.Error())
	// }

	// var jobs []Job
	// body := res.Body
	// err = json.NewDecoder(body).Decode(&jobs)
	// if err != nil {
	// 	panic(err)
	// }
	jobs := getJobs(url)
	fmt.Printf("jobs=%+v\n", jobs)
}

func test0() {
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

type Job struct {
	Id          int    `json:"id`
	Duplex      bool   `json:"duplex`
	Color       bool   `json:"color`
	Paper       string `json:"paper`
	Orientation string `json:"orientation`
	Name        string `json:"name`
}

func getJobs(url string) []Job {
	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}

	var jobs []Job
	body := res.Body
	err = json.NewDecoder(body).Decode(&jobs)
	if err != nil {
		panic(err)
	}
	return jobs
}

const text = `[{"id":1,
                "duplex":true,
                "color":true,
                "copies":1,
                "paper":"a4",
                "orientation":"portrait",
                "name":"Mein Kempf",
                "user":"Adolf.h@gmail.com"}]`

func decodeJson(text string) []Job {
	var jobs []Job
	// b := bytes.NewBufferString(text)
	b := strings.NewReader(text)
	err := json.NewDecoder(b).Decode(&jobs)
	if err != nil {
		panic(err)
	}
	return jobs
}
