package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/savioxavier/termlink"
	"github.com/spf13/viper"
	_ "modernc.org/sqlite"
)

var path string
var terminal_mode bool = false
var currentDate = time.Now().Format("2006-01-02")
var currentDir string

type RSS struct {
	url      string
	limit    int
	disabled bool
}

func check(e error) {
	if e != nil {
		panic(e)
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
		// default favicon #FIXME
		src = "https://www.cloudflare.com/favicon.ico"

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

func main() {
	currentDir, direrr := os.Getwd()
	if direrr != nil {
		log.Println(direrr)
	}
	generateConfigFile(currentDir)
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Print(err)
		panic("Error reading config.yaml file. Please create config.yaml file.")
	}

	if viper.IsSet("markdown_dir_path") {
		path = viper.Get("markdown_dir_path").(string)
	} else {
		path = currentDir
	}

	myMap := []RSS{}
	feeds := viper.Get("feeds")
	lat := viper.Get("weather_latitude").(float64)
	lon := viper.Get("weather_longitude").(float64)
	googleNewsKeywords := url.QueryEscape(viper.Get("google_news_keywords").(string))

	var limit int
	for _, feed := range feeds.([]any) {
		chopped := strings.Split(feed.(string), " ")
		if len(chopped) > 1 {
			limit, err = strconv.Atoi(chopped[1])
			if err != nil {
				check(err)
			}
		} else {
			limit = 10
		}

		myMap = append(myMap, RSS{url: chopped[0], limit: limit})
	}
	if googleNewsKeywords != "" {
		//FIXME
		googleNewsUrl := "https://news.google.com/rss/search?hl=en-US&gl=US&ceid=US%3Aen&oc=11&q=" + strings.Join(strings.Split(googleNewsKeywords, "%2C"), "%20%7C%20")
		myMap = append(myMap, RSS{url: googleNewsUrl, limit: 15}) // #FIXME make it configurable
	}
	instapaper := viper.GetBool("instapaper")
	terminal_mode = viper.GetBool("terminal_mode")

	// if -t parameter is passed overwrite terminal_mode setting in config.yml
	flag.BoolVar(&terminal_mode, "t", terminal_mode, "run in terminal mode")
	flag.Parse()

	dbPath, err := os.UserConfigDir()
	check(err)

	if _, err := os.Stat(dbPath + "/brew"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dbPath+"/brew", os.ModePerm)
		if err != nil {
			check(err)
		}
	}

	db, err := sql.Open("sqlite", dbPath+"/brew/matcha.db")
	check(err)
	// create new table on database
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS seen (url TEXT, date TEXT)")
	check(err)
	defer db.Close()

	if !terminal_mode {
		err := os.Remove(path + "/" + currentDate + ".md")
		if err != nil {
			// fmt.Println("INFO: Coudn't remove old file: ", err)
		}
	}

	// Start writing to markdown
	// Display weather
	writeToMarkdown(getWeather(lat, lon))

	fp := gofeed.NewParser()
	for _, rss := range myMap {
		if rss.disabled {
			continue
		}
		feed, err := fp.ParseURL(rss.url)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Failed to parse feed: " + rss.url)
		}

		if feed == nil {
			continue
		}

		items := ""
		for index, item := range feed.Items {
			if index == rss.limit {
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

			// Mastodon RSS has not Title, use Description instead
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

	}
}
