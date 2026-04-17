package cmd

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/gbolanos-dev/worklog/claude"
	"github.com/gbolanos-dev/worklog/config"
	"github.com/gbolanos-dev/worklog/prompts"
	"github.com/gbolanos-dev/worklog/store"
	"github.com/spf13/cobra"
)

var chatFrom string

var chatIssues []string
var chatPRs []string

var ChatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Chat with Claude",
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

		entriesText := buildEntriesText(entries)

		ticketsText, err := buildTicketsText(cfg, chatIssues)
		if err != nil {
			return err
		}

		prsText, err := buildPRsText(cfg, chatPRs)
		if err != nil {
			return err
		}

		context := fmt.Sprintf(
			"%s\n\nHere is my work context:\n\nEntries:\n%s\nTickets:\n%s\nPull Requests:\n%s",
			prompts.Chat,
			entriesText,
			ticketsText,
			prsText)

		// Start history with chat instructions, work context, and Claude's acknowledgment.
		history := []claude.Message{
			{Role: "user", Content: context},
			{Role: "assistant", Content: "I have your work context loaded. How can I help?"},
		}

		if chatFrom != "" {
			if chatFrom != "standup" {
				return fmt.Errorf("unsupported --from value %q (valid options: standup)", chatFrom)
			}

			standup, err := generateStandup(cfg, entries, chatIssues, chatPRs)
			if err != nil {
				return err
			}

			history = append(history, claude.Message{Role: "assistant", Content: standup})
		}

		client := claude.NewClient(cfg.Anthropic.APIKey)
		scanner := bufio.NewScanner(os.Stdin)
		if chatFrom == "standup" {
			fmt.Println("Chat started from generated standup. Type 'exit' to quit.")
		} else {
			fmt.Println("Chat started. Type 'exit' to quit.")
		}
		for {
			fmt.Print("> ")
			if !scanner.Scan() {
				break
			}
			input := scanner.Text()
			if input == "exit" || input == "quit" || input == "q" {
				break
			}

			// Add user message history
			history = append(history, claude.Message{Role: "user", Content: input})

			// Send the full history to Claude
			resp, err := client.Chat(history)
			if err != nil {
				return err
			}

			// Add Claude's response to the history
			history = append(history, claude.Message{Role: "assistant", Content: resp})

			fmt.Println(resp)
		}

		return nil
	},
}

func init() {
	ChatCmd.Flags().StringArrayVarP(&chatIssues, "issue", "i", nil, "YouTrack issue IDs")
	ChatCmd.Flags().StringArrayVarP(&chatPRs, "pr", "p", nil, "GitHub PR number")
	ChatCmd.Flags().StringVarP(&chatFrom, "from", "f", "", "Seed chat with a generated artifact (currently: standup)")
}
