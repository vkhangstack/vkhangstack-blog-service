package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/validate"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
	"github.com/vkhangstack/hexagonal-architecture/internal/logger"
)

// RoleHandler implements the Part B role permission management endpoints.
type RoleHandler struct {
	svc ports.RoleService
}

func NewRoleHandler(svc ports.RoleService) *RoleHandler {
	return &RoleHandler{svc: svc}
}

// ListRoles handles GET /v1/account/roles
func (h *RoleHandler) ListRoles(ctx *gin.Context) {
	HandleSuccess(ctx, h.svc.ListRoles(), "success")
}

// GetRolePermissions handles GET /v1/account/roles/:role/permissions
func (h *RoleHandler) GetRolePermissions(ctx *gin.Context) {
	role := getRole(ctx)
	perms, err := h.svc.GetRolePermissions(role)
	if err != nil {
		logger.Log.WithError(err).Error("GetRolePermissions failed")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	HandleSuccess(ctx, perms, "success")
}

// UpdateRolePermissions handles PUT /v1/account/roles/:role/permissions
func (h *RoleHandler) UpdateRolePermissions(ctx *gin.Context) {
	role := getRole(ctx)
	var req domain.UpdateRolePermissionsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}
	if err := h.svc.UpdateRolePermissions(role, req); err != nil {
		logger.Log.WithError(err).Error("UpdateRolePermissions failed")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	HandleSuccess(ctx, nil, "Role permissions updated")
}
