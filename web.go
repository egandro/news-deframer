package service

import (
	"context"

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
func (s *websrvc) Feed(ctx context.Context, p *web.FeedPayload) (res string, err error) {
	log.Printf(ctx, "web.feed")
	res = "feed " + p.FeedID
	return
}
