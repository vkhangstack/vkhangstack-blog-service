package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
)

func (h *UserHandler) LoginUser(ctx *gin.Context) {
	var user domain.Customer
	if err := ctx.ShouldBindJSON(&user); err != nil {
		HandleError(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	response, err := h.svc.LoginUser(user.Email, "")
	if err != nil {
		HandleError(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	data := map[string]interface{}{
		"id":            response.ID,
		"email":         response.Email,
		"access_token":  response.AccessToken,
		"refresh_token": response.RefreshToken,
	}

	HandleSuccess(ctx, data, "Login success!")
}
