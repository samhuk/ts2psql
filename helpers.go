package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// Reads a file and returns the text of the file
func readFile(path string) string {
	fileText, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error occured reading in file %v:\n%v\n", path, err)
	}
	return string(fileText)
}

// Converts the given string into snake case, i.e. "PascalCase" -> "pascal_case"
func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
