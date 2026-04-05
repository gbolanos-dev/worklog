package cmd

import (
	"fmt"
	"time"

	"github.com/gbolanos-dev/worklog/store"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all work entries",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		today := time.Now().Format("2006-01-02")
		entries, err := store.GetEntriesForDate(today)
		if err != nil {
			return err
		}

		for i, entry := range entries {
			fmt.Printf("%d. %s\n", i+1, entry.Entry)
		}
		return nil
	},
}
