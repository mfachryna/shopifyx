package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Croazt/shopifyx/db/connection"
	"github.com/Croazt/shopifyx/db/migrations"
	"github.com/Croazt/shopifyx/routes"
	"github.com/fasthttp/router"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
)

var db *sql.DB

func main() {
	var (
		err            error
		migrateCommand string
	)

	flag.StringVar(&migrateCommand, "migrate", "up", "migration")
	flag.Parse()

	if godotenv.Load() != nil {
		log.Fatal("error loading .env file")
	}

	db, err = connection.OpenPg()
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
			log.Fatalf("error in ListenAndServe: %s", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	fmt.Println("shutting down gracefully...")
	if err := s.Shutdown(); err != nil {
		log.Fatalf("error in Server Shutdown: %s", err)
	}
	fmt.Println("server stopped")
}
