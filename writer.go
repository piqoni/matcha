package main

import (
	"github.com/mmcdole/gofeed"
)

type Writer interface {
	Write(body string)
	WriteLink(title string, url string, newLine bool, readingTime string) string
	WriteSummary(content string, newLine bool) string
	WriteHeader(feed *gofeed.Feed) string
}
