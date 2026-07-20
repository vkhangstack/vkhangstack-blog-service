package repository

import (
	"context"
	"fmt"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

func (u *DB) CreateTimetableEntry(ctx context.Context, timetable domain.TimetableEntry) (*domain.TimetableEntry, error) {

	timetable.ID = u.snowflakeNode.GenerateID()
	_, err := u.db.NewInsert().Model(&timetable).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &timetable, nil
}

func (u *DB) GetTimetableEntryByID(ctx context.Context, authorID, id string) (*domain.TimetableEntry, error) {
	var timetable domain.TimetableEntry
	err := u.db.NewSelect().Model(&timetable).Where("id = ? AND author_id = ?", id, authorID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &timetable, nil
}

func (u *DB) UpdateTimetableEntry(ctx context.Context, id string, updates domain.TimetableEntry) error {
	_, err := u.db.NewUpdate().Model(&updates).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (u *DB) DeleteTimetableEntry(ctx context.Context, id string) error {
	_, err := u.db.NewDelete().Model(&domain.TimetableEntry{}).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (u *DB) ListTimetableEntries(ctx context.Context, authorID string, filter domain.TimetableEntryFilter) ([]*domain.TimetableEntry, int, error) {
	timetableEntries := make([]*domain.TimetableEntry, 0)
	query := u.db.NewSelect().Model(&timetableEntries)

	if filter.DayOfWeek != nil {
		query = query.Where("te.day_of_week = ?", filter.DayOfWeek)
	}

	total, err := query.Order("te.day_of_week ASC").Offset((filter.Page - 1) * filter.Limit).Limit(filter.Limit).ScanAndCount(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("timetable entries not found: %v", err)
	}

	return timetableEntries, total, nil
}

func (u *DB) ListTimetableEntriesCursor(ctx context.Context, authorID string, filter domain.TimetableEntryFilter, cursor string, limit int) ([]*domain.TimetableEntry, *string, int, error) {
	var cursorID string
	if cursor != "" {
		id, err := utils.DecodeCursor(cursor)
		if err != nil {
			return nil, nil, 0, err
		}
		cursorID = id
	}

	var timetableEntries []*domain.TimetableEntry
	query := u.db.NewSelect().Model(&timetableEntries).Where("te.author_id = ?", authorID)

	if filter.DayOfWeek != nil {
		query = query.Where("te.day_of_week = ?", filter.DayOfWeek)
	}
	queryCount := query.Clone()
	query = query.Order("te.created_at DESC", "te.id DESC")

	if cursorID != "" {
		cursorTimetableEntry := &domain.TimetableEntry{}
		err := u.db.NewSelect().Model(cursorTimetableEntry).Where("te.id = ? AND te.author_id = ?", cursorID, authorID).Limit(1).Scan(ctx)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("cursor timetable entry not found: %v", err)
		}
		query = query.Where("(te.created_at, te.id) < (?, ?)", cursorTimetableEntry.CreatedAt, cursorID)
	}

	err := query.Limit(limit + 1).Scan(ctx)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("timetable entries not found: %v", err)
	}

	var nextCursor *string
	if len(timetableEntries) > limit {
		timetableEntries = timetableEntries[:limit]
		nextCursor = utils.StringPtr(utils.EncodeCursor(timetableEntries[len(timetableEntries)-1].ID))
	}

	total, _ := queryCount.Count(ctx)
	return timetableEntries, nextCursor, total, nil
}
