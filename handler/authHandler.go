package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/Croazt/shopifyx/repository"
	"github.com/Croazt/shopifyx/utils/response"
	apierror "github.com/Croazt/shopifyx/utils/response/error"
	apisuccess "github.com/Croazt/shopifyx/utils/response/success"
	"github.com/valyala/fasthttp"
)

type AuthHandler struct {
	db *sql.DB
}

// NewUserHandler creates a new instance of UserHandler
func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{
		db: db,
	}
}

// Register registers a new user
func (uh *AuthHandler) Register(ctx *fasthttp.RequestCtx) {
	var user repository.UserRegister
	if err := json.Unmarshal(ctx.PostBody(), &user); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	rows := uh.db.QueryRow(`SELECT name FROM users WHERE username = $1;`, user.Username)
	var exists bool

	if err := rows.Scan(&exists); err == nil {
		response.Error(ctx, apierror.ClientAlreadyExists())
		return
	} else if !exists {
		date := time.Now()
		_, err := uh.db.Exec(`INSERT INTO users (username,name,password,created_at,updated_at) VALUES ($1,$2,$3,$4,$5)`, user.Name, user.Username, user.Password, date, date)
		if err != nil {
			log.Printf("%v", err)
			response.Error(ctx, apierror.ServerError())
			return
		}

		response.Success(ctx, apisuccess.RegisterResponse(&user))
		return
	}

}
