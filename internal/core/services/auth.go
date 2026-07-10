package services

import "github.com/vkhangstack/hexagonal-architecture/internal/utils"

// ValidateToken delegates to utils.ValidateToken.
// Returns userID, role, and any validation error.
func ValidateToken(authHeader string, jwtSecret string) (string, string, error) {
	return utils.ValidateToken(authHeader, jwtSecret)
}
