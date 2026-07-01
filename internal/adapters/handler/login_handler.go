package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/validate"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/services"
	"github.com/vkhangstack/hexagonal-architecture/internal/logger"
)

type LoginHandler struct {
	sva services.AccountService
}

func NewLoginHandler(sva services.AccountService) *LoginHandler {
	return &LoginHandler{
		sva: sva,
	}
}

func (h *LoginHandler) LoginAccount(ctx *gin.Context) {
	var user domain.LoginRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		logger.Log.WithError(err).Error("LoginAccount: Invalid request payload")
		HandleError(ctx, http.StatusBadRequest, validate.FormatValidationError(err), "Invalid request payload")
		return
	}

	response, err := h.sva.LoginAccount(user.Username, user.Password)
	if err != nil {
		logger.Log.WithError(err).Error("LoginAccount: Failed to login account")
		HandleError(ctx, domain.ErrorCodeInvalidCredentials, nil, "Username or password is incorrect")
		return
	}
	profile, err := h.sva.ProfileAccount(response.ID)
	if err != nil {
		logger.Log.WithError(err).Error("LoginAccount: Failed to get user profile")
		HandleError(ctx, domain.ErrorCodeInvalidCredentials, nil, "Username or password is incorrect")
		return
	}

	data := domain.LoginResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		User:         &domain.Profile{ID: profile.ID, Username: profile.Username, FullName: profile.FullName},
	}

	HandleSuccess(ctx, data, "Login success!")
}
