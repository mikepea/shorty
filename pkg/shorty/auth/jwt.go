package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// Claims represents the JWT claims
type Claims struct {
	UserID     uint   `json:"user_id"`
	Email      string `json:"email"`
	SystemRole string `json:"system_role"`
	jwt.RegisteredClaims
}

// getJWTSecret returns the JWT secret from environment or a default for development
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Default for development only - should be set in production
		secret = "shorty-dev-secret-change-in-production"
	}
	return []byte(secret)
}

// getTokenDuration returns the token validity duration
func getTokenDuration() time.Duration {
	// Default to 24 hours
	return 24 * time.Hour
}

// GenerateToken creates a new JWT token for a user
func GenerateToken(userID uint, email string, systemRole string) (string, error) {
	claims := &Claims{
		UserID:     userID,
		Email:      email,
		SystemRole: systemRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(getTokenDuration())),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "shorty",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return getJWTSecret(), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
