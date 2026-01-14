package main

import (
	"fmt"
	"no-as-a-service/internal/helper"
	"no-as-a-service/internal/middleware"
	"no-as-a-service/internal/router"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	var err error

	// Load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()

	// Read config from .env or use defaults
	port := helper.GetEnv("PORT", "8080")
	rateLimitMax := helper.GetEnvInt("RATE_LIMIT_MAX", 120)
	rateLimitWindow := helper.GetEnvInt("RATE_LIMIT_WINDOW_SECONDS", 60)

	// Initialize server engine
	engine := gin.New()

	// Set trusted proxies
	err = engine.SetTrustedProxies([]string{"::1"})
	if err != nil {
		panic(fmt.Errorf("Could not setup trusted proxies: %w", err))
	}

	// Add Recovery as middleware (recover from panics)
	engine.Use(gin.Recovery())

	// Add Rate Limiting middleware
	rateLimiter := middleware.NewRateLimiter(rateLimitMax, time.Duration(rateLimitWindow)*time.Second)
	engine.Use(rateLimiter.Middleware())

	// Setup API routes
	router.Setup(engine)

	// Start server
	err = engine.Run(":" + port)
	if err != nil {
		panic(fmt.Errorf("Could not start server: %w", err))
	}
}
