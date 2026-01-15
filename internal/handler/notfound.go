package handler

import (
	"net/http"
	"no-as-a-service/internal/response"

	"github.com/gin-gonic/gin"
)

func NotFound(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"Status": response.StatusResponse{
			Status:  http.StatusNotFound,
			Message: "The requested resource was not found",
			Code:    "RT_404_NOT_FOUND",
		},
	})
}
