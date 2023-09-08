package main

import (
	"fmt"

	"github.com/mmcdole/gofeed"
	"github.com/savioxavier/termlink"
)

type TerminalWriter struct{}

func (w TerminalWriter) write(body string) {
	fmt.Println(body)
}

func (w TerminalWriter) writeLink(title string, url string, newline bool, readingTime string) string {
	var content string
	content = termlink.Link(title, url)
	if readingTime != "" {
		content += " (" + readingTime + ")"
	}
	if newline {
		content += "\n"
	}
	return content
}

func (w TerminalWriter) writeSummary(content string, newline bool) string {
	if newline {
		content += "\n"
	}
	return content
}

func (w TerminalWriter) writeFavicon(s *gofeed.Feed) string {
	return ""
}
