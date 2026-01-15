package router

import (
	"no-as-a-service/internal/handler"

	"github.com/gin-gonic/gin"
)

func Setup(engine *gin.Engine) {
	engine.GET("/healthcheck", handler.Healthcheck)
	engine.GET("/reason", handler.Reason)

	engine.NoRoute(handler.NotFound)
}
