package goutils

import (
	"strconv"
	"unicode"
)

// assorted sort functions for slice.Sort()

// Function to compare two mixed strings
// Splits the string into alternating alphabetic and numeric segments. For example:
//
//	"a2b3" → ["a", "2", "b", "3"]
//	"a10b2" → ["a", "10", "b", "2"]
//
// It detects transitions between alphabetic and numeric parts and starts a new segment accordingly.
func CompareMixedNumericStrings(s1, s2 string) bool {
	segments1 := splitIntoSegments(s1)
	segments2 := splitIntoSegments(s2)

	// Compare corresponding segments
	for i := 0; i < len(segments1) && i < len(segments2); i++ {
		seg1 := segments1[i]
		seg2 := segments2[i]

		// Check if both segments are numeric
		num1, num1Err := strconv.Atoi(seg1)
		num2, num2Err := strconv.Atoi(seg2)

		if num1Err == nil && num2Err == nil {
			// Compare numerically
			if num1 != num2 {
				return num1 < num2
			}
		} else {
			// Compare lexicographically
			if seg1 != seg2 {
				return seg1 < seg2
			}
		}
	}

	// If all compared segments are equal, the shorter string is smaller
	return len(segments1) < len(segments2)
}

// Helper function to split a string into alternating alphabetic and numeric segments
func splitIntoSegments(s string) []string {
	var segments []string
	current := ""

	for _, char := range s {
		if unicode.IsDigit(char) {
			// Start a new segment if switching from alphabetic to numeric
			if current != "" && !unicode.IsDigit(rune(current[len(current)-1])) {
				segments = append(segments, current)
				current = ""
			}
		} else {
			// Start a new segment if switching from numeric to alphabetic
			if current != "" && unicode.IsDigit(rune(current[len(current)-1])) {
				segments = append(segments, current)
				current = ""
			}
		}
		current += string(char)
	}

	// Add the last segment
	if current != "" {
		segments = append(segments, current)
	}

	return segments
}
