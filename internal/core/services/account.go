package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

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

func (a *AccountService) CreateAccountRoot(ctx context.Context) error {
	password, err := utils.GeneratePassword(12)
	if err != nil {
		return err
	}
	account := domain.Account{
		Username:  "root",
		Password:  password, // In production, ensure to hash passwords and use secure practices
		Role:      domain.RoleRoot,
		Email:     nil,
		FirstName: "Super",
		LastName:  "Admin",
		FullName:  "Super Admin",
		Status:    string(domain.AccountStatusActive),
		IsActive:  true,
	}
	existingAccount, err := a.repo.FindAccountByUsername(ctx, account.Username)

	if existingAccount != nil {
		logger.Log.Info("root account already exists")
		return nil
	}
	_, err = a.repo.CreateAccount(ctx, account)
	if err != nil {
		logger.Log.WithError(err).Error("failed to create root account")
	} else {
		logger.Log.WithFields(map[string]interface{}{
			"username": domain.RoleRoot,
			"password": password,
		}).Info("root account created")
	}
	return err
}

func (a *AccountService) LoginAccount(ctx context.Context, username, password string) (*domain.LoginResponse, error) {
	userID, err := a.repo.LoginAccount(ctx, username, password)
	if err != nil {
		return nil, err
	}
	if userID == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Fetch account profile to get role for JWT and response
	account, err := a.repo.ProfileAccount(ctx, *userID)
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

func (a *AccountService) ProfileAccount(ctx context.Context, userID string) (*domain.Account, error) {
	return a.repo.ProfileAccount(ctx, userID)
}

func toAccountResponse(account *domain.Account) domain.AccountResponse {
	email := ""
	if account.Email != nil {
		email = *account.Email
	}
	phoneNumber := ""
	if account.PhoneNumber != nil {
		phoneNumber = *account.PhoneNumber
	}
	return domain.AccountResponse{
		ID:          account.ID,
		FirstName:   account.FirstName,
		LastName:    account.LastName,
		Username:    account.Username,
		Email:       email,
		PhoneNumber: phoneNumber,
		Status:      account.Status,
		Role:        account.Role,
		CreatedAt:   account.CreatedAt,
		UpdatedAt:   account.UpdatedAt,
	}
}

// isRootAccount reports whether an account is the hardcoded root superadmin account,
// which must stay unlockable and can't be edited/deleted/impersonated through the API.
func isRootAccount(account *domain.Account) bool {
	return account.Role == domain.RoleRoot
}

// ListAccounts returns all accounts for the /v1/account list view.
func (a *AccountService) ListAccounts(ctx context.Context) ([]domain.AccountResponse, error) {
	accounts, err := a.repo.ListAccounts(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]domain.AccountResponse, 0, len(accounts))
	for _, acc := range accounts {
		if isRootAccount(acc) {
			// Hide the root account from the list view, so it can't be edited/deleted/impersonated.
			continue
		}
		out = append(out, toAccountResponse(acc))
	}
	return out, nil
}

// GetAccount fetches a single account by ID.
func (a *AccountService) GetAccount(ctx context.Context, id string) (*domain.AccountResponse, error) {
	if id == "" {
		return nil, errors.New("account id is required")
	}
	account, err := a.repo.GetAccountByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}
	resp := toAccountResponse(account)
	return &resp, nil
}

// CreateAccount creates a new account (Part A CRUD). Only the root superadmin
// (actorIsRoot) may create another account with the root role.
func (a *AccountService) CreateAccount(ctx context.Context, req domain.CreateAccountRequest, actorIsRoot bool) (*domain.AccountResponse, error) {
	if req.Role == domain.RoleRoot && !actorIsRoot {
		return nil, fmt.Errorf("root role cannot be assigned")
	}
	if exists, _ := a.repo.CheckAccountExists(ctx, req.Username); exists {
		return nil, fmt.Errorf("username already exists")
	}
	if req.Email != nil {
		if existing, _ := a.repo.FindAccountByEmail(ctx, *req.Email); existing != nil {
			return nil, fmt.Errorf("email already exists")
		}
	}

	status := req.Status
	if status == "" {
		status = domain.AccountStatusActive
	}
	role := req.Role
	if role == "" {
		role = domain.RoleUser
	}

	account := domain.Account{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Username:    req.Username,
		Password:    req.Password,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		FullName:    strings.TrimSpace(req.FirstName + " " + req.LastName),
		Role:        role,
		Status:      string(status),
		IsActive:    status == domain.AccountStatusActive,
	}

	created, err := a.repo.CreateAccount(ctx, account)
	if err != nil {
		return nil, err
	}
	resp := toAccountResponse(created)
	return &resp, nil
}

