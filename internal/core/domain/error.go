package domain

import (
	"errors"
)

var (
	// ErrInternalServerError is an error for when the service fails to process the request.
	ErrInternalServerError = errors.New("internal server error")
	// ErrEmailAlreadyInUse is an error for when user's email conflicts with another.
	ErrEmailAlreadyInUse = errors.New("email already in use")
	// ErrUsernameAlreadyInUse is an error for when user's username conflicts with another.
	ErrUsernameAlreadyInUse = errors.New("username already exist")
	// ErrWrongCredentials is an error for when user's credentials are wrong.
	ErrWrongCredentials = errors.New("wrong credentials")
	// ErrUserNotFound is an error for when user couldn't be found.
	ErrUserNotFound = errors.New("user not found")
	// ErrNoFieldsToUpdate is an error for when all update fields are empty(nil)
	ErrNoFieldsToUpdate = errors.New("no fields to update")
	// ErrInvalidQuery is an error for when passed query for fetching is invalid.
	ErrInvalidQuery = errors.New("invalid query")
	// ErrLimitNotSet is an error for when fetching data requires limit to be set.
	ErrLimitNotSet = errors.New("limit not set")
	// ErrInvalidCursorFormat is an error for when cursor format is unexpected.
	ErrInvalidCursorFormat = errors.New("invalid cursor format")
	// ErrInvalidToken is an error for when token parsing fails.
	ErrInvalidToken = errors.New("invalid token")
	// ErrInvalidTokenType is an error for when user use invalid token type.
	ErrInvalidTokenType = errors.New("invalid token type")
	// ErrInvalidTokenRole is an error for when user tries to access a resource,
	// to which he doesn't have a precision to access.
	ErrInvalidTokenRole = errors.New("invalid token role")
	// ErrMalformedToken is an error for when token data is invalid.
	ErrMalformedToken = errors.New("invalid token subject")
	// ErrTokenNotFound  is an error for when token is expired or already user.
	ErrTokenNotFound = errors.New("token expired or already user")
)
