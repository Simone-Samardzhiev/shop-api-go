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
		zapLevel := zap.InfoLevel
		msg := "Handling request"
		statusCode := c.Writer.Status()

		switch {
		case statusCode >= 500:
			zapLevel = zap.ErrorLevel
			msg = "Internal server error"
		case statusCode >= 400:
			zapLevel = zap.WarnLevel
			msg = "Client error"
		}

		zap.L().Log(zapLevel, msg,
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
			zap.Int("status", c.Writer.Status()),
			zap.String("user-agent", c.Request.UserAgent()),
		)
	}
}
