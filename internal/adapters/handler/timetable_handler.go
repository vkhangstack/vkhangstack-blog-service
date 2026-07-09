package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/http"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/validate"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/services"
	"github.com/vkhangstack/hexagonal-architecture/internal/logger"
)

type TimetableHandler struct {
	timetableSvc *services.TimetableService
}

func NewTimetableHandler(timetableSvc *services.TimetableService) *TimetableHandler {
	return &TimetableHandler{timetableSvc: timetableSvc}
}

// CreateTimetableEntry handles the creation of a new timetable entry.
func (h *TimetableHandler) CreateTimetableEntry(ctx *gin.Context) {
	var req domain.CreateTimetableEntryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("Failed to create timetable entry")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}
	userID, err := http.GetUserID(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to get user ID")
		HandleError(ctx, domain.ErrorCodeForbidden, err.Error(), "Failed to get user ID")
		return
	}
	timetable, err := h.timetableSvc.CreateTimetableEntry(ctx.Request.Context(), userID, req)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to create timetable entry")
		HandleError(ctx, domain.ErrorCodeInternalServerError, err.Error(), "Failed to create timetable entry")
		return
	}
	HandleSuccess(ctx, timetable, "Timetable entry created successfully")
}

// GetTimetableEntry handles the retrieval of a timetable entry by ID.
func (h *TimetableHandler) GetTimetableEntry(ctx *gin.Context) {
	timetableID, err := parseIDParam(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to parse timetable ID")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, "Invalid timetable ID")
		return
	}
	userID, err := http.GetUserID(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to get user ID")
		HandleError(ctx, domain.ErrorCodeForbidden, nil, "Failed to get user ID")
		return
	}
	timetable, err := h.timetableSvc.GetTimetableEntry(ctx.Request.Context(), userID, timetableID)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to get timetable entry")
		HandleError(ctx, domain.ErrorCodeInternalServerError, err.Error(), "Failed to get timetable entry")
		return
	}
	HandleSuccess(ctx, timetable, "Timetable entry retrieved successfully")
}

// ListTimetableEntries handles the retrieval of timetable entries with optional filtering.
func (h *TimetableHandler) ListTimetableEntries(ctx *gin.Context) {
	var filter domain.TimetableEntryFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		logger.Log.WithError(err).Error("Failed to bind query parameters")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid query parameters")
		return
	}
	userID, err := http.GetUserID(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to get user ID")
		HandleError(ctx, domain.ErrorCodeForbidden, nil, "Failed to get user ID")
		return
	}
	timetables, total, err := h.timetableSvc.ListTimetableEntries(ctx.Request.Context(), userID, filter)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to list timetable entries")
		HandleError(ctx, domain.ErrorCodeInternalServerError, err.Error(), "Failed to list timetable entries")
		return
	}
	HandleSuccess(ctx, gin.H{"timetables": timetables, "total": total}, "Timetable entries retrieved successfully")
}

// ListTimetableEntriesCursor handles the retrieval of timetable entries with cursor-based pagination.
func (h *TimetableHandler) ListTimetableEntriesCursor(ctx *gin.Context) {
	var filter domain.TimetableEntryFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		logger.Log.WithError(err).Error("Failed to bind query parameters")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid query parameters")
		return
	}
	cursor := getCursor(ctx)
	limit, err := getLimit(ctx)

	if err != nil {
		logger.Log.WithError(err).Error("Failed to parse limit parameter")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, "Invalid limit parameter")
		return
	}
	userID, err := http.GetUserID(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to get user ID")
		HandleError(ctx, domain.ErrorCodeForbidden, nil, "Failed to get user ID")
		return
	}
	timetables, nextCursor, total, err := h.timetableSvc.ListTimetableEntriesCursor(ctx.Request.Context(), userID, filter, cursor, limit)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to list timetable entries with cursor")
		HandleError(ctx, domain.ErrorCodeInternalServerError, err.Error(), "Failed to list timetable entries with cursor")
		return
	}
	HandleSuccess(ctx, gin.H{"timetables": timetables, "next_cursor": nextCursor, "total": total}, "Timetable entries retrieved successfully")
}

func (h *TimetableHandler) UpdateTimetableEntry(ctx *gin.Context) {
	timetableID, err := parseIDParam(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to parse timetable ID")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, "Invalid timetable ID")
		return
	}
	var req domain.UpdateTimetableEntryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("Failed to bind request body")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request body")
		return
	}
	userID, err := http.GetUserID(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to get user ID")
		HandleError(ctx, domain.ErrorCodeForbidden, nil, "Failed to get user ID")
		return
	}
	timetable, err := h.timetableSvc.UpdateTimetableEntry(ctx.Request.Context(), userID, timetableID, req)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to update timetable entry")
		HandleError(ctx, domain.ErrorCodeInternalServerError, err.Error(), "Failed to update timetable entry")
		return
	}
	HandleSuccess(ctx, timetable, "Timetable entry updated successfully")
}

func (h *TimetableHandler) DeleteTimetableEntry(ctx *gin.Context) {
	timetableID, err := parseIDParam(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to parse timetable ID")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, "Invalid timetable ID")
		return
	}
	userID, err := http.GetUserID(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to get user ID")
		HandleError(ctx, domain.ErrorCodeForbidden, nil, "Failed to get user ID")
		return
	}
	err = h.timetableSvc.DeleteTimetableEntry(ctx.Request.Context(), userID, timetableID)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to delete timetable entry")
		HandleError(ctx, domain.ErrorCodeInternalServerError, err.Error(), "Failed to delete timetable entry")
		return
	}
	HandleSuccess(ctx, nil, "Timetable entry deleted successfully")
}
