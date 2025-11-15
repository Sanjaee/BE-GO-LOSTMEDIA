package utils

import (
	"crypto/rand"
	"fmt"
)

// GenerateOTP generates a 6-digit OTP code
func GenerateOTP() (string, error) {
	// Generate random 6-digit number
	max := 999999
	min := 100000

	// Generate random bytes
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Convert to integer between min and max
	num := int(b[0])<<24 | int(b[1])<<16 | int(b[2])<<8 | int(b[3])
	if num < 0 {
		num = -num
	}
	num = min + (num % (max - min + 1))

	return fmt.Sprintf("%06d", num), nil
}
