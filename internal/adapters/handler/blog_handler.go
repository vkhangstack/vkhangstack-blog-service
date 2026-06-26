package handler

import (
	"github.com/gin-gonic/gin"
	customhttp "github.com/vkhangstack/hexagonal-architecture/internal/adapters/http"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/services"
)

type BlogHandler struct {
	categorySvc *services.BlogCategoryService
	postSvc     *services.BlogPostService
}

func NewBlogHandler(categorySvc *services.BlogCategoryService, postSvc *services.BlogPostService) *BlogHandler {
	return &BlogHandler{categorySvc: categorySvc, postSvc: postSvc}
}

// CreateCategory handles POST /v1/cms/categories
func (h *BlogHandler) CreateCategory(ctx *gin.Context) {
	var req domain.CreateBlogCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	category, err := h.categorySvc.CreateCategory(req)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeBlogCategoryNotFound, nil, err.Error())
		return
	}
	HandleSuccess(ctx, category, "Category created")
}

// GetCategory handles GET /v1/cms/categories/:id
func (h *BlogHandler) GetCategory(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	category, err := h.categorySvc.GetCategory(id)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeBlogCategoryNotFound, nil, err.Error())
		return
	}
	HandleSuccess(ctx, category, "Success")
}

// ListCategories handles GET /v1/cms/categories and GET /v1/blog/categories
func (h *BlogHandler) ListCategories(ctx *gin.Context) {
	categories, err := h.categorySvc.ListCategories()
	if err != nil {
		HandleError(ctx, domain.ErrorCodeBlogCategoryNotFound, nil, err.Error())
		return
	}
	HandleSuccess(ctx, categories, "Success")
}

// UpdateCategory handles PUT /v1/cms/categories/:id
func (h *BlogHandler) UpdateCategory(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	var req domain.UpdateBlogCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	category, err := h.categorySvc.UpdateCategory(id, req)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeBlogCategoryNotFound, nil, err.Error())
		return
	}
	HandleSuccess(ctx, category, "Category updated")
}

// DeleteCategory handles DELETE /v1/cms/categories/:id
func (h *BlogHandler) DeleteCategory(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	if err := h.categorySvc.DeleteCategory(id); err != nil {
		HandleError(ctx, domain.ErrorCodeBlogCategoryNotFound, nil, err.Error())
		return
	}
	HandleSuccess(ctx, nil, "Category deleted")
}

// CreatePost handles POST /v1/cms/posts
func (h *BlogHandler) CreatePost(ctx *gin.Context) {
	authorID, err := getAuthorID(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeUnAuthorization, nil, err.Error())
		return
	}
	var req domain.CreateBlogPostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	post, err := h.postSvc.CreatePost(authorID, req)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeBlogPostNotFound, nil, err.Error())
		return
	}
	HandleSuccess(ctx, post, "Post created")
}

// GetPost handles GET /v1/cms/posts/:id
func (h *BlogHandler) GetPost(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	post, err := h.postSvc.GetPost(id)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeBlogPostNotFound, nil, err.Error())
		return
	}
	HandleSuccess(ctx, post, "Success")
}

// ListPosts handles GET /v1/cms/posts
func (h *BlogHandler) ListPosts(ctx *gin.Context) {
	var filter domain.BlogPostFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	posts, total, err := h.postSvc.ListPosts(filter)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeBlogPostNotFound, nil, err.Error())
		return
	}
	HandleSuccess(ctx, domain.BlogPostListResponse{Total: total, Posts: posts}, "Success")
}

// UpdatePost handles PUT /v1/cms/posts/:id
func (h *BlogHandler) UpdatePost(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	var req domain.UpdateBlogPostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	post, err := h.postSvc.UpdatePost(id, req)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeBlogPostNotFound, nil, err.Error())
		return
	}
	HandleSuccess(ctx, post, "Post updated")
}

// DeletePost handles DELETE /v1/cms/posts/:id
func (h *BlogHandler) DeletePost(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	if err := h.postSvc.DeletePost(id); err != nil {
		HandleError(ctx, domain.ErrorCodeBlogPostNotFound, nil, err.Error())
		return
	}
	HandleSuccess(ctx, nil, "Post deleted")
}

// PublishPost handles POST /v1/cms/posts/:id/publish
func (h *BlogHandler) PublishPost(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	post, err := h.postSvc.PublishPost(id)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeBlogPostNotFound, nil, err.Error())
		return
	}
	HandleSuccess(ctx, post, "Post published")
}

// GetPostBySlug handles GET /v1/blog/posts/:slug
func (h *BlogHandler) GetPostBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	post, err := h.postSvc.GetPostBySlug(slug)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeBlogPostNotFound, nil, err.Error())
		return
	}
	HandleSuccess(ctx, post, "Success")
}

// ListPublishedPosts handles GET /v1/blog/posts (public, published only)
func (h *BlogHandler) ListPublishedPosts(ctx *gin.Context) {
	var filter domain.BlogPostFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	filter.Status = string(domain.PostStatusPublished)
	posts, total, err := h.postSvc.ListPosts(filter)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeBlogPostNotFound, nil, err.Error())
		return
	}
	HandleSuccess(ctx, domain.BlogPostListResponse{Total: total, Posts: posts}, "Success")
}

func parseIDParam(ctx *gin.Context) (string, error) {
	return ctx.Param("id"), nil
}

func getAuthorID(ctx *gin.Context) (string, error) {
	idStr, err := customhttp.GetUserID(ctx)
	if err != nil {
		return "", err
	}
	return idStr, nil
}
