package validation

import (
	"fmt"
	"mime/multipart"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidation(v *validator.Validate) error {
	if err := v.RegisterValidation("fileformat", fileFormatValidator); err != nil {
		return fmt.Errorf("failed to register file format validation: %s", err)
	}

	if err := v.RegisterValidation("imageMaxSize", imageMaxSizeValidator); err != nil {
		return fmt.Errorf("failed to register image max size validation: %s", err)
	}

	if err := v.RegisterValidation("isBool", validateIsBool); err != nil {
		return fmt.Errorf("failed to register is boolean validation: %s", err)
	}

	return nil
}

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

func ValidateFile(v *validator.Validate, fileHeader *multipart.FileHeader) error {
	if err := v.Var(fileHeader, "fileformat"); err != nil {
		return fmt.Errorf("file format must be JPG or JPEG")
	}
	if err := v.Var(fileHeader, "imageMaxSize"); err != nil {
		return fmt.Errorf("image size cannot exceed 2MB")
	}
	return nil
}

func fileFormatValidator(fl validator.FieldLevel) bool {
	file, ok := fl.Top().Interface().(*multipart.FileHeader)
	if !ok {
		return false
	}
	ext := strings.ToLower(file.Filename[strings.LastIndex(file.Filename, ".")+1:])
	return ext == "jpg" || ext == "jpeg"
}

func imageMaxSizeValidator(fl validator.FieldLevel) bool {
	file, ok := fl.Top().Interface().(*multipart.FileHeader)
	if !ok {
		return false
	}
	maxSize := int64(2 * 1024 * 1024) // 2MB
	return file.Size <= maxSize
}

func validateIsBool(fl validator.FieldLevel) bool {
	return fl.Field().Kind() == reflect.Bool
}
