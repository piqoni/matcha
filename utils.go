package main

import "regexp"

const regex = `<.*?>`

// This method uses a regular expresion to remove HTML tags.
func stripHtmlRegex(s string) string {
	r := regexp.MustCompile(regex)
	return r.ReplaceAllString(s, "")
}
