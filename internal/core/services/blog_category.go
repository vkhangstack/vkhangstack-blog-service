package services

import (
	"context"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
)

type BlogCategoryService struct {
	repo ports.BlogCategoryRepository
}

func NewBlogCategoryService(repo ports.BlogCategoryRepository) *BlogCategoryService {
	return &BlogCategoryService{repo: repo}
}

func (s *BlogCategoryService) CreateCategory(ctx context.Context, req domain.CreateBlogCategoryRequest) (*domain.BlogCategory, error) {
	category := domain.BlogCategory{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ParentID:    req.ParentID,
		IsActive:    true,
	}
	return s.repo.CreateCategory(ctx, category)
}

func (s *BlogCategoryService) GetCategory(ctx context.Context, id string) (*domain.BlogCategory, error) {
	return s.repo.GetCategory(ctx, id)
}

func (s *BlogCategoryService) ListCategories(ctx context.Context) ([]*domain.BlogCategoryWithPostCount, error) {
	return s.repo.ListCategories(ctx)
}

func (s *BlogCategoryService) ListCategoriesCursor(ctx context.Context, cursor string, limit int) ([]*domain.BlogCategoryWithPostCount, *string, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return s.repo.ListCategoriesCursor(ctx, cursor, limit)
}

func (s *BlogCategoryService) UpdateCategory(ctx context.Context, id string, req domain.UpdateBlogCategoryRequest) (*domain.BlogCategory, error) {
	existing, err := s.repo.GetCategory(ctx, id)
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
	return s.repo.UpdateCategory(ctx, *existing)
}

func (s *BlogCategoryService) DeleteCategory(ctx context.Context, id string) error {
	return s.repo.DeleteCategory(ctx, id)
}
