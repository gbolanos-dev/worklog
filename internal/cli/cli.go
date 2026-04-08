package cli

import (
	"github.com/gbolanos-dev/worklog/cmd"
	"github.com/spf13/cobra"
)

func Run(args []string) int {
	var rootCmd = &cobra.Command{
		Use: "worklog",
	}
	rootCmd.AddCommand(cmd.AddCmd)
	rootCmd.AddCommand(cmd.ListCmd)
	rootCmd.AddCommand(cmd.StandupCmd)

	rootCmd.SetArgs(args)

	err := rootCmd.Execute()
	if err != nil {
		return 1
	}
	return 0
}
