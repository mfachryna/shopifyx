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
		r.Get("/", productHandler.Index)
		r.Post("/", productHandler.Create)

		r.Route("/{productId}", func(r chi.Router) {
			r.Patch("/", productHandler.Update)
			r.Delete("/", productHandler.Delete)
			r.Get("/stock", productHandler.Stock)
			r.Get("/", productHandler.Show)

			paymentHandler := handler.NewPaymentHandler(db, validator)
			r.Post("/buy", paymentHandler.Create)
		})
	})
}
func BankAccountRoute(r chi.Router, db *sql.DB, validator *validator.Validate) {
	bankAccountHandler := handler.NewBankAccountHandler(db, validator)
	r.Route("/bank/account", func(r chi.Router) {
		r.Use(middleware.JwtMiddleware)
		r.Get("/", bankAccountHandler.Index)
		r.Post("/", bankAccountHandler.Create)
		r.Patch("/{bankAccountId}", bankAccountHandler.Update)
		r.Delete("/{bankAccountId}", bankAccountHandler.Delete)
	})
}
