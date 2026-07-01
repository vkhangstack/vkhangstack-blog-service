package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	customhttp "github.com/vkhangstack/hexagonal-architecture/internal/adapters/http"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/validate"
	"github.com/vkhangstack/hexagonal-architecture/internal/logger"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/services"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

type BlogHandler struct {
	categorySvc     *services.BlogCategoryService
	postSvc         *services.BlogPostService
	searchEngineSvc *services.SearchEngineService
}

func NewBlogHandler(categorySvc *services.BlogCategoryService, postSvc *services.BlogPostService, searchEngineSvc *services.SearchEngineService) *BlogHandler {
	return &BlogHandler{categorySvc: categorySvc, postSvc: postSvc, searchEngineSvc: searchEngineSvc}
}

// CreateCategory handles POST /v1/cms/categories
func (h *BlogHandler) CreateCategory(ctx *gin.Context) {
	var req domain.CreateBlogCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}
	category, err := h.categorySvc.CreateCategory(req)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to create category")
		HandleError(ctx, domain.ErrorCodeBlogCategoryNotFound, nil, err.Error())
		return
	}
	HandleSuccess(ctx, category, "Category created")
}

// GetCategory handles GET /v1/cms/categories/:id
func (h *BlogHandler) GetCategory(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to parse ID param for GetCategory")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	category, err := h.categorySvc.GetCategory(id)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to get category")
		HandleError(ctx, domain.ErrorCodeBlogCategoryNotFound, nil, err.Error())
		return
	}
	HandleSuccess(ctx, category, "Success")
}

// ListCategories handles GET /v1/cms/categories and GET /v1/blog/categories
func (h *BlogHandler) ListCategories(ctx *gin.Context) {
	categories, err := h.categorySvc.ListCategories()
	if err != nil {
		logger.Log.WithError(err).Error("Failed to list categories")
		HandleError(ctx, domain.ErrorCodeBlogCategoryNotFound, nil, err.Error())
		return
	}
	HandleSuccess(ctx, categories, "Success")
}

// ListCategoriesCursor handles GET /v1/cms/categories/cursor with cursor-based pagination
func (h *BlogHandler) ListCategoriesCursor(ctx *gin.Context) {
	cursor := ctx.Query("cursor")
	limit := 20
	if l := ctx.Query("limit"); l != "" {
		if parsed, err := parseLimit(l); err == nil {
			limit = parsed
		}
	}

	categories, nextCursor, err := h.categorySvc.ListCategoriesCursor(cursor, limit)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to list categories with cursor")
		HandleError(ctx, domain.ErrorCodeBlogCategoryNotFound, nil, err.Error())
		return
	}

	response := utils.CursorPaginationResponse{
		Items:      categories,
		NextCursor: nextCursor,
		HasMore:    nextCursor != nil,
	}
	HandleSuccess(ctx, response, "Success")
}

// UpdateCategory handles PUT /v1/cms/categories/:id
func (h *BlogHandler) UpdateCategory(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to parse ID param for UpdateCategory")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	var req domain.UpdateBlogCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("Failed to bind JSON for UpdateCategory")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}
	category, err := h.categorySvc.UpdateCategory(id, req)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to update category")
		HandleError(ctx, domain.ErrorCodeBlogCategoryNotFound, nil, err.Error())
		return
	}
	HandleSuccess(ctx, category, "Category updated")
}

// DeleteCategory handles DELETE /v1/cms/categories/:id
func (h *BlogHandler) DeleteCategory(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to parse ID param for DeleteCategory")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	if err := h.categorySvc.DeleteCategory(id); err != nil {
		logger.Log.WithError(err).Error("Failed to delete category")
		HandleError(ctx, domain.ErrorCodeBlogCategoryNotFound, nil, err.Error())
		return
	}
	HandleSuccess(ctx, nil, "Category deleted")
}

