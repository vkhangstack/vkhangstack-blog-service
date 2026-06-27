package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/uptrace/bun"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

// --- Tag ---

func (u *DB) CreateTag(tag domain.Tag) (*domain.Tag, error) {
	ctx := context.Background()
	tag.ID = u.snowflakeNode.GenerateID()
	_, err := u.db.NewInsert().Model(&tag).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("tag not saved: %v", err)
	}
	return &tag, nil
}

func (u *DB) GetTagBySlug(slug string) (*domain.Tag, error) {
	ctx := context.Background()
	tag := &domain.Tag{}
	err := u.db.NewSelect().Model(tag).Where("t.slug = ?", slug).Limit(1).Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, errors.New("tag not found")
	}
	return tag, err
}

func (u *DB) ListTags() ([]*domain.Tag, error) {
	ctx := context.Background()
	var tags []*domain.Tag
	err := u.db.NewSelect().Model(&tags).Order("t.name ASC").Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("tags not found: %v", err)
	}
	return tags, nil
}

func (u *DB) AttachTags(postID string, tagIDs []string) error {
	if len(tagIDs) == 0 {
		return nil
	}
	ctx := context.Background()
	joins := make([]domain.BlogPostTag, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		joins = append(joins, domain.BlogPostTag{PostID: postID, TagID: tagID})
	}
	_, err := u.db.NewInsert().Model(&joins).On("CONFLICT DO NOTHING").Exec(ctx)
	return err
}

func (u *DB) DetachTags(postID string) error {
	ctx := context.Background()
	_, err := u.db.NewDelete().Model((*domain.BlogPostTag)(nil)).Where("post_id = ?", postID).Exec(ctx)
	return err
}

// --- Category ---

func (u *DB) CreateCategory(category domain.BlogCategory) (*domain.BlogCategory, error) {
	ctx := context.Background()
	category.ID = u.snowflakeNode.GenerateID()
	_, err := u.db.NewInsert().Model(&category).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("category not saved: %v", err)
	}
	return &category, nil
}

func (u *DB) GetCategory(id string) (*domain.BlogCategory, error) {
	ctx := context.Background()
	category := &domain.BlogCategory{}
	err := u.cache.Get(utils.CacheKeyCategoryPrefix+id, category)
	if err == nil {
		return category, nil
	}

	err = u.db.NewSelect().Model(category).Where("bc.id = ?", id).Limit(1).Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, errors.New("category not found")
	}
	u.cache.Set(utils.CacheKeyCategoryPrefix+id, category, time.Duration(utils.CacheTTLOneWeek)*time.Second)
	return category, err
}

func (u *DB) GetCategoryBySlug(slug string) (*domain.BlogCategory, error) {
	ctx := context.Background()
	category := &domain.BlogCategory{}
	err := u.db.NewSelect().Model(category).Where("bc.slug = ?", slug).Limit(1).Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, errors.New("category not found")
	}
	return category, err
}

func (u *DB) ListCategories() ([]*domain.BlogCategory, error) {
	ctx := context.Background()
	var categories []*domain.BlogCategory
	err := u.db.NewSelect().Model(&categories).Order("bc.created_at ASC").Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("categories not found: %v", err)
	}
	return categories, nil
}

func (u *DB) UpdateCategory(category domain.BlogCategory) (*domain.BlogCategory, error) {
	ctx := context.Background()
	_, err := u.db.NewUpdate().Model(&category).WherePK().Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("category not updated: %v", err)
	}
	u.cache.Delete(utils.CacheKeyCategoryPrefix + category.ID) // Invalidate cache
	return &category, nil
}

func (u *DB) DeleteCategory(id string) error {
	ctx := context.Background()
	category := &domain.BlogCategory{}
	res, err := u.db.NewDelete().Model(category).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("category not found")
	}
	u.cache.Delete(utils.CacheKeyCategoryPrefix + id) // Invalidate cache
	return nil
}

// --- Post ---

