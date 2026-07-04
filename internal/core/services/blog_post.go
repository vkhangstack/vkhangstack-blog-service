package services

import (
	"context"
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

func (s *BlogPostService) CreatePost(ctx context.Context, authorID string, req domain.CreateBlogPostRequest) (*domain.BlogPost, error) {
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
		LexicalState:  req.LexicalState,
	}
	if post.Status == domain.PostStatusPublished {
		now := time.Now()
		post.PublishedAt = &now
	}
	return s.repo.CreatePost(ctx, post, req.TagIDs)
}

func (s *BlogPostService) GetPost(ctx context.Context, id string) (*domain.BlogPost, error) {
	return s.repo.GetPost(ctx, id)
}

func (s *BlogPostService) GetPostBySlug(ctx context.Context, slug string) (*domain.BlogPostBySlugResponse, error) {
	post, err := s.repo.GetPostBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	_ = s.repo.IncrementViewCount(post.ID)
	return &domain.BlogPostBySlugResponse{
		ID:            post.ID,
		Title:         post.Title,
		Slug:          post.Slug,
		Excerpt:       post.Excerpt,
		Content:       post.Content,
		CoverImageURL: post.CoverImageURL,
		CategoryID:    post.CategoryID,
		Status:        post.Status,
		ViewCount:     post.ViewCount,
		AuthorID:      post.AuthorID,
		Type:          post.Type,
		Visibility:    post.Visibility,
		Locale:        post.Locale,
	}, nil
}

func (s *BlogPostService) ListPosts(ctx context.Context, filter domain.BlogPostFilter) ([]*domain.BlogPost, int, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	return s.repo.ListPosts(ctx, filter)
}

func (s *BlogPostService) ListPostsCursor(ctx context.Context, filter domain.BlogPostFilter, cursor string, limit int) ([]*domain.BlogPost, *string, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return s.repo.ListPostsCursor(ctx, filter, cursor, limit)
}

func (s *BlogPostService) UpdatePost(ctx context.Context, id string, req domain.UpdateBlogPostRequest) error {
	existing, err := s.repo.GetPost(ctx, id)
	if err != nil {
		return err
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
	return s.repo.UpdatePost(ctx, *existing, req.TagIDs)
}

func (s *BlogPostService) DeletePost(ctx context.Context, id string) error {
	return s.repo.DeletePost(ctx, id)
}

func (s *BlogPostService) PublishPost(ctx context.Context, id string) error {
	existing, err := s.repo.GetPost(ctx, id)
	if err != nil {
		return err
	}
	if existing.Status == domain.PostStatusPublished {
		return errors.New("post is already published")
	}
	now := time.Now()
	existing.Status = domain.PostStatusPublished
	existing.PublishedAt = &now
	return s.repo.UpdatePost(ctx, *existing, nil)
}
