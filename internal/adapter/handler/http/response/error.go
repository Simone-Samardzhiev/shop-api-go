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

// ErrorResponse represent a meaningful response when an error occurs.
type ErrorResponse struct {
	Code       string   `json:"code" example:"INTERNAL_SERVER_ERROR"`
	Messages   []string `json:"messages" example:"Server cannot process the request."`
	statusCode int
}

var errMap = map[error]ErrorResponse{
	domain.ErrInternal: {
		Code:       "INTERNAL_SERVER_ERROR",
		Messages:   []string{"Server cannot process the request."},
		statusCode: http.StatusInternalServerError,
	}, domain.ErrEmailAlreadyInUse: {
		Code:       "EMAIL_ALREADY_IN_USE",
		Messages:   []string{"Email is already in use."},
		statusCode: http.StatusConflict,
	}, domain.ErrUsernameAlreadyInUse: {
		Code:       "USER_ALREADY_IN_USE",
		Messages:   []string{"Username is already in use."},
		statusCode: http.StatusConflict,
	}, domain.ErrWrongCredentials: {
		Code:       "WRONG_CREDENTIALS",
		Messages:   []string{"Wrong credentials."},
		statusCode: http.StatusUnauthorized,
	}, domain.ErrUserNotFound: {
		Code:       "USER_NOT_FOUND",
		Messages:   []string{"User not found."},
		statusCode: http.StatusNotFound,
	}, domain.ErrNoFieldsToUpdate: {
		Code:       "NO_FIELDS_TO_UPDATE",
		Messages:   []string{"No fields to update."},
		statusCode: http.StatusBadRequest,
	}, domain.ErrInvalidQuery: {
		Code:       "INVALID_QUERY",
		Messages:   []string{"Invalid query."},
		statusCode: http.StatusBadRequest,
	}, domain.ErrInvalidCursor: {
		Code:       "INVALID_CURSOR_FORMAT",
		Messages:   []string{"Invalid cursor format."},
		statusCode: http.StatusBadRequest,
	}, domain.ErrInvalidUUID: {
		Code:       "INVALID_UUID",
		Messages:   []string{"Invalid uuid."},
		statusCode: http.StatusBadRequest,
	}, domain.ErrLimitNotSet: {
		Code:       "LIMIT_NOT_SET",
		Messages:   []string{"Please provide a limit when fetching a list of objects."},
		statusCode: http.StatusBadRequest,
	}, domain.ErrInvalidToken: {
		Code:       "INVALID_TOKEN",
		Messages:   []string{"Token is invalid."},
		statusCode: http.StatusUnauthorized,
	}, domain.ErrInvalidTokenType: {
		Code:       "INVALID_TOKEN_TYPE",
		Messages:   []string{"Token type is invalid."},
		statusCode: http.StatusForbidden,
	}, domain.ErrInvalidTokenRole: {
		Code:       "INVALID_TOKEN_ROLE",
		Messages:   []string{"Token role is invalid."},
		statusCode: http.StatusForbidden,
	}, domain.ErrTokenNotFound: {
		Code:       "TOKEN_NOT_FOUND",
		Messages:   []string{"Token not found."},
		statusCode: http.StatusNotFound,
	},
}

// HandleError parses the error and return a proper message to the client.
func HandleError(c *gin.Context, err error) {
	res, ok := errMap[err]
	if !ok {
		res = ErrorResponse{
			Code:       "INTERNAL_SERVER_ERROR",
			Messages:   []string{"Server cannot process the request."},
			statusCode: http.StatusInternalServerError,
		}
	}
	c.JSON(res.statusCode, res)
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
		"Code":     "INVALID_ENTITY",
		"messages": messages,
	})
}
