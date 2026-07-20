package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/validate"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/services"
	"github.com/vkhangstack/hexagonal-architecture/internal/logger"
)

type ChatHandler struct {
	svc *services.ChatService
}

func NewChatHandler(svc *services.ChatService) *ChatHandler {
	return &ChatHandler{svc: svc}
}

// ListConversations handles GET /v1/chats — the backoffice Chats inbox (one row per Zalo chat).
func (h *ChatHandler) ListConversations(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	limit, _ := strconv.Atoi(ctx.Query("limit"))

	conversations, total, err := h.svc.ListConversations(ctx, page, limit)
	if err != nil {
		logger.Log.WithError(err).Error("ListConversations: Failed to list conversations")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, domain.BotChatConversationListResponse{
		Total:         total,
		Conversations: conversations,
	}, "Success")
}

// ListMessages handles GET /v1/chats/:chatID/messages — one conversation's message history.
func (h *ChatHandler) ListMessages(ctx *gin.Context) {
	chatID := ctx.Param("chatID")
	page, _ := strconv.Atoi(ctx.Query("page"))
	limit, _ := strconv.Atoi(ctx.Query("limit"))

	messages, total, err := h.svc.ListMessages(ctx, chatID, page, limit)
	if err != nil {
		logger.Log.WithError(err).Error("ListMessages: Failed to list messages")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, domain.BotChatMessageListResponse{
		Total:    total,
		Messages: messages,
	}, "Success")
}

// SendReply handles POST /v1/chats/:chatID/reply — an admin pushes a manual reply to the user
// through the bot.
func (h *ChatHandler) SendReply(ctx *gin.Context) {
	chatID := ctx.Param("chatID")

	var req domain.ChatReplyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("SendReply: Invalid request payload")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}

	if err := h.svc.SendReply(ctx, chatID, req.Text); err != nil {
		logger.Log.WithError(err).Error("SendReply: Failed to send reply")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, nil, "Reply sent")
}
