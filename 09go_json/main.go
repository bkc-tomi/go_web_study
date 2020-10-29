package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type pl struct {
	Language  string `json:"language"`
	Birthyear int    `json:"birthyear"`
	Character string `json:"character"`
}

var strB = `
[
	{"language": "Go", "birthyear": 2009, "character": "gopher"},
	{"language": "JavaScript", "birthyear": 1997, "character": null},
	{"language": "Java", "birthyear": 1995, "character":"Duke"},
	{"language": "D", "birthyear": 2007, "character":"D-man"},
	{"language": "PHP", "birthyear": 1995, "character":"elePHPant"},
	{"language": "LISP", "birthyear": 1958, "character":"lien"},
	{"language": "Python", "birthyear": 1990, "character":"python"}
]
`

func main() {
	f, fOpenErr := os.Open("source/example.json")
	if fOpenErr != nil {
		fmt.Println("file open error:", fOpenErr)
		return
	}
	defer f.Close()
	b, rErr := ioutil.ReadAll(f)
	if rErr != nil {
		fmt.Println("file read error:", rErr)
		return
	}
	var pls []pl
	if err := json.Unmarshal(b, &pls); err != nil {
		fmt.Println(err)
	}

	for i, p := range pls {
		fmt.Println("language:", p.Language, ", birthyear:", p.Birthyear, ", character:", p.Character)
		pls[i].Language = strings.ToUpper(pls[i].Language)
		pls[i].Character = strings.ToUpper(pls[i].Character)
	}
	b, err := json.Marshal(pls)
	if err != nil {
		fmt.Println("encode error:", err)
	}

	if err := ioutil.WriteFile("source/upperExample.json", b, 0664); err != nil {
		fmt.Println("write file error:", err)

	}
}
