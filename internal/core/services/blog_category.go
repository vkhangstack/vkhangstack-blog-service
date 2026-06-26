package services

import (
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
)

type BlogCategoryService struct {
	repo ports.BlogCategoryRepository
}

func NewBlogCategoryService(repo ports.BlogCategoryRepository) *BlogCategoryService {
	return &BlogCategoryService{repo: repo}
}

func (s *BlogCategoryService) CreateCategory(req domain.CreateBlogCategoryRequest) (*domain.BlogCategory, error) {
	category := domain.BlogCategory{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ParentID:    req.ParentID,
		IsActive:    true,
	}
	return s.repo.CreateCategory(category)
}

func (s *BlogCategoryService) GetCategory(id string) (*domain.BlogCategory, error) {
	return s.repo.GetCategory(id)
}

func (s *BlogCategoryService) ListCategories() ([]*domain.BlogCategory, error) {
	return s.repo.ListCategories()
}

func (s *BlogCategoryService) UpdateCategory(id string, req domain.UpdateBlogCategoryRequest) (*domain.BlogCategory, error) {
	existing, err := s.repo.GetCategory(id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Slug != nil {
		existing.Slug = *req.Slug
	}
	if req.Description != nil {
		existing.Description = req.Description
	}
	if req.ParentID != nil {
		existing.ParentID = req.ParentID
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	return s.repo.UpdateCategory(*existing)
}

func (s *BlogCategoryService) DeleteCategory(id string) error {
	return s.repo.DeleteCategory(id)
}
