package cmd

import (
	"fmt"

	"github.com/gbolanos-dev/worklog/store"
	"github.com/spf13/cobra"
)

var tags []string

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Log a work entry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := store.AddEntry(args[0], tags)
		if err != nil {
			return err
		}
		fmt.Println("Logged: ", args[0])
		return nil
	},
}

func init() {
	AddCmd.Flags().StringArrayVarP(&tags, "tag", "t", nil, "Tags to associate with the entry")
}
