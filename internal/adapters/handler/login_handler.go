package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/services"
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
		HandleError(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	response, err := h.sva.LoginAccount(user.Username, user.Password)
	if err != nil {
		HandleError(ctx, http.StatusBadRequest, nil, "Username or password is incorrect")
		return
	}
	profile, err := h.sva.ProfileAccount(response.ID)
	if err != nil {
		HandleError(ctx, http.StatusBadRequest, nil, "Username or password is incorrect")
		return
	}

	data := domain.LoginResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		User:         &domain.Profile{ID: profile.ID, Username: profile.Username, FullName: profile.FullName},
	}

	HandleSuccess(ctx, data, "Login success!")
}
