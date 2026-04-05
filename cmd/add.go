package cmd

import (
	"fmt"

	"github.com/gbolanos-dev/worklog/store"
	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Log a work entry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := store.AddEntry(args[0])
		if err != nil {
			return err
		}
		fmt.Println("Logged: ", args[0])
		return nil
	},
}
