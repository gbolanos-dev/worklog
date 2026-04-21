package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Entry struct {
	ID    string   `json:"id"`
	Date  string   `json:"date"`
	Entry string   `json:"entry"`
	Tags  []string `json:"tags"`
}

// loadAll returns all entries and the worklog directory path.
func loadAll() ([]Entry, string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, "", err
	}

	dir := filepath.Join(home, ".worklog")
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, "", err
	}

	entries, err := loadEntries(dir)
	if err != nil {
		return nil, "", err
	}

	return entries, dir, nil
}

// saveEntries marshals and writes entries back to log.json.
func saveEntries(dir string, entries []Entry) error {
	data, err := json.MarshalIndent(entries, "", "")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "log.json"), data, 0666)
}

// findByPrefix returns the index of the entry matching the ID prefix.
// Returns an error if no match or ambiguous match.
func findByPrefix(entries []Entry, id string) (int, error) {
	var matches []int
	for i, entry := range entries {
		if strings.HasPrefix(entry.ID, id) {
			matches = append(matches, i)
		}
	}

	if len(matches) == 0 {
		return -1, fmt.Errorf("no entry found with ID prefix %q", id)
	}
	if len(matches) > 1 {
		return -1, fmt.Errorf("ambiguous ID prefix %q matches %d entries, use a longer prefix", id, len(matches))
	}

	return matches[0], nil
}

func AddEntry(text string, tags []string) error {
	entries, dir, err := loadAll()
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
	return saveEntries(dir, entries)
}

func GetEntriesForDate(date string) ([]Entry, error) {
	entries, _, err := loadAll()
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
	entries, _, err := loadAll()
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

func FindEntryById(id string) (*Entry, error) {
	entries, _, err := loadAll()
	if err != nil {
		return nil, err
	}

	idx, err := findByPrefix(entries, id)
	if err != nil {
		return nil, err
	}

	return &entries[idx], nil
}

func DeleteEntry(id string) error {
	entries, dir, err := loadAll()
	if err != nil {
		return err
	}

	idx, err := findByPrefix(entries, id)
	if err != nil {
		return err
	}

	entries = append(entries[:idx], entries[idx+1:]...)
	return saveEntries(dir, entries)
}

func EditEntry(id, newText string) error {
	entries, dir, err := loadAll()
	if err != nil {
		return err
	}

	idx, err := findByPrefix(entries, id)
	if err != nil {
		return err
	}

	entries[idx].Entry = newText
	return saveEntries(dir, entries)
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
