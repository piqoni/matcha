package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/savioxavier/termlink"
	_ "modernc.org/sqlite"
)

var path string
var terminal_mode bool = false
var currentDate = time.Now().Format("2006-01-02")
var lat float64
var lon float64
var instapaper bool
var myMap []RSS
var db *sql.DB
var feeds chan *gofeed.Feed
var wg sync.WaitGroup

type RSS struct {
	url      string
	limit    int
	disabled bool
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func writeLink(title string, url string) string {
	if terminal_mode {
		return termlink.Link(title, url)
	} else {
		return "[" + title + "](" + url + ")"
	}
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
		f, err := os.OpenFile(path+"/"+currentDate+".md", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func parseFeedURL(fp *gofeed.Parser, url string) {
	feed, err := fp.ParseURL(url)
	if err == nil {
		feeds <- feed
		wg.Done()
	}
}

func main() {
	bootstrapConfig()
	// Start writing to markdown
	// Display weather
	writeToMarkdown(getWeather(lat, lon))

	feeds = make(chan *gofeed.Feed, len(myMap))

	wg.Add(len(myMap))
	// Set up a slice to store the parsed feeds
	var parsedFeeds []*gofeed.Feed

	fp := gofeed.NewParser()
	for _, feed := range myMap {
		go parseFeedURL(fp, feed.url)
	}
	go func() {
		for feed := range feeds {
			parsedFeeds = append(parsedFeeds, feed)
		}
	}()
	wg.Wait()

	for _, feed := range parsedFeeds {

		items := ""
		for index, item := range feed.Items {
			if index == 10 {
				break
			}
			var url string
			var date string
			err := db.QueryRow("SELECT url, date FROM seen where url=?", item.Link).Scan(&url, &date)
			if err != nil && err != sql.ErrNoRows {
				fmt.Println(err)
			}
			if url != "" && date == currentDate {
				// fmt.Println("Already seen: " + item.Title)
				// Article is already in the database and it is for today's date so skip inserting it
			} else if url != "" && date != currentDate {
				// fmt.Println("Skipping: " + item.Link)
				continue
			} else {
				stmt, err := db.Prepare("INSERT INTO seen(url, date) values(?,?)")
				check(err)
				res, err := stmt.Exec(item.Link, currentDate)
				check(err)
				_ = res
				stmt.Close()
			}

			if strings.Contains(feed.Title, "Hacker News") {
				// Find Comments URL
				first_index := strings.Index(item.Description, "Comments URL") + 23
				comments_url := item.Description[first_index : first_index+45]
				// Find Comments number
				first_comments_index := strings.Index(item.Description, "Comments:") + 10
				// replace </p> with empty string
				comments_number := strings.Replace(item.Description[first_comments_index:], "</p>\n", "", -1)
				comments_number_int, _ := strconv.Atoi(comments_number)
				if comments_number_int < 100 {
					items += writeLink("💬 ", comments_url)
				} else {
					items += writeLink("🔥 ", comments_url)
				}
			}
			if instapaper && !terminal_mode {
				items += "[<img height=\"16\" src=\"https://staticinstapaper.s3.dualstack.us-west-2.amazonaws.com/img/favicon.png\">](https://www.instapaper.com/hello2?url=" + item.Link + ")"
			}

			title := item.Title
			link := item.Link

			// Mastodon RSS has no Title, use Description instead
			if title == "" {
				title = stripHtmlRegex(item.Description)
			}
			items += writeLink(title, link)
			items += "\n"
		}

		if items != "" {
			writeToMarkdown("### " + favicon(feed) + "  " + feed.Title + "\n")
			writeToMarkdown(items)
		}
		defer db.Close()

	}
}
