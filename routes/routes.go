package routes

import (
	"database/sql"

	"github.com/Croazt/shopifyx/handler"
	"github.com/Croazt/shopifyx/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func AuthRoute(r chi.Router, db *sql.DB, validator *validator.Validate) {
	authHandler := handler.NewAuthHandler(db, validator)
	r.Route("/user", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})
}

func ImageRoute(r chi.Router, validator *validator.Validate) {
	imageHandler := handler.NewImageHandler(validator)
	r.Route("/image", func(r chi.Router) {
		r.Use(middleware.JwtMiddleware)
		r.Post("/", imageHandler.Store)
	})
}

func ProductRoute(r chi.Router, db *sql.DB, validator *validator.Validate) {
	productHandler := handler.NewProductHandler(db, validator)
	r.Route("/product", func(r chi.Router) {
		r.Use(middleware.JwtMiddleware)
		r.Post("/", productHandler.Create)
		r.Patch("/{productId}", productHandler.Update)
	})
}
