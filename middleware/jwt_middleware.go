package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Croazt/shopifyx/utils/response"
	apierror "github.com/Croazt/shopifyx/utils/response/error"
	"github.com/golang-jwt/jwt"
	"github.com/valyala/fasthttp"
)

func JwtMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		authHeader := ctx.Request.Header.Peek("Authorization")
		if authHeader == nil {
			response.Error(ctx, apierror.CustomError(http.StatusUnauthorized, "token not found"))
			return
		}

		tokenString := string(authHeader)
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			validationErr, ok := err.(*jwt.ValidationError)
			if ok {
				if validationErr.Errors == jwt.ValidationErrorExpired {
					response.Error(ctx, apierror.ClientAccessExpired())
					return
				}
			}
			response.Error(ctx, apierror.CustomServerError("invalid token"))
			return
		}

		if !token.Valid {
			response.Error(ctx, apierror.CustomError(http.StatusUnauthorized, "invalid token claims"))
			return
		}

		ctx.SetUserValue("user_id", token.Claims.(jwt.MapClaims)["user_id"])
		next(ctx)
	}
}
