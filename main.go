package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Croazt/shopifyx/db/connection"
	"github.com/Croazt/shopifyx/routes"
	"github.com/fasthttp/router"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
)

var db *sql.DB

func main() {
	var err error
	if godotenv.Load() != nil {
		log.Fatal("Error loading .env file")
	}

	db, err = connection.OpenPg()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	r := router.New()

	s := &fasthttp.Server{
		Handler:          r.Handler,
		DisableKeepalive: true,
		ReadTimeout:      5 * time.Second,
		WriteTimeout:     5 * time.Second,
		IdleTimeout:      10 * time.Second,
	}

	routes.AuthRoute(r, db)

	go func() {
		if err := s.ListenAndServe(":8000"); err != nil {
			log.Fatalf("Error in ListenAndServe: %s", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	fmt.Println("Shutting down gracefully...")
	if err := s.Shutdown(); err != nil {
		log.Fatalf("Error in Server Shutdown: %s", err)
	}
	fmt.Println("Server stopped")
}
