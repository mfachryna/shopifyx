package error

import (
	"net/http"

	"github.com/Croazt/shopifyx/repository"
)

type Success struct {
	HttpStatus int         `json:"-"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func RegisterResponse(user *repository.UserRegister) Success {
	return Success{
		HttpStatus: http.StatusCreated,
		Message:    "User registered successfully ",
		Data:       user,
	}
}
