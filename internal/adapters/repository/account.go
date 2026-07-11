package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

func (u *DB) CreateAccount(ctx context.Context, account domain.Account) (*domain.Account, error) {
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

func (u *DB) FindAccountByUsername(ctx context.Context, username string) (*domain.Account, error) {
	var account domain.Account
	err := u.db.NewSelect().Model(&account).Where("username = ?", username).Limit(1).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (u *DB) LoginAccount(ctx context.Context, username, password string) (*string, error) {
	if isBlocked, _ := u.CheckAccountTemporarilyBlocked(ctx, username); isBlocked {
		return nil, fmt.Errorf("account is temporarily blocked due to multiple failed login attempts")
	}

	var account domain.Account
	err := u.db.NewSelect().Model(&account).Where("username = ? AND is_active = ?", username, true).Limit(1).Scan(ctx)
	if err != nil {
		return nil, err
	}

	if account.FailedLoginAttempts >= int(domain.FailedLoginAttemptsNumberMax) {
		if account.BlockedAt == nil || account.BlockedAt.Add(time.Minute*time.Duration(domain.FailedLoginAttemptsNumberBlockMinutes)).After(time.Now().UTC()) {
			err := u.SetAccountBlocked(ctx, username, true)
			u.SetAccountTemporarilyBlocked(ctx, username, time.Minute*time.Duration(domain.FailedLoginAttemptsNumberBlockMinutes))
			if err != nil {
				return nil, fmt.Errorf("failed to set account blocked: %w", err)
			}
		}
	}
	if account.BlockedAt != nil && account.BlockedAt.Add(time.Minute*time.Duration(domain.FailedLoginAttemptsNumberBlockMinutes)).After(time.Now().UTC()) {
		u.IncrementFailedLoginAttempts(ctx, username)
		return nil, fmt.Errorf("account is temporarily blocked due to multiple failed login attempts")
	}

	if err = utils.VerifyPassword(account.Password, password); err != nil {
		u.IncrementFailedLoginAttempts(ctx, username)
		return nil, err
	}
	u.ResetFailedLoginAttempts(ctx, username)
	return &account.ID, nil
}

func (u *DB) ProfileAccount(ctx context.Context, userID string) (*domain.Account, error) {
	var account domain.Account
	err := u.db.NewSelect().Model(&account).Where("id = ?", userID).Limit(1).Scan(ctx)
	if err != nil {
		return nil, err
	}
	account.Password = ""
	return &account, nil
}

func (u *DB) CheckAccountExists(ctx context.Context, username string) (bool, error) {
	count, err := u.db.NewSelect().Model((*domain.Account)(nil)).Where("username = ?", username).Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (u *DB) CheckAccountIsBlocked(ctx context.Context, username string) (bool, error) {
	var account domain.Account
	err := u.db.NewSelect().Model(&account).Where("username = ?", username).Limit(1).Scan(ctx)
	if err != nil {
		return false, err
	}
	return !account.IsActive, nil
}

func (u *DB) CheckAccountTemporarilyBlocked(_ context.Context, username string) (bool, error) {
	_, err := u.cache.GetString(utils.CacheKeyTemporarilyBlockedPrefix + username)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *DB) SetAccountTemporarilyBlocked(_ context.Context, username string, duration time.Duration) error {
	err := u.cache.SetString(utils.CacheKeyTemporarilyBlockedPrefix+username, "1", duration)
	if err != nil {
		return err
	}
	return nil
}

func (u *DB) SetAccountBlocked(ctx context.Context, username string, blocked bool) error {
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
func (u *DB) IncrementFailedLoginAttempts(ctx context.Context, username string) error {
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

func (u *DB) ResetFailedLoginAttempts(ctx context.Context, username string) error {
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

func (u *DB) GetRoleByUserID(ctx context.Context, userID string) (string, error) {
	var account domain.Account
	err := u.db.NewSelect().Column("role").Model(&account).Where("id = ?", userID).Limit(1).Scan(ctx)
	if err != nil {
		return "", err
	}
	return account.Role, nil
}

func (u *DB) FindAccountByEmail(ctx context.Context, email string) (*domain.Account, error) {
	var account domain.Account
	err := u.db.NewSelect().Model(&account).Where("email = ?", email).Limit(1).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (u *DB) ListAccounts(ctx context.Context) ([]*domain.Account, error) {
	var accounts []*domain.Account
	err := u.db.NewSelect().Model(&accounts).Where("username != ?", domain.RoleRoot).Order("created_at DESC").Scan(ctx)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (u *DB) GetAccountByID(ctx context.Context, id string) (*domain.Account, error) {
	var account domain.Account
	err := u.db.NewSelect().Model(&account).Where("id = ?", id).Limit(1).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (u *DB) UpdateAccountRecord(ctx context.Context, account domain.Account) (*domain.Account, error) {
	_, err := u.db.NewUpdate().Model(&account).WherePK().Exec(ctx)
	if err != nil {
		return nil, err
	}
	return u.GetAccountByID(ctx, account.ID)
}

func (u *DB) DeleteAccount(ctx context.Context, id string) error {
	account := &domain.Account{ID: id}
	_, err := u.db.NewDelete().Model(account).WherePK().Exec(ctx)
	return err
}
