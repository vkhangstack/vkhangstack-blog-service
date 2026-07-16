package handler

import (
	"github.com/gin-gonic/gin"
	customhttp "github.com/vkhangstack/hexagonal-architecture/internal/adapters/http"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/validate"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
	"github.com/vkhangstack/hexagonal-architecture/internal/logger"
)

// AccountHandler implements the Part A account CRUD endpoints.
type AccountHandler struct {
	svc ports.AccountService
}

func NewAccountHandler(svc ports.AccountService) *AccountHandler {
	return &AccountHandler{svc: svc}
}

// callerIsRoot reports whether the authenticated caller holds the root role,
// per their JWT claim, so root-only actions (e.g. granting the root role to
// another account) can be permitted.
func callerIsRoot(ctx *gin.Context) bool {
	role, err := customhttp.GetUserRole(ctx)
	return err == nil && role == domain.RoleRoot
}

// ListAccounts handles GET /v1/account
func (h *AccountHandler) ListAccounts(ctx *gin.Context) {
	accounts, err := h.svc.ListAccounts(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("ListAccounts failed")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, accounts, "success")
}

// GetAccount handles GET /v1/account/:id
func (h *AccountHandler) GetAccount(ctx *gin.Context) {
	id := getParamID(ctx)
	account, err := h.svc.GetAccount(ctx, id)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeAccountNotFound, nil, "account not found")
		return
	}
	HandleSuccess(ctx, account, "success")
}

// GetMe handles GET /v1/account/me
func (h *AccountHandler) GetMe(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeAccountNotFound, nil, "profile me not found")
		return
	}
	account, err := h.svc.GetAccount(ctx, userID)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeAccountNotFound, nil, "profile me not found")
		return
	}
	HandleSuccess(ctx, account, "Get profile me success")
}

// CreateAccount handles POST /v1/account
func (h *AccountHandler) CreateAccount(ctx *gin.Context) {
	var req domain.CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}
	account, err := h.svc.CreateAccount(ctx, req, callerIsRoot(ctx))
	if err != nil {
		logger.Log.WithError(err).Error("CreateAccount failed")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	HandleSuccess(ctx, account, "Account created")
}

// UpdateAccount handles PUT /v1/account/:id
func (h *AccountHandler) UpdateAccount(ctx *gin.Context) {
	id := getParamID(ctx)
	var req domain.UpdateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.Debugf("Invalid UpdateAcount %s", err.Error())
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}
	account, err := h.svc.UpdateAccount(ctx, id, req, callerIsRoot(ctx))
	if err != nil {
		logger.Log.WithError(err).Error("UpdateAccount failed")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	HandleSuccess(ctx, account, "Account updated")
}

// DeleteAccount handles DELETE /v1/account/:id
func (h *AccountHandler) DeleteAccount(ctx *gin.Context) {
	id := getParamID(ctx)
	if err := h.svc.DeleteAccount(ctx, id); err != nil {
		logger.Log.WithError(err).Error("DeleteAccount failed")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, err.Error())
		return
	}
	HandleSuccess(ctx, nil, "Account deleted")
}

// InviteAccount handles POST /v1/account/invite
func (h *AccountHandler) InviteAccount(ctx *gin.Context) {
	var req domain.InviteAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.Debugf("Invalid request payload %s", err.Error())
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}
	account, err := h.svc.InviteAccount(ctx, req, callerIsRoot(ctx))
	if err != nil {
		logger.Log.WithError(err).Error("InviteAccount failed")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	HandleSuccess(ctx, account, "Invite sent")
}
