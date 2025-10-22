package response

import (
	"errors"
	"fmt"
	"net/http"
	"shop-api-go/internal/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// errorResponse represent a meaningful response when an error occurs.
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
		messages:   []string{"Username is already in use."},
		statusCode: http.StatusConflict,
	},
	domain.ErrWrongCredentials: {
		code:       "WRONG_CREDENTIALS",
		messages:   []string{"Wrong credentials."},
		statusCode: http.StatusUnauthorized,
	}, domain.ErrUserNotFound: {
		code:       "USER_NOT_FOUND",
		messages:   []string{"User not found."},
		statusCode: http.StatusNotFound,
	}, domain.ErrNoFieldsToUpdate: {
		code:       "NO_FIELDS_TO_UPDATE",
		messages:   []string{"No fields to update."},
		statusCode: http.StatusBadRequest,
	}, domain.ErrInvalidQuery: {
		code:       "INVALID_QUERY",
		messages:   []string{"Invalid query."},
		statusCode: http.StatusBadRequest,
	}, domain.ErrInvalidCursorFormat: {
		code:       "INVALID_CURSOR_FORMAT",
		messages:   []string{"Invalid cursor format."},
		statusCode: http.StatusBadRequest,
	},
	domain.ErrLimitNotSet: {
		code:       "LIMIT_NOT_SET",
		messages:   []string{"Please provide a limit when fetching a list of objects."},
		statusCode: http.StatusBadRequest,
	},
	domain.ErrInvalidToken: {
		code:       "INVALID_TOKEN",
		messages:   []string{"Token is invalid."},
		statusCode: http.StatusUnauthorized,
	},
	domain.ErrInvalidTokenType: {
		code:       "INVALID_TOKEN_TYPE",
		messages:   []string{"Token type is invalid."},
		statusCode: http.StatusForbidden,
	}, domain.ErrInvalidTokenRole: {
		code:       "INVALID_TOKEN_ROLE",
		messages:   []string{"Token role is invalid."},
		statusCode: http.StatusForbidden,
	},
	domain.ErrMalformedToken: {
		code:       "MALFORMED_TOKEN",
		messages:   []string{"Token is malformed."},
		statusCode: http.StatusUnauthorized,
	},
	domain.ErrTokenNotFound: {
		code:       "TOKEN_NOT_FOUND",
		messages:   []string{"Token not found."},
		statusCode: http.StatusNotFound,
	},
}

// HandleError parses the error and return a proper message to the client.
func HandleError(c *gin.Context, err error) {
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

// HandleBindingError parses the error and returns a proper message to the client.
func HandleBindingError(c *gin.Context, err error) {
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
			case "user_role":
				messages = append(messages, fmt.Sprintf("%s is not a valid user role", e.Field()))
			case "min":
				messages = append(messages, fmt.Sprintf("%s must be more than %s", e.Field(), e.Param()))
			case "max":
				messages = append(messages, fmt.Sprintf("%s must be less than %s", e.Field(), e.Param()))
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
		"code":     "INVALID_ENTITY",
		"messages": messages,
	})
}
