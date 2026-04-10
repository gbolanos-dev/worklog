package update

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Result struct {
	Available bool
	Latest    string
}

type cache struct {
	CheckedAt time.Time `json:"checked_at"`
	Latest    string    `json:"latest"`
}

type response struct {
	TagName string `json:"tag_name"`
}

func Check(currentVersion, cacheDir string) *Result {
	if strings.EqualFold(currentVersion, "dev") {
		return nil
	}

	data, err := os.ReadFile(filepath.Join(cacheDir, "update-check.json"))
	if err != nil {
		return nil
	}

	var c cache
	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil
	}

	if time.Since(c.CheckedAt) > time.Hour*24 {
		return nil
	}

	if newer(currentVersion, c.Latest) {
		return &Result{
			Available: true,
			Latest:    c.Latest,
		}
	}

	return &Result{
		Available: false,
		Latest:    c.Latest,
	}
}

func Refresh(cacheDir string) {
	client := &http.Client{
		Timeout: time.Second * 3,
	}

	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/repos/gbolanos-dev/worklog/releases/latest",
		nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	var result response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return
	}

	c := cache{
		CheckedAt: time.Now().UTC(),
		Latest:    result.TagName,
	}

	data, err := json.Marshal(c)
	if err != nil {
		return
	}

	_ = os.WriteFile(filepath.Join(cacheDir, "update-check.json"), data, 0666)
}

func DoUpdate(module string) error {
	cmd := exec.Command("go", "install", module+"@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
