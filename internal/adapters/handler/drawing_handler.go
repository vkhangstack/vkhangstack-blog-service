package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/http"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/services"
	"github.com/vkhangstack/hexagonal-architecture/internal/logger"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

type DrawingHandler struct {
	svc services.DrawingService
}

func NewDrawingHandler(svc services.DrawingService) *DrawingHandler {
	return &DrawingHandler{svc: svc}
}

// CreateDrawing handles POST /v1/cms/drawings
func (h *DrawingHandler) CreateDrawing(ctx *gin.Context) {
	var req domain.CreateDrawingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("CreateDrawing: Invalid request payload")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	userID, err := http.GetUserID(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("CreateDrawing: Failed to get user ID")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	drawing, err := h.svc.CreateDrawing(ctx, userID, req)
	if err != nil {
		logger.Log.WithError(err).Error("CreateDrawing: Failed to create drawing")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	HandleSuccess(ctx, drawing, "Drawing created")
}

// ListDrawings handles GET /v1/cms/drawings
func (h *DrawingHandler) ListDrawings(ctx *gin.Context) {
	var filter domain.DrawingFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		logger.Log.WithError(err).Error("ListDrawings: Invalid query parameters")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	drawings, total, err := h.svc.ListDrawings(ctx, filter)
	if err != nil {
		logger.Log.WithError(err).Error("ListDrawings: Failed to list drawings")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	drawingsResponse := make([]*domain.DrawingResponse, len(drawings))
	for i, drawing := range drawings {
		drawingsResponse[i] = &domain.DrawingResponse{
			ID:        drawing.ID,
			Title:     drawing.Title,
			CreatedAt: drawing.CreatedAt,
			UpdatedAt: drawing.UpdatedAt,
		}
	}
	drawingResponse := &domain.ListDrawingResponse{
		Drawings: drawingsResponse,
		Total:    total,
	}
	HandleSuccess(ctx, drawingResponse, "Success")
}

// UpdateDrawing handles PUT /v1/cms/drawings/:id
func (h *DrawingHandler) UpdateDrawing(ctx *gin.Context) {
	id := ctx.Param("id")
	var req domain.UpdateDrawingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("UpdateDrawing: Invalid request payload")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	err := h.svc.UpdateDrawing(ctx, id, req)
	if err != nil {
		logger.Log.WithError(err).Error("UpdateDrawing: Failed to update drawing")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	HandleSuccess(ctx, nil, "Drawing updated")
}

// DeleteDrawing handles DELETE /v1/cms/drawings/:id
func (h *DrawingHandler) DeleteDrawing(ctx *gin.Context) {
	id := ctx.Param("id")
	err := h.svc.DeleteDrawing(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("DeleteDrawing: Failed to delete drawing")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	HandleSuccess(ctx, nil, "Drawing deleted")
}

// ListDrawingsCursor handles GET /v1/cms/drawings/cursor
func (h *DrawingHandler) ListDrawingsCursor(ctx *gin.Context) {
	var filter domain.DrawingFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		logger.Log.WithError(err).Error("ListDrawingsCursor: Invalid query parameters")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	cursor := ctx.Query("cursor")
	limit := filter.Limit
	drawings, nextCursor, total, err := h.svc.ListDrawingsCursor(ctx, filter, cursor, limit)
	if err != nil {
		logger.Log.WithError(err).Error("ListDrawingsCursor: Failed to list drawings")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	drawingsDataResponse := make([]*domain.DrawingResponse, len(drawings))
	for i, drawing := range drawings {
		drawingsDataResponse[i] = &domain.DrawingResponse{
			ID:        drawing.ID,
			Title:     drawing.Title,
			CreatedAt: drawing.CreatedAt,
			UpdatedAt: drawing.UpdatedAt,
		}
	}
	response := utils.CursorPaginationResponse{
		Items:      drawingsDataResponse,
		NextCursor: nextCursor,
		HasMore:    nextCursor != nil,
		Total:      &total,
	}

	HandleSuccess(ctx, response, "Success")
}

// GetDrawing handles GET /v1/cms/drawings/:id
func (h *DrawingHandler) GetDrawing(ctx *gin.Context) {
	id := ctx.Param("id")
	drawing, err := h.svc.GetDrawing(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("GetDrawing: Failed to get drawing")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	HandleSuccess(ctx, drawing, "Success")
}
