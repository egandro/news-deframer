package source

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseString(t *testing.T) {
	feedJSON := `
	[
		{
			"rss_url": "https://example.com/rss",
			"language": "en"
		},
		{
			"rss_url": "https://example.com/rss2",
			"language": "fr"
		}
	]
`

	feeds, err := ParseString(feedJSON)
	assert.NoError(t, err)
	assert.Len(t, feeds, 2)

	feed1 := feeds[0]
	assert.EqualValues(t, "https://example.com/rss", feed1.RSSURL)
	assert.EqualValues(t, "en", feed1.Language)

	feed2 := feeds[1]
	assert.EqualValues(t, "https://example.com/rss2", feed2.RSSURL)
	assert.EqualValues(t, "fr", feed2.Language)
}
