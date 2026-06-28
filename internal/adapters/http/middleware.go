package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vkhangstack/hexagonal-architecture/internal/config"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/services"
	"github.com/vkhangstack/hexagonal-architecture/internal/logger"
)

// Middleware to set Content-Type as application/json
func JSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}

// Middleware to handle CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Next()
	}
}

// TracingMiddleware adds distributed tracing capabilities to track requests
func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generate or extract trace ID from request header
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.Must(uuid.NewV7()).String()
		}

		// Set trace ID in context for downstream use
		c.Set("trace_id", traceID)
		// Add trace ID to response headers
		c.Writer.Header().Set("X-Trace-ID", traceID)

		logger.Log.WithFields(map[string]interface{}{
			"trace_id": traceID,
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"client":   c.ClientIP(),
		}).Debug("request started")

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		entry := logger.Log.WithFields(map[string]interface{}{
			"trace_id": traceID,
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"status":   statusCode,
			"duration": duration.String(),
		})
		if statusCode >= 500 {
			entry.Error("request completed")
		} else if statusCode >= 400 {
			entry.Warn("request completed")
		} else {
			entry.Info("request completed")
		}
	}
}

func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		loadConfig := config.LoadConfig().App
		userID, err := services.ValidateToken(c.Request.Header.Get("Authorization"), loadConfig.JWTSecret)
		if err != nil {
			if errors.Is(err, errors.New("token is expired")) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":   domain.ErrorCodeUnAuthorization,
					"message": "Unauthorized: Token is expired",
					"data":    nil,
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusForbidden, gin.H{
				"error":   domain.ErrorCodeForbidden,
				"message": "Unauthorized: Invalid token",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

func GetUserID(c *gin.Context) (string, error) {
	userID, ok := c.Get("user_id")
	if !ok {
		return "", errors.New("user_id not found in context")
	}
	return userID.(string), nil
}
