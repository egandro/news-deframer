package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	web "github.com/egandro/news-deframer/gen/web"
	"goa.design/clue/log"
)

// web service example implementation.
// The example methods log the requests and return zero values.
type websrvc struct{}

// NewWeb returns the web service implementation.
func NewWeb() web.Service {
	return &websrvc{}
}

// Returns the index page in HTML
func (s *websrvc) Index(ctx context.Context) (res string, err error) {
	log.Printf(ctx, "web.index")
	res = "hello world"
	return
}

// Returns the feed with the given xml
func (s *websrvc) Feed(ctx context.Context, p *web.FeedPayload) (res *web.FeedResult, resp io.ReadCloser, err error) {
	res = &web.FeedResult{}
	log.Printf(ctx, "web.feed")
	//res = "feed " + p.FeedID

	url := "https://www.tagesschau.de/index~rss2.xml"

	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("HTTP request failed with status %d", response.StatusCode))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	rssXML := string(body)

	// res.Type = "application/rss+xml"
	res.Type = "application/xml;charset=UTF-8"
	res.Length = int64(len(rssXML))

	// resp is the HTTP response body stream.
	resp = io.NopCloser(strings.NewReader(rssXML))

	return
}
