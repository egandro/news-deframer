package service

import (
	"context"
	"fmt"
	"io"
	"strings"
	"text/template"

	web "github.com/egandro/news-deframer/gen/web"
	"github.com/egandro/news-deframer/pkg/deframer"
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

	const tpl = `
<!DOCTYPE html>
<html>
<head><title>{{.Title}}</title></head>
<body>
	<h1>{{.Heading}}</h1>
	{{if .Items}}
		<ul>
			{{range .Items}}
				<li><a href="{{.Href}}">{{.Title}}</a></li>
			{{end}}
		</ul>
	{{else}}
		<p>No feeds available.</p>
	{{end}}
</body>
</html>`

	d, err := deframer.NewDeframer(ctx)
	if err != nil {
		return "", err
	}

	type Item struct {
		Href  string
		Title string
	}

	items := []Item{}

	caches, err := d.FindAllCaches()
	if err != nil {
		return "", err
	}

	for _, cache := range caches {
		items = append(items, Item{
			Href:  fmt.Sprintf("/feed/%v", cache.ID),
			Title: cache.Title,
		})
	}

	data := map[string]any{
		"Title":   "Deframer",
		"Heading": "Deframed RSS Feeds",
		"Items":   items,
	}

	res, err = renderTemplate(tpl, data)

	return
}

// Returns the feed with the given xml
func (s *websrvc) Feed(ctx context.Context, p *web.FeedPayload) (res *web.FeedResult, resp io.ReadCloser, err error) {
	res = &web.FeedResult{}
	log.Printf(ctx, "web.feed")
	//res = "feed " + p.FeedID

	d, err := deframer.NewDeframer(ctx)
	if err != nil {
		return res, resp, err
	}

	entry, err := d.FindCacheByID(p.FeedID)
	if err != nil {
		return res, resp, err
	}

	res.Type = "application/xml;charset=UTF-8"
	res.Length = int64(len(entry.Cache))

	// resp is the HTTP response body stream.
	resp = io.NopCloser(strings.NewReader(entry.Cache))

	return
}

// renderTemplate takes an template string and some data,
// and returns the rendered template as a string.
func renderTemplate(tpl string, data any) (string, error) {
	// Parse the template string
	t, err := template.New("page").Parse(tpl)
	if err != nil {
		return "", err
	}

	// Execute into a string builder
	var sb strings.Builder
	if err := t.Execute(&sb, data); err != nil {
		return "", err
	}

	return sb.String(), nil
}
