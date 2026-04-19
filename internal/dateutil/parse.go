package dateutil

import (
	"fmt"
	"strings"
	"time"
)

// Parse converts a user-friendly date input into YYYY-MM-DD.
func Parse(input string) (string, error) {
	now := time.Now()
	value := strings.ToLower(strings.TrimSpace(input))

	// Handle keywords first.
	switch value {
	case "today":
		return now.Format("2006-01-02"), nil
	case "yesterday":
		return now.AddDate(0, 0, -1).Format("2006-01-02"), nil
	}

	// Handle weekday names and short forms.
	if weekday, ok := parseWeekday(value); ok {
		return mostRecentWeekday(now, weekday).Format("2006-01-02"), nil
	}

	// Try short formats and inject the current year.
	shortFormats := []string{
		"1-2",
		"01-02",
		"1/2",
		"01/02",
	}

	for _, layout := range shortFormats {
		if t, err := time.Parse(layout, value); err == nil {
			withYear := time.Date(
				now.Year(),
				t.Month(),
				t.Day(),
				0, 0, 0, 0,
				now.Location(),
			)
			return withYear.Format("2006-01-02"), nil
		}
	}

	// Final attempt: full date.
	if t, err := time.Parse("2006-01-02", value); err == nil {
		return t.Format("2006-01-02"), nil
	}

	return "", fmt.Errorf("invalid date input %q", input)
}

func parseWeekday(input string) (time.Weekday, bool) {
	switch input {
	case "m", "mon", "monday":
		return time.Monday, true
	case "t", "tu", "tue", "tues", "tuesday":
		return time.Tuesday, true
	case "w", "wed", "wednesday":
		return time.Wednesday, true
	case "r", "th", "thu", "thur", "thurs", "thursday":
		return time.Thursday, true
	case "f", "fri", "friday":
		return time.Friday, true
	case "sa", "sat", "saturday":
		return time.Saturday, true
	case "su", "sun", "sunday":
		return time.Sunday, true
	default:
		return 0, false
	}
}

func mostRecentWeekday(from time.Time, target time.Weekday) time.Time {
	diff := (int(from.Weekday()) - int(target) + 7) % 7
	return from.AddDate(0, 0, -diff)
}
