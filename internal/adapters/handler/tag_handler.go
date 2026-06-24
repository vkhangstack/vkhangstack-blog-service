package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/services"
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
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	tag, err := h.svc.CreateTag(req)
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	HandleSuccess(ctx, tag, "Tag created")
}

// ListTags handles GET /v1/cms/tags
func (h *TagHandler) ListTags(ctx *gin.Context) {
	tags, err := h.svc.ListTags()
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	HandleSuccess(ctx, tags, "Success")
}
