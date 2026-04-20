package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gbolanos-dev/worklog/store"
	"github.com/spf13/cobra"
)

var force bool

var DeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a work entry by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		entry, err := store.FindEntryById(args[0])
		if err != nil {
			return err
		}

		fmt.Printf("  Entry: %s\n", entry.Entry)

		if !force {
			ok, err := confirmDelete()
			if err != nil {
				return err
			}
			if !ok {
				fmt.Println("Cancelled.")
				return nil
			}
		}

		err = store.DeleteEntry(entry.ID)
		if err != nil {
			return err
		}

		fmt.Printf("Deleted: %s\n", entry.ID[:8])
		return nil
	},
}

func confirmDelete() (bool, error) {
	fmt.Print("  Delete this entry? [y/N] ")
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return false, err
	}

	switch strings.ToLower(strings.TrimSpace(input)) {
	case "y", "yes":
		return true, nil
	default:
		return false, nil
	}
}

func init() {
	DeleteCmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation prompt")
}
