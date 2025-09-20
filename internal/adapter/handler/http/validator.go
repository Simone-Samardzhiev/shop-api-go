package http

import (
	"strconv"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// validateMinBytesLength is a function that implement validator.FieldLevel interface
// and varifies that length of the word is not less than the provided length.
func validateMinBytesLength(fl validator.FieldLevel) bool {
	word := fl.Field().String()
	params := fl.Param()

	count, err := strconv.Atoi(params)
	if err != nil {
		return false
	}
	if len(word) < count {
		return false
	}
	return true
}

// validateMaxBytesLength is a function that implement validator.FieldLevel interface
// and varifies that length of the word is not more than the provided length.
func validateMaxBytesLength(fl validator.FieldLevel) bool {
	word := fl.Field().String()
	params := fl.Param()

	count, err := strconv.Atoi(params)
	if err != nil {
		return false
	}
	if len(word) > count {
		return false
	}
	return true
}

// isValidPassword checks whether the given password meets the following criteria:
//
//  1. Length between 8 and 72 characters
//  2. Contains at least one lowercase letter
//  3. Contains at least one uppercase letter
//  4. Contains at least one number
//  5. Contains at least one punctuation/special character
//  6. Does not contain any whitespace
func isValidPassword(password string) bool {
	if len(password) < 8 || len(password) > 72 {
		return false
	}
	var (
		lower   = false
		upper   = false
		number  = false
		special = false
	)

	for _, char := range password {
		if unicode.IsSpace(char) {
			return false
		}

		switch {
		case unicode.IsLower(char):
			lower = true
		case unicode.IsUpper(char):
			upper = true
		case unicode.IsNumber(char):
			number = true
		case unicode.IsPunct(char):
			special = true
		}
	}

	return lower && upper && number && special
}

// validatePassword is a wrapper for isValidPassword that implements the
// validator.FieldLevel interface.
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return isValidPassword(password)
}
