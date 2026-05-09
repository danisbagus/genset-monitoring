package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps the go-playground validator.
type Validator struct {
	validate *validator.Validate
}

// New creates a new Validator instance.
func New() *Validator {
	return &Validator{validate: validator.New()}
}

// Validate validates a struct and returns a map of field errors.
func (v *Validator) Validate(s interface{}) map[string]string {
	err := v.validate.Struct(s)
	if err == nil {
		return nil
	}

	errs := make(map[string]string)
	var ve validator.ValidationErrors
	if ok := isValidationErrors(err, &ve); ok {
		for _, fe := range ve {
			field := strings.ToLower(fe.Field())
			errs[field] = formatMessage(fe)
		}
	}
	return errs
}

func isValidationErrors(err error, ve *validator.ValidationErrors) bool {
	if e, ok := err.(validator.ValidationErrors); ok {
		*ve = e
		return true
	}
	return false
}

func formatMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", fe.Field(), fe.Param())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", fe.Field())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("%s failed validation: %s", fe.Field(), fe.Tag())
	}
}
