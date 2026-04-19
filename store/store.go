package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Entry struct {
	ID    string   `json:"id"`
	Date  string   `json:"date"`
	Entry string   `json:"entry"`
	Tags  []string `json:"tags"`
}

func AddEntry(text string, tags []string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dir := filepath.Join(home, ".worklog")
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	entries, err := loadEntries(dir)
	if err != nil {
		return err
	}

	entry := Entry{
		ID:    fmt.Sprintf("%x", time.Now().UnixNano()),
		Date:  time.Now().Format("2006-01-02"),
		Entry: text,
		Tags:  tags,
	}

	entries = append(entries, entry)

	data, err := json.MarshalIndent(entries, "", "")
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(dir, "log.json"), data, 0666)
	if err != nil {
		return err
	}
	return nil
}

func GetEntriesForDate(date string) ([]Entry, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(home, ".worklog")

	entries, err := loadEntries(dir)
	if err != nil {
		return nil, err
	}

	var filtered []Entry
	for _, entry := range entries {
		if entry.Date == date {
			filtered = append(filtered, entry)
		}
	}
	return filtered, nil
}

func GetEntriesSince(since string) ([]Entry, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(home, ".worklog")

	entries, err := loadEntries(dir)
	if err != nil {
		return nil, err
	}

	cutoff, err := time.Parse("2006-01-02", since)
	if err != nil {
		return nil, err
	}

	var filtered []Entry
	for _, entry := range entries {
		entryDate, err := time.Parse("2006-01-02", entry.Date)
		if err != nil {
			return nil, err
		}
		if !entryDate.Before(cutoff) {
			filtered = append(filtered, entry)
		}
	}
	return filtered, nil
}

func FilterUntil(entries []Entry, until string) ([]Entry, error) {
	cutoff, err := time.Parse("2006-01-02", until)
	if err != nil {
		return nil, err
	}

	var filtered []Entry
	for _, entry := range entries {
		entryDate, err := time.Parse("2006-01-02", entry.Date)
		if err != nil {
			return nil, err
		}
		if !entryDate.After(cutoff) {
			filtered = append(filtered, entry)
		}
	}
	return filtered, nil
}

func FilterByTag(entries []Entry, tag string) []Entry {
	var filtered []Entry
	for _, entry := range entries {
		for _, t := range entry.Tags {
			if t == tag {
				filtered = append(filtered, entry)
				break
			}
		}
	}
	return filtered
}

func loadEntries(dir string) ([]Entry, error) {
	data, err := os.ReadFile(filepath.Join(dir, "log.json"))
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	var entries []Entry

	if len(data) > 0 {
		err = json.Unmarshal(data, &entries)
		if err != nil {
			return nil, err
		}
	}

	return entries, nil
}
