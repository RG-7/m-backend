package helpers

import (
	"strings"
	"time"
)

// GenerateSubgroups extracts subgroups from the given range (e.g., "1A1A-1A1E").
func GenerateSubgroups(subgroupRange string) []string {
	parts := strings.Split(subgroupRange, "-")
	if len(parts) != 2 {
		return nil
	}

	prefix := parts[0][:3]
	startLetter := parts[0][3]
	endLetter := parts[1][3]

	if startLetter > endLetter {
		return nil
	}

	subgroups := []string{}
	for c := startLetter; c <= endLetter; c++ {
		subgroups = append(subgroups, prefix+string(c))
	}

	return subgroups
}

// ParseDayOfWeek converts a string day to a time.Weekday.
func ParseDayOfWeek(day string) time.Weekday {
	switch strings.ToLower(day) {
	case "monday":
		return time.Monday
	case "tuesday":
		return time.Tuesday
	case "wednesday":
		return time.Wednesday
	case "thursday":
		return time.Thursday
	case "friday":
		return time.Friday
	case "saturday":
		return time.Saturday
	case "sunday":
		return time.Sunday
	default:
		return -1
	}
}

// Function to parse date
func parseDate(dateStr string) (time.Time, error) {
	layout := "2006-01-02"
	return time.Parse(layout, dateStr)
}

// GetDuration returns the duration based on session type.
func GetDuration(sessionType string) int {
	switch strings.ToLower(sessionType) {
	case "l", "t":
		return 50
	case "p":
		return 100
	default:
		return 0
	}
}
