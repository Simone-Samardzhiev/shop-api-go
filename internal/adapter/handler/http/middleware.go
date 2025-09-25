package http

import (
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/port"
	"strings"

	"github.com/gin-gonic/gin"
)

// jwtMiddleware is a middleware used to authenticate user by JWT.
//
// Note: The middleware support role specification and token type specification.
func jwtMiddleware(generator port.TokenGenerator, role domain.UserRole, tokenType domain.TokenType, key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			_ = c.Error(domain.ErrMalformedToken)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")
		token, err := generator.ParseToken(tokenString)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		if token.UserRole != role {
			_ = c.Error(domain.ErrMalformedToken)
			c.Abort()
			return
		}

		if token.TokenType != tokenType {
			_ = c.Error(domain.ErrInvalidToken)
			c.Abort()
			return
		}

		c.Set(key, token)
		c.Next()
	}
}
