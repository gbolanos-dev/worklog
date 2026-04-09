package cli

import (
	"github.com/gbolanos-dev/worklog/cmd"
	"github.com/spf13/cobra"
)

var Version = "dev"

func Run(args []string) int {
	var rootCmd = &cobra.Command{
		Use:     "worklog",
		Version: Version,
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
