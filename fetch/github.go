package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GitHubPR struct {
	Number       int    `json:"number"`
	Title        string `json:"title"`
	Body         string `json:"body"`
	FilesChanged []string
}

type ghFile struct {
	Filename string `json:"filename"`
}

func FetchPR(token, owner, repo string, number int) (*GitHubPR, error) {
	base := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d", owner, repo, number)

	var pr GitHubPR
	if err := ghGet(token, base, &pr); err != nil {
		return nil, err
	}

	var files []ghFile
	if err := ghGet(token, base+"/files", &files); err != nil {
		return nil, err
	}

	for _, f := range files {
		pr.FilesChanged = append(pr.FilesChanged, f.Filename)
	}

	return &pr, nil
}

func ghGet(token, url string, target any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("github: %d %s", resp.StatusCode, body)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}
