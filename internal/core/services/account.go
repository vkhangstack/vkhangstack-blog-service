package services

import (
	"fmt"

	"github.com/vkhangstack/hexagonal-architecture/internal/config"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
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
	password := "random@123" // In production, ensure to hash passwords and use secure practices
	account := domain.Account{
		Username: "root",
		Password: password, // In production, ensure to hash passwords and use secure practices
		Role:     domain.RoleRoot,
		Email:    nil,
		FullName: "Root User",
	}
	existingAccount, err := a.repo.FindAccountByUsername(account.Username)

	if existingAccount != nil {
		fmt.Println("Root account already exists")
		return nil
	}
	_, err = a.repo.CreateAccount(account)
	if err != nil {
		fmt.Printf("Error creating account: %v\n", err)
	} else {
		fmt.Println("Root account created successfully")
		fmt.Printf("Username: %s\n", "root")
		fmt.Printf("Password: %s\n", password)
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
	apiCfg := config.LoadConfig()

	accessToken, err := utils.GenerateAccessToken(*userID, apiCfg.App.JWTSecret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(*userID, apiCfg.App.JWTSecret)
	if err != nil {
		return nil, err
	}
	// In a real implementation, you would generate and return login tokens here
	return &domain.LoginResponse{
		ID:           *userID,
		Email:        "",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *AccountService) ProfileAccount(userID string) (*domain.Account, error) {
	return a.repo.ProfileAccount(userID)
}
