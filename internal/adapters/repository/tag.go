package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
)

// --- Tag ---
func (u *DB) CreateTag(ctx context.Context, tag domain.Tag) (*domain.Tag, error) {
	tag.ID = u.snowflakeNode.GenerateID()
	_, err := u.db.NewInsert().Model(&tag).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("tag not saved: %v", err)
	}
	return &tag, nil
}

func (u *DB) GetTagBySlug(ctx context.Context, slug string) (*domain.Tag, error) {
	tag := &domain.Tag{}
	err := u.db.NewSelect().Model(tag).Where("t.slug = ?", slug).Limit(1).Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, errors.New("tag not found")
	}
	return tag, err
}

func (u *DB) ListTags(ctx context.Context, authorID string, tagType string) ([]*domain.Tag, error) {
	var tags []*domain.Tag
	query := u.db.NewSelect().Model(&tags).Where("t.author_id = ?", authorID).Order("t.name ASC")
	if tagType != "" {
		query = query.Where("t.type = ?", tagType)
	}
	err := query.Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("tags not found: %v", err)
	}
	return tags, nil
}

func (u *DB) AttachTags(ctx context.Context, postID string, tagIDs []string) error {
	if len(tagIDs) == 0 {
		return nil
	}
	joins := make([]domain.BlogPostTag, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		joins = append(joins, domain.BlogPostTag{PostID: postID, TagID: tagID})
	}
	_, err := u.db.NewInsert().Model(&joins).On("CONFLICT DO NOTHING").Exec(ctx)
	return err
}

func (u *DB) DetachTags(ctx context.Context, postID string) error {
	_, err := u.db.NewDelete().Model((*domain.BlogPostTag)(nil)).Where("post_id = ?", postID).Exec(ctx)
	return err
}
