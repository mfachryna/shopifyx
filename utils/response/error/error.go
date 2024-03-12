package error

import "net/http"

type Error struct {
	HttpStatus int    `json:"-"`
	Message    string `json:"message"`
}

func CustomError(status int, message string) Error {
	return Error{
		HttpStatus: status,
		Message:    message,
	}
}
func ClientBadRequest() Error {
	return Error{
		HttpStatus: http.StatusBadRequest,
		Message:    "failed to parse request",
	}
}

func ClientNotFound(resourceName string) Error {
	return Error{
		HttpStatus: http.StatusNotFound,
		Message:    resourceName + " not found",
	}
}

func ClientUnauthorized() Error {
	return Error{
		HttpStatus: http.StatusUnauthorized,
		Message:    "given security scheme is invalid",
	}
}

func ClientInvalidCredential() Error {
	return Error{
		HttpStatus: http.StatusUnauthorized,
		Message:    "email or password is incorect",
	}
}

func ClientAccessExpired() Error {
	return Error{
		HttpStatus: http.StatusUnauthorized,
		Message:    "given security scheme is valid, but the lifetime has been expired or revoked.",
	}
}

func ClientForbidden() Error {
	return Error{
		HttpStatus: http.StatusForbidden,
		Message:    "we already sent email, please wait for a minute",
	}
}

func ClientInvalidToken() Error {
	return Error{
		HttpStatus: http.StatusUnauthorized,
		Message:    "token is invalid",
	}
}

func ClientInactiveUser() Error {
	return Error{
		HttpStatus: http.StatusBadRequest,
		Message:    "request account is inactive",
	}
}

func ClientAlreadyExists() Error {
	return Error{
		HttpStatus: http.StatusConflict,
		Message:    "username is already exists",
	}
}

func ServerError() Error {
	return Error{
		HttpStatus: http.StatusInternalServerError,
		Message:    "server has internal error",
	}
}
func CustomServerError(message string) Error {
	return Error{
		HttpStatus: http.StatusInternalServerError,
		Message:    message,
	}
}
