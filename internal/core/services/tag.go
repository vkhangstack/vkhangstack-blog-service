package services

import (
	"context"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
)

type TagService struct {
	repo ports.TagRepository
}

func NewTagService(repo ports.TagRepository) *TagService {
	return &TagService{repo: repo}
}

func (s *TagService) CreateTag(ctx context.Context, authorID string, req domain.CreateTagRequest) (*domain.Tag, error) {
	return s.repo.CreateTag(ctx, domain.Tag{
		Name:     req.Name,
		Slug:     req.Slug,
		Type:     req.Type,
		AuthorID: &authorID,
	})
}

func (s *TagService) ListTags(ctx context.Context, authorID string, tagType string) ([]*domain.Tag, error) {
	return s.repo.ListTags(ctx, authorID, tagType)
}

func (s *TagService) GetTagByID(ctx context.Context, id string) (*domain.Tag, error) {
	return s.repo.GetTagByID(ctx, id)
}

func (s *TagService) DeleteTag(ctx context.Context, id string) error {
	return s.repo.DeleteTag(ctx, id)
}
