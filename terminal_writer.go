package main

import (
	"fmt"

	"github.com/mmcdole/gofeed"
	"github.com/savioxavier/termlink"
)

type TerminalWriter struct{}

func (w TerminalWriter) Write(body string) {
	fmt.Println(body)
}

func (w TerminalWriter) WriteLink(title string, url string, newline bool, readingTime string) string {
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

func (w TerminalWriter) WriteSummary(content string, newline bool) string {

	if newline {
		content += "\n"
	}
	return content
}

func (w TerminalWriter) WriteHeader(feed *gofeed.Feed) string {
	return fmt.Sprintf("\n### 🍵 %s\n", feed.Title)
}
