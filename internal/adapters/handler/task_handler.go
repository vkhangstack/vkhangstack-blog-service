package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/http"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/validate"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/services"
	"github.com/vkhangstack/hexagonal-architecture/internal/logger"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

type TaskHandler struct {
	svc             *services.TaskService
	searchEngineSvc *services.SearchEngineService
}

func NewTaskHandler(svc *services.TaskService, searchEngineSvc *services.SearchEngineService) *TaskHandler {
	return &TaskHandler{
		svc:             svc,
		searchEngineSvc: searchEngineSvc,
	}
}

// CreateTask handles POST /v1/cms/tasks
func (h *TaskHandler) CreateTask(ctx *gin.Context) {
	var req domain.CreateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}
	userID, err := http.GetUserID(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("CreateTask: Failed to get user ID from context")
		HandleError(ctx, domain.ErrorCodeForbidden, nil, "Unauthorized")
		return
	}
	task, err := h.svc.CreateTask(ctx, userID, req)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	go func() {
		if err := h.searchEngineSvc.IndexDocument(string(domain.SearchEngineIndexNameTasks), task); err != nil {
			logger.Log.WithError(err).Error("CreateTask: Failed to index task in search engine")
		}
	}()
	HandleSuccess(ctx, task, "Task created")
}

// GetTask handles GET /v1/cms/tasks/:id
func (h *TaskHandler) GetTask(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	task, err := h.svc.GetTask(ctx, id)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, task, "Success")
}

// ListTasks handles GET /v1/cms/tasks
func (h *TaskHandler) ListTasks(ctx *gin.Context) {
	var filter domain.TaskFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid query parameters")
		return
	}
	userID, err := http.GetUserID(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("ListTasks: Failed to get user ID from context")
		HandleError(ctx, domain.ErrorCodeForbidden, nil, "Unauthorized")
		return
	}
	filter.CreatedBy = userID

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	tasks, total, err := h.svc.ListTasks(ctx, filter)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}

	response := domain.TaskListResponse{
		Total: total,
		Tasks: tasks,
	}
	HandleSuccess(ctx, response, "Success")
}

// ListTasksCursor handles GET /v1/cms/tasks/cursor with cursor-based pagination
func (h *TaskHandler) ListTasksCursor(ctx *gin.Context) {
	cursor := ctx.Query("cursor")
	limit := 10
	if l := ctx.Query("limit"); l != "" {
		if parsed, err := parseTaskLimit(l); err == nil {
			limit = parsed
		}
	}
	userID, err := http.GetUserID(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("ListTasksCursor: Failed to get user ID from context")
		HandleError(ctx, domain.ErrorCodeForbidden, nil, "Unauthorized")
		return
	}

	var filter domain.TaskFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid query parameters")
		return
	}

	filter.CreatedBy = userID

	tasks, nextCursor, total, err := h.svc.ListTasksCursor(ctx, filter, cursor, limit)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}

	response := utils.CursorPaginationResponse{
		Items:      tasks,
		NextCursor: nextCursor,
		HasMore:    nextCursor != nil,
		Total:      &total,
	}
	HandleSuccess(ctx, response, "Success")
}

// UpdateTask handles PUT /v1/cms/tasks/:id
func (h *TaskHandler) UpdateTask(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("UpdateTask: Invalid task ID")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	var req domain.UpdateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("UpdateTask: Invalid request payload")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}
	task, err := h.svc.UpdateTask(ctx, id, req)
	if err != nil {
		logger.Log.WithError(err).Error("UpdateTask: Failed to update task")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	go func() {
		if err := h.searchEngineSvc.UpdateDocument(string(domain.SearchEngineIndexNameTasks), task); err != nil {
			logger.Log.WithError(err).Error("UpdateTask: Failed to update task in search engine")
		}
	}()
	HandleSuccess(ctx, task, "Task updated")
}

// DeleteTask handles DELETE /v1/cms/tasks/:id
func (h *TaskHandler) DeleteTask(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("DeleteTask: Invalid task ID")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	if err := h.svc.DeleteTask(ctx, id); err != nil {
		logger.Log.WithError(err).Error("DeleteTask: Failed to delete task")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	go func() {
		if err := h.searchEngineSvc.DeleteDocument(string(domain.SearchEngineIndexNameTasks), id); err != nil {
			logger.Log.WithError(err).Error("DeleteTask: Failed to delete task from search engine")
		}
	}()
	HandleSuccess(ctx, nil, "Task deleted")
}

// GetTaskStatistics handles GET /v1/cms/tasks/statistics
func (h *TaskHandler) GetTaskStatistics(ctx *gin.Context) {
	stats, err := h.svc.GetTaskStatistics(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("GetTaskStatistics: Failed to get task statistics")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, stats, "Success")
}

func (h *TaskHandler) SearchTasks(ctx *gin.Context) {
	query := ctx.Query("q")
	limitStr := ctx.DefaultQuery("limit", "10")

	if query == "" {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, "Query parameter 'q' is required")
		return
	}

	limit, err := parseTaskLimit(limitStr)
	if err != nil || limit <= 0 {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, "Invalid limit parameter")
		return
	}

	if limit > 10 {
		limit = 10
	}

	tasks, err := h.searchEngineSvc.Search(string(domain.SearchEngineIndexNameTasks), query, limit)
	if err != nil {
		logger.Log.WithError(err).Error("SearchTasks: Failed to search tasks")
		HandleSuccess(ctx, []interface{}{}, "Empty task")
		return
	}
	HandleSuccess(ctx, tasks, "Success")
}

func parseTaskLimit(limitStr string) (int, error) {
	var limit int
	_, err := fmt.Sscanf(limitStr, "%d", &limit)
	return limit, err
}
