package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Croazt/shopifyx/domain"
	"github.com/Croazt/shopifyx/utils/jwt"
	"github.com/Croazt/shopifyx/utils/response"
	apierror "github.com/Croazt/shopifyx/utils/response/error"
	apisuccess "github.com/Croazt/shopifyx/utils/response/success"
	"github.com/Croazt/shopifyx/utils/validation"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db        *sql.DB
	validator *validator.Validate
}

// NewUserHandler creates a new instance of UserHandler
func NewAuthHandler(db *sql.DB, validator *validator.Validate) *AuthHandler {
	return &AuthHandler{
		db:        db,
		validator: validator,
	}
}

// Register registers a new user
func (uh *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var registerData domain.UserRegister
	if err := json.NewDecoder(r.Body).Decode(&registerData); err != nil {

		response.Error(w, apierror.ClientBadRequest())
		return
	}

	if err := uh.validator.Struct(registerData); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			response.Error(w, apierror.CustomError(http.StatusBadRequest, validation.CustomError(e)))
			return
		}
	}

	registerData.Username = strings.ToLower(registerData.Username)

	var count int
	err := uh.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", registerData.Username).Scan(&count)

	if err != nil {
		response.Error(w, apierror.CustomServerError(err.Error()))
		return
	}

	if count > 0 {
		response.Error(w, apierror.ClientAlreadyExists())
		return
	}

	date := time.Now()

	hashedPasswordChan := make(chan string)
	go func() {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerData.Password), bcrypt.DefaultCost)
		if err != nil {
			response.Error(w, apierror.CustomServerError(err.Error()))
			return
		}
		hashedPasswordChan <- string(hashedPassword)
	}()

	var id string
	uuid := uuid.New()
	if err := uh.db.QueryRow(`INSERT INTO users (id,username,name,password,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`, uuid, registerData.Username, registerData.Name, <-hashedPasswordChan, date, date).Scan(&id); err != nil {
		response.Error(w, apierror.CustomServerError(err.Error()))
		return
	}

	tokenString, err := jwt.SignedToken(jwt.Claim{
		UserId: id,
	})
	if err != nil {
		response.Error(w, apierror.CustomServerError("Failed to generate access token"))
		return
	}

	res := &domain.UserAuthResponse{
		Name:        registerData.Name,
		Username:    registerData.Username,
		AccessToken: tokenString,
	}

	response.Success(w, apisuccess.RegisterResponse(res))
}

// Register registers a new user
func (uh *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginData domain.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		response.Error(w, apierror.ClientBadRequest())
		return
	}

	if err := uh.validator.Struct(loginData); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			response.Error(w, apierror.CustomError(http.StatusBadRequest, validation.CustomError(e)))
			return
		}
	}

	loginData.Username = strings.ToLower(loginData.Username)

	var user domain.User
	err := uh.db.QueryRow("SELECT id,username,name,password FROM users WHERE username = $1 LIMIT 1;", loginData.Username).Scan(&user.ID, &user.Username, &user.Name, &user.Password)

	if err != nil {
		response.Error(w, apierror.ClientNotFound("Username"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		response.Error(w, apierror.CustomError(400, "Password missmatched"))
		return
	}

	tokenString, err := jwt.SignedToken(jwt.Claim{
		UserId: user.ID,
	})
	if err != nil {
		response.Error(w, apierror.CustomServerError("Failed to generate access token"))
		return
	}

	res := &domain.UserAuthResponse{
		Name:        user.Name,
		Username:    user.Username,
		AccessToken: tokenString,
	}

	response.Success(w, apisuccess.LoginResponse(res))
}
