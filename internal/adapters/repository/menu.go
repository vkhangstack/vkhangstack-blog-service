package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
)

func (u *DB) CreateMenu(ctx context.Context, entry domain.MenuEntry) (*domain.MenuEntry, error) {
	id := u.snowflakeNode.GenerateID()
	entry.ID = id
	now := time.Now().UTC()
	entry.CreatedAt = now

	_, err := u.db.NewInsert().Model(&entry).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("menu: create: %w", err)
	}
	return &entry, nil
}

func (u *DB) GetMenuByID(ctx context.Context, id string) (*domain.MenuEntry, error) {
	var entry domain.MenuEntry
	err := u.db.NewSelect().Model(&entry).Where("id = ? AND deleted_at IS NULL", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("menu: get: %w", err)
	}
	return &entry, nil
}

func (u *DB) UpdateMenu(ctx context.Context, id string, updates domain.MenuEntry) error {
	now := time.Now().UTC()
	updates.UpdatedAt = &now
	_, err := u.db.NewUpdate().Model(&updates).Where("id = ? AND deleted_at IS NULL", id).OmitZero().Exec(ctx)
	if err != nil {
		return fmt.Errorf("menu: update: %w", err)
	}
	return nil
}

func (u *DB) DeleteMenu(ctx context.Context, id string) error {
	now := time.Now().UTC()
	_, err := u.db.NewUpdate().
		TableExpr("menus").
		Set("deleted_at = ?", now).
		Where("id = ? AND deleted_at IS NULL", id).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("menu: delete: %w", err)
	}
	return nil
}

func (u *DB) ListMenus(ctx context.Context) ([]*domain.MenuEntry, error) {
	var entries []*domain.MenuEntry
	err := u.db.NewSelect().
		Model(&entries).
		Where("deleted_at IS NULL").
		OrderExpr("group_title ASC, sort_order ASC, created_at ASC").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("menu: list: %w", err)
	}
	return entries, nil
}
