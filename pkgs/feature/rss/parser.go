package rss

import (
	"encoding/xml"
	"time"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel *Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Language    string `xml:"language"`
	PubDate     string `xml:"pubDate"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
	Author      string `xml:"author"`
	Category    string `xml:"category"`
}

func Parse(data []byte) (*RSS, error) {
	var rss RSS
	err := xml.Unmarshal(data, &rss)
	if err != nil {
		return nil, err
	}
	return &rss, nil
}

func (i *Item) ParsePubDate() (time.Time, error) {
	layouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"Mon, 02 Jan 2006 15:04:05 MST",
		time.RFC3339,
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, i.PubDate); err == nil {
			return t, nil
		}
	}

	return time.Time{}, nil
}
