package main

import (
	"log"
	"os"
	"time"

	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
)

// https://www.machinet.net/tutorial-de/create-rss-feed-generator-go

func main() {
	// Read feed XML
	data, err := os.ReadFile("input.xml")
	if err != nil {
		log.Fatal(err)
	}

	parser := gofeed.NewParser()
	feed, err := parser.ParseString(string(data))
	if err != nil {
		log.Fatal(err)
	}

	// Update channel title with prefix
	prefix := "[Prefix] "
	feed.Title = prefix + feed.Title

	newFeed := &feeds.Feed{
		Title: feed.Title,
		Link: &feeds.Link{
			Href: feed.Link,
			Type: feed.FeedType,
		},
		Description: feed.Description,
		Author:      &feeds.Author{},
		// // Id:
		// // Subtitle:
		// // Items:
		Copyright: feed.Copyright,
		// Image:
	}

	if feed.Author != nil {
		newFeed.Author.Name = feed.Author.Name
		newFeed.Author.Email = feed.Author.Email
	}

	if feed.PublishedParsed != nil {
		newFeed.Created = *feed.PublishedParsed
	}

	if feed.UpdatedParsed != nil {
		newFeed.Updated = *feed.UpdatedParsed
	}

	if feed.Image != nil {
		newFeed.Image = &feeds.Image{}
		newFeed.Image.Url = feed.Image.URL
		newFeed.Image.Title = feed.Image.Title
	}

	item := &feeds.Item{
		Title:       "Ihr Post Titel",
		Link:        &feeds.Link{Href: "http://example.com/post-url"},
		Description: "Eine kurze Beschreibung zu Ihrem Post",
		Author:      &feeds.Author{Name: "Your Name", Email: "yourname@example.com"},
		Created:     time.Now(),
		Id:          "my id",
	}
	newFeed.Add(item)

	for _, current := range feed.Items {
		item := &feeds.Item{
			Title:       "T: " + current.Title,
			Link:        &feeds.Link{Href: current.Link},
			Description: "D: " + current.Description,
			Content:     "C: " + current.Content,
			Id:          current.GUID,
		}

		if current.PublishedParsed != nil {
			item.Created = *current.PublishedParsed
		}

		if current.UpdatedParsed != nil {
			item.Updated = *current.UpdatedParsed
		}

		if len(current.Authors) > 0 {
			item.Author = &feeds.Author{
				Name:  current.Authors[0].Name,
				Email: current.Authors[0].Email,
			}
		}

		newFeed.Add(item)
	}

	rss, err := newFeed.ToRss()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("output.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.WriteString(rss)
	if err != nil {
		log.Fatal(err)
	}
}
