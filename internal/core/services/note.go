package services

import (
	"context"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
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
		CreatedBy:   utils.StringPtr(authorID),
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
		err = n.repo.AttachNoteTags(ctx, res.ID, req.TagIDs)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (n *NoteService) GetNote(ctx context.Context, id string) (*domain.Note, error) {
	return n.repo.GetNoteByID(ctx, id)
}

func (n *NoteService) ListNotes(ctx context.Context, filter domain.NoteFilter) ([]*domain.NoteHasTag, int, error) {
	return n.repo.ListNotes(ctx, filter)
}

func (n *NoteService) ListNotesCursor(ctx context.Context, filter domain.NoteFilter, cursor string, limit int) ([]*domain.NoteHasTag, *string, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	return n.repo.ListNotesCursor(ctx, filter, cursor, limit)
}

func (n *NoteService) UpdateNote(ctx context.Context, id string, req domain.UpdateNoteRequest) error {
	noteInfo, err := n.repo.GetNoteByID(ctx, id)
	if err != nil {
		return err
	}

	utils.SetIfNotNil(&noteInfo.Title, req.Title)
	if req.SourceURL != nil {
		noteInfo.SourceUrl = req.SourceURL
	}
	utils.SetIfNotNil(&noteInfo.Status, req.Status)
	if req.HTML != nil {
		noteInfo.HTML = req.HTML
	}
	if req.Lexical != nil {
		noteInfo.Lexical = req.Lexical
	}
	if req.Description != nil {
		noteInfo.Description = req.Description
	}
	noteInfo.UpdatedBy = utils.StringPtr(ctx.Value("user_id").(string))

	return n.repo.UpdateNote(ctx, id, *noteInfo)
}

func (n *NoteService) DeleteNote(ctx context.Context, id string) error {
	n.repo.DetachNoteTags(ctx, id)
	return n.repo.DeleteNote(ctx, id)
}
