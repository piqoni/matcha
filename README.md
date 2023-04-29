<img align="right" src="https://github.com/piqoni/matcha/actions/workflows/test.yml/badge.svg">
<h1 align="center"> 🍵 Matcha </h1>
<div align="center"><p>
    <a href="https://github.com/piqoni/matcha/releases/latest">
      <img alt="Latest release" src="https://img.shields.io/github/v/release/piqoni/matcha?style=for-the-badge&logo=starship&color=C9CBFF&logoColor=D9E0EE&labelColor=302D41" />
    </a>
    <a href="https://github.com/piqoni/matcha/pulse">
      <img alt="Last commit" src="https://img.shields.io/github/last-commit/piqoni/matcha?style=for-the-badge&logo=starship&color=8bd5ca&logoColor=D9E0EE&labelColor=302D41"/>
    </a>
    <a href="https://github.com/piqoni/matcha/blob/main/LICENSE">
      <img alt="License" src="https://img.shields.io/github/license/piqoni/matcha?style=for-the-badge&logo=starship&color=ee999f&logoColor=D9E0EE&labelColor=302D41" />
    </a>
    <a href="https://github.com/piqoni/matcha/stargazers">
      <img alt="Stars" src="https://img.shields.io/github/stars/piqoni/matcha?style=for-the-badge&logo=starship&color=c69ff5&logoColor=D9E0EE&labelColor=302D41" />
    </a>
</div>

Matcha is a daily digest generator for your RSS feeds and interested topics/keywords. By using any markdown file viewer (such as [Obsidian](https://obsidian.md/)) or directly from terminal (-t option), you can read your RSS articles whenever you want at your pace, thus avoiding FOMO throughout the day. 
### In Obsidian
<img width="900" alt="image" src="https://user-images.githubusercontent.com/3144671/219786799-55db70c1-5860-4d4b-9df4-b81a89f8161d.png">

### On the terminal 

<img width="596" alt="image" src="https://user-images.githubusercontent.com/3144671/208323296-af2d6a51-7d33-42a9-a827-0e96a4a383fd.png">

## Features
 - RSS daily **digest**, it will show only articles not previously seen
 - Optional summary of articles from OpenAI for selected feeds
 - Weather for the next 12 Hours (from [YR](https://www.yr.no/))
 - Quick bookmarking of articles to Instapaper
 - Interested Topics/Keywords to follow (through [Google News](https://news.google.com/))
 - Hacker News comments direct link and distinguishing mostly dicussed posts 🔥
 - Terminal Mode by calling `./matcha -t` 

 
## Installation / Usage
1. Since Matcha generates markdown, any markdown reader should do the job. Currently it has been tested on [Obsidian](https://obsidian.md/) so you need a markdown reader before moving on, unless you will use terminal mode (-t option), then a markdown reader is not needed. 
2. **Download the [corresponding binary](https://github.com/piqoni/matcha/releases)** based on your OS and after executing (if on mac/linux run `chmod +x matcha-darwin-amd64` to make it executable), a sample `config.yml` will be generated, which you can add in your rss feeds, keywords and the `markdown_dir_path` where you want the markdown files to be generated (if left empty, it will generate the daily digest on current dir). 
3. You can either execute matcha on-demand (a terminal alias) or set a cron to run matcha as often as you want. Even if you set it to execute every hour, matcha will still generate daily digests, one file per day, and will add more articles to it if new articles are published throughout the day. 

## Configuration
On first execution, Matcha will generate the following config.yaml and a markdown file on the same directory as the application. Change the 'feeds' to your actual RSS feeds, and google_news_keywords to the keywords you are interested in. And if you want to change where the markdown files are generated, set the full directory path in `markdown_dir_path`. 

```yaml 
markdown_dir_path: 
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
summary_feeds: 
```
### GPT-3 Summary of Articles
To use the summary feature, you need to an OpanAI account ([sign up here](https://openai.com/api/)) then get your API key [here](https://openai.com/api/).

Then update the configuration with the feeds you want to be summarized under "summary_feeds" setting, example:
```
openai_api_key: xxxxxxxxxxxxxxxxx
summary_feeds:
    - http://hnrss.org/best
```
### Command line Options
Run matcha with --help option to see current cli options:
```
  -c filepath
    	Config file path (if you want to override the current directory config.yaml)
  -o filepath
    	OPML file path to append feeds from opml files
  -t	Run Matcha in Terminal Mode, no markdown files will be created
```  

#### OPML Import
To use OPML files (exported from other services), rename your file to `config.opml` and leave it in the directory where matcha is located. The other option is to run the command with -o option pointing to the opml filepath. 

<div align="center">
<a href="https://www.buymeacoffee.com/matchaRSS" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/default-orange.png" alt="Buy Me A Coffee" height="41" width="174"></a>
</div>
