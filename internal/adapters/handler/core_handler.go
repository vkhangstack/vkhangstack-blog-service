package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	customhttp "github.com/vkhangstack/hexagonal-architecture/internal/adapters/http"
)

func parseIDParam(ctx *gin.Context) (string, error) {
	return ctx.Param("id"), nil
}

func parseLimit(limitStr string) (int, error) {
	var limit int
	_, err := fmt.Sscanf(limitStr, "%d", &limit)
	return limit, err
}

func getAuthorID(ctx *gin.Context) (string, error) {
	idStr, err := customhttp.GetUserID(ctx)
	if err != nil {
		return "", err
	}
	return idStr, nil
}

func getUserID(ctx *gin.Context) (string, error) {
	idStr, err := customhttp.GetUserID(ctx)
	if err != nil {
		return "", err
	}
	return idStr, nil
}

func getCursor(ctx *gin.Context) string {
	return ctx.Query("cursor")
}

func getLimit(ctx *gin.Context) (int, error) {
	limitStr := ctx.Query("limit")
	if limitStr == "" {
		return 10, nil // Default limit
	}
	return parseLimit(limitStr)
}

func getRole(ctx *gin.Context) string {
	role := ctx.Param("role")
	return role
}

func getParamID(ctx *gin.Context) string {
	return ctx.Param("id")
}
