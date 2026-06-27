package validate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func FormatValidationError(err error) any {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		if err.Error() == "EOF" {
			return []FieldError{
				{Field: "body", Message: "request body is required"},
			}
		}
		return []FieldError{
			{Field: "request", Message: err.Error()},
		}
	}

	out := make([]FieldError, 0, len(ve))
	for _, fe := range ve {
		out = append(out, FieldError{
			Field:   jsonFieldName(fe),
			Message: customMessage(fe),
		})
	}
	return out
}

func customMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", jsonFieldName(fe))
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", jsonFieldName(fe), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", jsonFieldName(fe), fe.Param())
	case "email":
		return fmt.Sprintf("%s is not a valid email", jsonFieldName(fe))
	case "phone_number":
		return fmt.Sprintf("%s is not a valid phone number", jsonFieldName(fe))
	case "url":
		return fmt.Sprintf("%s is not a valid URL", jsonFieldName(fe))
	case "uuid":
		return fmt.Sprintf("%s is not a valid UUID", jsonFieldName(fe))
	case "datetime":
		return fmt.Sprintf("%s is not a valid datetime", jsonFieldName(fe))
	case "number":
		return fmt.Sprintf("%s is not a valid number", jsonFieldName(fe))
	case "string":
		return fmt.Sprintf("%s is not a valid string", jsonFieldName(fe))
	case "boolean":
		return fmt.Sprintf("%s is not a valid boolean", jsonFieldName(fe))
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", jsonFieldName(fe), fe.Param())
	default:
		return fmt.Sprintf("%s is invalid", jsonFieldName(fe))
	}
}

func jsonFieldName(fe validator.FieldError) string {
	field := fe.Field()
	return toSnake(field)
}

func toSnake(s string) string {
	var b strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			b.WriteByte('_')
		}
		b.WriteRune(r)
	}
	return strings.ToLower(b.String())
}
