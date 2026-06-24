package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
)

// GenerateAccessToken creates a short-lived JWT access token (1 hour) for the given user ID.
func GenerateAccessToken(userID, jwtSecret string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    domain.JwtIssuerAccess,
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour).UTC()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// GenerateRefreshToken creates a long-lived JWT refresh token (7 days) for the given user ID.
func GenerateRefreshToken(userID, jwtSecret string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    domain.JwtIssuerRefresh,
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour).UTC()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// ValidateToken parses and validates a Bearer token from the Authorization header.
// Returns the user ID (subject claim) on success.
func ValidateToken(authHeader, jwtSecret string) (string, error) {
	if authHeader == "" {
		return "", errors.New("token not found")
	}
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", errors.New("malformed authorization header")
	}

	tokenString := authHeader[7:]
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("token is not valid")
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}
	if claims.ExpiresAt == nil || claims.ExpiresAt.Before(time.Now().UTC()) {
		return "", errors.New("token has expired")
	}
	if claims.Issuer == domain.JwtIssuerRefresh {
		return "", errors.New("token is a refresh token, please use access token")
	}

	return claims.Subject, nil
}
