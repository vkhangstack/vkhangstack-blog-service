package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/http"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/services"
	"github.com/vkhangstack/hexagonal-architecture/internal/logger"
)

type TagHandler struct {
	svc *services.TagService
}

func NewTagHandler(svc *services.TagService) *TagHandler {
	return &TagHandler{svc: svc}
}

// CreateTag handles POST /v1/cms/tags
func (h *TagHandler) CreateTag(ctx *gin.Context) {
	var req domain.CreateTagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("CreateTag: Invalid request payload")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	tag, err := h.svc.CreateTag(ctx, ctx.MustGet("user_id").(string), req)
	if err != nil {
		logger.Log.WithError(err).Error("CreateTag: Failed to create tag")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	HandleSuccess(ctx, tag, "Tag created")
}

// ListTags handles GET /v1/cms/tags
func (h *TagHandler) ListTags(ctx *gin.Context) {
	userID, err := http.GetUserID(ctx)
	tagType := ctx.Query("type")
	if err != nil {
		logger.Log.WithError(err).Error("ListTags: Failed to get user ID")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	tags, err := h.svc.ListTags(ctx, userID, tagType)
	if err != nil {
		logger.Log.WithError(err).Error("ListTags: Failed to list tags")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	HandleSuccess(ctx, tags, "Success")
}
