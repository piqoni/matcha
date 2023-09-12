package main

import (
	"log"
	"os"
	"path/filepath"
)

type MarkdownWriter struct{}

func (w MarkdownWriter) write(body string) {
	markdown_file_name := mdPrefix + currentDate + mdSuffix + ".md"
	f, err := os.OpenFile(filepath.Join(markdownDirPath, markdown_file_name), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte(body)); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func (w MarkdownWriter) writeLink(title string, url string, newline bool, readingTime string) string {
	var content string
	content = "[" + title + "](" + url + ")"
	if readingTime != "" {
		content += " (" + readingTime + ")"
	}
	if newline {
		content += "\n"
	}
	return content
}

func (w MarkdownWriter) writeSummary(content string, newline bool) string {
	if content == "" {
		return content
	}
	if newline {
		content += "  \n\n"
	}
	return content
}
