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

	note := domain.Note{
		UpdatedBy: ctx.Value("user_id").(string),
	}
	utils.SetIfNotNil(&note.Title, req.Title)
	utils.SetIfNotNil(&note.SourceUrl, &req.SourceURL)
	utils.SetIfNotNil(&note.Status, req.Status)
	utils.SetIfNotNil(&note.HTML, &req.HTML)
	utils.SetIfNotNil(&note.Lexical, &req.Lexical)
	utils.SetIfNotNil(&note.Description, &req.Description)
	utils.SetIfNotNil(&note.CreatedBy, &noteInfo.CreatedBy)

	return n.repo.UpdateNote(ctx, id, note)
}

func (n *NoteService) DeleteNote(ctx context.Context, id string) error {
	n.repo.DetachNoteTags(ctx, id)
	return n.repo.DeleteNote(ctx, id)
}
