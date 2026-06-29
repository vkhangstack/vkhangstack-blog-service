package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/validate"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/services"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

type TaskHandler struct {
	svc *services.TaskService
}

func NewTaskHandler(svc *services.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

// CreateTask handles POST /v1/cms/tasks
func (h *TaskHandler) CreateTask(ctx *gin.Context) {
	var req domain.CreateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}
	task, err := h.svc.CreateTask(req)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, task, "Task created")
}

// GetTask handles GET /v1/cms/tasks/:id
func (h *TaskHandler) GetTask(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	task, err := h.svc.GetTask(id)
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

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	tasks, total, err := h.svc.ListTasks(filter)
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

	var filter domain.TaskFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid query parameters")
		return
	}

	tasks, nextCursor, total, err := h.svc.ListTasksCursor(filter, cursor, limit)
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
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	var req domain.UpdateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}
	task, err := h.svc.UpdateTask(id, req)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, task, "Task updated")
}

// DeleteTask handles DELETE /v1/cms/tasks/:id
func (h *TaskHandler) DeleteTask(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	if err := h.svc.DeleteTask(id); err != nil {
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, nil, "Task deleted")
}

// GetTaskStatistics handles GET /v1/cms/tasks/statistics
func (h *TaskHandler) GetTaskStatistics(ctx *gin.Context) {
	stats, err := h.svc.GetTaskStatistics()
	if err != nil {
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, stats, "Success")
}

func parseTaskLimit(limitStr string) (int, error) {
	var limit int
	_, err := fmt.Sscanf(limitStr, "%d", &limit)
	return limit, err
}
