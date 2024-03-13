package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func CustomError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "min", "max":
		return fmt.Sprintf("%s too short or long", e.Field())
	default:
		return e.Error()
	}
}
