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

type FlashcardHandler struct {
	svc services.FlashcardService
}

func NewFlashcardHandler(svc services.FlashcardService) *FlashcardHandler {
	return &FlashcardHandler{svc: svc}
}

// CreateDeck handles POST /v1/flashcards/decks
func (h *FlashcardHandler) CreateDeck(ctx *gin.Context) {
	var req domain.CreateFlashcardDeckRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("CreateDeck: Invalid request payload")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}
	userId, err := http.GetUserID(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("CreateDeck: Failed to get user ID from context")
		HandleError(ctx, domain.ErrorCodeForbidden, nil, "Unauthorized")
		return
	}
	deck, err := h.svc.CreateDeck(ctx, userId, req)
	if err != nil {
		logger.Log.WithError(err).Error("CreateDeck: Failed to create deck")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, deck, "Deck created")
}

// GetDeck handles GET /v1/flashcards/decks/:id
func (h *FlashcardHandler) GetDeck(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("GetDeck: Invalid deck ID")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	deck, err := h.svc.GetDeck(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("GetDeck: Failed to get deck")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, deck, "Success")
}

// ListDecksCursor handles GET /v1/flashcards/decks/cursor with cursor-based pagination
func (h *FlashcardHandler) ListDecksCursor(ctx *gin.Context) {
	cursor := ctx.Query("cursor")
	limit := 10
	if l := ctx.Query("limit"); l != "" {
		if parsed, err := utils.ParseLimit(l); err == nil {
			limit = parsed
		}
	}

	var filter domain.FlashcardDeckFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		logger.Log.WithError(err).Error("ListDecksCursor: Invalid query parameters")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid query parameters")
		return
	}

	userId, err := http.GetUserID(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("ListDecksCursor: Failed to get user ID from context")
		HandleError(ctx, domain.ErrorCodeForbidden, nil, "Unauthorized")
		return
	}

	decks, nextCursor, total, err := h.svc.ListDecksCursor(ctx, userId, filter, cursor, limit)
	if err != nil {
		logger.Log.WithError(err).Error("ListDecksCursor: Failed to list decks")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}

	response := utils.CursorPaginationResponse{
		Items:      decks,
		NextCursor: nextCursor,
		HasMore:    nextCursor != nil,
		Total:      &total,
	}
	HandleSuccess(ctx, response, "Success")
}

// UpdateDeck handles PUT /v1/flashcards/decks/:id
func (h *FlashcardHandler) UpdateDeck(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("UpdateDeck: Invalid deck ID")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	var req domain.UpdateFlashcardDeckRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("UpdateDeck: Invalid request payload")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}

	if err := h.svc.UpdateDeck(ctx, id, req); err != nil {
		logger.Log.WithError(err).Error("UpdateDeck: Failed to update deck")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, nil, "Deck updated")
}

// DeleteDeck handles DELETE /v1/flashcards/decks/:id
func (h *FlashcardHandler) DeleteDeck(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("DeleteDeck: Invalid deck ID")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	if err := h.svc.DeleteDeck(ctx, id); err != nil {
		logger.Log.WithError(err).Error("DeleteDeck: Failed to delete deck")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, nil, "Deck deleted")
}

// GenerateFlashcards handles POST /v1/flashcards/generate — drafts a deck via the AI provider.
func (h *FlashcardHandler) GenerateFlashcards(ctx *gin.Context) {
	var req domain.GenerateFlashcardsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("GenerateFlashcards: Invalid request payload")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}

	draft, err := h.svc.GenerateCards(ctx, req)
	if err != nil {
		logger.Log.WithError(err).Error("GenerateFlashcards: Failed to generate flashcards")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, draft, "Flashcards generated")
}

// CreateCard handles POST /v1/flashcards/decks/:id/cards
func (h *FlashcardHandler) CreateCard(ctx *gin.Context) {
	deckID, err := parseIDParam(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("CreateCard: Invalid deck ID")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	var req domain.CreateFlashcardCardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("CreateCard: Invalid request payload")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}

	card, err := h.svc.CreateCard(ctx, deckID, req)
	if err != nil {
		logger.Log.WithError(err).Error("CreateCard: Failed to create card")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, card, "Card created")
}

// UpdateCard handles PUT /v1/flashcards/cards/:id
func (h *FlashcardHandler) UpdateCard(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("UpdateCard: Invalid card ID")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	var req domain.UpdateFlashcardCardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("UpdateCard: Invalid request payload")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}

	if err := h.svc.UpdateCard(ctx, id, req); err != nil {
		logger.Log.WithError(err).Error("UpdateCard: Failed to update card")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, nil, "Card updated")
}

// DeleteCard handles DELETE /v1/flashcards/cards/:id
func (h *FlashcardHandler) DeleteCard(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("DeleteCard: Invalid card ID")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	if err := h.svc.DeleteCard(ctx, id); err != nil {
		logger.Log.WithError(err).Error("DeleteCard: Failed to delete card")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, nil, "Card deleted")
}

// ListDueCards handles GET /v1/flashcards/decks/:id/study — returns the cards due now for the caller.
func (h *FlashcardHandler) ListDueCards(ctx *gin.Context) {
	deckID, err := parseIDParam(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("ListDueCards: Invalid deck ID")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	limit := 20
	if l := ctx.Query("limit"); l != "" {
		if parsed, err := utils.ParseLimit(l); err == nil {
			limit = parsed
		}
	}
	userId, err := http.GetUserID(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("ListDueCards: Failed to get user ID from context")
		HandleError(ctx, domain.ErrorCodeForbidden, nil, "Unauthorized")
		return
	}

	cards, err := h.svc.ListDueCards(ctx, userId, deckID, limit)
	if err != nil {
		logger.Log.WithError(err).Error("ListDueCards: Failed to list due cards")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, cards, "Success")
}

// SubmitReview handles POST /v1/flashcards/cards/:id/review
func (h *FlashcardHandler) SubmitReview(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("SubmitReview: Invalid card ID")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	var req domain.SubmitFlashcardReviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("SubmitReview: Invalid request payload")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}
	userId, err := http.GetUserID(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("SubmitReview: Failed to get user ID from context")
		HandleError(ctx, domain.ErrorCodeForbidden, nil, "Unauthorized")
		return
	}

	result, err := h.svc.SubmitReview(ctx, userId, id, req)
	if err != nil {
		logger.Log.WithError(err).Error("SubmitReview: Failed to submit review")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, result, "Review submitted")
}
