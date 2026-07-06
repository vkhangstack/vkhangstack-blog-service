package services

import (
	"context"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

type DrawingService struct {
	repo ports.DrawingRepository
}

func NewDrawingService(repo ports.DrawingRepository) *DrawingService {
	return &DrawingService{repo: repo}
}

func (s *DrawingService) CreateDrawing(ctx context.Context, req domain.CreateDrawingRequest) (*domain.Drawing, error) {
	drawing := domain.Drawing{
		Title: req.Title,
	}
	return s.repo.CreateDrawing(ctx, drawing)
}

func (s *DrawingService) GetDrawing(ctx context.Context, id string) (*domain.Drawing, error) {
	return s.repo.GetDrawingByID(ctx, id)
}

func (s *DrawingService) ListDrawings(ctx context.Context, filter domain.DrawingFilter) ([]*domain.Drawing, int, error) {
	return s.repo.ListDrawings(ctx, filter)
}

func (s *DrawingService) UpdateDrawing(ctx context.Context, id string, req domain.UpdateDrawingRequest) error {
	drawing, err := s.repo.GetDrawingByID(ctx, id)
	if err != nil {
		return err
	}

	utils.SetIfNotNil(&drawing.Title, req.Title)

	return s.repo.UpdateDrawing(ctx, id, *drawing)
}

func (s *DrawingService) DeleteDrawing(ctx context.Context, id string) error {
	return s.repo.DeleteDrawing(ctx, id)
}

func (s *DrawingService) ListDrawingsCursor(ctx context.Context, filter domain.DrawingFilter, cursor string, limit int) ([]*domain.Drawing, *string, int, error) {
	return s.repo.ListDrawingsCursor(ctx, filter, cursor, limit)
}
