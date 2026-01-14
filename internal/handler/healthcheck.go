package handler

import (
	"net/http"
	"no-as-a-service/internal/response"

	"github.com/gin-gonic/gin"
)

func Healthcheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, response.DefaultResponse{
		Payload: nil,
		Status: response.StatusResponse{
			Status:  http.StatusOK,
			Message: "Service is healthy",
			Code:    "HC_200_OK",
		},
	})
}
