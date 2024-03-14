package error

import (
	"net/http"

	"github.com/Croazt/shopifyx/domain"
)

type Success struct {
	HttpStatus int         `json:"-"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}
type SuccessWithMeta struct {
	Success
	Meta domain.Meta `json:"meta"`
}

func RegisterResponse(user *domain.UserAuthResponse) Success {
	return Success{
		HttpStatus: http.StatusCreated,
		Message:    "User registered successfully",
		Data:       user,
	}
}

func LoginResponse(user *domain.UserAuthResponse) Success {
	return Success{
		HttpStatus: http.StatusOK,
		Message:    "User logged successfully",
		Data:       user,
	}
}

func CustomResponse(status int, message string, data interface{}) Success {
	return Success{
		HttpStatus: status,
		Message:    message,
		Data:       data,
	}
}

func IndexResponse(status int, message string, data interface{}, meta domain.Meta) SuccessWithMeta {
	return SuccessWithMeta{
		Success: Success{
			HttpStatus: status,
			Message:    message,
			Data:       data,
		},
		Meta: meta,
	}
}
