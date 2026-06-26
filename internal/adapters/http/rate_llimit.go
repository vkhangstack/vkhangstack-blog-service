package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/services"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

func RateLimitMiddleware(rateLimiter *services.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("x-api-key")

		// Get the client's IP address
		if key == "" {
			key = c.ClientIP()
		}

		// Check if the client is allowed to make a request
		ok, retryAfter := rateLimiter.Allow(key)
		if !ok {
			c.Header("Retry-After", utils.Uint64ToString(uint64(retryAfter)))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "Too many requests. Please try again later.",
				"error":   domain.ErrorCodeTooManyRequests,
			})
			return
		}

		c.Next()
	}
}
