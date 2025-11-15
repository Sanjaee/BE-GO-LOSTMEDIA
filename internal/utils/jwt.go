package utils

import (
	"errors"
	"sync"
	"time"

	"lostmediago/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret     []byte
	jwtSecretOnce sync.Once
)

// getJWTSecret returns JWT secret, initializing it once
func getJWTSecret() []byte {
	jwtSecretOnce.Do(func() {
		if config.AppConfig != nil {
			jwtSecret = []byte(config.AppConfig.JWT.Secret)
		} else {
			// Fallback to default secret if config not loaded
			jwtSecret = []byte("D8D3DA7A75F61ACD5A4CD579EDBBC")
		}
	})
	return jwtSecret
}

type Claims struct {
	UserId string `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token
func GenerateToken(userId, email, role string) (string, error) {
	var expirationTime time.Time
	if config.AppConfig != nil {
		expirationTime = time.Now().Add(config.AppConfig.JWT.Expiration)
	} else {
		expirationTime = time.Now().Add(24 * time.Hour) // Default 24 hours
	}

	claims := &Claims{
		UserId: userId,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(getJWTSecret())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateRefreshToken generates a refresh token
func GenerateRefreshToken(userId, email, role string) (string, error) {
	var expirationTime time.Time
	if config.AppConfig != nil {
		expirationTime = time.Now().Add(config.AppConfig.JWT.RefreshExpiration)
	} else {
		expirationTime = time.Now().Add(168 * time.Hour) // Default 7 days
	}

	claims := &Claims{
		UserId: userId,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(getJWTSecret())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
