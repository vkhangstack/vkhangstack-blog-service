package handler

import (
	"github.com/gin-gonic/gin"
	httpAdapter "github.com/vkhangstack/hexagonal-architecture/internal/adapters/http"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
	"github.com/vkhangstack/hexagonal-architecture/internal/logger"
)

// MenuHandler serves the navigation menu endpoint.
type MenuHandler struct {
	menuSvc ports.MenuService
}

func NewMenuHandler(menuSvc ports.MenuService) *MenuHandler {
	return &MenuHandler{menuSvc: menuSvc}
}

// GetMenu returns the role-filtered navigation menu for the authenticated account.
// GET /v1/account/menu
func (h *MenuHandler) GetMenu(c *gin.Context) {
	role, err := httpAdapter.GetUserRole(c)
	if err != nil {
		logger.Log.WithError(err).Error("GetMenu: failed to get user role")
		HandleError(c, domain.ErrorCodeForbidden, nil, "missing role")
		return
	}

	menu, err := h.menuSvc.GetMenu(c.Request.Context(), role)
	if err != nil {
		logger.Log.WithError(err).Error("GetMenu: failed to build menu")
		HandleError(c, domain.ErrorCodeInternalServerError, nil, "failed to build menu")
		return
	}

	HandleSuccess(c, menu, "success")
}
