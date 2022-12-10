# 🍵 Matcha
Matcha is a simple daily digest generator for your RSS feeds and interested topics/keywords. Using a markdown reader (such as Obsidian), one can navigate through articles at their own pace.

<img width="700" alt="image" src="https://user-images.githubusercontent.com/3144671/206862015-9a325a14-cd8b-4ac3-97bc-55c81008c0df.png">

## Features
 - RSS daily *digest*, it will show only articles not previously generated
 - Weather for the next 12 Hours (from [YR](https://www.yr.no/))
 - Optionally bookmark articles to Instapaper
 - Interested Topics/Keywords to follow (through [Google News](https://news.google.com/))
 - Hacker News comments direct link and distinguishing hot articles 🔥

 
## Installation / Usage
1. Since Matcha generates markdown, any markdown reader should do the job. Currently it has been tested on [Obsidian](https://obsidian.md/).
2. Download the corresponding binary based on your OS and add to the config.yml file your rss feeds, keywords and the markdown_dir_path where you want the markdown files to be generated. (if left empty, it will generate the daily digest on current dir)
3. You can either execute matcha on-demand or set a cron to run matcha as often as you want. Even if you set it to execute every hour, matcha will still generate daily digests, one file per day, and will add more articles to it if new articles are published throughout the day. 
