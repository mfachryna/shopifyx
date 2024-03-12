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
	"github.com/google/uuid"
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
	var registerData domain.UserRegister
	if err := json.Unmarshal(ctx.PostBody(), &registerData); err != nil {

		response.Error(ctx, apierror.ClientBadRequest())
		return
	}

	var count int
	err := uh.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", registerData.Username).Scan(&count)

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
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerData.Password), bcrypt.DefaultCost)
		if err != nil {
			response.Error(ctx, apierror.CustomServerError(err.Error()))
			return
		}
		hashedPasswordChan <- string(hashedPassword)
	}()

	var id string
	uuid := uuid.New()
	if err := uh.db.QueryRow(`INSERT INTO users (id,username,name,password,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`, uuid, registerData.Username, registerData.Name, <-hashedPasswordChan, date, date).Scan(&id); err != nil {
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

	res := &domain.UserAuthResponse{
		Name:        registerData.Name,
		Username:    registerData.Username,
		AccessToken: tokenString,
	}

	response.Success(ctx, apisuccess.RegisterResponse(res))
	return

}

// Register registers a new user
func (uh *AuthHandler) Login(ctx *fasthttp.RequestCtx) {
	var loginData domain.UserLogin
	if err := json.Unmarshal(ctx.PostBody(), &loginData); err != nil {
		response.Error(ctx, apierror.ClientBadRequest())
		return
	}

	var user domain.User
	err := uh.db.QueryRow("SELECT id,username,name,password FROM users WHERE username = $1 LIMIT 1;", loginData.Username).Scan(&user.ID, &user.Username, &user.Name, &user.Password)

	if err != nil {
		response.Error(ctx, apierror.ClientNotFound("Username"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		response.Error(ctx, apierror.CustomError(400, "Password missmatched"))
	}

	tokenString, err := jwt.SignedToken(jwt.Claim{
		UserId: user.ID,
	})
	if err != nil {
		response.Error(ctx, apierror.CustomServerError("Failed to generate access token"))
		return
	}

	res := &domain.UserAuthResponse{
		Name:        user.Name,
		Username:    user.Username,
		AccessToken: tokenString,
	}

	response.Success(ctx, apisuccess.LoginResponse(res))
	return

}
