package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/http"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/services"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

type UserHandler struct {
	svc services.CustomerService
	fvc services.FirebaseService
}

func NewUserHandler(CustomerService services.CustomerService, FirebaseService services.FirebaseService) *UserHandler {
	return &UserHandler{
		svc: CustomerService,
		fvc: FirebaseService,
	}
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	var user domain.Customer
	if err := ctx.ShouldBindJSON(&user); err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}

	_, err := h.svc.CreateUser(user.Email, "")
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	HandleSuccess(ctx, nil, "User created successfully")
}

func (h *UserHandler) ReadUser(ctx *gin.Context) {
	id := ctx.Param("id")

	userId, _ := utils.ParseUint64(id)
	user, err := h.svc.ReadUser(userId)

	if err != nil {
		HandleError(ctx, 4000, nil, err.Error())
		return
	}

	data := map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
	}
	HandleSuccess(ctx, data, "Success")
}

func (h *UserHandler) ReadUsers(ctx *gin.Context) {

	// users, err := h.svc.ReadUsers()
	// if err != nil {
	// 	HandleError(ctx, 4000, fmt.Errorf("user not found"))
	// 	return
	// }
	data := h.fvc.GetUser(context.Background(), "uid")
	HandleSuccess(ctx, data, "Success")
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	// Get user ID from token
	userID, err := http.GetUserID(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeUnAuthorization, nil, err.Error())
		return
	}

	// Update user
	var user domain.Customer
	if err := ctx.ShouldBindJSON(&user); err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}

	err = h.svc.UpdateUser(userID, user.Email, "")
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	HandleSuccess(ctx, nil, "Success")
}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	userID, err := http.GetUserID(ctx)
	if err != nil {
		HandleError(ctx, domain.ErrorCodeUnAuthorization, nil, err.Error())
		return
	}
	userId, _ := utils.ParseUint64(userID)

	err = h.svc.DeleteUser(userId)
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	HandleSuccess(ctx, nil, "Success")
}
