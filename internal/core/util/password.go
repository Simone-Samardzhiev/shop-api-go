package util

import (
	"shop-api-go/internal/core/domain"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes input password using bcrypt.
func HashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, domain.ErrInternalServerError
	}
	return hash, nil
}

// ComparePassword compares input password and hash.
func ComparePassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
