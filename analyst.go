package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/mmcdole/gofeed"
)

const (
	analystTag   = "#analyst"
	defaultLimit = 20
)

func RunAnalyst(cfg *Config, store *Storage, llm *LLMClient, writer Writer, fp *gofeed.Parser) {
	if len(cfg.AnalystFeeds) == 0 || cfg.AnalystPrompt == "" {
		return
	}
	if llm == nil {
		return
	}

	articles := collectArticlesForAnalysis(cfg, store, fp)
	if len(articles) == 0 {
		return
	}

	analysis := llm.Analyze(articles)

	if analysis != "" {
		writer.Write("\n## Daily Analysis:\n")
		writer.Write(analysis + "\n")

		if cfg.NotificationTrigger != "" &&
			cfg.NotificationWebhookURL != "" &&
			strings.TrimSpace(analysis) == strings.TrimSpace(cfg.NotificationTrigger) &&
			!store.WasFeedNotifiedToday(cfg.AnalystFeeds[0]) {

			if err := sendNotification(cfg.NotificationWebhookURL, analysis); err == nil {
				_ = store.MarkFeedNotified(cfg.AnalystFeeds[0])
			}
		}
	}
}

func collectArticlesForAnalysis(cfg *Config, store *Storage, fp *gofeed.Parser) []string {
	var articles []string

	for _, feedURL := range cfg.AnalystFeeds {
		feed, err := fp.ParseURL(feedURL)
		if err != nil {
			continue
		}

		// Limit items
		limit := defaultLimit
		if len(feed.Items) > limit {
			feed.Items = feed.Items[:limit]
		}

		for _, item := range feed.Items {
			articleLink := item.Link + analystTag
			seen, seenToday, summary := store.IsSeen(articleLink)

			if seen {
				continue
			}

			articleText := item.Title + ": " + item.Description
			articles = append(articles, articleText)

			if !seenToday {
				if err := store.MarkAsSeen(articleLink, summary); err != nil {
					continue
				}
			}
		}
	}

	return articles
}

func sendNotification(url, message string) error {
	// Detect Slack webhook
	if strings.Contains(url, "hooks.slack.com") {
		// Slack requires JSON payload
		payload := map[string]string{"text": message}
		body, _ := json.Marshal(payload)

		resp, err := http.Post(url, "application/json", bytes.NewReader(body))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		io.Copy(io.Discard, resp.Body)
		return nil
	}

	// Default: ntfy.sh or generic text webhook
	resp, err := http.Post(url, "text/plain", strings.NewReader(message))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	return nil
}
