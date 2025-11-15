package utils

import (
	"regexp"
	"strings"
)

// ValidateEmail validates an email address
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidateUsername validates a username (3-50 chars, alphanumeric and underscore)
func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 50 {
		return false
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return usernameRegex.MatchString(username)
}

// ValidatePassword validates a password (min 8 chars, at least one letter and one number)
func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	return hasLetter && hasNumber
}

// SanitizeString removes leading/trailing whitespace and converts to lowercase
func SanitizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
