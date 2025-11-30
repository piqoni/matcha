package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	readability "github.com/go-shiori/go-readability"
	"github.com/mmcdole/gofeed"
)

func ProcessFeed(rss RSS, cfg *Config, store *Storage, llm *LLMClient, w Writer, fp *gofeed.Parser) {
	feed, err := fp.ParseURL(rss.url)
	if err != nil {
		log.Printf("Error parsing %s: %v", rss.url, err)
		return
	}

	if len(feed.Items) > rss.limit {
		feed.Items = feed.Items[:rss.limit]
	}

	var outputBuffer string
	itemsFound := false

	for _, item := range feed.Items {
		// Check DB for seen status
		seen, seenToday, prevSummary := store.IsSeen(item.Link)

		// If we've seen it before (and not today), skip it
		if seen {
			continue
		}

		itemsFound = true

		title := item.Title
		if title == "" {
			title = stripHtmlRegex(item.Description)
		}

		summary := prevSummary

		if summary == "" && rss.summarize {
			summary = getSummary(llm, item, cfg)
		}

		var readingTime string
		if cfg.ReadingTime {
			readingTime = getReadingTime(item.Link)
		}

		if strings.Contains(feed.Link, "news.ycombinator.com") {
			outputBuffer += formatHackerNewsLinks(w, item)
		}

		if cfg.Instapaper && !cfg.TerminalMode {
			outputBuffer += getInstapaperLink(item.Link)
		}

		outputBuffer += w.WriteLink(title, item.Link, true, readingTime)

		if rss.summarize {
			outputBuffer += w.WriteSummary(summary, true)
		}

		if cfg.ShowImages && !cfg.TerminalMode {
			img := extractImageTagFromHTML(item.Content)
			if img != "" {
				outputBuffer += img + "\n"
			}
		}

		if !seenToday {
			store.MarkAsSeen(item.Link, summary)
		}
	}

	if itemsFound && outputBuffer != "" {
		header := w.WriteHeader(feed)
		w.Write(header + outputBuffer)
	}
}

func getSummary(llm *LLMClient, item *gofeed.Item, cfg *Config) string {
	fmt.Println("outside")
	if llm != nil {
		scrapedText, err := readability.FromURL(item.Link, 30*time.Second)
		content := item.Description
		if err == nil {
			content = scrapedText.TextContent
		}

		fmt.Println("we have crawled")

		return llm.Summarize(content)
	}

	return item.Description
}

func getReadingTime(link string) string {
	article, err := readability.FromURL(link, 30*time.Second)
	if err != nil {
		return ""
	}

	words := strings.Fields(article.TextContent)
	if len(words) == 0 {
		return ""
	}

	// 200 wpm
	minutes := len(words) / 200
	if minutes == 0 {
		return "" // < 1 min
	}
	return strconv.Itoa(minutes) + " min"
}

func extractImageTagFromHTML(htmlText string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlText))
	if err != nil {
		return ""
	}

	imgTags := doc.Find("img")
	if imgTags.Length() == 0 {
		return ""
	}

	firstImgTag := imgTags.First()

	// Resize logic
	width := firstImgTag.AttrOr("width", "")
	height := firstImgTag.AttrOr("height", "")

	if width != "" && height != "" {
		wInt, _ := strconv.Atoi(width)
		hInt, _ := strconv.Atoi(height)
		if wInt > 0 && hInt > 0 {
			aspectRatio := float64(wInt) / float64(hInt)
			const maxWidth = 400
			if wInt > maxWidth {
				wInt = maxWidth
				hInt = int(float64(wInt) / aspectRatio)
			}
			firstImgTag.SetAttr("width", fmt.Sprintf("%d", wInt))
			firstImgTag.SetAttr("height", fmt.Sprintf("%d", hInt))
		}
	}

	html, err := goquery.OuterHtml(firstImgTag)
	if err != nil {
		return ""
	}
	return html
}

func formatHackerNewsLinks(w Writer, item *gofeed.Item) string {
	desc := item.Description

	commentsURL := ""
	if start := strings.Index(desc, "Comments URL"); start != -1 {
		safeStart := start + 23
		if safeStart+45 < len(desc) {
			commentsURL = desc[safeStart : safeStart+45]
		}
	}

	// Find count
	count := 0
	if start := strings.Index(desc, "Comments:"); start != -1 {
		s := desc[start+10:]
		s = strings.Replace(s, "</p>\n", "", -1)
		s = strings.TrimSpace(s)
		count, _ = strconv.Atoi(s)
	}

	icon := "💬 "
	if count >= 100 {
		icon = "🔥 "
	}

	// If parsing failed, default to item.Link (often the comments page for text posts)
	if commentsURL == "" {
		commentsURL = item.Link
	}

	return w.WriteLink(icon, commentsURL, false, "")
}

func getInstapaperLink(link string) string {
	return fmt.Sprintf(`[<img height="16" src="https://staticinstapaper.s3.dualstack.us-west-2.amazonaws.com/img/favicon.png">](https://www.instapaper.com/hello2?url=%s)`, link)
}

func stripHtmlRegex(s string) string {
	const regex = `<.*?>`
	r := regexp.MustCompile(regex)
	return r.ReplaceAllString(s, "")
}
