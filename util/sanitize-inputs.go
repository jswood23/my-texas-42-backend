package util

import "strings"

// Sanitize fixes input strings to prevent SQL injection
func Sanitize(input string) string {
	// Replace all single quotes with two single quotes
	output := strings.ReplaceAll(input, "'", "''")

	// Remove semicolons
	output = strings.ReplaceAll(output, ";", "")

	// Remove newlines
	output = strings.ReplaceAll(output, "\n", "")

	return output
}
