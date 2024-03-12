package response

import (
	"encoding/json"

	apierror "github.com/Croazt/shopifyx/utils/response/error"
	apisuccess "github.com/Croazt/shopifyx/utils/response/success"
	"github.com/valyala/fasthttp"
)

func GenerateResponse(ctx *fasthttp.RequestCtx, status int, data interface{}) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(status)
	if data == nil {
		return
	}
	json.NewEncoder(ctx).Encode(data)
}

func Error(ctx *fasthttp.RequestCtx, e apierror.Error) {
	GenerateResponse(ctx, e.HttpStatus, e)
}

func Success(ctx *fasthttp.RequestCtx, e apisuccess.Success) {
	GenerateResponse(ctx, e.HttpStatus, e)
}
