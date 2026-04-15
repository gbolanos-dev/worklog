package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gbolanos-dev/worklog/claude"
	"github.com/gbolanos-dev/worklog/config"
	"github.com/gbolanos-dev/worklog/fetch"
	"github.com/gbolanos-dev/worklog/prompts"
	"github.com/gbolanos-dev/worklog/store"
	"github.com/spf13/cobra"
)

var issues []string
var prs []string
var standupTag string

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
		if standupTag != "" {
			entries, err = store.GetEntriesByTag(standupTag)
			if err != nil {
				return err
			}
		}
		if len(entries) == 0 {
			fmt.Println("No entries were found.")
			return nil
		}

		entriesText := buildEntriesText(entries)

		ticketsText, err := buildTicketsText(cfg, issues)
		if err != nil {
			return err
		}

		prsText, err := buildPRsText(cfg, prs)
		if err != nil {
			return err
		}

		prompt := fmt.Sprintf(prompts.Standup, entriesText, ticketsText, prsText)

		client := claude.NewClient(cfg.Anthropic.APIKey)
		resp, err := client.Complete(prompt)
		if err != nil {
			return err
		}
		fmt.Println(resp)
		return nil
	},
}

func init() {
	StandupCmd.Flags().StringArrayVarP(&issues, "issue", "i", nil, "YouTrack issue IDs")
	StandupCmd.Flags().StringArrayVarP(&prs, "pr", "p", nil, "GitHub PR number")
	StandupCmd.Flags().StringVarP(&standupTag, "tag", "t", "", "Filter entries by tag")
}

func buildEntriesText(entries []store.Entry) string {
	var s string
	for i, entry := range entries {
		s += fmt.Sprintf("%d. %s\n", i+1, entry.Entry)
	}
	return s
}

func buildTicketsText(cfg *config.Config, ticketIDs []string) (string, error) {
	var s string
	for _, id := range ticketIDs {
		ticket, err := fetch.FetchTicket(cfg.YouTrack.BaseURL, cfg.YouTrack.Token, id)
		if err != nil {
			return "", err
		}
		s += fmt.Sprintf("Ticket: %s - %s\nDescription: %s\n", id, ticket.Summary, ticket.Description)
		for _, c := range ticket.Comments {
			s += fmt.Sprintf("  Comment (%s): %s\n", c.Author.Login, c.Text)
		}
		s += "\n"
	}
	return s, nil
}

func buildPRsText(cfg *config.Config, prNums []string) (string, error) {
	var s string
	parts := strings.Split(cfg.GitHub.DefaultRepo, "/")
	for _, prNum := range prNums {
		num, err := strconv.Atoi(prNum)
		if err != nil {
			return "", err
		}
		pr, err := fetch.FetchPR(cfg.GitHub.Token, parts[0], parts[1], num)
		if err != nil {
			return "", err
		}
		s += fmt.Sprintf("PR #%d: %s\n%s\nFiles: %s\n\n", num, pr.Title, pr.Body, strings.Join(pr.FilesChanged, ", "))
	}
	return s, nil
}