func (u *DB) CreatePost(post domain.BlogPost, tagIDs []string) (*domain.BlogPost, error) {
	ctx := context.Background()
	post.ID = u.snowflakeNode.GenerateID()

	_, err := u.db.NewInsert().Model(&post).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("post not saved: %v", err)
	}
	if err := u.AttachTags(post.ID, tagIDs); err != nil {
		return nil, fmt.Errorf("tags not attached: %v", err)
	}
	return u.GetPost(post.ID)
}

func (u *DB) GetPost(id string) (*domain.BlogPost, error) {
	ctx := context.Background()
	post := &domain.BlogPost{}
	err := u.cache.Get(utils.CacheKeyPostPrefix+id, post)
	if err == nil {
		return post, nil
	}

	err = u.db.NewSelect().Model(post).
		Where("bp.id = ?", id).
		Limit(1).Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, errors.New("post not found")
	}
	u.cache.Set(utils.CacheKeyPostPrefix+id, post, time.Duration(utils.CacheTTLOneWeek)*time.Second)
	return post, err
}

func (u *DB) GetPostBySlug(slug string) (*domain.BlogPost, error) {
	ctx := context.Background()
	post := &domain.BlogPost{}
	err := u.cache.Get(utils.CacheKeyPostPrefix+slug, post)
	if err == nil {
		return post, nil
	}

	err = u.db.NewSelect().Model(post).
		Where("bp.slug = ?", slug).
		Limit(1).Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, errors.New("post not found")
	}
	u.cache.Set(utils.CacheKeyPostPrefix+post.Slug, post, time.Duration(utils.CacheTTLOneWeek)*time.Second)
	return post, err
}

func (u *DB) ListPosts(filter domain.BlogPostFilter) ([]*domain.BlogPost, int, error) {
	ctx := context.Background()

	applyFilters := func(q *bun.SelectQuery) *bun.SelectQuery {
		if filter.Status != "" {
			q = q.Where("bp.status = ?", filter.Status)
		}
		if filter.CategoryID != nil {
			q = q.Where("bp.category_id = ?", *filter.CategoryID)
		}
		if filter.Tag != "" {
			q = q.Join("JOIN blog_post_tags bpt ON bpt.post_id = bp.id").
				Join("JOIN tags tg ON tg.id = bpt.tag_id").
				Where("tg.slug = ?", filter.Tag)
		}
		return q
	}

	total, err := applyFilters(u.db.NewSelect().Model((*domain.BlogPost)(nil))).Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	var posts []*domain.BlogPost
	offset := (filter.Page - 1) * filter.Limit
	err = applyFilters(u.db.NewSelect().Model(&posts)).
		Order("bp.created_at DESC").
		Limit(filter.Limit).
		Offset(offset).
		Scan(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("posts not found: %v", err)
	}
	return posts, total, nil
}

func (u *DB) UpdatePost(post domain.BlogPost, tagIDs []string) error {
	ctx := context.Background()
	_, err := u.db.NewUpdate().Model(&post).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("post not updated: %v", err)
	}
	if tagIDs != nil {
		if err := u.DetachTags(post.ID); err != nil {
			return fmt.Errorf("failed to detach tags: %v", err)
		}
		if err := u.AttachTags(post.ID, tagIDs); err != nil {
			return fmt.Errorf("failed to attach tags: %v", err)
		}
	}
	u.cache.Delete(utils.CacheKeyPostPrefix + post.Slug) // Invalidate cache
	u.cache.Delete(utils.CacheKeyPostPrefix + post.ID)   // Invalidate cache
	return nil
}

func (u *DB) DeletePost(id string) error {
	ctx := context.Background()
	post := &domain.BlogPost{}
	res, err := u.db.NewDelete().Model(post).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("post not found")
	}
	u.cache.Delete(utils.CacheKeyPostPrefix + post.Slug) // Invalidate cache
	u.cache.Delete(utils.CacheKeyPostPrefix + post.ID)   // Invalidate cache
	return nil
}

func (u *DB) IncrementViewCount(id string) error {
	ctx := context.Background()
	_, err := u.db.NewUpdate().TableExpr("blog_posts").
		Set("view_count = view_count + 1").
		Where("id = ?", id).
		Exec(ctx)
	return err
}
