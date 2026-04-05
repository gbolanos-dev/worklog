package main

import (
	"github.com/gbolanos-dev/worklog/cmd"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use: "worklog",
	}
	rootCmd.AddCommand(cmd.AddCmd)
	rootCmd.AddCommand(cmd.ListCmd)
	err := rootCmd.Execute()
	if err != nil {
		return
	}
}
