package dto

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) Validate(s interface{}) []FieldError {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var errors []FieldError
	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, FieldError{
			Field:   err.Field(),
			Message: v.getTagMessage(err),
		})
	}
	return errors
}

func (v *Validator) getTagMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "field is required"
	case "max":
		return "field exceeds maximum length"
	case "min":
		return "field is below minimum"
	default:
		return "invalid value"
	}
}

func (v *Validator) ParseAndValidate(data []byte, target interface{}) error {
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	if errors := v.Validate(target); len(errors) > 0 {
		return fmt.Errorf("validation failed: %v", errors)
	}
	return nil
}
