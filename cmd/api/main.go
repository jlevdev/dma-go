package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const port = 3000

type application struct {
	Domain string
	DB     *sql.DB
}

func main() {
	var app application

	app.Domain = "localhost"

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
