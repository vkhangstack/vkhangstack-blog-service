package repository

import (
	"context"
	"fmt"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

func (u *DB) CreateDrawing(ctx context.Context, drawing domain.Drawing) (*domain.Drawing, error) {
	drawing.ID = u.snowflakeNode.GenerateID()
	_, err := u.db.NewInsert().Model(&drawing).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("drawing not saved: %v", err)
	}
	return &drawing, nil
}

func (u *DB) GetDrawingByID(ctx context.Context, id string) (*domain.Drawing, error) {
	drawing := &domain.Drawing{}
	err := u.db.NewSelect().Model(drawing).Where("d.id = ?", id).Limit(1).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("drawing not found: %v", err)
	}
	return drawing, nil
}

func (u *DB) UpdateDrawing(ctx context.Context, id string, updates domain.Drawing) error {
	updates.ID = id
	_, err := u.db.NewUpdate().Model(&updates).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("drawing not updated: %v", err)
	}
	return nil
}

func (u *DB) DeleteDrawing(ctx context.Context, id string) error {
	drawing := &domain.Drawing{ID: id}
	_, err := u.db.NewDelete().Model(drawing).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("drawing not deleted: %v", err)
	}
	return nil
}

func (u *DB) ListDrawings(ctx context.Context, filter domain.DrawingFilter) ([]*domain.Drawing, int, error) {
	drawings := make([]*domain.Drawing, 0)
	query := u.db.NewSelect().Model(&drawings)

	if filter.Title != nil {
		query = query.Where("d.title = ?", *filter.Title)
	}

	total, err := query.Order("d.id DESC").Offset((filter.Page - 1) * filter.Limit).Limit(filter.Limit).ScanAndCount(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list drawings: %v", err)
	}

	return drawings, total, nil
}

func (u *DB) ListDrawingsCursor(ctx context.Context, filter domain.DrawingFilter, cursor string, limit int) ([]*domain.Drawing, *string, int, error) {
	var cursorID string
	if cursor != "" {
		id, err := utils.DecodeCursor(cursor)
		if err != nil {
			return nil, nil, 0, err
		}
		cursorID = id
	}

	drawings := make([]*domain.Drawing, 0)
	query := u.db.NewSelect().Model(&drawings)

	if filter.Title != nil {
		query = query.Where("d.title = ?", *filter.Title)
	}

	queryCount := query.Clone()
	query = query.Order("t.created_at DESC", "t.id DESC")

	if cursorID != "" {
		cursorTask := &domain.Task{}
		err := u.db.NewSelect().Model(cursorTask).Where("t.id = ?", cursorID).Limit(1).Scan(ctx)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("cursor task not found: %v", err)
		}
		query = query.Where("(t.created_at, t.id) < (?, ?)", cursorTask.CreatedAt, cursorID)
	}
	err := query.Limit(limit + 1).Scan(ctx)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("tasks not found: %v", err)
	}

	var nextCursor *string
	if len(drawings) > limit {
		drawings = drawings[:limit]
		nextCursor = utils.StringPtr(utils.EncodeCursor(drawings[len(drawings)-1].ID))
	}

	total, _ := queryCount.Count(ctx)
	return drawings, nextCursor, total, nil
}
