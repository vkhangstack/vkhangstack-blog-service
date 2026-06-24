package services

import (
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
)

type TagService struct {
	repo ports.TagRepository
}

func NewTagService(repo ports.TagRepository) *TagService {
	return &TagService{repo: repo}
}

func (s *TagService) CreateTag(req domain.CreateTagRequest) (*domain.Tag, error) {
	return s.repo.CreateTag(domain.Tag{
		Name: req.Name,
		Slug: req.Slug,
	})
}

func (s *TagService) ListTags() ([]*domain.Tag, error) {
	return s.repo.ListTags()
}
