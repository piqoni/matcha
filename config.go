package main

import (
	"database/sql"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

var config string = `markdown_dir_path: 
feeds:
  - http://hnrss.org/best 10
  - https://waitbutwhy.com/feed
  - http://tonsky.me/blog/atom.xml
  - http://www.joelonsoftware.com/rss.xml
  - https://www.youtube.com/feeds/videos.xml?channel_id=UCHnyfMqiRRG1u-2MsSQLbXA
google_news_keywords: George Hotz,ChatGPT,Copenhagen 
instapaper: true 
weather_latitude: 37.77
weather_longitude: 122.41
terminal_mode: false
opml_file_path: 
markdown_file_prefix: 
markdown_file_suffix:
reading_time: false 
openai_api_key: 
summary_feeds: `

func parseOPML(xmlContent []byte) []RSS {
	o := Opml{}
	OpmlSlice := []RSS{}
	decoder := xml.NewDecoder(strings.NewReader(string(xmlContent)))
	decoder.Strict = false
	if err := decoder.Decode(&o); err != nil {
		log.Println(err)
	}
	for _, outline := range o.Body.Outline {
		if outline.XmlUrl != "" {
			OpmlSlice = append(OpmlSlice, RSS{url: outline.XmlUrl, limit: 20})
		}
		for _, feed := range outline.Outline {
			if feed.XmlUrl != "" {
				OpmlSlice = append(OpmlSlice, RSS{url: feed.XmlUrl, limit: 20})
			}
		}
	}
	return OpmlSlice
}

func getFeedAndLimit(feedURL string) (string, int) {
	var limit = 20 // default limit
	chopped := strings.Split(feedURL, " ")
	if len(chopped) > 1 {
		var err error
		limit, err = strconv.Atoi(chopped[1])
		if err != nil {
			check(err)
		}
	}
	return chopped[0], limit
}

func bootstrapConfig() {
	currentDir, direrr := os.Getwd()
	if direrr != nil {
		log.Println(direrr)
	}
	// if -t parameter is passed overwrite terminal_mode setting in config.yml
	flag.BoolVar(&terminal_mode, "t", terminal_mode, "Run Matcha in Terminal Mode, no markdown files will be created")
	configFile := flag.String("c", "", "Config file path (if you want to override the current directory config.yaml)")
	opmlFile := flag.String("o", "", "OPML file path to append feeds from opml files")
	build := flag.Bool("build", false, "Dev: Build matcha binaries in the bin directory")
	flag.Parse()

	if *build {
		buildBinaries()
		os.Exit(0)
	}

	// if -c parameter is passed overwrite config.yaml setting in config.yaml
	if len(*configFile) > 0 {
		viper.SetConfigFile(*configFile)
	} else {
		viper.AddConfigPath(".")
		generateConfigFile(currentDir)
		viper.SetConfigName("config")
	}

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Print(err)
		panic("Error reading yaml configuration file")
	}

	if viper.IsSet("markdown_dir_path") {
		markdown_dir_path = viper.Get("markdown_dir_path").(string)
	} else {
		markdown_dir_path = currentDir
	}
	myFeeds = []RSS{}
	feeds := viper.Get("feeds")
	if viper.IsSet("weather_latitude") {
		lat = viper.Get("weather_latitude").(float64)
	}
	if viper.IsSet("weather_longitude") {
		lon = viper.Get("weather_longitude").(float64)
	}
	if viper.IsSet("markdown_file_prefix") {
		mdPrefix = viper.Get("markdown_file_prefix").(string)
	}
	if viper.IsSet("markdown_file_suffix") {
		mdSuffix = viper.Get("markdown_file_suffix").(string)
	}
	if viper.IsSet("openai_api_key") {
		openaiApiKey = viper.Get("openai_api_key").(string)
	}

	if viper.IsSet("summary_feeds") {
		summaryFeeds := viper.Get("summary_feeds")

		for _, summaryFeed := range summaryFeeds.([]any) {
			url, limit := getFeedAndLimit(summaryFeed.(string))
			myFeeds = append(myFeeds, RSS{url: url, limit: limit, summarize: true})
		}
	}

	for _, feed := range feeds.([]any) {
		url, limit := getFeedAndLimit(feed.(string))
		myFeeds = append(myFeeds, RSS{url: url, limit: limit})
	}

	if viper.IsSet("google_news_keywords") {
		googleNewsKeywords := url.QueryEscape(viper.Get("google_news_keywords").(string))
		if googleNewsKeywords != "" {
			googleNewsUrl := "https://news.google.com/rss/search?hl=en-US&gl=US&ceid=US%3Aen&oc=11&q=" + strings.Join(strings.Split(googleNewsKeywords, "%2C"), "%20%7C%20")
			myFeeds = append(myFeeds, RSS{url: googleNewsUrl, limit: 15}) // #FIXME make it configurable
		}
	}

	// Import any config.opml file on current direcotory
	configPath := currentDir + "/" + "config.opml"
	if _, err := os.Stat(configPath); err == nil {
		xmlContent, _ := ioutil.ReadFile(currentDir + "/" + "config.opml")
		myFeeds = append(myFeeds, parseOPML(xmlContent)...)
	}
	// Append any opml file added by -o parameter
	if len(*opmlFile) > 0 {
		xmlContent, _ := ioutil.ReadFile(*opmlFile)
		myFeeds = append(myFeeds, parseOPML(xmlContent)...)
	}

	// Append opml file from config.yml
	if viper.IsSet("opml_file_path") {
		xmlContent, _ := ioutil.ReadFile(viper.Get("opml_file_path").(string))
		myFeeds = append(myFeeds, parseOPML(xmlContent)...)
	}

	instapaper = viper.GetBool("instapaper")
	reading_time = viper.GetBool("reading_time")

	// Overwrite terminal_mode from config file only if its not set through -t flag
	if !terminal_mode {
		terminal_mode = viper.GetBool("terminal_mode")
	}

	databaseFilePath := viper.GetString("database_file_path")
	if databaseFilePath == "" {
		databaseDirPath, err := os.UserConfigDir()
		check(err)
		databaseFilePath = filepath.Join(databaseDirPath, "brew", "matcha.db")
		check(os.MkdirAll(filepath.Dir(databaseFilePath), os.ModePerm))
	}

	db, err = sql.Open("sqlite", databaseFilePath)
	check(err)
	err = applyMigrations(db)
	if err != nil {
		log.Println("Coudn't apply migrations:", err)
	}

	if !terminal_mode {
		markdown_file_name := mdPrefix + currentDate + mdSuffix + ".md"
		err := os.Remove(filepath.Join(markdown_dir_path, markdown_file_name))
		if err != nil {
			// fmt.Println("INFO: Coudn't remove old file: ", err)
		}
	}
}

func generateConfigFile(currentDir string) {
	configPath := currentDir + "/" + "config.yaml"
	if _, err := os.Stat(configPath); err == nil {
		// File exists, dont do anything
		return
	}
	f, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}

	if _, err := f.Write([]byte(config)); err != nil {
		log.Fatal(err)
	}
}
