package repository

import "github.com/vkhangstack/hexagonal-architecture/internal/core/domain"

func (c *DB) CreateCheckoutSession(userID string, payment domain.Payment) error {
	return nil
}
