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

func NewMarkdownWriter(cfg *Config, dateStr string) *MarkdownWriter {
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}
	fname := cfg.MarkdownFilePrefix + dateStr + cfg.MarkdownFileSuffix + ".md"
	return &MarkdownWriter{
		FilePath: filepath.Join(cfg.MarkdownDirPath, fname),
	}
}

func (w MarkdownWriter) WriteFeedHeaderRaw(title, articleURL string) string {
	// Logic to reconstruct the Favicon HTML without a gofeed object
	var src string
	if strings.Contains(title, "Hacker News") {
		src = "https://news.ycombinator.com/favicon.ico"
	} else if articleURL == "" {
		return fmt.Sprintf("\n### 🍵 %s\n", title)
	} else {
		u, err := url.Parse(articleURL)
		if err != nil {
			return fmt.Sprintf("\n### 🍵 %s\n", title)
		}
		// Google S2 Favicon service
		src = "https://www.google.com/s2/favicons?sz=32&domain=" + u.Hostname()
	}

	faviconImg := fmt.Sprintf(`<img src="%s" width="32" height="32" />`, src)
	return fmt.Sprintf("\n### %s %s\n", faviconImg, title)
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
		content += " \n"
	}
	return content
}

func (w MarkdownWriter) WriteSummary(content string, newLine bool) string {
	if content == "" {
		return ""
	}
	if newLine {
		content += " \n\n"
	}
	return content
}

func (w MarkdownWriter) WriteHeader(feed *gofeed.Feed) string {
	favicon := w.getFaviconHTML(feed)
	return fmt.Sprintf("\n### %s %s\n", favicon, escapeMd(feed.Title))
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

func escapeMd(raw string) (safe string) {
	out := strings.Builder{}
	out.Grow(len(raw))
	i := 0
	for i < len(raw) {
		b := raw[i] // we're only escaping ascii characters, so iterating over bytes is fine
		switch b {
		case '\\', '<', '>', '[', ']':
			if err := out.WriteByte('\\'); err != nil {
				panic(err)
			}
		}
		if err := out.WriteByte(b); err != nil {
			panic(err)
		}
	}

	return out.String()
}