// CreatePost handles POST /v1/cms/posts
func (h *BlogHandler) CreatePost(ctx *gin.Context) {
	authorID, err := getAuthorID(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to get author ID from context")
		HandleError(ctx, domain.ErrorCodeUnAuthorization, nil, err.Error())
		return
	}
	var req domain.CreateBlogPostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("Failed to bind JSON for CreatePost")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}
	post, err := h.postSvc.CreatePost(authorID, req)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to create post")
		HandleError(ctx, domain.ErrorCodeBlogPostNotFound, nil, err.Error())
		return
	}
	if req.Status == domain.PostStatusPublished {
		go func() {
			h.searchEngineSvc.IndexDocument(string(domain.SearchEngineIndexNamePosts), post)
		}()
	}
	HandleSuccess(ctx, post, "Post created")
}

// GetPost handles GET /v1/cms/posts/:id
func (h *BlogHandler) GetPost(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to parse ID param for GetPost")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	post, err := h.postSvc.GetPost(id)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to get post")
		HandleError(ctx, domain.ErrorCodeBlogPostNotFound, nil, err.Error())
		return
	}
	HandleSuccess(ctx, post, "Success")
}

// ListPosts handles GET /v1/cms/posts
func (h *BlogHandler) ListPosts(ctx *gin.Context) {
	var filter domain.BlogPostFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		logger.Log.WithError(err).Error("Failed to bind query parameters for ListPosts")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	posts, total, err := h.postSvc.ListPosts(filter)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to list posts")
		HandleError(ctx, domain.ErrorCodeBlogPostNotFound, nil, err.Error())
		return
	}
	HandleSuccess(ctx, domain.BlogPostListResponse{Total: total, Posts: posts}, "Success")
}

// ListPostsCursor handles GET /v1/cms/posts/cursor with cursor-based pagination
func (h *BlogHandler) ListPostsCursor(ctx *gin.Context) {
	cursor := ctx.Query("cursor")
	limit := 10
	if l := ctx.Query("limit"); l != "" {
		if parsed, err := parseLimit(l); err == nil {
			limit = parsed
		}
	}

	var filter domain.BlogPostFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		logger.Log.WithError(err).Error("Failed to bind query parameters for ListPostsCursor")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}

	posts, nextCursor, total, err := h.postSvc.ListPostsCursor(filter, cursor, limit)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to list posts with cursor")
		HandleError(ctx, domain.ErrorCodeBlogPostNotFound, nil, err.Error())
		return
	}

	response := utils.CursorPaginationResponse{
		Items:      posts,
		NextCursor: nextCursor,
		HasMore:    nextCursor != nil,
		Total:      &total,
	}
	HandleSuccess(ctx, response, "Success")
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
		logger.Log.Errorf("Failed to bind JSON for UpdatePost: %v", err)
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}
	err = h.postSvc.UpdatePost(id, req)
	if err != nil {
		logger.Log.Errorf("Failed to update post: %v", err)
		HandleError(ctx, domain.ErrorCodeBlogPostNotFound, nil, err.Error())
		return
	}
	go func() {
		post, err := h.postSvc.GetPost(id)
		if err != nil {
			logger.Log.Errorf("Failed to get post: %v", err)
			return
		}
		if post.Status != domain.PostStatusPublished {
			err = h.searchEngineSvc.DeleteDocument(string(domain.SearchEngineIndexNamePosts), id)
		} else {
			err = h.searchEngineSvc.IndexDocument(string(domain.SearchEngineIndexNamePosts), post)
		}
		if err != nil {
			logger.Log.Errorf("Failed to update post in search engine: %v\n", err)
		}
	}()
	HandleSuccess(ctx, nil, "Post updated")
}

