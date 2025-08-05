package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type openAI struct {
	client *openai.Client
	model  string
}

// OpenAI handles AI operation
type OpenAI interface {
	Query(ctx context.Context, user string, system string) (string, error)
	FuzzyParseJSON(input string) (interface{}, error)
}

// NewAI
func NewAI(url string, model string, token string) OpenAI {
	res := &openAI{
		model: model,
	}

	if token != "" {
		// openai
		res.client = openai.NewClient(token)
	} else {
		// LM Studio or similar
		config := openai.DefaultConfig("")
		config.BaseURL = url
		// Optional: use a custom HTTP client (e.g., no TLS verification)
		config.HTTPClient = &http.Client{}
		res.client = openai.NewClientWithConfig(config)
	}

	return res
}

func (a *openAI) Query(ctx context.Context, user string, system string) (string, error) {
	resp, err := a.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: a.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: system,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: user,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func (a *openAI) FuzzyParseJSON(input string) (interface{}, error) {
	cleaned := cleanInput(input)

	// Step 1: Try to isolate the part where JSON begins
	idx := strings.IndexAny(cleaned, "{[\"tfn0123456789") // JSON starts with one of these
	if idx == -1 {
		return nil, fmt.Errorf("no JSON start found")
	}

	jsonCandidate := strings.TrimSpace(cleaned[idx:])

	// Step 2: Use a decoder to extract a single JSON value
	dec := json.NewDecoder(bytes.NewReader([]byte(jsonCandidate)))
	dec.UseNumber() // keeps numbers flexible

	var result interface{}
	if err := dec.Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w", err)
	}

	return result, nil
}

func cleanInput(input string) string {
	input = strings.TrimSpace(input)

	// Remove markdown ```json or ``` code blocks
	input = strings.TrimPrefix(input, "```json")
	input = strings.TrimPrefix(input, "```")
	input = strings.TrimSuffix(input, "```")

	input = strings.TrimSpace(input)

	// Match and remove "json" or "JSON" at the beginning, even if followed directly by JSON
	re := regexp.MustCompile(`(?i)^json\s*`)
	input = re.ReplaceAllString(input, "")

	return strings.TrimSpace(input)
}
