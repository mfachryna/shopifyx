package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Croazt/shopifyx/db/connection/postgresql"
	"github.com/Croazt/shopifyx/db/migrations"
	"github.com/Croazt/shopifyx/middleware"
	"github.com/Croazt/shopifyx/routes"
	"github.com/Croazt/shopifyx/utils/validation"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var db *sql.DB

func main() {
	var (
		err            error
		migrateCommand string
		validate       *validator.Validate
	)

	flag.StringVar(&migrateCommand, "migrate", "up", "migration")
	flag.Parse()

	if godotenv.Load() != nil {
		log.Fatal("error loading .env file")
	}

	db, err = postgresql.OpenPg()
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer db.Close()

	if migrateCommand != "" {
		err = migrations.Migrate(db, migrateCommand)
		if err != nil {
			log.Fatalf("error migrating to schema: %v", err)
		}
	}

	validate = validator.New()
	if err := validation.RegisterCustomValidation(validate); err != nil {
		log.Fatalf("error register custom validation")
	}

	r := chi.NewRouter()

	r.Handle("/metrics", promhttp.Handler())
	r.Route("/v1", func(r chi.Router) {
		r.Use(middleware.PrometheusMiddleware)
		routes.AuthRoute(r, db, validate)
		routes.ImageRoute(r, validate)
		routes.ProductRoute(r, db, validate)
		routes.BankAccountRoute(r, db, validate)

	})
	s := &http.Server{
		Addr:    ":8000",
		Handler: r,
	}

	go func() {
		fmt.Println("Listen and Serve at port 8000")
		if err := s.ListenAndServe(); err != nil {
			log.Fatalf("error in ListenAndServe: %s", err)
		}
	}()
	log.Print("Server Started")

	stopped := make(chan os.Signal, 1)
	signal.Notify(stopped, os.Interrupt)
	<-stopped

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("shutting down gracefully...")
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("error in Server Shutdown: %s", err)
	}
	fmt.Println("server stopped")
}
