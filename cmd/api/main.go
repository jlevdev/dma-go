package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const port = 3000

type application struct {
	Domain         string
	DB             *sql.DB
	JWTSecret      string
	GoogleClientID string
}

func main() {
	var app application

	app.Domain = "localhost"

	app.JWTSecret = os.Getenv("JWT_SECRET")
	if app.JWTSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	app.GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	if app.GoogleClientID == "" {
		log.Fatal("GOOGLE_CLIENT_ID environment variable is required")
	}

	db, err := sql.Open("pgx", "postgres://dma:dma@localhost:5432/dma")
	if err != nil {
		log.Fatal("failed to open db:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("failed to connect to db:", err)
	}

	app.DB = db

	log.Println("Starting application on port", port, " for domain ", app.Domain)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
