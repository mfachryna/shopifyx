package routes

import (
	"database/sql"

	"github.com/Croazt/shopifyx/handler"
	"github.com/fasthttp/router"
	"github.com/go-playground/validator/v10"
)

func AuthRoute(r *router.Router, db *sql.DB, validator *validator.Validate) {
	authHandler := handler.NewAuthHandler(db, validator)
	r.POST("/v1/user/register", authHandler.Register)
	r.POST("/v1/user/login", authHandler.Login)
}

func ImageRoute(r *router.Router, validator *validator.Validate) {
	imageHandler := handler.NewImageHandler(validator)
	r.POST("/v1/image", imageHandler.Store)
}
