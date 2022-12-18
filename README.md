# 🍵 Matcha
Matcha is a daily digest generator for your RSS feeds and interested topics/keywords. By using any markdown file viewer (such as [Obsidian](https://obsidian.md/)), you can read your RSS articles whenever you want at your pace, thus avoiding FOMO throughout the day. 

<img width="900" alt="image" src="https://user-images.githubusercontent.com/3144671/206862015-9a325a14-cd8b-4ac3-97bc-55c81008c0df.png">

Terminal Mode: 

<img width="596" alt="image" src="https://user-images.githubusercontent.com/3144671/208323296-af2d6a51-7d33-42a9-a827-0e96a4a383fd.png">

## Features
 - RSS daily **digest**, it will show only articles not previously seen
 - Weather for the next 12 Hours (from [YR](https://www.yr.no/))
 - Quick bookmarking of articles to Instapaper
 - Interested Topics/Keywords to follow (through [Google News](https://news.google.com/))
 - Hacker News comments direct link and distinguishing mostly dicussed posts 🔥
 - Terminal Mode by calling `./matcha -t` 

 
## Installation / Usage
1. Since Matcha generates markdown, any markdown reader should do the job. Currently it has been tested on [Obsidian](https://obsidian.md/) so you need a markdown reader before moving on. 
2. **Download the [corresponding binary](https://github.com/piqoni/matcha/releases)** based on your OS and after executing, a sample `config.yml` will be generated, which you can add in your rss feeds, keywords and the `markdown_dir_path` where you want the markdown files to be generated (if left empty, it will generate the daily digest on current dir). 
3. You can either execute matcha on-demand (a terminal alias) or set a cron to run matcha as often as you want. Even if you set it to execute every hour, matcha will still generate daily digests, one file per day, and will add more articles to it if new articles are published throughout the day. 
