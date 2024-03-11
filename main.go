package main

import (
	"log"

	"github.com/Croazt/shopifyx/connection"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := connection.SetupPg()
}
