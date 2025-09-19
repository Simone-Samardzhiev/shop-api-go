package validation

import (
	"strconv"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// validateMinBytesLength is a function that implement validator.FieldLevel interface
// and varifies that length of the word is not less than the provided length. It is registered as a custom "min_bytes" validator
// tag in Gin's binding.Validator.
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
// and varifies that length of the word is not more than the provided length. It is registered as a custom "min_bytes" validator
// tag in Gin's binding.Validator.
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

// validateMinBytesLength is a function that implement validator.FieldLevel interface
// and varifies if length of the word is less than the provided param. It is registered as a custom "max_bytes" validator
// tag in Gin's binding.Validator.
func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("min_bytes", validateMinBytesLength)
		_ = v.RegisterValidation("max_bytes", validateMaxBytesLength)
	}
}
