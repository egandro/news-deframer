package deframer

import (
	"fmt"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/egandro/news-deframer/pkg/database"
	"github.com/egandro/news-deframer/pkg/downloader"
	"github.com/egandro/news-deframer/pkg/openai"
	"github.com/egandro/news-deframer/pkg/source"
	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
)

const maxAge = time.Minute * 90

type deframer struct {
	db         *database.Database
	ai         openai.OpenAI
	src        *source.Source
	downloader downloader.Downloader
	prompts    map[string]source.Prompt
}

type Deframer interface {
	UpdateFeeds() (int, error)
	Deframe(feed *gofeed.Feed) (string, error)
}

// NewDeframer initializes a new deframer
func NewDeframer() (Deframer, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	db, err := database.NewDatabase(cfg.DatabaseFile)
	if err != nil {
		return nil, err
	}

	ai := openai.NewAI(cfg.AI_URL, cfg.AI_Model, "")

	src, err := source.ParseFile(cfg.Source)
	if err != nil {
		return nil, err
	}

	downloader := downloader.NewDownloader()

	prompts := make(map[string]source.Prompt)
	for _, prompt := range src.Prompts {
		prompts[prompt.Language] = prompt
	}

	res := &deframer{
		db:         db,
		ai:         ai,
		src:        src,
		downloader: downloader,
		prompts:    prompts,
	}

	return res, nil
}

func (d *deframer) UpdateFeeds() (int, error) {
	numberOfDownloads := 0
	parser := gofeed.NewParser()

	for _, feed := range d.src.Feeds {
		cache, err := d.db.FindCacheByFeedUrl(feed.RSS_URL, maxAge)
		if err != nil {
			return numberOfDownloads, err
		}

		if cache != nil {
			continue
		}

		data, err := d.downloader.DownloadRSSFeed(feed.RSS_URL)
		if err != nil {
			return numberOfDownloads, err
		}

		parsedData, err := parser.ParseString(string(data))
		if err != nil {
			return numberOfDownloads, err
		}

		title := parsedData.Title
		if title == "" {
			// some fallback
			title = feed.RSS_URL
		}

		title = fmt.Sprintf("%v (%v)", title, feed.Language)

		unframed, err := d.Deframe(parsedData)
		if err != nil {
			return numberOfDownloads, err
		}

		cache = &database.Cache{
			FeedUrl: feed.RSS_URL,
			Title:   title,
			Cache:   unframed,
		}

		err = d.db.CreateCache(cache)
		if err != nil {
			return numberOfDownloads, err
		}

		numberOfDownloads++
	}
	return numberOfDownloads, nil
}

func (d *deframer) Deframe(feed *gofeed.Feed) (string, error) {
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

	result, err := newFeed.ToRss()
	if err != nil {
		return "", err
	}

	return result, nil
}
