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

var week bool
var format string
var summaryTag string

var SummaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Summary of work entries for the week",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		since := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
		entries, err := store.GetEntriesSince(since)
		if err != nil {
			return err
		}
		if summaryTag != "" {
			entries = store.FilterByTag(entries, summaryTag)
		}
		if len(entries) == 0 {
			fmt.Println("No entries were found.")
			return nil
		}

		entriesText := buildEntriesText(entries)

		client := claude.NewClient(cfg.Anthropic.APIKey)

		var prompt string
		if format == "promo" {
			prompt = fmt.Sprintf(prompts.Promo, entriesText)
		} else {
			prompt = fmt.Sprintf(prompts.Summary, entriesText)
		}

		resp, err := client.Complete(prompt)
		if err != nil {
			return err
		}

		fmt.Println(resp)
		return nil
	},
}

func init() {
	SummaryCmd.Flags().BoolVar(&week, "week", false, "Show summary for the week")
	SummaryCmd.Flags().StringVarP(&format, "format", "f", "", "Output format")
	SummaryCmd.Flags().StringVarP(&summaryTag, "tag", "t", "", "Filter entries by tag")
}
