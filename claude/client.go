package claude

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	ApiKey  string
	baseUrl string
}

func NewClient(apiKey string) *Client {
	return &Client{
		ApiKey:  apiKey,
		baseUrl: "https://api.anthropic.com/v1",
	}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []Message `json:"messages"`
}

type Response struct {
	Content []ContentBlock `json:"content"`
}

type ContentBlock struct {
	Text string `json:"text"`
}

func (c *Client) Complete(prompt string) (string, error) {
	return c.Chat([]Message{{Role: "user", Content: prompt}})
}

func (c *Client) Chat(messages []Message) (string, error) {
	req := Request{
		Model:     "claude-haiku-4-5-20251001",
		MaxTokens: 1024,
		Messages:  messages,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest("POST", c.baseUrl+"/messages", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	request.Header.Set("x-api-key", c.ApiKey)
	request.Header.Set("anthropic-version", "2023-06-01")
	request.Header.Set("content-type", "application/json")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error %d: %s", resp.StatusCode, body)
	}

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Content) == 0 {
		return "", fmt.Errorf("empty response from API")
	}
	return result.Content[0].Text, nil
}
