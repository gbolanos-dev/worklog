package cmd

import (
	"fmt"

	"github.com/gbolanos-dev/worklog/store"
	"github.com/spf13/cobra"
)

var EditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit an entry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		entry, err := store.FindEntryById(args[0])
		if err != nil {
			return err
		}

		fmt.Printf("  Entry: %s\n", entry.Entry)


		return nil
	},
}