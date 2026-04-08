package cmd

import (
	"fmt"
	"time"

	"github.com/gbolanos-dev/worklog/claude"
	"github.com/gbolanos-dev/worklog/config"
	"github.com/gbolanos-dev/worklog/prompts"
	"github.com/gbolanos-dev/worklog/store"
	"github.com/spf13/cobra"
)

var StandupCmd = &cobra.Command{
	Use:   "standup",
	Short: "Standup (Y/T/B)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		today := time.Now().Format("2006-01-02")
		entries, err := store.GetEntriesForDate(today)
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			fmt.Println("No entries were found.")
			return nil
		}

		client := claude.NewClient(cfg.Anthropic.APIKey)

		var entriesText string
		for i, entry := range entries {
			entriesText += fmt.Sprintf("%d. %s\n", i+1, entry.Entry)
		}

		prompt := fmt.Sprintf(prompts.Standup, entriesText)

		resp, err := client.Complete(prompt)
		if err != nil {
			return err
		}
		fmt.Println(resp)
		return nil
	},
}
