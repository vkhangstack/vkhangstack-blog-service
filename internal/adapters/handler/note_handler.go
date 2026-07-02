package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/http"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/validate"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/services"
	"github.com/vkhangstack/hexagonal-architecture/internal/logger"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

type NoteHandler struct {
	svc services.NoteService
}

func NewNoteHandler(svc services.NoteService) *NoteHandler {
	return &NoteHandler{svc: svc}
}

// CreateNote handles POST /v1/cms/notes
func (h *NoteHandler) CreateNote(ctx *gin.Context) {
	var req domain.CreateNoteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("CreateNote: Invalid request payload")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}
	userId, err := http.GetUserID(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("CreateNote: Failed to get user ID from context")
		HandleError(ctx, domain.ErrorCodeForbidden, nil, "Unauthorized")
		return
	}
	note, err := h.svc.CreateNote(ctx, userId, req)
	if err != nil {
		logger.Log.WithError(err).Error("CreateNote: Failed to create note")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, note, "Note created")
}

// GetNote handles GET /v1/cms/notes/:id
func (h *NoteHandler) GetNote(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("GetNote: Invalid note ID")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	note, err := h.svc.GetNote(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("GetNote: Failed to get note")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, note, "Success")
}

// ListNotes handles GET /v1/cms/notes
func (h *NoteHandler) ListNotes(ctx *gin.Context) {
	var filter domain.NoteFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		logger.Log.WithError(err).Error("ListNotes: Invalid query parameters")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid query parameters")
		return
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	notes, total, err := h.svc.ListNotes(ctx, filter)
	if err != nil {
		logger.Log.WithError(err).Error("ListNotes: Failed to list notes")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}

	response := domain.NoteListResponse{
		Total: total,
		Notes: notes,
	}
	HandleSuccess(ctx, response, "Success")
}

// ListNotesCursor handles GET /v1/cms/notes/cursor with cursor-based pagination
func (h *NoteHandler) ListNotesCursor(ctx *gin.Context) {
	cursor := ctx.Query("cursor")
	limit := 10
	if l := ctx.Query("limit"); l != "" {
		if parsed, err := parseNoteLimit(l); err == nil {
			limit = parsed
		}
	}

	var filter domain.NoteFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		logger.Log.WithError(err).Error("ListNotesCursor: Invalid query parameters")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid query parameters")
		return
	}

	notes, nextCursor, total, err := h.svc.ListNotesCursor(ctx, filter, cursor, limit)
	if err != nil {
		logger.Log.WithError(err).Error("ListNotesCursor: Failed to list notes")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}

	response := utils.CursorPaginationResponse{
		Items:      notes,
		NextCursor: nextCursor,
		HasMore:    nextCursor != nil,
		Total:      &total,
	}
	HandleSuccess(ctx, response, "Success")
}

// UpdateNote handles PUT /v1/cms/notes/:id
func (h *NoteHandler) UpdateNote(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("UpdateNote: Invalid note ID")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	var req domain.UpdateNoteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("UpdateNote: Invalid request payload")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}

	note, err := h.svc.UpdateNote(ctx, id, req)
	if err != nil {
		logger.Log.WithError(err).Error("UpdateNote: Failed to update note")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, note, "Note updated")
}

// DeleteNote handles DELETE /v1/cms/notes/:id
func (h *NoteHandler) DeleteNote(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("DeleteNote: Invalid note ID")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	if err := h.svc.DeleteNote(ctx, id); err != nil {
		logger.Log.WithError(err).Error("DeleteNote: Failed to delete note")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, nil, "Note deleted")
}

func parseNoteLimit(limitStr string) (int, error) {
	return utils.ParseLimit(limitStr)
}
