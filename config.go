package main

import (
	"database/sql"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
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
terminal_mode: false`

func parseOPML(xmlContent []byte) []RSS {
	o := Opml{}
	OpmlSlice := []RSS{}
	decoder := xml.NewDecoder(strings.NewReader(string(xmlContent)))
	decoder.Strict = false
	if err := decoder.Decode(&o); err != nil {
		log.Println(err)
	}
	for _, outline := range o.Body.Outline {
		OpmlSlice = append(OpmlSlice, RSS{url: outline.XmlUrl, limit: 20})
		for _, feed := range outline.Outline {
			OpmlSlice = append(OpmlSlice, RSS{url: feed.XmlUrl, limit: 20})
		}
	}
	return OpmlSlice
}

func bootstrapConfig() {
	currentDir, direrr := os.Getwd()
	if direrr != nil {
		log.Println(direrr)
	}
	// if -t parameter is passed overwrite terminal_mode setting in config.yml
	flag.BoolVar(&terminal_mode, "t", terminal_mode, "run in terminal mode")
	configDir := flag.String("c", "", "Config directory (if you dont want to use current dir")
	flag.Parse()
	if len(*configDir) > 0 {
		viper.AddConfigPath(*configDir)
	} else {
		viper.AddConfigPath(".")
	}

	flag.Parse()
	generateConfigFile(currentDir)
	viper.SetConfigName("config")

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
	myMap = []RSS{}
	feeds := viper.Get("feeds")
	lat = viper.Get("weather_latitude").(float64)
	lon = viper.Get("weather_longitude").(float64)
	googleNewsKeywords := url.QueryEscape(viper.Get("google_news_keywords").(string))
	//var err error
	var limit int
	for _, feed := range feeds.([]any) {
		chopped := strings.Split(feed.(string), " ")
		if len(chopped) > 1 {
			limit, err = strconv.Atoi(chopped[1])
			if err != nil {
				check(err)
			}
		} else {
			limit = 20
		}

		myMap = append(myMap, RSS{url: chopped[0], limit: limit})
	}
	if googleNewsKeywords != "" {
		googleNewsUrl := "https://news.google.com/rss/search?hl=en-US&gl=US&ceid=US%3Aen&oc=11&q=" + strings.Join(strings.Split(googleNewsKeywords, "%2C"), "%20%7C%20")
		myMap = append(myMap, RSS{url: googleNewsUrl, limit: 15}) // #FIXME make it configurable
	}

	// Import any config.opml file on current direcot
	configPath := currentDir + "/" + "config.opml"
	if _, err := os.Stat(configPath); err == nil {
		xmlContent, _ := ioutil.ReadFile(currentDir + "/" + "config.opml")
		myMap = append(myMap, parseOPML(xmlContent)...)
	}

	instapaper = viper.GetBool("instapaper")

	// Overwrite terminal_mode from config file only if its not set through -t flag
	if !terminal_mode {
		terminal_mode = viper.GetBool("terminal_mode")
	}

	dbPath, err := os.UserConfigDir()
	check(err)

	if _, err := os.Stat(dbPath + "/brew"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dbPath+"/brew", os.ModePerm)
		if err != nil {
			check(err)
		}
	}

	db, err = sql.Open("sqlite", dbPath+"/brew/matcha.db")
	check(err)
	// create new table on database
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS seen (url TEXT, date TEXT)")
	check(err)

	if !terminal_mode {
		err := os.Remove(path + "/" + currentDate + ".md")
		if err != nil {
			// fmt.Println("INFO: Coudn't remove old file: ", err)
		}
	}
}

func generateConfigFile(currentDir string) {
	configPath := currentDir + "/" + "config.yaml"
	if _, err := os.Stat(configPath); err == nil {
		// File exists, dont do anything
	} else {
		f, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
		if err != nil {
			fmt.Println(err)
		}
		if _, err := f.Write([]byte(config)); err != nil {
			log.Fatal(err)
		}
	}

}
