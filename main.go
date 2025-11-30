package main

import (
	"log"
	"os"

	"github.com/mmcdole/gofeed"
	_ "modernc.org/sqlite"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	store, err := NewStorage(cfg.DatabaseFilePath)
	if err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}
	defer store.Close()

	llm := NewLLMClient(cfg)

	var writer Writer
	if cfg.TerminalMode {
		writer = TerminalWriter{}
	} else {
		mw := NewMarkdownWriter(cfg)
		os.Remove(mw.FilePath)
		writer = mw
	}

	fp := gofeed.NewParser()

	DisplayWeather(cfg, writer)
	DisplaySunriseSunset(cfg, writer)

	RunAnalyst(cfg, store, llm, writer, fp)

	for _, feedConfig := range cfg.Feeds {
		ProcessFeed(feedConfig, cfg, store, llm, writer, fp)
	}

}
