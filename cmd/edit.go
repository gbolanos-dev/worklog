package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gbolanos-dev/worklog/store"
	"github.com/spf13/cobra"
)

var EditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit an entry",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		entry, err := store.FindEntryById(args[0])
		if err != nil {
			return err
		}

		switch len(args) {
		case 1:
			newText, err := editInEditor(entry)
			if err != nil {
				return err
			}
			if newText == entry.Entry {
				fmt.Println("No changes")
				return nil
			}
			err = store.EditEntry(entry.ID, newText)
			if err != nil {
				return err
			}
		case 2:
			err := store.EditEntry(args[0], args[1])
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("expected 1 or 2 args, got %d", len(args))
		}

		return nil
	},
}

func getEditor() string {
	if e := os.Getenv("VISUAL"); e != "" {
		return e
	}
	if e := os.Getenv("EDITOR"); e != "" {
		return e
	}
	return "vi"
}

func openFileInEditor(filename string) error {
	editor := getEditor()
	cmd := exec.Command(editor, filename)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func editInEditor(entry *store.Entry) (string, error) {
	tmp, err := os.CreateTemp("", "worklog-edit-*.txt")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmp.Name())

	_, err = tmp.WriteString(entry.Entry)
	if err != nil {
		return "", err
	}
	tmp.Close()

	err = openFileInEditor(tmp.Name())
	if err != nil {
		return "", err
	}

	newText, err := os.ReadFile(tmp.Name())
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(newText)), nil
}
