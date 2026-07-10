package services

import (
	"context"
	"fmt"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
)

type MenuService struct {
	repo ports.MenuRepository
}

func NewMenuService(repo ports.MenuRepository) *MenuService {
	return &MenuService{repo: repo}
}

func (s *MenuService) CreateMenu(ctx context.Context, req domain.CreateMenuRequest) (*domain.MenuEntry, error) {
	entry := domain.MenuEntry{
		GroupTitle: req.GroupTitle,
		ParentID:   req.ParentID,
		Title:      req.Title,
		URL:        req.URL,
		Icon:       req.Icon,
		Badge:      req.Badge,
		Resource:   req.Resource,
		SortOrder:  req.SortOrder,
		IsActive:   req.IsActive,
	}
	result, err := s.repo.CreateMenu(ctx, entry)
	if err != nil {
		return nil, fmt.Errorf("menu service: create: %w", err)
	}
	return result, nil
}

func (s *MenuService) GetMenu(ctx context.Context, id string) (*domain.MenuEntry, error) {
	return s.repo.GetMenuByID(ctx, id)
}

func (s *MenuService) UpdateMenu(ctx context.Context, id string, req domain.UpdateMenuRequest) error {
	existing, err := s.repo.GetMenuByID(ctx, id)
	if err != nil {
		return err
	}

	// Apply partial updates
	if req.GroupTitle != nil {
		existing.GroupTitle = *req.GroupTitle
	}
	if req.ParentID != nil {
		existing.ParentID = req.ParentID
	}
	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.URL != nil {
		existing.URL = req.URL
	}
	if req.Icon != nil {
		existing.Icon = req.Icon
	}
	if req.Badge != nil {
		existing.Badge = req.Badge
	}
	if req.Resource != nil {
		existing.Resource = req.Resource
	}
	if req.SortOrder != nil {
		existing.SortOrder = *req.SortOrder
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	return s.repo.UpdateMenu(ctx, id, *existing)
}

func (s *MenuService) DeleteMenu(ctx context.Context, id string) error {
	return s.repo.DeleteMenu(ctx, id)
}

func (s *MenuService) ListMenus(ctx context.Context) ([]*domain.MenuEntry, error) {
	return s.repo.ListMenus(ctx)
}
