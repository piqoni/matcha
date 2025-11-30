package main

import (
	"encoding/xml"
	"flag"
	"fmt"
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
sunrise_sunset: false
openai_api_key:
openai_base_url:
openai_model:
summary_feeds:
summary_prompt:
show_images: false
analyst_feeds:
  - https://feeds.bbci.co.uk/news/business/rss.xml
analyst_prompt:
analyst_model:
`

type Config struct {
	MarkdownDirPath        string
	MarkdownFilePrefix     string
	MarkdownFileSuffix     string
	Feeds                  []RSS
	GoogleNewsKeywords     string
	Instapaper             bool
	WeatherLat             float64
	WeatherLon             float64
	TerminalMode           bool
	ReadingTime            bool
	SunriseSunset          bool
	ShowImages             bool
	OpenAIKey              string
	OpenAIBaseURL          string
	OpenAIModel            string
	SummaryPrompt          string
	AnalystFeeds           []string
	AnalystPrompt          string
	AnalystModel           string
	DatabaseFilePath       string
	NotificationTrigger    string
	NotificationWebhookURL string
}

type RSS struct {
	url       string
	limit     int
	summarize bool
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("limit", 20)

	terminalMode := flag.Bool("t", false, "Run Matcha in Terminal Mode, no markdown files will be created")
	configFile := flag.String("c", "", "Config file path (if you want to override the current directory config.yaml)")
	opmlFile := flag.String("o", "", "OPML file path to append feeds from opml files")
	build := flag.Bool("build", false, "Dev: Build matcha binaries in the bin directory")
	flag.Parse()

	if *build {
		buildBinaries()
		os.Exit(0)
	}

	if *configFile != "" {
		viper.SetConfigFile(*configFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		// Generate default if not exists
		if _, err := os.Stat("./config.yaml"); os.IsNotExist(err) {
			generateConfigFile()
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	cfg := &Config{
		MarkdownDirPath:        viper.GetString("markdown_dir_path"),
		MarkdownFilePrefix:     viper.GetString("markdown_file_prefix"),
		MarkdownFileSuffix:     viper.GetString("markdown_file_suffix"),
		GoogleNewsKeywords:     viper.GetString("google_news_keywords"),
		Instapaper:             viper.GetBool("instapaper"),
		WeatherLat:             viper.GetFloat64("weather_latitude"),
		WeatherLon:             viper.GetFloat64("weather_longitude"),
		TerminalMode:           viper.GetBool("terminal_mode") || *terminalMode,
		ReadingTime:            viper.GetBool("reading_time"),
		SunriseSunset:          viper.GetBool("sunrise_sunset"),
		ShowImages:             viper.GetBool("show_images"),
		OpenAIKey:              viper.GetString("openai_api_key"),
		OpenAIBaseURL:          viper.GetString("openai_base_url"),
		OpenAIModel:            viper.GetString("openai_model"),
		SummaryPrompt:          viper.GetString("summary_prompt"),
		AnalystFeeds:           viper.GetStringSlice("analyst_feeds"),
		AnalystPrompt:          viper.GetString("analyst_prompt"),
		AnalystModel:           viper.GetString("analyst_model"),
		DatabaseFilePath:       viper.GetString("database_file_path"),
		NotificationTrigger:    viper.GetString("notification_trigger"),
		NotificationWebhookURL: viper.GetString("notification_webhook_url"),
	}

	if cfg.MarkdownDirPath == "" {
		wd, _ := os.Getwd()
		cfg.MarkdownDirPath = wd
	}
	if cfg.DatabaseFilePath == "" {
		cfg.DatabaseFilePath = getDefaultDBPath()
	}

	cfg.Feeds = loadFeeds(cfg, *opmlFile)

	return cfg, nil
}

func loadFeeds(cfg *Config, flagOpml string) []RSS {
	var feeds []RSS
	// Summary Feeds
	rawSumFeeds := viper.Get("summary_feeds")
	if rawSumFeeds != nil {
		for _, f := range rawSumFeeds.([]interface{}) {
			url, limit := getFeedAndLimit(f.(string))
			feeds = append(feeds, RSS{url: url, limit: limit, summarize: true})
		}
	}

	// Standard Feeds
	rawFeeds := viper.Get("feeds")
	if rawFeeds != nil {
		for _, f := range rawFeeds.([]interface{}) {
			url, limit := getFeedAndLimit(f.(string))
			feeds = append(feeds, RSS{url: url, limit: limit})
		}
	}

	// OPML Files (Config dir + Flag)
	opmlPaths := []string{"config.opml", viper.GetString("opml_file_path"), flagOpml}
	for _, path := range opmlPaths {
		if path != "" {
			if content, err := os.ReadFile(path); err == nil {
				feeds = append(feeds, parseOPML(content)...)
			}
		}
	}

	if cfg.GoogleNewsKeywords != "" {
		escaped := url.QueryEscape(cfg.GoogleNewsKeywords)
		googleNewsUrl := "https://news.google.com/rss/search?hl=en-US&gl=US&ceid=US%3Aen&oc=11&q=" + strings.Join(strings.Split(escaped, "%2C"), "%20%7C%20") // TODO
		feeds = append(feeds, RSS{url: googleNewsUrl, limit: 15})                                                                                             // #FIXME make it configurable
	}

	return feeds
}

func getDefaultDBPath() string {
	dir, _ := os.UserConfigDir()
	path := filepath.Join(dir, "brew", "matcha.db")
	_ = os.MkdirAll(filepath.Dir(path), 0755)
	return path
}

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
			log.Fatalf("Error getting limit on feed: %v", err)
		}
	}
	return chopped[0], limit
}

func generateConfigFile() {
	currentDir, _ := os.Getwd()
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
