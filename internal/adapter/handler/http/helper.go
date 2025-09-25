package http

import (
	"errors"
	"fmt"
	"net/http"
	"shop-api-go/internal/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type errorResponse struct {
	code       string
	messages   []string
	statusCode int
}

var errMap = map[error]errorResponse{
	domain.ErrInternalServerError: {
		code:       "INTERNAL_SERVER_ERROR",
		messages:   []string{"Server cannot process the request."},
		statusCode: http.StatusInternalServerError,
	},
	domain.ErrEmailAlreadyInUse: {
		code:       "EMAIL_ALREADY_IN_USE",
		messages:   []string{"Email is already in use."},
		statusCode: http.StatusConflict,
	},
	domain.ErrUsernameAlreadyInUse: {
		code:       "USER_ALREADY_IN_USE",
		messages:   []string{"User is already in use."},
		statusCode: http.StatusConflict,
	},
	domain.ErrWrongCredentials: {
		code:       "WRONG_CREDENTIALS",
		messages:   []string{"Wrong credentials."},
		statusCode: http.StatusUnauthorized,
	},
	domain.ErrInvalidToken: {
		code:       "INVALID_TOKEN",
		messages:   []string{"Token is invalid."},
		statusCode: http.StatusUnauthorized,
	},
	domain.ErrInvalidTokenType: {
		code:       "INVALID_TOKEN_TYPE",
		messages:   []string{"Token type is invalid."},
		statusCode: http.StatusUnauthorized,
	},
	domain.ErrMalformedToken: {
		code:       "MALFORMED_TOKEN",
		messages:   []string{"Token is malformed."},
		statusCode: http.StatusUnauthorized,
	},
}

func handleError(c *gin.Context, err error) {
	res, ok := errMap[err]
	if !ok {
		res = errorResponse{
			code:       "INTERNAL_SERVER_ERROR",
			messages:   []string{"Server cannot process the request."},
			statusCode: http.StatusInternalServerError,
		}
	}

	zap.L().Error(
		"HTTP request error",
		zap.Int("status", res.statusCode),
		zap.Strings("messages", res.messages),
		zap.String("code", res.code),
		zap.Error(err),
	)

	c.JSON(res.statusCode, gin.H{
		"code":     res.code,
		"messages": res.messages,
	})
}

func handleBindingError(c *gin.Context, err error) {
	var validationsErrors validator.ValidationErrors
	messages := make([]string, 0, len(validationsErrors))

	if errors.As(err, &validationsErrors) {
		for _, e := range validationsErrors {
			switch e.Tag() {
			case "required":
				messages = append(messages, fmt.Sprintf("%s is required.", e.Field()))
			case "email":
				messages = append(messages, fmt.Sprintf("%s is not a valid email", e.Field()))
			case "password":
				messages = append(messages, fmt.Sprintf("%s is not a valid password", e.Field()))
			case "min_bytes":
				messages = append(messages, fmt.Sprintf("%s length must be more than %s", e.Field(), e.Param()))
			case "max_bytes":
				messages = append(messages, fmt.Sprintf("%s length must be less than %s", e.Field(), e.Param()))
			default:
				messages = append(messages, fmt.Sprintf("%s is not a valid type", e.Field()))
			}
		}
	} else {
		messages = append(messages, err.Error())
	}

	zap.L().Error("HTTP request error",
		zap.Int("status", http.StatusBadRequest),
		zap.Strings("messages", messages),
		zap.Error(err),
	)

	c.JSON(http.StatusBadRequest, gin.H{
		"code":    "INVALID_ENTITY",
		"message": messages,
	})
}
