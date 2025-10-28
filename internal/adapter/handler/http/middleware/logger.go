package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ZapLogger is a middleware used to log incoming requests.
func ZapLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)

		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
			zap.Int("status", c.Writer.Status()),
			zap.String("user-agent", c.Request.UserAgent()),
		}
		switch {
		case c.Writer.Status() >= 500:
			zap.L().Error("Server error", fields...)
		case c.Writer.Status() >= 400:
			zap.L().Warn("Client error", fields...)
		default:
			zap.L().Debug("Incoming request", fields...)
		}
	}
}
