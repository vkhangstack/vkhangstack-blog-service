package services

import (
	"context"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
)

type NoteService struct {
	repo ports.NoteRepository
}

func NewNoteService(repo ports.NoteRepository) *NoteService {
	return &NoteService{
		repo: repo,
	}
}

func (n *NoteService) CreateNote(ctx context.Context, authorID string, req domain.CreateNoteRequest) (*domain.Note, error) {
	note := domain.Note{
		CreatedBy:   authorID,
		Title:       req.Title,
		SourceUrl:   req.SourceURL,
		Status:      req.Status,
		HTML:        req.HTML,
		Lexical:     req.Lexical,
		Description: req.Description,
	}

	res, err := n.repo.CreateNote(ctx, note)
	if err != nil {
		return nil, err
	}
	if len(req.TagIDs) > 0 {
		err = n.repo.AttachNoteTags(res.ID, req.TagIDs)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (n *NoteService) GetNote(ctx context.Context, id string) (*domain.Note, error) {
	return n.repo.GetNoteByID(ctx, id)
}

func (n *NoteService) ListNotes(ctx context.Context, filter domain.NoteFilter) ([]*domain.Note, int, error) {
	return n.repo.ListNotes(ctx, filter)
}

func (n *NoteService) ListNotesCursor(ctx context.Context, filter domain.NoteFilter, cursor string, limit int) ([]*domain.Note, *string, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	return n.repo.ListNotesCursor(ctx, filter, cursor, limit)
}

func (n *NoteService) UpdateNote(ctx context.Context, id string, req domain.UpdateNoteRequest) (*domain.Note, error) {
	note := domain.Note{
		ID:          id,
		Title:       *req.Title,
		SourceUrl:   req.SourceURL,
		Status:      *req.Status,
		HTML:        req.HTML,
		Lexical:     req.Lexical,
		Description: req.Description,
		UpdatedBy:   ctx.Value("user_id").(string),
	}
	return n.repo.UpdateNote(ctx, id, note)
}

func (n *NoteService) DeleteNote(ctx context.Context, id string) error {
	n.repo.DetachNoteTags(id)
	return n.repo.DeleteNote(ctx, id)
}
