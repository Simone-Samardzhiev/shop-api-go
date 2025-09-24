package domain

import (
	"errors"
	"net/http"
)

// Error represent a domain-level error.
type Error struct {
	Code       string
	StatusCode int
	Err        error
}

// NewError creates a new Error instance.
func NewError(code string, err error, statusCode int) *Error {
	return &Error{Code: code, Err: err, StatusCode: statusCode}
}

func (e *Error) Error() string {
	return e.Err.Error()
}

var (
	// ErrInternalServerError is an error for when the service fails to process the request.
	ErrInternalServerError = NewError("INTERNAL_SERVER_ERROR", errors.New("internal server error"), http.StatusInternalServerError)
	// ErrEmailAlreadyInUse is an error for when user's email conflicts with another.
	ErrEmailAlreadyInUse = NewError("EMAIL_ALREADY_EXIST", errors.New("email already in use"), http.StatusConflict)
	// ErrUsernameAlreadyInUse is an error for when user's username conflicts with another.
	ErrUsernameAlreadyInUse = NewError("USERNAME_ALREADY_EXIST", errors.New("username already exist"), http.StatusConflict)
	// ErrInvalidToken is an error for when token parsing fails.
	ErrInvalidToken = NewError("INVALID_TOKEN", errors.New("invalid token"), http.StatusBadRequest)
	// ErrInvalidTokenType is an error for when user use invalid token type.
	ErrInvalidTokenType = NewError("INVALID_TOKEN_TYPE", errors.New("invalid token type"), http.StatusBadRequest)
	// ErrMalformedToken is an error for when token data is invalid.
	ErrMalformedToken = NewError("INVALID_TOKEN_SUBJECT", errors.New("invalid token subject"), http.StatusBadRequest)
)
