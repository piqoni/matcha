package main

import (
	"github.com/mmcdole/gofeed"
	_ "modernc.org/sqlite"
)

func main() {
	bootstrapConfig()

	fp := gofeed.NewParser()
	writer := getWriter()
	displayWeather(writer)

	for _, feed := range myFeeds {
		parsedFeed := parseFeed(fp, feed.url, feed.limit)

		if parsedFeed == nil {
			continue
		}

		items := generateFeedItems(writer, parsedFeed, feed)
		if items != "" {
			writeFeed(writer, parsedFeed, items)
		}
	}

	// Close the database connection after processing all the feeds
	defer db.Close()
}
