package services

import "github.com/vkhangstack/hexagonal-architecture/internal/utils"

// ValidateToken delegates to utils.ValidateToken.
// Kept here for backward compatibility with middleware that imports this package.
func ValidateToken(authHeader string, jwtSecret string) (string, error) {
	return utils.ValidateToken(authHeader, jwtSecret)
}
