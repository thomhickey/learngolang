package validator

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Get returns the validator instance
func Get() *validator.Validate {
	return validate
}

// Validate validates a struct and returns validation errors
func Validate(s interface{}) error {
	return validate.Struct(s)
}