// UpdateAccount applies a partial update to an existing account. Only the root
// superadmin (actorIsRoot) may promote another account to the root role.
func (a *AccountService) UpdateAccount(ctx context.Context, id string, req domain.UpdateAccountRequest, actorIsRoot bool) (*domain.AccountResponse, error) {
	if id == "" {
		return nil, errors.New("account id is required")
	}
	existing, err := a.repo.GetAccountByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}
	if isRootAccount(existing) {
		return nil, fmt.Errorf("root account cannot be modified")
	}
	if req.Role != nil && *req.Role == domain.RoleRoot && !actorIsRoot {
		return nil, fmt.Errorf("root role cannot be assigned")
	}

	if req.Username != nil && *req.Username != existing.Username {
		if exists, _ := a.repo.CheckAccountExists(ctx, *req.Username); exists {
			return nil, fmt.Errorf("username already exists")
		}
	}
	if req.Email != nil && (existing.Email == nil || *req.Email != *existing.Email) {
		if other, _ := a.repo.FindAccountByEmail(ctx, *req.Email); other != nil {
			return nil, fmt.Errorf("email already exists")
		}
	}

	utils.SetIfNotNil(&existing.FirstName, req.FirstName)
	utils.SetIfNotNil(&existing.LastName, req.LastName)
	utils.SetIfNotNil(&existing.Username, req.Username)
	utils.SetIfNotNil(&existing.Role, req.Role)
	if req.Email != nil {
		existing.Email = req.Email
	}
	if req.PhoneNumber != nil {
		existing.PhoneNumber = req.PhoneNumber
	}
	if req.Status != nil {
		existing.Status = string(*req.Status)
		existing.IsActive = *req.Status == domain.AccountStatusActive
	}
	if req.FirstName != nil || req.LastName != nil {
		existing.FullName = strings.TrimSpace(existing.FirstName + " " + existing.LastName)
	}
	if req.Password != nil && *req.Password != "" {
		hashed, err := utils.HashPassword(*req.Password)
		if err != nil {
			return nil, err
		}
		existing.Password = hashed
	}

	updated, err := a.repo.UpdateAccountRecord(ctx, *existing)
	if err != nil {
		return nil, err
	}
	resp := toAccountResponse(updated)
	return &resp, nil
}

// DeleteAccount soft-deletes an account.
func (a *AccountService) DeleteAccount(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("account id is required")
	}
	existing, err := a.repo.GetAccountByID(ctx, id)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}
	if isRootAccount(existing) {
		return fmt.Errorf("root account cannot be deleted")
	}
	return a.repo.DeleteAccount(ctx, id)
}

// InviteAccount creates an inactive/invited account with no password set yet.
// Only the root superadmin (actorIsRoot) may invite another account as root.
func (a *AccountService) InviteAccount(ctx context.Context, req domain.InviteAccountRequest, actorIsRoot bool) (*domain.AccountResponse, error) {
	if req.Role == domain.RoleRoot && !actorIsRoot {
		return nil, fmt.Errorf("root role cannot be assigned")
	}
	if existing, _ := a.repo.FindAccountByEmail(ctx, req.Email); existing != nil {
		return nil, fmt.Errorf("email already exists")
	}

	role := req.Role
	if role == "" {
		role = domain.RoleUser
	}

	username := strings.SplitN(req.Email, "@", 2)[0]
	if exists, _ := a.repo.CheckAccountExists(ctx, username); exists {
		username = fmt.Sprintf("%s_%d", username, time.Now().UnixNano()%10000)
	}

	account := domain.Account{
		Username: username,
		Password: "",
		Email:    &req.Email,
		FullName: username,
		Role:     role,
		Status:   string(domain.AccountStatusInvited),
		IsActive: false,
	}

	created, err := a.repo.CreateAccount(ctx, account)
	if err != nil {
		return nil, err
	}
	resp := toAccountResponse(created)
	return &resp, nil
}
