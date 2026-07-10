package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
	"github.com/vkhangstack/hexagonal-architecture/internal/logger"
)

// MenuAdminHandler handles CRUD operations for CMS-managed menu entries.
type MenuAdminHandler struct {
	svc ports.MenuAdminService
}

func NewMenuAdminHandler(svc ports.MenuAdminService) *MenuAdminHandler {
	return &MenuAdminHandler{svc: svc}
}

// CreateMenu godoc
// POST /v1/cms/menus
func (h *MenuAdminHandler) CreateMenu(c *gin.Context) {
	var req domain.CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	entry, err := h.svc.CreateMenu(c.Request.Context(), req)
	if err != nil {
		logger.Log.WithError(err).Error("CreateMenu failed")
		HandleError(c, domain.ErrorCodeInternalServerError, nil, "failed to create menu")
		return
	}
	HandleSuccess(c, entry, "Menu created")
}

// ListMenus godoc
// GET /v1/cms/menus
func (h *MenuAdminHandler) ListMenus(c *gin.Context) {
	entries, err := h.svc.ListMenus(c.Request.Context())
	if err != nil {
		logger.Log.WithError(err).Error("ListMenus failed")
		HandleError(c, domain.ErrorCodeInternalServerError, nil, "failed to list menus")
		return
	}
	HandleSuccess(c, entries, "success")
}

// GetMenu godoc
// GET /v1/cms/menus/:id
func (h *MenuAdminHandler) GetMenuByID(c *gin.Context) {
	id := c.Param("id")
	entry, err := h.svc.GetMenu(c.Request.Context(), id)
	if err != nil {
		HandleError(c, http.StatusNotFound, nil, "menu not found")
		return
	}
	HandleSuccess(c, entry, "success")
}

// UpdateMenu godoc
// PUT /v1/cms/menus/:id
func (h *MenuAdminHandler) UpdateMenu(c *gin.Context) {
	id := c.Param("id")
	var req domain.UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	if err := h.svc.UpdateMenu(c.Request.Context(), id, req); err != nil {
		logger.Log.WithError(err).Error("UpdateMenu failed")
		HandleError(c, domain.ErrorCodeInternalServerError, nil, "failed to update menu")
		return
	}
	HandleSuccess(c, nil, "Menu updated")
}

// DeleteMenu godoc
// DELETE /v1/cms/menus/:id
func (h *MenuAdminHandler) DeleteMenu(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteMenu(c.Request.Context(), id); err != nil {
		logger.Log.WithError(err).Error("DeleteMenu failed")
		HandleError(c, domain.ErrorCodeInternalServerError, nil, "failed to delete menu")
		return
	}
	HandleSuccess(c, nil, "Menu deleted")
}
