package http

import (
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/port"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// newJwtMiddleware is a middleware used to authenticate user by JWT.
//
// Note: Key value sets the key where the token will be stored in the context.
func newJwtMiddleware(generator port.TokenGenerator, key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			handleError(c, domain.ErrMalformedToken)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")
		token, err := generator.ParseToken(tokenString)
		if err != nil {
			handleError(c, domain.ErrInvalidToken)
			c.Abort()
			return
		}

		c.Set(key, token)
		c.Next()
	}
}

// zapLogger is a middleware used to log incoming requests.
func zapLogger() gin.HandlerFunc {
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
