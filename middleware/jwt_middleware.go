package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Croazt/shopifyx/utils/response"
	apierror "github.com/Croazt/shopifyx/utils/response/error"
	"github.com/golang-jwt/jwt"
)

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			fmt.Println("token not found")
			response.Error(w, apierror.CustomError(http.StatusUnauthorized, "token not found"))
			return
		}

		tokenString := string(authHeader)
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				fmt.Printf("unexpected signing method: %v \n", t.Header["alg"])
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			validationErr, ok := err.(*jwt.ValidationError)
			if ok {
				if validationErr.Errors == jwt.ValidationErrorExpired {
					fmt.Println(err.Error())
					response.Error(w, apierror.ClientAccessExpired())
					return
				}
			}
			fmt.Println(err.Error())
			response.Error(w, apierror.ClientInvalidToken())
			return
		}

		if !token.Valid {
			fmt.Println("invalid token claims")
			response.Error(w, apierror.CustomError(http.StatusUnauthorized, "invalid token claims"))
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", token.Claims.(jwt.MapClaims)["user_id"])
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
func OptionalJwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
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
					next.ServeHTTP(w, r)
					return
				}
			}
			next.ServeHTTP(w, r)
			return
		}

		if !token.Valid {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", token.Claims.(jwt.MapClaims)["user_id"])
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
