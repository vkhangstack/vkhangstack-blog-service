package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

func (u *DB) CreateAccount(account domain.Account) (*domain.Account, error) {
	ctx := context.Background()
	if account.Password != "" {
		hashedPassword, err := utils.HashPassword(account.Password)
		if err != nil {
			return nil, fmt.Errorf("password not hashed: %w", err)
		}
		account.Password = hashedPassword
	}
	account.ID = u.snowflakeNode.GenerateID()
	_, err := u.db.NewInsert().Model(&account).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (u *DB) FindAccountByUsername(username string) (*domain.Account, error) {
	ctx := context.Background()
	var account domain.Account
	err := u.db.NewSelect().Model(&account).Where("username = ?", username).Limit(1).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (u *DB) LoginAccount(username, password string) (*string, error) {
	ctx := context.Background()
	var account domain.Account
	err := u.db.NewSelect().Model(&account).Where("username = ? AND is_active = ?", username, true).Limit(1).Scan(ctx)
	if err != nil {
		return nil, err
	}

	if account.FailedLoginAttempts >= 5 {
		if account.BlockedAt == nil || account.BlockedAt.Add(time.Minute*15).After(time.Now().UTC()) {
			err := u.SetAccountBlocked(username, true)
			if err != nil {
				return nil, fmt.Errorf("failed to set account blocked: %w", err)
			}
		}
	}
	if account.BlockedAt != nil && account.BlockedAt.Add(time.Minute*15).After(time.Now().UTC()) {
		u.IncrementFailedLoginAttempts(username)
		return nil, fmt.Errorf("account is temporarily blocked due to multiple failed login attempts")
	}

	if err = utils.VerifyPassword(account.Password, password); err != nil {
		u.IncrementFailedLoginAttempts(username)
		return nil, err
	}
	u.ResetFailedLoginAttempts(username)
	return &account.ID, nil
}

func (u *DB) ProfileAccount(userID string) (*domain.Account, error) {
	ctx := context.Background()
	var account domain.Account
	err := u.db.NewSelect().Model(&account).Where("id = ?", userID).Limit(1).Scan(ctx)
	if err != nil {
		return nil, err
	}
	account.Password = ""
	return &account, nil
}

func (u *DB) CheckAccountExists(username string) (bool, error) {
	ctx := context.Background()
	count, err := u.db.NewSelect().Model((*domain.Account)(nil)).Where("username = ?", username).Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (u *DB) CheckAccountIsBlocked(username string) (bool, error) {
	ctx := context.Background()
	var account domain.Account
	err := u.db.NewSelect().Model(&account).Where("username = ?", username).Limit(1).Scan(ctx)
	if err != nil {
		return false, err
	}
	return !account.IsActive, nil
}

func (u *DB) CheckAccountTemporarilyBlocked(username string) (bool, error) {
	_, err := u.cache.GetString(utils.CacheKeyTemporarilyBlockedPrefix + username)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *DB) SetAccountTemporarilyBlocked(username string, duration time.Duration) error {
	err := u.cache.SetString(utils.CacheKeyTemporarilyBlockedPrefix+username, "1", duration)
	if err != nil {
		return err
	}
	return nil
}

func (u *DB) SetAccountBlocked(username string, blocked bool) error {
	ctx := context.Background()
	var account domain.Account
	err := u.db.NewSelect().Model(&account).Where("username = ?", username).Limit(1).Scan(ctx)
	if err != nil {
		return err
	}
	account.BlockedAt = nil
	if blocked {
		now := time.Now()
		account.BlockedAt = &now
	}
	_, err = u.db.NewUpdate().Model(&account).Where("id = ?", account.ID).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
func (u *DB) IncrementFailedLoginAttempts(username string) error {
	ctx := context.Background()
	var account domain.Account
	err := u.db.NewSelect().Model(&account).Where("username = ?", username).Limit(1).Scan(ctx)
	if err != nil {
		return err
	}
	account.FailedLoginAttempts++
	_, err = u.db.NewUpdate().Model(&account).Where("id = ?", account.ID).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (u *DB) ResetFailedLoginAttempts(username string) error {
	ctx := context.Background()
	var account domain.Account
	err := u.db.NewSelect().Model(&account).Where("username = ?", username).Limit(1).Scan(ctx)
	if err != nil {
		return err
	}
	account.FailedLoginAttempts = 0
	account.BlockedAt = nil
	_, err = u.db.NewUpdate().Model(&account).Where("id = ?", account.ID).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
