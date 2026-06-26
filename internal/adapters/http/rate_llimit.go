package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/services"
)

func RateLimitMiddleware(rateLimiter *services.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the client's IP address
		clientIP := c.ClientIP()

		// Check if the client is allowed to make a request
		if !rateLimiter.Allow(clientIP) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "Too many requests. Please try again later.",
				"error":   domain.ErrorCodeTooManyRequests,
			})
			return
		}

		c.Next()
	}
}
