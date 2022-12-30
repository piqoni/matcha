package main

import (
	"testing"
)

func TestParseOpmlWithoutCategories(t *testing.T) {
	data := `<opml version="1.0">
	<head>
	<title>Sample OPML file</title>
	</head>
		<body>
			<outline title="UserJS.org news" text="UserJS.org news" type="rss" version="RSS" xmlUrl="http://userjs.org/subscribe/news" htmlUrl="http://userjs.org/"/>
			<outline title="Olli's blog" text="Olli's blog" type="rss" version="RSS" xmlUrl="http://my.opera.com/olli/xml/atom/blog/" htmlUrl="http://my.opera.com/olli/"/>
			<outline title="Astronomy Picture of the Day" text="Astronomy Picture of the Day" type="rss" version="RSS" fix091="yes" xmlUrl="http://www.jwz.org/cheesegrater/RSS/apod.rss" htmlUrl="http://antwrp.gsfc.nasa.gov/apod/"/>
		</body>
	</opml>
	`

	var expected []RSS
	expected = append(expected, RSS{url: "http://userjs.org/subscribe/news", limit: 20})
	rssOpmlList := parseOPML([]byte(data))
	if len(rssOpmlList) != 3 {
		t.Fatalf("Wrong number of feeds: %d instead of %d", len(rssOpmlList), 3)
	}

	if rssOpmlList[0].url != expected[0].url {
		t.Errorf(`RSS Url is different: "%v" vs "%v"`, rssOpmlList[0], expected[0])
	}

}

func TestParseOpmlWithCategories(t *testing.T) {
	data := `<?xml version="1.0" encoding="UTF-8"?>
	<opml version="1.0">
			<head>
					<title>Subscription with Categories</title>
			</head>
			<body>
					<outline text="People" title="People">
							<outline type="rss" text="tonsky.me" title="tonsky.me" xmlUrl="http://tonsky.me/blog/atom.xml" htmlUrl="https://tonsky.me/"/>
							<outline type="rss" text="Armin Ronacher's Thoughts and Writings" title="Armin Ronacher's Thoughts and Writings" xmlUrl="http://lucumr.pocoo.org/feed.atom" htmlUrl="http://lucumr.pocoo.org/"/>
							<outline type="rss" text="Antirez weblog" title="Antirez weblog" xmlUrl="http://antirez.com/rss" htmlUrl="http://antirez.com"/>
					</outline>
					<outline text="Youtube" title="Youtube">
							<outline type="rss" text="LastWeekTonight" title="LastWeekTonight" xmlUrl="https://www.youtube.com/feeds/videos.xml?channel_id=UC3XTzVzaHQEd30rQbuvCtTQ" htmlUrl="https://www.youtube.com/channel/UC3XTzVzaHQEd30rQbuvCtTQ"/>
							<outline type="rss" text="Vsauce" title="Vsauce" xmlUrl="https://www.youtube.com/feeds/videos.xml?channel_id=UC6nSFpj9HTCZ5t-N3Rm3-HA" htmlUrl="https://www.youtube.com/channel/UC6nSFpj9HTCZ5t-N3Rm3-HA"/>
							<outline type="rss" text="suckerpinch" title="suckerpinch" xmlUrl="https://www.youtube.com/feeds/videos.xml?channel_id=UC3azLjQuz9s5qk76KEXaTvA" htmlUrl="https://www.youtube.com/channel/UC3azLjQuz9s5qk76KEXaTvA"/>
							<outline type="rss" text="Tom Scott" title="Tom Scott" xmlUrl="https://www.youtube.com/feeds/videos.xml?channel_id=UCBa659QWEk1AI4Tg--mrJ2A" htmlUrl="https://www.youtube.com/channel/UCBa659QWEk1AI4Tg--mrJ2A"/>
							<outline type="rss" text="Veritasium" title="Veritasium" xmlUrl="https://www.youtube.com/feeds/videos.xml?channel_id=UCHnyfMqiRRG1u-2MsSQLbXA" htmlUrl="https://www.youtube.com/channel/UCHnyfMqiRRG1u-2MsSQLbXA"/>
					</outline>
					<outline text="Tech / Main" title="Tech / Main">
							<outline type="rss" text="Wait But Why" title="Wait But Why" xmlUrl="http://waitbutwhy.com/feed" htmlUrl="https://waitbutwhy.com/"/>
							<outline type="rss" text="The Mad Ned Memo" title="The Mad Ned Memo" xmlUrl="https://madned.substack.com/feed/" htmlUrl="https://madned.substack.com"/>
							<outline type="rss" text="DHH – Signal v. Noise" title="DHH – Signal v. Noise" xmlUrl="https://m.signalvnoise.com/author/dhh/feed/" htmlUrl="https://m.signalvnoise.com"/>
							<outline type="rss" text="programming is terrible" title="programming is terrible" xmlUrl="http://programmingisterrible.com/rss" htmlUrl="https://programmingisterrible.com/"/>
							<outline type="rss" text="Coding Horror" title="Coding Horror" xmlUrl="http://feeds.feedburner.com/codinghorror/" htmlUrl="https://blog.codinghorror.com/"/>
							<outline type="rss" text="Astral Codex Ten" title="Astral Codex Ten" xmlUrl="https://astralcodexten.substack.com/feed/" htmlUrl="https://astralcodexten.substack.com"/>
							<outline type="rss" text="Hundred Rabbits" title="Hundred Rabbits" xmlUrl="https://100r.co/links/rss.xml" htmlUrl="https://100r.co"/>
					</outline>
			</body>
	</opml>
	`
	var expected []RSS
	expected = append(expected, RSS{url: "https://www.youtube.com/feeds/videos.xml?channel_id=UC3XTzVzaHQEd30rQbuvCtTQ", limit: 20})
	rssOpmlList := parseOPML([]byte(data))
	if len(rssOpmlList) != 15 {
		t.Fatalf("Wrong number of feeds: %d instead of %d", len(rssOpmlList), 15)
	}

	if rssOpmlList[3].url != expected[0].url {
		t.Errorf(`RSS Url is different: "%v" vs "%v"`, rssOpmlList[0], expected[0])
	}

}
