package domain

import (
	"errors"
)

var (
	// ErrInternal represents a generic internal service failure.
	ErrInternal = errors.New("internal error")

	// ErrEmailAlreadyInUse indicates a user's email conflicts with another.
	ErrEmailAlreadyInUse = errors.New("email already in use")

	// ErrUsernameAlreadyInUse indicates a user's username conflicts with another.
	ErrUsernameAlreadyInUse = errors.New("username already in use")

	// ErrWrongCredentials indicates the provided credentials are incorrect.
	ErrWrongCredentials = errors.New("wrong credentials")

	// ErrUserNotFound indicates the requested user could not be found.
	ErrUserNotFound = errors.New("user not found")

	// ErrNoFieldsToUpdate indicates that no fields were provided for an update.
	ErrNoFieldsToUpdate = errors.New("no fields to update")

	// ErrInvalidQuery indicates that the provided query parameters are invalid.
	ErrInvalidQuery = errors.New("invalid query")

	// ErrLimitNotSet indicates that a required limit parameter was missing.
	ErrLimitNotSet = errors.New("limit not set")

	// ErrInvalidCursor indicates that the provided cursor format is invalid.
	ErrInvalidCursor = errors.New("invalid cursor format")

	// ErrInvalidUUID indicates that a provided UUID is not valid.
	ErrInvalidUUID = errors.New("invalid uuid")

	// ErrInvalidToken indicates a failure to parse or validate a token.
	ErrInvalidToken = errors.New("invalid token")

	// ErrInvalidTokenType indicates the token type is not allowed for this action.
	ErrInvalidTokenType = errors.New("invalid token type")

	// ErrInvalidTokenRole indicates that the token's role does not have permission.
	ErrInvalidTokenRole = errors.New("invalid token role")

	// ErrTokenNotFound indicates that the expected token was not found.
	ErrTokenNotFound = errors.New("token not found")
)
