package rss

import(
  "encoding/xml"
  "encoding/json"
)

type RSS struct {
  Channel Channel `xml:"channel" json:"channel"`
}

type Channel struct {
	Title string `xml:"title" json:"title"`
	Link string `xml:"link" json:"link"`
	Description string `xml:"description" json:"description"`
	Language string `xml:"language" json:"language"`
  Copyright string `xml:"copyright" json:"copyright"`
	LastBuildDate string `xml:"lastBuildDate" json:"lastBuildDate"`
	Item []Item `xml:"item" json:"item"`
}

type ItemEnclosure struct {
  URL  string `xml:"url,attr" json:"url"`
  Type string `xml:"type,attr" json:"type"`
}

type Item struct {
  Title string `xml:"title" json:"title"`
  Link string  `xml:"link" json:"link"`
  Comments string `xml:"comments" json:"comments,omitempty"`
  PubDate string `xml:"pubDate" json:"pubDate"`
  GUID string `xml:"guid" json:"uid"`
  Category []string `xml:"category" json:"category"`
  Enclosure []ItemEnclosure `xml:"enclosure" json:"enclosure,omitempty"`
  Description string `xml:"description" json:"description,omitempty"`
  Author string `xml:"author" json:"author,omitempty"`
  // Content and namespacing is messed up .. needs work to use it`
}

func Unmarshal(b []byte, v interface{}) error {
  var feed RSS

  err := xml.Unmarshal(b, &feed)
  if err != nil {
    return err
  }

  jb, err := json.Marshal(feed);
  if err != nil {
    return err
  }

  err = json.Unmarshal(jb, &v)
  if err != nil {
    return err
  }

  return nil
}
