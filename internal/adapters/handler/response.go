package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Error int         `json:"error"`
	Data  interface{} `json:"data"`
	Msg   string      `json:"message"`
}

func HandleError(ctx *gin.Context, error int, data any, message string) {
	response := &Response{
		Error: error,
		Data:  data,
		Msg:   message,
	}
	ctx.JSON(http.StatusOK, response)
}

func HandleSuccess(ctx *gin.Context, data any, msg string) {
	response := &Response{
		Error: 0,
		Data:  data,
		Msg:   msg,
	}

	ctx.JSON(http.StatusOK, response)
}
