package utils

import (
	"fmt"
	"glk-web-app/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MagicLinkClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateMagicLinkToken creates a short-lived JWT for logging in.
func GenerateMagicLinkToken(email string) (string, error) {
	secret := config.GetEnv("JWT_SECRET", "default_secret")
	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &MagicLinkClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// VerifyMagicLinkToken verifies the JWT and extracts the email.
func VerifyMagicLinkToken(tokenString string) (string, error) {
	secret := config.GetEnv("JWT_SECRET", "default_secret")
	claims := &MagicLinkClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid or expired token")
	}

	return claims.Email, nil
}
