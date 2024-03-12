package handler

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Croazt/shopifyx/domain"
	"github.com/Croazt/shopifyx/utils/jwt"
	"github.com/Croazt/shopifyx/utils/response"
	apierror "github.com/Croazt/shopifyx/utils/response/error"
	apisuccess "github.com/Croazt/shopifyx/utils/response/success"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/bcrypt"
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
	var userData domain.UserRegister
	if err := json.Unmarshal(ctx.PostBody(), &userData); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	var count int
	err := uh.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", userData.Username).Scan(&count)

	if err != nil {
		response.Error(ctx, apierror.CustomServerError(err.Error()))
		return
	}

	if count > 0 {
		response.Error(ctx, apierror.ClientAlreadyExists())
		return
	}

	date := time.Now()

	hashedPasswordChan := make(chan string)
	go func() {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
		if err != nil {
			response.Error(ctx, apierror.CustomServerError(err.Error()))
			return
		}
		hashedPasswordChan <- string(hashedPassword)
	}()

	var id int
	if err := uh.db.QueryRow(`INSERT INTO users (username,name,password,created_at,updated_at) VALUES ($1,$2,$3,$4,$5) RETURNING id`, userData.Username, userData.Name, <-hashedPasswordChan, date, date).Scan(&id); err != nil {
		response.Error(ctx, apierror.CustomServerError(err.Error()))
		return
	}

	tokenString, err := jwt.SignedToken(jwt.Claim{
		UserId: id,
	})
	if err != nil {
		response.Error(ctx, apierror.CustomServerError("Failed to generate access token"))
		return
	}

	res := &domain.UserRegisterResponse{
		Name:        userData.Name,
		Username:    userData.Username,
		AccessToken: tokenString,
	}

	response.Success(ctx, apisuccess.RegisterResponse(res))
	return

}
