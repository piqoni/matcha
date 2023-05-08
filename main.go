package main

import (
	"github.com/mmcdole/gofeed"
	_ "modernc.org/sqlite"
)

func main() {
	bootstrapConfig()
	displayWeather()

	fp := gofeed.NewParser()
	for _, feed := range myFeeds {
		parsedFeed := parseFeed(fp, feed.url, feed.limit)

		if parsedFeed == nil {
			continue
		}

		items := generateFeedItems(parsedFeed, feed)
		if items != "" {
			writeFeedToMarkdown(parsedFeed, items)
		}
	}

	// Close the database connection after processing all the feeds
	defer db.Close()
}
