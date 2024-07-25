package util

import (
	"regexp"
	"strings"
)

func IsUsernameValid(username string) bool {
	return isNonWhitespace(username) &&
		lessThanChar(username, 25) &&
		greaterThanChar(username, 4) &&
		!containsSpecialCharacters(username)
}

func isNonWhitespace(input string) bool {
	return len(strings.TrimSpace(input)) > 0
}

func lessThanChar(input string, max int) bool {
	// we use <= on the frontend as well
	return len(input) <= max
}

func greaterThanChar(input string, min int) bool {
	// we use >= on the frontend as well
	return len(input) >= min
}

func containsSpecialCharacters(input string) bool {
	pattern := `[^a-zA-Z0-9\s]`
	re := regexp.MustCompile(pattern)
	containsSpecialCharacters := re.MatchString(input)
	return containsSpecialCharacters
}
