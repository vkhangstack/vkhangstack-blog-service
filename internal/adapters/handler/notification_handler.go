package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/http"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/validate"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/services"
	"github.com/vkhangstack/hexagonal-architecture/internal/logger"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

type NotificationHandler struct {
	svc        *services.NotificationService
	settingSvc *services.NotificationSettingService
	zaloBot    ports.ZaloBotClient
}

func NewNotificationHandler(svc *services.NotificationService, settingSvc *services.NotificationSettingService, zaloBot ports.ZaloBotClient) *NotificationHandler {
	return &NotificationHandler{svc: svc, settingSvc: settingSvc, zaloBot: zaloBot}
}

// CreateNotification handles POST /v1/notifications
func (h *NotificationHandler) CreateNotification(ctx *gin.Context) {
	userID, err := http.GetUserID(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeForbidden, nil, "Unauthorized")
		return
	}

	var req domain.CreateNotificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("CreateNotification: Invalid request payload")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}

	notification, err := h.svc.CreateNotification(ctx, userID, req)
	if err != nil {
		logger.Log.WithError(err).Error("CreateNotification: Failed to create notification")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, notification, "Notification created")
}

// ListNotifications handles GET /v1/notifications
func (h *NotificationHandler) ListNotifications(ctx *gin.Context) {
	userID, err := http.GetUserID(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeForbidden, nil, "Unauthorized")
		return
	}

	var filter domain.NotificationFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		logger.Log.WithError(err).Error("ListNotifications: Invalid query parameters")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid query parameters")
		return
	}

	notifications, total, err := h.svc.ListNotifications(ctx, userID, filter)
	if err != nil {
		logger.Log.WithError(err).Error("ListNotifications: Failed to list notifications")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}

	response := domain.NotificationListResponse{
		Total:         total,
		Notifications: notifications,
	}
	HandleSuccess(ctx, response, "Success")
}

// ListNotificationsCursor handles GET /v1/notifications/cursor with cursor-based pagination
func (h *NotificationHandler) ListNotificationsCursor(ctx *gin.Context) {
	userID, err := http.GetUserID(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeForbidden, nil, "Unauthorized")
		return
	}

	cursor := ctx.Query("cursor")
	limit := 10
	if l := ctx.Query("limit"); l != "" {
		if parsed, err := utils.ParseLimit(l); err == nil {
			limit = parsed
		}
	}

	var filter domain.NotificationFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		logger.Log.WithError(err).Error("ListNotificationsCursor: Invalid query parameters")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid query parameters")
		return
	}

	notifications, nextCursor, total, err := h.svc.ListNotificationsCursor(ctx, userID, filter, cursor, limit)
	if err != nil {
		logger.Log.WithError(err).Error("ListNotificationsCursor: Failed to list notifications")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}

	response := utils.CursorPaginationResponse{
		Items:      notifications,
		NextCursor: nextCursor,
		HasMore:    nextCursor != nil,
		Total:      &total,
	}
	HandleSuccess(ctx, response, "Success")
}

// GetNotification handles GET /v1/notifications/:id
func (h *NotificationHandler) GetNotification(ctx *gin.Context) {
	userID, err := http.GetUserID(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeForbidden, nil, "Unauthorized")
		return
	}

	id, err := parseIDParam(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}

	notification, err := h.svc.GetNotification(ctx, userID, id)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeNotificationNotFound, nil, err.Error())
		return
	}
	HandleSuccess(ctx, notification, "Success")
}

// UpdateNotification handles PUT /v1/notifications/:id (e.g. marking as read)
func (h *NotificationHandler) UpdateNotification(ctx *gin.Context) {
	userID, err := http.GetUserID(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeForbidden, nil, "Unauthorized")
		return
	}

	id, err := parseIDParam(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}

	var req domain.UpdateNotificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("UpdateNotification: Invalid request payload")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}

	if err := h.svc.UpdateNotification(ctx, userID, id, req); err != nil {
		logger.Log.WithError(err).Error("UpdateNotification: Failed to update notification")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, nil, "Notification updated")
}

// DeleteNotification handles DELETE /v1/notifications/:id
func (h *NotificationHandler) DeleteNotification(ctx *gin.Context) {
	userID, err := http.GetUserID(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeForbidden, nil, "Unauthorized")
		return
	}

	id, err := parseIDParam(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}

	if err := h.svc.DeleteNotification(ctx, userID, id); err != nil {
		logger.Log.WithError(err).Error("DeleteNotification: Failed to delete notification")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, nil, "Notification deleted")
}

// GetNotificationSetting handles GET /v1/notifications/settings
func (h *NotificationHandler) GetNotificationSetting(ctx *gin.Context) {
	userID, err := http.GetUserID(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeForbidden, nil, "Unauthorized")
		return
	}

	setting, err := h.settingSvc.GetNotificationSetting(ctx, userID)
	if err != nil {
		logger.Log.WithError(err).Error("GetNotificationSetting: Failed to get notification setting")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, setting, "Success")
}

// UpdateNotificationSetting handles PUT /v1/notifications/settings
// Body selects which channels are enabled and the token/address for each.
func (h *NotificationHandler) UpdateNotificationSetting(ctx *gin.Context) {
	userID, err := http.GetUserID(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeForbidden, nil, "Unauthorized")
		return
	}

	var req domain.UpdateNotificationSettingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("UpdateNotificationSetting: Invalid request payload")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}

	setting, err := h.settingSvc.UpdateNotificationSetting(ctx, userID, req)
	if err != nil {
		logger.Log.WithError(err).Error("UpdateNotificationSetting: Failed to update notification setting")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, setting, "Notification setting updated")
}

// RequestChannelVerification handles POST /v1/notifications/settings/channels/verify
// Generates a code the client shows the user to copy and send from within the target
// channel's own app (e.g. as a Zalo message to our bot) to link that channel.
func (h *NotificationHandler) RequestChannelVerification(ctx *gin.Context) {
	userID, err := http.GetUserID(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeForbidden, nil, "Unauthorized")
		return
	}

	var req domain.CreateChannelVerificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("RequestChannelVerification: Invalid request payload")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}

	verification, err := h.settingSvc.RequestChannelVerification(ctx, userID, req)
	if err != nil {
		logger.Log.WithError(err).Error("RequestChannelVerification: Failed to generate verification code")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	HandleSuccess(ctx, verification, "Verification code generated")
}
