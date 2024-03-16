package validation

import (
	"fmt"
	"mime/multipart"
	"reflect"
	"regexp"
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
		return fmt.Errorf("failed to register boolean validation: %s", err)
	}
	if err := v.RegisterValidation("noSpace", validateNoSpace); err != nil {
		return fmt.Errorf("failed to register username has space: %s", err)
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

	// fmt.Println(fileHeader)
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
	fmt.Println(ext)
	return ext == "jpg" || ext == "jpeg"
}

func imageMaxSizeValidator(fl validator.FieldLevel) bool {
	file, ok := fl.Top().Interface().(*multipart.FileHeader)
	if !ok {
		return false
	}
	maxSize := int64(2) // 2MB
	fmt.Println(file.Size, maxSize)
	return file.Size <= maxSize
}

func validateIsBool(fl validator.FieldLevel) bool {
	return fl.Field().Kind() == reflect.Bool
}

func validateNoSpace(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	return !strings.Contains(field, " ")
}

func UrlValidation(url string) error {
	pattern := `^(https?|ftp):\/\/[^\s\/$.?#].[^\s]*$`

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("failed to compile pattern %v", err)
	}

	if !regex.MatchString(url) {
		return fmt.Errorf("url is not valid")
	}

	return nil
}

func UuidValidation(uuid string) error {
	pattern := `^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[1-5][a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("failed to compile pattern %v", err)
	}

	if !regex.MatchString(uuid) {
		return fmt.Errorf("uuid is not valid")
	}

	return nil
}
