package http

import (
	"shop-api-go/internal/core/domain"
	"shop-api-go/internal/core/port"
	"strings"

	"github.com/gin-gonic/gin"
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
