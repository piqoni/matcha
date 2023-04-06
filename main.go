package main

import (
	"github.com/mmcdole/gofeed"
	_ "modernc.org/sqlite"
)

func main() {
	bootstrapConfig()
	displayWeather()


	fp := gofeed.NewParser()
	for _, rss := range myMap {
		feed := parseFeed(fp, rss.url, rss.limit)

		if feed == nil {
			continue
		}

		items := generateFeedItems(feed, &rss)
		if items != "" {
			writeFeedToMarkdown(feed, items)
		}
	}

	// Close the database connection after processing all the feeds
	defer db.Close()
}
