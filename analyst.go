package main

import (
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
