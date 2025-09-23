package http

import (
	"errors"
	"fmt"
	"net/http"
	"shop-api-go/internal/core/domain"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

// ErrorResponse represents a json response when an error occurs.
type ErrorResponse struct {
	Code     string   `json:"code"`
	Messages []string `json:"messages"`
}

// parseValidationError transform a validator.FieldError into meaningful message.
func parseValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return fmt.Sprintf("%s is not a valid email", err.Field())
	case "password":
		return fmt.Sprintf("%s is not a valid password", err.Field())
	case "min_bytes":
		return fmt.Sprintf("%s length must be more than %s", err.Field(), err.Param())
	case "max_bytes":
		return fmt.Sprintf("%s length must be less than %s", err.Field(), err.Param())
	default:
		return err.Error()
	}
}

// errorHandler is globar error handler that responds with propper error response if an error occurred.
func errorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			var validationErrors validator.ValidationErrors
			var domainError *domain.Error
			switch {
			case errors.As(c.Errors.Last().Err, &validationErrors):
				messages := make([]string, 0, len(validationErrors))
				for _, e := range validationErrors {
					messages = append(messages, parseValidationError(e))
				}
				c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
					Code:     "INVALID_ENTITY",
					Messages: messages,
				})
			case errors.As(c.Errors.Last().Err, &domainError):
				c.JSON(domainError.StatusCode, ErrorResponse{
					Code:     domainError.Code,
					Messages: []string{domainError.Error()},
				})
			default:
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Code:     "INTERNAL_SERVER_ERROR",
					Messages: []string{"unknown error"},
				})
			}

		}
	}
}

// jwtMiddleware middleware is used to validate a JWT token and pass it down the chain.
func jwtMiddleware(key string, secret []byte, role domain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			_ = c.Error(domain.ErrInvalidToken)
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenStr, &domain.Token{}, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})
		if err != nil {
			_ = c.Error(fmt.Errorf("%w: %v", domain.ErrInvalidToken, err))
			c.Abort()
			return
		}

		parsedToken, ok := token.Claims.(*domain.Token)
		if !ok || !token.Valid || parsedToken.UserRole != role {
			_ = c.Error(fmt.Errorf("%w: %v", domain.ErrInvalidToken, "invalid token"))
			c.Abort()
			return
		}

		c.Set(key, parsedToken)
		c.Next()
	}
}
