package main

import (
	"fmt"
	"regexp"
)

const regex = `<.*?>`

// This method uses a regular expresion to remove HTML tags.
func stripHtmlRegex(s string) string {
	r := regexp.MustCompile(regex)
	return r.ReplaceAllString(s, "")
}

func stringToInt(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}
