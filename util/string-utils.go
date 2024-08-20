package util

import (
	"my-texas-42-backend/models"
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

func IsEmailValid(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	isEmailValid := re.MatchString(email)
	return isEmailValid
}

func IsDominoNameValid(d models.DominoName) bool {
	return regexp.MustCompile(`^[0-6]-[0-6]$`).MatchString(string(d))
}
