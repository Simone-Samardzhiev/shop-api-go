package middleware

import (
	"shop-api-go/internal/adapter/handler/http/response"
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/port"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTMiddleware is a middleware used to authenticate user by JWT.
//
// Note: Key value sets the key where the token will be stored in the context.
func JWTMiddleware(generator port.TokenGenerator, key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.HandleError(c, domain.ErrInvalidToken)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")
		token, err := generator.ParseToken(tokenString)
		if err != nil {
			response.HandleError(c, domain.ErrInvalidToken)
			c.Abort()
			return
		}

		c.Set(key, token)
		c.Next()
	}
}
