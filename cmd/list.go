package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gbolanos-dev/worklog/internal/dateutil"
	"github.com/gbolanos-dev/worklog/store"
	"github.com/spf13/cobra"
)

var date string
var since string
var until string
var tag string
var jsonOutput bool

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all work entries",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate: --until requires --since
		if until != "" && since == "" {
			return fmt.Errorf("--until requires --since to also be set")
		}

		var entries []store.Entry
		var err error

		if since != "" {
			parsedSince, err := dateutil.Parse(since)
			if err != nil {
				return err
			}

			if until != "" {
				parsedUntil, err := dateutil.Parse(until)
				if err != nil {
					return err
				}

				// Validate: --until must be on or after --since
				sinceTime, _ := time.Parse("2006-01-02", parsedSince)
				untilTime, _ := time.Parse("2006-01-02", parsedUntil)
				if untilTime.Before(sinceTime) {
					return fmt.Errorf("--until (%s) must be on or after --since (%s)", parsedUntil, parsedSince)
				}

				entries, err = store.GetEntriesSince(parsedSince)
				if err != nil {
					return err
				}
				entries, err = store.FilterUntil(entries, parsedUntil)
				if err != nil {
					return err
				}
			} else {
				entries, err = store.GetEntriesSince(parsedSince)
				if err != nil {
					return err
				}
			}
		} else {
			targetDate := time.Now().Format("2006-01-02")
			if date != "" {
				targetDate, err = dateutil.Parse(date)
				if err != nil {
					return err
				}
			}
			entries, err = store.GetEntriesForDate(targetDate)
			if err != nil {
				return err
			}
		}

		if tag != "" {
			entries = store.FilterByTag(entries, tag)
		}

		if jsonOutput {
			out, err := json.MarshalIndent(entries, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(out))
			return nil
		}

		for i, entry := range entries {
			tags := ""
			if len(entry.Tags) > 0 {
				tags = " [" + strings.Join(entry.Tags, ", ") + "]"
			}
			fmt.Printf("%d. [%s] %s%s\n", i+1, entry.ID[:8], entry.Entry, tags)
		}
		return nil
	},
}

func init() {
	ListCmd.Flags().StringVarP(&date, "date", "d", "", "Filter entries by date")
	ListCmd.Flags().StringVarP(&since, "since", "s", "", "Filter entries from date (inclusive)")
	ListCmd.Flags().StringVar(&until, "until", "", "Filter entries until date (inclusive)")
	ListCmd.Flags().StringVarP(&tag, "tag", "t", "", "Filter entries by tag")
	ListCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output entries in JSON format")
	ListCmd.MarkFlagsMutuallyExclusive("date", "since")
	ListCmd.MarkFlagsMutuallyExclusive("date", "until")
}
