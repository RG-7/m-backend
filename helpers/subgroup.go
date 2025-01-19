package helpers

import (
	"strings"
	"time"
)

// GenerateSubgroups extracts subgroups from the given range (e.g., "1A1A-1A1E").
func GenerateSubgroups(subgroupRange string) []string {
	// Check if the range contains a hyphen indicating a range
	if strings.Contains(subgroupRange, "-") {
		parts := strings.Split(subgroupRange, "-")
		if len(parts) != 2 {
			// Return an empty slice if the range is invalid
			return []string{}
		}

		// Extract prefix and range letters
		prefix := parts[0][:len(parts[0])-1] // Assuming the prefix is everything except the last letter
		startLetter := parts[0][len(parts[0])-1]
		endLetter := parts[1][len(parts[1])-1]

		// Ensure start letter is not greater than end letter
		if startLetter > endLetter {
			return []string{}
		}

		// Generate subgroups for the range
		subgroups := []string{}
		for c := startLetter; c <= endLetter; c++ {
			subgroups = append(subgroups, prefix+string(c))
		}

		return subgroups
	}

	// If the range does not contain a hyphen, it's a single subgroup, so return it as a list
	return []string{subgroupRange}
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
