package response

import (
	"encoding/json"
	"net/http"

	apierror "github.com/Croazt/shopifyx/utils/response/error"
	apisuccess "github.com/Croazt/shopifyx/utils/response/success"
)

func GenerateResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data == nil {
		return
	}
	json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, e apierror.Error) {
	GenerateResponse(w, e.HttpStatus, e)
}

func Success(w http.ResponseWriter, e apisuccess.Success) {
	GenerateResponse(w, e.HttpStatus, e)
}
func SuccessMeta(w http.ResponseWriter, e apisuccess.SuccessWithMeta) {
	GenerateResponse(w, e.HttpStatus, e)
}
