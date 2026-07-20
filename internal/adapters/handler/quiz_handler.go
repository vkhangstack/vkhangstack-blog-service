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

type QuizHandler struct {
	svc services.QuizService
}

func NewQuizHandler(svc services.QuizService) *QuizHandler {
	return &QuizHandler{svc: svc}
}

// CreateQuiz handles POST /v1/quizzes
func (h *QuizHandler) CreateQuiz(ctx *gin.Context) {
	var req domain.CreateQuizRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("CreateQuiz: Invalid request payload")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}
	userId, err := http.GetUserID(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("CreateQuiz: Failed to get user ID from context")
		HandleError(ctx, domain.ErrorCodeForbidden, nil, "Unauthorized")
		return
	}
	quiz, err := h.svc.CreateQuiz(ctx, userId, req)
	if err != nil {
		logger.Log.WithError(err).Error("CreateQuiz: Failed to create quiz")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, quiz, "Quiz created")
}

// GetQuiz handles GET /v1/quizzes/:id
func (h *QuizHandler) GetQuiz(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("GetQuiz: Invalid quiz ID")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	quiz, err := h.svc.GetQuiz(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("GetQuiz: Failed to get quiz")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, quiz, "Success")
}

// ListQuizzesCursor handles GET /v1/quizzes/cursor with cursor-based pagination
func (h *QuizHandler) ListQuizzesCursor(ctx *gin.Context) {
	cursor := ctx.Query("cursor")
	limit := 10
	if l := ctx.Query("limit"); l != "" {
		if parsed, err := utils.ParseLimit(l); err == nil {
			limit = parsed
		}
	}

	var filter domain.QuizFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		logger.Log.WithError(err).Error("ListQuizzesCursor: Invalid query parameters")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid query parameters")
		return
	}

	quizzes, nextCursor, total, err := h.svc.ListQuizzesCursor(ctx, filter, cursor, limit)
	if err != nil {
		logger.Log.WithError(err).Error("ListQuizzesCursor: Failed to list quizzes")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}

	response := utils.CursorPaginationResponse{
		Items:      quizzes,
		NextCursor: nextCursor,
		HasMore:    nextCursor != nil,
		Total:      &total,
	}
	HandleSuccess(ctx, response, "Success")
}

// UpdateQuiz handles PUT /v1/quizzes/:id
func (h *QuizHandler) UpdateQuiz(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("UpdateQuiz: Invalid quiz ID")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	var req domain.UpdateQuizRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("UpdateQuiz: Invalid request payload")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}

	err = h.svc.UpdateQuiz(ctx, id, req)
	if err != nil {
		logger.Log.WithError(err).Error("UpdateQuiz: Failed to update quiz")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, nil, "Quiz updated")
}

// DeleteQuiz handles DELETE /v1/quizzes/:id
func (h *QuizHandler) DeleteQuiz(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("DeleteQuiz: Invalid quiz ID")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	if err := h.svc.DeleteQuiz(ctx, id); err != nil {
		logger.Log.WithError(err).Error("DeleteQuiz: Failed to delete quiz")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, nil, "Quiz deleted")
}

// GenerateQuiz handles POST /v1/quizzes/generate — drafts a quiz via the AI provider.
func (h *QuizHandler) GenerateQuiz(ctx *gin.Context) {
	var req domain.GenerateQuizRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("GenerateQuiz: Invalid request payload")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}

	draft, err := h.svc.GenerateQuiz(ctx, req)
	if err != nil {
		logger.Log.WithError(err).Error("GenerateQuiz: Failed to generate quiz")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, draft, "Quiz generated")
}

// SubmitAttempt handles POST /v1/quizzes/:id/attempts
func (h *QuizHandler) SubmitAttempt(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("SubmitAttempt: Invalid quiz ID")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	var req domain.SubmitQuizAttemptRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("SubmitAttempt: Invalid request payload")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}
	userId, err := http.GetUserID(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("SubmitAttempt: Failed to get user ID from context")
		HandleError(ctx, domain.ErrorCodeForbidden, nil, "Unauthorized")
		return
	}

	result, err := h.svc.SubmitAttempt(ctx, id, userId, req)
	if err != nil {
		logger.Log.WithError(err).Error("SubmitAttempt: Failed to submit attempt")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, result, "Attempt submitted")
}

// ListAttemptsCursor handles GET /v1/quizzes/:id/attempts/cursor
func (h *QuizHandler) ListAttemptsCursor(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("ListAttemptsCursor: Invalid quiz ID")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	cursor := ctx.Query("cursor")
	limit := 10
	if l := ctx.Query("limit"); l != "" {
		if parsed, err := utils.ParseLimit(l); err == nil {
			limit = parsed
		}
	}

	attempts, nextCursor, total, err := h.svc.ListQuizAttemptsCursor(ctx, id, cursor, limit)
	if err != nil {
		logger.Log.WithError(err).Error("ListAttemptsCursor: Failed to list attempts")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}

	response := utils.CursorPaginationResponse{
		Items:      attempts,
		NextCursor: nextCursor,
		HasMore:    nextCursor != nil,
		Total:      &total,
	}
	HandleSuccess(ctx, response, "Success")
}
