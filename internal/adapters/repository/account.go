package repository

import (
	"context"
	"fmt"

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
	account.ID = uint64(u.snowflakeNode.GenerateIDInt64())
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
