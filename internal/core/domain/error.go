package domain

import "errors"

var (
	// ErrInternalServerError is an error for when the service fails to process the request.
	ErrInternalServerError = errors.New("internal server error")
	// ErrEmailAlreadyInUse is an error for when user's email conflicts with another.
	ErrEmailAlreadyInUse = errors.New("email already in use")
	// ErrUsernameAlreadyInUse is an error for when user's username conflicts with another.
	ErrUsernameAlreadyInUse = errors.New("username already in use")
)
