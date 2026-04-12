package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type YouTrackTicket struct {
	ID          string      `json:"id"`
	Summary     string      `json:"summary"`
	Description string      `json:"description"`
	Comments    []ytComment `json:"comments"`
}

type ytComment struct {
	Text   string   `json:"text"`
	Author ytAuthor `json:"author"`
}

type ytAuthor struct {
	Login string `json:"login"`
}

func FetchTicket(baseURL, token, ticketID string) (*YouTrackTicket, error) {
	url := fmt.Sprintf(
		"%s/api/issues/%s?fields=id,summary,description,comments(text,author(login))",
		baseURL,
		ticketID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("youtrack: %d %s", resp.StatusCode, body)
	}

	var ticket YouTrackTicket
	if err := json.NewDecoder(resp.Body).Decode(&ticket); err != nil {
		return nil, err
	}

	return &ticket, nil
}
