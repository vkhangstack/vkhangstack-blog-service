package handler

import (
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HealthCheck(ctx *gin.Context) {
	data := map[string]interface{}{
		"status": "ok",
		"time":   time.Unix(time.Now().Unix(), 0),
		"trace":  ctx.GetString("trace_id"),
	}
	HandleSuccess(ctx, data, "Service is healthy")
}
