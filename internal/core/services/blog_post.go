package services

import (
	"errors"
	"time"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
)

type BlogPostService struct {
	repo ports.BlogPostRepository
}

func NewBlogPostService(repo ports.BlogPostRepository) *BlogPostService {
	return &BlogPostService{repo: repo}
}

func (s *BlogPostService) CreatePost(authorID string, req domain.CreateBlogPostRequest) (*domain.BlogPost, error) {
	status := req.Status
	if status == "" {
		status = domain.PostStatusDraft
	}
	post := domain.BlogPost{
		Title:         req.Title,
		Slug:          req.Slug,
		Excerpt:       req.Excerpt,
		Content:       req.Content,
		CoverImageURL: req.CoverImageURL,
		CategoryID:    req.CategoryID,
		Status:        status,
		AuthorID:      authorID,
		TagIDs:        req.TagIDs,
		LexicalState:  req.LexicalState,
	}
	if post.Status == domain.PostStatusPublished {
		now := time.Now()
		post.PublishedAt = &now
	}
	return s.repo.CreatePost(post, req.TagIDs)
}

func (s *BlogPostService) GetPost(id string) (*domain.BlogPost, error) {
	return s.repo.GetPost(id)
}

func (s *BlogPostService) GetPostBySlug(slug string) (*domain.BlogPost, error) {
	post, err := s.repo.GetPostBySlug(slug)
	if err != nil {
		return nil, err
	}
	_ = s.repo.IncrementViewCount(post.ID)
	return post, nil
}

func (s *BlogPostService) ListPosts(filter domain.BlogPostFilter) ([]*domain.BlogPost, int, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	return s.repo.ListPosts(filter)
}

func (s *BlogPostService) UpdatePost(id string, req domain.UpdateBlogPostRequest) (*domain.BlogPost, error) {
	existing, err := s.repo.GetPost(id)
	if err != nil {
		return nil, err
	}
	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Slug != nil {
		existing.Slug = *req.Slug
	}
	if req.Excerpt != nil {
		existing.Excerpt = req.Excerpt
	}
	if req.Content != nil {
		existing.Content = *req.Content
	}
	if req.CoverImageURL != nil {
		existing.CoverImageURL = req.CoverImageURL
	}
	if req.CategoryID != nil {
		existing.CategoryID = req.CategoryID
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}
	if req.ScheduledAt != nil {
		existing.ScheduledAt = req.ScheduledAt
	}
	if req.LexicalState != nil {
		existing.LexicalState = req.LexicalState
	}
	return s.repo.UpdatePost(*existing, req.TagIDs)
}

func (s *BlogPostService) DeletePost(id string) error {
	return s.repo.DeletePost(id)
}

func (s *BlogPostService) PublishPost(id string) (*domain.BlogPost, error) {
	existing, err := s.repo.GetPost(id)
	if err != nil {
		return nil, err
	}
	if existing.Status == domain.PostStatusPublished {
		return nil, errors.New("post is already published")
	}
	now := time.Now()
	existing.Status = domain.PostStatusPublished
	existing.PublishedAt = &now
	return s.repo.UpdatePost(*existing, nil)
}
