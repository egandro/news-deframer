package source

import (
	"encoding/json"
	"os"
)

type Feed struct {
	RSSURL   string `json:"rss_url"`
	Language string `json:"language"` // e.g. 'en' for English
}

// ParseString parses the feed from a JSON string and returns feeds
func ParseString(feedJSON string) ([]Feed, error) {
	var feeds []Feed
	err := json.Unmarshal([]byte(feedJSON), &feeds)
	return feeds, err
}

// ParseFile parses the feed from a JSON file and returns feeds
func ParseFile(filePath string) ([]Feed, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return ParseString(string(data))
}
