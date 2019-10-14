package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"

	"gopkg.in/yaml.v2"
)

const fmMagic = "---\n"

var frontMatterMatcher = regexp.MustCompile(`(?s)^---\n(.+?\n)---\n+`)
var emptyFontMatterMatcher = regexp.MustCompile(`(?s)^---\n+---\n+`)

// T struct for unmarshal
type T struct {
	// Layout     string
	Title      string
	Date       string
	Tags       []string
	Permalink  string
	Published  bool
	Category   string
	Categories []string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	source, err := ioutil.ReadFile("./example/2019-09-08-programming-language-performance.markdown")
	check(err)
	fm := T{}

	source = bytes.Replace(source, []byte("\r\n"), []byte("\n"), -1)
	if match := frontMatterMatcher.FindSubmatchIndex(source); match != nil {
		if err = yaml.Unmarshal(source[match[2]:match[3]], &fm); err != nil {
			return
		}
	}
	// err = yaml.Unmarshal([]byte(dat), &t)
	// check(err)
	// fmt.Print(string(dat))
	fmt.Printf("--- t:\n%v\n\n", fm)
	d, err := yaml.Marshal(&fm)
	check(err)
	fmt.Printf("--- m dump:\n%s\n\n", string(d))
}
