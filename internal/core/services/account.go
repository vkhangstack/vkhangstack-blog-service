package services

import (
	"fmt"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
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
	password := "random" // In production, ensure to hash passwords and use secure practices
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