// DeletePost handles DELETE /v1/cms/posts/:id
func (h *BlogHandler) DeletePost(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	if err := h.postSvc.DeletePost(id); err != nil {
		logger.Log.Errorf("Failed to delete post error %s", err.Error())
		HandleError(ctx, domain.ErrorCodeBlogPostNotFound, validate.FormatValidationError(err), "Invalid request payload")
		return
	}
	go func() {
		err := h.searchEngineSvc.DeleteDocument(string(domain.SearchEngineIndexNamePosts), id)
		if err != nil {
			logger.Log.Errorf("Failed to delete post from search engine: %v\n", err)
		}
	}()
	HandleSuccess(ctx, nil, "Post deleted")
}

// PublishPost handles POST /v1/cms/posts/:id/publish
func (h *BlogHandler) PublishPost(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		logger.Log.Errorf("Failed to parse ID param for PublishPost: %v", err)
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}
	err = h.postSvc.PublishPost(id)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to publish post")
		HandleError(ctx, domain.ErrorCodeBlogPostNotFound, nil, err.Error())
		return
	}
	HandleSuccess(ctx, nil, "Post published")
}

// GetPostBySlug handles GET /v1/blog/posts/:slug
func (h *BlogHandler) GetPostBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	post, err := h.postSvc.GetPostBySlug(slug)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to get post by slug")
		HandleError(ctx, domain.ErrorCodeBlogPostNotFound, nil, err.Error())
		return
	}
	HandleSuccess(ctx, post, "Success")
}

// ListPublishedPosts handles GET /v1/blog/posts (public, published only)
func (h *BlogHandler) ListPublishedPosts(ctx *gin.Context) {
	var filter domain.BlogPostFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		logger.Log.WithError(err).Error("Failed to bind query parameters for ListPublishedPosts")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	filter.Status = string(domain.PostStatusPublished)
	posts, total, err := h.postSvc.ListPosts(filter)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to list published posts")
		HandleError(ctx, domain.ErrorCodeBlogPostNotFound, nil, err.Error())
		return
	}
	HandleSuccess(ctx, domain.BlogPostListResponse{Total: total, Posts: posts}, "Success")
}

func (h *BlogHandler) SearchPosts(ctx *gin.Context) {
	query := ctx.Query("query")
	limit := 10
	if l := ctx.Query("limit"); l != "" {
		if parsed, err := parseLimit(l); err == nil {
			limit = parsed
		}
	}

	results, err := h.searchEngineSvc.Search(string(domain.SearchEngineIndexNamePosts), query, limit)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeBlogPostNotFound, nil, err.Error())
		return
	}
	HandleSuccess(ctx, results, "Success")
}

func (h *BlogHandler) SearchBlogPostsPublic(ctx *gin.Context) {
	query := ctx.Query("query")
	limit := 10
	if l := ctx.Query("limit"); l != "" {
		if parsed, err := parseLimit(l); err == nil {
			limit = parsed
		}
	}

	results, err := h.searchEngineSvc.Search(string(domain.SearchEngineIndexNamePosts), query, limit)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to search blog posts")
		HandleError(ctx, domain.ErrorCodeBlogPostNotFound, nil, err.Error())
		return
	}
	response := make([]*domain.BlogUserSearchResult, 0, len(results))
	for _, result := range results {
		post := &domain.BlogUserSearchResult{}
		if err := utils.MapToStruct(result, post); err != nil {
			logger.Log.WithError(err).Error("Failed to map search result to BlogUserSearchResult")
			HandleError(ctx, domain.ErrorCodeInternalServerError, nil, fmt.Sprintf("Failed to map search result to BlogUserSearchResult: %v", err))
			return
		}
		response = append(response, post)
	}
	HandleSuccess(ctx, response, "Success")
}

func parseIDParam(ctx *gin.Context) (string, error) {
	return ctx.Param("id"), nil
}

func parseLimit(limitStr string) (int, error) {
	var limit int
	_, err := fmt.Sscanf(limitStr, "%d", &limit)
	return limit, err
}

func getAuthorID(ctx *gin.Context) (string, error) {
	idStr, err := customhttp.GetUserID(ctx)
	if err != nil {
		return "", err
	}
	return idStr, nil
}
