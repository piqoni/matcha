package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	readability "github.com/go-shiori/go-readability"
	"github.com/mmcdole/gofeed"
	"github.com/savioxavier/termlink"
	_ "modernc.org/sqlite"
)

var markdown_dir_path string
var mdPrefix, mdSuffix string
var terminal_mode bool = false
var currentDate = time.Now().Format("2006-01-02")
var lat, lon float64
var instapaper bool
var openaiApiKey string
var reading_time bool
var myMap []RSS
var db *sql.DB

type RSS struct {
	url       string
	limit     int
	summarize bool
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func writeLink(title string, url string, newline bool, readingTime string) string {
	var content string
	if terminal_mode {
		content = termlink.Link(title, url)
	} else {
		content = "[" + title + "](" + url + ")"
	}
	if readingTime != "" {
		content += " (" + readingTime + ")"
	}
	if newline {
		content += "\n"
	}
	return content
}

func getReadingTime(link string) string {
	article, err := readability.FromURL(link, 30*time.Second)
	if err != nil {
		return "" // Just dont display any reading time if can't get the article text
	}

	// get number of words in a string
	words := strings.Fields(article.TextContent)

	// assuming average reading time is 200 words per minute calculate reading time of the article
	readingTime := float64(len(words)) / float64(200)
	minutes := int(readingTime)

	// if minutes is zero return an empty string
	if minutes == 0 {
		return ""
	}

	return strconv.Itoa(minutes) + " min"
}

func writeSummary(content string, newline bool) string {
	if content == "" {
		return content
	}
	if terminal_mode {
		if newline {
			content += "\n"
		}
	} else {
		if newline {
			content += "  \n\n"
		}
	}
	return content
}

func favicon(s *gofeed.Feed) string {
	if terminal_mode {
		return ""
	}
	var src string
	if s.FeedLink == "" {
		// default feed favicon
		return "🍵"

	} else {
		u, err := url.Parse(s.FeedLink)
		if err != nil {
			fmt.Println(err)
		}
		src = "https://www.google.com/s2/favicons?sz=32&domain=" + u.Hostname()
	}
	// if s.Title contains "hacker news"
	if strings.Contains(s.Title, "Hacker News") {
		src = "https://news.ycombinator.com/favicon.ico"
	}

	//return html image tag of favicon
	return fmt.Sprintf("<img src=\"%s\" width=\"32\" height=\"32\" />", src)
}

func writeToMarkdown(body string) {
	if terminal_mode {
		fmt.Println(body)
	} else {
		markdown_file_name := mdPrefix + currentDate + mdSuffix + ".md"
		f, err := os.OpenFile(filepath.Join(markdown_dir_path, markdown_file_name), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
}

func parseFeedURL(fp *gofeed.Parser, url string) *gofeed.Feed {
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil
	}
	return feed
}

func main() {
	bootstrapConfig()

	// Display weather if lat and lon are set
	if lat != 0 && lon != 0 {
		writeToMarkdown(getWeather(lat, lon))
	}

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

// Parses the feed URL and returns the feed object
func parseFeed(fp *gofeed.Parser, url string, limit int) *gofeed.Feed {
	feed, err := fp.ParseURL(url)
	if err != nil {
		fmt.Printf("Error parsing %s with error: %s", url, err)
		return nil
	}

	if len(feed.Items) > limit {
		feed.Items = feed.Items[:limit]
	}

	return feed
}

// Generates the feed items and returns them as a string
func generateFeedItems(feed *gofeed.Feed, rss *RSS) string {
	var items string

	for _, item := range feed.Items {
		seen, seen_today, summary := isSeenArticle(item)
		if seen {
			continue
		}
		title, link := getFeedTitleAndLink(item)
		if summary == "" {
			summary = getSummary(rss, item, link)
		}
		// Add the comments link if it's a Hacker News feed
		if strings.Contains(feed.Title, "Hacker News") {
			commentsLink, commentsCount := getCommentsInfo(item)
			if commentsCount < 100 {
				items += writeLink("💬 ", commentsLink, false, "")
			} else {
				items += writeLink("🔥 ", commentsLink, false, "")
			}
		}

		// Add the Instapaper link if enabled
		if instapaper && !terminal_mode {
			items += getInstapaperLink(item.Link)
		}

		// Support RSS with no Title (such as Mastodon), use Description instead
		if title == "" {
			title = stripHtmlRegex(item.Description)
		}

		timeInMin := ""
		if reading_time {
			timeInMin = getReadingTime(link)
		}

		items += writeLink(title, link, true, timeInMin)
		if rss.summarize {
			items += writeSummary(summary, true)
		}

		// Add the item to the seen table if not seen today
		if !seen_today {
			addToSeenTable(item.Link, summary)
		}
	}

	return items
}

// Writes the feed and its items to the markdown file
func writeFeedToMarkdown(feed *gofeed.Feed, items string) {
	writeToMarkdown(fmt.Sprintf("\n### %s  %s\n%s", favicon(feed), feed.Title, items))
}

// Returns the title and link for the given feed item
func getFeedTitleAndLink(item *gofeed.Item) (string, string) {
	return item.Title, item.Link
}

// Returns the summary for the given feed item
func getSummary(rss *RSS, item *gofeed.Item, link string) string {
	if !rss.summarize {
		return ""
	}

	summary := getSummaryFromLink(link)
	if summary == "" {
		summary = item.Description
	}

	return summary
}

// Returns the comments link and count for the given feed item
func getCommentsInfo(item *gofeed.Item) (string, int) {
	first_index := strings.Index(item.Description, "Comments URL") + 23
	comments_url := item.Description[first_index : first_index+45]
	// Find Comments number
	first_comments_index := strings.Index(item.Description, "Comments:") + 10
	// replace </p> with empty string
	comments_number := strings.Replace(item.Description[first_comments_index:], "</p>\n", "", -1)
	comments_number_int, _ := strconv.Atoi(comments_number)
	// return the link and the number of comments
	return comments_url, comments_number_int
}

func addToSeenTable(link string, summary string) {
	stmt, err := db.Prepare("INSERT INTO seen(url, date, summary) values(?,?,?)")
	check(err)
	res, err := stmt.Exec(link, currentDate, summary)
	check(err)
	_ = res
	stmt.Close()
}

func getInstapaperLink(link string) string {
	return "[<img height=\"16\" src=\"https://staticinstapaper.s3.dualstack.us-west-2.amazonaws.com/img/favicon.png\">](https://www.instapaper.com/hello2?url=" + link + ")"
}

func isSeenArticle(item *gofeed.Item) (seen bool, today bool, summaryText string) {
	var url string
	var date string
	var summary sql.NullString
	err := db.QueryRow("SELECT url, date, summary FROM seen WHERE url=?", item.Link).Scan(&url, &date, &summary)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		return false, false, ""
	}

	if summary.Valid {
		summaryText = summary.String
	} else {
		summaryText = ""
	}

	seen = url != "" && date != currentDate
	today = url != "" && date == currentDate
	return seen, today, summaryText
}
