package bcrypt

import (
	"shop-api-go/internal/core/domain"

	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher implements port.PasswordHasher and provides password hashing with bcrypt.
type PasswordHasher struct{}

func (p *PasswordHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", domain.ErrInternalServerError
	}
	return string(hash), nil
}

func (p *PasswordHasher) Compare(password, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return domain.ErrWrongCredentials
	}
	return nil
}
