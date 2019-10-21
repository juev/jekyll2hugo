package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"gopkg.in/yaml.v2"
)

const fmMagic = "---\n"

var frontMatterMatcher = regexp.MustCompile(`(?s)^---\n(.+?\n)---\n+(.*)`)
var emptyFontMatterMatcher = regexp.MustCompile(`(?s)^---\n+---\n+`)

var src = flag.String("src", "_posts", "Source directory with jekyll posts")
var dst = flag.String("dst", "hugo", "Destination directory with hugo posts")

// FrontMatter struct for unmarshal
type FrontMatter struct {
	// Layout     string
	Title     string
	Date      string
	Tags      []string
	Permalink string
	// Published   bool
	Category    string
	Categories  []string
	Keywords    string
	Description string
	Image       string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// FileHasFrontMatter returns a bool indicating whether the
// file looks like it has frontmatter.
func FileHasFrontMatter(source string) (bool, error) {
	return string(source[:4]) == fmMagic, nil
}

func parse(f string) string {
	source, err := ioutil.ReadFile(f)
	check(err)
	fm := FrontMatter{}
	source = bytes.Replace(source, []byte("\r\n"), []byte("\n"), -1)
	hasFM, err := FileHasFrontMatter(string(source))
	check(err)
	var body string
	if hasFM {
		if match := frontMatterMatcher.FindSubmatch(source); match != nil {
			if err = yaml.Unmarshal(match[1], &fm); err != nil {
				return string(source)
			}
			body = string(match[2])
		}
	}
	moscow, err := time.LoadLocation("Europe/Moscow")
	date, _ := time.ParseInLocation("2006-01-02 15:04", fm.Date, moscow)
	result := "---\n" +
		"title: \"" + fm.Title + "\"\n" +
		"date: \"" + date.Format("2006-01-02T15:04:05-0700") + "\"\n"
	if fm.Tags != nil {
		result += "tags:" + "\n"
		for _, tag := range fm.Tags {
			result += "  - " + tag + "\n"
		}
	}
	if fm.Permalink != "" {
		result += "permalink: " + fm.Permalink + "\n"
	}
	if fm.Category != "" {
		result += "category: " + fm.Category + "\n"
	}
	if fm.Categories != nil {
		result += "category: " + "\n"
		for _, category := range fm.Categories {
			result += "  - " + category + "\n"
		}
	}
	if fm.Keywords != "" {
		result += "keywords: " + fm.Keywords + "\n"
	}
	if fm.Description != "" {
		result += "description: " + fm.Description + "\n"
	}
	if fm.Image != "" {
		result += "image: " + fm.Image + "\n"
	}
	result += "---\n" + body
	return result
}

func main() {
	flag.Parse()
	if _, err := os.Stat(*src); os.IsNotExist(err) {
		fmt.Println(*src + ": file or directory not found")
		os.Exit(1)
	}

	files, err := ioutil.ReadDir(*src)
	check(err)
	os.MkdirAll(*dst, os.ModePerm)
	for _, file := range files {
		full, _ := filepath.Abs(*src + file.Name())
		content := parse(full)
		err := ioutil.WriteFile(*dst+"/"+file.Name(), []byte(content), 0644)
		check(err)
	}
}
