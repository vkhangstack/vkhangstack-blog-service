package services

import (
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
)

type CustomerService struct {
	repo ports.CustomerRepository
}

func NewCustomerService(repo ports.CustomerRepository) *CustomerService {
	return &CustomerService{
		repo: repo,
	}
}

func (u *CustomerService) CreateUser(email, password string) (*domain.Customer, error) {
	return u.repo.CreateUser(email, password)
}

func (u *CustomerService) ReadUser(id uint64) (*domain.Customer, error) {
	return u.repo.ReadUser(id)
}

func (u *CustomerService) ReadUsers() ([]*domain.Customer, error) {
	return u.repo.ReadUsers()
}

func (u *CustomerService) UpdateUser(id, email, password string) error {
	return u.repo.UpdateUser(id, email, password)
}

func (u *CustomerService) DeleteUser(id uint64) error {
	return u.repo.DeleteUser(id)
}

func (u *CustomerService) LoginUser(email, password string) (*domain.LoginResponse, error) {
	return u.repo.LoginUser(email, password)
}

func (u *CustomerService) UpdateMembershipStatus(id uint64, status bool) error {
	return u.repo.UpdateMembershipStatus(id, status)
}
