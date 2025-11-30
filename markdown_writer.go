package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

type MarkdownWriter struct {
	FilePath string
}

func NewMarkdownWriter(cfg *Config) *MarkdownWriter {
	date := time.Now().Format("2006-01-02")
	fname := cfg.MarkdownFilePrefix + date + cfg.MarkdownFileSuffix + ".md"
	return &MarkdownWriter{
		FilePath: filepath.Join(cfg.MarkdownDirPath, fname),
	}
}

func (w MarkdownWriter) Write(body string) {
	f, err := os.OpenFile(w.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening file %s: %v", w.FilePath, err)
		return
	}
	defer f.Close()

	if _, err := f.Write([]byte(body)); err != nil {
		log.Printf("Error writing to file: %v", err)
	}
}

func (w MarkdownWriter) WriteLink(title string, url string, newLine bool, readingTime string) string {
	content := fmt.Sprintf("[%s](%s)", title, url)
	if readingTime != "" {
		content += fmt.Sprintf(" (%s)", readingTime)
	}

	if newLine {
		content += "  \n"
	}
	return content
}

func (w MarkdownWriter) WriteSummary(content string, newLine bool) string {
	if content == "" {
		return ""
	}

	if newLine {
		content += "  \n\n"
	}
	return content
}

func (w MarkdownWriter) WriteHeader(feed *gofeed.Feed) string {
	favicon := w.getFaviconHTML(feed)
	return fmt.Sprintf("\n### %s  %s\n", favicon, feed.Title)
}

// Helper method specifically for MarkdownWriter
func (w MarkdownWriter) getFaviconHTML(s *gofeed.Feed) string {
	var src string

	// Hacker news is a special case
	if strings.Contains(s.Title, "Hacker News") {
		src = "https://news.ycombinator.com/favicon.ico"
	} else if s.FeedLink == "" {
		return "🍵"
	} else {
		u, err := url.Parse(s.FeedLink)
		if err != nil {
			// If URL parsing fails, just return emoji
			return "🍵"
		}
		src = "https://www.google.com/s2/favicons?sz=32&domain=" + u.Hostname()
	}

	return fmt.Sprintf(`<img src="%s" width="32" height="32" />`, src)
}
