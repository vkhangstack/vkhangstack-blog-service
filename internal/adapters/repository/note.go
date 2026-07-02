package repository

import (
	"context"
	"fmt"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

func (u *DB) CreateNote(ctx context.Context, note domain.Note) (*domain.Note, error) {
	note.ID = u.snowflakeNode.GenerateID()
	_, err := u.db.NewInsert().Model(&note).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("note not saved: %v", err)
	}

	return &note, nil
}

func (u *DB) GetNoteByID(ctx context.Context, id string) (*domain.Note, error) {
	note := &domain.Note{}
	err := u.db.NewSelect().Model(note).Where("n.id = ?", id).Limit(1).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("note not found: %v", err)
	}
	return note, nil
}

func (u *DB) ListNotes(ctx context.Context, filter domain.NoteFilter) ([]*domain.Note, int, error) {
	var notes []*domain.Note
	query := u.db.NewSelect().Model(&notes)

	if filter.Status != nil {
		query = query.Where("n.status = ?", *filter.Status)
	}

	if filter.CreatedBy != nil && *filter.CreatedBy != "" {
		query = query.Where("n.created_by = ?", *filter.CreatedBy)
	}

	if filter.Title != nil && *filter.Title != "" {
		query = query.Where("n.title ILIKE ?", "%"+*filter.Title+"%")
	}

	total, err := query.Order("n.created_at DESC").Offset((filter.Page - 1) * filter.Limit).Limit(filter.Limit).ScanAndCount(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list notes: %v", err)
	}

	return notes, total, nil
}

func (u *DB) UpdateNote(ctx context.Context, id string, updates domain.Note) (*domain.Note, error) {
	updates.ID = id
	_, err := u.db.NewUpdate().Model(&updates).WherePK().Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("note not updated: %v", err)
	}
	return u.GetNoteByID(ctx, id)
}

func (u *DB) DeleteNote(ctx context.Context, id string) error {
	note := &domain.Note{ID: id}
	_, err := u.db.NewDelete().Model(note).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("note not deleted: %v", err)
	}
	return nil
}

func (u *DB) ListNotesCursor(ctx context.Context, filter domain.NoteFilter, cursor string, limit int) ([]*domain.Note, *string, int, error) {
	var cursorID string
	if cursor != "" {
		id, err := utils.DecodeCursor(cursor)
		if err != nil {
			return nil, nil, 0, err
		}
		cursorID = id
	}

	var notes []*domain.Note
	query := u.db.NewSelect().Model(&notes)

	if filter.Status != nil {
		query = query.Where("n.status = ?", *filter.Status)
	}
	if filter.CreatedBy != nil && *filter.CreatedBy != "" {
		query = query.Where("n.created_by = ?", *filter.CreatedBy)
	}
	if filter.Title != nil && *filter.Title != "" {
		query = query.Where("n.title ILIKE ?", "%"+*filter.Title+"%")
	}

	queryCount := query.Clone()
	query = query.Order("n.created_at DESC", "n.id DESC")

	if cursorID != "" {
		cursorNote := &domain.Note{}
		err := u.db.NewSelect().Model(cursorNote).Where("n.id = ?", cursorID).Limit(1).Scan(ctx)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("cursor note not found: %v", err)
		}
		query = query.Where("(n.created_at, n.id) < (?, ?)", cursorNote.CreatedAt, cursorID)
	}

	err := query.Limit(limit + 1).Scan(ctx)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("notes not found: %v", err)
	}

	var nextCursor *string
	if len(notes) > limit {
		notes = notes[:limit]
		nextCursor = utils.StringPtr(utils.EncodeCursor(notes[len(notes)-1].ID))
	}

	total, _ := queryCount.Count(ctx)
	return notes, nextCursor, total, nil
}

func (u *DB) AttachNoteTags(noteID string, tagIDs []string) error {
	if len(tagIDs) == 0 {
		return nil
	}
	ctx := context.Background()
	joins := make([]domain.NoteTag, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		joins = append(joins, domain.NoteTag{NoteID: noteID, TagID: tagID})
	}
	_, err := u.db.NewInsert().Model(&joins).On("CONFLICT DO NOTHING").Exec(ctx)
	return err
}

func (u *DB) DetachNoteTags(noteID string) error {
	ctx := context.Background()
	_, err := u.db.NewDelete().Model((*domain.NoteTag)(nil)).Where("note_id = ?", noteID).Exec(ctx)
	return err
}
