package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strconv"
	"time"
)

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GetJWTSecret() ([]byte, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET environment variable is not set")
	}
	if len(secret) < 32 {
		return nil, errors.New("JWT_SECRET must be at least 32 characters long")
	}
	return []byte(secret), nil
}

func GetJWTExpiration() (time.Duration, error) {
	expStr := os.Getenv("JWT_EXPIRATION_HOURS")
	if expStr == "" {
		return 0, errors.New("JWT_EXPIRATION_HOURS environment variable is not set")
	}
	expHours, err := strconv.Atoi(expStr)
	if err != nil {
		return 0, fmt.Errorf("invalid JWT_EXPIRATION_HOURS value: %w", err)
	}
	return time.Duration(expHours) * time.Hour, nil
}

func GenerateJWTToken(userID uint, email, username string) (string, error) {
	expirationTime, err := GetJWTExpiration()
	if err != nil {
		return "", fmt.Errorf("failed to get JWT expiration: %w", err)
	}
	expirationTimeUnix := time.Now().Add(expirationTime)

	claims := JWTClaims{
		UserID:   userID,
		Email:    email,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTimeUnix),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "SlotSwapper",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret, err := GetJWTSecret()
	if err != nil {
		return "", fmt.Errorf("failed to get JWT secret: %w", err)
	}

	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %w", err)
	}

	return signedToken, nil
}

func ValidateJWTToken(tokenStr string) (*JWTClaims, error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return GetJWTSecret()
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT token: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("invalid JWT token")
	}
	return claims, nil
}
