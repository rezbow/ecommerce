package authentication

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	UserId  uuid.UUID
	IsAdmin bool
	jwt.RegisteredClaims
}

func NewJWTToken(
	userId uuid.UUID,
	isAdmin bool,
	secretKey string,
) (string, error) {
	expiresAt := time.Now().Add(time.Hour * 24)
	claims := &UserClaims{
		UserId:  userId,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ValidateToken(tokenString string, secretKey string) (*UserClaims, error) {
	claims := &UserClaims{}

	// Parse the token, providing the secret key callback
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token's signing method is what we expect (HMAC)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token validation error: %w", err)
	}

	// Final check: Is the token valid?
	if !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	return claims, nil
}
