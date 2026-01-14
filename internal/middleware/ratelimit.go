package middleware

import (
	"net/http"
	"no-as-a-service/internal/response"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter tracks request counts per IP
type RateLimiter struct {
	requests map[string]*clientInfo
	mutex    sync.Mutex
	limit    int
	window   time.Duration
}

type clientInfo struct {
	count     int
	expiresAt time.Time
}

// NewRateLimiter creates a rate limiter with given limit per window
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string]*clientInfo),
		limit:    limit,
		window:   window,
	}
}

// Middleware returns a Gin middleware for rate limiting
func (rateLimiter *RateLimiter) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		now := time.Now()

		rateLimiter.mutex.Lock()
		info, exists := rateLimiter.requests[ip]

		// Reset if window expired or new client
		if !exists || now.After(info.expiresAt) {
			rateLimiter.requests[ip] = &clientInfo{count: 1, expiresAt: now.Add(rateLimiter.window)}
			rateLimiter.mutex.Unlock()
			ctx.Next()
			return
		}

		info.count++
		if info.count > rateLimiter.limit {
			rateLimiter.mutex.Unlock()
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"Status": response.StatusResponse{
					Status:  http.StatusTooManyRequests,
					Message: "Rate limit exceeded. Try again later.",
					Code:    "RL_429_TOO_MANY_REQUESTS",
				},
			})
			return
		}

		rateLimiter.mutex.Unlock()
		ctx.Next()
	}
}
