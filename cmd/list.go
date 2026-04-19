package cmd

import (
	"fmt"
	"time"

	"github.com/gbolanos-dev/worklog/internal/dateutil"
	"github.com/gbolanos-dev/worklog/store"
	"github.com/spf13/cobra"
)

var date string
var tag string

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all work entries",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		targetDate := time.Now().Format("2006-01-02")
		if date != "" {
			parsed, err := dateutil.Parse(date)
			if err != nil {
				return err
			}
			targetDate = parsed
		}
		entries, err := store.GetEntriesForDate(targetDate)
		if err != nil {
			return err
		}
		if tag != "" {
			entries = store.FilterByTag(entries, tag)
		}

		for i, entry := range entries {
			fmt.Printf("%d. %s\n", i+1, entry.Entry)
		}
		return nil
	},
}

func init() {
	ListCmd.Flags().StringVarP(&date, "date", "d", "", "Filter entries by date")
	ListCmd.Flags().StringVarP(&tag, "tag", "t", "", "Filter entries by tag")
}
