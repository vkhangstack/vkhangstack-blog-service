package services

import (
	"fmt"

	"github.com/vkhangstack/hexagonal-architecture/internal/config"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
	"github.com/vkhangstack/hexagonal-architecture/internal/logger"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

type AccountService struct {
	repo ports.AccountRepository
}

func NewAccountService(repo ports.AccountRepository) *AccountService {
	return &AccountService{
		repo: repo,
	}
}

func (a *AccountService) CreateAccountRoot() error {
	password, err := utils.GeneratePassword(12)
	if err != nil {
		return err
	}
	account := domain.Account{
		Username: "root",
		Password: password, // In production, ensure to hash passwords and use secure practices
		Role:     domain.RoleRoot,
		Email:    nil,
		FullName: "Super Admin",
		IsActive: true,
	}
	existingAccount, err := a.repo.FindAccountByUsername(account.Username)

	if existingAccount != nil {
		logger.Log.Info("root account already exists")
		return nil
	}
	_, err = a.repo.CreateAccount(account)
	if err != nil {
		logger.Log.WithError(err).Error("failed to create root account")
	} else {
		logger.Log.WithFields(map[string]interface{}{
			"username": "root",
			"password": password,
		}).Info("root account created")
	}
	return err
}

func (a *AccountService) LoginAccount(username, password string) (*domain.LoginResponse, error) {
	userID, err := a.repo.LoginAccount(username, password)
	if err != nil {
		return nil, err
	}
	if userID == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Fetch account profile to get role for JWT and response
	account, err := a.repo.ProfileAccount(*userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch account profile: %w", err)
	}

	apiCfg := config.LoadConfig()

	accessToken, err := utils.GenerateAccessToken(*userID, account.Role, apiCfg.App.JWTSecret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(*userID, apiCfg.App.JWTSecret)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		ID:           *userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: &domain.Profile{
			ID:       account.ID,
			Username: account.Username,
			FullName: account.FullName,
			Role:     account.Role,
		},
	}, nil
}

func (a *AccountService) ProfileAccount(userID string) (*domain.Account, error) {
	return a.repo.ProfileAccount(userID)
}
