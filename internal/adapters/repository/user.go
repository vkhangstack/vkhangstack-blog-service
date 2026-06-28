package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/vkhangstack/hexagonal-architecture/internal/config"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/logger"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

func (u *DB) CreateUser(email, password string) (*domain.Customer, error) {
	ctx := context.Background()

	existing := &domain.Customer{}
	err := u.db.NewSelect().Model(existing).Where("email = ?", email).Limit(1).Scan(ctx)
	if err == nil {
		return nil, errors.New("user already exists")
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

	user := &domain.Customer{
		ID:    uint64(u.snowflakeNode.GenerateIDInt64()),
		Email: email,
	}
	_, err = u.db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("user not saved: %v", err)
	}
	return user, nil
}

func (u *DB) ReadUser(id uint64) (*domain.Customer, error) {
	ctx := context.Background()
	user := &domain.Customer{}
	cacheKey := utils.Uint64ToString(id)

	if err := u.cache.Get(cacheKey, user); err == nil {
		return user, nil
	}

	err := u.db.NewSelect().Model(user).Where("id = ?", id).Limit(1).Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	if cacheErr := u.cache.Set(cacheKey, user, time.Minute*10); cacheErr != nil {
		logger.Log.WithError(cacheErr).Warn("error storing user in cache")
	}
	return user, nil
}

func (u *DB) ReadUsers() ([]*domain.Customer, error) {
	ctx := context.Background()
	var users []*domain.Customer
	err := u.db.NewSelect().Model(&users).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("users not found: %v", err)
	}
	return users, nil
}

func (u *DB) UpdateUser(id, email, password string) error {
	ctx := context.Background()

	user := &domain.Customer{}
	err := u.db.NewSelect().Model(user).Where("id = ?", id).Limit(1).Scan(ctx)
	if err == sql.ErrNoRows {
		return errors.New("user not found")
	}
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user.Email = email
	user.Password = hashedPassword

	res, err := u.db.NewUpdate().Model(user).Column("email", "password").Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("unable to update user :(")
	}

	if cacheErr := u.cache.Delete(id); cacheErr != nil {
		logger.Log.WithError(cacheErr).Warn("error deleting user from cache")
	}
	return nil
}

func (u *DB) DeleteUser(id uint64) error {
	ctx := context.Background()
	user := &domain.Customer{}
	res, err := u.db.NewDelete().Model(user).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("user not found")
	}
	if cacheErr := u.cache.Delete(utils.Uint64ToString(id)); cacheErr != nil {
		logger.Log.WithError(cacheErr).Warn("error deleting user from cache")
	}
	return nil
}

func (u *DB) LoginUser(email, password string) (*domain.LoginResponse, error) {
	apiCfg := config.LoadConfig()
	user, err := u.findUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if err = utils.VerifyPassword(user.Password, password); err != nil {
		return nil, err
	}

	accessToken, err := utils.GenerateAccessToken(utils.Uint64ToString(user.ID), apiCfg.App.JWTSecret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(utils.Uint64ToString(user.ID), apiCfg.App.JWTSecret)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *DB) UpdateMembershipStatus(id uint64, membership bool) error {
	ctx := context.Background()

	user := &domain.Customer{}
	err := u.db.NewSelect().Model(user).Where("id = ?", id).Limit(1).Scan(ctx)
	if err == sql.ErrNoRows {
		return errors.New("user not found")
	}
	if err != nil {
		return err
	}

	res, err := u.db.NewUpdate().Model(user).OmitZero().Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("unable to update membership status :(")
	}
	return nil
}

func (u *DB) findUserByEmail(email string) (*domain.Customer, error) {
	ctx := context.Background()
	user := &domain.Customer{}
	err := u.db.NewSelect().Model(user).Where("email = ?", email).Limit(1).Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
