package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomToken generates a random token for email verification and password reset
func GenerateRandomToken(length int) (string, error) {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateEmailVerificationToken generates a token for email verification
func GenerateEmailVerificationToken() (string, error) {
	return GenerateRandomToken(64)
}

// GeneratePasswordResetToken generates a token for password reset
func GeneratePasswordResetToken() (string, error) {
	return GenerateRandomToken(64)
}

