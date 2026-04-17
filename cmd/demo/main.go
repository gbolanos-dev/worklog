package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var (
	success  = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))  // green
	warning  = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))  // yellow
	header   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14")) // bold cyan
	muted    = lipgloss.NewStyle().Foreground(lipgloss.Color("245")) // mid-gray (predictable across themes)
	errStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9")) // red bold
	tag      = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))  // blue (better contrast than magenta)
)

func main() {
	fmt.Println()
	fmt.Println(header.Render("═══ worklog UI preview ═══"))
	fmt.Println()

	// --- add command ---
	fmt.Println(muted.Render("$ worklog add --tag backend \"refactored auth middleware\""))
	fmt.Println(success.Render("Logged:") + " refactored auth middleware " + tag.Render("[backend]"))
	fmt.Println()

	// --- add another ---
	fmt.Println(muted.Render("$ worklog add --tag frontend \"fixed table rendering bug\""))
	fmt.Println(success.Render("Logged:") + " fixed table rendering bug " + tag.Render("[frontend]"))
	fmt.Println()

	// --- list command (table format) ---
	fmt.Println(muted.Render("$ worklog list"))
	fmt.Println()

	t := table.New().
		Headers("#", "ID", "Date", "Entry", "Tags").
		Row("1", muted.Render("a1b2c3d4"), "2026-04-16", "refactored auth middleware", tag.Render("backend")).
		Row("2", muted.Render("e5f6a7b8"), "2026-04-16", "fixed table rendering bug", tag.Render("frontend")).
		Row("3", muted.Render("c9d0e1f2"), "2026-04-16", "reviewed PR #42", "").
		Row("4", muted.Render("34ab56cd"), "2026-04-16", "updated API rate limiting docs", tag.Render("backend")).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14")).Padding(0, 1)
			}
			return lipgloss.NewStyle().Padding(0, 1)
		})

	fmt.Println(t.Render())
	fmt.Println()

	// --- delete command ---
	fmt.Println(muted.Render("$ worklog delete a1b2c3d4"))
	fmt.Println("  " + muted.Render("Entry:") + " refactored auth middleware")
	fmt.Println(warning.Render("  Delete this entry? [y/N]") + " y")
	fmt.Println(success.Render("Deleted:") + " a1b2c3d4")
	fmt.Println()

	// --- error example ---
	fmt.Println(muted.Render("$ worklog standup"))
	fmt.Println(errStyle.Render("Error:") + " config: anthropic.api_key is required " + muted.Render("(run \"worklog init\" to set up)"))
	fmt.Println()

	// --- stats command ---
	fmt.Println(muted.Render("$ worklog stats --since 2026-04-10"))
	fmt.Println()
	fmt.Println(header.Render("Stats") + muted.Render(" (2026-04-10 → 2026-04-16)"))
	fmt.Println()
	fmt.Printf("  Total entries: %s\n", success.Render("24"))
	fmt.Printf("  Most active:   %s %s\n", "2026-04-14", muted.Render("(8 entries)"))
	fmt.Println()
	fmt.Println(header.Render("  Entries per day"))
	days := []struct{ date string; count int }{
		{"Apr 10", 3}, {"Apr 11", 5}, {"Apr 12", 2},
		{"Apr 13", 0}, {"Apr 14", 8}, {"Apr 15", 4}, {"Apr 16", 2},
	}
	for _, d := range days {
		bar := success.Render(strings.Repeat("█", d.count))
		fmt.Printf("  %s %s %s\n", muted.Render(d.date), bar, muted.Render(fmt.Sprintf("%d", d.count)))
	}
	fmt.Println()
	fmt.Println(header.Render("  Tags"))
	tags := []struct{ name string; count int }{
		{"backend", 12}, {"frontend", 6}, {"docs", 4}, {"infra", 2},
	}
	for _, tg := range tags {
		bar := tag.Render(strings.Repeat("█", tg.count))
		fmt.Printf("  %-12s %s %s\n", tg.name, bar, muted.Render(fmt.Sprintf("%d", tg.count)))
	}
	fmt.Println()

	// --- no entries ---
	fmt.Println(muted.Render("$ worklog list --date 2026-01-01"))
	fmt.Println(warning.Render("No entries found."))
	fmt.Println()
}
