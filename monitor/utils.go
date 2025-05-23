package monitor

import (
	"fmt"
	"time"
)

// formatDurationHelper formats a duration into a human-readable string
func formatDurationHelper(d time.Duration) string {
	switch {
	case d >= 365*24*time.Hour:
		years := int(d.Hours() / (24 * 365))
		return fmt.Sprintf("%dY", years)
	case d >= 30*24*time.Hour:
		months := int(d.Hours() / (24 * 30))
		return fmt.Sprintf("%dM", months)
	case d >= 24*time.Hour:
		days := int(d.Hours() / 24)
		return fmt.Sprintf("%dd", days)
	case d >= time.Hour:
		hours := int(d.Hours())
		minutes := int(d.Minutes()) % 60
		return fmt.Sprintf("%dh%dm", hours, minutes)
	case d >= time.Minute:
		minutes := int(d.Minutes())
		return fmt.Sprintf("%dm", minutes)
	case d >= time.Second:
		seconds := int(d.Seconds())
		return fmt.Sprintf("%ds", seconds)
	default:
		milliseconds := int(d.Nanoseconds() / int64(time.Millisecond))
		return fmt.Sprintf("%dms", milliseconds)
	}
}

// formatTimeAgo returns a human-friendly string representing how long ago the given time was
func formatTimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)
	return formatDurationHelper(diff)
}

// formatDuration returns a human-friendly string for duration, showing whole units
func formatDuration(d time.Duration) string {
	return formatDurationHelper(d)
}
