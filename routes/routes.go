package routes

import (
	"database/sql"

	"github.com/Croazt/shopifyx/handler"
	"github.com/fasthttp/router"
)

func AuthRoute(r *router.Router, db *sql.DB) {
	authHandler := handler.NewAuthHandler(db)
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
}
