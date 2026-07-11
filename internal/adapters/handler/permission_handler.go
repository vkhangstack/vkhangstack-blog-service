package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	httpAdapter "github.com/vkhangstack/hexagonal-architecture/internal/adapters/http"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
	"github.com/vkhangstack/hexagonal-architecture/internal/logger"
)

// PermissionHandler manages per-user Casbin permission grants.
type PermissionHandler struct {
	authzSvc   ports.AuthorizationService
	accountSvc ports.AccountService
}

func NewPermissionHandler(authzSvc ports.AuthorizationService, accountSvc ports.AccountService) *PermissionHandler {
	return &PermissionHandler{authzSvc: authzSvc, accountSvc: accountSvc}
}

// isRootUser looks up the target account's actual role (source of truth, independent of
// whether the user has logged in yet and so has a synced Casbin grouping row).
func (h *PermissionHandler) isRootUser(ctx context.Context, userID string) bool {
	account, err := h.accountSvc.ProfileAccount(ctx, userID)
	if err != nil || account == nil {
		return false
	}
	return account.Role == domain.RoleRoot
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
	if h.isRootUser(c.Request.Context(), req.UserID) {
		HandleError(c, domain.ErrorCodeForbidden, nil, "root account permissions cannot be modified")
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
	if h.isRootUser(c.Request.Context(), req.UserID) {
		HandleError(c, domain.ErrorCodeForbidden, nil, "root account permissions cannot be modified")
		return
	}
	if err := h.authzSvc.RevokePermission(req.UserID, req.Resource, req.Action); err != nil {
		logger.Log.WithError(err).Error("RevokePermission failed")
		HandleError(c, domain.ErrorCodeInternalServerError, nil, "failed to revoke permission")
		return
	}
	HandleSuccess(c, nil, "Permission revoked")
}

type roleAssignmentRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Role   string `json:"role"    binding:"required"`
}

// AssignRole grants a specific user membership in an RBAC role, so they immediately
// inherit that role's whole permission set (in addition to their existing role).
// POST /v1/cms/permissions/assign-role
func (h *PermissionHandler) AssignRole(c *gin.Context) {
	var req roleAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	if !domain.IsKnownRole(req.Role) {
		HandleError(c, domain.ErrorCodePayloadBadRequest, nil, "unknown role")
		return
	}
	if req.Role == domain.RoleRoot && !callerIsRoot(c) {
		HandleError(c, domain.ErrorCodeForbidden, nil, "root role cannot be assigned")
		return
	}
	if h.isRootUser(c.Request.Context(), req.UserID) {
		HandleError(c, domain.ErrorCodeForbidden, nil, "root account roles cannot be modified")
		return
	}
	if err := h.authzSvc.AssignRole(req.UserID, req.Role); err != nil {
		logger.Log.WithError(err).Error("AssignRole failed")
		HandleError(c, domain.ErrorCodeInternalServerError, nil, "failed to assign role")
		return
	}
	HandleSuccess(c, nil, "Role assigned")
}

// RemoveRole revokes a specific user's membership in an RBAC role.
// POST /v1/cms/permissions/revoke-role
func (h *PermissionHandler) RemoveRole(c *gin.Context) {
	var req roleAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	if h.isRootUser(c.Request.Context(), req.UserID) {
		HandleError(c, domain.ErrorCodeForbidden, nil, "root account roles cannot be modified")
		return
	}
	if err := h.authzSvc.UnassignRole(req.UserID, req.Role); err != nil {
		logger.Log.WithError(err).Error("UnassignRole failed")
		HandleError(c, domain.ErrorCodeInternalServerError, nil, "failed to remove role")
		return
	}
	HandleSuccess(c, nil, "Role removed")
}

// GetUserRoles returns the RBAC roles a specific user is currently a member of.
// GET /v1/cms/permissions/:user_id/roles
func (h *PermissionHandler) GetUserRoles(c *gin.Context) {
	userID := c.Param("user_id")
	HandleSuccess(c, h.authzSvc.GetUserRoles(userID), "success")
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
