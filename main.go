package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Author struct {
	Name string `xml:"name"`
}

type Link struct {
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
	Href string `xml:"href,attr"`
}

type Entry struct {
	ID        string    `xml:"id"`
	Published time.Time `xml:"published"`
	Updated   time.Time `xml:"updated"`
	Author    Author    `xml:"author"`
	Content   string    `xml:"content"`
	Link      []Link    `xml:"link"`
	Summary   string    `xml:"summary"`
}

type Atom struct {
	XMLName  xml.Name  `xml:"feed"`
	ID       string    `xml:"id"`
	Title    string    `xml:"title"`
	Subtitle string    `xml:"subtitle"`
	Updated  time.Time `xml:"updated"`
	Author   Author    `xml:"author"`
	Link     []Link    `xml:"link"`
	Entries  []Entry   `xml:"entry"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
}
type RSS struct {
	XMLName     xml.Name `xml:"rss"`
	Title       string   `xml:"title"`
	Link        []Link   `xml:"link"`
	Description string   `xml:"description"`
	Item        []Item   `xml:"item"`
}

func main() {
	// RSS & Atom feeds example
	// using docker related feeds for example
	sources := []string{
		// RSS
		"http://blog.docker.com/feed/",
		"https://serversforhackers.com/feed",
		"http://dtrace.org/blogs/wesolows/feed/",

		// Atom
		"http://googlecloudplatform.blogspot.com/feeds/posts/default",
		// "http://www.goinggo.net/feeds/posts/default",
	}
	for _, feed := range sources {
		// fetch feed
		fmt.Printf("Feed: %v\n", feed)
		resp, err := http.Get(feed)
		if err != nil {
			fmt.Printf("Error fetch feed: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		// TODO: we can use Last-Modified or Etag response header
		// to know the feed updated or not
		// TODO: benchmark content pake string dan byte saat encode
		// TODO: remove html tag on content

		fmt.Printf("media type: %v\n", resp.Header.Get("Content-Type"))
		contentType := resp.Header.Get("Content-Type")
		// bagaimana kalo tau itu atom?
		s := strings.Split(contentType, ";")
		if s[0] == "application/atom+xml" {
			fmt.Printf("\n\nFeed: %v jenis ATOM maka decode\n", feed)
			// TODO: what is a factor that make decoder so slow? last measured is 10.948143721s
			start := time.Now()
			var atom Atom
			dec := xml.NewDecoder(resp.Body)
			err := dec.Decode(&atom)
			if err != nil {
				fmt.Printf("Error saat decode %v\n", err)
				return
			}
			fmt.Printf("Encode %v\n", time.Since(start))
			fmt.Printf("%#v\n", atom)
			// fmt.Printf("Atom name : %s\n", atom.XMLName)
			// fmt.Printf("Atom ID : %s\n", atom.ID)
			// fmt.Printf("Atom Title : %s\n", atom.Title)
			// fmt.Printf("Atom Subtitle : %s\n", atom.Subtitle)
			// fmt.Printf("Atom Updated : %v\n", atom.Updated)
			// fmt.Printf("Atom Author : %#v\n", atom.Author)
			// fmt.Printf("Atom Link : %#v\n", atom.Link)
			// fmt.Println("Atom entries:")
			// for _, entry := range atom.Entries {
			// 	fmt.Printf("Author : %#v\n", entry.Author)
			// 	fmt.Printf("Published : %v\n", entry.Published)
			// 	fmt.Printf("Updated : %v\n", entry.Updated)
			// 	fmt.Printf("Content : %s\n", entry.Content)
			// 	fmt.Printf("Link : %s\n", entry.Link)
			// 	fmt.Printf("Summary : %s\n\n", entry.Summary)
			// }
		}

		if s[0] == "text/xml" {
			fmt.Printf("%v type RSS\n", feed)
			start := time.Now()
			var rss RSS
			dec := xml.NewDecoder(resp.Body)
			err := dec.Decode(&rss)
			if err != nil {
				fmt.Printf("Error saat decode %v\n", err)
				continue
			}
			fmt.Printf("Encode %v\n", time.Since(start))
			fmt.Printf("%#v\n", rss)
		}

	}
}
