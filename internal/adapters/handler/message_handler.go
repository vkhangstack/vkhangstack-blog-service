package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/http"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/services"
)

type MessageHandler struct {
	svc services.MessengerService
}

func NewMessageHandler(MessengerService services.MessengerService) *MessageHandler {
	return &MessageHandler{
		svc: MessengerService,
	}
}

func (h *MessageHandler) CreateMessage(ctx *gin.Context) {
	// Validate token
	userID, err := http.GetUserID(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}

	var message domain.Message
	message.UserID = userID

	if err := ctx.ShouldBindJSON(&message); err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}

	err = h.svc.CreateMessage(userID, message)
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}

	HandleSuccess(ctx, nil, "Message created successfully")
}

func (h *MessageHandler) ReadMessage(ctx *gin.Context) {
	id := ctx.Param("id")
	message, err := h.svc.ReadMessage(id)

	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	HandleSuccess(ctx, message, "Message read successfully")
}

func (h *MessageHandler) ReadMessages(ctx *gin.Context) {

	messages, err := h.svc.ReadMessages()

	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}

	HandleSuccess(ctx, messages, "Messages read successfully")
}

func (h *MessageHandler) UpdateMessage(ctx *gin.Context) {
	userID, err := http.GetUserID(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	// check if userID match with message.UserID
	id := ctx.Param("id")
	msg, err := h.svc.ReadMessage(id)
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	if msg.UserID != userID {
		HandleError(ctx, domain.ErrorCodeUnAuthorization, nil, fmt.Errorf("you are not authorized to update this message").Error())
		return
	}

	var message domain.Message
	if err := ctx.ShouldBindJSON(&message); err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}

	err = h.svc.UpdateMessage(id, message)
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	HandleSuccess(ctx, nil, "Message updated successfully")
}

func (h *MessageHandler) DeleteMessage(ctx *gin.Context) {

	userID, err := http.GetUserID(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeUnAuthorization, nil, err.Error())
		return
	}

	// check if userID match with message.UserID
	id := ctx.Param("id")
	message, err := h.svc.ReadMessage(id)
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	if message.UserID != userID {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, fmt.Errorf("you are not authorized to delete this message").Error())
		return
	}

	err = h.svc.DeleteMessage(id)
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	HandleSuccess(ctx, nil, "Message deleted successfully")
}
