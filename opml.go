package main

import (
	"encoding/xml"
)

type Opml struct {
	XMLName xml.Name `xml:"opml"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Head    struct {
		Text  string `xml:",chardata"`
		Title string `xml:"title"`
	} `xml:"head"`
	Body struct {
		Text    string `xml:",chardata"`
		Outline []struct {
			Text     string `xml:",chardata"`
			AttrText string `xml:"text,attr"`
			Title    string `xml:"title,attr"`
			Type     string `xml:"type,attr"`
			XmlUrl   string `xml:"xmlUrl,attr"`
			HtmlUrl  string `xml:"htmlUrl,attr"`
			Outline  []struct {
				Text     string `xml:",chardata"`
				Type     string `xml:"type,attr"`
				AttrText string `xml:"text,attr"`
				Title    string `xml:"title,attr"`
				XmlUrl   string `xml:"xmlUrl,attr"`
				HtmlUrl  string `xml:"htmlUrl,attr"`
			} `xml:"outline"`
		} `xml:"outline"`
	} `xml:"body"`
}
