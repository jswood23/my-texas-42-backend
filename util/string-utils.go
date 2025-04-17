package util

import (
	"math/rand"
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

func SliceContains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func containsSpecialCharacters(input string) bool {
	pattern := `[^a-zA-Z0-9\s]`
	re := regexp.MustCompile(pattern)
	containsSpecialCharacters := re.MatchString(input)
	return containsSpecialCharacters
}

// IsEmailValid checks if the email is in the correct format
func IsEmailValid(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	isEmailValid := re.MatchString(email)
	return isEmailValid
}

// IsDominoNameValid checks if the domino name is in the correct format (e.g. "0-0")
func IsDominoNameValid(d models.DominoName) bool {
	return regexp.MustCompile(`^[0-6]-[0-6]$`).MatchString(string(d))
}

// GenerateInviteCode generates a random 6 character string with all uppercase letters
func GenerateInviteCode() models.InviteCode {
	inviteCode := make([]byte, 6)
	for i := range inviteCode {
		inviteCode[i] = 'A' + byte(rand.Intn(26))
	}
	return models.InviteCode(inviteCode)
}

// IsInviteCodeValid checks if the invite code has 6 characters and is all uppercase letters
func IsInviteCodeValid(inviteCode string) bool {
	return regexp.MustCompile(`^[A-Z]{6}$`).MatchString(inviteCode)
}
