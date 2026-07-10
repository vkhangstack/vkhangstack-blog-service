package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpAdapter "github.com/vkhangstack/hexagonal-architecture/internal/adapters/http"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
	"github.com/vkhangstack/hexagonal-architecture/internal/logger"
)

// PermissionHandler manages per-user Casbin permission grants.
type PermissionHandler struct {
	authzSvc ports.AuthorizationService
}

func NewPermissionHandler(authzSvc ports.AuthorizationService) *PermissionHandler {
	return &PermissionHandler{authzSvc: authzSvc}
}

type permissionRequest struct {
	UserID   string `json:"user_id"  binding:"required"`
	Resource string `json:"resource" binding:"required"`
	Action   string `json:"action"   binding:"required"`
}

// GrantPermission adds a direct policy for a user.
// POST /v1/cms/permissions/grant
func (h *PermissionHandler) GrantPermission(c *gin.Context) {
	var req permissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	if err := h.authzSvc.GrantPermission(req.UserID, req.Resource, req.Action); err != nil {
		logger.Log.WithError(err).Error("GrantPermission failed")
		HandleError(c, domain.ErrorCodeInternalServerError, nil, "failed to grant permission")
		return
	}
	HandleSuccess(c, nil, "Permission granted")
}

// RevokePermission removes a direct policy for a user.
// POST /v1/cms/permissions/revoke
func (h *PermissionHandler) RevokePermission(c *gin.Context) {
	var req permissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	if err := h.authzSvc.RevokePermission(req.UserID, req.Resource, req.Action); err != nil {
		logger.Log.WithError(err).Error("RevokePermission failed")
		HandleError(c, domain.ErrorCodeInternalServerError, nil, "failed to revoke permission")
		return
	}
	HandleSuccess(c, nil, "Permission revoked")
}

// GetUserPermissions returns direct policies for the requesting user.
// GET /v1/account/permissions
func (h *PermissionHandler) GetMyPermissions(c *gin.Context) {
	userID, err := httpAdapter.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	perms := h.authzSvc.GetUserPermissions(userID)
	HandleSuccess(c, perms, "success")
}
