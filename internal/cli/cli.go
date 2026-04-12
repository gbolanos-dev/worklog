package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gbolanos-dev/worklog/cmd"
	"github.com/gbolanos-dev/worklog/internal/update"
	"github.com/spf13/cobra"
)

var Version = "dev"
var updateFlag bool

func Run(args []string) int {
	var rootCmd = &cobra.Command{
		Use:     "worklog",
		Version: Version,
		RunE: func(cmd *cobra.Command, args []string) error {
			if updateFlag {
				return update.DoUpdate("github.com/gbolanos-dev/worklog/cmd/worklog")
			}
			return cmd.Help()
		},
	}

	// Bind the update flag
	rootCmd.PersistentFlags().BoolVarP(
		&updateFlag,
		"update",
		"u",
		false,
		"Update worklog to the latest version")

	rootCmd.AddCommand(cmd.AddCmd)
	rootCmd.AddCommand(cmd.ListCmd)
	rootCmd.AddCommand(cmd.StandupCmd)

	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	if err != nil {
		return 1
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}

	cacheDir := filepath.Join(home, ".worklog")

	result := update.Check(Version, cacheDir)
	if result == nil {
		update.Refresh(cacheDir)
	}

	if result != nil && result.Available {
		_, _ = fmt.Fprintf(os.Stderr, "Update available: %s -> run: worklog --update\n", result.Latest)
	}

	return 0
}
