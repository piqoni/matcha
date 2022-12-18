package main

import (
	"fmt"
	"log"
	"os"
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
